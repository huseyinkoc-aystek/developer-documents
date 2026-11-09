**Media**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Media

Purpose
Provides management and access endpoints for Media resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `media:create`
- `media:read`
- `media:update`
- `media:delete`
- `media:execute`
- `media:approve`
- `media:assign_roles`
- `media:activate`
- `media:deactivate`
- `media:reorder`
- `media:publish`
- `media:close`
- `media:responses_read`
- `media:responses_delete`
- `media:enable`
- `media:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/media`
- `GET /admin/media/filter`
- `POST /admin/media/upload`
- `DELETE /admin/media/{id}`
- `GET /admin/media/{id}`
- `PUT /admin/media/{id}`
- `GET /public/media`
- `GET /public/media/filter`
- `POST /public/media/upload`
- `DELETE /public/media/{id}`
- `GET /public/media/{id}`
- `PUT /public/media/{id}`

Example Request
```http
GET /admin/media HTTP/1.1
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
