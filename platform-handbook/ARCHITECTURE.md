**Architecture**
AYSTEK Core v1.0 is a multi-tenant, multi-site platform with a single backend authority. This document aligns to the canonical architecture and references the OpenAPI spec at `api-docs/swagger.yaml`.

**System Overview**
```text
                 +---------------------------+
                 |     aystek-core-next      |
                 |  Admin UI (read-only)     |
                 +-------------+-------------+
                               |
                               v
+------------------------------------------------------------------+
|                           aystek-core                            |
|  Auth + Tenant/Site + ModuleGate + RBAC + Plugins + Audit + SSE   |
+--------------+---------------------------+------------------------+
               |                           |
               v                           v
        +------+-----+              +------+-----+
        | MongoDB    |              | Event Bus  |
        | Redis      |              | SSE Streams|
        +------------+              +------------+

Public Sites / Integrators -> /public APIs -> aystek-core
```

**Tenant / Site Model**
```text
Tenant
  -> Sites (domains / business contexts)
       -> Modules enabled per site
       -> Plugins enabled per site
```

**API Surfaces (OpenAPI)**
- Authentication: `/auth/*` (see `swagger.yaml#/tag/Auth`)
- Admin: `/admin/*` (see module tags under `swagger.yaml`)
- Public: `/public/*` (see module tags under `swagger.yaml`)

**Example Admin Operations**
- `/admin/stores` (site management)
- `/admin/permissions` (RBAC)
- `/admin/plugins/metadata` (plugin state and gating)

**Example Public Operations**
- `/public/settings`
- `/public/authors`
- `/public/ads`
