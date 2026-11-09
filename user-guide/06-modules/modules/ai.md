**Ai**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Ai

Purpose
Provides management and access endpoints for Ai resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `ai:create`
- `ai:read`
- `ai:update`
- `ai:delete`
- `ai:execute`
- `ai:approve`
- `ai:assign_roles`
- `ai:activate`
- `ai:deactivate`
- `ai:reorder`
- `ai:publish`
- `ai:close`
- `ai:responses_read`
- `ai:responses_delete`
- `ai:enable`
- `ai:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `POST /admin/ai/chat`
- `POST /admin/ai/events`
- `GET /admin/ai/insights`
- `GET /admin/ai/jobs`
- `POST /admin/ai/jobs`
- `GET /admin/ai/jobs/{id}`
- `GET /admin/ai/recommendations`
- `PUT /admin/ai/recommendations/{id}`
- `GET /api/analytics/aggregations`
- `GET /api/analytics/anomalies`
- `GET /api/analytics/conversions`
- `POST /api/analytics/events`
- `GET /api/analytics/insights`
- `GET /api/analytics/products/{productID}/metrics`
- `GET /api/analytics/revenue`
- `GET /api/analytics/trends`
- `GET /api/analytics/users/{userID}/events`

Example Request
```http
POST /admin/ai/chat HTTP/1.1
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
