**Erp Integrations**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Erp Integrations

Purpose
Provides management and access endpoints for Erp Integrations resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `erp_integrations:write`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/erp-integrations/config`
- `POST /admin/erp-integrations/config`
- `DELETE /admin/erp-integrations/config/{id}`
- `PUT /admin/erp-integrations/config/{id}`
- `GET /admin/erp-integrations/health`
- `POST /admin/erp-integrations/jobs/cancel`
- `GET /admin/erp-integrations/jobs/dlq`
- `POST /admin/erp-integrations/jobs/dlq/retry`
- `POST /admin/erp-integrations/jobs/enqueue`
- `GET /admin/erp-integrations/jobs/list`
- `POST /admin/erp-integrations/jobs/mapping`
- `GET /admin/erp-integrations/jobs/metrics`
- `GET /admin/erp-integrations/jobs/status/{job_id}`
- `GET /admin/erp-integrations/jobs/sync-state/{erp}/{scope}`
- `GET /admin/erp-integrations/logs`
- `GET /admin/erp-integrations/reconciliation`
- `POST /admin/erp-integrations/{erp}/trigger`

Example Request
```http
GET /admin/erp-integrations/config HTTP/1.1
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
