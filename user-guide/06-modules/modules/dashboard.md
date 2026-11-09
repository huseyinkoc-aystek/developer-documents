**Dashboard**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Dashboard

Purpose
Provides management and access endpoints for Dashboard resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `dashboard:create`
- `dashboard:read`
- `dashboard:update`
- `dashboard:delete`
- `dashboard:execute`
- `dashboard:approve`
- `dashboard:assign_roles`
- `dashboard:activate`
- `dashboard:deactivate`
- `dashboard:reorder`
- `dashboard:publish`
- `dashboard:close`
- `dashboard:responses_read`
- `dashboard:responses_delete`
- `dashboard:enable`
- `dashboard:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/dashboard`
- `GET /admin/dashboard/analytics`
- `GET /admin/dashboard/analytics-modules`
- `GET /admin/dashboard/categories`
- `GET /admin/dashboard/health`
- `GET /admin/dashboard/metrics`
- `GET /admin/dashboard/payment-methods`
- `GET /admin/dashboard/realtime`
- `GET /admin/dashboard/top-customers`
- `GET /admin/dashboard/top-products`

Example Request
```http
GET /admin/dashboard HTTP/1.1
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
