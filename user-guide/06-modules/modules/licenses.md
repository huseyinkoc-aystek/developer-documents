**Licenses**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Licenses

Purpose
Provides management and access endpoints for Licenses resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `licenses:create`
- `licenses:read`
- `licenses:update`
- `licenses:delete`
- `licenses:execute`
- `licenses:approve`
- `licenses:assign_roles`
- `licenses:activate`
- `licenses:deactivate`
- `licenses:reorder`
- `licenses:publish`
- `licenses:close`
- `licenses:responses_read`
- `licenses:responses_delete`
- `licenses:enable`
- `licenses:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/licenses/access-grants`
- `POST /admin/licenses/keys`
- `GET /admin/licenses/pools`
- `POST /admin/licenses/pools`
- `DELETE /admin/licenses/pools/{id}`
- `PUT /admin/licenses/pools/{id}`
- `POST /public/licenses/access`

Example Request
```http
GET /admin/licenses/access-grants HTTP/1.1
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
