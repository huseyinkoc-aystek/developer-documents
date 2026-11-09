# Telematics Wave A Phase-6 Vendor Activation Automation And Field Exception Handling Execution Plan

Date: 2026-03-27
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- [TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md)
- [TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md)
- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)

## Objective

Phase-6 hedefi, Teltonika ağırlıklı provisioning ve activation akışını daha vendor-aware, saha istisnalarına dayanıklı ve policy-driven hale getirmektir.

Bu faz şu boşlukları kapatır:
- activation-ready görünürlüğü var ama vendor-specific automation branches sınırlı
- saha kaynaklı exception state’leri auditleniyor ama structured field exception workbench eksik
- activation başarısızlıkları için controlled retry, hold ve escalate ayrımı tam ürünleşmiş değil

## Product Outcome

Teslim edilecek ana yüzey:
- `Vendor Activation Automation And Field Exception Control`

Bu yüzey şunları bağlar:
- vendor-specific activation transitions
- activation hold / retry / manual handoff policy
- field exception classification
- field exception remediation macros
- activation exception governance, closure trend and export visibility

## In Scope

1. vendor-specific activation policy branches
2. activation retry / hold / manual handoff model
3. field exception taxonomy
4. exception remediation note and action history
5. activation failure and field exception notice chain
6. exception governance export and closure trend
7. Teltonika-first vendor playbook hardening

## Out Of Scope

1. full technician mobile workflow
2. warehouse/installer dispatch optimization
3. FOTA command orchestration
4. multi-vendor full parity rollout
5. carrier contract automation

## Delivery Tracks

### Track P6-A: Vendor Activation Automation

1. vendor activation branch model
2. auto-advance / hold / rollback path
3. activation retry budget and cooldown
4. activation closure quality

### Track P6-B: Field Exception Workbench

1. field exception category model
2. exception severity and ownership
3. remediation action log
4. manual handoff / return-to-queue lifecycle

### Track P6-C: Governance

1. activation exception KPI
2. exception notice chain
3. export and scheduled governance visibility
4. recurrence and closure trend

## Core Entities

1. `device_activation_branch`
2. `device_activation_exception`
3. `device_field_exception`
4. `device_exception_remediation`
5. `device_activation_exception_notice`

## Proposed Exception States

1. `awaiting_activation`
2. `activation_retry_scheduled`
3. `activation_hold`
4. `manual_handoff_required`
5. `field_exception_open`
6. `field_exception_acknowledged`
7. `field_exception_resolved`
8. `field_exception_recurring`

## Backend Work Plan

1. add activation branch write-model under `telematics/devices`
2. add vendor-specific activation retry/hold policy
3. persist field exception history and remediation actions
4. add activation exception notice and recurrence chain
5. add field exception governance export/report pack

## Frontend Work Plan

1. expand `admin/telematics/devices` with activation branch and exception filters
2. add field exception card to device detail
3. add activation retry / hold / handoff macros
4. add activation exception governance dashboard block
5. add closure trend and export pressure visibility

## Governance And Audit

1. activation branch audit stream
2. field exception remediation audit
3. activation exception notice chain
4. export/report visibility
5. recurrence and closure KPI

## Acceptance Criteria

1. vendor-specific activation branches are policy-driven
2. activation retry / hold / manual handoff is auditable
3. field exceptions are batch-remediable and severity-classified
4. activation exception governance is exportable and closure-tracked
5. Teltonika-first vendor activation hardening is operationally visible

## Immediate Implementation Order

1. activation branch summary model
2. field exception projection and device detail card
3. activation retry / hold write-model
4. exception notice and recurrence chain
5. governance export and scheduled visibility

## Risks

1. field exceptions can be noisy and overlap with provisioning issues
2. over-automation may hide real installer-side defects
3. vendor-specific branches can fragment operational UX

## Mitigation

1. keep exception taxonomy explicit and severity-scoped
2. separate retry automation from manual handoff
3. retain centralized audit/export/governance contract

## Exit Criteria

Phase-6 closes when:
- vendor activation automation is policy-driven
- field exceptions are structured and remediation-ready
- activation exception pressure is visible and batch-remediable
- governance export/report chain is live

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md)
- [TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md)
- [TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md)
