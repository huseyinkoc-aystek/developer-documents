**Comments**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Comments

Purpose
Provides management and access endpoints for Comments resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `comments:create`
- `comments:read`
- `comments:update`
- `comments:delete`
- `comments:execute`
- `comments:approve`
- `comments:assign_roles`
- `comments:activate`
- `comments:deactivate`
- `comments:reorder`
- `comments:publish`
- `comments:close`
- `comments:responses_read`
- `comments:responses_delete`
- `comments:enable`
- `comments:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/comments`
- `POST /admin/comments`
- `GET /admin/comments/`
- `POST /admin/comments/bulk/delete`
- `POST /admin/comments/bulk/status`
- `DELETE /admin/comments/{id}`
- `GET /admin/comments/{id}`
- `PUT /admin/comments/{id}`
- `POST /admin/comments/{id}/like`
- `POST /admin/comments/{id}/reaction`
- `POST /admin/comments/{id}/reply`
- `PUT /admin/comments/{id}/seo`
- `PUT /admin/comments/{id}/status`
- `GET /comments`
- `GET /comments/`
- `POST /comments/bulk/delete`
- `POST /comments/bulk/status`
- `DELETE /comments/{id}`
- `GET /comments/{id}`
- `PUT /comments/{id}`
- `POST /comments/{id}/like`
- `POST /comments/{id}/reaction`
- `POST /comments/{id}/reply`
- `PUT /comments/{id}/seo`
- `PUT /comments/{id}/status`
- `GET /public/comments`
- `POST /public/comments`
- `GET /public/comments/`
- `POST /public/comments/bulk/delete`
- `POST /public/comments/bulk/status`
- `DELETE /public/comments/{id}`
- `GET /public/comments/{id}`
- `PUT /public/comments/{id}`
- `POST /public/comments/{id}/like`
- `POST /public/comments/{id}/reaction`
- `POST /public/comments/{id}/reply`
- `PUT /public/comments/{id}/seo`
- `PUT /public/comments/{id}/status`

Example Request
```http
GET /admin/comments HTTP/1.1
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
