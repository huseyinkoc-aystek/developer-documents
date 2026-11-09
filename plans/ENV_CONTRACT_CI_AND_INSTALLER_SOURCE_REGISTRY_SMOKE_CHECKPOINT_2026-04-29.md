# Env Contract CI and Installer Source/Registry Smoke Checkpoint - 2026-04-29

## Scope

This checkpoint closes the first automation pass for environment contract drift and installer deployment source validation across:

- `aystek-core`
- `aystek-core-next`
- `aysteknoloji-com`
- `aystek-core-installer`

## Completed

- Added env contract checks for backend env examples.
- Added backend Go code env scanning. A key used in Go code must now be present in `.env.example` or explicitly listed in `scripts/env-contract.optional.keys`.
- Added env contract checks for admin panel `.env.example` against actual `process.env.*` usage.
- Added env contract checks for public site `.env.example` and `.env.local.example` against actual `process.env.*` usage.
- Added installer env template contract checks for generated `core.env`, `panel.env`, and Docker stack `.env` templates.
- Added CI wiring so env contract drift fails early.
- Added release workflow env contract checks before GHCR package/image publishing.
- Added installer registry-mode dry-run smoke.
- Kept existing local-source smoke in place for server-local source directory installs.
- Rebuilt and verified installer release binaries for `linux/amd64` and `linux/arm64`.

## Installer Deployment Modes Covered

- `docker + registry`
- `docker + local-source`
- `native + local-source`
- `native + local-artifacts`

## Important Behavior

- If an operator enters a server-local source directory path, installer treats it as `local-source` and does not ask for registry credentials.
- If an operator configures GHCR/private image tags, installer treats it as registry mode and validates registry credentials/token path in dry-run smoke.
- `aystek-core-installer-linux-amd64` is now built through `make release-binaries` and verified as an x86-64 ELF.

## Verification

- `cd /workspace/go/aystek-core && bash scripts/check-env-contract.sh && go test ./cmd/api ./internal/platform/config ./internal/seeds`
- `cd /workspace/node/aystek-core-next && npm run check-env`
- `cd /workspace/node/aysteknoloji-com && npm run check-env`
- `cd /workspace/go/aystek-core-installer && make env-contract integration-source-modes integration-local-source-smoke integration-registry-smoke`
- `cd /workspace/go/aystek-core-installer && make verify-binaries`
- `cd /workspace/go/aystek-core-installer && ./bin/aystek-core-installer-linux-amd64 version`
- `cd /workspace/go/aystek-core-installer && readelf -h bin/aystek-core-installer-linux-amd64`
