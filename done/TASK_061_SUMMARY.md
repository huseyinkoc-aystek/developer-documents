# TASK 061 - Executive Summary

## Outcome

The `sports` module now supports multiple sports (football, basketball,
volleyball), multiple tiers (top/second/international) per sport, and an
opt-in, per-competition auto-sync engine that pulls standings and recent
results from a real external provider — selectable by the operator, with
no code change to switch providers. Verified end-to-end against real
production traffic: Süper Lig and Trendyol 1. Lig now update automatically
from TheSportsDB and the result renders correctly on kureselhat.com's live
homepage and `/spor` pages.

## What actually works today

- **Football, current season, free:** TheSportsDB. Confirmed live.
- **Football, 2022-2024 seasons only, free:** api-sports.io (their whole
  free plan is capped to those seasons, platform-wide — confirmed via
  direct query, not documentation-only).
- **Basketball/volleyball, current season:** no working free source found
  on either integrated provider. Match results exist for basketball on
  TheSportsDB even where standings don't; volleyball has no usable data at
  all on the free key. This is disclosed, not silently degraded — the
  sync engine correctly refuses to overwrite a snapshot when a provider
  returns nothing, so these competitions simply stay manually-managed.

## A real bug the live testing caught that no unit test did

The auto-sync scheduler was wired, tested, and merged — and never once
actually ran a tick in production, from the moment it was deployed until
this task. Its top-level "should I run" check read a tenant-scoped setting
using a context with no tenant, which always fails; the failure branch
said "so don't run" instead of "figure it out per scope like the working
`FXSyncScheduler` does." Fixed by matching the fail-open pattern that
`StartFXSyncScheduler` already used correctly. This is exactly the class
of bug that only shows up when you actually watch a real value change in a
real database over real wall-clock time — worth remembering next time a
scheduler "passes review."

## Collateral fixes (found via the same live-as-a-real-editor testing pass)

- **Editor role couldn't open Sports at all** — `sports` was never added
  to the editor permission allowlist (`role_seed.go`), despite the menu
  item being visible. Fixed.
- **A stale-vs-live menu-visibility list let 76 unrelated menu keys leak
  into every site's editor view** — `view_profiles_seed.go`'s hand-copied
  `allMenuKeys` had drifted behind the real menu catalog (`menu_seed.go`).
  Restructured so it's now *computed from* the real catalog instead of
  copied, closing this class of drift permanently. Also expanded the
  news-editor's own visible-menu allowlist by 19 modules that already had
  permission grants but no menu path (galleries, ads, comments, events,
  and 15 others) — a real news-site editor can now actually reach the
  content types they're already allowed to manage.
- **Tags leaked across sites within the same tenant** — no `site_id` on
  `Tag` at all, so kureselhat.com and aystekhaber.com shared one tag pool.
  Added site scoping for new tags; the 26 existing tags stay tenant-wide
  (a safe default, not a data migration) rather than risking breaking
  existing post↔tag references by guessing which site each one "belongs"
  to.
- **A stale-cross-user browser-session bug** — logging in as a second,
  more narrowly-scoped user in the same browser tab could still show the
  previous, broader user's site selection and menu content, because
  `router.replace` (a client-side transition) never remounts React
  context providers, and their refetch logic only triggers on a tenant
  change, not a raw user switch. Fixed with a hard navigation after login.
- **kureselhat.com was missing two whole feature commits in production**
  — an earlier same-day hotfix had patched one file on top of a stale
  checkout instead of rebuilding from the current commit, so `/spor` and
  the multi-sport homepage widget code existed on `origin/main` but not on
  the live site. Full rebuild restored both.
- **Cover-image thumbnails looked broken (or silently missing) in the
  admin posts list and editor** — a relative image path was resolved
  against the wrong origin, and the interim fix (hide on error) made a
  broken image indistinguishable from "no image was ever set." Now shows
  a clear placeholder icon distinct from an empty state.
- **Win/draw/loss form letters showed the wrong color and the wrong
  language** — a provider-supplied `'W'` matched none of the hardcoded
  Turkish-letter checks and rendered as a loss (red) regardless of the
  actual result; now normalized and localized to the site's active
  language.

## Deployment

All fixes across `aystek-core`, `aystek-core-next` (admin panel), and
`aystekhaber-com`/`kureselhat.com` (shared public renderer) were committed,
pushed to their respective GitHub repositories, rebuilt, and deployed to
production during this task. Full backend test suite (`go test ./...`)
green at close. See `TASK_061_SPORTS_AUTO_SYNC.md` for the technical
reference, file-by-file change list, and what a future session should pick
up next (primarily: a working basketball/volleyball data source, and an
integration test for the scheduler bug class).
