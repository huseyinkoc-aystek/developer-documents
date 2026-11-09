---
Title: Error Code Reference
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md

---
# Table of Contents
1. [Canonical Payload](#1-canonical-payload)
2. [Context Error Codes](#2-context-error-codes)
3. [Security Event Labels](#3-security-event-labels)
4. [Runtime Emission Guarantee](#4-runtime-emission-guarantee)

## 1. Canonical Payload

```json
{
  "error": "<ErrorCode>",
  "message": "<Human readable message>",
  "trace_id": "<request id>"
}
```

## 2. Context Error Codes

| Error | HTTP | Description |
|---|---|---|
| `TenantMissing` | 400 | tenant context is required but missing |
| `SiteMissing` | 400 | site context is required but missing |
| `SiteTenantMismatch` | 403 | site does not belong to resolved tenant |
| `InvalidTenant` | 400 | tenant header/context is invalid |
| `InvalidSite` | 400 | site/store header/context is invalid |
| `CrossTenantAccessAttempt` | 403 | authenticated tenant differs from resolved tenant |
| `HostFallbackDisabled` | 400 | host-based fallback is disabled in production |

## 3. Security Event Labels

1. `TenantMissing`
2. `SiteMissing`
3. `SiteTenantMismatch`
4. `InvalidTenant`
5. `InvalidSite`
6. `CrossTenantAccessAttempt`
7. `CrossSiteAccessAttempt`
8. `HostFallbackDisabled`
9. `ContextInvariantViolation`
10. `ContextMutationDetected`
11. `ContextResolutionTimeout`

## 4. Runtime Emission Guarantee

Security event emission is covered by middleware tests and includes timeout/mutation branches.  
Test file: `internal/platform/middleware/security_events_test.go`
