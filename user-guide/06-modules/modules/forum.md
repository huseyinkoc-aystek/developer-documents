**Forum**
Relevant API
Refer to: https://<BASE_URL>/api-docs/swagger.yaml#/tag/Forum

Purpose
Provides management and access endpoints for Forum resources.

Dependencies
- Tenant and site context (TenantIsolation + SiteIsolation middleware).
- ModuleGate and PermissionMiddleware (admin and protected surfaces).
- MongoDB persistence for module data.
- Audit logging for admin actions where configured.

Permissions Required
- `forum:create`
- `forum:read`
- `forum:update`
- `forum:delete`
- `forum:execute`
- `forum:approve`
- `forum:assign_roles`
- `forum:activate`
- `forum:deactivate`
- `forum:reorder`
- `forum:publish`
- `forum:close`
- `forum:responses_read`
- `forum:responses_delete`
- `forum:enable`
- `forum:disable`
- `forum:sponsor.manage`
- `forum:ads.manage`
- `forum:reputation.manage`
- `forum:badges.manage`

Multi-Tenant Behavior
All requests run through tenant resolution and isolation. Requests without tenant context are rejected.

Site-Level Behavior
All requests run through site resolution and isolation. Requests without site context are rejected.

Permission Gating
ModuleGate and RBAC permissions are enforced server-side before any business logic executes.

Plugin Extension Points
Plugins may request capabilities that map to this module’s permissions, but execution is still gated by ModuleGate, RBAC, license, and feature flags.

API Endpoints
- `GET /admin/forum/access-plans`
- `POST /admin/forum/access-plans`
- `GET /admin/forum/ads`
- `POST /admin/forum/ads`
- `DELETE /admin/forum/ads/{id}`
- `GET /admin/forum/ai-decisions`
- `PUT /admin/forum/ai-decisions/{id}/approve`
- `PUT /admin/forum/ai-decisions/{id}/reject`
- `PUT /admin/forum/ai-decisions/{id}/revert`
- `GET /admin/forum/ai-policies`
- `POST /admin/forum/ai-policies`
- `GET /admin/forum/ai-policies/{id}`
- `GET /admin/forum/analytics`
- `GET /admin/forum/badge-rules`
- `POST /admin/forum/badge-rules`
- `DELETE /admin/forum/badge-rules/{id}`
- `GET /admin/forum/badges`
- `POST /admin/forum/badges`
- `GET /admin/forum/badges/users`
- `DELETE /admin/forum/badges/{id}`
- `POST /admin/forum/badges/{id}/grant`
- `POST /admin/forum/badges/{id}/revoke`
- `GET /admin/forum/bans`
- `POST /admin/forum/bans`
- `PUT /admin/forum/bans/{id}/revoke`
- `GET /admin/forum/categories`
- `POST /admin/forum/categories`
- `DELETE /admin/forum/categories/{id}`
- `GET /admin/forum/categories/{id}`
- `PUT /admin/forum/categories/{id}`
- `GET /admin/forum/config`
- `GET /admin/forum/forums`
- `POST /admin/forum/forums`
- `DELETE /admin/forum/forums/{id}`
- `GET /admin/forum/forums/{id}`
- `PUT /admin/forum/forums/{id}`
- `POST /admin/forum/indexes/verify`
- `POST /admin/forum/indexes/verify-all`
- `GET /admin/forum/moderation-actions`
- `GET /admin/forum/posts`
- `POST /admin/forum/posts`
- `DELETE /admin/forum/posts/{id}`
- `GET /admin/forum/posts/{id}`
- `PUT /admin/forum/posts/{id}`
- `GET /admin/forum/reports`
- `PUT /admin/forum/reports/{id}/assign`
- `PUT /admin/forum/reports/{id}/status`
- `GET /admin/forum/reputation/ledger`
- `GET /admin/forum/reputation/rules`
- `POST /admin/forum/reputation/rules`
- `DELETE /admin/forum/reputation/rules/{id}`
- `GET /admin/forum/subscriptions`
- `POST /admin/forum/subscriptions`
- `GET /admin/forum/threads`
- `POST /admin/forum/threads`
- `DELETE /admin/forum/threads/{id}`
- `GET /admin/forum/threads/{id}`
- `PUT /admin/forum/threads/{id}`
- `POST /admin/forum/threads/{id}/merge`
- `PUT /admin/forum/threads/{id}/move`
- `POST /admin/forum/threads/{id}/split`
- `PUT /admin/forum/threads/{id}/sponsor`
- `PUT /admin/forum/threads/{id}/status`
- `PUT /admin/forum/threads/{id}/type`
- `GET /public/forums/ads`
- `GET /public/forums/badges`
- `GET /public/forums/badges/users/{user_id}`
- `GET /public/forums/categories`
- `GET /public/forums/forums`
- `GET /public/forums/forums/{id}`
- `GET /public/forums/pm/conversations`
- `POST /public/forums/pm/conversations`
- `GET /public/forums/pm/conversations/{id}/messages`
- `POST /public/forums/pm/conversations/{id}/messages`
- `POST /public/forums/pm/conversations/{id}/participants`
- `DELETE /public/forums/pm/conversations/{id}/participants/{user_id}`
- `GET /public/forums/posts`
- `POST /public/forums/posts`
- `DELETE /public/forums/posts/{id}/reactions`
- `GET /public/forums/posts/{id}/reactions`
- `POST /public/forums/posts/{id}/reactions`
- `POST /public/forums/reports`
- `GET /public/forums/reputation/{user_id}`
- `GET /public/forums/search/posts`
- `GET /public/forums/search/threads`
- `POST /public/forums/subscriptions`
- `GET /public/forums/threads`
- `POST /public/forums/threads`
- `GET /public/forums/threads/{id}`

Example Request
```http
GET /admin/forum/access-plans HTTP/1.1
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
