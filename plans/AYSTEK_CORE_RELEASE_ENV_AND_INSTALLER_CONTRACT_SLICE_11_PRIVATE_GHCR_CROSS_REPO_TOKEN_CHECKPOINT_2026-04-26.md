# AYSTEK Core Release Env And Installer Contract
## Slice 11 Checkpoint: Private GHCR Cross-Repo Token

Date: 2026-04-26

### Problem
- `aystek-core` workflow failed in `Verify Matching Panel Tag`
- GHCR returned:
  - `403 Forbidden`
- Root cause:
  - `GITHUB_TOKEN` from one private repo cannot reliably inspect a private GHCR package published by another private repo

### Completed
- Core release workflow now prefers `GHCR_READ_TOKEN` and falls back to `GITHUB_TOKEN`:
  - `/workspace/go/aystek-core/.github/workflows/release.yml`
- Panel release workflow now prefers `GHCR_READ_TOKEN` for the reciprocal visibility check:
  - `/workspace/node/aystek-core-next/.github/workflows/release.yml`
- Runbook documents the required shared secret and scope:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Operator Action
Define repository secret `GHCR_READ_TOKEN` in:
- `huseyinkoc-aystek/aystek-core`
- `huseyinkoc-aystek/admin-panel-next`

Required scope:
- `read:packages`

### Result
- Cross-repo private GHCR parity checks can authenticate consistently.
- Matching tag visibility and tag-lock checks no longer depend on repository-local `GITHUB_TOKEN` behavior.
