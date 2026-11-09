# AYSTEK Core — Execution Diagnostic Report

**Scope:** Evidence-based snapshot of current state vs Canonical Architecture and AI standards.  
**Authoritative references:**
- `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`
- `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`

---

## Evidence Sources (Observed)

- **Canonical Architecture:** `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`
- **AI Standards:** `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`
- **Codebase:** `/workspace/go/aystek-core`
- **Git working tree:** large set of modified files across platform middleware, seeds, and modules (`git status --short`). Key relevant paths include:
  - `internal/platform/middleware/auth.go`
  - `internal/platform/middleware/store_context.go`
  - `internal/platform/middleware/site_isolation.go`
  - `internal/platform/middleware/module_gate.go`
  - `internal/platform/middleware/permission.go`
  - `internal/platform/middleware/impersonation_audit.go`
  - `internal/platform/router/groups.go`
  - `internal/platform/db/router.go`
  - `internal/modules/menus/handler/admin.go`
  - `internal/modules/menus/handler/menu_gate.go`
  - `internal/modules/menus/routes.go`
  - `internal/modules/translations/handler/public.go`
  - `internal/seeds/permission_seed.go`
  - `internal/seeds/role_seed.go`
  - `internal/seeds/menu_seed.go`
  - `internal/seeds/translation_seed.go`
  - `internal/modules/activity_logs/service/service.go`
  - `internal/modules/forum/repository/mongo_thread.go`

---

## Canonical Execution Steps (1–8) — Evidence Map

### Step 1 — Auth & Identity Resolution
**Status:** ✅ Present

**Evidence:**
- `internal/platform/middleware/auth.go`: `AuthContextMiddleware` validates access tokens, resolves user, sets `user_type` (`platform` vs `tenant`) and `is_platform_user`.
- `internal/platform/router/groups.go`: `ProtectedChain` applies `AuthContextMiddleware(..., true)` before other middleware.

**Risk if missing:** Unauthorized access and identity ambiguity (Canonical §8).

---

### Step 2 — Tenant Context Resolution
**Status:** ✅ Present

**Evidence:**
- `internal/platform/middleware/store_context.go`: `StoreContextMiddleware` resolves tenant from headers (`X-Tenant-ID` / `X-Tenant-Code`) or store-domain mapping, populates `tenant_id`.
- `internal/platform/middleware/gateway.go`: `TenantIsolation` enforces `tenant_id` and prevents mismatched `auth_tenant_id`.
- `internal/platform/middleware/request_context.go`: propagates `tenant_id` into context.

**Risk if missing:** Cross-tenant access, data leakage.

---

### Step 3 — Site Context Resolution
**Status:** ✅ Present

**Evidence:**
- `internal/platform/middleware/store_context.go`: resolves site via store domain (`X-Store-Domain` or host), sets `site_id`/`store_id` and `site`.
- `internal/platform/middleware/site_isolation.go`: rejects requests without `site_id` and enforces site–tenant binding.

**Risk if missing:** Cross-site data leakage, invalid site access.

---

### Step 4 — Module Gate (AccessGate / ModuleGate)
**Status:** ✅ Present

**Evidence:**
- `internal/platform/middleware/module_gate.go`: `ModuleGateByPathMiddleware` checks module enablement with hierarchical settings keys:
  - `site.<siteID>.module.<module>.enabled`
  - `tenant.<tenantID>.module.<module>.enabled`
  - `module.<module>.enabled`
- `internal/platform/router/groups.go`: `ModuleGateByPathMiddleware` applied before `PermissionMiddleware` in protected/public chains.

**Risk if missing:** Disabled modules remain callable (policy breach).

---

### Step 5 — Permission & Policy Enforcement
**Status:** ✅ Present

**Evidence:**
- `internal/platform/middleware/permission.go`: `PermissionMiddleware` enforces module-level permission (`module:*` or `*:*`), platform user bypass only if not impersonating.
- `internal/platform/middleware/auth.go`: `resolvePermissions` from role-permission assignments.
- `internal/seeds/permission_seed.go`: seeds module-action permissions.
- `internal/seeds/role_seed.go`: `admin` role owns all permissions.

**Risk if missing:** RBAC bypass, admin/user privilege escalation.

---

### Step 6 — Repository Layer Enforcement (tenant/site)
**Status:** ⚠️ Partial

**Evidence (positive):**
- `internal/modules/forum/repository/mongo_thread.go`: all queries require `tenant_id` (via `tenantfilter.RequireTenantID`) and include `{tenant_id: ...}` filters.
- `internal/platform/db/router.go`: `DataRouter` abstraction exists; currently shared DB implementation.

**Gaps:**
- No repo-wide proof that *all* modules enforce `tenant_id`/`site_id` (only confirmed for selected modules).
- Site-level enforcement in repositories is not consistently evident.

**Risk if missing:** Cross-tenant/site leakage at data layer.

---

### Step 7 — Admin UI Dynamic Rendering Support
**Status:** ⚠️ Partial

**Evidence (positive):**
- `internal/modules/menus/handler/admin.go`: `/admin/menus` returns menus with tree filters by role and module gate.
- `internal/modules/menus/handler/menu_gate.go`: filters menu tree based on roles and module availability.
- `internal/modules/menus/routes.go`: `/admin/menus` endpoints defined.
- `internal/seeds/menu_seed.go`: menu entries are seeded.
- `internal/seeds/translation_seed.go`: 28-language menu label keys exist.

**Gaps:**
- No explicit “enabled modules” or “permission set” metadata endpoint confirmed for UI consumption beyond menus.

**Risk if missing:** UI drift and inconsistent visibility logic.

---

### Step 8 — Observability & Audit Logging (Minimum Set)
**Status:** ⚠️ Partial

**Evidence (positive):**
- `internal/platform/middleware/impersonation_audit.go`: logs impersonation via activity logs.
- `internal/modules/activity_logs/service/service.go`: tenant/site-scoped audit log creation.

**Gaps:**
- No evidence of audit logging for module enable/disable, permission changes, or site creation updates.

**Risk if missing:** Compliance and forensic gaps; privileged actions not auditable.

---

## Intent Checklist Coverage (from user confirmation)

- **Platform (not a website):** ✅ Documented in Canonical Architecture §1–2.
- **Tenant = customer organization:** ✅ Canonical §4.1.
- **Site = domain + purpose:** ✅ Canonical §4.2; store model includes `Domain`.
- **Modules attached to site (not tenant):** ✅ Canonical §4.2; ModuleGate supports site-scoped settings.
- **Users belong to tenant; optional site scope:** ⚠️ User model includes `tenant_id` but no explicit site scope fields found in `internal/modules/users/model/user.go`.
- **PlatformUser full access + impersonation audit:** ✅ `auth.go` + `impersonation_audit.go` evidence.
- **Single admin UI, backend-driven:** ✅ Canonical §2.2; backend menu endpoints and filtering exist.
- **Public sites external API consumers; domain→site→tenant binding:** ✅ Store context resolution and site isolation enforce domain and tenant binding.
- **Shared MongoDB + tenant_id (+ site_id where applicable):** ⚠️ Evidence of tenant_id in several models; site_id not consistently verified across repositories.
- **Canonical execution order:** ✅ Middleware chain order in `internal/platform/router/groups.go` aligns with required order.

---

## Gaps (Numbered)

**G-01 — Repository tenant/site enforcement not proven across all modules**
- **Missing:** Global enforcement verification for tenant_id and site_id in all repositories.
- **Where it should live:** Repository interfaces and implementations across `internal/modules/*/repository`.
- **Verify completion:** Audit all repository queries for tenant/site filters; add tests to ensure cross-tenant access fails.

**G-02 — Site-scoped users not represented**
- **Missing:** User model or user-site relationship for site-scoped access.
- **Where it should live:** `internal/modules/users` (model + repository/service) or a join/association module.
- **Verify completion:** Tests demonstrating user access restricted to specific site IDs.

**G-03 — Admin metadata endpoints incomplete**
- **Missing:** Explicit backend endpoints for module availability and permission sets for UI rendering (beyond menus).
- **Where it should live:** `internal/modules/permissions` or `internal/modules/settings` admin handlers.
- **Verify completion:** `/admin/*` endpoints returning modules/permissions tied to current tenant/site.

**G-04 — Audit coverage incomplete**
- **Missing:** Audit logs for module enable/disable, permission changes, site creation/updates.
- **Where it should live:** Relevant services/handlers (settings, permissions, stores/site) using `activity_logs` service.
- **Verify completion:** Create/Update actions emit logs; tests or logs show entries.

**G-05 — Site enforcement in data models not fully evidenced**
- **Missing:** site_id presence and enforced usage in models/repositories across modules.
- **Where it should live:** `internal/modules/*/model` and repository filters.
- **Verify completion:** Model schema includes `site_id`; repository queries include `site_id`.

---

## Verification Checklist (How to Confirm Completion)

- **Tenant/site middleware:** Confirm with integration tests that missing tenant/site returns 400 (middleware already present).
- **Repository enforcement:** Inspect all repo queries for tenant/site filters; add tests for cross-tenant access.
- **Admin metadata:** Verify `/admin/menus` returns filtered tree and additional endpoints exist for modules/permissions.
- **Audit events:** Confirm activity log entries for impersonation + administrative actions.

---

## Summary (State Snapshot)

- Core middleware chain matches canonical order and enforces auth → tenant → site → module gate → RBAC.
- DataRouter abstraction exists for future dedicated DB routing.
- Backend-driven menus and translations are seeded and served; filtering by roles and module gate is implemented.
- Audit logging exists for impersonation, but broader admin action audit coverage is not evidenced.
- Repository-level tenant/site enforcement is confirmed in forum module, but not proven across all modules.

