# Multi Tenant Login, Tenant Selector, and Runtime Verification Checkpoint

Date: 2026-04-11

## Closed in this slice

- Auth layout no longer forces `publicOnly` translation preload.
- Server-side translation preload now prefers canonical admin store resolution on localhost auth/admin surfaces.
- Client axios now injects canonical admin store context for localhost auth/bootstrap/public/languages requests when no explicit site context exists.
- Bootstrap mutations no longer attempt CSRF prefetch on bootstrap paths.
- Backend `StoreContextMiddleware` bug was fixed:
  - in production mode it previously returned `HostFallbackDisabled` before honoring `X-Store-Code`
  - it now blocks host fallback only when both `X-Store-Code` and `X-Store-Domain` are absent
- Added regression test for explicit store-code resolution in production middleware path.

## Key root cause found

The tenant/site bootstrap failures were not only a frontend problem.

There was a backend middleware defect:

- requests with `X-Store-Code: admin-panel`
- and no `X-Store-Domain`
- were still treated as a host-fallback attempt
- and failed with `HostFallbackDisabled`

This explains:

- login failing on localhost
- public translation fetches showing raw keys
- tenant selection persistence failures
- site creation flow not completing reliably

## Runtime note

The actual listener on `:9090` was not the temporary `go run ./cmd/api` process.
It was an existing compiled binary:

- `/workspace/go/aystek-core/api`

That binary was restarted with the updated code.

## Verification completed

- `go test ./internal/platform/middleware ./internal/modules/system/... ./internal/seeds ./cmd/api`
- Public CMS verification by explicit domain header still returns seeded public menu data for `aysteknolojileri.com`

## Verification still in progress

Full runtime smoke against the restarted backend is waiting for the current seed/startup cycle to complete and bind `:9090`.

Pending live checks:

- login with admin user
- `/me/context`
- tenant select
- site list for `default`
- site list for `aystek`
- public site spot-checks

