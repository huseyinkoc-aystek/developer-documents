---
Title: Request Context Contract
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ./ENTERPRISE_CONTEXT_GUARANTEE.md
  - ./MIDDLEWARE_CHAIN_INVARIANTS.md

---
# Table of Contents
1. [Overview](#overview)

# REQUEST CONTEXT CONTRACT — Tenant/Site Invariants

## Invariant Summary
- `tenant_id` is REQUIRED for all routes except `/healthz`.
- `site_id` is REQUIRED for admin routes.
- `site_id` is OPTIONAL for public routes (`/public/*`, `/languages`, `/auth/*`, `/healthz`).
- Downstream services and repositories must read tenant/site ONLY from RequestContext.
- Headers must be read ONLY in middleware (StoreContext).

## Context Guarantees
- `RequestContextMiddleware` injects tenant/site into context keys.
- `contextkeys.IsPublic(ctx)` indicates whether the request is public.

## Required Helper Functions
- `tenantfilter.RequireTenant(ctx)`
- `tenantfilter.RequireSite(ctx)`
- `tenantfilter.IsPublic(ctx)`

## Dev-Mode Assertions
- Missing tenant/site in non-public requests will panic in dev mode.
- Public requests allow missing site but still enforce tenant.

## Enforcement Tests
- `TestRepositoryRequiresTenant`
- `TestAdminServiceRequiresSite`
- `TestPublicServiceAllowsEmptySite`
- `TestAllRepositoriesIncludeTenantFilter`
- `TestNoHeaderAccessOutsideMiddleware`

