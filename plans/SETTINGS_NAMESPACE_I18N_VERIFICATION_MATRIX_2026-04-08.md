# Settings Namespace I18n Verification Matrix

Date: 2026-04-08

## Scope

This verification pass covers the `settings.*` namespace used by the Next admin settings surfaces.

Scanned areas:

- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/page.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/dashboard/page.tsx`

Seed source:

- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Result

- UI-used `settings.*` keys: `164`
- Missing keys in seed before remediation: `7`
- Missing keys after remediation: `0`

## Missing Keys That Were Added

- `settings.brand.logoPreview`
- `settings.brand.faviconPreview`
- `settings.brand.ogPreview`
- `settings.localization.titleRequired`
- `settings.maintenance.error`
- `settings.notice.loadFail`
- `settings.notice.loadFailTitle`

## File Coverage

- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/page.tsx`
  - key count: `145`
  - missing after remediation: `0`

- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/dashboard/page.tsx`
  - key count: `19`
  - missing after remediation: `0`

## Notes

- This matrix verifies seed presence, not runtime cache state.
- Runtime visibility still depends on the translation reconcile path being executed against the active database.
- The earlier admin hydration issue was already addressed by ensuring admin surfaces do not fetch `publicOnly` translations.

## Verification Commands

- `comm -23 /workspace/settings_ui_keys.txt <(python3 - <<'PY' ... )`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/settings/page.tsx' --file 'app/(dashboard)/admin/settings/dashboard/page.tsx' --file 'app/components/LanguageInitializer.tsx' --file 'lib/stores/i18n.ts'`
