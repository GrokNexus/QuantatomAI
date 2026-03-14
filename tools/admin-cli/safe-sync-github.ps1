[CmdletBinding()]
param(
    [string]$RemoteName = "origin",
    [string]$BranchName = "main",
    [string]$BackupPrefix = "backup/local-pre-sync",
    [switch]$SkipBackupPush
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-Git {
    param(
        [Parameter(Mandatory = $true)]
        [string[]]$Arguments
    )

    $output = & git @Arguments 2>&1
    if ($LASTEXITCODE -ne 0) {
        $message = ($output | Out-String).Trim()
        throw "git $($Arguments -join ' ') failed: $message"
    }

    return ($output | Out-String).Trim()
}

function Write-Section {
    param([string]$Message)
    Write-Host ""
    Write-Host "== $Message ==" -ForegroundColor Cyan
}

Write-Section "Verify repository"
$insideWorkTree = Invoke-Git -Arguments @("rev-parse", "--is-inside-work-tree")
if ($insideWorkTree -ne "true") {
    throw "This script must be run inside a Git working tree."
}

$workingTreeStatus = Invoke-Git -Arguments @("status", "--porcelain")
if ($workingTreeStatus) {
    throw "Working tree is not clean. Commit or stash local changes before syncing."
}

$currentBranch = Invoke-Git -Arguments @("rev-parse", "--abbrev-ref", "HEAD")
if ($currentBranch -ne $BranchName) {
    throw "Current branch is '$currentBranch'. Switch to '$BranchName' before syncing."
}

Write-Section "Refresh remote state"
Invoke-Git -Arguments @("fetch", "--prune", $RemoteName) | Out-Null

$localRef = $BranchName
$remoteRef = "$RemoteName/$BranchName"

$localHead = Invoke-Git -Arguments @("rev-parse", $localRef)
$remoteHead = Invoke-Git -Arguments @("rev-parse", $remoteRef)

Write-Host "Local HEAD : $localHead"
Write-Host "Remote HEAD: $remoteHead"

Write-Section "Create local safety branch"
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$shortHead = Invoke-Git -Arguments @("rev-parse", "--short", $localRef)
$backupBranch = "$BackupPrefix-$timestamp-$shortHead"
Invoke-Git -Arguments @("branch", $backupBranch, $localRef) | Out-Null
Write-Host "Created local backup branch: $backupBranch"

if (-not $SkipBackupPush) {
    Write-Section "Push safety branch to GitHub"
    Invoke-Git -Arguments @("push", $RemoteName, "$backupBranch`:$backupBranch") | Out-Null
    Write-Host "Pushed backup branch: $RemoteName/$backupBranch"
}

Write-Section "Compare local and remote"
$counts = Invoke-Git -Arguments @("rev-list", "--left-right", "--count", "$localRef...$remoteRef")
$parts = $counts -split "\s+"
if ($parts.Count -lt 2) {
    throw "Unable to parse ahead/behind counts from: $counts"
}

[int]$aheadCount = $parts[0]
[int]$behindCount = $parts[1]

Write-Host "Ahead : $aheadCount"
Write-Host "Behind: $behindCount"

if ($aheadCount -eq 0 -and $behindCount -eq 0) {
    Write-Section "Result"
    Write-Host "Local and GitHub are already in sync. No push needed." -ForegroundColor Green
    exit 0
}

if ($aheadCount -gt 0 -and $behindCount -eq 0) {
    Write-Section "Push local commits"
    Invoke-Git -Arguments @("push", $RemoteName, $BranchName) | Out-Null
    Write-Host "Pushed local '$BranchName' to '$RemoteName/$BranchName'." -ForegroundColor Green
    exit 0
}

if ($aheadCount -eq 0 -and $behindCount -gt 0) {
    Write-Section "Fast-forward local branch"
    Invoke-Git -Arguments @("pull", "--ff-only", $RemoteName, $BranchName) | Out-Null
    Write-Host "Fast-forwarded local '$BranchName' to match '$RemoteName/$BranchName'." -ForegroundColor Green
    exit 0
}

Write-Section "Diverged histories"
$reconcileBranch = "sync/$timestamp-$BranchName-reconcile"
Invoke-Git -Arguments @("switch", "-c", $reconcileBranch) | Out-Null
Write-Host "Created reconcile branch: $reconcileBranch"
Write-Host "Local and GitHub both contain unique commits."
Write-Host "Next step: merge '$remoteRef' into '$reconcileBranch', resolve conflicts there, test, then push the reconcile branch and open a PR."
Write-Host "Suggested commands:"
Write-Host "  git merge $remoteRef"
Write-Host "  git push $RemoteName $reconcileBranch"
exit 0
