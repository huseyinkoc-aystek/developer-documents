**Marketplaces**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Marketplaces

Purpose
Provides management and access endpoints for Marketplaces resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `marketplaces:create`
- `marketplaces:read`
- `marketplaces:update`
- `marketplaces:delete`
- `marketplaces:execute`
- `marketplaces:approve`
- `marketplaces:assign_roles`
- `marketplaces:activate`
- `marketplaces:deactivate`
- `marketplaces:reorder`
- `marketplaces:publish`
- `marketplaces:close`
- `marketplaces:responses_read`
- `marketplaces:responses_delete`
- `marketplaces:enable`
- `marketplaces:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/marketplaces`
- `POST /admin/marketplaces`
- `GET /admin/marketplaces/webhooks`
- `GET /admin/marketplaces/webhooks/{id}`
- `DELETE /admin/marketplaces/{id}`
- `GET /admin/marketplaces/{id}`
- `PUT /admin/marketplaces/{id}`
- `GET /admin/marketplaces/{id}/categories`
- `GET /admin/marketplaces/{id}/categories/{categoryId}/attributes`
- `GET /admin/marketplaces/{id}/orders`
- `POST /admin/marketplaces/{id}/orders/status`
- `POST /admin/marketplaces/{id}/sync/catalog`
- `POST /admin/marketplaces/{id}/sync/price`
- `POST /admin/marketplaces/{id}/sync/product`
- `POST /admin/marketplaces/{id}/sync/products`
- `POST /admin/marketplaces/{id}/sync/stock`
- `GET /admin/marketplaces/{id}/test-connection`
- `POST /public/webhooks/marketplaces/{provider}`

Example Request
```http
GET /admin/marketplaces HTTP/1.1
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
