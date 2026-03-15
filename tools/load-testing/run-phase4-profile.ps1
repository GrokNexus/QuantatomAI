[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("A", "B", "C", "D")]
    [string]$Profile,

    [string]$Command,
    [string]$ResultsRoot = "tools/load-testing/results",
    [int]$TenantCount = 4,
    [int]$DurationSeconds = 300,
    [switch]$DryRun
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Get-GitHash {
    $hash = & git rev-parse HEAD 2>$null
    if ($LASTEXITCODE -ne 0) {
        return "unknown"
    }

    return ($hash | Out-String).Trim()
}

function New-EvidenceSummary {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,
        [Parameter(Mandatory = $true)]
        [pscustomobject]$SelectedProfile,
        [Parameter(Mandatory = $true)]
        [string]$GitHash,
        [Parameter(Mandatory = $true)]
        [string]$RunId,
        [Parameter(Mandatory = $true)]
        [string]$ExecutedCommand,
        [Parameter(Mandatory = $true)]
        [bool]$WasDryRun,
        [Parameter(Mandatory = $true)]
        [int]$TenantCount,
        [Parameter(Mandatory = $true)]
        [int]$DurationSeconds
    )

    $thresholdLines = foreach ($property in $SelectedProfile.successThresholds.PSObject.Properties) {
        "- {0}: {1}" -f $property.Name, $property.Value
    }

    $controlLines = foreach ($control in $SelectedProfile.governanceControls) {
        "- $control"
    }

    $content = @(
        "# Phase 4 Evidence Summary",
        "",
        "## Run Metadata",
        "- Run ID: $RunId",
        "- Profile: $($SelectedProfile.id) - $($SelectedProfile.name)",
        "- Git Hash: $GitHash",
        "- Tenant Count: $TenantCount",
        "- Duration Seconds: $DurationSeconds",
        "- Dry Run: $WasDryRun",
        "- Command: $ExecutedCommand",
        "",
        "## Governance Controls",
        $controlLines,
        "",
        "## Threshold Targets",
        $thresholdLines,
        "",
        "## Observed Metrics",
        "- p50 latency ms:",
        "- p95 latency ms:",
        "- p99 latency ms:",
        "- throughput ops/sec:",
        "- tenant fairness ratio:",
        "- audit amplification:",
        "- replay invalid rows:",
        "",
        "## Result",
        "- Status:",
        "- Risks:",
        "- Notes:"
    ) -join [Environment]::NewLine

    Set-Content -Path $FilePath -Value $content -Encoding ascii
}

$workspaceRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$profilesPath = Join-Path $PSScriptRoot "phase4-profiles.json"
$profilesDocument = Get-Content -Path $profilesPath -Raw | ConvertFrom-Json
$selectedProfile = $profilesDocument.profiles | Where-Object { $_.id -eq $Profile } | Select-Object -First 1

if (-not $selectedProfile) {
    throw "Profile '$Profile' was not found in $profilesPath"
}

$gitHash = Get-GitHash
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$runId = "phase4-{0}-{1}" -f $Profile.ToLowerInvariant(), $timestamp
$resultsDirectory = Join-Path $workspaceRoot $ResultsRoot
$runDirectory = Join-Path $resultsDirectory $runId

New-Item -ItemType Directory -Path $runDirectory -Force | Out-Null

$executedCommand = if ([string]::IsNullOrWhiteSpace($Command)) { "not-provided" } else { $Command }
$manifest = [ordered]@{
    runId = $runId
    createdAtUtc = (Get-Date).ToUniversalTime().ToString("o")
    profile = $selectedProfile
    gitHash = $gitHash
    tenantCount = $TenantCount
    durationSeconds = $DurationSeconds
    dryRun = [bool]$DryRun
    command = $executedCommand
    commandExitCode = $null
    commandOutputFile = if ([string]::IsNullOrWhiteSpace($Command)) { $null } else { "command-output.txt" }
}

if (-not $DryRun -and -not [string]::IsNullOrWhiteSpace($Command)) {
    Push-Location $workspaceRoot
    try {
        $commandOutput = Invoke-Expression $Command 2>&1 | Out-String
        $manifest.commandExitCode = $LASTEXITCODE
        Set-Content -Path (Join-Path $runDirectory "command-output.txt") -Value $commandOutput -Encoding ascii
    }
    finally {
        Pop-Location
    }
}

$manifestPath = Join-Path $runDirectory "run-manifest.json"
$summaryPath = Join-Path $runDirectory "evidence-summary.md"
$manifest | ConvertTo-Json -Depth 6 | Set-Content -Path $manifestPath -Encoding ascii

New-EvidenceSummary \
    -FilePath $summaryPath \
    -SelectedProfile $selectedProfile \
    -GitHash $gitHash \
    -RunId $runId \
    -ExecutedCommand $executedCommand \
    -WasDryRun ([bool]$DryRun) \
    -TenantCount $TenantCount \
    -DurationSeconds $DurationSeconds

Write-Host "Created Phase 4 evidence bundle: $runDirectory"
Write-Host "Manifest: $manifestPath"
Write-Host "Summary : $summaryPath"

if ($manifest.commandExitCode -ne $null) {
    Write-Host "Command exit code: $($manifest.commandExitCode)"
}