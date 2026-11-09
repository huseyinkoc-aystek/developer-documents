# Auth, I18n, Tenant Selector And Public Repo Hygiene Checkpoint

Date: 2026-04-11

## Closed In This Slice

- `tenant-selector` site list parsing now reads the real `/me/bootstrap-sites` contract.
- auth/login flow now sends single-site users directly to dashboard instead of forcing a redundant site step.
- logout-time language loading no longer escalates a 401/403 into an unhandled runtime failure.
- `tr` is the default i18n baseline in Next admin runtime and the focused i18n hooks/API files were lint-cleaned.
- `aysteknoloji-com` now has a committed ESLint baseline and its current lint issues in blog pages were fixed.

## Important Technical Notes

- tenant/site selector visibility bugs were caused by frontend payload parsing, not by missing seeded sites.
- the broader “which user can access which tenant and which site” requirement is still a dedicated authorization feature slice, not a small bugfix.
- public repo runtime port was left on `9090` because that is still the canonical local backend port in the shared workspace; the temporary `9091` usage was only for isolated smoke validation.

## Files Touched

- `/workspace/node/aystek-core-next/lib/api/me.ts`
- `/workspace/node/aystek-core-next/lib/api/languages.ts`
- `/workspace/node/aystek-core-next/hooks/useI18n.ts`
- `/workspace/node/aystek-core-next/hooks/useLanguages.ts`
- `/workspace/node/aystek-core-next/contexts/SiteContext.tsx`
- `/workspace/node/aystek-core-next/app/(auth)/login/page.tsx`
- `/workspace/node/aysteknoloji-com/.eslintrc.json`
- `/workspace/node/aysteknoloji-com/components/pages/BlogDetailPageClient.tsx`
- `/workspace/node/aysteknoloji-com/components/pages/BlogPageClient.tsx`

## Validation

- focused Next lint for auth/i18n/site-context files
- `CI=1 npm run lint` in `/workspace/node/aysteknoloji-com`

## Next Correct Step

- run browser-level verification again for:
  - `/login`
  - `/tenant-selector?step=tenant`
  - `/tenant-selector?step=site`
- then decide whether to open the dedicated tenant-site authorization execution slice for per-user tenant/site access control
