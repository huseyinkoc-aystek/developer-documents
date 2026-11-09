**Stock**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Stock

Purpose
Provides management and access endpoints for Stock resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `stock:create`
- `stock:read`
- `stock:update`
- `stock:delete`
- `stock:execute`
- `stock:approve`
- `stock:assign_roles`
- `stock:activate`
- `stock:deactivate`
- `stock:reorder`
- `stock:publish`
- `stock:close`
- `stock:responses_read`
- `stock:responses_delete`
- `stock:enable`
- `stock:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/stock/export`
- `GET /admin/stock/jobs`
- `GET /admin/stock/jobs/{id}`
- `GET /admin/stock/variants`
- `POST /admin/stock/variants/bulk`
- `POST /admin/stock/variants/bulk-async`
- `PATCH /admin/stock/variants/{variant_id}`
- `GET /admin/stock/warehouses`
- `POST /admin/stock/warehouses`
- `POST /admin/stock/warehouses/transfer/{from_warehouse_id}/{to_warehouse_id}/variants/{variant_id}`
- `GET /admin/stock/warehouses/{warehouse_id}/inventory/{variant_id}`
- `POST /admin/stock/warehouses/{warehouse_id}/variants/{variant_id}/adjust`
- `GET /api/checkouts/{checkout_id}/reservation`
- `GET /api/reservations`
- `POST /api/reservations`
- `GET /api/reservations/{id}`
- `PATCH /api/reservations/{id}`
- `POST /api/reservations/{id}/cancel`
- `POST /api/reservations/{id}/confirm`
- `POST /api/reservations/{id}/release`
- `GET /health/reservations`

Example Request
```http
GET /admin/stock/export HTTP/1.1
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
