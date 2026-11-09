# Telematics Wave B Sprint Backlog

Date: 2026-03-29
Owner: AI Assistant (Codex)
Status: Active
Depends on:
- [TELEMATICS_WAVE_B_MOBILE_AND_FIELD_OPERATIONS_EXECUTION_PLAN_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_B_MOBILE_AND_FIELD_OPERATIONS_EXECUTION_PLAN_2026-03-29.md)

## Sprint Cadence

- Sprint length: 2 weeks
- Definition of done: code + tests + lint + docs + audit trace
- Environments: local dev, staging pilot, controlled rollout

## Sprint-1 (Foundation And Contracts)

Target:
1. mobile auth/session contract hardening
2. manager/driver role shell API contracts
3. field workbench ticket model and lifecycle state machine

Primary deliverables:
1. tenant/site-safe mobile session bootstrap endpoint
2. role-scoped mobile navigation payload contract
3. install ticket create/assign/start/close API lifecycle
4. audit log entries for each field ticket state transition

## Sprint-2 (Driver And Notifications)

Target:
1. driver trip timeline and violation feed
2. notification routing for assignment/escalation/alert
3. acknowledge/recheck response loop for mobile actions

Primary deliverables:
1. driver trip list/detail contracts (completed 2026-03-29, B1-S2-01)
2. push queue retry and dead-letter controls (completed 2026-03-29, B1-S2-02)
3. notification action audit stream (completed 2026-03-29, B1-S2-03)
4. acknowledge/recheck response loop contracts for mobile notification actions (completed 2026-03-29, B1-S2-04)

## Sprint-3 (Field Evidence And Pilot Readiness)

Target:
1. replacement/RMA evidence and closure checklist
2. installer workbench pilot controls
3. UAT runbook for field + fleet operator

Primary deliverables:
1. evidence attachment metadata contract (completed 2026-03-29, B1-S3-01)
2. closure checklist with mandatory evidence gating (completed 2026-03-29, B1-S3-02)
3. pilot readiness scorecard and exit checklist (completed 2026-03-29, B1-S3-03)

## Sprint-4 (Stabilization And Controlled Pilot)

Target:
1. resilience and security hardening
2. telemetry/notification performance checks
3. controlled pilot and issue burndown

Primary deliverables:
1. rate-limit and tenant validation regression pack
2. performance baseline report (mobile + field routes)
3. pilot completion checkpoint and go/no-go note

## Sprint-2 Checkpoint (2026-03-29)

Completed under B1-S2-01:

1. Driver trip list contract implemented: `GET /admin/telematics/mobile/driver/trips`
2. Driver trip detail contract implemented: `GET /admin/telematics/mobile/driver/trips/:id`
3. Driver role gating + tenant/site scoped service context applied on both routes
4. Handler tests added for success and error matrix (bad ids, forbidden role, not found)
5. Validation run: `go test ./internal/modules/telematics/handler ./internal/modules/telematics/... ./cmd/api` (successful)

Completed under B1-S2-02:

1. Push notification retry lifecycle extended with `queued` and `dlq` statuses.
2. Retry metadata fields added (`retry_count`, `max_retries`, `next_retry_at`, `moved_to_dlq_at`).
3. Failure flow now schedules exponential retry and moves to DLQ after max retries.
4. Admin endpoints added for DLQ/retry controls:
	- `GET /admin/push-notifications/dlq`
	- `POST /admin/push-notifications/dlq/:id/retry`
	- `POST /admin/push-notifications/retries/process`
5. Validation run: `go test ./internal/modules/notifications/push/...` (successful)

Completed under B1-S2-03:

1. Push notification action audit stream eklendi (`audit_trail`).
2. Audit eventleri lifecycle boyunca kaydedilir hale getirildi:
	- `created`
	- `scheduled_retry`
	- `moved_dlq`
	- `retry_requested`
	- `sent`
3. Audit stream okuma endpointi eklendi:
	- `GET /admin/push-notifications/:id/audit`
4. Push service testleri audit stream davranisini dogrulayacak sekilde genisletildi.
5. Validation run: `go test ./internal/modules/notifications/push/...` (successful)

Sprint-2 hardening/regression pass (2026-03-29):

1. Push admin handler test coverage eklendi:
	- DLQ list filter/pagination
	- DLQ retry (invalid id + success)
	- Retry queue process counter response
	- Audit endpoint response contract
2. Validation run: `go test ./internal/modules/notifications/push/...` (successful)

Completed under B1-S2-04:

1. Mobile action response loop endpointleri eklendi:
	- `POST /notifications/:id/acknowledge`
	- `POST /notifications/:id/recheck`
2. Inbox read model response alanlari eklendi:
	- `acknowledged_at`
	- `rechecked_at`
	- `last_action`
	- `action_note`
3. Service katmaninda mobile action eventleri eklendi:
	- `notification.acknowledged`
	- `notification.rechecked`
4. Notification service test kapsaminda acknowledge/recheck loop testleri eklendi.
5. Validation run: `go test ./internal/modules/notifications/...` (successful)

## Immediate Next Action

Continue Sprint-3 execution complete — move to Sprint-4:

1. Sprint-3 items B1-S3-01, B1-S3-02, B1-S3-03 all completed 2026-03-29
2. Next: Sprint-4 (Stabilization And Controlled Pilot)

## Sprint-4 Checkpoint (2026-03-29)

Completed under B1-S4-01 + B1-S4-02 (rate-limit and tenant validation regression pack):

1. `sprint4_regression_test.go` created — 9 tests covering:
   - `TestRateLimitBlocksExcessRequests` — 3 req/min, 4th request → 429 + Retry-After
   - `TestRateLimitBypassedForHealthEndpoints` — /healthz never rate-limited
   - `TestRateLimitPerClientIsolation` — independent counters per remote IP
   - `TestRateLimitFieldTicketWriteEndpoint` — POST /mobile/field-tickets is rate-limited
   - `TestMobileEndpointRejectsMissingTenant` — 401 when tenant_id absent
   - `TestMobileEndpointRejectsMissingSite` — 403 when site_id absent
   - `TestFieldTicketCrossTenantIsolation` — Tenant B gets 404 on Tenant A's ticket
   - `TestFieldTicketCrossSiteIsolation` — site-2 gets 404 on site-1 ticket (same tenant)
   - `TestPermissionMiddlewareForbidsWithoutPermission` — 403 with empty permissions
2. Validation: all 9 pass (`go test ./internal/modules/telematics/handler/... -run "TestRateLimit|TestMobileEndpoint|TestFieldTicketCross|TestPermission"`)

Completed under B1-S4-03 (performance baseline report):

1. Handler: `AdminHandler.GetMobilePerfBaseline` added to `field_tickets.go`
   - Route: `GET /admin/telematics/mobile/perf-baseline`
   - Response: `routes[]` (11 mobile endpoints with p50/p95/p99 latency buckets), `sample_count`, `latency_profile`, `generated_at`
2. Tests: `TestPerfBaselineReportShape`, `TestPerfBaselineReportRouteCount` (≥10 routes)
3. Validation: both pass

Completed under B1-S4-04 (pilot go/no-go checkpoint):

1. Handler: `AdminHandler.GetPilotGoNoGo` added to `field_tickets.go`
   - Route: `GET /admin/telematics/mobile/pilot/go-no-go`
   - Aggregates `PilotReadinessSummary` + layer-status flags → single `verdict: "go"|"no_go"`
   - Exit criteria for "go": `total_tickets >= 5` AND `readiness_score >= 80`
   - `blocking_items[]` lists the specific exit criteria that are not yet met
2. Tests: `TestPilotGoNoGoEmptyState` (no_go), `TestPilotGoNoGoWithSufficientClosedTickets` (go with 5 fully-closed tickets)
3. Validation: both pass

Sprint-4 full regression run (2026-03-29):

- 30 handler tests PASS (21 from Sprint-2/3 + 9 S4 regression pack)
- `go test ./internal/modules/telematics/...` — all 11 packages OK, zero failures

## Sprint-3 Checkpoint (2026-03-29)

Completed under B1-S3-01:

1. `EvidenceAttachment` metadata struct added to field ticket model.
2. `Evidences []EvidenceAttachment` field added to `FieldTicket`.
3. `AppendEvidence()` method added to repository interface, memory, and mongo implementations.
4. Handler: `POST /admin/telematics/mobile/field-tickets/:id/evidences` endpoint added.
5. `evidence_attached` flag auto-set to `true` when evidence appended.

Completed under B1-S3-02:

1. `ClosureChecklistItem` struct added to field ticket model.
2. `DefaultClosureChecklist()` returns 3 required items: `device_power_check`, `gps_signal_check`, `photo_evidence`.
3. `Checklist []ClosureChecklistItem` field added to `FieldTicket`; seeded on ticket creation.
4. `UpdateChecklistItem()` method added to repository interface, memory, and mongo implementations.
5. `TransitionFieldTicket` (close action) now gates on: (a) closure note + evidence_attached, (b) all required checklist items completed.
6. Handlers added: `GET /admin/telematics/mobile/field-tickets/:id/checklist`, `PUT /admin/telematics/mobile/field-tickets/:id/checklist/:key`.

Completed under B1-S3-03:

1. `PilotReadinessSummary` struct added to model with: total/closed/open tickets, evidence/checklist counts, readiness score (0-100), exit criteria flag.
2. `ListTickets()` method added to repository interface, memory, and mongo.
3. `computePilotReadiness()` computes score: 70% from closure rate + 20% evidence + 10% checklist completion.
4. Exit criteria: total_tickets >= 5 AND score >= 80.
5. Handler: `GET /admin/telematics/mobile/field-tickets/pilot-readiness` endpoint registered.
6. Tests added for all Sprint-3 scenarios: evidence attach, checklist get/update, close gating, pilot readiness (empty + with closed ticket).
7. Validation run: `go test ./internal/modules/telematics/handler/...` — all 21 tests PASS.
