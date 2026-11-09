# Telematics Wave A Phase-5 First-Signal And Provisioning Hardening Execution Plan

Date: 2026-03-27
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Depends on:
- `TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md`
- `TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md`
- `TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md`

## Objective

Phase-5 hedefi, Teltonika onboarding ve assignment otomasyonundan sonra kalan provisioning hardening boşluklarını kapatmaktır.

Bu dalga şu boşluğu kapatır:
- assignment automation var ama first-signal timeout akışı tam policy-driven değil
- activation-ready görünürlüğü var ama otomatik geçiş ve closure lifecycle eksik
- vendor-specific provisioning playbook var ama hard fail / timeout / escalation zinciri sınırlı

## Product Outcome

Teslim edilecek ana ürün yüzeyi:
- `Teltonika First-Signal And Provisioning Hardening`

Bu yüzey şunları bağlar:
- first-signal timeout detection
- activation-ready automation
- provisioning blocked escalation
- vendor-specific readiness gate enforcement
- provisioning closure trend and governance

## In Scope

1. first-signal timeout evaluator
2. activation-ready transition automation
3. SIM/operator readiness hardening hooks
4. provisioning blocked severity model
5. vendor-specific Teltonika provisioning playbook expansion
6. provisioning governance KPI, notice and recurrence chain
7. audit/export/scheduled governance visibility

## Out Of Scope

1. Full firmware/FOTA orchestration
2. Queclink production rollout
3. Installer mobile workflow
4. Multi-vendor provisioning parity
5. Deep carrier procurement integration

## Delivery Tracks

### Track P5-A: First-Signal Control

1. first-signal deadline policy
2. timeout buckets and severity
3. retry/remediation macro
4. timeout closure classification

### Track P5-B: Activation Automation

1. activation-ready state machine hardening
2. auto-transition trigger
3. manual override / rollback path
4. activation audit and closure quality

### Track P5-C: Provisioning Playbooks

1. Teltonika-specific readiness steps
2. SIM/operator readiness enrichment
3. first payload verification hardening
4. provisioning blocked remediation notes

### Track P5-D: Governance

1. first-signal governance KPI
2. provisioning blocked notice chain
3. export/report and auto-report visibility
4. recurrence and closure trend

## Core Entities

1. `device_first_signal_window`
2. `device_provisioning_gate`
3. `device_activation_transition`
4. `device_provisioning_notice`
5. `device_provisioning_batch_outcome`

## Proposed Provisioning States

1. `awaiting_signal`
2. `signal_timeout_risk`
3. `signal_timeout`
4. `readiness_pending`
5. `activation_ready`
6. `activation_in_progress`
7. `activated`
8. `blocked`
9. `manually_overridden`

## Backend Work Plan

1. add first-signal timeout write-model under `telematics/devices`
2. persist signal deadline and timeout severity snapshots
3. add activation-ready automation writer and rollback trace
4. add provisioning blocked notice and recurrence policy
5. add provisioning governance export/report pack

## Frontend Work Plan

1. expand `admin/telematics/devices` with first-signal and activation filters
2. add provisioning hardening card to device detail
3. add timeout and activation remediation macros
4. add provisioning governance dashboard block
5. add closure trend and export pressure visibility

## Governance And Audit

1. first-signal timeout audit stream
2. activation-ready transition audit
3. provisioning blocked notice chain
4. export/report visibility
5. recurrence and closure KPI

## Acceptance Criteria

1. first-signal timeout is policy-driven and visible
2. activation-ready transition is automatable and auditable
3. provisioning blocked devices are batch-remediable
4. provisioning governance is exportable and closure-tracked
5. vendor-specific readiness enforcement is operationally visible

## Immediate Implementation Order

1. first-signal summary model
2. timeout projection and device detail hardening card
3. activation-ready automation
4. provisioning notice and recurrence chain
5. governance export and scheduled visibility

## Risks

1. signal latency can vary by environment and vendor deployment quality
2. operator-dependent provisioning steps may still create false positives
3. aggressive automation may hide field-level issues without strong audit

## Mitigation

1. keep timeout policy tenant/system configurable
2. separate timeout detection from activation execution
3. require closure audit for auto-transition and retries

## Exit Criteria

Phase-5 closes when:
- first-signal timeout is governed end-to-end
- activation-ready transition is automatable and reversible
- provisioning blocked pressure is visible and batch-remediable
- governance export/report chain is live

## References

- [TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_DEVICE_AND_PROTOCOL_PRODUCTIZATION_EXECUTION_PLAN_2026-03-25.md)
- [TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md)
- [TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md)
