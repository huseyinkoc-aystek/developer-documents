**Languages**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Languages

Purpose
Provides management and access endpoints for Languages resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `languages:create`
- `languages:read`
- `languages:update`
- `languages:delete`
- `languages:execute`
- `languages:approve`
- `languages:assign_roles`
- `languages:activate`
- `languages:deactivate`
- `languages:reorder`
- `languages:publish`
- `languages:close`
- `languages:responses_read`
- `languages:responses_delete`
- `languages:enable`
- `languages:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Site-scoped languages
The tenant owns the language catalog. Each site may narrow it with its own selection (ordered language codes plus a default), stored in `site_language_configs`. A site without a selection follows the tenant-wide active languages. `GET /public/languages/active` and `GET /public/languages/default` return the selected site's result, so public websites never need a hard-coded language list. `GET /public/languages` (paged list) stays tenant-wide.

Admin screens: `/admin/languages` (tenant catalog, "on this site" column) and `/admin/languages/site` (selection for the currently selected site; menu key `site-languages`). Selection rules: at least one language, every code must be an active tenant language, the default must be one of the selected languages.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/languages`
- `POST /admin/languages`
- `GET /admin/languages/active`
- `GET /admin/languages/default`
- `GET /admin/languages/summary`
- `DELETE /admin/languages/{id}`
- `GET /admin/languages/{id}`
- `PUT /admin/languages/{id}`
- `POST /admin/languages/{id}/activate`
- `POST /admin/languages/{id}/deactivate`
- `POST /admin/languages/{id}/default`
- `POST /admin/languages/{id}/order`
- `GET /admin/languages/site`
- `PUT /admin/languages/site`
- `DELETE /admin/languages/site`
- `GET /languages`
- `GET /languages/`
- `GET /languages/active`
- `GET /languages/default`
- `GET /public/languages`
- `POST /public/languages`
- `GET /public/languages/active`
- `GET /public/languages/default`
- `GET /public/languages/summary`
- `DELETE /public/languages/{id}`
- `GET /public/languages/{id}`
- `PUT /public/languages/{id}`
- `POST /public/languages/{id}/activate`
- `POST /public/languages/{id}/deactivate`
- `POST /public/languages/{id}/default`
- `POST /public/languages/{id}/order`

Example Request
```http
GET /admin/languages HTTP/1.1
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
