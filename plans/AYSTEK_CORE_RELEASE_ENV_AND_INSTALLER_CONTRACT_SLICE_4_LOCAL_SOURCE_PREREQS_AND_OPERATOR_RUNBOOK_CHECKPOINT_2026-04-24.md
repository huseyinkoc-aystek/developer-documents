# AYSTEK Core Release Env And Installer Contract
## Slice 4 Checkpoint: Local Source Prereqs And Operator Runbook

Date: 2026-04-24
Status: Completed
Scope: `aystek-core-installer`, docs/runbooks

## Closed Items

- local-source prerequisite checks are now enforced in deployers
  - native:
    - `go`
    - `npm`
    - `go.mod`
    - `package.json`
  - docker:
    - `docker`
    - `deployments/Dockerfile`
    - `Dockerfile`
    - output image tags
- source-mode contract matrix now includes `tenant-site + native + local-source`
- operator runbook added for choosing between:
  - `docker + registry`
  - `docker + local-source`
  - `native + local-artifacts`
  - `native + local-source`
- installer usage guide and server setup guide updated with local-source prerequisite rules

## Verification

```bash
cd /workspace/go/aystek-core-installer && make integration-source-modes
cd /workspace/go/aystek-core-installer && go test ./...
cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds
```

## Outcome

Installer no longer treats local-source as a soft convention. It is now a verified deployment mode with explicit prerequisite enforcement and operator-facing runbook guidance.
