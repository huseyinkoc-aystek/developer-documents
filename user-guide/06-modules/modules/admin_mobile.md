**Admin Mobile**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Admin Mobile

Purpose
Provides management and access endpoints for Admin Mobile resources.

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
- `GET /admin/mobile/dashboard/summary`
- `GET /admin/mobile/notifications`
- `GET /admin/mobile/notifications/device-tokens`
- `POST /admin/mobile/notifications/device-tokens`
- `POST /admin/mobile/notifications/device-tokens/revoke-all`
- `POST /admin/mobile/notifications/device-tokens/{id}/revoke`
- `POST /admin/mobile/notifications/read-all`
- `GET /admin/mobile/notifications/unread-count`
- `POST /admin/mobile/notifications/{id}/read`
- `GET /admin/mobile/orders`
- `GET /admin/mobile/orders/{id}`
- `GET /admin/mobile/sessions`
- `POST /admin/mobile/sessions/revoke-all`
- `POST /admin/mobile/sessions/{id}/revoke`
- `GET /admin/mobile/stock/critical`

Example Request
```http
GET /admin/mobile/dashboard/summary HTTP/1.1
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
