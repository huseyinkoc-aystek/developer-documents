**Taxonomy**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Taxonomy

Purpose
Provides management and access endpoints for Taxonomy resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `taxonomy:create`
- `taxonomy:read`
- `taxonomy:update`
- `taxonomy:delete`
- `taxonomy:execute`
- `taxonomy:approve`
- `taxonomy:assign_roles`
- `taxonomy:activate`
- `taxonomy:deactivate`
- `taxonomy:reorder`
- `taxonomy:publish`
- `taxonomy:close`
- `taxonomy:responses_read`
- `taxonomy:responses_delete`
- `taxonomy:enable`
- `taxonomy:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/categories`
- `POST /admin/categories`
- `POST /admin/categories/bulk/delete`
- `POST /admin/categories/bulk/locales/add`
- `POST /admin/categories/bulk/locales/remove`
- `POST /admin/categories/reorder`
- `GET /admin/categories/slug/{slug}`
- `GET /admin/categories/tree`
- `DELETE /admin/categories/{id}`
- `GET /admin/categories/{id}`
- `PUT /admin/categories/{id}`
- `GET /admin/tags`
- `POST /admin/tags`
- `POST /admin/tags/bulk/delete`
- `POST /admin/tags/bulk/locale/add`
- `POST /admin/tags/bulk/locale/remove`
- `DELETE /admin/tags/{id}`
- `GET /admin/tags/{id}`
- `PUT /admin/tags/{id}`
- `GET /public/categories`
- `POST /public/categories`
- `POST /public/categories/bulk/delete`
- `POST /public/categories/bulk/locales/add`
- `POST /public/categories/bulk/locales/remove`
- `POST /public/categories/reorder`
- `GET /public/categories/slug/{slug}`
- `GET /public/categories/tree`
- `DELETE /public/categories/{id}`
- `GET /public/categories/{id}`
- `PUT /public/categories/{id}`
- `GET /public/tags`
- `POST /public/tags`
- `POST /public/tags/bulk/delete`
- `POST /public/tags/bulk/locale/add`
- `POST /public/tags/bulk/locale/remove`
- `DELETE /public/tags/{id}`
- `GET /public/tags/{id}`
- `PUT /public/tags/{id}`

Example Request
```http
GET /admin/categories HTTP/1.1
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
