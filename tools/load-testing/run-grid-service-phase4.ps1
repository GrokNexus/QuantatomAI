[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("A", "B", "C", "D")]
    [string]$Profile,

    [switch]$PrepareFixtures,
    [switch]$DryRun,
    [string]$DatabaseHost = "localhost",
    [string]$DatabaseUser = "quantatomai",
    [string]$DatabaseName = "quantatomai",
    [int]$TenantCount = 4,
    [int]$DurationSeconds = 300
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$workspaceRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$phase4Runner = Join-Path $PSScriptRoot "run-phase4-profile.ps1"
$fixturePreparer = Join-Path $PSScriptRoot "prepare-grid-service-phase4-fixtures.ps1"

$commandMap = @{
    "A" = "go test -run ^$ -bench BenchmarkGridQueryServiceSingleRecord -benchmem ./services/grid-service/pkg/orchestration"
    "B" = "go test -run ^$ -bench BenchmarkGridQueryServiceLargeRecord -benchmem ./services/grid-service/pkg/orchestration"
    "C" = "go test -run ^$ -bench BenchmarkLWWElementSetMergeSequential -benchmem ./services/grid-service/pkg/sync"
    "D" = "go test -run ^$ -bench BenchmarkLWWElementSetMergeConflictHeavy -benchmem ./services/grid-service/pkg/sync"
}

$selectedCommand = $commandMap[$Profile]
if (-not $selectedCommand) {
    throw "No grid-service command is mapped for profile '$Profile'"
}

Push-Location $workspaceRoot
try {
    if ($PrepareFixtures) {
        $fixtureArgs = @{
            Profile = $Profile
            DatabaseHost = $DatabaseHost
            DatabaseUser = $DatabaseUser
            DatabaseName = $DatabaseName
        }
        if ($DryRun) {
            $fixtureArgs.DryRun = $true
        }

        & $fixturePreparer @fixtureArgs
        if ($LASTEXITCODE -ne 0) {
            throw "Fixture preparation exited with code $LASTEXITCODE"
        }
    }

    $runnerArgs = @{
        Profile = $Profile
        Command = $selectedCommand
        TenantCount = $TenantCount
        DurationSeconds = $DurationSeconds
    }

    if ($DryRun) {
        $runnerArgs.DryRun = $true
    }

    & $phase4Runner @runnerArgs
    if ($LASTEXITCODE -ne 0) {
        throw "Phase 4 runner exited with code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}