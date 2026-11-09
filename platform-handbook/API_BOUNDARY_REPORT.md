**Admin/Public Boundary Report**

This report lists modules with public routes and highlights root duplicates that blur the admin/public boundary.

**Modules Exposing Public Routes**

| Module (Tag) | Public Operations | Example Public Path |
| --- | --- | --- |
| Ads | 5 | `/public/ads` |
| Authors | 11 | `/public/authors` |
| Billing | 26 | `/public/billing/accounts` |
| Catalog | 13 | `/public/catalog` |
| Chatbot | 3 | `/public/chatbot/context` |
| Comments | 13 | `/public/comments` |
| Contacts | 5 | `/public/contacts` |
| Crm | 3 | `/public/webhooks/hubspot` |
| Forum | 25 | `/public/forums/ads` |
| Languages | 12 | `/public/languages` |
| Licenses | 1 | `/public/licenses/access` |
| Localized Content | 6 | `/public/localized-content` |
| Marketplaces | 1 | `/public/webhooks/marketplaces/{provider}` |
| Media | 6 | `/public/media` |
| Notifications | 16 | `/public/notifications` |
| Pages | 13 | `/public/pages` |
| Payments | 13 | `/public/payment-intents/callbacks/{provider}` |
| Posts | 12 | `/public/posts` |
| References | 12 | `/public/references` |
| Settings | 2 | `/public/settings` |
| Shipping | 1 | `/public/shipping/track/{tracking_number}` |
| Sliders | 7 | `/public/sliders` |
| Surveys | 7 | `/public/surveys` |
| Taxonomy | 19 | `/public/categories` |
| Translations | 1 | `/public/translations` |
| Units | 5 | `/public/units` |

**Modules With Root (Non-Prefixed) Routes**

| Module (Tag) | Root Operations | Example Root Path |
| --- | --- | --- |
| Ai | 9 | `/api/analytics/aggregations` |
| Approvals | 16 | `/api/v1/approvals/escalations/check` |
| Authors | 10 | `/authors` |
| B2B | 36 | `/api/admin/b2b/dealers` |
| Carts | 12 | `/cart` |
| Checkout | 6 | `/checkout/{id}` |
| Comments | 12 | `/comments` |
| Deployment | 4 | `/deployment/deploy` |
| Events | 1 | `/events/stream` |
| Fx | 10 | `/api/admin/fx-health` |
| Golive | 9 | `/go-live/beta-customer` |
| Graphql | 3 | `/graphql` |
| Languages | 4 | `/languages` |
| Maintenance | 1 | `/maintenance` |
| Menus | 5 | `/menus` |
| Notifications | 17 | `/notifications` |
| Ops | 3 | `/healthz` |
| Performance | 6 | `/performance/cache/clear` |
| Pricing | 17 | `/api/v1/pricing/analytics/{id}` |
| Products | 22 | `/api/v1/products/{productID}/reviews` |
| Reports | 18 | `/executions/{id}` |
| Security | 14 | `/security/account/unlock` |
| Stock | 9 | `/api/checkouts/{checkout_id}/reservation` |
| Uat | 7 | `/uat/bug-report` |

**Admin/Public Root Duplicates**

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
- `/admin/languages` ⇄ `/languages`
- `/admin/languages/active` ⇄ `/languages/active`
- `/admin/languages/default` ⇄ `/languages/default`
- `/admin/menus` ⇄ `/menus`
- `/admin/menus/` ⇄ `/menus/`
- `/admin/menus/{id}` ⇄ `/menus/{id}`
- `/admin/notifications` ⇄ `/notifications`
- `/admin/notifications/create` ⇄ `/notifications/create`
- `/admin/notifications/list` ⇄ `/notifications/list`
- `/admin/notifications/me` ⇄ `/notifications/me`
- `/admin/notifications/preferences` ⇄ `/notifications/preferences`
- `/admin/notifications/push` ⇄ `/notifications/push`
- `/admin/notifications/push/send` ⇄ `/notifications/push/send`
- `/admin/notifications/push/test` ⇄ `/notifications/push/test`
- `/admin/notifications/read-all` ⇄ `/notifications/read-all`
- `/admin/notifications/receipts` ⇄ `/notifications/receipts`
- `/admin/notifications/{id}` ⇄ `/notifications/{id}`
- `/admin/notifications/{id}/read` ⇄ `/notifications/{id}/read`
- `/admin/subscriptions/{id}` ⇄ `/subscriptions/{id}`
- `/public/authors` ⇄ `/authors`
- `/public/authors/` ⇄ `/authors/`
- `/public/authors/by-id/{id}` ⇄ `/authors/by-id/{id}`
- `/public/authors/{id}` ⇄ `/authors/{id}`
- `/public/authors/{id}/toggle` ⇄ `/authors/{id}/toggle`
- `/public/comments` ⇄ `/comments`
- ... and 24 more
