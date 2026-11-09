**Discounts**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Discounts

Purpose
Provides management and access endpoints for Discounts resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `discounts:create`
- `discounts:read`
- `discounts:update`
- `discounts:delete`
- `discounts:execute`
- `discounts:approve`
- `discounts:assign_roles`
- `discounts:activate`
- `discounts:deactivate`
- `discounts:reorder`
- `discounts:publish`
- `discounts:close`
- `discounts:responses_read`
- `discounts:responses_delete`
- `discounts:enable`
- `discounts:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/coupons`
- `POST /admin/coupons`
- `DELETE /admin/coupons/{id}`
- `GET /admin/coupons/{id}`
- `PUT /admin/coupons/{id}`
- `GET /admin/discounts`
- `POST /admin/discounts`
- `DELETE /admin/discounts/{id}`
- `GET /admin/discounts/{id}`
- `PUT /admin/discounts/{id}`

Example Request
```http
GET /admin/coupons HTTP/1.1
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
