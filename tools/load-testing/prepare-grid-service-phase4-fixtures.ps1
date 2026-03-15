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

    $canExecuteDirectly = $null -ne (Get-Command -Name $executable -ErrorAction SilentlyContinue)
    $useDockerPsqlFallback = (-not $canExecuteDirectly) -and ($executable -eq "psql")
    if (-not $canExecuteDirectly -and -not $useDockerPsqlFallback) {
        throw "Fixture step '$($Step.name)' requires executable '$executable', but it is not available in PATH"
    }

    $priorEnvironment = @{}
    foreach ($environmentVariable in $environmentVariables.GetEnumerator()) {
        $priorEnvironment[$environmentVariable.Key] = [Environment]::GetEnvironmentVariable($environmentVariable.Key, 'Process')
        [Environment]::SetEnvironmentVariable($environmentVariable.Key, $environmentVariable.Value, 'Process')
    }

    $nativePreferenceVariable = Get-Variable -Name PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue
    $priorNativePreference = $null
    if ($null -ne $nativePreferenceVariable) {
        $priorNativePreference = [bool]$nativePreferenceVariable.Value
        $script:PSNativeCommandUseErrorActionPreference = $false
    }

    try {
        if ($useDockerPsqlFallback) {
            Write-Host "Executable '$executable' not found. Falling back to docker exec against gridservice-postgres."

            $resolvedArguments = @($arguments)
            for ($index = 0; $index -lt $resolvedArguments.Count; $index++) {
                if ($resolvedArguments[$index] -eq "-f" -and $index + 1 -lt $resolvedArguments.Count) {
                    $hostScriptPath = $resolvedArguments[$index + 1]
                    $absoluteHostScriptPath = if ([System.IO.Path]::IsPathRooted($hostScriptPath)) {
                        $hostScriptPath
                    }
                    else {
                        Join-Path $workspaceRoot $hostScriptPath
                    }

                    if (-not (Test-Path $absoluteHostScriptPath)) {
                        throw "Fixture SQL file '$hostScriptPath' was not found at '$absoluteHostScriptPath'"
                    }

                    $containerScriptPath = "/tmp/{0}" -f ([System.IO.Path]::GetFileName($absoluteHostScriptPath))
                    & docker cp $absoluteHostScriptPath ("gridservice-postgres:{0}" -f $containerScriptPath)
                    $copyExitCodeVariable = Get-Variable -Name LASTEXITCODE -ErrorAction SilentlyContinue
                    $copyExitCode = if ($null -eq $copyExitCodeVariable) { 0 } else { [int]$copyExitCodeVariable.Value }
                    if ($copyExitCode -ne 0) {
                        throw "Failed to copy fixture SQL file '$absoluteHostScriptPath' into gridservice-postgres"
                    }

                    $resolvedArguments[$index + 1] = $containerScriptPath
                }
            }

            $dockerArguments = @("exec", "-i")
            foreach ($environmentVariable in $environmentVariables.GetEnumerator()) {
                $dockerArguments += "-e"
                $dockerArguments += ("{0}={1}" -f $environmentVariable.Key, $environmentVariable.Value)
            }

            $dockerArguments += "gridservice-postgres"
            $dockerArguments += "psql"
            $dockerArguments += $resolvedArguments

            & docker @dockerArguments
        }
        else {
            & $executable @arguments
        }

        $stepExitCodeVariable = Get-Variable -Name LASTEXITCODE -ErrorAction SilentlyContinue
        $stepExitCode = if ($null -eq $stepExitCodeVariable) { 0 } else { [int]$stepExitCodeVariable.Value }
        if ($stepExitCode -ne 0) {
            throw "Fixture step '$($Step.name)' failed with code $stepExitCode"
        }
    }
    finally {
        if ($null -ne $nativePreferenceVariable) {
            $script:PSNativeCommandUseErrorActionPreference = $priorNativePreference
        }

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