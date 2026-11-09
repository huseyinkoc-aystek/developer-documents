# Telematics Wave A Phase-10 Multi-Vendor Effectiveness Parity And Decision Automation Execution Plan

Date: 2026-03-29
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- [TELEMATICS_WAVE_A_PHASE_9_COMPLETION_REPORT_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_9_COMPLETION_REPORT_2026-03-29.md)
- [TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md)
- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)

## Objective

Phase-10 hedefi, Phase-9 ile tamamlanan freeze closure ve recovery effectiveness governance zincirini Teltonika-odakli tek-vendor dogasindan cikarip cok-vendor policy parity ve karar otomasyonu seviyesine tasimaktir.

Bu faz su bosluklari kapatir:
- effectiveness scoring var ama vendor bazli threshold ve yorumlama farklari sinirli
- closure quality KPI var ama vendor-segment policy parity ve karsilastirmali trend zayif
- escalation zinciri var ama operatora otomatik sonraki en iyi aksiyon onerisi sinirli
- governance export var ama vendor-comparative closure quality snapshot standardi eksik

## Product Outcome

Teslim edilecek ana yuzey:
- `Multi-Vendor Freeze Closure And Recovery Decision Desk`

Bu yuzey sunlari baglar:
- vendor-specific effectiveness threshold matrix
- policy parity and drift detection
- action recommendation engine (next-best-action)
- comparative closure and recurrence trend
- governance export and scheduled parity snapshot

## In Scope

1. vendor-aware effectiveness threshold matrix
2. policy drift detection for closure and recurrence
3. next-best-action recommendation model
4. vendor-comparative closure/effectiveness KPI and trend
5. decision confidence visibility and override reason capture
6. governance export and scheduled parity snapshot
7. Teltonika + Queclink parity-ready model skeleton

## Out Of Scope

1. full autonomous field execution
2. technician dispatch optimization
3. billing/commercial optimization linkage
4. mobile field acceptance workflow
5. full multi-vendor feature parity beyond governance layer

## Delivery Tracks

### Track P10-A: Vendor Parity Model

1. vendor threshold profile entity
2. closure/effectiveness policy matrix
3. vendor parity drift evaluator
4. parity status projection

### Track P10-B: Decision Automation

1. next-best-action rule set
2. confidence scoring and explainability payload
3. operator override reason model
4. recommendation outcome feedback loop

### Track P10-C: Comparative Governance

1. vendor comparative KPI blocks
2. recurrence and closure trend by vendor
3. parity drift notices and recurrence policy
4. export/report parity visibility and scheduled snapshot

## Core Entities

1. `device_vendor_effectiveness_profile`
2. `device_vendor_policy_parity_state`
3. `device_recovery_recommendation`
4. `device_recommendation_outcome`
5. `device_vendor_governance_snapshot`

## Proposed States

1. `parity_unknown`
2. `parity_aligned`
3. `parity_drift_detected`
4. `recommendation_available`
5. `recommendation_applied`
6. `recommendation_overridden`
7. `recommendation_ineffective`
8. `vendor_escalation_open`

## Backend Work Plan

1. add vendor parity read/write model under `telematics/devices`
2. persist vendor threshold and policy profile history
3. add parity drift evaluator and notice chain
4. add recommendation engine and confidence score payload
5. add recommendation apply/override write-model and feedback capture
6. extend governance export/report pack with vendor parity snapshots

## Frontend Work Plan

1. expand device detail with recommendation and confidence card
2. add vendor parity block into `admin/telematics/dashboard`
3. add comparative vendor filters and parity drift widgets
4. add recommendation apply/override actions and audit deep-links
5. add stock-side parity snapshot visibility for exports/reports

## Governance And Audit

1. vendor parity drift audit stream
2. recommendation decision and override audit
3. comparative closure/effectiveness recurrence notices
4. export parity pressure visibility
5. scheduled governance parity snapshots

## Acceptance Criteria

1. operator can compare closure/effectiveness quality between vendors
2. recommendation engine provides explainable next-best-action per device
3. recommendation apply/override path is persistent and auditable
4. parity drift generates notice/escalation with recurrence controls
5. export/report chain exposes vendor parity closure quality

## Immediate Implementation Order

1. vendor profile and parity summary model
2. comparative dashboard KPI and vendor filters
3. recommendation write-model and action endpoints
4. parity notice and recurrence policy
5. export/report parity snapshot visibility

## Risks

1. recommendation quality can degrade with sparse vendor data
2. parity thresholds can overfit early vendor behavior
3. comparative dashboards can create noisy interpretation if not normalized

## Mitigation

1. keep thresholds explicit and versioned by vendor profile
2. require confidence and evidence payload for each recommendation
3. publish normalized KPI definitions with alert cooldown discipline

## Exit Criteria

Phase-10 closes when:
- vendor-level closure/effectiveness parity is measurable
- next-best-action recommendations are explainable and auditable
- parity drift notices are recurrently governed
- export/report chain carries scheduled vendor parity snapshots

## References

- [TELEMATICS_WAVE_A_PHASE_9_COMPLETION_REPORT_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_9_COMPLETION_REPORT_2026-03-29.md)
- [TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md)
- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)

## Execution Status (2026-03-29)

Status: implemented + validated

Completed implementation items:
1. Phase-10 vendor parity and decision automation fields/events integrated in telematics devices backend and dashboard.
2. Phase-10 stock reporting chain connected with new governance endpoints:
	- `POST /admin/stock/reports/governance/phase10/acknowledge`
	- `POST /admin/stock/reports/governance/phase10/recheck`
3. Phase-10 stock reporting snapshot integrated:
	- new report type `phase10_governance`
	- dashboard KPI snapshot totals/status/closure metrics
	- export row builder for phase10 governance stream
4. Frontend stock ecosystem API/types and telematics dashboard wired for Phase-10 export notice ack/recheck and snapshot cards.

Validation evidence:
1. `cd /workspace/go/aystek-core && go test ./internal/modules/stock/reporting/...`
2. `cd /workspace/node/aystek-core-next && npx eslint 'app/(dashboard)/admin/telematics/dashboard/page.tsx' 'lib/api/stock-ecosystem.ts' 'lib/types/stock-ecosystem.ts' 'app/(dashboard)/admin/stock/reports/exports/jobs/page.tsx'`