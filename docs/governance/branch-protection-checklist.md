# Branch Protection Checklist (GitHub)

Apply this to `main` in **Settings → Branches → Add branch protection rule**.

## Rule pattern

- Branch name pattern: `main`

## Required settings (recommended now)

- ✅ Require a pull request before merging
- ✅ Require status checks to pass before merging
- ✅ Require branches to be up to date before merging
- ✅ Require conversation resolution before merging
- ✅ Require linear history
- ✅ Do not allow force pushes
- ✅ Do not allow deletions
- ✅ Include administrators

## Solo mode adjustment (to avoid blocking yourself)

- Keep required reviewers at `0` or allow self-review if you are the only maintainer.
- Do **not** enable "Require review from Code Owners" until your `@quantatom/*` teams exist and are staffed.

## Team mode upgrade (later)

- Require at least `1` approval.
- Enable "Require review from Code Owners".
- Add deployment protection rules for production environments.

## Status checks to include

Use the checks that already run in this repository (names may vary by workflow):

- Build and test checks
- Service-specific checks (grid, frontend, atom engine)
- Security scan checks

Tip: merge one PR, then return here and mark exact check names as required once they appear in GitHub UI.
