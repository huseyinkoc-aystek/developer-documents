# Device Command Governance And Approval Slice Checkpoint

Date: 2026-04-19

## Scope

This slice advanced the stock device command flow from basic policy gating to a more explicit governance and approval-like model.

## Completed

- Added policy audit/history persistence for device command policy changes and override executions.
- Added a dedicated policy audit read endpoint and exposed the history on the admin registry page.
- Tightened override execution rules.
  - Matching route policy must allow override.
  - `override_reason` is required.
  - `override_ticket_ref` is required.
  - `stock_sms:override_approve` permission is required in addition to route override capability.
- Persisted override governance metadata on command executions.
  - `override_used`
  - `override_reason`
  - `override_ticket_ref`
- Tightened dashboard and launchpad visibility with membership scope plus command permissions.

## User-Facing Outcome

- Command operators now see a visible policy audit/history surface instead of silent policy changes.
- Override is no longer a casual fallback path. It behaves like an approval-style exception flow.
- Device command administration surfaces appear only when the user has both the correct permission set and the relevant scope context.

## Verification

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/new/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/[id]/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'app/(dashboard)/page.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'
```

## Next Step

- Extend governance from policy/event history toward a fuller approval workflow.
- Add richer audit filtering and operator-facing traceability for policy and override decisions.
- Tighten visibility and actionability further with combined membership, permission, and command-governance scope.
