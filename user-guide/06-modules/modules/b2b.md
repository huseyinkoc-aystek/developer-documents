**B2B**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/B2B

Purpose
Provides management and access endpoints for B2B resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `b2b_dealers:create`
- `b2b_dealers:read`
- `b2b_dealers:update`
- `b2b_dealers:delete`
- `b2b_dealers:execute`
- `b2b_dealers:approve`
- `b2b_dealers:assign_roles`
- `b2b_dealers:activate`
- `b2b_dealers:deactivate`
- `b2b_dealers:reorder`
- `b2b_dealers:publish`
- `b2b_dealers:close`
- `b2b_dealers:responses_read`
- `b2b_dealers:responses_delete`
- `b2b_dealers:enable`
- `b2b_dealers:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/b2b/dealers`
- `POST /admin/b2b/dealers`
- `DELETE /admin/b2b/dealers/{id}`
- `GET /admin/b2b/dealers/{id}`
- `PUT /admin/b2b/dealers/{id}`
- `GET /admin/b2b/dealers/{id}/credit`
- `POST /admin/b2b/dealers/{id}/credit`
- `GET /admin/b2b/dealers/{id}/pricing-rules`
- `POST /admin/b2b/dealers/{id}/pricing-rules`
- `DELETE /admin/b2b/dealers/{id}/pricing-rules/{rule_id}`
- `PUT /admin/b2b/dealers/{id}/pricing-rules/{rule_id}`
- `GET /admin/b2b/purchase-orders`
- `POST /admin/b2b/purchase-orders`
- `GET /admin/b2b/purchase-orders/{id}`
- `PUT /admin/b2b/purchase-orders/{id}`
- `POST /admin/b2b/purchase-orders/{id}/approve`
- `POST /admin/b2b/purchase-orders/{id}/cancel`
- `POST /admin/b2b/purchase-orders/{id}/convert`
- `POST /admin/b2b/purchase-orders/{id}/reject`
- `POST /admin/b2b/purchase-orders/{id}/submit`
- `POST /api/b2b-pricing/calculate`
- `GET /api/b2b-pricing/contracts`
- `POST /api/b2b-pricing/contracts`
- `DELETE /api/b2b-pricing/contracts/{id}`
- `GET /api/b2b-pricing/contracts/{id}`
- `PUT /api/b2b-pricing/contracts/{id}`
- `GET /api/b2b-pricing/dealer-tiers`
- `POST /api/b2b-pricing/dealer-tiers`
- `DELETE /api/b2b-pricing/dealer-tiers/{id}`
- `GET /api/b2b-pricing/dealer-tiers/{id}`
- `PUT /api/b2b-pricing/dealer-tiers/{id}`
- `GET /api/b2b-pricing/prices`
- `POST /api/b2b-pricing/prices`
- `DELETE /api/b2b-pricing/prices/{id}`
- `GET /api/b2b-pricing/prices/{id}`
- `PUT /api/b2b-pricing/prices/{id}`

Example Request
```http
GET /admin/b2b/dealers HTTP/1.1
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
