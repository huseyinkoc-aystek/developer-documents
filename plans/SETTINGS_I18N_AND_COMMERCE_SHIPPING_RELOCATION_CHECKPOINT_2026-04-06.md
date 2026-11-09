# Settings I18n And Commerce Shipping Relocation Checkpoint

Date: 2026-04-06

## Root Cause

- Admin translations were intermittently rendering raw keys because admin pages could be hydrated with a `publicOnly` translation dictionary.
- `shipping.selection_weights` UI lived in the generic system settings screen even though it belongs to commerce/shipping settings.
- The `settings.shippingWeights.*` keys were not present in the seed set, so that block had no translation source at all.

## Changes

- Fixed admin-side language initialization so `/admin` routes request full translation dictionaries.
- Added missing shipping weights translation keys to `translation_seed.go`.
- Added new commerce menu entry:
  - `menu.commerce.settings.shipping_settings`
- Moved shipping weights management UI from:
  - `/admin/settings`
  to:
  - `/admin/shipping/settings`
- Reduced settings-page lint debt by removing remaining explicit `any` usage in the touched files and stabilizing hook dependencies.

## Files

- `/workspace/node/aystek-core-next/app/components/LanguageInitializer.tsx`
- `/workspace/node/aystek-core-next/lib/stores/i18n.ts`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/page.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/settings/dashboard/page.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/shipping/settings/page.tsx`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`
- `/workspace/go/aystek-core/internal/seeds/menu_seed.go`

## Verification

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/settings/page.tsx' --file 'app/(dashboard)/admin/settings/dashboard/page.tsx' --file 'app/(dashboard)/admin/shipping/settings/page.tsx' --file 'app/components/LanguageInitializer.tsx' --file 'lib/stores/i18n.ts' --file 'lib/api/settings.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/settings/... ./internal/seeds ./cmd/api`
