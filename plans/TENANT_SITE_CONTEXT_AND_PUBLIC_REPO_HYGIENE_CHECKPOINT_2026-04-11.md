# Tenant/Site Context and Public Repo Hygiene Checkpoint — 2026-04-11

## Fixed
- Login and tenant selector no longer depend on a false empty-site state caused by frontend parsing of `/me/bootstrap-sites`.
- Site context loading now uses bootstrap-authorized site lists instead of broad admin store listing.
- Logout-time language loading is hardened against 401/403 and falls back safely.
- Header now exposes a direct return path to the full tenant/site selection screen.
- Turkish is now the default i18n fallback baseline instead of English.

## Verified
- `go test ./internal/seeds ./cmd/api` passes.
- Public site repo now has a committed `.eslintrc.json` baseline.
- `aysteknoloji-com` production build succeeds.

## Open
- `aystek-core-next` still has pre-existing lint debt in `components/layout/Header.tsx` and adjacent shared files.
- `aysteknoloji-com` lint must now be read after the first committed ESLint baseline; next step is to fix reported issues instead of interactive setup.
- Full tenant/site authorization model (which users can access which tenant/sites) still needs a dedicated feature slice; current work improves selection and context flow, not the full authorization domain.
