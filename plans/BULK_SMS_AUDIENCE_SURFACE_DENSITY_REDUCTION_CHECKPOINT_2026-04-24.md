# Bulk SMS Audience Surface Density Reduction Checkpoint

Date: 2026-04-24

## Context

`/admin/integrations/bulk-sms` audience surface still felt crowded and visually unstable. The main issue was not a single spacing bug; it was the information architecture:

- audience filtering
- page selection / pagination
- refresh / reset operations
- active locale / active preset context
- audience table

were competing inside a compressed layout.

## What Changed

- Audience tab was restructured into three clear layers:
  - audience summary
  - filter workspace
  - filtered audience table
- The old sidebar-style filter stack was removed.
- Search, approval filter, and language filter now sit in a wider top row.
- Pagination, selection actions, and operational actions are grouped into separate cards instead of one cramped toolbar.
- Active locale and active preset context were moved into a calmer header summary.

## i18n Alignment

New keys added:

- `sms.bulk.actions.selection`
- `sms.bulk.actions.operations`
- `sms.bulk.actions.open_composer`

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Result

The audience surface now behaves more like an operations workspace and less like a compressed form toolbar. This reduces visual noise and makes the selected-tab contract clearer.
