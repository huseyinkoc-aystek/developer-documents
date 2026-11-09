**Chatbot**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Chatbot

Purpose
Provides management and access endpoints for Chatbot resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `chatbot:create`
- `chatbot:read`
- `chatbot:update`
- `chatbot:delete`
- `chatbot:execute`
- `chatbot:approve`
- `chatbot:assign_roles`
- `chatbot:activate`
- `chatbot:deactivate`
- `chatbot:reorder`
- `chatbot:publish`
- `chatbot:close`
- `chatbot:responses_read`
- `chatbot:responses_delete`
- `chatbot:enable`
- `chatbot:disable`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/chatbot/canned-responses`
- `POST /admin/chatbot/canned-responses`
- `DELETE /admin/chatbot/canned-responses/{id}`
- `GET /admin/chatbot/canned-responses/{id}`
- `PUT /admin/chatbot/canned-responses/{id}`
- `GET /admin/chatbot/intent-schemas`
- `POST /admin/chatbot/intent-schemas`
- `POST /admin/chatbot/intent-schemas/validate`
- `DELETE /admin/chatbot/intent-schemas/{id}`
- `GET /admin/chatbot/intent-schemas/{id}`
- `PUT /admin/chatbot/intent-schemas/{id}`
- `GET /admin/chatbot/intents`
- `GET /admin/chatbot/intents/{intentId}`
- `POST /public/chatbot/context`
- `GET /public/chatbot/health`
- `POST /public/chatbot/message`

Example Request
```http
GET /admin/chatbot/canned-responses HTTP/1.1
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
