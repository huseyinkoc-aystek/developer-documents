**Carts**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Carts

Purpose
Provides management and access endpoints for Carts resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `carts:create`
- `carts:read`
- `carts:update`
- `carts:delete`
- `carts:execute`
- `carts:approve`
- `carts:assign_roles`
- `carts:activate`
- `carts:deactivate`
- `carts:reorder`
- `carts:publish`
- `carts:close`
- `carts:responses_read`
- `carts:responses_delete`
- `carts:enable`
- `carts:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `DELETE /cart`
- `GET /cart`
- `POST /cart/apply-coupon`
- `POST /cart/checkout`
- `POST /cart/checkout/draft`
- `POST /cart/checkout/preview`
- `POST /cart/checkout/quote`
- `POST /cart/currency`
- `POST /cart/items`
- `DELETE /cart/items/{variant_id}`
- `PATCH /cart/items/{variant_id}`
- `POST /cart/validate`

Example Request
```http
DELETE /cart HTTP/1.1
Host: api.example.com
Authorization: Bearer <access_token>
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
