# AYSTEK Core API Versioning Policy (v1.0)

This document defines the official API versioning strategy for AYSTEK Core. It is a policy and documentation artifact only. It does not change runtime behavior.

---

## 1) Current State (v1.0)

- **Implicit v1:** All current API routes without a version prefix are treated as **v1**.
- **No `/v1` prefix:** There is no active `/v1` routing prefix and no change is required for current consumers.
- **Swagger:** The `swagger.yaml` describes the implicit v1 surface.

---

## 2) Backward Compatibility Rules

Backward compatibility is mandatory for all v1 routes unless explicitly declared deprecated.

**Allowed without breaking clients**
1. Adding new endpoints.
2. Adding optional request fields or headers.
3. Adding new response fields.
4. Adding new enum values if clients are not required to handle them.
5. Fixing server bugs that produce incorrect results.

**Breaking changes (not allowed in v1 without deprecation cycle)**
1. Removing endpoints or response fields.
2. Renaming endpoints, fields, or parameters.
3. Changing data types or required fields.
4. Changing error codes or response shapes.
5. Tightening auth/permission requirements for existing endpoints.
6. Changing tenant/site context requirements for existing endpoints.

---

## 3) Deprecation Window Policy

All breaking changes must go through a **soft deprecation** window before removal.

**Required steps**
1. Mark the legacy endpoint as deprecated in Swagger.
2. Emit deprecation headers in runtime responses:
   - `Deprecation: true`
   - `Sunset: <RFC1123 date>`
   - `Link: <...>; rel="successor-version"`
3. Document the canonical replacement in platform docs.
4. Maintain full runtime behavior during the deprecation window.

**Minimum deprecation window**
- **12 months** for general endpoints.
- Longer windows may be defined for high‑impact enterprise integrations.

---

## 4) v2 Introduction Strategy

When v2 is introduced, it will use **explicit versioning**:

- **Canonical base path:** `/api/v2/`
- v1 remains active at the implicit root paths.
- v2 will be added **without breaking existing v1 routes**.
- Only v2 will contain the new canonicalized endpoints that are not compatible with v1.

**Routing rules**
1. v1 routes remain unprefixed.
2. v2 routes use `/api/v2/...`.
3. No `/v1` prefix is introduced.

---

## 5) Plugin Compatibility Policy

Plugins must declare compatibility ranges relative to core versions.

**Rules**
1. Plugins must declare `min_core_version` and `max_core_version`.
2. v2-only plugins must not be enabled on v1 cores.
3. Marketplace listings must show compatible core versions.
4. Runtime gating must refuse activation outside the declared range.

**Compatibility guarantees**
- v1 plugins are assumed compatible with the implicit v1 API surface.
- v2 introduces explicit compatibility checks and may require plugin upgrades.

---

## 6) Swagger Version Tagging Rules

Swagger must reflect API versioning accurately:

1. v1 spec remains `api-docs/swagger.yaml` and documents unprefixed routes.
2. v2 spec will be emitted as `api-docs/swagger-v2.yaml` (or a separate file) and must include `/api/v2/...` paths only.
3. Deprecated v1 endpoints must include `deprecated: true` and `x-deprecation-replacement`.
4. Swagger `info.version` must match the core API version it describes.

---

## 7) Migration Example (v1 → v2)

**Scenario:** a legacy endpoint `/security/audit-log` is deprecated and its canonical v2 equivalent becomes `/api/v2/admin/security/audit-log`.

**Steps**
1. Add `/api/v2/admin/security/audit-log`.
2. Mark `/security/audit-log` as deprecated in Swagger with:
   - `x-deprecation-replacement: /api/v2/admin/security/audit-log`
3. Emit deprecation headers on `/security/audit-log` responses.
4. Maintain both endpoints during the deprecation window.
5. After sunset date, remove `/security/audit-log` from v1.

---

## 8) Deprecation Lifecycle Timeline (Example)

```
T0: 2026-02-11  Legacy endpoint remains canonical (v1)
T1: 2026-03-01  Mark deprecated; add deprecation headers
T2: 2026-03-15  Publish v2 endpoint
T3: 2027-03-01  Sunset date reached (12 months)
T4: 2027-04-01  Remove legacy endpoint (v1)
```

---

## 9) Marketplace Compatibility Implications

Marketplace and plugin distribution must enforce version compatibility:

1. A plugin package must declare compatible core version ranges.
2. Marketplace listings must display compatibility and deprecation status.
3. Deprecated endpoints in v1 must be indicated in plugin docs if used.
4. v2-only plugins must be rejected by v1 cores at install/enable time.

---

## Summary of Future Migration Strategy

AYSTEK Core will maintain implicit v1 for all current unprefixed routes. v2 will be introduced under `/api/v2/` without breaking existing clients. Deprecated v1 endpoints will remain functional for a minimum of 12 months and will emit standard deprecation headers with a published sunset date and a clear replacement path. Plugins and marketplace listings will enforce compatibility based on core version ranges.
