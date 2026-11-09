**Subscriptions**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Subscriptions

Purpose
Provides management and access endpoints for Subscriptions resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `subscriptions:create`
- `subscriptions:read`
- `subscriptions:update`
- `subscriptions:delete`
- `subscriptions:execute`
- `subscriptions:approve`
- `subscriptions:assign_roles`
- `subscriptions:activate`
- `subscriptions:deactivate`
- `subscriptions:reorder`
- `subscriptions:publish`
- `subscriptions:close`
- `subscriptions:responses_read`
- `subscriptions:responses_delete`
- `subscriptions:enable`
- `subscriptions:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/subscriptions`
- `POST /admin/subscriptions`
- `DELETE /admin/subscriptions/{id}`
- `GET /admin/subscriptions/{id}`
- `PUT /admin/subscriptions/{id}`
- `PATCH /admin/subscriptions/{id}/status`

Example Request
```http
GET /admin/subscriptions HTTP/1.1
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
