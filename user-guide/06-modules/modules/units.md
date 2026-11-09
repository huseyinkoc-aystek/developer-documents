**Units**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Units

Purpose
Provides management and access endpoints for Units resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `units:create`
- `units:read`
- `units:update`
- `units:delete`
- `units:execute`
- `units:approve`
- `units:assign_roles`
- `units:activate`
- `units:deactivate`
- `units:reorder`
- `units:publish`
- `units:close`
- `units:responses_read`
- `units:responses_delete`
- `units:enable`
- `units:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/units`
- `POST /admin/units`
- `DELETE /admin/units/{id}`
- `GET /admin/units/{id}`
- `PUT /admin/units/{id}`
- `GET /public/units`
- `POST /public/units`
- `DELETE /public/units/{id}`
- `GET /public/units/{id}`
- `PUT /public/units/{id}`

Example Request
```http
GET /admin/units HTTP/1.1
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
