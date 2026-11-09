**Plugins**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Plugins

Purpose
Provides management and access endpoints for Plugins resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `plugins:create`
- `plugins:read`
- `plugins:update`
- `plugins:delete`
- `plugins:execute`
- `plugins:approve`
- `plugins:assign_roles`
- `plugins:activate`
- `plugins:deactivate`
- `plugins:reorder`
- `plugins:publish`
- `plugins:close`
- `plugins:responses_read`
- `plugins:responses_delete`
- `plugins:enable`
- `plugins:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/marketplace/plugins`
- `GET /admin/marketplace/plugins/{key}`
- `GET /admin/plugins`
- `POST /admin/plugins`
- `GET /admin/plugins/extensions`
- `GET /admin/plugins/metadata`
- `DELETE /admin/plugins/{id}`
- `GET /admin/plugins/{id}`
- `PUT /admin/plugins/{id}`
- `POST /admin/plugins/{id}/disable`
- `POST /admin/plugins/{id}/enable`

Example Request
```http
GET /admin/marketplace/plugins HTTP/1.1
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
