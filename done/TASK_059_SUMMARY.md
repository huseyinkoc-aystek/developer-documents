# TASK 059 Summary - Site-Scoped Active Languages

Each public site can now choose which of the tenant's languages it offers,
managed from the admin panel. A site without its own selection keeps using the
tenant-wide active languages.

Completed items:

- `site_language_configs` collection + `SiteLanguageConfig` model (`Apply` drops
  languages no longer active, keeps order, guarantees a default)
- Repository (`tenantfilter.WithTenantAndSite`), service validation, admin
  endpoints `GET/PUT/DELETE /admin/languages/site`
- `/public/languages/active` and `/default` are site-aware
- `LanguageMiddleware` cache is per tenant+site (was one global cache)
- Admin: new `/admin/languages/site` page (menu key `site-languages`, own view
  profile handling) + banner and "Bu Sitede" column on `/admin/languages`
- View profiles: key registered in `allMenuKeys`, haber editorial and Gürsel
  araç takip visible sets, Gürsel menu list
- 28-language translation keys for all new UI text
- aystekhaber-com: language resolution, switcher, sitemap, newsletter and
  Pages-based label fallback (en -> tr) driven by Core instead of hard-coded
  `tr`/`en`
- aysteknoloji-com: sitemap limited to the site's enabled languages
- Standards: README_FOR_AI_ASSISTANTS v1.2.0 (RULE #1.1 site-scoped languages,
  RULE #2.2 view-profile sync, checklist items), canonical architecture 4.2.1,
  operator guides (`languages.md`, `aystek_haber.md`)
- Pre-existing red checks fixed: `go vet ./...` (8 packages of stale fakes),
  stale seed-count test, `MaskSensitiveURL`, setter-audit allowlist, lint errors
  in the touched admin files

Key decisions:

- Override lives next to the language catalog (same module) instead of
  extending the wide settings `Service` interface.
- Subset-of-catalog model: names/localizations stay defined once per tenant.
- No override = previous behavior, nothing changes until an admin opts in.
- Site languages got their own page/menu key (not a card in Tenant & Site) so
  view profiles can govern them.

Validation: unit tests, live two-site scenario, and a real-browser pass over the
admin screens (see the technical reference).

Second pass (tenant/context audits, duplicates, docs):

- All `tenantfilter` audits green: ~25 `context.Background()` uses moved to
  audited `contextkeys.System/ForTenantSite/Detached`; the duplicate `api_sharing`
  operation/scope repository, service and handler implementation (dead code)
  removed; global catalogs and platform-level `system_ops` allowlisted with
  documented rationale
- OpenAPI generator fixed (module path, nil-guarded route blocks, `//`, probe headers,
  deterministic output) and `api-docs/swagger.yaml` regenerated (803 -> 1551 paths)
- 29 new UI keys translated into all 28 languages (overlay file + tests)
- README v1.2.0 additions: background-context standard, allowlist rule, "duplicate
  implementations are a defect" rule

Third pass (leftovers):

- Tenant Admin role no longer receives `system:*` / `system_ops:*` (measured 200 -> 403 for a
  customer tenant admin; platform users unaffected) - README RULE #2.1
- `PublicSiteContextMiddleware` reviewed: load-bearing (admin bootstrap, payment links), documented,
  intentionally unchanged
- Production builds pass for aystekhaber-com and aysteknoloji-com
