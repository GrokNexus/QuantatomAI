# Contributing to QuantatomAI

This project uses a safety-first workflow to prevent code loss while keeping changes professional and reviewable.

## Core Rules

- Never delete backup refs (`backup/*`, `autosnap-*`) unless intentionally pruning history.
- Prefer small, focused commits with clear messages.
- Keep code and docs together in the same PR when both are impacted.
- Use feature branches for all non-trivial work.
- Keep `main` deployable.

## Branch Strategy

- `main`: stable, protected branch.
- `feat/<scope>-<short-name>`: feature work.
- `fix/<scope>-<short-name>`: bug fixes.
- `chore/<scope>-<short-name>`: maintenance, CI, docs, or tooling.

Example:

```bash
git checkout -b feat/grid-metadata-cache
```

## Daily Safe Workflow

```bash
git fetch --all --tags
git checkout main
git pull --rebase origin main
git checkout -b feat/<scope>-<name>

# work
git add -A
git commit -m "feat(scope): short description"
git push -u origin feat/<scope>-<name>
```

Open a PR from your branch to `main`.

## Commit Message Convention

Use Conventional Commit style:

- `feat:` new functionality
- `fix:` bug fix
- `docs:` documentation only
- `chore:` maintenance/tooling
- `refactor:` internal restructuring without behavior change
- `test:` tests only

Examples:

- `feat(grid): add metadata hierarchy resolver`
- `fix(ui): correct grid query pagination`
- `docs(architecture): clarify WRM flow`

## Pull Request Quality Checklist

Before requesting review:

- Scope is focused and minimal.
- CI workflows pass.
- No unrelated file churn.
- Docs updated for behavior/config changes.
- Rollback path is clear.

## Backup & Recovery

Automated safety snapshots run on every push to `main`:

- Moving backup branch: `backup/main-latest`
- Immutable snapshot tags: `autosnap-YYYYMMDD-HHMMSS-<sha7>`

Recovery procedures are documented in:

- `docs/dev-guides/github-recovery-playbook.md`

## Force Push Policy

- Avoid force-push on shared branches.
- If absolutely required, use:

```bash
git push --force-with-lease
```

Never use raw `--force` on `main`.

