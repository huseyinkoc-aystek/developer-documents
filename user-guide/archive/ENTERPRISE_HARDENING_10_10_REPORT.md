---
Title: Enterprise Hardening 10/10 Report
Version: 2.0.0
Status: Archive
Last Updated: 2026-02-12
Related Documents:
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md

---
# Table of Contents
1. [Overview](#overview)

# Enterprise Hardening 10/10 Report

## What Changed

1. OpenAPI contract alignment:
   - Added `X-Site-ID` parameter to OpenAPI components and operation parameter sets.
   - Aligned `x-tenant-context` / `x-site-context` with context matrix by route scope.
   - Canonicalized `ErrorResponse` schema to `{ error, message, trace_id }`.
   - Canonicalized all OpenAPI error examples (`x-example-errors` and 4xx/5xx response examples) to the same schema.
   - Added OpenAPI guard tests for matrix + header + schema + examples consistency.
2. Repository isolation proof:
   - Expanded repository audit to inspect Mongo operations:
     - `Find`, `FindOne`, `Aggregate`, `UpdateOne`, `UpdateMany`, `ReplaceOne`, `DeleteOne`, `DeleteMany`, `BulkWrite`.
   - Added site-scope evidence enforcement for site-scoped repositories.
   - Retained minimal allowlist with explicit rationale documentation.
   - Closed one real tenant scope gap in authors repository (`settings` lookup now tenant-filtered).
3. Public site-optional tolerance verification:
   - Added chain-level tests proving tenant-required, site-optional behavior for:
     - `/public/translations`
     - `/public/settings`
   - Verified public flow does not break in ModuleGate/Permission layers when `site_id` is absent.
4. Observability and audit reinforcement:
   - Settings service now propagates actor identity into audit logs.
   - Tests now prove privileged mutations emit audit logs with tenant/site/user context for:
     - site create/update
     - module enable/disable
     - permission update
     - role update
5. Production hardening and deterministic security events:
   - Added `CrossTenantAccessAttempt` deterministic `403` for auth/header tenant tampering.
   - Added security-event test coverage for:
     - `TenantMissing`
     - `SiteMissing`
     - `InvalidTenant`
     - `InvalidSite`
     - `SiteTenantMismatch`
     - `CrossTenantAccessAttempt`
     - `ContextInvariantViolation`
     - `ContextMutationDetected`
     - `ContextResolutionTimeout`
     - `HostFallbackDisabled`
   - Added repo-wide context setter restriction guard (allowlisted files only) for `tenant_id`/`site_id` writes.
6. Integration (env-gated, real Mongo):
   - Added integration suite that seeds tenant + site fixture and validates `200/400/403` behavior against matrix for representative public/admin/graphql routes.
   - Expanded matrix coverage with high-risk admin mutation endpoints:
     - `PUT /admin/roles/:id`
     - `PUT /admin/settings/:key`

## Why It Was Needed

1. Swagger contract was missing `X-Site-ID` and canonical error schema consistency.
2. Static repository guard needed operation-level coverage, not only file-level heuristics.
3. Public route site-optional contract required explicit proof.
4. Privileged config mutation audit needed stronger user-context attribution guarantees.
5. Security events and tampering outcomes needed deterministic, test-proven outcomes.
6. End-to-end matrix behavior needed real Mongo verification.

## How It Is Proven

1. OpenAPI guards:
   - `internal/platform/openapi/openapi_test.go`
2. Repository scope guards:
   - `internal/platform/tenantfilter/tenantfilter_audit_test.go`
3. Public route tolerance:
   - `internal/platform/router/public_context_test.go`
   - `internal/platform/middleware/module_gate_test.go`
4. Audit context propagation:
   - `internal/modules/settings/service/service_test.go`
   - `internal/modules/stores/service/service_test.go`
   - `internal/modules/permissions/service/service_test.go`
5. Security/tampering guards:
   - `internal/platform/middleware/security_events_test.go`
   - `internal/platform/middleware/context_mutation_guard_test.go`
   - `internal/platform/middleware/store_context_test.go`
6. Integration matrix:
    - `internal/tests/integration/tenant_site_matrix_test.go`

## Remaining Known Risks

1. Repository audit is static-analysis based and conservative; runtime integration checks remain the final guarantee.

## Final Completion Pass (This Update)

1. Repository filter enforcement tightened:
   - `tenantfilter.WithTenantID` and `tenantfilter.WithTenantAndSite` now reject caller-supplied mismatched tenant/site filters (`ErrTenantFilterMismatch`, `ErrSiteFilterMismatch`).
   - Tenant/site values are enforced from request context instead of silently accepting caller values.
   - Aggregate entry-point guard added: `ValidateAggregateTenantStage` requires first `$match` stage with matching `tenant_id`.
2. Protected payload hardening:
   - Product service rejects payload/query tenant override mismatch with deterministic error (`ErrTenantContextMismatch`).
   - Admin product handler maps these mismatches to deterministic `403 CrossTenantAccessAttempt`.
   - FX sync admin handler removed implicit `"default"` tenant fallback and requires resolved context tenant.
   - Admin mobile stock list removed tenant query fallback and blocks mismatch with `403 CrossTenantAccessAttempt`.
3. Authorization fail-open closure:
   - `ModulePermissionMiddleware` now fails closed when permission context is missing/invalid.
   - Platform-user bypass is explicit and isolated.
4. GraphQL hardening:
   - Added query guardrails: max query length, max variables size, max operation count, max depth, max field complexity.
   - Added execution timeout (`5s`) via context timeout wrapping.
5. Observability hardening:
   - Request ID header now validated and bounded before context propagation.
   - Invalid or oversized `X-Request-ID` values are replaced with server-generated IDs.
6. CI/PR guard reinforcement:
   - Added context setter audit to ensure `tenant_id`/`site_id` can only be set from approved middleware files.

## Additional Proof Added

1. `internal/platform/tenantfilter/tenantfilter_hardening_test.go`
2. `internal/platform/middleware/module_permission_test.go`
3. `internal/platform/middleware/request_id_test.go`
4. `internal/platform/middleware/context_setter_audit_test.go`
5. `internal/modules/graphql/handler/admin_security_test.go`
6. `internal/modules/products/handler/admin_tenant_test.go`

## Final Lock Closure (This Pass)

1. AST repository dominance audit implemented:
   - Added `internal/platform/tenantfilter/tenantfilter_ast_dominance_test.go`.
   - Parses repository files with `go/ast`, detects Mongo callsites, and requires tenant dominance proof (helper-dominated filter, tenant-scoped filter expression, or aggregate tenant-stage validation).
   - Covers `Find`, `FindOne`, `CountDocuments`, `Distinct`, `UpdateOne`, `UpdateMany`, `ReplaceOne`, `DeleteOne`, `DeleteMany`, `FindOneAnd*`, `Aggregate`, `BulkWrite`.
2. Service/repository `context.Background()` CI gate implemented:
   - Added `internal/platform/tenantfilter/context_background_audit_test.go`.
   - Fails when `context.Background()` appears in non-test service/repository code.
   - Cleaned existing violations to `context.TODO()` or request-derived detached context where appropriate.
3. Concurrency hardening tests added:
   - Event replay race test (`internal/platform/events/events_security_test.go`) proves single acceptance + replay rejection under concurrent publish.
   - Cache poisoning concurrency test (`internal/platform/cache/service_security_test.go`) proves tenant-key isolation under concurrent access.
4. Security gates consolidated:
   - Added `make security-gates` target in `/workspace/go/aystek-core/Makefile`.
   - Includes tenantfilter audit, header-access audit, context setter restriction audit, and OpenAPI contract guards.
5. Live integration concurrency proof added:
   - `internal/tests/integration/tenant_site_hardening_integration_test.go::TestTenantSiteIsolationConcurrencyIntegration`
   - Concurrent admin/public/tamper request mix validates no tenant/site bleed under load.
6. Multi-env drift proof added:
   - `internal/tests/integration/tenant_site_hardening_integration_test.go::TestTenantSiteMultiEnvDriftIntegration`
   - Header-scoped behavior is stable across `dev/staging/production`; only host-fallback behavior differs by contract.
7. Tenant tombstone replay state-machine guard added:
   - `tenantstate.GuardEventDelivery` enforces status-aware delivery checks.
   - Tombstoned replay rejection error introduced: `ErrTenantTombstonedReplay`.
   - Applied to event publish and Redis subscriber ingest path.
   - Covered by:
     - `internal/platform/tenantstate/state_test.go`
     - `internal/platform/events/events_security_test.go`
     - `internal/platform/jobs/runner/worker_security_test.go`
