**Mailtemplates**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Mailtemplates

Purpose
Provides management and access endpoints for Mailtemplates resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `mailtemplates:create`
- `mailtemplates:read`
- `mailtemplates:update`
- `mailtemplates:delete`
- `mailtemplates:execute`
- `mailtemplates:approve`
- `mailtemplates:assign_roles`
- `mailtemplates:activate`
- `mailtemplates:deactivate`
- `mailtemplates:reorder`
- `mailtemplates:publish`
- `mailtemplates:close`
- `mailtemplates:responses_read`
- `mailtemplates:responses_delete`
- `mailtemplates:enable`
- `mailtemplates:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/mail-templates`
- `POST /admin/mail-templates`
- `GET /admin/mail-templates/{id}`
- `PUT /admin/mail-templates/{id}`

Example Request
```http
GET /admin/mail-templates HTTP/1.1
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
