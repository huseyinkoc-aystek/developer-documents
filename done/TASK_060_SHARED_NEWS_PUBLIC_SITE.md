# TASK 060 - Shared News Public-Site Runtime and Küresel Hat Provisioning

Technical reference for the reusable news/magazine renderer, its Core site
identity, and its installer deployment contract.

## Objective

The Figma-derived `aystekhaber-com` application must remain one local source
project for every domain that uses the same design. `aystekhaber.com` and
`kureselhat.com` are separate Core sites and separate server deployments.

The implementation must satisfy all of these constraints:

- one Core backend authority;
- one reusable frontend codebase and Docker image for this design;
- one tenant/site identity per domain;
- no content, brand, SEO or cache leakage between sites;
- no duplicate local project for a same-design domain;
- separate server source/install/config/log directories per domain;
- separate runtime environment, port, container and Nginx/TLS vhost;
- dynamic/admin-managed content through existing Core modules;
- no demo content copied to a newly provisioned production site;
- idempotent targeted provisioning that preserves administrator changes.

## Canonical topology

```text
Local development
  /workspace/node/aystekhaber-com
    process 1 -> port 3001 -> aystekhaber.com identity
    process 2 -> port 3002 -> kureselhat.com identity

Server
  /opt/aystek-public/aystekhaber-com/source
    -> /opt/aystek-public/aystekhaber-com
    -> /etc/aystek-public/aystekhaber-com
    -> /var/log/aystek-public/aystekhaber-com
    -> 127.0.0.1:3001

  /opt/aystek-public/kureselhat-com/source
    -> /opt/aystek-public/kureselhat-com
    -> /etc/aystek-public/kureselhat-com
    -> /var/log/aystek-public/kureselhat-com
    -> 127.0.0.1:3002

Both deployments
  -> one immutable aystekhaber-com image/version
  -> one aystek-core service
  -> distinct store/site records and site-scoped content
```

A different visual design is a different project in local development and on
the server. This task does not turn unrelated designs into one conditional UI.

## Ownership boundaries

### AYSTEK Core

Core owns tenant/site identity, language availability, pages, navigation,
settings, posts, authors, galleries, taxonomy, post placements, sport snapshots,
newsletter data and access isolation.

### Shared public renderer

The renderer owns presentation and typed adaptation of public Core APIs. It
does not infer a site from an untrusted Host and does not persist business data.

### Installer

The installer owns filesystem paths, environment rendering, image selection,
compose identity, loopback port, Nginx vhost, certificate lifecycle and state.

## Core site provisioning

The site definition lives in `internal/seeds/news_sites_seed.go`:

```text
tenant_id: aystek
site code: kureselhat
domain: kureselhat.com
alias: www.kureselhat.com
type: public_news_portal
renderer: aystekhaber-com
public_site_ready: false
```

The tenant is the existing stable AYSTEK tenant, not a newly created tenant
record.

Targeted execution:

```bash
cd /workspace/go/aystek-core
go run ./cmd/seed --news-site=kureselhat
```

The command loads normal Mongo configuration, requires the active `aystek`
tenant, uses a bounded system context and exits without invoking the full seed
runner or reset prompt.

### Created defaults

- one active store/site;
- canonical domain plus `www` alias metadata;
- one published homepage document with TR/EN localizations;
- nine public navigation documents;
- RSS-only social-links setting;
- empty weather/market/live ticker setting;
- one site-scoped editorial view profile;
- references to the existing tenant-scoped news taxonomy.

### Deliberately not created

- posts;
- authors;
- galleries;
- post placements;
- sport snapshots;
- fake market/weather values;
- fake social profiles;
- fake app store links;
- corporate or legal copy presented as final business information.

### Idempotency policy

Production-style news-site defaults use insert-only semantics.

- Existing IDs remain stable.
- Existing administrator edits remain stable.
- Existing `updated_at` values are not touched by a repeat run.
- Soft-deleted defaults are not resurrected or duplicated.
- An inactive, deleted or domain-changed store requires operator review.
- A canonical-domain or alias collision aborts before writes.
- A site cannot be moved across tenants by seed execution.

AYSTEK Haber's legacy demo seed retains its existing update behavior. The new
targeted Küresel Hat command never runs that demo catalog.

## Shared seed refactoring

Existing AYSTEK Haber helpers were separated into reusable template functions:

- `seedNewsTaxonomy`;
- `newsHomePageDefaults`;
- `seedNewsNavigation`;
- `seedNewsPublicSettings`;
- `seedNewsViewProfile`.

The site code prefixes navigation keys. Homepage brand names are parameters.
Demo author selectors and placeholder app-store links are opt-in and remain
limited to the legacy demo. The new site receives empty selectors/stores.

## Public Core isolation

`SiteIsolation` now distinguishes unscoped public bootstrap access from an
explicit public-site selection.

An explicit selection is present when the request carries one of:

- `X-Site-ID`;
- `X-Store-Code`;
- `X-Store-Domain`.

When explicit, the request must resolve to an active store and valid site ID.
If a domain is supplied, it must match the canonical domain or a registered
alias of the resolved store. A site ID or code cannot override an incompatible
domain.

The integration matrix covers settings, active languages, pages and posts with:

- canonical domain;
- registered alias;
- domain-only selection;
- unknown domain;
- inactive site;
- tenant mismatch;
- site-ID override attempt.

Unknown explicit public domains now fail closed instead of receiving tenant or
global settings.

## Site-only branding API

The public settings response retains legacy branding fields for compatibility
and adds `site_branding`.

`site_branding` is resolved only from keys scoped to the current site ID. It
does not inherit a tenant/global logo. This prevents a new news publication
without a logo from showing the AYSTEK corporate logo or another publication's
logo.

The frontend checks that `site_branding.site_id` equals the homepage site ID
before applying logo or favicon values.

## Runtime identity contract

The public renderer requires:

| Variable | Responsibility |
|---|---|
| `SITE_URL` | canonical HTTP(S) origin at runtime |
| `SITE_DOMAIN_ALIASES` | accepted additional Host values |
| `AYSTEK_STORE_DOMAIN` | canonical Core store-domain header |
| `AYSTEK_TENANT_ID` | Core tenant header |
| `API_BASE_URL` | server-only Core API base URL |
| `NEXT_DIST_DIR` | optional local concurrent build directory |

`NEXT_PUBLIC_SITE_URL` is not used as the news renderer's identity. Public
Next variables are build-time values and would bake the first site into a
shared image.

`getDeploymentSite()` validates required values and permits only a clean
HTTP(S) origin. User info, path, query, fragment and non-HTTP protocols fail.

## Host admission

Host is an admission check, never an identity selector.

Accepted:

- hostname from `SITE_URL`;
- `AYSTEK_STORE_DOMAIN`;
- explicit `SITE_DOMAIN_ALIASES`;
- loopback names used by health/installer checks.

Rejected:

- another news deployment's domain;
- an unknown domain;
- suffix attacks such as `kureselhat.com.attacker.test`;
- comma-separated forwarded-host injection;
- user-info style values;
- malformed host strings.

The middleware returns 421 for a host that belongs to another or unknown
deployment. Missing/malformed mandatory runtime identity returns 503.

## Core request adapter

Every server-side Core request:

- sets `X-Store-Domain` from deployment configuration;
- sets `X-Tenant-ID` from deployment configuration;
- removes caller-provided `X-Site-ID`;
- removes caller-provided `X-Store-Code`;
- overwrites caller-provided domain and tenant headers;
- uses `cache: no-store` by default;
- uses a ten-second timeout by default.

This prevents route code or an inbound request from changing site authority.

## Homepage identity

The existing site-scoped homepage document is the identity source. Its
localized title supplies the brand name. `site_id` links that identity to
site-only settings.

Localization lookup order is requested language, English, then Turkish. The
site's available/default language still comes from the site-aware language API.

An empty or stale `content.brand.name` can no longer erase the homepage title.
The typed adapter overwrites template brand name with the localized page title.

Only 404 permits the homepage slug fallback from `tr/anasayfa` to `en/home`.
Rate limiting, invalid deployment identity and Core outages remain operational
errors; they no longer masquerade as a Next.js not-found response.

## Brand rendering

`NewsBrandLogo` supports a site-specific asset or a text wordmark. When no
site-scoped logo exists, the publication name is rendered as text. A tenant or
other site's logo is never used as a fallback.

The header, mobile menu, footer and accessible labels all receive the same
runtime brand configuration.

Generic placeholder and favicon assets are used where a publication-specific
asset is absent. The AYSTEK logo asset remains available only when the AYSTEK
Haber site setting explicitly references it.

## Removal of fabricated runtime data

The shared renderer no longer displays made-up operational/product values.

- market values default to an empty list;
- weather defaults to absent;
- live stream defaults to absent;
- social links default to absent except seeded RSS;
- mobile app stores default to empty;
- the app-download panel is hidden without real links;
- author fallback wording is publication-neutral.

Administrators can populate supported values through existing settings/pages.

## SEO contract

`createSiteSeo(site)` owns metadata and structured-data creation.

The runtime site supplies:

- publisher name;
- canonical origin;
- OpenGraph site name;
- article and gallery URLs;
- JSON-LD publisher;
- optional publisher logo.

The following surfaces are site-aware:

- root layout metadata;
- homepage metadata;
- category landing metadata;
- post detail metadata;
- author list/detail metadata;
- gallery list/detail/photo metadata;
- static corporate/legal page metadata;
- canonical URLs;
- OpenGraph payloads;
- NewsArticle JSON-LD;
- gallery JSON-LD.

No metadata helper stores a process-global mutable site identity.

## Feed and discovery endpoints

Runtime site identity now drives:

- `/rss.xml`;
- `/sitemap.xml`;
- `/robots.txt`;
- `/manifest.webmanifest`.

The sitemap uses enabled site languages and site-owned content. Corporate page
slugs come from the homepage configuration and are verified against one
site-scoped published-page listing before inclusion. Posts, authors, navigation,
galleries and homepage localizations are each fetched as a shared multilingual
data set instead of once per language. A deployment never emits another
publication's origin.

Newsletter subscription source also uses the runtime domain rather than a
hard-coded AYSTEK Haber string.

## Build and image model

The Docker build does not receive site identity build arguments. It produces a
generic news-template image.

The runtime stage contains:

- Next standalone server;
- `.next/static` assets;
- public assets;
- non-root `nextjs` user.

`.dockerignore` excludes environment files, build artifacts and local logs.

For manual standalone execution, use `npm start`. Its wrapper copies `.next`
static and public assets into the standalone tree before starting `server.js`.
Starting a newly built raw standalone server without these copies produces an
unstyled page and is not a valid verification method.

## Local concurrent development

Example AYSTEK Haber process:

```bash
SITE_URL=http://localhost:3001 \
AYSTEK_STORE_DOMAIN=aystekhaber.com \
AYSTEK_TENANT_ID=aystek \
API_BASE_URL=http://localhost:9090 \
SITE_DOMAIN_ALIASES=www.aystekhaber.com \
NEXT_DIST_DIR=.next-aystekhaber \
npm run dev -- --port 3001
```

Example Küresel Hat process:

```bash
SITE_URL=http://localhost:3002 \
AYSTEK_STORE_DOMAIN=kureselhat.com \
AYSTEK_TENANT_ID=aystek \
API_BASE_URL=http://localhost:9090 \
SITE_DOMAIN_ALIASES=www.kureselhat.com \
NEXT_DIST_DIR=.next-kureselhat \
npm run dev -- --port 3002
```

The separate `NEXT_DIST_DIR` values prevent development build artifacts from
colliding. They are not production identity.

## Installer contract

The public-site environment template emits both legacy `NEXT_PUBLIC_SITE_URL`
for other templates and runtime `SITE_URL`/`SITE_DOMAIN_ALIASES` for the shared
news renderer.

Both the source template and embedded asset template were updated and tested.
The embedded copy is what the compiled installer reads.

Provided profiles:

- `docs/examples/aystekhaber-local-source.json`;
- `docs/examples/aystekhaber-registry.json`;
- `docs/examples/kureselhat-local-source.json`;
- `docs/examples/kureselhat-registry.json`.

The two local-source profiles intentionally share the image tag because the
design/source is the same. They intentionally differ in:

- domain and aliases;
- store domain;
- port;
- source path;
- install path;
- config path;
- log path;
- compose/container identity;
- installer state-file path.

The installer test renders both source and embedded env templates. It asserts
runtime identity values, shared image identity and non-overlapping deployment
state.

## Scheduler safety discovered during verification

Starting the entire Core binary for a read-only browser test also starts normal
background jobs. The personal-assistant agency scheduler immediately wrote
eight agency articles and six story threads.

The exact 14 test-created records were identified by IDs and creation window,
then deleted. The same time window was verified to contain zero remaining
agency articles and story threads. The delete is not recoverable from MongoDB;
source articles can be fetched again.

Core now supports:

```bash
PA_AGENCY_SCHEDULER_ENABLED=false
```

Default remains `true`, preserving production compatibility. A unit test
asserts explicit false opt-out and true default behavior. Browser verification
used the safe opt-out and confirmed no new agency/article/thread records.

## Rate-limit finding

Next.js server rendering fans one browser page into multiple Core public API
requests. Core's public limiter keys anonymous requests by network identity.
Depending on proxy topology, all requests from one renderer container can share
that identity.

During multi-page browser smoke testing, the default 60 requests/minute was
exhausted. The homepage identity request returned 429, which exposed the older
404-masking issue.

Actions:

- the frontend now treats 429 as an operational failure;
- the first broad browser verification used a temporary higher test limit;
- sitemap generation no longer probes each corporate page or reloads posts,
  authors and navigation for every enabled language;
- a focused production-build smoke test then ran sitemap, homepage and alias
  checks consecutively at the default 60 requests/minute without a Core 429;
- operations docs require traffic-based `PUBLIC_RATE_REQ` sizing;
- production monitoring must alert on 429 rate;
- a later capacity task may introduce a trusted renderer identity or further
  reduce normal page render-time API fan-out; this task does not weaken Core's
  default globally.

## Verification evidence

### Core targeted tests

```bash
NEWS_SITE_TEST_MONGO_URI=mongodb://mongodb:27017 \
go test ./internal/seeds -run 'Test(NewsSite|NewsHome|AysTekHaber)' -count=1

go test ./cmd/api ./cmd/seed \
  ./internal/platform/middleware \
  ./internal/modules/settings/service \
  ./internal/platform/router \
  ./internal/platform/tenantfilter \
  ./internal/modules/stores/repository \
  ./internal/modules/view_profiles/repository \
  ./internal/platform/db
```

Both commands passed.

### Core full verification

`go test ./...` and `go build ./cmd/api ./cmd/seed` were executed after targeted
tests. Both passed; the seed package completed in 289.142 seconds.

### Frontend verification

```bash
npm run check-env
npm run test:site
npm run lint
npm run build
```

All passed. Contract suite result: eight tests, eight passes.

### Browser verification

Chromium exercised the production standalone build for `kureselhat.com`.

Verified:

- homepage HTTP 200;
- desktop layout and dynamic Küresel Hat wordmark;
- mobile layout without horizontal overflow;
- mobile navigation open/close and labels;
- authors page HTTP 200;
- galleries page HTTP 200;
- agenda category page HTTP 200;
- OpenGraph site name equals Küresel Hat;
- canonical hostname equals `kureselhat.com`;
- no AYSTEK Haber text/logo leak;
- no fake app-store CTA;
- RSS status and site title/origin;
- sitemap status and own origin only;
- robots sitemap URL;
- manifest publication name;
- unknown Host returns 421;
- AYSTEK Haber Host against Küresel deployment returns 421;
- `www.kureselhat.com` alias succeeds;
- no browser page errors.

After sitemap fan-out reduction, a second production-standalone smoke test used
Core's unchanged default 60 requests/minute limit. Consecutive sitemap, homepage
and alias requests returned 200/200/200, the cross-domain Host returned 421,
the sitemap contained 100 URLs, and Core logged no 429 response.

Screenshots were visually inspected at desktop and mobile dimensions. The
publication is intentionally empty and displays “Yayın hazırlanıyor”.

### Installer verification

```bash
go test ./...
cmp -s templates/env/publicsite.env.tmpl \
  assets/templates/env/publicsite.env.tmpl
git diff --check
```

All passed. The repository-wide `scripts/check-env-contract.sh` still reports
an unrelated pre-existing drift between the panel source and embedded templates;
this task neither created nor changed that panel-template difference.

### Development database result

Latest local result:

```text
site_id: 6ab7fe6da38e0fee294e27bd
store: 1
homepage: 1
public navigation: 9
public settings: 2 (site-ID encoded keys)
editorial profile: 1
posts: 0
authors: 0
galleries: 0
```

The site ID is environment-specific and must never be copied to installer or
application configuration.

## Files by repository

### aystek-core

- `cmd/api/main.go`;
- `cmd/api/background_jobs_config_test.go`;
- `cmd/seed/main.go`;
- `internal/modules/settings/service/interface.go`;
- `internal/modules/settings/service/service.go`;
- `internal/modules/settings/service/service_test.go`;
- `internal/platform/middleware/site_isolation.go`;
- `internal/platform/middleware/public_site_isolation_test.go`;
- `internal/seeds/news_sites_seed.go`;
- `internal/seeds/news_sites_seed_test.go`;
- `internal/seeds/aystek_haber_seed.go`;
- `internal/seeds/view_profiles_seed.go`;
- `internal/seeds/runner.go`.

### aystekhaber-com

- runtime identity and middleware;
- typed API adapters;
- site identity composition;
- dynamic brand components;
- page metadata and SEO helpers;
- discovery/feed routes;
- environment examples;
- Docker/runtime packaging;
- contract tests;
- project README.

The Git diff is the authoritative detailed file list.

### aystek-core-installer

- public-site source and embedded env templates;
- AYSTEK Haber local-source path correction;
- Küresel Hat local-source and registry profiles;
- shared-source/separate-deployment test;
- installer README and user guide.

### aystek-core-dev-docs

- multidomain workplan;
- server installation guide;
- operations guide;
- deployment scenarios;
- this technical reference;
- executive summary.

## Production launch checklist

Code readiness does not equal publication readiness.

- [ ] Create real Küresel Hat posts, authors and galleries.
- [ ] Review placement selections.
- [ ] Create verified about, masthead, contact and advertising pages.
- [ ] Create privacy, cookie and terms documents with legal review.
- [ ] Upload site-specific logo and favicon.
- [ ] Configure only real social links.
- [ ] Configure analytics and verification tokens if approved.
- [ ] Review tenant-inherited languages and set a site override if needed.
- [ ] Replace placeholder operator email in installer profile.
- [ ] Select immutable Core and public-site image versions.
- [ ] Size Core public rate limits from load test evidence.
- [ ] Provision separate server source/install/config/log paths.
- [ ] Verify unique loopback port and container name.
- [ ] Verify DNS A/AAAA for apex and aliases.
- [ ] Verify Nginx default-host rejection.
- [ ] Obtain and verify TLS SAN coverage.
- [ ] Run installer validate.
- [ ] Check Core health and background workers.
- [ ] Check canonical/OpenGraph/JSON-LD on the public domain.
- [ ] Check robots, sitemap, RSS and manifest externally.
- [ ] Check unknown/cross-domain rejection externally.
- [ ] Check newsletter subscription source and consent text.
- [ ] Record release, backup and rollback image tags.
- [ ] Set `public_site_ready=true` only after product-owner launch approval.

## Rollback

Frontend rollback changes only the affected deployment's immutable image tag,
then runs installer update and validation. It must not change the sibling site's
tag implicitly.

Seed rollback is not a destructive reverse seed. The seed is idempotent and
does not delete administrator content. If data restoration is required, use a
reviewed MongoDB backup/restore procedure scoped to the affected site.

DNS/TLS rollback remains an operator action. Keep the previous vhost and image
configuration available until post-deployment verification is complete.

## Known remaining work

The implementation is complete for shared-code runtime operation. Remaining
work is operational/product data:

- real content and editorial approval;
- corporate and legal content;
- publication brand assets;
- production DNS and certificate;
- production load test and rate-limit sizing;
- external monitoring and rollback rehearsal.

The translation-debt workplan is independent and remains active. This task did
not add a new admin module, menu key, permission or user-visible translation
key, so it does not create another 28-language seed obligation.
