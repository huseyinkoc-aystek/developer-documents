# AYSTEK Core Private GHCR And Local Source Operator Runbook

Date: 2026-04-24
Status: Active
Scope: `aystek-core`, `aystek-core-next`, `aystek-core-installer`

## Purpose

Operational guidance for choosing between private GHCR images and local source directories during installation or update.

## Decision Path

### Use `docker + registry` when

- release images are already pushed to `ghcr.io/aystek`
- server should not build code locally
- registry credentials are managed safely

Required:
- `registry_host`
- `registry_username`
- `registry_token` or `registry_token_file`
- `core_image_tag`
- `panel_image_tag`

### Use `docker + local-source` when

- the target server already has trusted local clones of:
  - `aystek-core`
  - `aystek-core-next`
- private GHCR access is unavailable or intentionally avoided
- Docker is available and build should happen on-server

Required:
- `core_source_path`
- `panel_source_path`
- `core_image_tag`
- `panel_image_tag`

Prerequisites:
- `docker`
- `deployments/Dockerfile` exists under core source tree
- `Dockerfile` exists under panel source tree

### Use `native + local-artifacts` when

- release binary and panel bundle are prepared externally
- target server should not compile code

Required:
- `core_binary_path`
- `panel_build_path`

### Use `native + local-source` when

- target server has local source trees
- native/systemd deployment is preferred

Required:
- `core_source_path`
- `panel_source_path`

Prerequisites:
- `go`
- `npm`
- `go.mod` exists under core source tree
- `package.json` exists under panel source tree

## Release Coordination

- Core release workflow publishes:
  - `ghcr.io/aystek/aystek-core:<version>`
- Panel release workflow must publish:
  - `ghcr.io/aystek/aystek-panel:<version>`
- Installer config should use matching core/panel tags for the same release line
- Mixed tag pairs are not recommended unless explicitly tested

## Operator Checklist

1. Choose `source_mode` before filling env/service details.
2. If using local source:
   - verify source directories exist on the server
   - verify toolchain prerequisites exist
3. If using GHCR:
   - verify token file exists and is readable
   - verify both core and panel images exist for the same tag
4. Run installer validation after deploy.

## Recommended Validation Commands

```bash
cd /workspace/go/aystek-core-installer && make integration-source-modes
cd /workspace/go/aystek-core-installer && go test ./...
cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds
```
