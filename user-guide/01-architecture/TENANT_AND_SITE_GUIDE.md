---
Title: Tenant and Site Context Guide
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ./Enterprise_Isolation_Whitepaper.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
  - ../04-api/ERROR_CODE_REFERENCE.md
---

# Table of Contents
1. [Concept Model](#1-concept-model)
2. [Resolution Inputs](#2-resolution-inputs)
3. [Route Scope Rules](#3-route-scope-rules)
4. [Deterministic Rejection Cases](#4-deterministic-rejection-cases)

## 1. Concept Model
- `tenant`: customer organization boundary.
- `site`: domain/runtime boundary within a tenant.
- Modules and capabilities are evaluated in tenant+site context.

Detailed architecture rationale lives in `./Enterprise_Isolation_Whitepaper.md`.

## 2. Resolution Inputs
Tenant resolution candidates:
- `X-Tenant-ID`
- `X-Tenant-Code`

Site resolution candidates:
- `X-Site-ID`
- `X-Store-Code`
- `X-Store-Domain`
- Host-based fallback only where runtime policy permits

## 3. Route Scope Rules
| Route Type | Tenant | Site |
|---|---|---|
| `/healthz` | Optional | Optional |
| `/public/*`, `/languages`, `/auth/*` | Required | Optional |
| `/admin/*`, `/graphql` | Required | Required |

## 4. Deterministic Rejection Cases
- Missing required tenant: `TenantMissing` (`400`)
- Missing required site: `SiteMissing` (`400`)
- Site does not belong to tenant: `SiteTenantMismatch` (`403`)
- Tenant tampering attempt: `CrossTenantAccessAttempt` (`403`)

Canonical error schema:
```json
{
  "error": "SiteTenantMismatch",
  "message": "site does not belong to tenant",
  "trace_id": "req_123"
}
```
