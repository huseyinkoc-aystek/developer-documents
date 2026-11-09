**Returns**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Returns

Purpose
Provides management and access endpoints for Returns resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `returns:create`
- `returns:read`
- `returns:update`
- `returns:delete`
- `returns:execute`
- `returns:approve`
- `returns:assign_roles`
- `returns:activate`
- `returns:deactivate`
- `returns:reorder`
- `returns:publish`
- `returns:close`
- `returns:responses_read`
- `returns:responses_delete`
- `returns:enable`
- `returns:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/returns`
- `POST /admin/returns`
- `GET /admin/returns/{id}`
- `POST /admin/returns/{id}/authorize`
- `POST /admin/returns/{id}/inspect`
- `POST /admin/returns/{id}/receive`
- `PATCH /admin/returns/{id}/refund`
- `PATCH /admin/returns/{id}/status`

Example Request
```http
GET /admin/returns HTTP/1.1
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
