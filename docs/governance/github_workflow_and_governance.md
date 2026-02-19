# GitHub Governance & Workflow Strategy: QuantatomAI

To maintain "Ultra-Diamond" code quality across Dev, QA, and Production, we implement a strict governance model using GitHub's native protection rules and Actions.

## 1. Branching Strategy (GitFlow Optimized)

We adopt a modified GitFlow strategy optimized for Continuous Delivery.

| Environment | Branch | Purpose | Protection Rule |
| :--- | :--- | :--- | :--- |
| **Production** | `main` | Stable, deployed to Prod. | **Strict**: partial linear history, signed commits, require PR, require status checks. |
| **QA / Staging** | `release/*` | Release candidates (or dedicated `qa` branch). | **High**: Require PR, require integration tests. |
| **Development** | `develop` | Integration branch for features. | **Medium**: Require CI pass, squash merge. |
| **Feature** | `feat/*` | Individual work. | None (Developer freedom). |

## 1.1 Governance Phases (Scaling from Solo to Enterprise)
Since you are currently a **Solo Developer**, applying full enterprise rules (like mandatory peer reviews) will block you. We use a **Phased Approach**:

### Phase 1: Solo Mode (Current)
*   **Structure:** Keep the `main` / `develop` branch structure. It builds muscle memory.
*   **CI Checks:** **Keep enabled.** Automated tests (`test-go`, `test-ui`) are your "digital pair programmer."
*   **Code Owners:** Keep the file definition (for documentation), but **disable** "Require review from Code Owners" in GitHub settings.
*   **PRs:** Create PRs for your own code to trigger CI, but allow **Merge without Review** (or self-approval).

### Phase 2: Team Mode (Future)
*   **Enable** "Require 1 reviewer" on `main`.
*   **Enable** "Require review from Code Owners".
*   **Enforce** Segregation of Duties (Dev cannot deploy to Prod without approval).

## 2. Branch Protection Rules

### `main` (Production)
*   **Require pull request reviews before merging:**
    *   Required approvals: **1** (or 2 for high security).
    *   *Dismiss stale pull request approvals when new commits are pushed.*
    *   *Require review from Code Owners.*
*   **Require status checks to pass before merging:**
    *   `test-go (grid-service)`
    *   `test-rust (atom-engine)`
    *   `test-ui (web)`
    *   `lint-go` / `clippy`
    *   `security-scan`
*   **Require signed commits:** Yes.
*   **Require linear history:** Yes (Squash or Rebase merges only).
*   **Include administrators:** Yes.

### `develop` (Integration)
*   **Require status checks to pass:** Yes (Basic tests & lint).
*   **Require linear history:** Recommended.
*   **Allow force pushes:** No.

## 3. CI/CD Pipeline Structure

All PRs must pass the "Quality Gate" before merging.

### Quality Gate (CI)
Trigger: `pull_request` to `main`, `develop`.

1.  **Lint & Format:**
    *   Go: `golangci-lint`
    *   Rust: `cargo fmt --check`, `cargo clippy`
    *   TS/React: `eslint`, `prettier`
2.  **Unit Tests:**
    *   `go test ./... -race`
    *   `cargo test`
    *   `npm test`
3.  **Security Scan:**
    *   `gosec` (Go Security)
    *   `cargo-audit` (Rust Dependencies)
    *   `npm audit` (Frontend)
4.  **Build Verification:**
    *   `docker build` (Ensure container builds succeed).

### Deployment Gate (CD)
Trigger: `push` to `main` (Prod) or `release/*` (QA).

1.  **Build & Tag:**
    *   Build Docker images.
    *   Tag with semantic version + commit SHA.
    *   Push to ECR/Registry.
2.  **Deploy to Environment:**
    *   **QA:** Auto-deploy to QA cluster.
    *   **Prod:** Manual approval (GitHub Environment Protection) -> Deploy to Prod.

## 4. GitHub Security Settings

1.  **Dependabot:** Enable Version Updates and Security Updates.
2.  **Secret Scanning:** Enable.
3.  **Code Scanning (CodeQL):** Enable for Go, Rust, TypeScript.

## 5. Implementation Steps

1.  **Repo Settings:** Apply Branch Protection Rules (Manual).
2.  **Workflows:** Create `.github/workflows/ci.yml` (Unified or Per-Service).

## 6. Multi-Team Collaboration (Scaling the Monorepo)

When multiple teams (e.g., Grid Team, Compute Team, Frontend Team) work in the same repository, we use **Domain Ownership** to prevent stepping on each other's toes.

### A. The `CODEOWNERS` File
We assign specific directories to GitHub Teams. This ensures that a change to the *Grid Service* automatically requests a review from the *Grid Team*.

**Example `.github/CODEOWNERS`:**
```
# Global Owners (Architects)
*       @quantatom/architects

# Grid Service Team
/services/grid-service/   @quantatom/grid-team
/services/mdf-store/      @quantatom/grid-team

# Compute Team (Rust Engine)
/services/atom-engine/    @quantatom/compute-team
/compute/heliocalc/       @quantatom/compute-team

# Frontend Team
/ui/web/                  @quantatom/frontend-team
/ui/excel-addin/          @quantatom/frontend-team

# Infrastructure Team
/infra/                   @quantatom/sre-team
/.github/                 @quantatom/sre-team
```

### B. CI/CD Path Filtering
To speed up PRs, CI workflows should only run for the changed service.
*   *Example:* A change to `ui/web` should trigger `test-ui` but **not** `test-rust`.
*   We use GitHub Actions `paths` filter:
    ```yaml
    on:
      pull_request:
        paths:
          - 'services/grid-service/**'
    ```

### C. Release Train vs. Independent Release
1.  **Independent Release (Recommended):** Each service (`grid-service`, `atom-engine`) has its own version tag (e.g., `grid-v1.2.0`, `engine-v0.5.0`). Teams release at their own cadence.
2.  **Platform Release:** A "Release Train" happens monthly for the on-prem/enterprise bundle, aggregating the stable versions of all services.

### D. Breaking Changes Contract
If Team A (Grid) needs to change an API used by Team B (Frontend):
1.  **Deprecation Notice:** Mark the field `deprecated` in Protobuf/GraphQL.
2.  **Parallell Run:** Support both old and new fields for 1 release cycle.
3.  **Removal:** Remove old field after Team B has migrated.
