# Multi Tenant Bootstrap And Settings Namespace Closure Checkpoint

Date: 2026-04-10

## Completed In This Slice

- Tenant selector bootstrap endpoints were aligned with frontend context-guard rules.
- Tenant/site bootstrap mutations were marked as CSRF-skip requests on the Next client because the backend CSRF middleware is currently placeholder and the bootstrap flow must not depend on `/auth/csrf` prefetch.
- Bootstrap-related backend middleware exemptions were widened so these routes no longer require an already-selected site:
  - `/me/context`
  - `/me/select-tenant`
  - `/me/bootstrap-sites`
  - `/me/bootstrap-tenant`
  - `/me/bootstrap-site`
- Localhost store fallback now prefers `admin-panel` before `default`, which aligns local admin runtime with the seeded `core.aysteknolojileri.com` admin site contract.
- Login flow now resolves sites after tenant selection:
  - one site -> enter dashboard directly
  - zero/multiple sites -> route to `/tenant-selector?step=site`
- Tenant selector site bootstrap UI no longer leaves the create-site action disabled in the empty-state path.
- Tenant selector now prefills enterprise-safe bootstrap suggestions:
  - `default` tenant -> `AYSTEK Admin` / `admin-panel` / `core.aysteknolojileri.com`
  - `aystek` tenant -> `AYSTEK Mühendislik` / `aystek` / `aysteknolojileri.com`

## Validation

- Go:
  - `cd /workspace/go/aystek-core && go test ./internal/modules/system/... ./internal/platform/middleware/... ./internal/seeds ./cmd/api`
- Next:
  - `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/tenant-selector/page.tsx' --file 'app/(auth)/login/page.tsx' --file 'lib/api/me.ts' --file 'lib/api/requestContext.ts'`

## Settings Namespace Recheck

- Re-ran the `settings.*` namespace comparison with the tenant-site surface included.
- UI-used `settings.*` keys: `185`
- Missing seed keys after recheck: `0`

Artifacts:

- `/workspace/settings_ui_keys_recheck.txt`
- `/workspace/settings_seed_keys_recheck.txt`
- `/workspace/settings_missing_keys_recheck.txt`

## Important Runtime Note

The code path is now aligned, but a currently running dev instance may still need:

1. backend restart
2. seed reconcile pass
3. hard refresh in the browser

to fully reflect the fixed bootstrap and translation behavior.
