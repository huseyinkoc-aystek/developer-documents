# AYSTEK Core CI First-Run Validation Runbook

Date: 2026-04-24
Status: Ready for operator execution
Scope:
- `huseyinkoc-aystek/aystek-core-installer`
- `huseyinkoc-aystek/admin-panel-next`
- `huseyinkoc-aystek/aystek-core`
- `huseyinkoc-aystek/aysteknoloji-com`

## Objective

Validate the first real GitHub Actions runs for:
- installer local-source disposable smoke
- panel image release
- core package/image release with panel tag lock
- public site image release

This runbook is intended for the first operator-driven execution from GitHub UI.

## Workflow Targets

### 1. Installer
- Repository: `huseyinkoc-aystek/aystek-core-installer`
- Workflow: `Installer P0 Disposable Matrix`
- Trigger:
  - `workflow_dispatch` from GitHub UI
  - keep defaults unless a different core/panel repo must be tested

Expected inputs:
- `aystek_core_repo = huseyinkoc-aystek/aystek-core`
- `aystek_panel_repo = huseyinkoc-aystek/admin-panel-next`

Expected success criteria:
- `source-mode-contract` job passes
- `Local-source prerequisite smoke` step passes
- `docker + local-source` smoke runs on runner with Docker available
- `native + local-source` smoke passes
- `p0-disposable` job passes

Evidence to capture:
- workflow URL
- job names
- screenshots or copied status lines for:
  - `Validate source_mode fixtures`
  - `Local-source prerequisite smoke`
  - `Run CI disposable P0 subset`

### 2. Panel
- Repository: `huseyinkoc-aystek/admin-panel-next`
- Workflow: `Release - Panel Docker Image`
- Trigger:
  - push/create release tag from GitHub UI, example `vX.Y.Z`

Expected success criteria:
- panel image pushed to `ghcr.io/<namespace>/aystek-panel:X.Y.Z`
- release body shows:
  - matching core tag visibility flag
  - expected core image path

Evidence to capture:
- workflow URL
- image path used
- release URL
- release note snippet for `Core Tag Visibility`

### 3. Core
- Repository: `huseyinkoc-aystek/aystek-core`
- Workflow: `Release & Package`
- Trigger:
  - push/create matching release tag from GitHub UI, example `vX.Y.Z`

Expected success criteria:
- `Verify Matching Panel Tag` passes
- package artifact uploaded
- docker image pushed to `ghcr.io/<namespace>/aystek-core:X.Y.Z`
- release created with canonical installer snippet

Evidence to capture:
- workflow URL
- `Verify Matching Panel Tag` step output
- release asset list
- release body image/tag examples

### 4. Public Site
- Repository: `huseyinkoc-aystek/aysteknoloji-com`
- Workflow: `Release - Public Site Docker Image`
- Trigger:
  - push/create release tag from GitHub UI, example `vX.Y.Z`

Expected success criteria:
- public site image pushed to `ghcr.io/<namespace>/aysteknoloji-com:X.Y.Z`
- release body shows:
  - compatible core tag visibility flag
  - expected core image path

Evidence to capture:
- workflow URL
- image path used
- release URL
- release note snippet for `Core Tag Visibility`

## Required Preconditions

Before starting the first real run:

1. Compatibility variables are configured intentionally:
   - `aystek-core` → `PANEL_COMPAT_TAG`
   - `admin-panel-next` → optional `CORE_COMPAT_TAG`
   - `aysteknoloji-com` → optional `CORE_COMPAT_TAG`
2. If custom namespace is desired, set the same `GHCR_NAMESPACE` variable in all repos.
3. If no custom namespace is used, default namespace is `github.repository_owner`.
4. `aystek-core-installer` workflow inputs point to the correct private repos.

## Execution Order

1. Run installer workflow from GitHub UI
2. Run panel release
3. Run core release
4. Run public site release
5. Verify registry pulls:

```bash
docker pull ghcr.io/<namespace>/aystek-panel:<panel-tag>
docker pull ghcr.io/<namespace>/aystek-core:<core-tag>
docker pull ghcr.io/<namespace>/aysteknoloji-com:<public-tag>
```

## Post-Run Validation Notes

Record the following after first successful runs:

- installer workflow run URL:
- panel workflow run URL:
- core workflow run URL:
- public site workflow run URL:
- namespace used:
- validated image coordinates:
  - panel:
  - core:
  - public site:
- local-source smoke result:
- panel core-tag visibility result:
- core panel-tag lock result:

## Next Documentation Update

After the first successful runs:
- update the dual-repo runbook with `Validated in CI`
- replace `<namespace>` placeholders with the real chosen namespace if fixed
- link the exact first successful run URLs for operator reference
