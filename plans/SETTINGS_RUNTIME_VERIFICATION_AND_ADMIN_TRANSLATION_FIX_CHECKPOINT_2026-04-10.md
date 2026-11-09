# Settings Runtime Verification And Admin Translation Fix Checkpoint

Date: 2026-04-10

## Scope

- Runtime verification follow-up for `settings/dashboard`, `settings`, `settings/tenant-site`, `shipping/settings`
- Admin translation hydration root-cause analysis
- Seed reconcile confirmation

## Root Cause

Admin surfaces were not failing because `settings.*` keys were absent from the seed set.

The actual problem was translation loading order:

1. Client-side `translationsApi.fetch(...)` fetched `/public/translations` first.
2. If public translations returned any payload at all, the function returned early.
3. Admin-only keys such as `settings.*` remained unloaded even though they existed in DB.
4. Server-side dashboard preload had the same issue because it only used `/public/translations`.

Result: authenticated admin pages could render raw keys even when translation records were present and seeded.

## Fix Applied

### Next client runtime

Updated:

- `/workspace/node/aystek-core-next/lib/api/translations.ts`

Behavior now:

- `publicOnly=true`: still fetches only public translations
- admin/default mode: fetches public translations, then fetches admin translations, and merges them
- admin dictionary overrides public dictionary when both exist

### Next server preload

Updated:

- `/workspace/node/aystek-core-next/lib/api/translations.server.ts`

Behavior now:

- SSR preload for dashboard/admin surfaces merges `public/translations` and `admin/translations`
- admin pages no longer start with a public-only dictionary

## Runtime Verification Notes

### Seed reconcile

Seed reconcile was run through `cmd/api` with seed exit mode.

Observed log evidence:

- `Translations seeded (tenant=default, processed=8510, inserted=60, existing=8450)`
- `Translations seeded (tenant=aystek, processed=8510, inserted=60, existing=8450)`
- later tenant runs continued with the same insert pattern

This confirms newly added translation keys were reconciled into the database during the runtime verification pass.

### settings namespace state

The latest `settings.*` reverification matrix already showed:

- UI-used real `settings.*` keys: `186`
- missing seeded keys: `0`

Reference:

- `/workspace/aystek-core-dev-docs/plans/SETTINGS_NAMESPACE_I18N_REVERIFICATION_MATRIX_2026-04-10.md`

## Validation

### Next

- `npm run lint -- --file 'lib/api/translations.ts' --file 'lib/api/translations.server.ts' --file 'app/components/LanguageInitializer.tsx' --file 'app/(dashboard)/layout.tsx' --file 'app/(dashboard)/admin/settings/page.tsx' --file 'app/(dashboard)/admin/settings/dashboard/page.tsx' --file 'app/(dashboard)/admin/settings/tenant-site/page.tsx' --file 'app/(dashboard)/admin/shipping/settings/page.tsx'`

Result:

- passed

### Go

- `go test ./internal/seeds ./cmd/api`

Result:

- passed

## Residual Risk

If raw keys are still visible in a live browser after this fix, the remaining likely causes are operational rather than code/seed related:

- stale frontend bundle
- old browser tab state
- translation cache not yet refreshed in the current session
- backend instance not restarted after seed reconcile

## Conclusion

`settings.*` raw-key behavior is no longer explained by missing seed entries.

The root cause was admin surfaces receiving a public-only translation dictionary.
That behavior is now fixed in both client runtime loading and server-side preload.
