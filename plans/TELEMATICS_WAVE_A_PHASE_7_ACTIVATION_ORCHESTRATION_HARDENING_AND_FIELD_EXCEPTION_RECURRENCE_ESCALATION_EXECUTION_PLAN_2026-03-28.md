# Telematics Wave A Phase-7 Activation Orchestration Hardening And Field Exception Recurrence Escalation Execution Plan

Date: 2026-03-28
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- [TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md)
- [TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md)
- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)

## Objective

Phase-7 hedefi, Phase-6 ile açılmış activation automation ve field exception governance hattını daha kurallı orchestration, recurrence escalation ve freeze-like protection seviyesine taşımaktır.

Bu faz şu boşlukları kapatır:
- activation branch’ler var ama branch-to-branch orchestration kararları hâlâ sınırlı
- retry, hold ve manual handoff var ama escalation ladder henüz tam policy-driven değil
- field exception recurrence ölçülüyor ama second-order escalation ve fleet-wide priority desk eksik
- export-side governance zinciri var ama activation orchestration kararlarıyla birleşik closure board eksik

## Product Outcome

Teslim edilecek ana yüzey:
- `Activation Orchestration Control And Exception Escalation Desk`

Bu yüzey şunları bağlar:
- activation branch orchestration policies
- retry budget and branch cooldown ladder
- recurrence escalation and freeze recommendation
- fleet-level exception priority desk
- orchestration governance, export and closure trend

## In Scope

1. activation orchestration policy matrix
2. retry budget, cooldown and branch fallback model
3. recurrence escalation ladder for field exceptions
4. freeze-like recommendation and controlled recovery entrypoints
5. fleet-level exception priority and batch remediation macros
6. orchestration governance export and closure quality visibility
7. Teltonika-first activation orchestration hardening

## Out Of Scope

1. technician routing optimization
2. carrier contract automation
3. multi-vendor full parity rollout
4. full FOTA orchestration
5. mobile field app implementation

## Delivery Tracks

### Track P7-A: Activation Orchestration

1. branch policy matrix
2. retry budget and cooldown ladder
3. fallback and recovery transitions
4. orchestration closure quality

### Track P7-B: Exception Escalation

1. recurrence escalation model
2. freeze recommendation and protection state
3. fleet-wide priority ranking
4. controlled recovery entrypoints

### Track P7-C: Governance

1. orchestration KPI and notice chain
2. escalation closure trend
3. export and scheduled governance visibility
4. remediation effectiveness tracking

## Core Entities

1. `device_activation_orchestration_policy`
2. `device_activation_retry_budget`
3. `device_exception_escalation_state`
4. `device_exception_priority_rank`
5. `device_activation_recovery_decision`

## Proposed States

1. `retry_budget_available`
2. `retry_budget_exhausted`
3. `cooldown_active`
4. `fallback_branch_selected`
5. `escalation_notice_open`
6. `freeze_recommended`
7. `recovery_candidate`
8. `recovery_validated`

## Backend Work Plan

1. add activation orchestration read/write model under `telematics/devices`
2. add retry budget and cooldown policy evaluation
3. add exception recurrence escalation state
4. add freeze recommendation and controlled recovery hooks
5. extend governance export/report pack with orchestration closure quality

## Frontend Work Plan

1. expand device detail with orchestration card
2. add fleet priority desk to telematics dashboard
3. add orchestration batch macros
4. add recurrence escalation and recovery governance cards
5. add stock-side closure snapshot for orchestration governance

## Governance And Audit

1. activation orchestration audit stream
2. recurrence escalation notice chain
3. freeze recommendation audit
4. recovery validation timeline
5. export/report closure trend

## Acceptance Criteria

1. activation retry budget and fallback path are policy-driven
2. field exception recurrence escalations are auditable and batch-remediable
3. freeze recommendation and controlled recovery are governance-visible
4. orchestration governance is exportable and closure-tracked
5. Teltonika-first activation orchestration hardening is operationally visible

## Immediate Implementation Order

1. activation orchestration summary model
2. retry budget and cooldown projection
3. recurrence escalation write-model
4. fleet priority desk and batch macros
5. governance export and closure trend

## Risks

1. over-escalation may create noisy operations
2. retry budget policies may hide true field defects
3. recovery paths can become inconsistent across branches

## Mitigation

1. keep escalation thresholds explicit and cooldown-scoped
2. separate retry budget exhaustion from field exception severity
3. retain centralized audit/export/governance contract

## Exit Criteria

Phase-7 closes when:
- activation orchestration is policy-driven
- recurrence escalations are structured and batch-remediable
- freeze recommendation and recovery candidate states are visible
- governance export/report chain is live

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md)
- [TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md)
- [TELEMATICS_WAVE_A_PHASE_7_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_7_COMPLETION_REPORT_2026-03-28.md)
