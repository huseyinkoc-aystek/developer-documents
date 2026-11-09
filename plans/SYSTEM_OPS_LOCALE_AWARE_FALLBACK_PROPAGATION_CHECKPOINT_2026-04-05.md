# System Ops Locale-Aware Fallback Propagation Checkpoint

Date: 2026-04-05

## Context
- System Ops Next surfaces were partially translated because older `menu.system_ops.*` keys were already seeded, while newer `system_ops.*` keys could still fall back to English literals on the frontend.
- This created mixed-language screens under Turkish locale, especially on overview/dashboard-style surfaces.

## Completed
- Added locale-aware fallback helper usage across shared System Ops display and page surfaces.
- Extended locale-aware fallback propagation to:
  - dashboard
  - overview
  - health
  - node
  - version
  - backups
  - maintenance
  - policies
  - migrations
  - updates list
  - update start
  - update detail
  - rollback start
- Updated shared display normalization to use locale-aware TR/EN fallbacks for:
  - status labels
  - backup type labels
  - trigger source labels

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/overview/page.tsx' --file 'app/(dashboard)/admin/system/health/page.tsx' --file 'app/(dashboard)/admin/system/node/page.tsx' --file 'app/(dashboard)/admin/system/version/page.tsx' --file 'app/(dashboard)/admin/system/backups/page.tsx' --file 'app/(dashboard)/admin/system/maintenance/page.tsx' --file 'app/(dashboard)/admin/system/policies/page.tsx' --file 'app/(dashboard)/admin/system/migrations/page.tsx' --file 'app/(dashboard)/admin/system/updates/page.tsx' --file 'app/(dashboard)/admin/system/updates/start/page.tsx' --file 'app/(dashboard)/admin/system/updates/[id]/page.tsx' --file 'app/(dashboard)/admin/system/rollbacks/start/page.tsx' --file 'components/system-ops/SystemOpsStatusBadge.tsx' --file 'lib/system-ops/display.ts' --file 'lib/system-ops/text.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Result
- System Ops pages no longer depend solely on seed timing to display Turkish UI text.
- When DB translation rows are missing or not yet reconciled, frontend fallbacks now remain locale-aware instead of reverting to English by default.
