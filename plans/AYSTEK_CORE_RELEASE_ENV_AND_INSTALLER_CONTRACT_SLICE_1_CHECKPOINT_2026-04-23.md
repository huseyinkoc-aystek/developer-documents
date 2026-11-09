# AYSTEK Core Release Env And Installer Contract - Slice 1 Checkpoint

Date: 2026-04-23
Status: Completed

## Completed

### Canonical panel env contract aligned
The following keys are now emitted and documented consistently for panel runtime:
- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_API_URL`
- `NEXT_PUBLIC_SSE_BASE_URL`
- `NEXT_PUBLIC_APP_ENV`
- `NEXT_PUBLIC_SITE_URL`
- `API_URL`
- `INTERNAL_API_URL`

Updated locations:
- `aystek-core-installer/assets/templates/env/panel.env.tmpl`
- `aystek-core-installer/templates/env/panel.env.tmpl`
- `aystek-core/deployments/production/.env.example`
- `aystek-core/deployments/production/docker-compose.yml`
- `aystek-core/config/env/panel.env.example`
- `aystek-core-next/.env.example`
- `aystek-core-next/Dockerfile`
- `aystek-core-next/README.md`

### Core env contract tightened
- installer `core.env` now emits both `MONGO_URI` and `MONGODB_URI`
- local Mongo authSource fixed to `admin`
- `MONGODB_DB`, `FRONTEND_URL`, `ALLOWED_ORIGINS`, `CORS_ALLOWED_ORIGINS`, `APP_HOST`, `SERVER_HOST` added to generated env

Updated locations:
- `aystek-core-installer/assets/templates/env/core.env.tmpl`
- `aystek-core-installer/templates/env/core.env.tmpl`

### GHCR image defaults unified
Installer defaults and update wizard defaults were aligned to:
- `ghcr.io/aystek/aystek-core`
- `ghcr.io/aystek/aystek-panel`

Updated locations:
- `aystek-core-installer/internal/config/schema.go`
- `aystek-core-installer/internal/ui/wizard.go`

### Duplicate installer templates cleaned
`assets/templates/env/*` and `templates/env/*` were drifting. They are now synchronized again.

### Docs synchronized
Updated:
- `go/aystek-core/deployments/production/README.md`
- `go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`
- `aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`
- `aystek-core-dev-docs/user-guide/05-operations/docker.env.example`
- `aystek-core-dev-docs/user-guide/05-operations/docker-compose.aystek-core.example.yml`

## Validation
- `cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds`
- `cd /workspace/go/aystek-core-installer && go test ./...`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'lib/api/baseUrl.ts' --file 'lib/api/translations.server.ts'`

## Remaining slices
- validate native-mode external/local service combinations more deeply
- review if image naming should remain `aystek-panel` or move to a more explicit `aystek-core-next` release artifact label without breaking installer compatibility
- add release packaging docs/checklist for private GHCR build/push workflow
