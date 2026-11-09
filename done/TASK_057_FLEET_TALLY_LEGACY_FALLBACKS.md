# TASK 057 - Fleet Tally Legacy Fallbacks and Push Tracking

## Scope

Completed the remaining hardening items from the GPS Cetele migration audit.

## Changes

- Added legacy route-pass fallback to `point_in_circle_multi_window`.
  - If no point is inside the target circle, the algorithm checks nearby route segments.
  - Default behavior follows the legacy PHP flow:
    - fallback window: 30 minutes
    - fallback radius: `radius * 2`
    - minimum segment time: 60 seconds
    - minimum segment distance: configured value or original radius
- Added normalized status `route_passed`, mapped to legacy SAP code `10`.
- Added schedule-level transient verification retry for technical/unknown statuses.
  - Uses `VerificationSchedule.RetryPolicy.MaxAttempts`.
  - Uses `VerificationSchedule.RetryPolicy.BackoffSeconds`.
  - Persists verification retry count on the job.
- Added per-job push-back status tracking after connector push.
  - Successful batch push marks jobs as `sent`.
  - Failed batch push marks jobs as `failed`.
  - Push attempts are incremented per job.
- Verification status persistence now stores GPS provider and control timestamps with the job status update.
- Already verified and already pushed jobs are not pushed again.

## Validation

```bash
go test ./internal/modules/fleet_tally/...
go test ./internal/seeds ./cmd/api
```

Result: passed.
