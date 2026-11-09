# AYSTEK Core Release Env And Installer Contract

Date: 2026-04-24
Status: Slice 2 Completed

## Completed In This Slice

- installer config contract expanded with first-class `source_mode`
- fresh install and update wizard now support:
  - `docker + registry`
  - `docker + local-source`
  - `native + local-artifacts`
  - `native + local-source`
- docker deployer can build local images from:
  - local `aystek-core` repo
  - local `aystek-core-next` repo
- native deployer can build:
  - `aystek-core` binary from local source
  - `aystek-core-next` standalone bundle from local source
- post-install validation for external MongoDB / Redis / RabbitMQ is stricter:
  - URI/address parsing
  - real TCP reachability check instead of unconditional skip
- operator docs updated to explain local-source deployment contract

## Main Files

- `/workspace/go/aystek-core-installer/internal/config/schema.go`
- `/workspace/go/aystek-core-installer/internal/ui/wizard.go`
- `/workspace/go/aystek-core-installer/internal/deploy/native.go`
- `/workspace/go/aystek-core-installer/internal/deploy/docker.go`
- `/workspace/go/aystek-core-installer/internal/validate/health.go`
- `/workspace/go/aystek-core-installer/internal/i18n/i18n.go`
- `/workspace/go/aystek-core-installer/internal/i18n/en.go`
- `/workspace/go/aystek-core-installer/internal/i18n/tr.go`
- `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`

## Verification

- `cd /workspace/go/aystek-core-installer && go test ./...`
- `cd /workspace/go/aystek-core-installer && go build ./cmd/aystek-core-installer`
- `cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'lib/api/baseUrl.ts' --file 'lib/api/translations.server.ts'`

## Operational Notes

- `docker + local-source` no longer requires GHCR credentials
- `native + local-source` requires `go`, `node`, and `npm` on the target server
- custom host ports remain valid; external service health checks now validate configured endpoints rather than silently skipping them
