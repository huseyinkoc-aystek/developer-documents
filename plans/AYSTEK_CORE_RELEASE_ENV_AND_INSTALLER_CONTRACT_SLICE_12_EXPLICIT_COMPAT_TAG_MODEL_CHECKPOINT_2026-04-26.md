# AYSTEK Core Release Env And Installer Contract
## Slice 12 Checkpoint: Explicit Compatibility Tag Model

Date: 2026-04-26

### Problem
- `aystek-core` and `aystek-core-next` were incorrectly treated as if they shared the same release version line.
- This caused a false workflow design:
  - core release expected panel tag `v2.0.4`
  - while panel release legitimately used `v1.0.5`

### Completed
- `aystek-core` release workflow now uses explicit compatibility tag variable:
  - `PANEL_COMPAT_TAG`
  - `/workspace/go/aystek-core/.github/workflows/release.yml`
- `aystek-core-next` release workflow now uses optional visibility variable:
  - `CORE_COMPAT_TAG`
  - `/workspace/node/aystek-core-next/.github/workflows/release.yml`
- Dual-repo runbook now describes independent versioning and compatibility mapping:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Result
- Core and panel no longer need matching semantic version tags.
- Compatibility is operator-controlled and explicit.
- Release coordination remains strict without incorrectly coupling unrelated version lines.
