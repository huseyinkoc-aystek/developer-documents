**Settings**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Settings

Purpose
Provides management and access endpoints for Settings resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `settings:create`
- `settings:read`
- `settings:update`
- `settings:delete`
- `settings:execute`
- `settings:approve`
- `settings:assign_roles`
- `settings:activate`
- `settings:deactivate`
- `settings:reorder`
- `settings:publish`
- `settings:close`
- `settings:responses_read`
- `settings:responses_delete`
- `settings:enable`
- `settings:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/settings`
- `POST /admin/settings`
- `DELETE /admin/settings/{key}`
- `GET /admin/settings/{key}`
- `PUT /admin/settings/{key}`
- `GET /public/settings`
- `GET /public/settings/public`

Example Request
```http
GET /admin/settings HTTP/1.1
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
