**Reports**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Reports

Purpose
Provides management and access endpoints for Reports resources.

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
- `GET /executions/{id}`
- `GET /report-templates`
- `POST /report-templates/{id}/create`
- `GET /reports`
- `POST /reports`
- `POST /reports/preview`
- `POST /reports/validate-query`
- `DELETE /reports/{id}`
- `GET /reports/{id}`
- `PUT /reports/{id}`
- `GET /reports/{id}/analytics`
- `POST /reports/{id}/duplicate`
- `POST /reports/{id}/execute`
- `POST /reports/{id}/execute-async`
- `GET /reports/{id}/export`
- `POST /reports/{id}/share`
- `POST /reports/{id}/subscribe`
- `DELETE /subscriptions/{id}`

Example Request
```http
GET /executions/{id} HTTP/1.1
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
