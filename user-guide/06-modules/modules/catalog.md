**Catalog**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Catalog

Purpose
Provides management and access endpoints for Catalog resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `catalog:create`
- `catalog:read`
- `catalog:update`
- `catalog:delete`
- `catalog:execute`
- `catalog:approve`
- `catalog:assign_roles`
- `catalog:activate`
- `catalog:deactivate`
- `catalog:reorder`
- `catalog:publish`
- `catalog:close`
- `catalog:responses_read`
- `catalog:responses_delete`
- `catalog:enable`
- `catalog:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/catalog`
- `POST /admin/catalog`
- `POST /admin/catalog/create`
- `POST /admin/catalog/delete`
- `GET /admin/catalog/detail`
- `GET /admin/catalog/list`
- `PUT /admin/catalog/update`
- `DELETE /admin/catalog/{id}`
- `GET /admin/catalog/{id}`
- `PUT /admin/catalog/{id}`
- `POST /admin/catalog/{id}/feature`
- `POST /admin/catalog/{id}/sort`
- `POST /admin/catalog/{id}/status`
- `GET /public/catalog`
- `POST /public/catalog`
- `POST /public/catalog/create`
- `POST /public/catalog/delete`
- `GET /public/catalog/detail`
- `GET /public/catalog/list`
- `PUT /public/catalog/update`
- `DELETE /public/catalog/{id}`
- `GET /public/catalog/{id}`
- `PUT /public/catalog/{id}`
- `POST /public/catalog/{id}/feature`
- `POST /public/catalog/{id}/sort`
- `POST /public/catalog/{id}/status`

Example Request
```http
GET /admin/catalog HTTP/1.1
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
