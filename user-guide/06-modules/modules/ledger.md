**Ledger**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Ledger

Purpose
Provides management and access endpoints for Ledger resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `ledger:create`
- `ledger:read`
- `ledger:update`
- `ledger:delete`
- `ledger:execute`
- `ledger:approve`
- `ledger:assign_roles`
- `ledger:activate`
- `ledger:deactivate`
- `ledger:reorder`
- `ledger:publish`
- `ledger:close`
- `ledger:responses_read`
- `ledger:responses_delete`
- `ledger:enable`
- `ledger:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/ledger/accounts`
- `POST /admin/ledger/accounts`
- `DELETE /admin/ledger/accounts/{id}`
- `GET /admin/ledger/accounts/{id}`
- `GET /admin/ledger/journals`
- `POST /admin/ledger/journals`
- `GET /admin/ledger/journals/{id}`

Example Request
```http
GET /admin/ledger/accounts HTTP/1.1
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
