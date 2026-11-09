**Pages**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Pages

Purpose
Provides management and access endpoints for Pages resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `pages:create`
- `pages:read`
- `pages:update`
- `pages:delete`
- `pages:execute`
- `pages:approve`
- `pages:assign_roles`
- `pages:activate`
- `pages:deactivate`
- `pages:reorder`
- `pages:publish`
- `pages:close`
- `pages:responses_read`
- `pages:responses_delete`
- `pages:enable`
- `pages:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/pages`
- `POST /admin/pages`
- `DELETE /admin/pages/by-id/{id}`
- `GET /admin/pages/by-id/{id}`
- `PUT /admin/pages/by-id/{id}`
- `GET /admin/pages/filter`
- `GET /admin/pages/localized/{lang}/{slug}`
- `GET /admin/pages/slug/{slug}`
- `DELETE /admin/pages/{id}`
- `GET /admin/pages/{id}`
- `PUT /admin/pages/{id}`
- `PATCH /admin/pages/{id}/sort`
- `GET /admin/pages/{lang}/{slug}`
- `GET /public/pages`
- `POST /public/pages`
- `DELETE /public/pages/by-id/{id}`
- `GET /public/pages/by-id/{id}`
- `PUT /public/pages/by-id/{id}`
- `GET /public/pages/filter`
- `GET /public/pages/localized/{lang}/{slug}`
- `GET /public/pages/slug/{slug}`
- `DELETE /public/pages/{id}`
- `GET /public/pages/{id}`
- `PUT /public/pages/{id}`
- `PATCH /public/pages/{id}/sort`
- `GET /public/pages/{lang}/{slug}`

Example Request
```http
GET /admin/pages HTTP/1.1
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
