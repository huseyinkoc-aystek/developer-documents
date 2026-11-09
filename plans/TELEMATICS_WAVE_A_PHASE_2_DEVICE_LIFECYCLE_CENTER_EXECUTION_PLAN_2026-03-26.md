# Telematics Wave A Phase-2 Device Lifecycle Center Execution Plan

Date: 2026-03-26
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- `TELEMATICS_WAVE_A_PHASE_1_COMPLETION_REPORT_2026-03-26.md`
- `TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md`

## Objective

Phase-2 hedefi, adapter foundation ve certification operasyonunu gerçek cihaz yaşam döngüsü ürün yüzeyine çevirmektir.

Bu dalga şu boşluğu kapatır:
- cihaz onboard edildi ama operasyonel yaşam döngüsü görünmüyor
- provisioning, firmware, SIM, health ve certification aynı yüzeyde birleşmiyor
- vendor adapter artık sadece parser/certification nesnesi; henüz gerçek device center ürünü değil

## Product Outcome

Teslim edilecek ana ürün yüzeyi:
- `Device Lifecycle Center`

Bu yüzey tek cihaz ve cihaz filosu için şunları toplar:
- identity
- provisioning
- connectivity
- firmware
- SIM/operator
- certification
- protocol health
- last seen/session
- anomaly/remediation

## In Scope

1. Device lifecycle read model
2. Device detail center
3. Device provisioning state machine
4. Certification state and approval visibility
5. Firmware/FOTA metadata surface
6. SIM/operator and line health surface
7. Tenant/site scoped device health KPI
8. Device action timeline and remediation notes

## Out of Scope

1. Full FOTA binary delivery engine
2. Mobile installer application
3. Video telematics
4. Regulatory packs
5. Queclink production adapter

## Delivery Tracks

### Track P2-A: Lifecycle Domain

1. device lifecycle aggregate
2. provisioning states
3. lifecycle notes and manual override events
4. lifecycle health summary

### Track P2-B: Device Center UI

1. device list with lifecycle facets
2. device detail page
3. timeline and health cards
4. remediation and governance drawer

### Track P2-C: Provisioning And Certification

1. onboarding checklist
2. certification approval visibility
3. failed certification retry flow
4. freeze / unfreeze / approve transitions

### Track P2-D: Fleet Operations

1. fleet-level device lifecycle KPI
2. stale devices
3. uncertified devices
4. provisioning blocked devices

## Core Entities

1. `device_lifecycle_snapshot`
2. `device_lifecycle_event`
3. `device_provisioning_state`
4. `device_health_summary`
5. `device_certification_gate`

## Proposed Lifecycle States

1. `discovered`
2. `identified`
3. `assigned`
4. `provisioning_pending`
5. `provisioning_blocked`
6. `active`
7. `degraded`
8. `uncertified`
9. `frozen`
10. `retired`

## Backend Work Plan

1. Add lifecycle projection under `telematics/devices`
2. Persist lifecycle events from:
   - TCP listener
   - certification actions
   - freeze/unfreeze actions
   - provisioning actions
3. Add summary endpoint for device lifecycle KPI
4. Add device detail endpoint with timeline blocks
5. Add lifecycle note / acknowledge / override actions
6. Add certification approval state write-model

## Frontend Work Plan

1. Add `admin/telematics/devices` lifecycle list filters
2. Expand device detail into lifecycle center
3. Add provisioning checklist card
4. Add certification gate card
5. Add firmware/SIM/connectivity block
6. Add lifecycle timeline with remediation actions

## Governance And Audit

1. lifecycle action audit stream
2. provisioning block / unblock audit
3. certification approval / denial stream
4. stale device escalation KPI
5. uncertified device pressure KPI

## Acceptance Criteria

1. Device list can separate active / stale / uncertified / frozen devices
2. Device detail shows last seen, certification, firmware, SIM and provisioning state together
3. Manual lifecycle actions are audited
4. Fleet dashboard shows lifecycle health pressure
5. Protocol adapter certification remains linked to concrete device lifecycle state

## Immediate Implementation Order

1. device lifecycle summary model
2. fleet KPI and filtered device list
3. device detail lifecycle timeline
4. provisioning checklist and notes
5. certification gate and approval visibility

## Risks

1. Device truth may split across telemetry, stock and device registry
2. Provisioning semantics may vary by vendor
3. Certification and provisioning may drift out of sync

## Mitigation

1. Keep a single lifecycle projection with source attribution
2. Model provisioning state vendor-agnostic first
3. Emit lifecycle events from every certification state transition

## Exit Criteria

Phase-2 closes when:
- device lifecycle center is live
- fleet KPI exists
- provisioning + certification are unified
- lifecycle remediation is auditable

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_1_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_1_COMPLETION_REPORT_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_2_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_2_COMPLETION_REPORT_2026-03-26.md)
