# Telematics Wave A Phase-3 Teltonika Onboarding And Provisioning Automation Execution Plan

Date: 2026-03-26
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- `TELEMATICS_WAVE_A_PHASE_2_COMPLETION_REPORT_2026-03-26.md`
- `TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md`

## Objective

Phase-3 hedefi, tamamlanan device lifecycle center’ı gerçek Teltonika onboarding ve provisioning otomasyonuna taşımaktır.

Bu dalga şu boşluğu kapatır:
- cihaz lifecycle center var ama onboarding akışı hâlâ yarı-manuel
- provisioning checklist görünür ama vendor playbook otomasyonu yok
- Teltonika cihazı için install-to-active dönüşümü operasyonel makroya bağlanmış değil

## Product Outcome

Teslim edilecek ana ürün yüzeyi:
- `Teltonika Onboarding And Provisioning Automation`

Bu yüzey şunları bağlar:
- discovered device
- IMEI identity verification
- stock/device assignment
- SIM/operator readiness
- provisioning checklist completion
- first telemetry validation
- certification gate pass
- activation handoff

## In Scope

1. Teltonika onboarding state machine
2. IMEI to stock/device assignment automation
3. provisioning checklist automation hooks
4. first telemetry / first heartbeat validation
5. SIM/operator readiness visibility
6. provisioning batch actions and failure buckets
7. vendor-specific lifecycle playbook registry
8. onboarding audit and export governance

## Out Of Scope

1. Full FOTA package delivery
2. Queclink production implementation
3. Mobile installer app
4. Regulatory packs
5. Video / sensor stacks

## Delivery Tracks

### Track P3-A: Teltonika Onboarding Domain

1. onboarding aggregate
2. onboarding states
3. assignment and readiness checks
4. playbook step status model

### Track P3-B: Provisioning Automation

1. stock-to-device binding helpers
2. SIM/operator validation hooks
3. first payload validation
4. activation-ready transition

### Track P3-C: Operations UI

1. onboarding queue
2. failed onboarding buckets
3. batch provisioning actions
4. playbook step detail drawer

### Track P3-D: Governance

1. onboarding lag KPI
2. provisioning failure recurrence
3. onboarding export/report pack
4. Teltonika-specific remediation notices

## Core Entities

1. `device_onboarding_snapshot`
2. `device_onboarding_event`
3. `device_onboarding_playbook`
4. `device_assignment_resolution`
5. `device_first_signal_validation`

## Proposed Onboarding States

1. `discovered`
2. `imei_verified`
3. `inventory_bound`
4. `sim_ready`
5. `provisioning_in_progress`
6. `waiting_first_signal`
7. `certification_pending`
8. `activation_ready`
9. `activated`
10. `blocked`

## Backend Work Plan

1. add onboarding projection under `telematics/devices`
2. persist onboarding events from:
   - Teltonika TCP listener
   - stock assignment actions
   - provisioning checklist completion
   - first telemetry validation
3. add playbook registry and resolution strategy
4. add onboarding queue and failure bucket summary
5. add batch remediation endpoints for blocked/waiting devices

## Frontend Work Plan

1. add onboarding queue facet into `admin/telematics/devices`
2. expand device detail with onboarding block
3. add Teltonika-specific playbook card
4. add onboarding failure and readiness filters
5. add batch onboarding remediation controls

## Governance And Audit

1. onboarding action audit stream
2. assignment mismatch audit
3. first-signal timeout notice
4. provisioning blocked recurrence KPI
5. onboarding export/report visibility

## Acceptance Criteria

1. Teltonika devices can be filtered by onboarding state
2. IMEI verification, assignment, first signal and certification are visible in one flow
3. blocked onboarding devices can be batch-remediated
4. onboarding lag and first-signal failure are audited and exportable
5. vendor-specific playbook guidance exists for Teltonika lifecycle steps

## Immediate Implementation Order

1. onboarding summary model
2. onboarding queue and filters
3. Teltonika playbook state block
4. first-signal validation and blocked buckets
5. batch remediation and governance

## Risks

1. stock/device ownership truth may drift during onboarding
2. SIM readiness may vary by operator workflow
3. first signal timeout semantics may differ by device class

## Mitigation

1. keep assignment resolution source-attributed
2. separate inventory readiness from telemetry readiness
3. model first-signal policy as configurable per playbook

## Exit Criteria

Phase-3 closes when:
- Teltonika onboarding queue is live
- provisioning automation exists for common happy path
- blocked/waiting devices are batch-remediable
- onboarding lag and first-signal pressure are governed

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_2_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_2_COMPLETION_REPORT_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md)
