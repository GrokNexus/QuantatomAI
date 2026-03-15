[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("A", "B", "C", "D")]
    [string]$Profile,

    [string]$DatabaseHost = "localhost",
    [string]$DatabaseUser = "quantatomai",
    [string]$DatabaseName = "quantatomai",
    [switch]$DryRun
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$databasePassword = [Environment]::GetEnvironmentVariable("PGPASSWORD", "Process")
if ([string]::IsNullOrWhiteSpace($databasePassword)) {
    $databasePassword = "quantatomai"
}

function Resolve-StepToken {
    param(
        [AllowNull()]
        [string]$Value
    )

    if ($null -eq $Value) {
        return $null
    }

    $resolvedValue = $Value
    $resolvedValue = $resolvedValue.Replace("{{DatabaseHost}}", $DatabaseHost)
    $resolvedValue = $resolvedValue.Replace("{{DatabaseUser}}", $DatabaseUser)
    $resolvedValue = $resolvedValue.Replace("{{DatabaseName}}", $DatabaseName)
    $resolvedValue = $resolvedValue.Replace("{{DatabasePassword}}", $databasePassword)

    return $resolvedValue
}

function Invoke-StructuredFixtureStep {
    param(
        [Parameter(Mandatory = $true)]
        [pscustomobject]$Step
    )

    $executable = Resolve-StepToken -Value ([string]$Step.executable)
    $arguments = @()
    $argumentProperty = $Step.PSObject.Properties['arguments']
    if ($argumentProperty -and $null -ne $argumentProperty.Value) {
        foreach ($argument in $argumentProperty.Value) {
            $arguments += Resolve-StepToken -Value ([string]$argument)
        }
    }

    $environmentVariables = @{}
    $environmentProperty = $Step.PSObject.Properties['environment']
    if ($environmentProperty -and $null -ne $environmentProperty.Value) {
        foreach ($property in $environmentProperty.Value.PSObject.Properties) {
            $environmentVariables[$property.Name] = Resolve-StepToken -Value ([string]$property.Value)
        }
    }

    Write-Host "Executable: $executable"
    Write-Host "Arguments : $($arguments -join ' ')"
    foreach ($environmentVariable in $environmentVariables.GetEnumerator()) {
        $displayValue = if ($environmentVariable.Key -match 'PASSWORD|TOKEN|SECRET') { '<redacted>' } else { $environmentVariable.Value }
        Write-Host "Environment: $($environmentVariable.Key)=$displayValue"
    }

    if ($DryRun) {
        return
    }

    $priorEnvironment = @{}
    foreach ($environmentVariable in $environmentVariables.GetEnumerator()) {
        $priorEnvironment[$environmentVariable.Key] = [Environment]::GetEnvironmentVariable($environmentVariable.Key, 'Process')
        [Environment]::SetEnvironmentVariable($environmentVariable.Key, $environmentVariable.Value, 'Process')
    }

    try {
        & $executable @arguments
        if ($LASTEXITCODE -ne 0) {
            throw "Fixture step '$($Step.name)' failed with code $LASTEXITCODE"
        }
    }
    finally {
        foreach ($environmentVariable in $environmentVariables.GetEnumerator()) {
            [Environment]::SetEnvironmentVariable($environmentVariable.Key, $priorEnvironment[$environmentVariable.Key], 'Process')
        }
    }
}

$workspaceRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$fixtureManifestPath = Join-Path $PSScriptRoot "grid-service-phase4-fixtures.json"
$fixtureManifest = Get-Content -Path $fixtureManifestPath -Raw | ConvertFrom-Json
$selectedProfile = $fixtureManifest.profiles | Where-Object { $_.id -eq $Profile } | Select-Object -First 1

if (-not $selectedProfile) {
    throw "Fixture profile '$Profile' not found in $fixtureManifestPath"
}

Write-Host "Grid-service fixture profile: $($selectedProfile.id)"
Write-Host "Description: $($selectedProfile.description)"

if (-not $selectedProfile.requiresFixtures) {
    Write-Host "No fixture preparation required for this profile."
    exit 0
}

Push-Location $workspaceRoot
try {
    foreach ($step in $selectedProfile.steps) {
        Write-Host "Preparing step: $($step.name)"
        Invoke-StructuredFixtureStep -Step $step
    }
}
finally {
    Pop-Location
}