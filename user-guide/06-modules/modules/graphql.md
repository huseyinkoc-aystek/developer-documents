**Graphql**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Graphql

Purpose
Provides management and access endpoints for Graphql resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `graphql:create`
- `graphql:read`
- `graphql:update`
- `graphql:delete`
- `graphql:execute`
- `graphql:approve`
- `graphql:assign_roles`
- `graphql:activate`
- `graphql:deactivate`
- `graphql:reorder`
- `graphql:publish`
- `graphql:close`
- `graphql:responses_read`
- `graphql:responses_delete`
- `graphql:enable`
- `graphql:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `POST /graphql`
- `POST /graphql/`
- `GET /graphql/stats`

Example Request
```http
POST /graphql HTTP/1.1
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
