**Payments**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Payments

Purpose
Provides management and access endpoints for Payments resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `payments:create`
- `payments:read`
- `payments:update`
- `payments:delete`
- `payments:execute`
- `payments:approve`
- `payments:assign_roles`
- `payments:activate`
- `payments:deactivate`
- `payments:reorder`
- `payments:publish`
- `payments:close`
- `payments:responses_read`
- `payments:responses_delete`
- `payments:enable`
- `payments:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `POST /admin/installments/calculate`
- `GET /admin/installments/rules`
- `POST /admin/installments/rules`
- `DELETE /admin/installments/rules/{id}`
- `GET /admin/installments/rules/{id}`
- `PUT /admin/installments/rules/{id}`
- `GET /admin/payment-gateways`
- `POST /admin/payment-gateways`
- `GET /admin/payment-gateways/{id}`
- `PUT /admin/payment-gateways/{id}`
- `GET /admin/payment-intents`
- `POST /admin/payment-intents`
- `GET /admin/payment-intents/{id}`
- `POST /admin/payment-intents/{id}/3ds/complete`
- `POST /admin/payment-intents/{id}/cancel`
- `POST /admin/payment-intents/{id}/capture`
- `POST /admin/payment-intents/{id}/confirm`
- `POST /admin/payment-intents/{id}/refunds`
- `GET /admin/payment-links`
- `POST /admin/payment-links`
- `GET /admin/payment-links/{id}`
- `PUT /admin/payment-links/{id}`
- `GET /admin/payment-risk/events`
- `GET /admin/payment-risk/report`
- `GET /admin/payments`
- `POST /admin/payments`
- `DELETE /admin/payments/{id}`
- `GET /admin/payments/{id}`
- `PUT /admin/payments/{id}`
- `POST /admin/payments/{id}/capture`
- `POST /admin/payments/{id}/refund`
- `POST /public/payment-intents/callbacks/{provider}`
- `GET /public/payment-links/{token}`
- `POST /public/payment-links/{token}/3ds/complete`
- `POST /public/payment-links/{token}/capture`
- `POST /public/payment-links/{token}/checkout`
- `POST /public/payment-links/{token}/confirm`
- `GET /public/payments`
- `POST /public/payments`
- `DELETE /public/payments/{id}`
- `GET /public/payments/{id}`
- `PUT /public/payments/{id}`
- `POST /public/payments/{id}/capture`
- `POST /public/payments/{id}/refund`

Example Request
```http
POST /admin/installments/calculate HTTP/1.1
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
