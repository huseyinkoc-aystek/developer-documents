# AYSTEK Core

## Platform Architecture & Execution Standard (v1)

This document defines the **canonical architecture, execution order, and non-negotiable standards** of the AYSTEK Core platform.

It exists to ensure that:

- The platform remains coherent over time
- New developers or AI agents can continue work without context loss
- Architectural decisions are not re-litigated repeatedly

This document is **architecture-first**.
Code implementation must follow this document, not reinterpret it.

---

## 1. Platform Definition

AYSTEK Core is a **platform-level SaaS system**, not a single product.

Key characteristics:

- Multi-tenant
- Multi-site per tenant
- Multi-module per site
- Single backend authority
- Controlled admin interface
- Open public site ecosystem

---

## 2. System Components (Authoritative)

### 2.1 Backend (Single Source of Truth)

- Name: `aystek-core`
- Responsibilities:
  - Auth & identity
  - Tenant, site, module management
  - Permissions & policy enforcement
  - Data access & isolation
  - Audit & observability

There is **only one backend**.

### 2.2 Admin UI

- Name: `aystek-core-next`
- Purpose:
  - Operational control
  - Tenant & site administration
  - Module configuration
- Characteristics:
  - Single codebase
  - No per-customer customization
  - Fully driven by backend context (tenant, site, permissions)

UI **never** decides authority.
UI **only reflects** backend decisions.

### 2.3 Public Websites

- Built externally (customer or agency)
- Any technology stack
- Consume public APIs from `aystek-core`
- Must respect domain → site → tenant binding
- Must take the offered languages from `GET /public/languages/active` (site-aware) and never hard-code a language list

AYSTEK Core does **not** own public UI design.

---

## 3. Identity & User Model

### 3.1 Platform Users (AYSTEK Team)

- Not bound to a tenant
- Can access all tenants and sites
- Can impersonate tenant users (logged & audited)
- Use `aystek-core-next`

### 3.2 Tenant Users

- Always belong to exactly one tenant
- Can be:
  - Global (all sites)
  - Site-scoped (specific sites)
- Roles & permissions enforced server-side

---

## 4. Tenant & Site Model

### 4.1 Tenant

- Represents a customer organization
- Owns:
  - Users
  - Sites
  - Data
  - Configuration
- Can scale independently

### 4.2 Site

- Represents **domain + purpose**
- Belongs to one tenant
- Enables a subset of modules
- Is the primary runtime context

**Modules are attached to sites, not tenants.**

#### 4.2.1 Site-scoped languages

- The **tenant** owns the language catalog (names, native names, tenant-wide active/default).
- A **site** may narrow that catalog with an override (`site_language_configs`: ordered codes + default).
- No override means the site follows the tenant-wide active languages.
- Public sites read the result from `/public/languages/active` and `/public/languages/default`.
- Admin UI: `/admin/languages` (catalog) and `/admin/languages/site` (selected site's selection).

---

## 5. Module System

- Modules follow DDD boundaries (CMS, Forum, Commerce, etc.)
- Every module must be:
  - Tenant-aware
  - Site-aware
- Module availability is controlled via:
  - Site configuration
  - Licensing / quota rules
  - Permission checks

A module being disabled must:

- Hide UI entry points
- Reject API access server-side

---

## 6. Data Strategy (MongoDB)

### 6.1 Initial Strategy

- Shared MongoDB database
- Every collection includes:
  - tenant_id
  - site_id (where applicable)

Indexes must always include tenant_id.

### 6.2 Future Strategy (Planned, Not Immediate)

- Selective tenant isolation into dedicated databases
- Controlled by tenant configuration:
  - data_strategy: shared | dedicated
- Database routing handled via abstraction layer

**No business logic may assume a single database.**

---

## 7. Core Abstractions (Mandatory)

The following abstractions must exist to avoid future refactors:

- TenantResolver
- SiteResolver
- ModuleGate / AccessGate
- Permission Policy Engine
- DataRouter (DB routing)
- Audit & Event Bus

No module may bypass these abstractions.

---

## 8. Execution Order (Canonical)

Development must follow this order:

1. Auth & identity resolution
2. Tenant context resolution
3. Site context resolution
4. Module gating
5. Permission enforcement
6. Data access standardization
7. Admin UI dynamic rendering
8. Observability & audit logging

Skipping steps introduces systemic risk.

---

## 9. Over-Engineering Rules (What NOT to Do)

- Do NOT create per-tenant databases early
- Do NOT fork admin panels
- Do NOT hardcode module visibility in UI
- Do NOT bind public site tech choices to core
- Do NOT embed AI features tightly into core flows

---

## 10. Risk Management

| Risk                | Mitigation                             |
| ------------------- | -------------------------------------- |
| Tenant data leakage | Repository-level tenant enforcement    |
| Permission bypass   | API middleware + tests                 |
| UI / backend drift  | Backend-driven menu & permission seeds |
| Module mismatch     | ModuleGate enforcement                 |
| Future DB migration | DataRouter abstraction                 |

---

## 11. First Sprint Definition (Baseline)

Minimum platform-complete state:

- Tenant + Site context middleware
- ModuleGate & permission enforcement
- Repository layer tenant/site requirement
- Dynamic admin menu rendering
- Audit logging (impersonation, config changes)

---

## 12. Long-Term Evolution

Planned (not immediate):

- Tenant-dedicated databases
- Fine-grained policy injection
- AI-driven moderation & analytics
- Public integration contracts & SDKs

---

## 13. Final Principle

> AYSTEK Core prioritizes **architectural correctness over short-term speed**.
> Any implementation that violates this document is considered invalid, even if it works.
