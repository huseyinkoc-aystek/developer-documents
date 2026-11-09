# Tenant Site Authorization Slice 5 Checkpoint

Date: 2026-04-12

## Scope Closed

This slice made tenant/site authorization more operational on two fronts:

- users list now exposes a faster matrix-style access view for selected users
- header and context switching are hardened against out-of-scope tenant/site selections

## Frontend Changes

### Users list matrix view

The users list page now includes:

- bulk access panel with selected-user matrix summary
- per-user tenant/site access matrix cards
- tenant/site labels resolved from live context metadata instead of raw IDs where possible

Files:

- `/workspace/node/aystek-core-next/app/(dashboard)/admin/users/page.tsx`

### Context payload typing

The authenticated context payload now carries typed user memberships on the frontend.

Files:

- `/workspace/node/aystek-core-next/lib/api/me.ts`
- `/workspace/node/aystek-core-next/contexts/TenantContext.tsx`
- `/workspace/node/aystek-core-next/contexts/SiteContext.tsx`

### Context switch hardening

- tenant switch options are now aligned with membership-scoped context
- site switch no longer accepts arbitrary IDs outside known options
- site refresh is filtered against tenant membership scope before rendering options

Files:

- `/workspace/node/aystek-core-next/components/context/TenantSwitcher.tsx`
- `/workspace/node/aystek-core-next/components/context/SiteSwitcher.tsx`
- `/workspace/node/aystek-core-next/lib/tenant.ts`
- `/workspace/node/aystek-core-next/middleware.ts`

## Backend Changes

The store-context middleware now enforces tenant/site selections against effective user memberships.

- non-platform users can switch to selected tenants only if they are in effective memberships
- resolved site context is rejected if it falls outside the user’s allowed sites under the active tenant
- legacy non-user/system tests still preserve previous permissive behavior

Files:

- `/workspace/go/aystek-core/internal/platform/middleware/store_context.go`
- `/workspace/go/aystek-core/internal/platform/middleware/resolution_errors.go`

## Verification

Passed:

```bash
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/users/page.tsx' --file 'components/context/TenantSwitcher.tsx' --file 'components/context/SiteSwitcher.tsx' --file 'contexts/TenantContext.tsx' --file 'contexts/SiteContext.tsx' --file 'lib/api/me.ts' --file 'lib/tenant.ts' --file 'middleware.ts'
cd /workspace/go/aystek-core && go test ./internal/platform/middleware ./internal/modules/system/... ./internal/modules/users/... ./internal/seeds ./cmd/api
```

## Current Status

Closed:

- users list operational matrix visibility
- membership-aware header context boundary enforcement

Next:

- extend matrix editing into a denser permission-management surface
- add explicit platform-admin vs tenant-admin UI affordances
- bind dashboard module visibility more directly to tenant/site membership scope
