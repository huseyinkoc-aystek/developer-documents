# Multi-Tenant Login, Selector, and Runtime Checkpoint — 2026-04-11

## Fixed
- `me/bootstrap-sites` frontend parsing now reads canonical API payloads correctly, so seeded sites for `default` and `aystek` render instead of false empty-state.
- Auth layout preload keeps using `admin-panel` store context and no longer carries the local `any` typing debt in its payload helpers.
- Production runtime now starts successfully on `:9091`; login, `/me/context`, `/me/select-tenant`, and `/me/bootstrap-sites` respond correctly.

## Verified Runtime
- `default` tenant returns:
  - `admin-panel` -> `core.aysteknolojileri.com`
  - `default` -> `shop.local`
- `aystek` tenant returns:
  - `aystek` -> `aysteknolojileri.com`
- Public CMS spot-check for `aysteknolojileri.com` returns populated:
  - nav menus
  - published pages
  - posts
  - sliders

## Remaining Runtime Risk
- Browser hard-refresh is still needed on the active local session to flush old bundles that were showing stale login/selector behavior.
- `login/page.tsx` still carries the existing `no-img-element` lint warning; functional issue değil.
