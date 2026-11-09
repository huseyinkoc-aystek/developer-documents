**Posts**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Posts

Purpose
Provides management and access endpoints for Posts resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `posts:create`
- `posts:read`
- `posts:update`
- `posts:delete`
- `posts:execute`
- `posts:approve`
- `posts:assign_roles`
- `posts:activate`
- `posts:deactivate`
- `posts:reorder`
- `posts:publish`
- `posts:close`
- `posts:responses_read`
- `posts:responses_delete`
- `posts:enable`
- `posts:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/posts`
- `POST /admin/posts`
- `DELETE /admin/posts/by-id/{id}`
- `GET /admin/posts/by-id/{id}`
- `PUT /admin/posts/by-id/{id}`
- `POST /admin/posts/by-id/{id}/view`
- `GET /admin/posts/filter`
- `GET /admin/posts/lang/{lang}`
- `GET /admin/posts/localized/{lang}/{slug}`
- `GET /admin/posts/slug/{slug}`
- `PATCH /admin/posts/{id}/sort`
- `GET /admin/posts/{lang}/{slug}`
- `GET /public/posts`
- `POST /public/posts`
- `DELETE /public/posts/by-id/{id}`
- `GET /public/posts/by-id/{id}`
- `PUT /public/posts/by-id/{id}`
- `POST /public/posts/by-id/{id}/view`
- `GET /public/posts/filter`
- `GET /public/posts/lang/{lang}`
- `GET /public/posts/localized/{lang}/{slug}`
- `GET /public/posts/slug/{slug}`
- `PATCH /public/posts/{id}/sort`
- `GET /public/posts/{lang}/{slug}`

Example Request
```http
GET /admin/posts HTTP/1.1
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
