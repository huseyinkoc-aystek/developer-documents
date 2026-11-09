---
Title: Enterprise Context Guarantee
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ./MIDDLEWARE_CHAIN_INVARIANTS.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
  - ../04-api/ERROR_CODE_REFERENCE.md

---
# Table of Contents
1. [Overview](#overview)

# ENTERPRISE CONTEXT GUARANTEE

## Scope

This document defines runtime guarantees for tenant/site context enforcement in AYSTEK Core v1.0.

## Guarantees

1. `tenant_id` is required for all routes except `/healthz`.
2. `site_id` is required for `/admin/*` and `/graphql`.
3. Public routes (`/public/*`, `/languages`, `/auth/*`) are tenant-scoped and site-optional.
4. `site_id` must belong to `tenant_id`; mismatches are rejected with `403`.
5. Auth/header tenant tampering is rejected with `403 CrossTenantAccessAttempt`.
6. Context values are resolved in middleware and propagated to request context.
7. Downstream layers must use request context values, not direct headers.

## Production Safety

1. Host-based fallback resolution is disabled in production.
2. Any host fallback attempt in production emits `HostFallbackDisabled`.
3. Resolution errors return canonical JSON format:

```json
{
  "error": "<ErrorCode>",
  "message": "<Human readable message>",
  "trace_id": "<request id>"
}
```

## Dev vs Prod

| Rule | Dev | Production |
|---|---|---|
| Context invariant violation | panic | structured error log |
| Missing tenant in tenant-scoped operation | panic (tenantfilter) | error return |
| Missing site in site-scoped operation | panic (tenantfilter) | error return |
| Host fallback | allowed | disabled |
