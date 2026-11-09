# Telematics Wave B Sprint-1 Foundation And Contracts Execution Plan

Date: 2026-03-29
Owner: AI Assistant (Codex)
Status: In Progress
Depends on:
- [TELEMATICS_WAVE_B_SPRINT_BACKLOG_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_B_SPRINT_BACKLOG_2026-03-29.md)

## Sprint Objective

Establish secure and testable foundations for mobile manager/driver surfaces and field workbench lifecycle contracts.

## In Scope

1. mobile auth/session bootstrap contracts
2. role-scoped mobile shell payload contracts (manager/driver)
3. install ticket lifecycle contracts (create, assign, start, close)
4. field state transition audit requirements

## Sprint Tasks

### Task B1-S1-01: Mobile Session Bootstrap Contract

Deliverables:
1. endpoint contract for mobile session context (`tenant_id`, `site_id`, `role`, `permissions`)
2. invalid tenant/site response matrix (401/403/404)
3. session expiry and refresh payload fields

Acceptance:
1. invalid tenant/site never yields successful bootstrap
2. manager and driver roles resolve deterministic permission sets
3. contract tests cover success + error paths

### Task B1-S1-02: Role Shell Contract

Deliverables:
1. manager shell payload sections (`fleet_summary`, `alerts`, `workbench_shortcuts`)
2. driver shell payload sections (`trip_feed`, `score_summary`, `violation_actions`)
3. versioned response envelope (`version`, `generated_at`, `trace_id`)

Acceptance:
1. role mismatch returns forbidden contract
2. payload shape remains backward-compatible for one minor version
3. schema tests verify mandatory fields

### Task B1-S1-03: Field Ticket Lifecycle Contract

Deliverables:
1. ticket state machine: `created -> assigned -> in_progress -> closed`
2. transition guard rules and rejection codes
3. minimum metadata for create/close operations

Acceptance:
1. invalid transitions return deterministic 409 contract
2. close operation requires mandatory closure evidence flags
3. lifecycle history is queryable per ticket id

### Task B1-S1-04: Audit And Observability Contract

Deliverables:
1. audit event taxonomy for ticket and mobile auth actions
2. correlation id propagation requirement
3. basic dashboards for ticket throughput and transition failures

Acceptance:
1. each state transition emits exactly one audit event
2. failed transitions are searchable by `trace_id`
3. daily transition failure trend is exportable

## Validation Commands (Planned)

Backend:
1. `go test ./internal/modules/telematics/...`
2. `go test ./internal/modules/activity_logs/...`

Frontend/API client:
1. `npm run -s lint -- --file 'lib/api/telematics.ts'`
2. `npm run -s test`

## Execution Checkpoint (2026-03-29)

Task B1-S1-01 started and backend contract draft is implemented.

Implemented backend artifacts:

1. new endpoint handler:
	- `GET /admin/telematics/mobile/session/bootstrap`
	- file: `internal/modules/telematics/handler/mobile_session.go`
2. route registration:
	- file: `internal/modules/telematics/routes.go`
3. contract tests (success + 401/403/404 matrix):
	- file: `internal/modules/telematics/handler/mobile_session_test.go`

Contract payload includes:

1. `user_id`, `tenant_id`, `site_id`
2. `role`, `roles`, `permissions`
3. `issued_at`, `expires_at`
4. `error_matrix` (401/403/404)

Executed validation:

1. `go test ./internal/modules/telematics/handler ./internal/modules/telematics/...` -> passed
2. `go test ./internal/modules/stock/...` -> passed (regression guard)

Task B1-S1-04 completion checkpoint: daily transition failure trend export is implemented.

Implemented backend artifacts:

1. failure trend export endpoint:
	- `GET /admin/telematics/mobile/field-tickets/audit/failure-trend`
	- supports `days` (1..31), optional `ticket_id`, and `format=json|csv`
	- file: `internal/modules/telematics/handler/field_tickets.go`
2. route registration:
	- file: `internal/modules/telematics/routes.go`
3. repository contract/implementation extension for date-range audit queries:
	- `internal/modules/telematics/field_tickets/repository/interface.go`
	- `internal/modules/telematics/field_tickets/repository/mongo.go`
	- `internal/modules/telematics/field_tickets/repository/memory.go`
4. contract tests for json and csv trend export:
	- file: `internal/modules/telematics/handler/field_tickets_test.go`

Executed validation:

1. `go test ./internal/modules/telematics/handler ./internal/modules/telematics/field_tickets/... ./cmd/api` -> passed

Task B1-S1-05 started and field ticket persistence draft is migrated to repository-backed storage.

Implemented backend artifacts:

1. new field ticket domain model + repository contract:
	- `internal/modules/telematics/field_tickets/model/model.go`
	- `internal/modules/telematics/field_tickets/repository/interface.go`
2. repository implementations:
	- mongo persistence: `internal/modules/telematics/field_tickets/repository/mongo.go`
	- test fallback memory repository: `internal/modules/telematics/field_tickets/repository/memory.go`
3. handler migration from global in-memory maps to repository usage:
	- file: `internal/modules/telematics/handler/field_tickets.go`
4. DI wiring update:
	- `cmd/api/main.go` telematics admin now receives `telematicsFieldTicketsRepo`
	- `internal/modules/telematics/handler/admin.go` constructor signature updated
5. test isolation update:
	- `internal/modules/telematics/handler/field_tickets_test.go` now injects per-test memory repository

Executed validation:

1. `go test ./internal/modules/telematics/handler` -> passed
2. `go test ./internal/modules/telematics/field_tickets/...` -> passed
3. `go test ./cmd/api` -> passed

Task B1-S1-04 started and audit/trace contract draft is implemented.

Implemented backend artifacts:

1. trace_id enforcement on write routes:
	- create and transition flows require trace id
	- missing trace id returns validation error
2. audit taxonomy draft:
	- `field_ticket.created`
	- `field_ticket.transitioned`
	- `field_ticket.transition_rejected`
	- `field_ticket.validation_failed`
3. trace-searchable audit endpoint:
	- `GET /admin/telematics/mobile/field-tickets/audit?trace_id=...`
4. transition history now carries `trace_id` on state changes

Executed validation:

1. `go test ./internal/modules/telematics/handler ./internal/modules/telematics/...` -> passed
2. `go test ./internal/modules/stock/...` -> passed (regression guard)

Task B1-S1-02 started and backend contract draft is implemented.

Implemented backend artifacts:

1. new endpoint handler:
	- `GET /admin/telematics/mobile/shell`
	- file: `internal/modules/telematics/handler/mobile_session.go`
2. route registration:
	- file: `internal/modules/telematics/routes.go`
3. contract tests (manager + driver + role mismatch 403):
	- file: `internal/modules/telematics/handler/mobile_session_test.go`

Contract payload includes:

1. `envelope.version`
2. `envelope.generated_at`
3. `envelope.trace_id`
4. role-scoped shell blocks:
	- manager: `fleet_summary`, `alerts`, `workbench_shortcuts`
	- driver: `trip_feed`, `score_summary`, `violation_actions`

Task B1-S1-03 started and backend lifecycle draft is implemented.

Implemented backend artifacts:

1. field ticket lifecycle handler endpoints:
	- `POST /admin/telematics/mobile/field-tickets`
	- `GET /admin/telematics/mobile/field-tickets/:id`
	- `POST /admin/telematics/mobile/field-tickets/:id/transitions`
2. lifecycle state machine draft:
	- `created -> assigned -> in_progress -> closed`
	- invalid transition returns `409`
	- close transition requires `closure_note` and `evidence_attached=true`
3. route registration:
	- file: `internal/modules/telematics/routes.go`
4. contract tests:
	- file: `internal/modules/telematics/handler/field_tickets_test.go`

Executed validation:

1. `go test ./internal/modules/telematics/handler ./internal/modules/telematics/...` -> passed
2. `go test ./internal/modules/stock/...` -> passed (regression guard)

## Sprint Exit Criteria

1. all Sprint-1 contracts published in code and docs
2. tenant/site guardrails validated by automated tests
3. field ticket state machine rejects invalid transitions reliably
4. audit trace covers all lifecycle and auth transitions
