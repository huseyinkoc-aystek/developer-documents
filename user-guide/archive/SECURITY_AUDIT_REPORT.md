---
Title: Security Audit Report
Version: 2.0.0
Status: Archive
Last Updated: 2026-02-12
Related Documents:
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md

---
# Table of Contents
1. [Overview](#overview)

# SECURITY AUDIT REPORT — Tenant/Site Hardening

Scope:
- internal/platform/middleware/*
- internal/modules/**/service/
- internal/modules/**/repository/
- internal/modules/**/handler/
- internal/platform/request_context
- internal/platform/module_gate
- internal/platform/permission

## Summary
This audit focuses on multi-tenant and multi-site isolation risks, header usage, tenant/site resolution, and repository tenant scoping. Findings below are categorized by risk and include exploit scenario and recommended fix.

## Findings

| Risk Category | File | Function | Severity | Exploit Scenario | Recommended Fix |
|---|---|---|---|---|---|
| Header Access Outside Middleware | internal/modules/** | various | HIGH | Downstream code reading tenant/site headers can bypass deterministic resolution, leading to inconsistent enforcement. | Enforce “headers only in middleware” with static audit test and move any header access into StoreContextMiddleware. |
| Host Fallback in Production | internal/platform/middleware/store_context.go | resolveStore | HIGH | Implicit host-based resolution can route a request to an unintended tenant/site in production. | Disable host fallback outside dev; require explicit headers or store lookup. |
| Cross-Tenant Site Resolution | internal/platform/middleware/store_context.go | ResolveSite | CRITICAL | A site resolved via X-Site-ID or domain could belong to another tenant, enabling cross-tenant access. | Validate store.TenantID vs tenant_id; reject mismatches with 403. |
| Public Route Site Requirement | internal/platform/middleware/site_isolation.go | SiteIsolation | MEDIUM | Public routes fail when site_id missing, causing availability issues and forcing unsafe fallback. | Allow site optional for public routes while keeping tenant required. |
| Tenant Missing Errors | internal/platform/middleware/gateway.go | TenantIsolation | HIGH | Generic 400 errors hide root cause and reduce auditability. | Use structured error codes and deterministic messages. |
| Repository Tenant Filter Omission | internal/modules/**/repository | various | CRITICAL | Missing tenant filter leads to cross-tenant read/write. | Add tenantfilter.WithTenantID/WithTenantAndSite or explicit tenant_id filters; enforce via audit tests. |
| Site Assumption in Services | internal/modules/**/service | various | MEDIUM | Services assuming site_id for all routes can panic on public requests. | Add RequireSite for admin paths; allow empty site for public flows. |
| ModuleGate Assumes Site | internal/platform/middleware/module_gate.go | moduleEnabledWithScope | MEDIUM | Site missing could cause unintended behavior. | Allow site optional on public requests; ensure tenant still used. |
| Permission Middleware Site Assumption | internal/platform/middleware/permission.go | PermissionMiddleware | MEDIUM | Public requests could fail if permission logic assumes site. | Ensure permission checks are tenant-scoped and tolerant of empty site for public routes. |

## Residual Risks
- Some repositories may still lack explicit tenant filter enforcement; audit test is a guard but not a substitute for manual verification.
- Public endpoints still require tenant; missing tenant will return 400 with deterministic error code.

