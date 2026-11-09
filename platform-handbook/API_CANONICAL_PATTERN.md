**Canonical REST Pattern**

The platform will normalize to a single, consistent REST pattern without changing business logic. This is a planning document for the normalization sprint.

**Canonical CRUD Pattern (Admin)**
```text
GET    /admin/{resource}           -> list
POST   /admin/{resource}           -> create
GET    /admin/{resource}/{id}      -> get
PUT    /admin/{resource}/{id}      -> update (full)
PATCH  /admin/{resource}/{id}      -> update (partial)
DELETE /admin/{resource}/{id}      -> delete
```

**Canonical Read Pattern (Public)**
```text
GET    /public/{resource}          -> list
GET    /public/{resource}/{id|slug} -> get
```

**Canonical Rules**
- No trailing slashes.
- No `/by-id` variants (use `/{id}`).
- No root duplicates for admin/public (use the correct prefix only).
- Avoid mixing `/api/admin`, `/api/v1`, `/api/analytics` with `/admin` and `/public` unless versioning is explicitly required by contract.

**Non-Canonical Endpoint Categories (Detected)**

| Category | Count | Example |
| --- | --- | --- |
| Trailing slash variants | 10 | `/admin/authors/` |
| /by-id variants | 9 | `/admin/authors/by-id/{id}` |
| /api/admin prefix | 19 | `/api/admin/b2b/dealers` |
| /api/v1 prefix | 30 | `/api/v1/approvals/escalations/check` |
| /api/analytics prefix | 9 | `/api/analytics/aggregations` |
| /api/b2b-pricing prefix | 7 | `/api/b2b-pricing/calculate` |

**Migration Mapping (Old → Canonical)**

Admin + Root duplicates:
- `/authors` → `/admin/authors`
- `/authors/` → `/admin/authors/`
- `/authors/by-id/{id}` → `/admin/authors/by-id/{id}`
- `/authors/{id}` → `/admin/authors/{id}`
- `/authors/{id}/toggle` → `/admin/authors/{id}/toggle`
- `/comments` → `/admin/comments`
- `/comments/` → `/admin/comments/`
- `/comments/bulk/delete` → `/admin/comments/bulk/delete`
- `/comments/bulk/status` → `/admin/comments/bulk/status`
- `/comments/{id}` → `/admin/comments/{id}`
- `/comments/{id}/like` → `/admin/comments/{id}/like`
- `/comments/{id}/reaction` → `/admin/comments/{id}/reaction`
- `/comments/{id}/reply` → `/admin/comments/{id}/reply`
- `/comments/{id}/seo` → `/admin/comments/{id}/seo`
- `/comments/{id}/status` → `/admin/comments/{id}/status`
- `/languages` → `/admin/languages`
- `/languages/active` → `/admin/languages/active`
- `/languages/default` → `/admin/languages/default`
- `/menus` → `/admin/menus`
- `/menus/` → `/admin/menus/`
- ... and 14 more

Public + Root duplicates:
- `/authors` → `/public/authors`
- `/authors/` → `/public/authors/`
- `/authors/by-id/{id}` → `/public/authors/by-id/{id}`
- `/authors/{id}` → `/public/authors/{id}`
- `/authors/{id}/toggle` → `/public/authors/{id}/toggle`
- `/comments` → `/public/comments`
- `/comments/` → `/public/comments/`
- `/comments/bulk/delete` → `/public/comments/bulk/delete`
- `/comments/bulk/status` → `/public/comments/bulk/status`
- `/comments/{id}` → `/public/comments/{id}`
- `/comments/{id}/like` → `/public/comments/{id}/like`
- `/comments/{id}/reaction` → `/public/comments/{id}/reaction`
- `/comments/{id}/reply` → `/public/comments/{id}/reply`
- `/comments/{id}/seo` → `/public/comments/{id}/seo`
- `/comments/{id}/status` → `/public/comments/{id}/status`
- `/languages` → `/public/languages`
- `/languages/active` → `/public/languages/active`
- `/languages/default` → `/public/languages/default`
- `/notifications` → `/public/notifications`
- `/notifications/create` → `/public/notifications/create`
- ... and 10 more

by-id variants:
- `/admin/authors/by-id/{id}` → `/admin/authors/{id}`
- `/admin/pages/by-id/{id}` → `/admin/pages/{id}`
- `/authors/by-id/{id}` → `/authors/{id}`
- `/public/authors/by-id/{id}` → `/public/authors/{id}`
- `/public/pages/by-id/{id}` → `/public/pages/{id}`

Trailing slash variants:
- `/admin/authors/` → `/admin/authors`
- `/admin/comments/` → `/admin/comments`
- `/admin/menus/` → `/admin/menus`
- `/authors/` → `/authors`
- `/comments/` → `/comments`
- `/graphql/` → `/graphql`
- `/languages/` → `/languages`
- `/menus/` → `/menus`
- `/public/authors/` → `/public/authors`
- `/public/comments/` → `/public/comments`
