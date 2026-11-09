**Currencies**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Currencies

Purpose
Provides management and access endpoints for Currencies resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `currencies:create`
- `currencies:read`
- `currencies:update`
- `currencies:delete`
- `currencies:execute`
- `currencies:approve`
- `currencies:assign_roles`
- `currencies:activate`
- `currencies:deactivate`
- `currencies:reorder`
- `currencies:publish`
- `currencies:close`
- `currencies:responses_read`
- `currencies:responses_delete`
- `currencies:enable`
- `currencies:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/currencies`
- `POST /admin/currencies`
- `GET /admin/currencies/fx-rates`
- `POST /admin/currencies/fx-rates`
- `DELETE /admin/currencies/fx-rates/{id}`
- `GET /admin/currencies/fx-rates/{id}`
- `PUT /admin/currencies/fx-rates/{id}`
- `GET /admin/currencies/fx-sync/config`
- `PUT /admin/currencies/fx-sync/config`
- `POST /admin/currencies/fx-sync/run`
- `DELETE /admin/currencies/{id}`
- `PUT /admin/currencies/{id}`

Example Request
```http
GET /admin/currencies HTTP/1.1
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
