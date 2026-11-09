# Bulk SMS System Preset Rules And Time Window Governance Checkpoint

Date: 2026-04-23

## Completed

- Bulk SMS governance endpoint now accepts a time-window parameter (`days`).
- Governance time-series can now be refreshed for different windows from the UI.
- Bulk SMS governance summary now supports richer breakdowns already added in the previous slice, and the selected time window is now visible in the page workflow.
- System preset lifecycle rules are now explicitly visible in UI:
  - system presets are shown as automatically managed
  - tenant presets are shown as archivably governed presets
  - counts for system vs tenant presets are displayed

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx' --file 'lib/api/sms.ts' --file 'lib/types/sms.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/modules/users/... ./internal/seeds ./cmd/api`

## Current State

- Bulk SMS preset lifecycle is now more visible and less implicit.
- Governance reporting is more usable because the operator can switch the reporting window.
- Ownership, breakdowns, and time-window reporting now work together as one coherent operational surface.

## Next Logical Step

- Add provider/template trend comparison cards across selected windows.
- Add stronger governance drill-down filters for approval history.
- If needed, separate governance reporting into a dedicated analytics/report surface.
