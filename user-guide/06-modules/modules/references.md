**References**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/References

Purpose
Provides management and access endpoints for References resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `references:create`
- `references:read`
- `references:update`
- `references:delete`
- `references:execute`
- `references:approve`
- `references:assign_roles`
- `references:activate`
- `references:deactivate`
- `references:reorder`
- `references:publish`
- `references:close`
- `references:responses_read`
- `references:responses_delete`
- `references:enable`
- `references:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/references`
- `POST /admin/references`
- `POST /admin/references/create`
- `POST /admin/references/delete`
- `GET /admin/references/detail`
- `GET /admin/references/list`
- `PUT /admin/references/update`
- `DELETE /admin/references/{id}`
- `GET /admin/references/{id}`
- `PUT /admin/references/{id}`
- `POST /admin/references/{id}/order`
- `POST /admin/references/{id}/status`
- `GET /public/references`
- `POST /public/references`
- `POST /public/references/create`
- `POST /public/references/delete`
- `GET /public/references/detail`
- `GET /public/references/list`
- `PUT /public/references/update`
- `DELETE /public/references/{id}`
- `GET /public/references/{id}`
- `PUT /public/references/{id}`
- `POST /public/references/{id}/order`
- `POST /public/references/{id}/status`

Example Request
```http
GET /admin/references HTTP/1.1
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
