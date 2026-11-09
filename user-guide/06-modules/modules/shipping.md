**Shipping**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Shipping

Purpose
Provides management and access endpoints for Shipping resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `shipping:create`
- `shipping:read`
- `shipping:update`
- `shipping:delete`
- `shipping:execute`
- `shipping:approve`
- `shipping:assign_roles`
- `shipping:activate`
- `shipping:deactivate`
- `shipping:reorder`
- `shipping:publish`
- `shipping:close`
- `shipping:responses_read`
- `shipping:responses_delete`
- `shipping:enable`
- `shipping:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/shipping/carriers`
- `POST /admin/shipping/carriers`
- `DELETE /admin/shipping/carriers/{id}`
- `GET /admin/shipping/carriers/{id}`
- `PUT /admin/shipping/carriers/{id}`
- `GET /admin/shipping/delivery-options`
- `POST /admin/shipping/delivery-options`
- `DELETE /admin/shipping/delivery-options/{id}`
- `GET /admin/shipping/delivery-options/{id}`
- `PUT /admin/shipping/delivery-options/{id}`
- `POST /admin/shipping/quote`
- `GET /admin/shipping/shipments`
- `POST /admin/shipping/shipments`
- `GET /admin/shipping/shipments/{id}`
- `POST /admin/shipping/shipments/{id}/events`
- `POST /admin/shipping/shipments/{id}/label`
- `GET /public/shipping/track/{tracking_number}`

Example Request
```http
GET /admin/shipping/carriers HTTP/1.1
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
