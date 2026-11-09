---
Title: Repository Allowlist Rationale
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ./MIDDLEWARE_CHAIN_INVARIANTS.md

---
# Table of Contents
1. [Overview](#overview)

# Repository Tenant Scope Allowlist Rationale

## Scope

This document explains explicit allowlist exceptions used by repository tenant/site audit guards.

## Allowlist Entries

1. `internal/modules/stores/repository/lookup.go`
   - Rationale: This repository is used only during request context bootstrap in `StoreContextMiddleware`.
   - Constraint: It is intentionally tenant-agnostic to resolve store metadata before tenant/site context is finalized.
   - Safety: Any resolved site is validated against tenant context by middleware (`SiteTenantMismatch` protection).

2. `internal/modules/integrations/api_sharing/repository/operation_catalog_mongo.go`
   - Rationale: The API operation registry is a platform-wide catalog; `ApiOperation` has no tenant or site field.
   - Constraint: Only the operation catalog (collection `api_operations`) lives in this file. Tenant-scoped API sharing data (grants, consents, DPAs, audit logs, usage) stays in `mongo.go` and always receives an explicit tenant.
   - Safety: Catalog writes go through the API sharing admin surface (permission-gated); operations describe API shape, not tenant data.

3. `internal/modules/integrations/api_sharing/repository/scope_catalog_mongo.go`
   - Rationale: The API scope registry is a platform-wide catalog; `Scope` has no tenant field (`IsGlobal` only marks catalog visibility).
   - Constraint: Only the scope catalog (collection `api_scopes`) lives in this file; per-tenant permission to use a scope is expressed by tenant-scoped grants in `mongo.go`.
   - Safety: Scopes are definitions; a tenant can only use one through a grant that is tenant-scoped.

4. `internal/modules/system_ops/repository/mongo.go`
   - Rationale: system_ops stores platform lifecycle records of this installation (updates, migrations, backups, maintenance, guardrails, audit); none of its models has a tenant field.
   - Constraint: These records describe the installation, not a tenant's data, and must never hold tenant business data.
   - Safety: Access is permission-gated (`system_ops` permissions) and the menu is hidden from tenant view profiles.

## Maintenance Rule

An allowlist entry must have (a) a matching entry in `repositoryAuditAllowlist` (`internal/platform/tenantfilter/tenantfilter_audit_test.go`), (b) the exact rationale text above, and (c) a model that genuinely has no tenant field. Mixed files (global catalog + tenant data) must be split so that only the catalog file is allowlisted. Background contexts in repository/service layers must come from `internal/platform/contextkeys` (`System`, `ForTenantSite`, `Detached`), never a literal `context.Background()`.
