**Ads**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Ads

Purpose
Provides management and access endpoints for Ads resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `ads:create`
- `ads:read`
- `ads:update`
- `ads:delete`
- `ads:execute`
- `ads:approve`
- `ads:assign_roles`
- `ads:activate`
- `ads:deactivate`
- `ads:reorder`
- `ads:publish`
- `ads:close`
- `ads:responses_read`
- `ads:responses_delete`
- `ads:enable`
- `ads:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/ads`
- `POST /admin/ads`
- `DELETE /admin/ads/{id}`
- `GET /admin/ads/{id}`
- `PUT /admin/ads/{id}`
- `GET /public/ads`
- `POST /public/ads`
- `DELETE /public/ads/{id}`
- `GET /public/ads/{id}`
- `PUT /public/ads/{id}`

Example Request
```http
GET /admin/ads HTTP/1.1
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
