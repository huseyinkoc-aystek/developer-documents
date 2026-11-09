**Authors**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Authors

Purpose
Provides management and access endpoints for Authors resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `authors:create`
- `authors:read`
- `authors:update`
- `authors:delete`
- `authors:execute`
- `authors:approve`
- `authors:assign_roles`
- `authors:activate`
- `authors:deactivate`
- `authors:reorder`
- `authors:publish`
- `authors:close`
- `authors:responses_read`
- `authors:responses_delete`
- `authors:enable`
- `authors:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/authors`
- `POST /admin/authors`
- `GET /admin/authors/`
- `DELETE /admin/authors/by-id/{id}`
- `GET /admin/authors/by-id/{id}`
- `PUT /admin/authors/by-id/{id}`
- `DELETE /admin/authors/{id}`
- `GET /admin/authors/{id}`
- `PUT /admin/authors/{id}`
- `POST /admin/authors/{id}/toggle`
- `GET /admin/authors/{slug}`
- `GET /authors`
- `POST /authors`
- `GET /authors/`
- `DELETE /authors/by-id/{id}`
- `GET /authors/by-id/{id}`
- `PUT /authors/by-id/{id}`
- `DELETE /authors/{id}`
- `GET /authors/{id}`
- `PUT /authors/{id}`
- `POST /authors/{id}/toggle`
- `GET /public/authors`
- `POST /public/authors`
- `GET /public/authors/`
- `DELETE /public/authors/by-id/{id}`
- `GET /public/authors/by-id/{id}`
- `PUT /public/authors/by-id/{id}`
- `DELETE /public/authors/{id}`
- `GET /public/authors/{id}`
- `PUT /public/authors/{id}`
- `POST /public/authors/{id}/toggle`
- `GET /public/authors/{slug}`

Example Request
```http
GET /admin/authors HTTP/1.1
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
