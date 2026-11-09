# Telematics Wave A / Phase-9 Export Notice Checkpoint

Date: 2026-03-28

Checkpoint scope:
- `freeze_closure_governance` auto-report pressure icin export-side notice emission acildi
- stock settings notice policy ve stock dashboard last notice gorunurlugu eklendi
- merkezi audit focus ve telematics dashboard export notice KPI gorunurlugu eklendi

Implemented keys:
- `freeze_closure_governance.alert_enabled`
- `freeze_closure_governance.alert_queue_threshold`
- `freeze_closure_governance.alert_failed_threshold`
- `freeze_closure_governance.alert_cooldown_sec`
- `freeze_closure_governance.alert.last_state`
- `freeze_closure_governance.alert.last_at`

Event:
- `telematics.devices.freeze_closure_governance_report_alert_emitted`

Planned next step:
- export notice closure lifecycle: `acknowledge`, `recheck`, `resolved / still_open`
