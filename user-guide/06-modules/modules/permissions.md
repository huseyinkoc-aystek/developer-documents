**Permissions**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Permissions

Purpose
Provides management and access endpoints for Permissions resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `permissions:create`
- `permissions:read`
- `permissions:update`
- `permissions:delete`
- `permissions:execute`
- `permissions:approve`
- `permissions:assign_roles`
- `permissions:activate`
- `permissions:deactivate`
- `permissions:reorder`
- `permissions:publish`
- `permissions:close`
- `permissions:responses_read`
- `permissions:responses_delete`
- `permissions:enable`
- `permissions:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/permissions`
- `POST /admin/permissions`
- `GET /admin/permissions/modules`
- `DELETE /admin/permissions/{id}`
- `PUT /admin/permissions/{id}`
- `GET /admin/roles`
- `POST /admin/roles`
- `DELETE /admin/roles/{id}`
- `GET /admin/roles/{id}`
- `PUT /admin/roles/{id}`
- `GET /admin/roles/{id}/permissions`

Example Request
```http
GET /admin/permissions HTTP/1.1
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
