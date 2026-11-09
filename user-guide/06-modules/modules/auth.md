**Auth**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Auth

Purpose
Provides management and access endpoints for Auth resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `auth:create`
- `auth:read`
- `auth:update`
- `auth:delete`
- `auth:execute`
- `auth:approve`
- `auth:assign_roles`
- `auth:activate`
- `auth:deactivate`
- `auth:reorder`
- `auth:publish`
- `auth:close`
- `auth:responses_read`
- `auth:responses_delete`
- `auth:enable`
- `auth:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /auth/csrf`
- `POST /auth/login-by-email`
- `POST /auth/login-by-phone`
- `POST /auth/login-by-username`
- `POST /auth/logout`
- `POST /auth/refresh`
- `POST /auth/register`
- `POST /auth/request-password-reset`
- `POST /auth/reset-password`
- `POST /auth/send-verification/{userID}`
- `POST /auth/verify-email/{token}`

Example Request
```http
GET /auth/csrf HTTP/1.1
Host: api.example.com
X-Tenant-Code: acme
X-Store-Code: acme-store
```

Example Response
```json
{}
```

Example Error Responses
```json
{
  "success": false,
  "message": "bad_request"
}
```
