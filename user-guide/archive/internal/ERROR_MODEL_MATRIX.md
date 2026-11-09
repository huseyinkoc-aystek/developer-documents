---
Title: Error Model Matrix
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ../04-api/ERROR_CODE_REFERENCE.md

---
# Table of Contents
1. [Overview](#overview)

# ERROR MODEL MATRIX — Tenant/Site Resolution

## Standard Error Shape
```
{
  "error": "SiteTenantMismatch",
  "code": "SITE_TENANT_MISMATCH",
  "message": "site does not belong to tenant"
}
```

## Status Code Mapping

| Scenario | HTTP | error | code | message |
|---|---|---|---|---|
| Missing tenant_id | 400 | TenantMissing | TENANT_MISSING | tenant_id required |
| Missing site_id (admin) | 400 | SiteMissing | SITE_MISSING | site_id required |
| Invalid tenant header | 400 | InvalidTenant | INVALID_TENANT | invalid tenant header |
| Invalid site id | 400 | InvalidSite | INVALID_SITE | invalid site id |
| Site belongs to another tenant | 403 | SiteTenantMismatch | SITE_TENANT_MISMATCH | site does not belong to tenant |
| Auth failure | 401 | (auth errors) | (auth codes) | unauthorized |

## Determinism Rules
- 400 = missing required context or invalid input.
- 403 = cross-tenant/site mismatch.
- 401 = auth failure.
- 500 = internal errors only.

