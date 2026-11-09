# AYSTEK Core And Panel Dual Repo Release Runbook

Date: 2026-04-24
Status: Active
Scope: `aystek-core`, `aystek-core-next`, `aystek-core-installer`, `aysteknoloji-com`

Related:
- `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_CI_FIRST_RUN_VALIDATION_RUNBOOK_2026-04-24.md`

## Objective

Coordinate core, panel, installer, and public-site releases as one operational deployment line.

## Canonical Image Coordinates

- Default private namespace: `ghcr.io/<repository_owner>/...`
- Optional override namespace: repo variable `GHCR_NAMESPACE`
- Core: `ghcr.io/<namespace>/aystek-core:<version>`
- Panel: `ghcr.io/<namespace>/aystek-panel:<version>`
- Public site: `ghcr.io/<namespace>/aysteknoloji-com:<version>`

## Release Order

1. Prepare matching version tag for both repos.
2. Release panel repo:
   - verify image pushed to `ghcr.io/aystek/aystek-panel:<version>`
3. Release core repo:
   - workflow now verifies matching panel tag lock before continuing
   - verify archive and image pushed
4. Release public site repo if a new public deployment artifact is needed.
5. Update installer config / deployment env with matching core and panel tags.

## Why This Order

- panel image availability must be confirmed before core release instructions point operators to a matching tag
- installer operators should not be exposed to half-published tag pairs

## Operator Validation

```bash
docker pull ghcr.io/aystek/aystek-panel:<version>
docker pull ghcr.io/aystek/aystek-core:<version>
```

## Installer Modes

- `docker + registry`:
  - use matching core/panel tags
- `docker + local-source`:
  - no registry dependency
- `native + local-artifacts`:
  - release tarball + panel bundle
- `native + local-source`:
  - local repos on target host

## Coordination Checklist

1. Panel release workflow succeeded.
2. Matching panel image tag exists in GHCR before core release starts.
3. Core release workflow succeeded.
4. Installer `source_mode` selected intentionally.
5. If using registry mode, both images are pullable.
6. If using local-source mode, target host prerequisites are present.

## Tag Lock Rule

- `aystek-core`, `aystek-core-next`, and `aysteknoloji-com` are versioned independently.
- Compatibility is defined explicitly with repo variables, not by assuming equal semantic tags.
- `aystek-core` release workflow fails fast if the configured compatible panel image tag is not already published.
- `aystek-core-next` release workflow reports visibility for the configured compatible core image tag when provided.
- This prevents half-published release pairs from being documented to operators.
- This asymmetry is intentional; blocking both sides would create a release deadlock.

## Operational Flow

1. Decide release pair explicitly:
   - core tag, example: `v2.0.4`
   - panel tag, example: `v1.0.5`
   - public site tag, example: `v1.0.0`
2. In `aystek-core`, set repo variable `PANEL_COMPAT_TAG` to the chosen panel tag.
3. Optionally in `aystek-core-next`, set repo variable `CORE_COMPAT_TAG` to the chosen core tag for visibility reporting.
4. Optionally in `aysteknoloji-com`, set repo variable `CORE_COMPAT_TAG` to the chosen core tag for visibility reporting.
5. Release panel tag and wait for GHCR publish success.
6. Release core tag.
7. Let core workflow verify the configured compatible panel tag automatically.
8. Release public site tag when public deployment artifact is needed.
9. Only after required artifacts are available, proceed with installer-driven deployment.

## Namespace Rule

- If `GHCR_NAMESPACE` repo variable is not set, both workflows publish into `github.repository_owner`.
- If a dedicated package namespace is desired, set the same `GHCR_NAMESPACE` value in both repos.
- Mixed namespace configuration is invalid for coordinated release.

## Compatibility Variables

- In `aystek-core`:
  - `PANEL_COMPAT_TAG`
  - example: `v1.0.5`
- In `aystek-core-next`:
  - `CORE_COMPAT_TAG`
  - optional, for visibility/reporting
  - example: `v2.0.4`
- In `aysteknoloji-com`:
  - `CORE_COMPAT_TAG`
  - optional, for visibility/reporting
  - example: `v2.0.4`

## Private GHCR Access Rule

- `GITHUB_TOKEN` is sufficient for publishing the package of the current repository.
- It is not always sufficient for reading a private GHCR package published by a different private repository.
- For cross-repo tag visibility and tag-lock checks, define the same repository secret in both repos:
  - `GHCR_READ_TOKEN`
- Required token scope:
  - `read:packages`
- Recommended:
  - fine-grained or classic token limited to the package owner account and private package read access

Affected checks:
- `aystek-core` → `Verify Matching Panel Tag`
- `aystek-core-next` → `Check Matching Core Tag Visibility`

## CI Validation Status

Current state:
- workflow and runbook preparation complete
- waiting for first real GitHub Actions executions from GitHub UI

After first successful runs, append:
- `Validated in CI: yes`
- first successful installer workflow URL
- first successful panel workflow URL
- first successful core workflow URL
- final namespace and image coordinates used
