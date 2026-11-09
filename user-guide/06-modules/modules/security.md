**Security**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Security

Purpose
Provides management and access endpoints for Security resources.

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
- `POST /auth/login`
- `POST /security/account/unlock`
- `GET /security/audit-log`
- `POST /security/compliance/gdpr`
- `POST /security/compliance/hipaa`
- `POST /security/compliance/pci-dss`
- `GET /security/compliance/status`
- `GET /security/encryption/status`
- `GET /security/headers`
- `GET /security/high-severity-events`
- `POST /security/logout`
- `GET /security/vulnerabilities/open`
- `POST /security/vulnerabilities/scan`
- `POST /security/whitelist/manage`
- `GET /security/whitelist/status`

Example Request
```http
POST /auth/login HTTP/1.1
Host: api.example.com
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
