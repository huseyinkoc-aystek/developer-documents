# System Ops Module Hardening Round 15 Checkpoint

Date: 2026-04-04

## Completed In This Round

- Normalized remaining short hard-coded helper strings in legacy System Ops wrappers.
- Added missing translation keys to seed data for newly introduced shared helpers and operational text.

## Updated Frontend Components

- `components/system-ops/MaintenanceBanner.tsx`
- `components/system-ops/LiveLogTerminal.tsx`
- `components/system-ops/UpdateProgressCard.tsx`

## Updated Translation Seed

- `go/aystek-core/internal/seeds/translation_seed.go`

## Translation Areas Covered

- maintenance banner title and default fallback message
- manage / dismiss shared action labels
- log terminal title, live label, reconnect, waiting, and empty-output messages
- update progress card title and detail links
- shared form helper support keys introduced in recent rounds:
  - locale fallback help
  - channel help
  - rollback reason help
  - dashboard / overview module tile help
  - license key label

## Outcome

- System Ops multilingual consistency is now tighter across both page surfaces and legacy compatibility wrappers.
- Newly introduced shared helpers no longer depend on undeclared translation keys.
- Remaining plain strings in the module are mostly structural values, technical enums, route paths, or acceptable fallback text inside `t(...)`.

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'components/system-ops/MaintenanceBanner.tsx' --file 'components/system-ops/LiveLogTerminal.tsx' --file 'components/system-ops/UpdateProgressCard.tsx' --file 'app/(dashboard)/admin/system/maintenance/page.tsx' --file 'app/(dashboard)/admin/system/updates/start/page.tsx' --file 'app/(dashboard)/admin/system/rollbacks/start/page.tsx' --file 'app/(dashboard)/admin/system/license/page.tsx' --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/overview/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/system_ops/... ./internal/seeds ./cmd/api`

## Next Logical Step

- Run one more focused pass on `components/system-ops/` to decide which remaining legacy wrappers should be fully retired versus kept as compatibility shims.
- If the module will grow further, add a small i18n audit checklist into the System Ops execution plan so new helper components always ship with seed coverage.
