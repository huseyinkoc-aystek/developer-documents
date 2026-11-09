---
Title: Middleware Chain Invariants
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ./ENTERPRISE_CONTEXT_GUARANTEE.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md

---
# Table of Contents
1. [Overview](#overview)

# MIDDLEWARE CHAIN INVARIANTS

## Protected Chain Order

1. `AuthContextMiddleware`
2. `StoreContextMiddleware`
3. `ContextIntegrityValidator`
4. `TenantIsolation`
5. `SiteIsolation`
6. `RequestContextMiddleware`
7. `ImpersonationAuditMiddleware`
8. `ModuleGateByPathMiddleware`
9. `PermissionMiddleware`
10. `LanguageMiddleware`
11. `CSRFMiddleware`
12. `MaintenanceMiddleware`

## Invariants

1. Middleware order is immutable.
2. `StoreContextMiddleware` is the only resolver of raw tenant/site headers.
3. `ContextIntegrityValidator` detects post-resolution context mutation.
4. `TenantIsolation` enforces tenant presence and tenant mismatch rejection.
5. `SiteIsolation` enforces site requirements by route scope.
6. `RequestContextMiddleware` mirrors middleware context into request context.
7. No downstream module should read tenant/site headers directly.
8. Repository-wide tenant/site setters are restricted to approved files and enforced by test (PR gate).

## Environment Gates

| Gate | Dev | Production |
|---|---|---|
| Context integrity violation | panic | log only |
| Host fallback | enabled | disabled |
| Tenantfilter missing context | panic | error return |
