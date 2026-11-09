# AYSTEK Core Release Env And Installer Contract
## Slice 3 Checkpoint: GHCR Workflow And Source Mode Matrix

Date: 2026-04-24
Status: Completed
Scope: `aystek-core`, `aystek-core-installer`

## Closed Items

- `aystek-core` release workflow now uses canonical GHCR coordinates:
  - `ghcr.io/aystek/aystek-core`
  - `ghcr.io/aystek/aystek-panel`
- installer CI now validates source-mode contract fixtures before disposable P0 execution
- dedicated `integration-source-modes` target added to installer Makefile
- source-mode contract runner no longer depends on a real `/run/secrets/ghcr-token` file
- private GHCR vs local-source operator checklist added

## Verification

```bash
cd /workspace/go/aystek-core-installer && make integration-source-modes
cd /workspace/go/aystek-core-installer && go test ./...
cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds
```

## Outcome

The installer contract is now explicit for:
- `docker + registry`
- `docker + local-source`
- `native + local-artifacts`
- `native + local-source`

Private GHCR is now a supported first-class release path, but not a mandatory deployment dependency. Local source directories remain a canonical operator path and are covered by contract verification.
