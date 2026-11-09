# TASK 058 - AYSTEK Haber Public Site Integration

## Scope

AYSTEK Haber was integrated as a Core-backed public news site under tenant
`aystek` and store/domain `aystekhaber.com`. The implementation avoids parallel
CMS logic by using the existing AYSTEK Core posts, categories, menus, authors,
translations, view profiles and post placement ownership where those domains
already existed.

## Backend

- Added `aystekhaber.com` store/site bootstrap under the existing store seed.
- Added `SeedAysTekHaberCMS` for page, branding, categories, navigation and 70
  localized launch posts.
- Added explicit post placement seed entries for `homepage_hero`, `breaking`,
  `sidebar` and `trending` using the existing post placement module.
- Added `SeedAysTekHaberSports` for an admin-managed Super League snapshot.
- Added tenant/site scoped `sports` module with model, DTO, repository, service,
  admin handler, public handler and routes.
- Added sports indexes, permissions, menu seed, translations and view profile
  visibility.
- Updated `cmd/seed` to load `.env` consistently with `cmd/api`.
- Added five idempotent editorial author profiles using the existing authors
  collection. Standalone editorial profiles can now be managed without a
  synthetic login account; linked-user authors remain supported.
- Hardened public author policy: list requests require `published` and active,
  and detail requests reject a record if either condition is false.
- Separated sports administrative and public reads. Admin may retrieve an
  inactive snapshot to reactivate it; public traffic can only retrieve active
  snapshots.
- Added sports domain invariants for unique rank/team, non-negative values,
  consistent played totals, valid match status and normalized team/form values.

## Frontend - Public Site

- Created `aystekhaber-com` as a standalone Next.js public site backed by Core.
- Added server-side Core clients for posts, placements, nav menus, pages,
  authors and sports.
- Refactored the visual chrome into `NewsChrome` to match the Figma header,
  topbar and breaking ticker structure.
- Added `SportsSnapshotPanel` for latest matches and league standings.
- Completed the template sections: primary/sponsor ads, category grids, dark
  video block, city-tabbed local news, author cards, mobile app callout and
  responsive footer/newsletter form.
- Removed the duplicate author content array from page configuration. The page
  keeps only an ordered `authorSlugs` selector while author card data comes from
  the existing Core author API.
- Added `/yazar/[slug]`, `/rss.xml`, `/sitemap.xml`, `/robots.txt` and the web
  manifest while retaining localized `/haber/[slug]` rendering.
- Added a same-origin newsletter proxy so the browser does not need direct Core
  access or duplicated tenant resolution logic.
- Added `sharp` and a standalone start script that prepares `.next/static` and
  `public` inside `.next/standalone` for local production runs.

## Frontend - Admin Panel

- Added `/admin/sports` page using the required gradient hero and rounded-2xl
  card structure.
- Added `lib/api/sports.ts` admin API client.
- Added post placement event documentation to the posts list page with
  `EventKeysInfo`.
- Replaced the sports breadcrumb literal with the seeded `sports.breadcrumb`
  translation and exposed Turkish/English league names in the editor.

## Operations

- Added registry and local-source public-site installer profiles.
- Updated installer public-site update sequencing so public-site deployments
  rebuild/pull and recreate the container before validation.
- Updated installer and central server setup documentation for AYSTEK Haber
  install, update, validation and controlled rollback.
- Corrected public-site TLS ownership. Domain/SAN resolution now uses
  `public_site_domain` plus normalized `public_site_domain_aliases` rather than
  the unrelated Core API/panel domains.
- Public-site installation writes an HTTP vhost until the certificate exists,
  then rewrites the correct public vhost as HTTPS, validates Nginx and reloads.
- Let's Encrypt now fails early with a clear error when `tls_email` is missing.
- Added focused config tests covering public-site alias de-duplication, primary
  certificate lineage and preservation of admin API/panel domain behavior.

## Validation

- `go test ./internal/modules/sports/... ./internal/seeds -run 'TestAysTekHaber|TestUpsert|TestPublicGetReturnsSnapshotInDataEnvelope' -count=1`
- `go build ./cmd/api`
- `npm run check-env && npm run build` in `node/aystekhaber-com`
- `npm run build` in `node/aystek-core-next`
- `go test ./internal/config ./internal/tls ./internal/installer ./internal/deploy`
  in `go/aystek-core-installer`
- `npx next lint --file 'app/(dashboard)/admin/sports/page.tsx' --file 'lib/api/sports.ts'`
- Live Core API checks for posts, nav menus, placements and sports snapshot.
- Playwright desktop and mobile checks for styled render, ticker, sports blocks,
  no horizontal overflow, no console errors and no failed requests.

## Notes

- `aystek-core-next` has unrelated existing TypeScript/lint debt outside the
  sports files. The new sports admin files and edited posts file have no editor
  diagnostics.
- Docker CLI is not available in this dev container, so full installer container
  deploy was documented and structurally validated, while runtime verification
  used the local Core API and standalone Next server.

## Runtime Ownership and Data Flow

The browser talks to the Next.js public site. Server components attach
`X-Tenant-ID=aystek` and `X-Store-Domain=aystekhaber.com` to Core requests. Core
resolves the store/site context before repository access. Posts, placements,
navigation, page configuration, authors and sports remain separate bounded
contexts; `NewsPortal` composes their read models without persisting a second
copy.

The home page configuration owns presentation text, external links, city tabs,
ad labels, app store links, footer links and the ordered author slug selector.
It does not own article, author or league business entities. This boundary lets
editors change a writer once in Authors and have both the home card and writer
profile reflect the change.

## Security and Isolation Decisions

- Public API clients are server-only; `API_BASE_URL` is not shipped to browsers.
- Newsletter submission is proxied same-origin and Core remains responsible for
  persistence and request policy.
- Public authors and sports have explicit active-state filtering.
- Sports repository queries require both tenant and site context.
- The AYSTEK Haber view profile is allowlist-based. Marketing/ads/newsletter are
  visible, while orders, payments, telematics and system operations remain
  hidden from the editorial profile.
- Seed execution is additive/idempotent; production instructions explicitly
  prohibit full-reset mode.

## Operational Acceptance Matrix

| Area | Acceptance criterion |
|---|---|
| Tenant/store | Requests resolve `aystek` + `aystekhaber.com` only |
| CMS | 70 launch stories and eight categories are idempotently present |
| Placements | Hero 1, breaking 5, sidebar 3, trending 5 |
| Authors | Five selected slugs resolve to active published author entities |
| Sports | Public returns active snapshot; admin can edit inactive snapshot |
| Responsive UI | Template sections render without horizontal overflow |
| Localization | Turkish and English content/labels/routes resolve |
| Newsletter | Same-origin endpoint forwards to Core and reports localized state |
| SEO | robots, sitemap, RSS and metadata routes are available |
| TLS | Primary + aliases share SAN cert; temporary vhost is HTTP-only |
| Update | Public-site image rebuild/pull precedes health validation |
| Rollback | Immutable known-good image tag can be reapplied without data reset |

## Follow-up Boundary

Live market/weather/provider feeds, actual advertising creatives and mobile app
store URLs remain operator-managed content or future provider adapters. Their
presentation contracts are dynamic today; no vendor-specific integration was
invented as part of this task.
