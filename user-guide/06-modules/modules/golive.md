**Golive**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Golive

Purpose
Provides management and access endpoints for Golive resources.

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
- `POST /go-live/beta-customer`
- `GET /go-live/beta-customers`
- `GET /go-live/checklist`
- `PUT /go-live/checklist/{item_id}`
- `POST /go-live/incident`
- `GET /go-live/metrics`
- `GET /go-live/monitoring`
- `GET /go-live/risks`
- `GET /go-live/validate`

Example Request
```http
POST /go-live/beta-customer HTTP/1.1
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
