# Login Modal Locale And Membership Scope Affordance Checkpoint

Date: 2026-04-12

## Closed In This Slice

- Login language picker moved to modal-based selection flow.
- Tenant selector language picker moved to modal-based selection flow.
- Login auth error fallback no longer leaks raw key for known auth failures.
- Launchpad now falls back to root-level navigable entries instead of showing an empty state too early.
- Dashboard now exposes clearer platform-admin vs tenant-admin scope affordances.
- Users permission surface remains the operational center for tenant/site membership management.

## Key Files

- `/workspace/node/aystek-core-next/app/(auth)/login/page.tsx`
- `/workspace/node/aystek-core-next/app/tenant-selector/page.tsx`
- `/workspace/node/aystek-core-next/components/i18n/LocaleModalPicker.tsx`
- `/workspace/node/aystek-core-next/components/layout/LaunchpadMenu.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/page.tsx`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Verification

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleModalPicker.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'app/(dashboard)/page.tsx' --file 'app/(dashboard)/admin/users/page.tsx' --file 'hooks/useMenu.ts' --file 'components/context/SiteSwitcher.tsx' --file 'lib/api/translations.ts' --file 'lib/api/translations.server.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Runtime Note

- A visible browser result may still require backend restart plus hard refresh if the active tab is holding stale translation/cache state.
