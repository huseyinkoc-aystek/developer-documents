**Crm**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Crm

Purpose
Provides management and access endpoints for Crm resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- No module-specific permission entries in `permission_seed.go`. Access still passes through global permission checks.

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/crm/contacts`
- `GET /admin/crm/deals`
- `GET /admin/crm/integrations`
- `POST /admin/crm/integrations`
- `GET /admin/crm/integrations/{id}`
- `POST /admin/crm/integrations/{integrationId}/sync/contacts`
- `POST /admin/crm/integrations/{integrationId}/sync/deals`
- `POST /admin/crm/integrations/{integrationId}/sync/leads`
- `GET /admin/crm/leads`
- `POST /public/webhooks/hubspot`
- `POST /public/webhooks/salesforce`
- `GET /public/webhooks/status`

Example Request
```http
GET /admin/crm/contacts HTTP/1.1
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
