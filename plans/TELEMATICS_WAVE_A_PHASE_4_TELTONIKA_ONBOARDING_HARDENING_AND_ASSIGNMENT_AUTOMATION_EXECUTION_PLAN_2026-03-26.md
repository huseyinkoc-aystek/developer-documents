# Telematics Wave A Phase-4 Teltonika Onboarding Hardening And Assignment Automation Execution Plan

Date: 2026-03-26
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- `TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md`
- `TELEMATICS_WAVE_A_PHASE_3_TELTONIKA_ONBOARDING_AND_PROVISIONING_AUTOMATION_EXECUTION_PLAN_2026-03-26.md`
- `TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md`

## Objective

Phase-4 hedefi, tamamlanan Teltonika onboarding operasyon yüzeyini gerçek assignment automation ve hardening katmanına taşımaktır.

Bu dalga şu boşluğu kapatır:
- onboarding queue ve governance var ama assignment hâlâ operatör ağırlıklı
- playbook skeleton var ama vendor-specific remediation otomasyonu sınırlı
- first signal ve activation-ready geçişi görünür ama tam otomatik kapanış zinciri yok

## Product Outcome

Teslim edilecek ana ürün yüzeyi:
- `Teltonika Onboarding Hardening And Assignment Automation`

Bu yüzey şunları bağlar:
- IMEI to inventory resolution
- automatic stock/device binding
- assignment mismatch detection
- SIM/operator readiness confirmation
- first signal timeout automation
- activation-ready transition
- vendor-specific playbook enforcement

## In Scope

1. IMEI to stock/device assignment resolver
2. assignment confidence and mismatch model
3. Teltonika onboarding hardening policy
4. SIM/operator readiness hook surface
5. first signal timeout policy and remediation
6. activation-ready automation
7. vendor-specific onboarding playbook expansion
8. assignment governance, audit and export visibility

## Out Of Scope

1. Full FOTA or firmware rollout
2. Queclink production implementation
3. Mobile installer app
4. Full SIM procurement workflow
5. Multi-vendor provisioning parity

## Delivery Tracks

### Track P4-A: Assignment Automation

1. IMEI resolver
2. stock/device candidate scoring
3. automatic bind or operator-confirm path
4. assignment mismatch drift audit

### Track P4-B: Onboarding Hardening

1. first signal timeout evaluator
2. blocked bucket enrichment
3. assignment retry policy
4. activation-ready policy

### Track P4-C: Vendor Playbooks

1. Teltonika-specific playbook steps
2. SIM/operator readiness checklist
3. first payload verification step
4. activation handoff step

### Track P4-D: Governance

1. assignment pressure KPI
2. mismatch recurrence notice
3. assignment/export governance
4. auto-report and closure trend

## Core Entities

1. `device_assignment_resolution`
2. `device_assignment_candidate`
3. `device_assignment_mismatch`
4. `device_first_signal_policy`
5. `device_activation_gate`

## Proposed Assignment States

1. `unresolved`
2. `candidate_found`
3. `auto_bound`
4. `operator_confirmation_required`
5. `mismatch_detected`
6. `sim_pending`
7. `first_signal_pending`
8. `activation_ready`
9. `activated`
10. `blocked`

## Backend Work Plan

1. add assignment resolution write-model under `telematics/devices`
2. persist candidate scoring and chosen source attribution
3. add mismatch audit and recurrence policy
4. add first-signal timeout evaluator and automation hooks
5. add activation-ready transition writer
6. add assignment governance export/report pack

## Frontend Work Plan

1. expand `admin/telematics/devices` with assignment filters
2. add assignment confidence and mismatch cards to device detail
3. add Teltonika hardening playbook section
4. add assignment remediation macros
5. add assignment governance dashboard block

## Governance And Audit

1. assignment mismatch audit stream
2. first signal timeout notice
3. activation-ready transition audit
4. assignment export/report visibility
5. recurrence and closure KPI

## Acceptance Criteria

1. Teltonika devices can be auto-bound or clearly routed to operator confirmation
2. assignment mismatch and first-signal timeout are governed with notices
3. activation-ready transition is visible and automatable
4. assignment remediation is batch-operable
5. assignment governance is exportable and closure-tracked

## Immediate Implementation Order

1. assignment summary model
2. candidate scoring and mismatch projection
3. device detail assignment card
4. first-signal timeout automation
5. assignment governance and export chain

## Risks

1. inventory truth may drift while onboarding is active
2. IMEI reuse or malformed device IDs can degrade confidence
3. operator workflow differences may affect activation timing

## Mitigation

1. keep assignment source-attributed and reversible
2. separate `auto_bound` from `confirmation_required`
3. model first-signal and activation policy as vendor-configurable

## Exit Criteria

Phase-4 closes when:
- Teltonika assignment automation exists for the common happy path
- mismatch and first-signal timeout are governed
- activation-ready transition is observable and automatable
- assignment export/report chain is live

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_3_TELTONIKA_ONBOARDING_AND_PROVISIONING_AUTOMATION_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_3_TELTONIKA_ONBOARDING_AND_PROVISIONING_AUTOMATION_EXECUTION_PLAN_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md)
