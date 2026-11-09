**Contracts**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Contracts

Purpose
Provides management and access endpoints for Contracts resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `contracts:create`
- `contracts:read`
- `contracts:update`
- `contracts:delete`
- `contracts:execute`
- `contracts:approve`
- `contracts:assign_roles`
- `contracts:activate`
- `contracts:deactivate`
- `contracts:reorder`
- `contracts:publish`
- `contracts:close`
- `contracts:responses_read`
- `contracts:responses_delete`
- `contracts:enable`
- `contracts:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/contracts`
- `POST /admin/contracts`
- `DELETE /admin/contracts/{id}`
- `GET /admin/contracts/{id}`
- `PUT /admin/contracts/{id}`
- `GET /admin/contracts/{id}/contents`
- `POST /admin/contracts/{id}/contents`
- `GET /admin/contracts/{id}/contents/{contentId}`
- `PUT /admin/contracts/{id}/contents/{contentId}`
- `GET /admin/orders/{id}/contracts`
- `POST /admin/orders/{id}/contracts`

Example Request
```http
GET /admin/contracts HTTP/1.1
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
