**Isolation Model**
AYSTEK Core enforces strict tenant and site isolation in every request.

**Isolation Layers**
```text
Auth Context
  -> TenantIsolation (tenant_id enforced)
  -> SiteIsolation (site_id enforced)
  -> Data Access (tenant_id + site_id in collections)
```

**Tenant Isolation**
- Rejects requests without tenant context.
- Rejects mismatch between authenticated tenant and request tenant.

**Site Isolation**
- Rejects requests without site context.
- Ensures site belongs to tenant.

**Data Model Constraints**
Collections include:
- `tenant_id`
- `site_id` (where applicable)

**OpenAPI References**
Header parameters are defined globally in `api-docs/swagger.yaml` and required for most operations.

**Example**
```http
GET /admin/orders HTTP/1.1
Host: api.example.com
Authorization: Bearer <access_token>
X-Tenant-Code: acme
X-Store-Code: acme-store
```
