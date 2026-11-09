**Swagger Operations Summary**

| Metric | Value |
| --- | --- |
| Total unique paths | 768 |
| Total unique tags | 67 |

**Top 5 Tags By Operation Count**

| Tag | Operation Count |
| --- | --- |
| Forum | 89 |
| Products | 61 |
| Notifications | 56 |
| Payments | 44 |
| Billing | 40 |

**Suspicious Patterns**

- Admin + root duplicate: `/admin/authors/` and `/authors/`
- Admin + root duplicate: `/admin/authors/by-id/{id}` and `/authors/by-id/{id}`
- Admin + root duplicate: `/admin/authors/{id}/toggle` and `/authors/{id}/toggle`
- Admin + root duplicate: `/admin/authors/{id}` and `/authors/{id}`
- Admin + root duplicate: `/admin/authors` and `/authors`
- Admin + root duplicate: `/admin/comments/` and `/comments/`
- Admin + root duplicate: `/admin/comments/bulk/delete` and `/comments/bulk/delete`
- Admin + root duplicate: `/admin/comments/bulk/status` and `/comments/bulk/status`
- Admin + root duplicate: `/admin/comments/{id}/like` and `/comments/{id}/like`
- Admin + root duplicate: `/admin/comments/{id}/reaction` and `/comments/{id}/reaction`
- Admin + root duplicate: `/admin/comments/{id}/reply` and `/comments/{id}/reply`
- Admin + root duplicate: `/admin/comments/{id}/seo` and `/comments/{id}/seo`
- Admin + root duplicate: `/admin/comments/{id}/status` and `/comments/{id}/status`
- Admin + root duplicate: `/admin/comments/{id}` and `/comments/{id}`
- Admin + root duplicate: `/admin/comments` and `/comments`
- Admin + root duplicate: `/admin/languages/active` and `/languages/active`
- Admin + root duplicate: `/admin/languages/default` and `/languages/default`
- Admin + root duplicate: `/admin/languages` and `/languages`
- Admin + root duplicate: `/admin/menus/` and `/menus/`
- Admin + root duplicate: `/admin/menus/{id}` and `/menus/{id}`
- Admin + root duplicate: `/admin/menus` and `/menus`
- Admin + root duplicate: `/admin/notifications/create` and `/notifications/create`
- Admin + root duplicate: `/admin/notifications/list` and `/notifications/list`
- Admin + root duplicate: `/admin/notifications/me` and `/notifications/me`
- Admin + root duplicate: `/admin/notifications/preferences` and `/notifications/preferences`
- Admin + root duplicate: `/admin/notifications/push/send` and `/notifications/push/send`
- Admin + root duplicate: `/admin/notifications/push/test` and `/notifications/push/test`
- Admin + root duplicate: `/admin/notifications/push` and `/notifications/push`
- Admin + root duplicate: `/admin/notifications/read-all` and `/notifications/read-all`
- Admin + root duplicate: `/admin/notifications/receipts` and `/notifications/receipts`
- Admin + root duplicate: `/admin/notifications/{id}/read` and `/notifications/{id}/read`
- Admin + root duplicate: `/admin/notifications/{id}` and `/notifications/{id}`
- Admin + root duplicate: `/admin/notifications` and `/notifications`
- Admin + root duplicate: `/admin/subscriptions/{id}` and `/subscriptions/{id}`
- Duplicate by-id vs id variants: `/admin/authors/by-id/{id}` and `/admin/authors/{id}`
- Duplicate by-id vs id variants: `/admin/pages/by-id/{id}` and `/admin/pages/{id}`
- Duplicate by-id vs id variants: `/authors/by-id/{id}` and `/authors/{id}`
- Duplicate by-id vs id variants: `/public/authors/by-id/{id}` and `/public/authors/{id}`
- Duplicate by-id vs id variants: `/public/pages/by-id/{id}` and `/public/pages/{id}`
- Duplicate trailing slash variants: `/admin/authors` and `/admin/authors/`
- Duplicate trailing slash variants: `/admin/comments` and `/admin/comments/`
- Duplicate trailing slash variants: `/admin/menus` and `/admin/menus/`
- Duplicate trailing slash variants: `/authors` and `/authors/`
- Duplicate trailing slash variants: `/comments` and `/comments/`
- Duplicate trailing slash variants: `/graphql` and `/graphql/`
- Duplicate trailing slash variants: `/languages` and `/languages/`
- Duplicate trailing slash variants: `/menus` and `/menus/`
- Duplicate trailing slash variants: `/public/authors` and `/public/authors/`
- Duplicate trailing slash variants: `/public/comments` and `/public/comments/`
- Public + root duplicate: `/public/authors/` and `/authors/`
- Public + root duplicate: `/public/authors/by-id/{id}` and `/authors/by-id/{id}`
- Public + root duplicate: `/public/authors/{id}/toggle` and `/authors/{id}/toggle`
- Public + root duplicate: `/public/authors/{id}` and `/authors/{id}`
- Public + root duplicate: `/public/authors` and `/authors`
- Public + root duplicate: `/public/comments/` and `/comments/`
- Public + root duplicate: `/public/comments/bulk/delete` and `/comments/bulk/delete`
- Public + root duplicate: `/public/comments/bulk/status` and `/comments/bulk/status`
- Public + root duplicate: `/public/comments/{id}/like` and `/comments/{id}/like`
- Public + root duplicate: `/public/comments/{id}/reaction` and `/comments/{id}/reaction`
- Public + root duplicate: `/public/comments/{id}/reply` and `/comments/{id}/reply`
- Public + root duplicate: `/public/comments/{id}/seo` and `/comments/{id}/seo`
- Public + root duplicate: `/public/comments/{id}/status` and `/comments/{id}/status`
- Public + root duplicate: `/public/comments/{id}` and `/comments/{id}`
- Public + root duplicate: `/public/comments` and `/comments`
- Public + root duplicate: `/public/languages/active` and `/languages/active`
- Public + root duplicate: `/public/languages/default` and `/languages/default`
- Public + root duplicate: `/public/languages` and `/languages`
- Public + root duplicate: `/public/notifications/create` and `/notifications/create`
- Public + root duplicate: `/public/notifications/list` and `/notifications/list`
- Public + root duplicate: `/public/notifications/me` and `/notifications/me`
- Public + root duplicate: `/public/notifications/preferences` and `/notifications/preferences`
- Public + root duplicate: `/public/notifications/push/send` and `/notifications/push/send`
- Public + root duplicate: `/public/notifications/push/test` and `/notifications/push/test`
- Public + root duplicate: `/public/notifications/push` and `/notifications/push`
- Public + root duplicate: `/public/notifications/read-all` and `/notifications/read-all`
- Public + root duplicate: `/public/notifications/receipts` and `/notifications/receipts`
- Public + root duplicate: `/public/notifications/{id}/read` and `/notifications/{id}/read`
- Public + root duplicate: `/public/notifications/{id}` and `/notifications/{id}`
- Public + root duplicate: `/public/notifications` and `/notifications`
