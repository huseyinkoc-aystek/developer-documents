**Notifications**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Notifications

Purpose
Provides management and access endpoints for Notifications resources.

Domain Integrations
- Wallet bank-transfer declarations use the existing notification pipeline; no parallel inbox is maintained.
- After `POST /me/wallet/bank-transfer-declarations` creates a new `pending_reconciliation` journal, `wallets.Service` targets the wallet owner's `Person.LinkedUserID` and calls `notifications.Service.Create`.
- The resulting `notification.created` event produces the user receipt/read-model entry consumed by `GET /notifications`; configured push delivery continues through the same FCM/APNS path.
- Inbox event handling preserves tenant/site values but detaches from HTTP request cancellation with a bounded 15-second worker context.
- Repeating the same `(wallet_id, reference)` declaration is idempotent and does not create a second inbox notification.
- The bank-transfer message includes amount, currency, and reference in all 28 platform languages.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `notifications:read`
- `notifications:create`
- `notifications:update`
- `notifications:stream`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/notifications`
- `POST /admin/notifications/create`
- `GET /admin/notifications/list`
- `GET /admin/notifications/me`
- `GET /admin/notifications/preferences`
- `PUT /admin/notifications/preferences`
- `GET /admin/notifications/push`
- `GET /admin/notifications/push-tokens`
- `POST /admin/notifications/push-tokens`
- `POST /admin/notifications/push-tokens/revoke-all`
- `POST /admin/notifications/push-tokens/{id}/revoke`
- `POST /admin/notifications/push/send`
- `POST /admin/notifications/push/test`
- `POST /admin/notifications/read-all`
- `GET /admin/notifications/receipts`
- `DELETE /admin/notifications/{id}`
- `GET /admin/notifications/{id}`
- `PUT /admin/notifications/{id}`
- `PATCH /admin/notifications/{id}/read`
- `POST /admin/notifications/{id}/read`
- `GET /admin/push-notifications`
- `POST /admin/push-notifications/send`
- `POST /admin/push-notifications/test`
- `GET /notifications`
- `POST /notifications/create`
- `GET /notifications/list`
- `GET /notifications/me`
- `GET /notifications/preferences`
- `PUT /notifications/preferences`
- `GET /notifications/push`
- `POST /notifications/push/send`
- `POST /notifications/push/test`
- `POST /notifications/read-all`
- `GET /notifications/receipts`
- `GET /notifications/stream`
- `DELETE /notifications/{id}`
- `GET /notifications/{id}`
- `PUT /notifications/{id}`
- `PATCH /notifications/{id}/read`
- `POST /notifications/{id}/read`
- `GET /public/notifications`
- `POST /public/notifications/create`
- `GET /public/notifications/list`
- `GET /public/notifications/me`
- `GET /public/notifications/preferences`
- `PUT /public/notifications/preferences`
- `GET /public/notifications/push`
- `POST /public/notifications/push/send`
- `POST /public/notifications/push/test`
- `POST /public/notifications/read-all`
- `GET /public/notifications/receipts`
- `DELETE /public/notifications/{id}`
- `GET /public/notifications/{id}`
- `PUT /public/notifications/{id}`
- `PATCH /public/notifications/{id}/read`
- `POST /public/notifications/{id}/read`

Example Request
```http
GET /admin/notifications HTTP/1.1
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
