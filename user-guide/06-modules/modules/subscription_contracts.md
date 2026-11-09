**Subscription Contracts**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Subscription Contracts

Purpose
Provides management and access endpoints for Subscription Contracts resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `subscription_contracts:create`
- `subscription_contracts:read`
- `subscription_contracts:update`
- `subscription_contracts:delete`
- `subscription_contracts:execute`
- `subscription_contracts:approve`
- `subscription_contracts:assign_roles`
- `subscription_contracts:activate`
- `subscription_contracts:deactivate`
- `subscription_contracts:reorder`
- `subscription_contracts:publish`
- `subscription_contracts:close`
- `subscription_contracts:responses_read`
- `subscription_contracts:responses_delete`
- `subscription_contracts:enable`
- `subscription_contracts:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/subscription-contracts`
- `POST /admin/subscription-contracts`
- `GET /admin/subscription-contracts/{id}`
- `PUT /admin/subscription-contracts/{id}`
- `PATCH /admin/subscription-contracts/{id}/status`

Example Request
```http
GET /admin/subscription-contracts HTTP/1.1
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
