---
Title: API Reference (Content Ingestion & Content Ops)
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ./INTEGRATION_GUIDE.md
  - ./ERROR_CODE_REFERENCE.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
---

# Table of Contents
1. [Source of Truth](#1-source-of-truth)
2. [Context and Security Contract](#2-context-and-security-contract)
3. [API Feature Set](#3-api-feature-set)
4. [Connector APIs](#4-connector-apis)
5. [Ingestion APIs](#5-ingestion-apis)
6. [Moderation APIs](#6-moderation-apis)
7. [Content Ops Settings APIs](#7-content-ops-settings-apis)
8. [Error Model](#8-error-model)
9. [Permission Matrix](#9-permission-matrix)

## 1. Source of Truth
Detailed route and schema definitions are generated from runtime OpenAPI:
- `go/aystek-core/api-docs/swagger.yaml`

This document focuses on enterprise ingestion and content ops endpoints.

## 2. Context and Security Contract
Required headers by route scope:
- `Authorization: Bearer <token>`
- `X-Tenant-ID` or `X-Tenant-Code`
- `X-Site-ID` (admin site-scoped routes)
- `X-Request-ID` recommended for trace continuity

Rule: tenant/site are always resolved server-side; request body cannot override scope.

## 3. API Feature Set
- Server-side tenant + site scoping for all protected operations.
- Deterministic pagination/list response model across ingestion endpoints.
- Moderation workflow endpoints with explicit action semantics.
- Site-scoped settings APIs for policy/flags/limits/ad/injection governance.
- Canonical error contract with context and permission failure classes.

## 4. Connector APIs
Compatibility note:
- Marketplace provider connectors also use `/admin/connectors/:id/*`.
- Content ingestion connectors are namespaced as `/admin/ingestion/connectors/*` in runtime to avoid route collisions.
- If an alias `/admin/connectors` is enabled in a deployment, it must map to the same service contract below.

### 3.1 List Connectors
- `GET /admin/ingestion/connectors` (runtime)
- `GET /admin/connectors` (alias/contract)

Example request:
```http
GET /admin/ingestion/connectors?page=1&limit=25&status=active&type=rss
Authorization: Bearer <token>
X-Tenant-ID: <tenant>
X-Site-ID: <site>
```

Example response:
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "...",
        "name": "AA RSS",
        "connector_type": "rss",
        "status": "active",
        "health": { "status": "healthy", "error_count": 0 }
      }
    ]
  },
  "pagination": { "page": 1, "limit": 25, "total": 1, "total_pages": 1 }
}
```

### 3.2 Create Connector
- `POST /admin/ingestion/connectors` (runtime)
- `POST /admin/connectors` (alias/contract)

Example request:
```json
{
  "name": "Finance API",
  "slug": "finance-api",
  "connector_type": "api",
  "status": "active",
  "schedule": { "cron": "*/5 * * * *" },
  "fetch_policy": { "rate_limit_rpm": 60, "retries": 3 }
}
```

### 3.3 Run Connector
- `POST /admin/ingestion/connectors/:id/run` (runtime)
- `POST /admin/connectors/:id/run` (alias/contract)

Example request:
```json
{ "dry_run": true }
```

### 3.4 Disable Connector
- `POST /admin/ingestion/connectors/:id/disable` (runtime)
- `POST /admin/connectors/:id/disable` (alias/contract)

Example response:
```json
{ "success": true, "data": { "disabled": true } }
```

## 5. Ingestion APIs
### 4.1 Jobs
`GET /admin/ingestion/jobs`

Query parameters:
- `page`, `limit`
- `status`
- `connector_id`
- `date_from`, `date_to`

### 4.2 Items
`GET /admin/ingestion/items`

Query parameters:
- `page`, `limit`
- `status`
- `connector_id`
- `duplicate_group`
- `date_from`, `date_to`

Example response item:
```json
{
  "id": "...",
  "status": "pending_review",
  "source_meta": {
    "external_id": "src-123",
    "canonical_url": "https://source/news/123"
  },
  "normalized_content": { "title": "...", "summary": "..." }
}
```

## 6. Moderation APIs
### 5.1 Queue
`GET /admin/moderation/queue`

### 5.2 Approve
`POST /admin/moderation/:itemId/approve`

### 5.3 Reject
`POST /admin/moderation/:itemId/reject`

Example request:
```json
{ "notes": "Source duplication and low editorial quality" }
```

### 5.4 Edit
`PUT /admin/moderation/:itemId/edit`

Example request:
```json
{
  "content": { "title": "Revised title", "body_html": "<p>Edited body</p>" },
  "seo": { "meta_title": "...", "meta_description": "..." },
  "tags": ["economy"],
  "categories": ["..."]
}
```

### 5.5 Schedule
`POST /admin/moderation/:itemId/schedule`

Example request:
```json
{ "publish_at": "2026-02-18T15:00:00Z", "timezone": "Europe/Istanbul" }
```

## 7. Content Ops Settings APIs
Site-scoped routes:

- `GET /admin/sites/:siteId/policy`
- `PUT /admin/sites/:siteId/policy`
- `GET /admin/sites/:siteId/feature-flags`
- `PUT /admin/sites/:siteId/feature-flags`
- `GET /admin/sites/:siteId/content-limits`
- `PUT /admin/sites/:siteId/content-limits`
- `GET /admin/sites/:siteId/ad-policy`
- `PUT /admin/sites/:siteId/ad-policy`
- `GET /admin/sites/:siteId/injection-codes`
- `POST /admin/sites/:siteId/injection-codes`
- `PUT /admin/sites/:siteId/injection-codes/:id`
- `DELETE /admin/sites/:siteId/injection-codes/:id`

Injection create example:
```json
{
  "scope": "site",
  "environment": "prod",
  "device": "all",
  "key": "ga-tag",
  "label": "Google Analytics",
  "code": "<script>/* redacted */</script>",
  "is_active": true
}
```

## 8. Error Model
Canonical deterministic payload:
```json
{
  "success": false,
  "message": "errors.tenant.invalid",
  "error": {
    "code": "TENANT_ID_INVALID",
    "message_key": "errors.tenant.invalid"
  }
}
```

Common classes:
- Context errors: tenant/site missing or invalid
- Authorization errors: forbidden/insufficient permission
- Validation errors: invalid payload, invalid transition, size limits

## 9. Permission Matrix
| Endpoint Group | Minimum Permission |
| --- | --- |
| Connectors list/get | `ingestion.read` or module-equivalent read |
| Connector create/update/disable | `connector.manage` |
| Connector run | `ingestion.run` |
| Ingestion jobs/items list | `ingestion.read` |
| Moderation approve/reject | `moderation.approve` |
| Moderation edit/schedule/publish | `moderation.publish` |
| Content ops settings read | `content_ops_settings.read` |
| Content ops settings update/injection changes | `content_ops_settings.manage` |
