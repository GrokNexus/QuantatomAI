[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet("A", "B", "C", "D")]
    [string]$Profile,

    [switch]$PrepareFixtures,
    [switch]$DatabaseBacked,
    [switch]$DryRun,
    [string]$DatabaseHost = "localhost",
    [string]$DatabaseUser = "quantatomai",
    [string]$DatabaseName = "quantatomai",
    [int]$TenantCount = 4,
    [int]$DurationSeconds = 300
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Get-LastExitCodeOrZero {
    $lastExitCodeVariable = Get-Variable -Name LASTEXITCODE -ErrorAction SilentlyContinue
    if ($null -eq $lastExitCodeVariable) {
        return 0
    }

    return [int]$lastExitCodeVariable.Value
}

$workspaceRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$phase4Runner = Join-Path $PSScriptRoot "run-phase4-profile.ps1"
$fixturePreparer = Join-Path $PSScriptRoot "prepare-grid-service-phase4-fixtures.ps1"

$commandMap = @{
    "A" = "go -C services/grid-service test -run ^$ -bench BenchmarkGridQueryServiceSingleRecord -benchmem ./pkg/orchestration"
    "B" = "go -C services/grid-service test -run ^$ -bench BenchmarkGridQueryServiceLargeRecord -benchmem ./pkg/orchestration"
    "C" = "go -C services/grid-service test -run ^$ -bench BenchmarkLWWElementSetMergeSequential -benchmem ./pkg/sync"
    "D" = "go -C services/grid-service test -run ^$ -bench BenchmarkLWWElementSetMergeConflictHeavy -benchmem ./pkg/sync"
}

$databaseCommandMap = @{
    "B" = "docker cp services/grid-service/sql/validation/phase4_planning_workload_smoke_checks.sql gridservice-postgres:/tmp/phase4_planning_workload_smoke_checks.sql; docker cp services/grid-service/sql/validation/phase6_consolidation_domain_checks.sql gridservice-postgres:/tmp/phase6_consolidation_domain_checks.sql; docker cp services/grid-service/sql/validation/phase7_ai_inference_governance_checks.sql gridservice-postgres:/tmp/phase7_ai_inference_governance_checks.sql; docker exec -i gridservice-postgres psql -U quantatomai -d quantatomai -v ON_ERROR_STOP=1 -f /tmp/phase4_planning_workload_smoke_checks.sql; docker exec -i gridservice-postgres psql -U quantatomai -d quantatomai -v ON_ERROR_STOP=1 -f /tmp/phase6_consolidation_domain_checks.sql; docker exec -i gridservice-postgres psql -U quantatomai -d quantatomai -v ON_ERROR_STOP=1 -f /tmp/phase7_ai_inference_governance_checks.sql"
    "C" = "docker cp services/grid-service/sql/validation/phase4_fixture_smoke_checks.sql gridservice-postgres:/tmp/phase4_fixture_smoke_checks.sql; docker exec -i gridservice-postgres psql -U quantatomai -d quantatomai -v ON_ERROR_STOP=1 -f /tmp/phase4_fixture_smoke_checks.sql"
    "D" = "docker cp services/grid-service/sql/validation/phase4_fixture_smoke_checks.sql gridservice-postgres:/tmp/phase4_fixture_smoke_checks.sql; docker exec -i gridservice-postgres psql -U quantatomai -d quantatomai -v ON_ERROR_STOP=1 -f /tmp/phase4_fixture_smoke_checks.sql"
}

$selectedCommand = if ($DatabaseBacked -and $databaseCommandMap.ContainsKey($Profile)) {
    $databaseCommandMap[$Profile]
}
else {
    $commandMap[$Profile]
}

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
        $fixtureExitCode = Get-LastExitCodeOrZero
        if ($fixtureExitCode -ne 0) {
            throw "Fixture preparation exited with code $fixtureExitCode"
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
    $runnerExitCode = Get-LastExitCodeOrZero
    if ($runnerExitCode -ne 0) {
        throw "Phase 4 runner exited with code $runnerExitCode"
    }
}
finally {
    Pop-Location
}