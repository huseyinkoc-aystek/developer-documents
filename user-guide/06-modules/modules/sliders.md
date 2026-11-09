**Sliders**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Sliders

Purpose
Provides management and access endpoints for Sliders resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `sliders:create`
- `sliders:read`
- `sliders:update`
- `sliders:delete`
- `sliders:execute`
- `sliders:approve`
- `sliders:assign_roles`
- `sliders:activate`
- `sliders:deactivate`
- `sliders:reorder`
- `sliders:publish`
- `sliders:close`
- `sliders:responses_read`
- `sliders:responses_delete`
- `sliders:enable`
- `sliders:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/sliders`
- `POST /admin/sliders`
- `DELETE /admin/sliders/{id}`
- `GET /admin/sliders/{id}`
- `PUT /admin/sliders/{id}`
- `POST /admin/sliders/{id}/reorder`
- `POST /admin/sliders/{id}/status`
- `GET /public/sliders`
- `POST /public/sliders`
- `DELETE /public/sliders/{id}`
- `GET /public/sliders/{id}`
- `PUT /public/sliders/{id}`
- `POST /public/sliders/{id}/reorder`
- `POST /public/sliders/{id}/status`

Example Request
```http
GET /admin/sliders HTTP/1.1
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
