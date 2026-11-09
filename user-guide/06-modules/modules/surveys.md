**Surveys**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Surveys

Purpose
Provides management and access endpoints for Surveys resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `surveys:create`
- `surveys:read`
- `surveys:update`
- `surveys:delete`
- `surveys:execute`
- `surveys:approve`
- `surveys:assign_roles`
- `surveys:activate`
- `surveys:deactivate`
- `surveys:reorder`
- `surveys:publish`
- `surveys:close`
- `surveys:responses_read`
- `surveys:responses_delete`
- `surveys:enable`
- `surveys:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/surveys`
- `POST /admin/surveys`
- `DELETE /admin/surveys/{id}`
- `GET /admin/surveys/{id}`
- `POST /admin/surveys/{id}/close`
- `POST /admin/surveys/{id}/publish`
- `GET /admin/surveys/{slug}`
- `GET /public/surveys`
- `POST /public/surveys`
- `DELETE /public/surveys/{id}`
- `GET /public/surveys/{id}`
- `POST /public/surveys/{id}/close`
- `POST /public/surveys/{id}/publish`
- `GET /public/surveys/{slug}`

Example Request
```http
GET /admin/surveys HTTP/1.1
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
