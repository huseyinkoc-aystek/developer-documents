# Public CMS Site Scope Hardening Slice 1 Checkpoint

Date: 2026-04-10

## Scope
- `pages`
- `menus`
- `sliders`
- `media`

## Completed
- Added `site_id` to the domain models for all four modules.
- Upgraded repository filters from tenant-only to `tenant + site`.
- Required `site_id` during create/update/delete/find flows, aligned with `tenantfilter.WithTenantAndSite(...)`.
- Added Mongo index registrations for the four modules:
  - tenant/site/status and tenant/site/created variants
  - tenant/site/key for menus
  - tenant/site/file_path for media
- Updated default content/media/slider seeders to write `site_id`.
- Updated AYSTEK public CMS media seed to upsert by `tenant_id + site_id + file_path`.
- Added `SeedPublicCMSSiteScope` to backfill legacy tenant-scoped CMS documents with the correct site id.
- Moved the backfill step earlier in the seed runner so upgrades do not create duplicate site-scoped records before reconciliation.

## Validation
- `go test ./internal/modules/pages/... ./internal/modules/menus/... ./internal/modules/sliders/... ./internal/modules/media/... ./internal/modules/stores/repository ./internal/platform/middleware ./internal/seeds ./cmd/api`

## Result
- Public/admin reads for these modules now isolate content by `tenant_id + site_id`.
- Existing seeded CMS content is no longer left behind as tenant-only ghost data after the repository hardening.

## Remaining Work
- Add targeted repository tests for cross-site isolation on the four hardened modules.
- Re-run a full seed/startup smoke flow against a disposable database and verify:
  - `default/admin-panel`
  - `aystek/aystek`
- Verify `/workspace/node/aysteknoloji-com` against the hardened API runtime once non-interactive lint/test ergonomics are fixed in that repo.
