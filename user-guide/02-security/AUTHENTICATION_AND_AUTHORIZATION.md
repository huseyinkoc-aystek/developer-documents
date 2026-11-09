---
Title: Authentication and Authorization Guide
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ./TENANT_SITE_SECURITY_MATRIX.md
  - ../04-api/API_REFERENCE.md
  - ../04-api/ERROR_CODE_REFERENCE.md
  - ../01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md
---

# Table of Contents
1. [Authentication](#1-authentication)
2. [Authorization](#2-authorization)
3. [Context Contract](#3-context-contract)
4. [Injection Code Security Model](#4-injection-code-security-model)
5. [Connector Secret Handling](#5-connector-secret-handling)
6. [RBAC Enforcement](#6-rbac-enforcement)
7. [Multi-Tenant Isolation Model](#7-multi-tenant-isolation-model)
8. [Failure Model](#8-failure-model)

## 1. Authentication
Supported auth surfaces:
- `/auth/*` bootstrap and token lifecycle endpoints
- Bearer token (`Authorization`) for protected routes
- Session/cookie flows where configured by runtime

## 2. Authorization
Authorization is backend-enforced through:
1. Tenant/site isolation checks
2. ModuleGate checks
3. RBAC permission checks
4. License and feature-flag constraints (where applicable)

## 3. Context Contract
Headers used for context resolution:
- `X-Tenant-ID` or `X-Tenant-Code`
- `X-Site-ID` (or store headers where allowed by runtime policy)

Scope rules are authoritative in `./TENANT_SITE_SECURITY_MATRIX.md`.

## 4. Injection Code Security Model
Injection code registry is governed as privileged operational configuration, not generic CMS content.

### 4.1 Controls
- Sanitization on write path (allowlist enforcement)
- Byte-size ceiling per site (`content limits`)
- Versioning for deterministic rollback and traceability
- Confirm-before-save UX in admin panel
- Soft delete for lifecycle-safe deactivation

### 4.2 Exposure Rules
- Injection code is only accessible through authenticated admin endpoints with site context.
- No public or SSR-side exposure of raw code outside explicit admin form contexts.

## 5. Connector Secret Handling
Connector credentials are managed as secret references.

### 5.1 Secret Storage Policy
- Stored encrypted at rest
- Not persisted as plaintext in connector config
- Not part of generic `site_settings`

### 5.2 Rotation
- Rotation endpoint supported for credential lifecycle changes
- Audit entries required for secret change operations
- Secrets must be redacted from logs/events

## 6. RBAC Enforcement
Minimum critical permissions for this platform slice:
- `ingestion.read`
- `connector.manage`
- `moderation.approve`
- `content_ops_settings.manage`

Enforcement occurs on server handlers and service-level policy checks.

## 7. Multi-Tenant Isolation Model
- Tenant and site are derived from auth/context middleware.
- Request body cannot override tenant/site scope.
- Every repository query includes tenant/site scoping and soft-delete predicates.

## 8. Failure Model
Canonical error payload:
```json
{
  "error": "SiteMissing",
  "message": "site context required for this route",
  "trace_id": "req_123"
}
```
