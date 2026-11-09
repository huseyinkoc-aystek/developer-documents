# TASK 059 - Site-Scoped Active Languages

Technical reference. Owner module: `aystek-core/internal/modules/languages`.

## Problem

Languages were managed per **tenant** only (`languages` collection, unique
`tenant_id + code`). Every public site of a tenant therefore exposed the same
active languages. `aysteknoloji.com` and `aystekhaber.com` (both tenant
`aystek`) could not differ, e.g. 9 vs 5 languages. `aystekhaber-com` also had
`tr`/`en` hard-coded in every page and in its language switcher.

## Design

The tenant catalog stays the single source of truth (names, native names,
localized names, tenant-wide `is_active`/`enabled`/`is_default`). A site may
add an optional **override** that narrows the catalog:

| Site state                    | Languages the site exposes                          |
|-------------------------------|-----------------------------------------------------|
| No override record            | Tenant-wide active languages (previous behavior)    |
| Override record               | Only `codes`, in that order, `default_code` default |
| Override, nothing usable left | Falls back to tenant-wide languages                 |

Rules enforced in `service.UpdateSiteLanguages`:

- at least one language, no duplicates, codes trimmed/lower-cased
- every code must be an **active + enabled tenant language**
- `default_code` must be one of `codes` (empty = first code)

If a language is later deactivated tenant-wide it silently drops out of the
site's effective list (`SiteLanguageConfig.Apply`); the default moves to the
first remaining language.

## Data model

Collection `site_language_configs` (unique index `tenant_id + site_id`):

```go
type SiteLanguageConfig struct {
    ID          primitive.ObjectID
    TenantID    string   // bson: tenant_id
    SiteID      string   // bson: site_id
    Codes       []string // ordered
    DefaultCode string
    CreatedAt, UpdatedAt time.Time
    UpdatedBy   string
}
```

All repository operations go through `tenantfilter.WithTenantAndSite`.
Site is taken from the request context (`X-Site-ID` in admin, store domain /
`X-Store-Domain` on public sites).

## API

| Method | Path                    | Purpose                                                  |
|--------|-------------------------|----------------------------------------------------------|
| GET    | `/admin/languages/site` | Effective languages of the selected site + `inherited`   |
| PUT    | `/admin/languages/site` | `{ "codes": ["tr","en"], "default_code": "tr" }`         |
| DELETE | `/admin/languages/site` | Remove override, follow tenant-wide languages again      |
| GET    | `/public/languages/active`  | Now site-aware (override or tenant-wide)             |
| GET    | `/public/languages/default` | Now returns the site default                         |

`GET /public/languages` (paginated list, used by the admin panel) intentionally
stays tenant-wide.

Requests without a site context (site-exempt `/languages*` routes) ignore
overrides and behave exactly as before.

`LanguageMiddleware` now caches the public language list **per tenant+site**
(30 s). Previously one process-wide cache was shared by every tenant/site,
which would have leaked one site's list to another.

## Admin panel (aystek-core-next)

Two screens, one menu entry each (no duplicated management surface):

- `/admin/languages` (existing) - tenant catalog. Now also shows a banner for the
  selected site ("uses general settings" / "custom selection, N languages"), a
  **Bu Sitede** column per language and a "Site Dilleri" shortcut in the hero.
- `/admin/languages/site` (new, menu key `site-languages`, Settings group) -
  the selected site's selection: checkbox per active tenant language, default
  radio, up/down ordering, "Genel ayara dön" (confirmed with `NoticeModal`).
  Built to the admin visual standard (gradient hero, `rounded-2xl` cards,
  in-hero stat tiles). State lives in `hooks/useSiteLanguages.ts`.

An earlier iteration put this as a card inside Tenant & Site Settings; it was
moved to its own page so it has its own menu key/permission surface for view
profiles, and the card + component were deleted.

Translation keys: `settings.siteLanguages.*`, `languages.site.*`,
`languages.table.headers.site`, `menu.settings.site_languages` (28 languages;
9 real translations tr/en/de/fr/ar/fa/zh/ja/pt, the rest English via
`allLangsMulti`, same as the existing seeds).

## View profiles

`site-languages` was added to `allMenuKeys` (system default), to the
AYSTEK Haber editorial visible set (+ allowlist test) and to the Gürsel araç
takip visible set; `gursel_turizm_seed.go` has its own menu list and got the
entry too. Profiles built with `hiddenMenuItemsExcept` (AYSFİLO, AYPASS, Veli,
ComfortPay) hide it by default (least privilege). Settings menu orders after
`languages` were shifted by one. Verified in the seeded DB: not hidden in the
haber and Gürsel araç takip profiles.

## Public sites

- **aysteknoloji-com**: language picker and provider already consumed
  `/public/languages/active`, so they follow the site list automatically.
  `app/sitemap.ts` no longer force-adds `tr/en/es`; it lists only enabled
  languages (legacy trio only when the backend returns none).
- **aystekhaber-com**: was hard-coded `tr`/`en`. Now
  - `lib/api/languages.ts`: `fetchSiteLanguages()` / `resolveSiteLanguage()`
    (requested language if enabled, otherwise site default; TR/EN only when
    the backend is unreachable)
  - every page/metadata resolver, the sitemap and the newsletter route use it
  - `LanguageSwitcher` + `useSiteLanguages` render one link per enabled
    language (desktop and mobile chrome) via `/api/languages`
  - non-`tr` languages use English inline UI strings (was: only `en`)
  - portal labels / page copy come from the home page's per-language
    localization in the Pages module (admin-managed). A language without one
    falls back English -> Turkish (`lib/api/pages.ts`: one fetch of the
    canonical home document + `pickLocalization`; static pages try the
    requested language slug, then `en`, then `tr`). Before this, e.g. `de`
    rendered the Turkish defaults.

## Verification

- `go build ./...`, `go vet ./...` clean; `go test` for languages (model
  `Apply`, service resolution/validation/reset), middleware, seeds,
  httpclient.
- Live (local Mongo seeded with `cmd/seed`, backend :9090, admin :3000,
  haber :3010), tenant `aystek`:
  - both sites start with 11 languages (inherited)
  - `PUT` aystekhaber -> `tr,en,de,fr,ar`: aystekhaber public = 5 languages,
    aysteknolojileri still 11; `?lang=ja` -> `Content-Language: tr` on haber,
    `ja` on aysteknolojileri; default change; `DELETE` restores 11 and empties
    `site_language_configs`; invalid payloads -> 400
  - Browser (Playwright/Chromium against the real admin): new page renders in
    the standard, unchecking languages + save -> `PUT 200` + success modal,
    languages list shows banner + "Sunuluyor/Sunulmuyor" per row and all
    columns fit at 1360px, reset -> confirm modal -> `DELETE 200`
  - haber `/?lang=tr|en|de` -> Turkish / English / English labels;
    `/api/languages` returns only the site's languages
- `tsc --noEmit` and `next lint` clean for the touched files in
  aystek-core-next, aystekhaber-com and aysteknoloji-com.

## Cleanups made in the same work item

- Removed the `SiteLanguagesCard` (superseded by the page) and the unused
  `TenantSitePublicSocialLinks` import that failed lint in the tenant-site page.
- `languages/page.tsx`: removed `any` usages / hook-dependency warnings, fixed
  the delete-confirm placeholder mismatch (`{{name}}` vs `code`), typed RTL
  hint, bounded the localized-names column so the table fits.
- Fixed pre-existing red checks: stale test fakes (`go vet ./...` was failing in
  8 packages), stale "exactly 38/70 stories" assertions in
  `aystek_haber_seed_test.go`, `MaskSensitiveURL` userinfo masking,
  `public_site_context.go` added (with rationale) to the tenant/site setter
  audits, `cmd/updater` vet error (redundant newline), site-audit false positives.

## Tenant/context audit work (second pass)

All `internal/platform/tenantfilter` audits are green.

- **Background contexts:** ~25 literal `context.Background()` uses in
  repository/service layers were replaced with three audited helpers in
  `internal/platform/contextkeys` (`System`, `ForTenantSite`, `Detached`) with unit
  tests. Where a caller context existed it is now threaded through
  (`chain_runner` edge conditions, `system_ops` update/rollback/backup jobs and
  `execUpgradeScript`); detached goroutines keep the caller's tenant/site
  (`stock/reporting`, `api_sharing` consent cascade, PA media proxy); scheduled
  fleet-tally runs use `ForTenantSite`; constructor index creation and
  tenant-provisioning hooks use `System` (unchanged bare-root behavior).
- **`api_sharing` duplicate implementation removed:** `operation_mongo.go`,
  `scope_mongo.go`, `operation_service.go`, `scope_service.go`,
  `operation_handler.go`, `scope_handler.go`, the `indexes/` sub-package and the
  two repository test files (they tested private in-file mocks) were dead code -
  production wires only `MongoRepository`/`service.go`, indexes come from
  `api_sharing/indexes.go`. The wired operation/scope catalog methods were moved
  unchanged into `operation_catalog_mongo.go` / `scope_catalog_mongo.go`.
- **Allowlist (documented in `REPOSITORY_ALLOWLIST_RATIONALE.md`):** the two
  catalog files (models have no tenant field) and `system_ops/repository/mongo.go`
  (platform lifecycle records). `TestRepositoryTenantEnforcementEvidence` now
  shares the same allowlist instead of a private hard-coded exemption.

## Known limits / follow-ups

- `tools/openapi_gen.py` (`make openapi`) was broken and is fixed; `api-docs/swagger.yaml`
  is regenerated from it (see "Third pass").
- The 29 new UI keys are translated into all 28 languages
  (`seeds/site_language_translations_seed.go` overlays the 19 languages that
  `allLangsMulti` would fill with English; `az-cyrl` is a transliteration of the
  Azerbaijani text). A native-speaker review of the less common languages is
  still advisable.

## Third pass (leftovers)

- **OpenAPI generator fixed and document regenerated.** Root causes: (1) the
  import path was hard-coded as `github.com/aystek/aystek_core` while the module is
  `aystek-core`, so no route group prefix ever resolved (module routes came out as
  `/ads` instead of `/admin/ads`, `//device/...`); now read from `go.mod`;
  (2) `RegisterRoutes(adminGroup, adminH, nil)` calls emitted the `if public != nil`
  blocks too (bogus `/public/ads/{id}` DELETE, `/admin/languages/active`, ...);
  the generator now skips blocks guarded by a parameter the call passes as `nil`;
  (3) `r.Group("")` produced `//`; (4) probes (`/healthz`, `/readyz`, `/metrics`)
  now carry the optional context headers the contract test requires, and the
  contract matrix treats `/metrics` like the other probes; (5) output is
  `sort_keys` deterministic. Result: 803 -> 1551 paths, 2111 operations
  (`swagger.count` = operations), 112 bogus paths gone, no duplicate `//` paths,
  `make security-gates` and `internal/platform/openapi` green. Common operations keep
  identical operationId/tags/security/summary.
- **Platform-only permissions (security).** A customer tenant admin
  (`type=tenant`, non-`@aystek.com`) could read `/admin/system/{version,updates,backups}`
  because the seed gave the Admin role all 32 `system:*`/`system_ops:*` permissions
  (measured `200`). They are now excluded via `platformOnlyPermissionModules` in
  `role_seed.go` (README RULE #2.1 mechanism); measured `403` afterwards, tenant-scoped
  routes still `200`, AYSTEK platform users unaffected (RBAC bypass in
  `PermissionMiddleware`). Re-seeding updates existing admin roles.
- **`PublicSiteContextMiddleware` reviewed, intentionally unchanged.** An unresolved
  host on `/public/*` falls back to tenant `aystek` (measured: `unknown.example.com`
  returns aystek data). It is load-bearing: the admin panel calls
  `/public/translations|languages|settings` from an API host that is not a
  registered store domain, and the public payment-link pages have no tenant context.
  Removing or narrowing it would break admin login translations and payment links;
  the safe follow-up is an explicit per-route public-tenant policy, not a blanket
  change.
- **Production builds:** `next build` succeeds for aystekhaber-com and
  aysteknoloji-com; the admin build was OOM-killed by the sandbox (VS Code processes),
  the admin panel is verified by `tsc`, `next lint` and a real-browser pass.
- **`check-i18n` implemented** (`aystek-core-next/scripts/check-i18n.mjs`, ratchet with
  `scripts/i18n-baseline.json`; the old npm script pointed to a missing ts-node file):
  fails only for NEW `t("key")` usages no seed defines; 2259 legacy keys are baselined.
- **Translation coverage gate:** `buildTranslationEntries()` was extracted from
  `SeedTranslations` (behavior unchanged) so `translation_coverage_test.go` can check keys
  without a database: `fullyTranslatedPrefixes` (ratchet) + `TestTranslationCoverageReport`
  diagnostic. Measured legacy debt: `languages.*` 62 keys / 1200 gaps; 8307 keys / 106,230
  gaps overall.
- **Remaining legacy translation debt is planned, not ignored:**
  `I18N_TRANSLATION_DEBT_WORKPLAN_TR.md` (WP-1 `languages.*`, WP-2 seedless keys, WP-3
  fallback keys) with a ready prompt for other AI tools for WP-1.
