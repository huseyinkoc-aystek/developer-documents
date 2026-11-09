**API Route Audit**

| Category | Count |
| --- | --- |
| Total unique paths | 768 |
| Admin + root duplicates | 34 |
| Public + root duplicates | 30 |
| /by-id vs /:id duplicates | 5 |
| Trailing slash variants | 10 |

**Admin + Root Duplicate Examples**

- `/admin/authors` ⇄ `/authors`
- `/admin/authors/` ⇄ `/authors/`
- `/admin/authors/by-id/{id}` ⇄ `/authors/by-id/{id}`
- `/admin/authors/{id}` ⇄ `/authors/{id}`
- `/admin/authors/{id}/toggle` ⇄ `/authors/{id}/toggle`
- `/admin/comments` ⇄ `/comments`
- `/admin/comments/` ⇄ `/comments/`
- `/admin/comments/bulk/delete` ⇄ `/comments/bulk/delete`
- `/admin/comments/bulk/status` ⇄ `/comments/bulk/status`
- `/admin/comments/{id}` ⇄ `/comments/{id}`
- `/admin/comments/{id}/like` ⇄ `/comments/{id}/like`
- `/admin/comments/{id}/reaction` ⇄ `/comments/{id}/reaction`
- `/admin/comments/{id}/reply` ⇄ `/comments/{id}/reply`
- `/admin/comments/{id}/seo` ⇄ `/comments/{id}/seo`
- `/admin/comments/{id}/status` ⇄ `/comments/{id}/status`
- ... and 19 more

**Public + Root Duplicate Examples**

- `/public/authors` ⇄ `/authors`
- `/public/authors/` ⇄ `/authors/`
- `/public/authors/by-id/{id}` ⇄ `/authors/by-id/{id}`
- `/public/authors/{id}` ⇄ `/authors/{id}`
- `/public/authors/{id}/toggle` ⇄ `/authors/{id}/toggle`
- `/public/comments` ⇄ `/comments`
- `/public/comments/` ⇄ `/comments/`
- `/public/comments/bulk/delete` ⇄ `/comments/bulk/delete`
- `/public/comments/bulk/status` ⇄ `/comments/bulk/status`
- `/public/comments/{id}` ⇄ `/comments/{id}`
- `/public/comments/{id}/like` ⇄ `/comments/{id}/like`
- `/public/comments/{id}/reaction` ⇄ `/comments/{id}/reaction`
- `/public/comments/{id}/reply` ⇄ `/comments/{id}/reply`
- `/public/comments/{id}/seo` ⇄ `/comments/{id}/seo`
- `/public/comments/{id}/status` ⇄ `/comments/{id}/status`
- ... and 15 more

**by-id vs id Duplicate Examples**

- `/admin/authors/by-id/{id}` ⇄ `/admin/authors/{id}`
- `/admin/pages/by-id/{id}` ⇄ `/admin/pages/{id}`
- `/authors/by-id/{id}` ⇄ `/authors/{id}`
- `/public/authors/by-id/{id}` ⇄ `/public/authors/{id}`
- `/public/pages/by-id/{id}` ⇄ `/public/pages/{id}`

**Trailing Slash Duplicate Examples**

- `/admin/authors` ⇄ `/admin/authors/`
- `/admin/comments` ⇄ `/admin/comments/`
- `/admin/menus` ⇄ `/admin/menus/`
- `/authors` ⇄ `/authors/`
- `/comments` ⇄ `/comments/`
- `/graphql` ⇄ `/graphql/`
- `/languages` ⇄ `/languages/`
- `/menus` ⇄ `/menus/`
- `/public/authors` ⇄ `/public/authors/`
- `/public/comments` ⇄ `/public/comments/`

