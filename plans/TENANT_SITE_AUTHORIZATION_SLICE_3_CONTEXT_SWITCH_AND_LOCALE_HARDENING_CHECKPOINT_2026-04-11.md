# Tenant/Site Authorization Slice 3 Checkpoint

Date: 2026-04-11

## Scope
- Login and tenant-selector locale hardening
- Active language list deduplication
- Context persistence strengthened with cookies for SSR-aware flows
- Header tenant/site switch behavior aligned with selector-first context resolution

## Backend / Runtime Impact
- No backend contract change in this slice
- Next middleware now forwards tenant/site context from cookies into request headers

## Frontend Impact
- Shared modern locale pill switcher introduced
- Login page now uses locale-aware, preload-aware picker
- Tenant selector now uses active operational languages instead of ad hoc admin language payload
- Tenant/site context writes also persist cookies, not only localStorage
- Tenant switch from header routes back through site selection flow
- Site switch from header triggers route refresh so admin data reloads against the selected context

## Files
- `/workspace/node/aystek-core-next/components/i18n/LocalePillSwitcher.tsx`
- `/workspace/node/aystek-core-next/lib/stores/language.ts`
- `/workspace/node/aystek-core-next/lib/tenant.ts`
- `/workspace/node/aystek-core-next/middleware.ts`
- `/workspace/node/aystek-core-next/app/(auth)/login/page.tsx`
- `/workspace/node/aystek-core-next/app/tenant-selector/page.tsx`
- `/workspace/node/aystek-core-next/components/context/TenantSwitcher.tsx`
- `/workspace/node/aystek-core-next/components/context/SiteSwitcher.tsx`

## Expected Outcome
- Login page should not render raw `login.*` keys during normal runtime
- Tenant-selector language options should stay stable and reflect only active languages
- Language changes on tenant-selector should visibly update UI copy
- Default/admin-panel context should survive navigation more reliably because SSR now receives tenant/site cookies as request headers
