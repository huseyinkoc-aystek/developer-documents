**Products**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Products

Purpose
Provides management and access endpoints for Products resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `products:create`
- `products:read`
- `products:update`
- `products:delete`
- `products:execute`
- `products:approve`
- `products:assign_roles`
- `products:activate`
- `products:deactivate`
- `products:reorder`
- `products:publish`
- `products:close`
- `products:responses_read`
- `products:responses_delete`
- `products:enable`
- `products:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/products`
- `POST /admin/products`
- `DELETE /admin/products/{id}`
- `GET /admin/products/{id}`
- `PUT /admin/products/{id}`
- `GET /admin/products/{id}/attributes`
- `POST /admin/products/{id}/attributes`
- `GET /admin/products/{id}/attributes/categories/{category_id}/attributes`
- `POST /admin/products/{id}/attributes/categories/{category_id}/attributes`
- `DELETE /admin/products/{id}/attributes/categories/{category_id}/attributes/{attribute_id}`
- `GET /admin/products/{id}/attributes/values`
- `POST /admin/products/{id}/attributes/values`
- `DELETE /admin/products/{id}/attributes/values/{value_id}`
- `PUT /admin/products/{id}/attributes/values/{value_id}`
- `DELETE /admin/products/{id}/attributes/{id}`
- `GET /admin/products/{id}/attributes/{id}`
- `PUT /admin/products/{id}/attributes/{id}`
- `GET /admin/products/{id}/contents`
- `POST /admin/products/{id}/contents`
- `GET /admin/products/{id}/contents/{contentId}`
- `PUT /admin/products/{id}/contents/{contentId}`
- `GET /admin/products/{id}/media`
- `POST /admin/products/{id}/media`
- `GET /admin/products/{id}/media/gallery`
- `GET /admin/products/{id}/media/main`
- `POST /admin/products/{id}/media/reorder`
- `DELETE /admin/products/{id}/media/{media_id}`
- `GET /admin/products/{id}/media/{media_id}`
- `PUT /admin/products/{id}/media/{media_id}`
- `POST /admin/products/{id}/pricing`
- `GET /admin/products/{id}/pricing/active`
- `GET /admin/products/{id}/pricing/history`
- `PUT /admin/products/{id}/pricing/{price_id}`
- `POST /admin/products/{id}/pricing/{price_id}/expire`
- `GET /admin/products/{id}/variants`
- `POST /admin/products/{id}/variants`
- `POST /admin/products/{id}/variants/bulk`
- `DELETE /admin/products/{id}/variants/{variant_id}`
- `PUT /admin/products/{id}/variants/{variant_id}`
- `GET /api/v1/products/{productID}/reviews`
- `POST /api/v1/products/{productID}/reviews`
- `GET /api/v1/products/{productID}/reviews/statistics`
- `DELETE /api/v1/reviews/{reviewID}`
- `GET /api/v1/reviews/{reviewID}`
- `PUT /api/v1/reviews/{reviewID}`
- `POST /api/v1/reviews/{reviewID}/helpful`
- `POST /api/v1/reviews/{reviewID}/response`
- `PATCH /api/v1/reviews/{reviewID}/status`
- `POST /api/v1/reviews/{reviewID}/unhelpful`
- `GET /bundles`
- `POST /bundles`
- `GET /bundles/by-product/{product_id}`
- `GET /bundles/slug/{slug}`
- `DELETE /bundles/{id}`
- `GET /bundles/{id}`
- `PUT /bundles/{id}`
- `POST /bundles/{id}/check-availability`
- `GET /bundles/{id}/stats`
- `PATCH /bundles/{id}/status`
- `PATCH /bundles/{id}/stock`
- `POST /bundles/{id}/view`

Example Request
```http
GET /admin/products HTTP/1.1
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
