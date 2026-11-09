**Orders**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Orders

Purpose
Provides management and access endpoints for Orders resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `orders:create`
- `orders:read`
- `orders:update`
- `orders:delete`
- `orders:execute`
- `orders:approve`
- `orders:assign_roles`
- `orders:activate`
- `orders:deactivate`
- `orders:reorder`
- `orders:publish`
- `orders:close`
- `orders:responses_read`
- `orders:responses_delete`
- `orders:enable`
- `orders:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/orders`
- `POST /admin/orders`
- `GET /admin/orders/{id}`
- `POST /admin/orders/{id}/payment-intents`
- `PATCH /admin/orders/{id}/status`

Example Request
```http
GET /admin/orders HTTP/1.1
Host: api.example.com
Authorization: Bearer <access_token>
X-Tenant-Code: acme
X-Store-Code: acme-store
```

Example Response
```json
{
  "success": true,
  "data": {}
}
```

Example Error Responses
```json
{
  "success": false,
  "message": "bad_request"
}
```
