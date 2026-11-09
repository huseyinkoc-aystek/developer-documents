# TASK 061 - Sports Matches & Standings: Multi-Sport Model and Provider Auto-Sync

Technical reference for the sports module's extension from a single
hardcoded Süper Lig snapshot to a multi-sport, multi-tier competition
catalog with an opt-in, pluggable external-provider auto-sync engine.

## Objective

The existing `sports` module (one `LeagueSnapshot` document per competition,
edited entirely by hand through the admin panel) needed to:

- support more than one sport (football, basketball, volleyball) and more
  than one domestic tier (top division, second division) plus international
  competitions, without inventing a second module;
- let an editor optionally pull standings and recent results from a real,
  external, third-party sports-data API instead of always typing them in by
  hand, on a **per-competition** basis (manual and automated data must never
  silently overwrite each other);
- support more than one external provider, selectable by the operator, with
  no code change required to switch;
- reflect correctly on the public news site's homepage widget and its
  `/spor` hub/detail pages, in the site's active language.

No fabricated data is acceptable anywhere in this feature: if a provider
does not have data for a competition, the snapshot must stay whatever it
already was (manual or last-known-good auto) rather than being overwritten
with an empty or wrong result.

## Data model (`internal/modules/sports/model/league.go`)

`LeagueSnapshot` gained, on top of the pre-existing `Code`, `Season`, `Week`,
`Standings`, `Matches`, `Active` fields:

```go
type SportType string          // "football" | "basketball" | "volleyball"
type CompetitionTier string    // "top" | "second" | "international"
type DataSource string         // "manual" | "auto"

type ExternalRef struct {
    Provider string `bson:"provider"`   // e.g. "apisports", "thesportsdb"
    LeagueID string `bson:"league_id"`  // provider-specific competition id
    Season   string `bson:"season"`     // provider-specific season value
}
```

- `Sport`/`Tier` default to `football`/`top` on read via `Normalize()` so
  snapshots written before these fields existed keep resolving to sensible
  values — but see the "Known pitfall" section below: this default only
  applies to code that reads through `Normalize()`, not to a raw MongoDB
  query filter.
- `ExternalRef` is `nil` by default — a competition is **manual-only** until
  an admin explicitly opts it in from its own edit page. The sync job never
  touches a snapshot without one.
- `Source` records who last wrote the data (`manual` from the admin
  `Upsert` path, `auto` from the sync job's `SyncUpdate` path) plus
  `SourceSyncedAt`. These two write paths are structurally incapable of
  colliding: `Service.Upsert` always stamps `manual`, `SyncService.RunOnce`
  always stamps `auto`.
- `DisplayOrder int` controls which competition's tab the public homepage
  widget opens by default. Added after the multi-sport widget was tested
  live: the list endpoint's default sort (`sport, tier, code`) is
  alphabetical, so a `basketball` competition with no data yet would rank
  ahead of `football`'s real content purely by sport name. Lower shows
  first; ties still fall back to `sport, tier, code`.

## Provider-abstraction and sync engine

### `sports/service/sync.go` — `SyncService`

```go
type SyncConfig struct {
    Enabled bool; Provider string; APIKey string; IntervalMinutes int
}
type SyncProvider interface {
    FetchSnapshot(ctx, apiKey string, ref ExternalRef, sport SportType) ([]Standing, []Match, error)
}
```

`SyncConfig` is read from the generic `settings` collection under key
`sports.sync` (`scope: "module"`, tenant-scoped like every other settings
document — see "Known pitfall" for why that matters). `RunOnce(ctx)`:

1. loads the tenant's `SyncConfig`; returns `ErrSyncDisabled` /
   `ErrSyncConfigMissing` if not usable;
2. looks up the configured provider in the registry map;
3. lists every snapshot in the tenant/site scope carried by `ctx`;
4. for each snapshot with a non-nil `ExternalRef` whose `Provider` matches
   the active provider, calls `FetchSnapshot`, replaces `Standings`/
   `Matches`, stamps `Source=auto` + `SourceSyncedAt`, and writes it back
   through the same repository `Upsert` the admin path uses.

A snapshot with no `ExternalRef`, or one opted into a *different* provider
than the currently active one, is left untouched — so switching providers
never lets the old one race the new one, and un-opted-in competitions are
never silently converted to auto.

### Providers (`sports/provider/`)

Two adapters exist, selected by `SyncConfig.Provider`; the registry is a
plain `map[string]SyncProvider` in `cmd/api/main.go`, so adding a third
provider is adding one file plus one registry line — no other code changes.

**`apisports.go`** — api-sports.io family (`v3.football`, `v1.basketball`,
`v1.volleyball`). Auth via `x-apisports-key` header. Confirmed live and
working for football's *documented* shape (`/standings`, `/fixtures`).
**Free-plan limitation, confirmed live 2026-09-28: standings/fixtures are
only available for seasons 2022-2024 — the current season returns
`{"plan":"Free plans do not have access to this season, try from 2022 to
2024."}` for football, basketball and volleyball alike.** This applies
platform-wide, not per-sport — there is no free-tier path to current-season
data through this provider.

**`thesportsdb.go`** — thesportsdb.com. Auth via URL path segment
(`/api/v1/json/{key}/...`), not a header; falls back to TheSportsDB's own
published free/test key (`"3"`) when no key is configured, so a competition
can be smoke-tested before anyone pays for anything. Season format is
`"YYYY-YYYY"`, not a bare year. **Confirmed live and working for the
current 2026-2027 season** — this is the provider that actually satisfies
"pull real current-season data for free" for football. Confirmed
*confirmed not usable* for Turkish basketball/volleyball on the free key:
`lookuptable.php` returns a genuinely empty body (not an error, not JSON at
all) for `Turkish Basketbol Super Ligi` (league 4475) and the one generic
`Turkish Volleyball League` entry (4543) across every season tried;
`eventspastleague.php` (match results) does return *some* data for
basketball even where standings don't. TheSportsDB's paid tier
($9/month, per their own pricing page as of this check) raises rate limits
and unlocks their V2 API — their documentation does not promise it adds
league coverage that the free tier lacks, so it is not a confirmed fix for
this specific gap.

**Net result:** football has a fully working, free, current-season path
(TheSportsDB). Basketball and volleyball do not, on either integrated
provider's free tier, as of this check. Both remain fully usable
*manually* — this is a data-source availability gap, not a code
limitation, and the architecture (provider registry, `ExternalRef` opt-in)
requires no change to add a working source for these sports later.

### Scheduler (`internal/platform/jobs/sports_sync_scheduler.go`)

Mirrors `route_fares/jobs/expiring_fare_scheduler.go` /
`fx_sync_scheduler.go`'s aggregate-then-loop pattern: `ListDistinctTenantSiteScopes()`
enumerates every `(tenant_id, site_id)` pair with at least one snapshot,
and `RunOnce` is called once per scope with that scope's tenant/site on the
context. Retry: 2 attempts, 5s backoff, only for genuine provider errors —
`ErrSyncDisabled`/`ErrSyncConfigMissing` are treated as the normal state
for most scopes, not failures, and are neither retried nor recorded as an
ops failure metric.

### Known pitfall — found and fixed live, 2026-09-28

**The scheduler never actually ran a single tick from the day it was
written until this fix**, despite passing every unit test and code review.
`StartSportsSyncScheduler`'s outer loop read `sports.sync` once per tick
*before* knowing which tenant/site it was about to serve, using a context
with no tenant set, to decide the poll interval and whether to bother
calling `runSportsSyncTick` at all. `sports.sync` is a tenant-scoped
settings document (`tenantfilter.WithTenantID` requires a tenant in
context or returns an error) — that top-level call therefore *always*
failed, and the code's `else` branch on that failure was `shouldRun =
false`, skipping the tick unconditionally, every single time, forever.

`StartFXSyncScheduler` has the structurally identical top-level call and
does **not** have this bug, because its failure branch does nothing
(`shouldRun` keeps its pre-set default of `true`) — it fails *open*, not
closed. The sports scheduler was fixed to match: the top-level read only
ever adjusts the poll interval when it happens to succeed; whether to
actually sync is decided per-scope, inside `RunOnce`, against a properly
tenant/site-scoped context, which was always correct.

**Lesson for whoever touches a scheduler like this next:** a scheduler
that reads a tenant-scoped setting before it has resolved which tenant
it's about to act on cannot get a real answer from that read. Decide
per-scope, inside the loop over `ListDistinctTenantSiteScopes()`, not
before it — and if a "cheap top-level check" is kept anyway (e.g. to pick
a poll interval), its failure must never be interpreted as "so don't run,"
because for every scope except an accidental single-tenant coincidence, it
will always fail. This was caught by live end-to-end verification against
production (Süper Lig's `Source` staying `manual` no matter how the config
was set), not by any test — the existing unit tests mocked `SyncService`
at a level that never exercised the real scheduler loop or a real
`tenantfilter` call. **A "does the scheduler actually run" integration
test does not exist yet and is a good candidate for Faz 6 / next-session
hardening.**

## Credential handling (`sports/credcrypt/`)

A third near-identical AES-256-GCM single-value encryption package
(mirrors `fleet_tally/credcrypt`, `personal_assistant/credcrypto`) — this
codebase's established convention is one small credential-encryption
helper per module family rather than a shared platform package. Pass-
through (unencrypted storage) when `SPORTS_CRED_KEY` is unset, so dev
environments keep working; production should set it (documented in
`.env.example`).

## Admin UI (`aystek-core-next`)

- `/admin/sports` (list): tenant-wide sync settings card — enable toggle,
  provider **select** (not free text — see below), interval, API key.
  A configured key never round-trips to the browser: `GetConfig` returns
  only `api_key_configured: bool` and, since this task, the last four
  characters (`api_key_suffix`) plus `updated_at`/`updated_by`, so the
  field can show a masked "•••• ed5c, saved <time>" confirmation instead
  of an input that always renders empty and looks like the save silently
  failed — the concrete complaint that prompted this change.
- `/admin/sports/[code]` (editor): per-competition "Otomatik
  Senkronizasyon" section — opt-in checkbox + provider/league_id/season,
  plus the new `display_order` ("Sıra") field.
- Provider is a `<select>` sourced from `SYNC_PROVIDER_OPTIONS` in
  `lib/api/sports.ts` (one entry per registry entry in `cmd/api/main.go` —
  the two lists must be kept in sync by hand; there is no shared source of
  truth between the Go registry and this TS constant).

## Public site (`aystekhaber-com` / `kureselhat.com`)

- `SportsWidget.tsx`: homepage multi-sport tab switcher, `leagues[0]` as
  the default tab — now meaningful because of `DisplayOrder`.
- `SportsSnapshotPanel.tsx`: standings/matches panel, used sliced on the
  homepage and full on `/spor/[code]`.
  - **Form-letter localization bug, found and fixed live:** the win/draw/
    loss streak (`standing.form`, e.g. `"WWDWL"` or `"GBGGG"`) can now
    come from two different alphabets — manually-entered data historically
    used Turkish letters (G/B/M), external providers return the
    international W/D/L convention. The color-coding function checked only
    for literal `'G'`/`'B'`; a plain provider `'W'` matched neither and
    fell through to the *loss* color, and every letter was displayed
    exactly as the data source happened to write it regardless of the
    reader's language. Both alphabets now normalize to one canonical
    value (`normalizeFormResult`), which is both colored correctly and
    localized for display (`localizeFormResult`) against the page's active
    `lang`, not the data source's language.
- `/spor` (hub) and `/spor/[code]` (detail): confirmed **absent from the
  live kureselhat.com deployment** during this task, despite existing on
  `origin/main` for over a day — the running container had been built
  from a stale base (a single hand-patched file on top of an old
  checkout, not a fresh build from current source) for an earlier,
  unrelated nav-bar fix. Rebuilding fully from `origin/main` restored
  both routes (verified 200 OK). **Lesson: a single-file patch on a
  running site's source is a legitimate quick fix for that one file, but
  it does not catch up the rest of the source tree — the next real change
  to that site needs a full rebuild from the current commit, not another
  patch on top of the stale base, or entire feature commits silently stay
  missing from production indefinitely.**

## Verified end-to-end (production, 2026-09-28)

- Süper Lig (`super-lig`) opted into TheSportsDB (league 4339, season
  `2026-2027`): scheduler ran automatically after the fix, `Source`
  flipped to `auto`, `SourceSyncedAt` stamped, standings replaced with
  real current data (5 teams reported by the provider — see below), and
  the change was visible on the live kureselhat.com homepage and `/spor`
  hub within one tick.
- Trendyol 1. Lig (`tff-1-lig`, tier `second`) opted into the same
  provider (league 4676): synced successfully with real current data
  (Kayserispor leading, 8 matches played) — confirms the mechanism
  generalizes across competitions, not just the one originally wired up.
  Correctly absent from the homepage widget (which only shows `tier=top`)
  and correctly present on `/spor`.
- Basketbol Süper Ligi (`bsl`) opted into the same provider (league 4475):
  sync correctly **did not** overwrite the snapshot — `FetchSnapshot`
  returned `ErrTheSportsDBResponseShape` for the empty standings body,
  `RunOnce` reported `updated: false`, `Source` stayed `manual`. This is
  the "no fabricated fallback" contract working as designed, not a bug.

**Known, disclosed, unresolved data-completeness gap:** TheSportsDB's free
tier currently reports only 5 of Süper Lig's ~18-20 teams (confirmed via
direct query, independent of this codebase). This is not a parsing bug —
`parseTheSportsDBStandings` correctly reflects exactly what the provider
returns. Full-league coverage on a free tier was not found during this
task; it requires either a paid provider plan or manual completion of the
missing teams.

## Editor role and menu-visibility fixes (found during this task, unrelated to sports specifically but discovered via it)

Testing the sports feature as a real non-admin `editor`-role user (not
`admin`) surfaced two independent, pre-existing gaps affecting the whole
platform, not just sports:

1. **`role_seed.go`'s `editorRead`/`editorWrite` allowlists never included
   `"sports"`** — any editor got `"missing module permission"` opening
   `/admin/sports` even though the menu item itself was visible. Fixed by
   adding the module to both maps (read + full write, matching how
   `posts`/`galleries`/etc. are already granted).
2. **`view_profiles_seed.go`'s `allMenuKeys` was a hand-maintained list
   that had drifted 76 keys behind the real `menu_seed.go` catalog**
   (276 tracked vs. 352 actually registered, confirmed by test). Any menu
   key absent from this list defaults to *visible* in every site's
   editorial profile (the reverse of the intended "new menus stay hidden
   until reviewed" default) — this is why an editor scoped to a pure news
   site could see an unrelated `transport-platform-suite` (library/PMS)
   menu group. Root-caused and fixed structurally, not just patched: the
   menu-seed list (`menuSeeds`) was promoted to a package-level var, and
   `allMenuKeys` is now `allMenuKeysList()` — computed *from* `menuSeeds`,
   not a second copy of it — so this specific class of drift cannot recur.
   Separately, the news-editor allowlist (`aysTekHaberVisibleMenuKeys`)
   itself was expanded by 19 keys (galleries, ads, comments, events,
   classifieds, interviews, condolences, biographies, news sources/
   agencies, digital editions, sliders, surveys, localized content,
   references, contacts, company directory) — all of which already had
   editor permission grants and simply had no menu path to reach them.

A site's already-created `view_profiles` document is a point-in-time
snapshot computed at creation (or last explicit refresh) — it does **not**
pick up later `allMenuKeys`/allowlist changes automatically. A new
operation, `seeds.RefreshNewsSiteViewProfile(ctx, db, code)` (CLI:
`--news-site=<code> --refresh-view-profile`), force-recomputes and
persists a named site's profile on demand. **Any future change to
`menu_seed.go` or a site's visible-keys allowlist must be followed by
running this for every already-provisioned news site**, or the fix only
applies to sites provisioned after the code change.

## Tag site-isolation fix (same investigation, same root class of gap)

`taxonomy.Tag` had no `SiteID` at all — every tag query scoped by
`tenant_id` only, so kureselhat.com's tag list/picker showed
aystekhaber.com's tags and vice versa. Added `Tag.SiteID`; every new tag
now gets one from context. Existing tags (all 26, confirmed, all
predating this field) intentionally keep showing tenant-wide via an
explicit `$or` fallback (`site_id` = current site, OR absent, OR empty) —
mirrors the already-established `authors` module pattern
(`withTenantAndOptionalSite`/`strict` parameter) rather than a destructive
retroactive migration that would have had to guess which site each shared
tag "belongs" to. `Category` was deliberately left untouched — out of
this task's requested scope.

## Files changed

**aystek-core:** `internal/modules/sports/{model,dto,repository,service,
handler,provider,credcrypt}/*`, `internal/platform/jobs/
sports_sync_scheduler.go`, `internal/platform/ops/sports_sync_metrics.go`,
`internal/seeds/{menu_seed,view_profiles_seed,role_seed,
translation_seed,news_sites_seed,site_editor_seed}.go`,
`internal/modules/taxonomy/{model,repository}/*`, `cmd/api/main.go`,
`cmd/seed/main.go`, `.env.example`.

**aystek-core-next (admin panel):** `app/(dashboard)/admin/sports/
{page,[code]/page}.tsx`, `lib/api/sports.ts`, `app/(dashboard)/admin/
posts/{page,create/page}.tsx`, `app/components/CoverThumb.tsx` (new).

**aystekhaber-com (public, shared with kureselhat.com):**
`components/news/{SportsWidget,SportsSnapshotPanel,NewsChrome}.tsx`,
`app/(auth)/login/page.tsx` — unrelated stale-cross-user-context fix found
via the same live-testing pass (see summary doc).

## What Faz 6 / a future session should pick up

- A working free (or affordably paid) current-season data source for
  Turkish basketball and volleyball — not found during this task.
- An integration test that actually exercises
  `StartSportsSyncScheduler`'s goroutine end-to-end (the class of bug this
  task found is exactly what unit-level mocking cannot catch).
- Faz 0-4 items already closed are tracked in
  `SPORTS_MATCHES_STANDINGS_WORKPLAN_TR.md` (§4) with commit references;
  this document supersedes that workplan's §5 continuation notes now that
  Faz 3 and Faz 5 are both complete.
