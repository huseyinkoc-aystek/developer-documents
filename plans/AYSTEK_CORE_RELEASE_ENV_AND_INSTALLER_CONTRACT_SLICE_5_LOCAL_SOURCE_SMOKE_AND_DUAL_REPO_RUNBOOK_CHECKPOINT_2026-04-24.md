# AYSTEK Core Release Env And Installer Contract
## Slice 5 Checkpoint: Local Source Smoke And Dual Repo Runbook

Date: 2026-04-24
Status: Completed
Scope: `aystek-core`, `aystek-core-next`, `aystek-core-installer`

## Closed Items

- added local-source dry-run smoke harness
- added:
  - `IT-12` docker local-source prerequisite smoke
  - `IT-13` native local-source prerequisite smoke
- installer CI now validates local-source smoke against real sibling repo paths
- panel release workflow aligned to canonical GHCR namespace
- dual-repo release runbook added for coordinated core/panel tag publishing

## Verification

```bash
cd /workspace/go/aystek-core-installer && make integration-source-modes
cd /workspace/go/aystek-core-installer && make integration-local-source-smoke
cd /workspace/go/aystek-core-installer && go test ./...
cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds
```

## Outcome

Local-source is now covered at two levels:
- config contract
- deployer prerequisite smoke against real repo paths

Core and panel releases are now documented as one coordinated operator workflow rather than two unrelated image pushes.
