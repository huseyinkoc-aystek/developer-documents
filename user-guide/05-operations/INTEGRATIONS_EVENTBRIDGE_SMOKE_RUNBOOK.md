---
Title: Integrations EventBridge Smoke Runbook
Version: 1.1.0
Status: Production
Last Updated: 2026-03-07
Scope: All Modules (Module-Agnostic)
Related Documents:
  - ./OPERATIONS_GUIDE.md
  - ./INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
  - ./GITLAB_PROTECTED_ENV_SETUP.md
  - ../06-modules/modules/integrations.md
---

# Integrations EventBridge Smoke Runbook

## 1. Purpose
This runbook validates the end-to-end integration event pipeline for **any module**:

`module domain event -> integration hook lookup -> RabbitMQ publish -> worker consume -> execution status + retry/DLQ`

This is **not stock-specific**. It is valid for all modules that publish domain events.

## 2. Preconditions
1. API is healthy (`/healthz`, `/readyz`).
2. RabbitMQ is reachable.
3. At least one connector, action, and enabled hook exist.
4. Hook mode is set to `async` for queue-path validation.
5. Tenant/site context is valid for the selected test scope.

## 3. Dynamic Module Selection (UI)
In Hook edit/create screen (`/admin/integrations/hooks/:id`), select module events dynamically via:
1. `Olay Anahtarı (Event Key)` field (select or type).
2. Event key must match a published domain event, e.g.:
   - `stock.vehicle_asset.created`
   - `orders.order.created`
   - `payments.intent.succeeded`
3. Save hook with `Faz=Sonrası (post)` and `Mod=Asenkron (async)`.

This screen enables module-independent binding by event key, so integration wiring is reusable across domains.

## 4. Smoke Steps (CI/CD Post-Deploy)
1. Trigger a real domain action from selected module (create/update/delete/etc.).
2. Verify event enters queue path:
   - queue depth changes on `integration_events`.
3. Verify worker consumption:
   - worker logs show `Processing message: event_key=<selected_key>`.
4. Verify execution record:
   - one new execution exists for hook/event.
   - request/response persisted.
5. Verify success path:
   - status becomes `success`.
6. Verify retry path:
   - force temporary failure (target unavailable).
   - execution enters queued retry state with `next_retry_at`.
7. Verify DLQ path:
   - keep failure until retry limit exceeded.
   - execution becomes `failed`.
   - DLQ record appears (if DLQ enabled on hook).

## 5. Acceptance Criteria
1. Selected event key resolves exactly one intended hook path.
2. Async event is consumed by worker without tenant/site context errors.
3. Success case records full trace payloads.
4. Failure case follows retry policy and DLQ rules.
5. No silent drop in logs/metrics.

## 6. Failure Diagnostics
1. `No execution created`
   - event key mismatch
   - hook disabled
   - tenant/site mismatch
2. `Queued but never processed`
   - worker not running
   - RabbitMQ connection issue
3. `Processed but stays queued`
   - async execution semantics regression
4. `Immediate hard fail`
   - connector/auth/header/template error
5. `DLQ never fills`
   - hook DLQ disabled or retry policy misconfigured

## 7. Suggested CI Gate (Post-Deploy)
Minimum automated checks:
1. API health endpoints pass.
2. RabbitMQ queue presence check (`integration_events`, delayed queue if enabled).
3. Synthetic event trigger for one module.
4. Poll execution endpoint until `success` or timeout.
5. Negative test to confirm retry/DLQ behavior.

Pipeline fails if any required check fails.

## 8. Automation Script (bash/curl)
Script path:

`user-guide/05-operations/scripts/integrations_eventbridge_smoke.sh`

Backend repo runtime path (pipeline için):

`go/aystek-core/tools/ops/integrations_eventbridge_smoke.sh`

Required env vars:
1. `API_BASE_URL`
2. `AUTH_TOKEN`

Optional env vars:
1. `EVENT_KEY` (default: `stock.vehicle_asset.created`)
2. `HOOK_ID` (if set, hook lookup by event key is skipped)
3. `PAYLOAD_JSON` (JSON string)
4. `TIMEOUT_SEC` (default: `90`)
5. `POLL_INTERVAL_SEC` (default: `3`)
6. `TRIGGER_TYPE` (default: `pipeline_smoke`)
7. `EXPECTED_FINAL_STATUS` (default: `success`)

Example:

```bash
API_BASE_URL="https://staging-api.example.com" \
AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" \
EVENT_KEY="stock.vehicle_asset.created" \
PAYLOAD_JSON='{"asset_id":"smoke-1","vehicle_id":"v1","operation":"created"}' \
./user-guide/05-operations/scripts/integrations_eventbridge_smoke.sh
```

## 9. Example CI Job Snippet
```yaml
integrations_smoke:
  stage: post_deploy
  image: curlimages/curl:8.8.0
  script:
    - apk add --no-cache jq bash
    - chmod +x user-guide/05-operations/scripts/integrations_eventbridge_smoke.sh
    - API_BASE_URL="$STAGING_API_BASE_URL" AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" EVENT_KEY="stock.vehicle_asset.created" ./user-guide/05-operations/scripts/integrations_eventbridge_smoke.sh
```

## 10. Direct GitLab Integration
Backend reposunda aktif job dosyası eklendi:

`go/aystek-core/.gitlab-ci.yml`

Enable etmek için minimum değişkenler:
1. `INTEGRATIONS_SMOKE_ENABLED=true`
2. `INTEGRATIONS_SMOKE_TOKEN=<admin bearer token>`
3. `INTEGRATIONS_SMOKE_API_BASE_URL=<staging api url>` (veya `STAGING_API_BASE_URL`)
4. Production için protected environment/variable doğrulaması:
   - `go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`
