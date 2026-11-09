**Activity Logs**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Activity Logs

Purpose
Provides management and access endpoints for Activity Logs resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `activity_logs:create`
- `activity_logs:read`
- `activity_logs:update`
- `activity_logs:delete`
- `activity_logs:execute`
- `activity_logs:approve`
- `activity_logs:assign_roles`
- `activity_logs:activate`
- `activity_logs:deactivate`
- `activity_logs:reorder`
- `activity_logs:publish`
- `activity_logs:close`
- `activity_logs:responses_read`
- `activity_logs:responses_delete`
- `activity_logs:enable`
- `activity_logs:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/activity-logs`

Example Request
```http
GET /admin/activity-logs HTTP/1.1
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
