---
Title: Tenant Site Security Matrix
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md
  - ../archive/internal/MIDDLEWARE_CHAIN_INVARIANTS.md
  - ../04-api/ERROR_CODE_REFERENCE.md

---
# Table of Contents
1. [Overview](#overview)

# TENANT SITE SECURITY MATRIX

## Route Scope Matrix

| Route Type | Tenant | Site | Result if Missing Tenant | Result if Missing Site |
|---|---|---|---|---|
| `/healthz` | Optional | Optional | Allowed | Allowed |
| `/public/*` | Required | Optional | `400 TenantMissing` | Allowed |
| `/languages` | Required | Optional | `400 TenantMissing` | Allowed |
| `/auth/*` | Required | Optional | `400 TenantMissing` | Allowed |
| `/admin/*` | Required | Required | `400 TenantMissing` | `400 SiteMissing` |
| `/graphql` | Required | Required | `400 TenantMissing` | `400 SiteMissing` |

## Public Endpoint Determinism

| Endpoint | Tenant | Site | Expected |
|---|---|---|---|
| `/public/translations` | Required | Optional | `200` with tenant-only context |
| `/public/settings` | Required | Optional | `200` with tenant-only context |

## High-Risk Mutation Coverage

| Endpoint | Tenant | Site | Expected |
|---|---|---|---|
| `PUT /admin/roles/:id` | Required | Required | `400 SiteMissing` without site, `200` with tenant+site |
| `PUT /admin/settings/:key` | Required | Required | `400 SiteMissing` without site, `403 SiteTenantMismatch` on tenant/site mismatch |

## Cross-Isolation Enforcement

| Condition | Status | Error |
|---|---|---|
| Site belongs to different tenant | `403` | `SiteTenantMismatch` |
| Auth tenant differs from resolved tenant | `403` | `CrossTenantAccessAttempt` |
| Host fallback attempted in production | `400` | `HostFallbackDisabled` |

## Header Contract

| Header | Required | Max Length | Notes |
|---|---|---|---|
| `X-Tenant-ID` | tenant-scoped routes | 128 | primary tenant resolution key |
| `X-Tenant-Code` | optional alt | 128 | secondary tenant resolution key |
| `X-Site-ID` | admin/graphql | 128 | primary site resolution key |
| `X-Store-Code` | optional | 128 | secondary site/store resolution key |
| `X-Store-Domain` | optional | 255 | secondary site/store resolution key |
