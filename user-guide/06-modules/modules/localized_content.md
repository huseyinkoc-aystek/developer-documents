**Localized Content**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Localized Content

Purpose
Provides management and access endpoints for Localized Content resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `localized_content:create`
- `localized_content:read`
- `localized_content:update`
- `localized_content:delete`
- `localized_content:execute`
- `localized_content:approve`
- `localized_content:assign_roles`
- `localized_content:activate`
- `localized_content:deactivate`
- `localized_content:reorder`
- `localized_content:publish`
- `localized_content:close`
- `localized_content:responses_read`
- `localized_content:responses_delete`
- `localized_content:enable`
- `localized_content:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/localized-content`
- `POST /admin/localized-content`
- `DELETE /admin/localized-content/{id}`
- `GET /admin/localized-content/{id}`
- `PUT /admin/localized-content/{id}`
- `GET /admin/localized-content/{language}`
- `GET /public/localized-content`
- `POST /public/localized-content`
- `DELETE /public/localized-content/{id}`
- `GET /public/localized-content/{id}`
- `PUT /public/localized-content/{id}`
- `GET /public/localized-content/{language}`

Example Request
```http
GET /admin/localized-content HTTP/1.1
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
