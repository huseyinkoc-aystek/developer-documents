**Billing**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Billing

Purpose
Provides management and access endpoints for Billing resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `billing:create`
- `billing:read`
- `billing:update`
- `billing:delete`
- `billing:execute`
- `billing:approve`
- `billing:assign_roles`
- `billing:activate`
- `billing:deactivate`
- `billing:reorder`
- `billing:publish`
- `billing:close`
- `billing:responses_read`
- `billing:responses_delete`
- `billing:enable`
- `billing:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/invoices`
- `POST /admin/invoices`
- `POST /admin/invoices/bulk/eissue`
- `GET /admin/invoices/eissue/errors`
- `DELETE /admin/invoices/{id}`
- `GET /admin/invoices/{id}`
- `PUT /admin/invoices/{id}`
- `POST /admin/invoices/{id}/eissue/cancel`
- `GET /admin/invoices/{id}/eissue/info`
- `POST /admin/invoices/{id}/eissue/retry`
- `GET /admin/invoices/{id}/eissue/status`
- `POST /admin/invoices/{id}/erp`
- `POST /admin/invoices/{id}/issue`
- `POST /admin/invoices/{id}/pdf`
- `GET /public/billing/accounts`
- `POST /public/billing/accounts`
- `DELETE /public/billing/accounts/{id}`
- `GET /public/billing/accounts/{id}`
- `GET /public/billing/balance`
- `GET /public/billing/charges`
- `POST /public/billing/charges`
- `GET /public/billing/charges/{id}`
- `PUT /public/billing/charges/{id}`
- `GET /public/billing/payments`
- `GET /public/billing/plans`
- `POST /public/billing/plans`
- `GET /public/billing/plans/{id}`
- `PUT /public/billing/plans/{id}`
- `POST /public/billing/plans/{id}/run`
- `GET /public/billing/rules`
- `POST /public/billing/rules`
- `GET /public/billing/rules/{id}`
- `PUT /public/billing/rules/{id}`
- `POST /public/billing/rules/{id}/distribute`
- `GET /public/billing/statement`
- `POST /public/billing/survey-decisions/apply`
- `GET /public/billing/units/{unit_id}/charges`
- `GET /public/billing/units/{unit_id}/payments`
- `GET /public/billing/units/{unit_id}/statement`
- `GET /public/invoices/{id}`

Example Request
```http
GET /admin/invoices HTTP/1.1
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
