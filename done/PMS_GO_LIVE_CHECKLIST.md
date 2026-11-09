# PMS Go Live Checklist

## Scope
This checklist covers PMS production readiness for Sprint-15 hardening.

## Pre-Deployment
- [ ] PMS changes are isolated in a dedicated PR.
- [ ] Run the scope guard: `bash tools/pms_scope_guard.sh`.
- [ ] Confirm no secrets, target URLs, or payloads are logged in PMS paths.
- [ ] Review release notes and confirm no API breaking changes.

## Data Verification
- [ ] Call GET /admin/pms/webhooks/legacy-summary and record legacy_count.
- [ ] If legacy_count > 0, create an approved migration plan per tenant/site.
- [ ] Update legacy subscriptions with explicit event_types.
- [ ] Re-run legacy-summary until legacy_count reaches 0.

## Observability
- [ ] Verify pms_webhooks_legacy_allowlist_total gauge is present and alertable.
- [ ] Verify retry metrics exist for exports and webhook deliveries.
- [ ] Confirm audit logs record retry actions without sensitive data.

## E2E Simulation
- [ ] Run `go test ./internal/modules/pms/... -count=1` and confirm PASS.
- [ ] Validate export lifecycle E2E coverage (create -> process -> cooldown -> retry).
- [ ] Validate webhook lifecycle E2E coverage (dispatch -> fail -> cooldown -> retry -> sent).
- [ ] Confirm audit log entries are created for webhook delivery retry actions.

## Deployment
- [ ] Run `go build ./cmd/api`.
- [ ] Run PMS seed/permission checks (menu, roles, translations).
- [ ] Validate retention TTL indexes exist for PMS collections.
- [ ] Confirm environment variables in .env.example are populated in production.

## Rollback Plan
- [ ] Keep previous build artifacts available.
- [ ] Document rollback steps for API and database migrations.
- [ ] Capture a snapshot of PMS metrics and key dashboards before release.
