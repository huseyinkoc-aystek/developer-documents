# AYSTEK Core Private GHCR Release And Source Mode Checklist

Date: 2026-04-24
Status: Active
Scope: `aystek-core`, `aystek-core-installer`

## Objective

Make private GHCR release packaging and installer source-mode combinations operationally explicit and verifiable.

## Canonical Release Contract

- Core image: `ghcr.io/aystek/aystek-core:<version>`
- Panel image: `ghcr.io/aystek/aystek-panel:<version>`
- Native release archive:
  - `aystek-core-<version>-linux-amd64.tar.gz`
- Installer source modes:
  - `docker + registry`
  - `docker + local-source`
  - `native + local-artifacts`
  - `native + local-source`

## Pre-Release Checklist

1. `aystek-core` release workflow points to canonical GHCR namespace.
2. Installer defaults and docs use the same image coordinates.
3. `make integration-source-modes` passes.
4. Installer P0 disposable matrix passes.
5. Deployment docs mention local-source as a first-class alternative to private GHCR.
6. Registry auth is skipped for local-source mode.
7. Native local-source prerequisites are documented:
   - `go`
   - `node`
   - `npm`
8. Docker local-source prerequisites are documented:
   - `docker`
   - build access to local `aystek-core` and `aystek-core-next` directories

## Operator Decision Tree

- If private GHCR is available and credentials are managed:
  - prefer `docker + registry`
- If server already contains trusted local source trees:
  - use `docker + local-source` or `native + local-source`
- If release artifacts are staged manually:
  - use `native + local-artifacts`

## Validation Commands

```bash
cd /workspace/go/aystek-core-installer && make integration-source-modes
cd /workspace/go/aystek-core-installer && make integration-p0-ci
cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds
```
