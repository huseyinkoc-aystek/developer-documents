# TASK 056 - Fleet Tally Audit and Runtime Hardening

## Scope

GPS Cetele / Fleet Tally migration was reviewed against the legacy PHP and Python flow:

- `GpsserviceController.php`
- `Konum.php`
- `PersonelStartGPSControl.py`
- `SchootStartGPSControl.py`

The current AYSTEK Core module already follows the intended enterprise direction: connector adapters, GPS provider adapters, vehicle mappings, schedules, verification jobs, and pluggable algorithms live under one module instead of duplicated customer-specific modules.

## Fixes

- Fixed Fleet Tally cron worker so both standard 5-field expressions (`0 6 * * *`) and optional-second 6-field expressions (`0 0 6 * * *`) are accepted.
- Preserved SAP raw payload on `VerificationJob.RawExternalPayload` so provider-specific fields such as `ZgpsNo` / `ZgpsDatano` remain available during GPS execution.
- Resolved provider-specific vehicle external ID for all providers through vehicle mappings first, then SAP raw payload fallback.
- Applied GPS provider timezone offset to verification fetch requests, with Turkey default `+180` when no provider/schema offset is configured.
- Added idempotent schedule job creation by checking existing jobs through `schedule_id + external_ref.doc_no + external_ref.line_no` before insert.
- Mapped school service refs (`belge_no`, `belge_kalem_no`) into the normalized external ref fields.
- Reordered GPS provider execution by SAP `Zgps` hint when available, while keeping the configured provider priority as fallback.
- Fixed duplicate package declaration in `mobiliz_test.go`.
- Added regression tests for cron parsing and raw payload/external-ref normalization.

## Validation

```bash
go test ./internal/modules/fleet_tally/...
```

Result: passed.

## Remaining Follow-Up

- Legacy `SSAEmptyLocation = 10` route-pass fallback is only partially represented by the separate `route_intersection_circle` algorithm; it is not yet wired as a fallback after point-in-circle failure.
- Legacy retry behavior for transient unknown provider statuses should be modeled as a schedule retry policy rather than ad-hoc sleeps.
- Push-back status is currently best-effort at schedule level; per-job push success/failure tracking should be tightened.
- Device inventory semantics from the old `SatilanUrunler::cihazBilgisi` branch should be mapped deliberately to stock/vehicle assets or vehicle mappings.
