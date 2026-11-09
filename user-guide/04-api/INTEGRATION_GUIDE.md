---
Title: Integration Guide
Version: 2.3.0
Status: Enterprise Certified
Last Updated: 2026-03-09
Related Documents:
  - ./API_REFERENCE.md
  - ./ERROR_CODE_REFERENCE.md
  - ./API_CLIENT_ONBOARDING_BEST_PRACTICES.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
  - ../06-modules/modules/integrations.md
---

# Table of Contents
1. [Integration Flow](#1-integration-flow)
2. [Action Test Flow](#2-action-test-flow)
3. [Dynamic Response Capture](#3-dynamic-response-capture)
4. [Required Headers](#4-required-headers)
5. [Examples](#5-examples)
6. [Scheduled Pull API Flow](#6-scheduled-pull-api-flow)
7. [Error Handling](#7-error-handling)

## 1. Integration Flow
1. Resolve tenant context.
2. Resolve site context when route requires it.
3. Authenticate (for protected surfaces).
4. Build request payload using `request_mapping` + enrichment + defaults.
5. Execute action through connector policy (SSRF allowlist + auth + timeout).
6. Normalize response (`json` / `xml|soap` / `text`).
7. Evaluate success from **Response Schema** expectations (primary).
8. Persist execution record and request/response preview.

## 2. Action Test Flow
`POST /api/admin/integrations/actions/execute`

Purpose:
- Executes an action in test context without changing hook trigger semantics.
- Stores test execution trace with request/response preview.

Request payload (summary):
- `action_id`
- `input_payload`
- optional overrides: `test_url`, `headers`, `field_mapping`, `static_values`, `body_template`, `content_type`

Response payload (summary):
- `success`
- `http_status_code`
- `response_body`
- `request_preview` (`method`, `url`, `headers`, `body`, `payload_warnings`)

## 3. Dynamic Response Capture
Response handling is format-agnostic and active for REST/SOAP/XML/Text:

- JSON:
  - parsed object/array
  - `json.flat` leaf path map
  - promoted leaf values for easier predicates
- XML/SOAP:
  - `xml.flat` path map
  - `xml.values` leaf values
  - SOAP fault auto-detection (`faultstring`)
- Text:
  - raw payload capture

`_http` context is included for evaluation:
- `status_code`
- `content_type`
- normalized headers

### Primary Success Source: Response Schema
Action form "Yanıt Şeması" fields are now active runtime expectations.

Example:
```json
{
  "gpsUpdateResult": "true",
  "meta.status": "ok",
  "result.success": true
}
```

Behavior:
- All configured expectation keys are checked against normalized response.
- Mismatch => `RESPONSE_SCHEMA_MISMATCH` (non-retryable in test path).

Backward compatibility:
- `response_mapping.success_criteria` remains supported as fallback.

## 4. Required Headers
- `X-Tenant-ID` or `X-Tenant-Code`
- `X-Site-ID` for admin/graphql routes
- `Authorization` for protected routes
- `X-Request-ID` for trace continuity (recommended)

## 5. Examples
Public route:
```http
GET /public/settings HTTP/1.1
Host: api.example.com
X-Tenant-Code: acme
```

Admin route:
```http
GET /admin/menus HTTP/1.1
Host: api.example.com
Authorization: Bearer <token>
X-Tenant-ID: 65f0...
X-Site-ID: 65f1...
```

Integration action test route:
```http
POST /api/admin/integrations/actions/execute HTTP/1.1
Host: api.example.com
Authorization: Bearer <token>
X-Tenant-ID: 65f0...
X-Site-ID: 65f1...
Content-Type: application/json
```

## 6. Scheduled Pull API Flow
Use case: zamanlanmış dış veri çekimi yönetimi.

Base path:
- `/admin/integrations`

Endpoints:
1. `GET /pull-schedules`
2. `POST /pull-schedules`
3. `GET /pull-schedules/:id`
4. `PUT /pull-schedules/:id`
5. `DELETE /pull-schedules/:id`
6. `POST /pull-schedules/:id/test-run`
7. `GET /pull-runs`
8. `GET /pull-runs/:id`
9. `POST /pull-runs/:id/retry`

Operational flow:
1. Schedule oluştur/aktif et.
2. Test-run ile connector/action doğrula.
3. Pull run list/detail ile ingestion summary ve issues izle.
4. Başarısız run için manuel retry tetikle.

## 7. Error Handling
Always parse canonical error shape:
```json
{
  "error": "CrossTenantAccessAttempt",
  "message": "tenant mismatch",
  "trace_id": "req_456"
}
```
