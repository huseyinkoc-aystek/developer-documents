**Pricing**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Pricing

Purpose
Provides management and access endpoints for Pricing resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- No module-specific permission entries in `permission_seed.go`. Access still passes through global permission checks.

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /api/v1/pricing/analytics/{id}`
- `GET /api/v1/pricing/demand/forecast/{id}`
- `POST /api/v1/pricing/demand/update`
- `POST /api/v1/pricing/optimize`
- `GET /api/v1/pricing/products`
- `POST /api/v1/pricing/products`
- `DELETE /api/v1/pricing/products/{id}`
- `GET /api/v1/pricing/products/{id}`
- `PUT /api/v1/pricing/products/{id}`
- `GET /api/v1/pricing/products/{id}/calculate`
- `GET /api/v1/pricing/products/{id}/history`
- `POST /api/v1/pricing/recalculate`
- `GET /api/v1/pricing/rules`
- `POST /api/v1/pricing/rules`
- `DELETE /api/v1/pricing/rules/{id}`
- `GET /api/v1/pricing/rules/{id}`
- `PUT /api/v1/pricing/rules/{id}`

Example Request
```http
GET /api/v1/pricing/analytics/{id} HTTP/1.1
Host: api.example.com
Authorization: Bearer <access_token>
X-Tenant-Code: acme
X-Store-Code: acme-store
```

Example Response
```json
{}
```

Example Error Responses
```json
{
  "success": false,
  "message": "bad_request"
}
```
