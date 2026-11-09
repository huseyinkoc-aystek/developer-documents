**Stores**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Stores

Purpose
Provides management and access endpoints for Stores resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `stores:create`
- `stores:read`
- `stores:update`
- `stores:delete`
- `stores:execute`
- `stores:approve`
- `stores:assign_roles`
- `stores:activate`
- `stores:deactivate`
- `stores:reorder`
- `stores:publish`
- `stores:close`
- `stores:responses_read`
- `stores:responses_delete`
- `stores:enable`
- `stores:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/stores`
- `POST /admin/stores`
- `GET /admin/stores/current`
- `DELETE /admin/stores/{id}`
- `GET /admin/stores/{id}`
- `PUT /admin/stores/{id}`

Example Request
```http
GET /admin/stores HTTP/1.1
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
