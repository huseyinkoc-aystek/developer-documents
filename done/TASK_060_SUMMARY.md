# TASK 060 - Executive Summary

## Outcome

`aystekhaber-com` is now a reusable news/magazine renderer for sites that share
the same design. AYSTEK Haber and Küresel Hat use one local codebase and one
build/image, while keeping Core data and server runtime completely separate.

Küresel Hat was provisioned under tenant `aystek` as the independent site
`kureselhat.com`, with `www.kureselhat.com` as its alias. The targeted seed adds
only safe site defaults and copies no AYSTEK Haber demo content.

## Architecture decision

- Local: one `/workspace/node/aystekhaber-com` checkout, separate processes and
  separate `NEXT_DIST_DIR` values.
- Server: separate source/install/config/log folders, env files, ports,
  containers, compose state and Nginx/TLS vhosts per domain.
- Core: one backend, distinct tenant/site context and site-scoped records.
- Image: one generic immutable image; site identity is runtime configuration.
- Different visual design: separate project locally and on the server.

## Core changes

- Added targeted `--news-site=kureselhat` provisioning.
- Added idempotent, insert-only site defaults.
- Reused existing pages, menus, settings, taxonomy and view-profile models.
- Preserved administrator edits and deleted defaults.
- Added canonical-domain and alias collision protection.
- Added fail-closed explicit public-domain isolation.
- Added site-only public branding without tenant-logo fallback.
- Added tests for identity, isolation, collisions and repeat execution.
- Added `PA_AGENCY_SCHEDULER_ENABLED=false` for safe browser/API test startup.

Latest development database result: one store, one homepage, nine public menus,
two settings, one editorial profile, and zero posts/authors/galleries.

## Frontend changes

- Removed build-baked site identity.
- Added mandatory runtime site contract.
- Added Host admission middleware.
- Prevented request headers from overriding deployment site/tenant.
- Made brand, metadata, canonical URLs, OpenGraph and JSON-LD site-aware.
- Made RSS, sitemap, robots, manifest and newsletter source site-aware.
- Restricted logos/favicons to explicit site branding.
- Removed fabricated weather, market, social, live and app-store fallbacks.
- Preserved site-scoped language behavior.
- Made Core 429/outage an operational error instead of a false 404.
- Added six runtime/identity/SEO contract tests.

## Installer changes

- Added runtime `SITE_URL` and `SITE_DOMAIN_ALIASES` to source and embedded env
  templates.
- Added Küresel Hat local-source and registry examples.
- Corrected AYSTEK Haber server source path to its own deployment directory.
- Asserted shared image plus separate paths, ports, state and container identity.
- Updated installer and operations documentation.

## Browser result

The production standalone build for Küresel Hat passed Chromium desktop/mobile
verification:

- homepage and selected public pages: 200;
- dynamic Küresel Hat brand and canonical/OpenGraph identity;
- no AYSTEK Haber logo/name leak;
- responsive menu with no overflow;
- RSS, sitemap, robots and manifest correct;
- unknown and cross-site Host rejected with 421;
- registered `www` alias accepted;
- no browser page errors.

## Automated verification

- Core targeted Mongo integration tests: passed.
- Core middleware/settings/router/context tests: passed.
- Core full `go test ./...`: passed (including the 289.142-second seed package).
- Core API and seed builds: passed.
- Frontend env check: passed.
- Frontend contract tests: 8/8 passed.
- Frontend lint: passed with zero warnings/errors.
- Frontend production build: passed.
- Installer `go test ./...`: passed.
- Installer source/embedded public-site env parity: passed.
- All touched-repository `git diff --check` checks: passed.

The installer repository's general env-contract script exposes an unrelated
pre-existing panel-template drift. It is not caused by this task; the public-site
source and embedded templates modified here are byte-identical and tested.

## Safety incident and resolution

Starting the full Core binary initially activated the personal-assistant agency
scheduler. It created eight agency articles and six story threads during the
test window. The exact 14 test-created documents were identified and deleted;
the window was rechecked and contained zero remaining records. The deletion is
not recoverable from MongoDB, although source data can be fetched again.

The new scheduler opt-out was then used for all browser checks, and no new
agency article/thread was created.

## Important operational finding

Server-side rendering produces multiple Core public requests per browser page.
The default 60 requests/minute limit was exhausted during the multi-page smoke
test because requests shared the renderer network identity. Sitemap calls were
then collapsed from per-language/per-link fan-out into shared multilingual data
sets. A focused sitemap + homepage + alias production-build smoke test passed at
the unchanged default limit without a Core 429. Production must still set
`PUBLIC_RATE_REQ` from a real concurrency/load test and monitor 429 responses;
the focused smoke test is not a capacity result.

## Not yet a production launch

The code path is complete, but Küresel Hat should not be publicly launched until:

- real editorial content is approved;
- corporate/legal pages are completed;
- logo and favicon are uploaded;
- social/analytics values are verified;
- site language selection is reviewed;
- immutable release tags are chosen;
- rate-limit/load testing is complete;
- DNS, Nginx and TLS are installed and tested;
- external SEO/feed/isolation smoke tests pass;
- backup and rollback steps are recorded.

`public_site_ready=false` remains descriptive metadata, not a traffic block.
It should change only after product-owner launch approval.

## Completion checklist

- [x] Existing project standards and canonical architecture followed.
- [x] No duplicate local project created for Küresel Hat.
- [x] Server deployments separated by filesystem and runtime identity.
- [x] Existing Core domain models reused.
- [x] Targeted safe seed implemented and applied.
- [x] Public domain selection made fail-closed.
- [x] Shared frontend identity made runtime-dynamic.
- [x] Brand and SEO leakage prevented.
- [x] Discovery feeds made site-aware.
- [x] Installer profiles and tests added.
- [x] Browser desktop/mobile verification completed.
- [x] Operator documentation synchronized.
- [ ] Editorial and legal launch content completed.
- [ ] Production DNS/TLS deployment completed.
- [ ] Production load/capacity test completed.
- [ ] Product owner launch approval recorded.

Detailed reference:
`done/TASK_060_SHARED_NEWS_PUBLIC_SITE.md`.
