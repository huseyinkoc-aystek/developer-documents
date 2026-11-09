# AYSTEKNolojileri Public Site Spot-Check — 2026-04-11

## Verified
- Public site repo points to aystek-core with:
  - `API_BASE_URL=http://localhost:9090`
  - `AYSTEK_STORE_DOMAIN=aysteknolojileri.com`
  - `AYSTEK_TENANT_ID=aystek`
- Backend runtime verification on `:9091` returned populated public CMS data for:
  - nav menus
  - published pages
  - posts
  - sliders
- Multi-tenant selector/runtime chain now confirms:
  - `default/admin-panel`
  - `default/default`
  - `aystek/aystek`

## Gaps
- `aysteknoloji-com` still lacks a committed ESLint configuration; `npm run lint` opens the interactive Next.js setup prompt.
- `.env.local` still targets port `9090`; local verification in this slice used backend on `9091`.

## Conclusion
- Backend CMS/content contract for `aysteknolojileri.com` is alive.
- Remaining work is frontend repo hygiene and full browser-level end-to-end verification against the corrected runtime.
