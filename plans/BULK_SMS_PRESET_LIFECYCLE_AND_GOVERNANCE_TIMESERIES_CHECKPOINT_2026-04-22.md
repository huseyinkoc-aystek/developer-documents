# Bulk SMS Preset Lifecycle And Governance Timeseries Checkpoint

Date: 2026-04-22

## Completed

- Bulk SMS governance API now returns `series` data and Next parses it correctly.
- Bulk SMS page now sends `segment_preset_id` with bulk send requests.
- Segment presets now expose lifecycle metadata in UI:
  - status
  - lifecycle rule
  - usage count
  - last used at
- Non-system presets can now be archived/restored from the UI.
- Archived presets are separated from active presets in the page.
- Approval history UI now shows:
  - approver role
  - provider label
  - template name
- Governance time-series is now visible as a simple, readable per-point chart/list.
- Missing `sms.bulk.*` lifecycle/governance/history translation keys were added to seed.

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx' --file 'lib/api/sms.ts' --file 'lib/types/sms.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/modules/users/... ./internal/seeds ./cmd/api`

## Current State

- Bulk SMS is now beyond a basic send screen and has visible lifecycle/governance behavior.
- Segment presets are operationally manageable.
- Approval history is more readable and less ID-centric.
- Governance view now includes a time-distribution signal instead of only static counters.

## Next Logical Step

- Carry preset lifecycle into stricter governance:
  - system preset management rules
  - preset ownership visibility
  - richer governance reporting by provider/template/approver role
- If needed, move bulk campaign history and approval history into separate dedicated drill-down surfaces.
