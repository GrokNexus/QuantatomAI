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

function Get-FirstNumericMatch {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Text,
        [Parameter(Mandatory = $true)]
        [string]$Pattern,
        [switch]$Last
    )

    $matches = [System.Text.RegularExpressions.Regex]::Matches($Text, $Pattern, [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
    if ($matches.Count -eq 0) {
        return $null
    }

    $index = if ($Last) { $matches.Count - 1 } else { 0 }
    return [double]$matches[$index].Groups[1].Value
}

function Get-AutoExtractedMetrics {
    param(
        [Parameter(Mandatory = $true)]
        [string]$CommandOutput
    )

    $metrics = [ordered]@{}

    $nsPerOp = Get-FirstNumericMatch -Text $CommandOutput -Pattern 'Benchmark\S+\s+\d+\s+([0-9]+(?:\.[0-9]+)?)\s+ns/op' -Last
    if ($null -ne $nsPerOp -and $nsPerOp -gt 0) {
        $metrics.p95LatencyMs = [math]::Round($nsPerOp / 1000000.0, 3)
        $metrics.throughputOpsPerSec = [math]::Round(1000000000.0 / $nsPerOp, 2)
    }

    $tenantFairnessRatio = Get-FirstNumericMatch -Text $CommandOutput -Pattern 'tenant[_ ]fairness[_ ]ratio\s*(?:\||:|=)\s*([0-9]+(?:\.[0-9]+)?)' -Last
    if ($null -ne $tenantFairnessRatio) {
        $metrics.tenantFairnessRatio = [math]::Round($tenantFairnessRatio, 4)
    }

    $replayInvalidRows = Get-FirstNumericMatch -Text $CommandOutput -Pattern 'replay_invalid_rows\s*\n[-\s\+\|]*\n\s*([0-9]+)' -Last
    if ($null -eq $replayInvalidRows) {
        $replayInvalidRows = Get-FirstNumericMatch -Text $CommandOutput -Pattern 'replay[_ ]invalid[_ ]rows\s*(?:\||:|=)\s*([0-9]+)' -Last
    }
    if ($null -ne $replayInvalidRows) {
        $metrics.replayInvalidRows = [int]$replayInvalidRows
    }

    return $metrics
}

function Get-ThresholdEvaluation {
    param(
        [Parameter(Mandatory = $true)]
        [pscustomobject]$SelectedProfile,
        [Parameter(Mandatory = $true)]
        [System.Collections.IDictionary]$Metrics
    )

    $evaluation = @()
    foreach ($threshold in $SelectedProfile.successThresholds.PSObject.Properties) {
        $thresholdName = [string]$threshold.Name
        $thresholdTarget = [double]$threshold.Value
        $metricName = $null

        switch ($thresholdName) {
            'p95LatencyMs' { $metricName = 'p95LatencyMs' }
            'tenantFairnessRatioMax' { $metricName = 'tenantFairnessRatio' }
            'replayCorrectnessInvalidRowsMax' { $metricName = 'replayInvalidRows' }
            default { $metricName = $null }
        }

        $hasMetric = $false
        if ($null -ne $metricName) {
            if ($Metrics -is [System.Collections.IDictionary]) {
                $hasMetric = $Metrics.Contains($metricName)
            }
            else {
                $hasMetric = $null -ne $Metrics.PSObject.Properties[$metricName]
            }
        }

        if ($null -eq $metricName -or -not $hasMetric) {
            $evaluation += [ordered]@{
                threshold = $thresholdName
                target = $thresholdTarget
                status = 'not_evaluated'
                reason = "metric '$metricName' not found in output"
            }
            continue
        }

        $actualValue = [double]$Metrics[$metricName]
        $passed = $actualValue -le $thresholdTarget

        $evaluation += [ordered]@{
            threshold = $thresholdName
            target = $thresholdTarget
            actual = $actualValue
            status = if ($passed) { 'pass' } else { 'fail' }
            metric = $metricName
        }
    }

    return $evaluation
}

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
        [int]$DurationSeconds,
        [System.Collections.IDictionary]$AutoMetrics = @{},
        [array]$ThresholdEvaluation = @()
    )

    $thresholdLines = foreach ($property in $SelectedProfile.successThresholds.PSObject.Properties) {
        "- {0}: {1}" -f $property.Name, $property.Value
    }

    $controlLines = foreach ($control in $SelectedProfile.governanceControls) {
        "- $control"
    }

    $governanceBlock = if ($controlLines.Count -gt 0) { $controlLines -join [Environment]::NewLine } else { "- none" }
    $thresholdBlock = if ($thresholdLines.Count -gt 0) { $thresholdLines -join [Environment]::NewLine } else { "- none" }
    $autoMetricsBlock = if ($AutoMetrics.Count -eq 0) {
        "- none"
    }
    else {
        ($AutoMetrics.Keys | ForEach-Object { "- {0}: {1}" -f $_, $AutoMetrics[$_] }) -join [Environment]::NewLine
    }
    $thresholdEvaluationBlock = if ($ThresholdEvaluation.Count -eq 0) {
        "- none"
    }
    else {
        ($ThresholdEvaluation | ForEach-Object {
            $actualProperty = $_.PSObject.Properties['actual']
            $reasonProperty = $_.PSObject.Properties['reason']
            $actualDisplay = if ($null -ne $actualProperty) { $actualProperty.Value } else { 'n/a' }
            $reasonDisplay = if ($null -ne $reasonProperty) { [string]$reasonProperty.Value } else { '' }
            "- {0}: status={1}, target={2}, actual={3}{4}" -f $_.threshold, $_.status, $_.target, $actualDisplay, $(if ([string]::IsNullOrWhiteSpace($reasonDisplay)) { '' } else { ", reason=$reasonDisplay" })
        }) -join [Environment]::NewLine
    }

    $content = @"
# Phase 4 Evidence Summary

## Run Metadata
- Run ID: $RunId
- Profile: $($SelectedProfile.id) - $($SelectedProfile.name)
- Git Hash: $GitHash
- Tenant Count: $TenantCount
- Duration Seconds: $DurationSeconds
- Dry Run: $WasDryRun
- Command: $ExecutedCommand

## Governance Controls
$governanceBlock

## Threshold Targets
$thresholdBlock

## Observed Metrics
- p50 latency ms:
- p95 latency ms:
- p99 latency ms:
- throughput ops/sec:
- tenant fairness ratio:
- audit amplification:
- replay invalid rows:

## Auto-Extracted Metrics
$autoMetricsBlock

## Threshold Evaluation
$thresholdEvaluationBlock

## Result
- Status:
- Risks:
- Notes:
"@

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
    autoExtractedMetrics = @{}
    thresholdEvaluation = @()
}

if (-not $DryRun -and -not [string]::IsNullOrWhiteSpace($Command)) {
    Push-Location $workspaceRoot
    try {
        $nativePreferenceVariable = Get-Variable -Name PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue
        $priorNativePreference = $null
        $priorErrorActionPreference = $ErrorActionPreference
        if ($null -ne $nativePreferenceVariable) {
            $priorNativePreference = [bool]$nativePreferenceVariable.Value
            $script:PSNativeCommandUseErrorActionPreference = $false
        }

        try {
            $ErrorActionPreference = "Continue"
            $commandOutput = Invoke-Expression $Command 2>&1 | Out-String
        }
        finally {
            $ErrorActionPreference = $priorErrorActionPreference
            if ($null -ne $nativePreferenceVariable) {
                $script:PSNativeCommandUseErrorActionPreference = $priorNativePreference
            }
        }

        $manifest.commandExitCode = if ($null -ne (Get-Variable -Name LASTEXITCODE -ErrorAction SilentlyContinue)) { $LASTEXITCODE } else { 0 }
        $autoMetrics = Get-AutoExtractedMetrics -CommandOutput $commandOutput
        $manifest.autoExtractedMetrics = $autoMetrics
        $manifest.thresholdEvaluation = Get-ThresholdEvaluation -SelectedProfile $selectedProfile -Metrics $autoMetrics
        Set-Content -Path (Join-Path $runDirectory "command-output.txt") -Value $commandOutput -Encoding ascii
    }
    finally {
        Pop-Location
    }
}

$manifestPath = Join-Path $runDirectory "run-manifest.json"
$summaryPath = Join-Path $runDirectory "evidence-summary.md"
$manifest | ConvertTo-Json -Depth 6 | Set-Content -Path $manifestPath -Encoding ascii

$summaryArgs = @{
    FilePath = $summaryPath
    SelectedProfile = $selectedProfile
    GitHash = $gitHash
    RunId = $runId
    ExecutedCommand = $executedCommand
    WasDryRun = [bool]$DryRun
    TenantCount = $TenantCount
    DurationSeconds = $DurationSeconds
    AutoMetrics = $manifest.autoExtractedMetrics
    ThresholdEvaluation = [array]$manifest.thresholdEvaluation
}
New-EvidenceSummary @summaryArgs

Write-Host "Created Phase 4 evidence bundle: $runDirectory"
Write-Host "Manifest: $manifestPath"
Write-Host "Summary : $summaryPath"

if ($manifest.commandExitCode -ne $null) {
    Write-Host "Command exit code: $($manifest.commandExitCode)"
}