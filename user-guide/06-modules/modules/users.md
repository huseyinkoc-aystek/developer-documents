**Users**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Users

Purpose
Provides management and access endpoints for Users resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `users:create`
- `users:read`
- `users:update`
- `users:delete`
- `users:execute`
- `users:approve`
- `users:assign_roles`
- `users:activate`
- `users:deactivate`
- `users:reorder`
- `users:publish`
- `users:close`
- `users:responses_read`
- `users:responses_delete`
- `users:enable`
- `users:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/customers/{customerID}/addresses`
- `POST /admin/customers/{customerID}/addresses`
- `DELETE /admin/customers/{customerID}/addresses/{addressID}`
- `GET /admin/customers/{customerID}/addresses/{addressID}`
- `PATCH /admin/customers/{customerID}/addresses/{addressID}`
- `GET /admin/users`
- `POST /admin/users/create`
- `PUT /admin/users/language`
- `DELETE /admin/users/{id}`
- `GET /admin/users/{id}`
- `PUT /admin/users/{id}`
- `PATCH /admin/users/{id}/approve`
- `PATCH /admin/users/{id}/roles`

Example Request
```http
GET /admin/customers/{customerID}/addresses HTTP/1.1
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
