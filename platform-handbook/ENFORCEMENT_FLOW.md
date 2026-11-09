**Enforcement Flow**
This document describes the request enforcement pipeline used by AYSTEK Core.

**End-to-End Request Flow**
```text
Request
  -> Auth Context (access token / cookie)
  -> StoreContextResolver (tenant + site)
  -> TenantIsolation
  -> SiteIsolation
  -> ModuleGate
  -> RBAC PermissionMiddleware
  -> LicenseGate
  -> Feature Flags
  -> Service / Repository
  -> Audit Logging + Event Bus
```

**Auth Context**
Auth context is resolved from the `Authorization` header or `access_token` cookie.
OpenAPI reference: `swagger.yaml#/tag/Auth`

**Tenant / Site Context**
Required headers:
- `X-Tenant-ID` or `X-Tenant-Code`
- `X-Store-Code` or `X-Store-Domain`

These are defined as reusable parameters in OpenAPI.

**ModuleGate and RBAC**
Module availability and permissions are enforced server-side. ModuleGate prevents access to disabled modules; RBAC prevents unauthorized actions.
OpenAPI reference: `swagger.yaml#/tag/Permissions`

**License and Feature Flags**
License and feature flags refine behavior but never override ModuleGate or RBAC.
OpenAPI reference: `swagger.yaml#/tag/Plugins` for metadata.

**Sample Admin Flow**
```text
POST /admin/plugins/:id/enable
  -> Auth Context
  -> TenantIsolation + SiteIsolation
  -> ModuleGate (plugins)
  -> RBAC (plugins:enable)
  -> LicenseGate
  -> Feature Flags
  -> Plugin lifecycle update
```

**Canonical Evaluation Order**
See `platform-handbook/CANONICAL_EVALUATION_ORDER.md` for the authoritative sequence.
