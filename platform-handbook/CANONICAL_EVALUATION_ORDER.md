**Canonical Evaluation Order**
This document resolves ordering conflicts between enforcement flow and plugin runtime. The sequence below is the single authoritative evaluation order for AYSTEK Core v1.0 request handling.

1. **Auth / Identity Resolution**
   - Establishes requester identity, roles, and base permissions.
   - Required before any tenant or site checks that depend on user context.

2. **Tenant Resolution**
   - Resolves `tenant_id` from headers or domain.
   - Enforces tenant isolation (rejects missing or mismatched tenant).

3. **Site Resolution**
   - Resolves `site_id` from headers or domain.
   - Enforces site isolation (rejects missing site or cross-tenant site use).

4. **ModuleGate**
   - Prevents access to modules disabled for the resolved site.
   - Must run before RBAC and plugin evaluation to avoid leaking module surfaces.

5. **RBAC Permission Enforcement**
   - Validates action-level permissions (`module:action`) for the request.
   - Must run after ModuleGate so disabled modules never pass RBAC.

6. **LicenseGate**
   - Verifies that the tenant/site license grants access to module/plugin capabilities.
   - License is a hard override and must be evaluated before feature flags.

7. **Feature Flag Evaluation**
   - Applies site → tenant → global → default flag resolution.
   - Flags refine behavior but never enable a disabled module or plugin.

8. **Hook Evaluator (Plugin Runtime)**
   - For plugin-managed flows, evaluates hook execution eligibility using:
     - Effective plugin enablement
     - Allowed capabilities
     - License and flag results
   - Runs after gates and flags to ensure hooks cannot bypass enforcement.

9. **Repository Enforcement**
   - Data access and persistence with tenant/site constraints.
   - Final enforcement layer for isolation and consistency.

This order is aligned with the canonical architecture and the backend’s middleware chain.
