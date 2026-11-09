# Settings I18n Cache Fix And Tenant Site Public Branding Checkpoint

Date: 2026-04-06

## Closed

- Fixed the admin translation cache overwrite bug on Next.js.
  - Root cause: admin translation dictionaries were being replaced by narrower `publicOnly` dictionaries during client hydration.
  - Fix: i18n cache entries now merge instead of overwrite in `lib/stores/i18n.ts`.

- Added backend-managed tenant-site public branding settings.
  - New admin-protected endpoint:
    - `GET /admin/settings/tenant-site/public-branding`
    - `PUT /admin/settings/tenant-site/public-branding`
  - New permission:
    - `settings:tenant_site_public_branding`
  - Admin role receives this automatically because role seed force-syncs all system permissions.

- Added separate public branding behavior for public pages.
  - Public settings now prefer:
    1. `site.<site_id>.public.logo_url` / `favicon_url`
    2. `tenant.<tenant_id>.public.logo_url` / `favicon_url`
    3. `public.logo_url` / `favicon_url`
    4. legacy admin branding fallback keys

- Added Settings admin UI support for:
  - tenant-site public logo
  - tenant-site public favicon
  - scope badge
  - dedicated save action

## Files

- Go:
  - `internal/modules/settings/service/interface.go`
  - `internal/modules/settings/service/service.go`
  - `internal/modules/settings/service/service_test.go`
  - `internal/modules/settings/dto/dto.go`
  - `internal/modules/settings/handler/admin.go`
  - `internal/modules/settings/routes.go`
  - `internal/seeds/permission_seed.go`
  - `internal/seeds/translation_seed.go`

- Next:
  - `lib/stores/i18n.ts`
  - `lib/api/settings.ts`
  - `lib/api/public-settings.ts`
  - `lib/types.ts`
  - `app/(dashboard)/admin/settings/page.tsx`

## Verification

- Passed:
  - `cd /workspace/go/aystek-core && go test ./internal/modules/settings/... ./internal/seeds ./cmd/api`

- Partial:
  - Targeted Next lint on settings pages still reports pre-existing file-level lint debt in `app/(dashboard)/admin/settings/page.tsx` and broad historical `lib/types.ts` debt.
  - New runtime-facing files (`lib/api/settings.ts`, `lib/api/public-settings.ts`, `lib/stores/i18n.ts`) are structurally updated, but full page-file lint cleanliness was not achieved because of existing violations outside this slice's core changes.

## Residual

- If needed, the next cleanup slice should be:
  - `settings page lint debt reduction`
  - `settings dashboard locale-aware fallback tightening`
  - `public branding preview consumption on public shell/favicon head`
