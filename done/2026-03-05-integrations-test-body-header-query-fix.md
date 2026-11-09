# Integrations Test Pipeline Fix (Body Template + Header/Query)

## Scope
- Module: Integrations (backend executor/service + Next admin test panel)
- Date: 2026-03-05

## Problem
- Manual test execution did not reliably use saved `body_template` when request payload was sent.
- `header_template` / `query_template` were not rendered/applied consistently.
- Test UI could create transient `test_values` keys while typing new field-mapping keys.

## Implemented
1. Backend request builder:
- Added effective request mapping composition in `BuildHTTPRequestWithConfig`.
- Fallback to stored `request_mapping.body_template` when test payload does not pass `body_template`.
- Applied `query_template` to URL with template rendering.
- Applied `header_template` with template rendering.
- Ensured payload warnings are deduplicated.

2. Backend runtime parity:
- `BuildHTTPRequest` now also applies rendered `header_template` + `query_template`.

3. Direct test response observability:
- `DirectExecuteAction` now returns `request_preview`:
  - method, url, headers, rendered body, payload_warnings
- Dry-run and real test responses include this preview (or `request_preview_error`).

4. Next Admin test panel:
- Parses and displays `request_preview` (URL/body/headers/warnings).
- Supports faster diagnosis of body-template injection issues.

5. Field Mapping -> Test Input sync:
- Auto-add now only for stable mapping keys (`key` and mapped `value` both non-empty).
- Prevents per-keystroke transient keys (`c`, `cc`, `ccc`) from polluting test input values.

## Validation
- Go tests/build passed for integrations package:
  - `go test ./internal/modules/integrations/...` ✅
- Frontend global type-check has unrelated existing errors in repository; targeted integration changes compiled logically and are isolated.

## Files
- go/aystek-core/internal/modules/integrations/executor/request_builder.go
- go/aystek-core/internal/modules/integrations/service/service.go
- node/aystek-core-next/lib/api/integrations.ts
- node/aystek-core-next/components/integrations/ActionTestPanel.tsx
- node/aystek-core-next/app/(dashboard)/admin/integrations/actions/ActionForm.tsx

## Additional Update (same day)
- `ValidateMapping` now returns `rendered_body_preview` for text templates.
- Action test panel now displays rendered body preview under payload preview.
- Test URL behavior updated:
  - If manual `test_url` already includes a non-root path (e.g., `https://host/api/gpsUpdate`), action path is no longer appended.
  - If `test_url` is base-only (e.g., `https://host`), action path continues to append as before.

## SSRF Policy Upgrade (Mongo + Next Managed)
- Replaced env-based private-IP exception approach with connector-scoped policy managed in MongoDB.
- New connector-managed fields:
  - `outbound_allowlist[]` (domain/ip/cidr/port/scheme)
  - `allow_private_ips_for_allowlist_match` (boolean)
- Effective allowlist is built at runtime per request:
  - manual connector `outbound_allowlist`
  - auto-injected host rules from connector `base_url` and `test_url`
- Private IP access policy:
  - still blocked by default
  - allowed only when connector flag is enabled AND destination matches explicit allowlist
- Implemented request-scoped SSRF policy execution path (`DoWithSSRFPolicy`) so connector updates reflect immediately without service restart.

## Allowlist Hardening (UI + Backend Validation)
- Connector allowlist now has stricter validation on both client and server:
  - Domain/IP required
  - Domain pattern validation
  - IP/CIDR validation
  - Port range validation (0..65535)
  - Scheme validation (`http|https`)
  - Duplicate rule detection
- Next connector form shows per-row warnings, red invalid field highlighting, and blocks save when validation errors exist.
- Added/seeded translation keys for all new validation messages.
- Menu seed note: no new module/page route introduced in this change set, so no new `menu_seed.go` entry was required.

## Response Schema Upgrade (Enterprise)
- `response_mapping` is now persisted end-to-end for actions (DTO + handler + repository + Next payload).
- HTTP response parser now normalizes by format:
  - JSON: parsed object/array
  - XML/SOAP: parsed into `xml.flat` + `xml.values` and promoted leaf values (e.g. `gpsUpdateResult`)
  - Text: stored as `raw`
- Added `_response_meta` for all responses (`format`, `content_type`, `size_bytes`).
- Added active success evaluation from `response_mapping.success_criteria`:
  - `status_code_min/max`
  - `predicate_expr` (JMESPath) over normalized response body
- If HTTP is `200` but predicate is false (example: `gpsUpdateResult=false`), execution is marked failed with `SUCCESS_CRITERIA_PREDICATE_NOT_MET` (non-retryable).
- Added tests for XML/SOAP predicate success/failure paths.
- Updated GPS integration seed and live records to include:
  - `response_mapping.success_criteria.predicate_expr = "gpsUpdateResult == 'true'"`

## Dynamic Response Capture (REST + SOAP + XML + Text)
- Unified response normalization now supports all common formats:
  - `json` responses: parsed + leaf flatten (`json.flat`) + promoted leaf values
  - `xml/soap` responses: parsed + path flatten (`xml.flat`) + leaf values (`xml.values`)
  - `text` responses: raw capture
- Added `_http` context (`status_code`, normalized headers, content_type) into predicate evaluation scope.
- Added automatic SOAP fault detection (`faultstring`) even on HTTP 200 responses.
- `success_criteria.predicate_expr` can now be written against dynamic leaf names without hard-coding a single service field.

## Response Schema Activation (UI Marked Block as Primary Source)
- Action success/failure evaluation now first checks `response_schema` expectations (the "Yanıt Şeması" block in Action form).
- `response_schema` can be used as dynamic expected map (e.g. `gpsUpdateResult: "true"`, `meta.status: "ok"`).
- Existing `response_mapping.success_criteria` remains backward-compatible fallback, but primary active source is `response_schema`.
