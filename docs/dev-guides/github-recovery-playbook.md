# GitHub Recovery Playbook

This repository now includes automated safety snapshots via `.github/workflows/repo-safety-snapshots.yml`.

## What is protected

- Every push to `main` updates `backup/main-latest` to the same commit.
- Every push to `main` also creates an immutable tag in the format:
  - `autosnap-YYYYMMDD-HHMMSS-<sha7>`

## Quick checks

```bash
git fetch --all --tags
git branch -r | grep backup/main-latest
git tag -l "autosnap-*" | tail -n 10
```

## Recover main from backup branch

```bash
git fetch --all --tags
git checkout main
git reset --hard origin/backup/main-latest
git push origin main --force-with-lease
```

## Recover main from a snapshot tag

```bash
git fetch --all --tags
git checkout main
git reset --hard autosnap-YYYYMMDD-HHMMSS-<sha7>
git push origin main --force-with-lease
```

## Restore a deleted local branch safely

```bash
git fetch --all --tags
git checkout -b restore-work autosnap-YYYYMMDD-HHMMSS-<sha7>
```

## Guardrails to follow

- Prefer `git push --force-with-lease` over `--force`.
- Keep work on feature branches and merge to `main`.
- Do not delete `backup/main-latest` or `autosnap-*` tags unless intentionally pruning history.
