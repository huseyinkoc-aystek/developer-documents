# Flutter Kisisel Asistan - Current State Notes

**Date:** 2026-06-25  
**Project path:** `/workspace/flutter/kisisel_asistan`  
**Purpose:** AYSTEK Core integrated personal assistant app for people, items, locations, travels, reminders, news, agency feeds, and AI chat.

---

## Standing Standards For Future Work

All future work on this app must follow `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`.

Non-negotiable rules to apply on every change:
- Work as an enterprise senior engineer, not as a patch generator.
- Keep object-oriented/domain-oriented ownership clear.
- Remove obsolete or replaced code in the same work item.
- Prefer dynamic/user-admin managed behavior instead of hard-coded business behavior.
- Keep localization complete for every user-visible text.
- Keep tests and analyzer clean before considering work complete.

---

## Technical Snapshot

- Flutter: `3.44.2` stable
- Dart: `3.12.2`
- Package: `kisisel_asistan`
- App version: `1.0.0+1`
- State management: `flutter_riverpod`
- Navigation: `go_router` with `StatefulShellRoute.indexedStack`
- Local database: `sqflite`
- Secure token storage: `flutter_secure_storage`
- Background jobs: `workmanager` + `background_fetch`
- Notifications: `flutter_local_notifications`
- Backend integration: AYSTEK Core Personal Assistant module
- Dev backend URL currently hard-coded in app services: `http://192.168.1.3:9090`
- Tenant host header: `pers-assit.aysteknolojileri.com`

Source size:
- `lib`: 79 Dart files
- `test`: 2 Dart test files
- non-generated app Dart total: about 22,900 lines
- total Dart including generated l10n: about 40,800 lines

Largest non-generated files:
- `lib/screens/travels/travel_detail_screen.dart` - 1,866 lines
- `lib/screens/aystek/news_prefs_screen.dart` - 1,817 lines
- `lib/screens/ai/ai_chat_screen.dart` - 1,764 lines
- `lib/screens/home/home_screen.dart` - 1,178 lines
- `lib/screens/travels/add_edit_travel_screen.dart` - 1,025 lines

These files are feature-rich but now large enough to be refactoring candidates.

---

## App Architecture

Current structure:

```text
lib/
  main.dart
  core/
    constants/
    extensions/
    router/
    theme/
  database/
    database_helper.dart
  l10n/
  models/
  providers/
  screens/
    ai/
    aystek/
    home/
    items/
    locations/
    onboarding/
    persons/
    reminders/
    search/
    settings/
    travels/
  services/
  widgets/
```

The app is not a strict Clean Architecture implementation yet. It has useful domain models and services, but `DatabaseHelper`, Riverpod notifiers, and some screens carry too much responsibility. Future work should gradually move toward clearer object ownership:
- model: domain state and parsing
- repository: persistence/API contracts
- service: business orchestration
- provider/notifier: UI state coordination
- screen/widget: rendering and user interaction only

---

## Main Functional Areas

### Local personal assistant data

Local/offline features are backed by `sqflite` through `DatabaseHelper`.

Current local entities:
- Person
- Location
- Item
- ItemPhoto
- Travel
- TransportSegment
- AccommodationInfo
- VehicleRental
- TravelMemory
- Reminder
- AISettings
- WebSearchSettings
- Provider keys
- Chat messages

Database version is currently `7`.

### AYSTEK Core integration

Primary service:
- `lib/services/aystek_api_service.dart`

Current behavior:
- Uses singleton `AystekApiService.instance`.
- Stores access/refresh tokens in `FlutterSecureStorage`.
- Uses custom `_TenantAwareClient` to send `host: pers-assit.aysteknolojileri.com`.
- Supports auth, profile, news preferences, news feed, public source catalog, user-added RSS sources, agency providers, agency credentials, and agency feed.

Important concern:
- `_kBaseUrl` is hard-coded to dev IP in both `aystek_api_service.dart` and `news_notification_service.dart`.
- This should become environment/flavor/config driven before production.

### News and agency flow

Screens:
- `news_tab_screen.dart`
- `news_feed_screen.dart`
- `news_prefs_screen.dart`
- `agency_feed_screen.dart`
- `news_article_screen.dart`

Current navigation:
- `/news` opens a 3-tab news shell: feed, sources/preferences, agencies.
- `/settings/news-prefs` opens `NewsPrefsScreen` standalone.
- `/news/article` expects `NewsArticle` through `GoRouterState.extra`.
- `/news-feed` remains as standalone deep-link compatibility route.

Important concern:
- News sources and news agencies must remain separate domains.
- News sources are AYSTEK Core managed RSS / web / AI scan sources.
- News agencies are premium providers only; public RSS is legacy-only and must not be used for new agency credentials.
- Agency provider lists and credential types must come from AYSTEK Core, not hard-coded mobile fallback lists.
- Agency feed reads premium articles pulled/stored by AYSTEK Core integrations connector/action flows.
- The obsolete standalone `agency_credential_screen.dart` path was removed; agency account binding now goes through `news_prefs_screen.dart`.
- Haber Masası meeting/agenda prompt generation uses `services/news_agenda_prompt_builder.dart` and includes both AYSTEK Core news sources and premium agency articles as separate context sections.
- Mobile API parsing accepts both `data` and `items` response shapes for PA news/source/agency endpoints.
- AYSTEK Core `personal_assistant.GetAgencyProviders` returns only providers configured for PA via `metadata.pa_display`; generic/core-only news agencies are not shown in the mobile app.
- AYSTEK Core seeds now treat news agencies as premium API providers. Legacy PA RSS integration connector/schedule seed records are soft-deleted by `internal/seeds/pa_agency_integrations_seed.go`.
- AYSTEK Core `SeedNewsAgencies` now writes premium agency catalog records for both the default admin site and the `pers-assist` site. PA display/auth/fetch metadata lives on the `agency_providers` catalog record, so AYSTEK Core Next and the mobile app read the same source of truth.
- The old runtime `personal_assistant/service/agency_seed.go` duplicate seed was removed. Premium agency catalog/connector setup is owned by the AYSTEK Core seed layer, not by the PA runtime service.
- PA premium connector linking is site-aware (`tenant_id + site_id + slug`), preventing connector IDs from being attached to the wrong site when the same tenant has admin and personal-assistant sites.
- AYSTEK Core boot now runs an idempotent PA agency catalog/integration guard before the PA scheduler starts. This keeps the pers-assist agency catalog available even when full seed execution is disabled in non-dev environments.
- AYSTEK Core news agency admin routes are now under `/admin/news-agencies/providers`, matching AYSTEK Core Next. The older generic `/admin/providers` path was removed from this module ownership.
- Mobile agency credentials are only considered "Bağlı" when AYSTEK Core returns a server credential ID. Local-only credentials are shown as "Senkron bekliyor" and do not imply that premium agency news can be pulled yet.
- AYSTEK Core Next news agency admin API uses `/admin/news-agencies/providers`, matching the Go routes.
- AYSTEK Core Next news agency connector tab lists API connectors only and shows the PA auth/fetch action contract readiness for linked connectors.

### AI chat

Screens/services/models:
- `screens/ai/ai_chat_screen.dart`
- `screens/ai/chat_tab_screen.dart`
- `services/ai_service.dart`
- `services/news_agenda_prompt_builder.dart`
- `models/chat_mode.dart`
- `models/ai_settings.dart`
- `models/ai_action.dart`

Supported providers:
- OpenAI
- Claude
- Gemini
- Groq
- Perplexity
- Custom endpoint

Modes:
- assistant
- newsDesk
- coffeeFal
- tarotFal
- cooking
- suggestions

Important concern:
- `ChatMode` labels, prompts, hints, and many user-facing strings are hard-coded in Dart.
- Per AYSTEK standards, these should move to l10n/config where user-visible.
- AI model defaults are hard-coded in enum values. This may be acceptable as technical defaults, but user/admin-configurable provider/model management should be considered for enterprise quality.

### Background notifications

Primary service:
- `lib/services/news_notification_service.dart`

Current behavior:
- Android: WorkManager periodic task with 15 minute minimum interval.
- iOS/Android fallback: `background_fetch`.
- Uses secure storage token and Host header in background isolate.
- Supports breaking, developing, and agency notifications.
- Notification tap routes to `/news`.

Important implementation notes:
- Uses `ExistingPeriodicWorkPolicy.keep`, which is correct for `workmanager ^0.9.0`.
- Uses top-level `@pragma('vm:entry-point')` callbacks.

Analyzer notes:
- `flutter analyze` is clean as of the 2026-06-25 agency/news-source separation pass.

---

## Navigation Map

Main shell branches:
- `/home`
- `/news`
- `/chat`
- `/travels`
- `/settings`

Full-screen/special routes:
- `/onboarding`
- `/search`
- `/news-feed`
- `/news/article`

Nested routes include:
- `/persons`
- `/persons/add`
- `/persons/:personId`
- `/persons/:personId/edit`
- `/persons/:personId/locations/add`
- `/persons/:personId/locations/:locationId`
- `/persons/:personId/locations/:locationId/items/add`
- `/persons/:personId/locations/:locationId/items/:itemId`
- `/persons/:personId/ai`
- `/travels/add`
- `/travels/:travelId`
- `/travels/:travelId/edit`
- `/travels/:travelId/memories/add`
- `/settings/ai`
- `/settings/about`
- `/settings/aystek-account`
- `/settings/aystek-login`
- `/settings/aystek-register`
- `/settings/news-prefs`
- `/settings/reminders`

---

## Localization State

ARB files:
- `app_tr.arb`
- `app_en.arb`
- `app_de.arb`
- `app_fr.arb`
- `app_es.arb`
- `app_pt.arb`
- `app_ar.arb`
- `app_zh.arb`
- `app_ja.arb`
- `app_ru.arb`
- `app_hi.arb`

Current key counts:
- Turkish: 478
- Other 10 languages: 469 each

Observed mismatch:
- Turkish has metadata keys not present in English: `@homeItems`, `@homeLocations`, `@homeTravels`

Important concern:
- Key count parity is mostly good, but the app still has many hard-coded user-visible Turkish strings in Dart files. This violates the standing standard for user-visible text.

---

## Validation Results

### `flutter analyze`

Result: clean as of 2026-06-25.

Command:
- `flutter analyze`

Result:
- `No issues found!`

### `flutter test`

Result: passing as of 2026-06-25.

Command:
- `flutter test`

Result:
- `All tests passed!`

### Existing test coverage

Current tests:
- `test/widget_test.dart` - placeholder only
- `test/news_prefs_layout_test.dart` - news source and agency tab layout test
- `test/news_agenda_prompt_builder_test.dart` - verifies premium agency articles enter the AI agenda prompt as first-class context

Testing is currently too light for the size and business importance of the app.

---

## Platform/Branding State

Android:
- Active namespace/applicationId: `com.aysteknolojileri.kisisel_asistan`
- Android manifest label: `Kişisel Asistan`
- HTTP cleartext is allowed only for local/dev domains in `network_security_config.xml`.

Stale platform branding remains:
- README title still says `kt_kurum_portal`.
- Web manifest/title still says `kt_kurum_portal`.
- Windows/macOS/Linux project names still contain `kt_kurum_portal`.
- iOS/macOS bundle identifiers still contain `com.example.ktKurumPortal`.
- Windows metadata still contains `com.example`.
- Android has stale old package file: `android/app/src/main/kotlin/com/example/kt_kurum_portal/MainActivity.kt`.

These are cleanup candidates under the "remove obsolete work" rule.

---

## Notable Technical Debt / Risks

### 2026-07-05 Social / Stories Scan Notes

- Flutter bottom navigation now uses `/stories` instead of the old `/travels` branch.
- Flutter has local Story and StoryEntry models, local sqflite persistence, story list/detail/add-entry screens, and best-effort remote sync to `/pa/stories`.
- AYSTEK Core has protected `/pa/stories` endpoints, `pa_stories` / `pa_story_entries` collections, story DTOs, service methods, and repository methods.
- The current feature is a private user timeline/story journal, not yet a full social platform.
- Social-platform primitives are incomplete: `Follow` exists only as a Flutter model; there is no DB table, provider, backend model, route, permission, admin screen, feed, comment, like, privacy, moderation, or notification workflow yet.
- The story share bar currently communicates a social intent, but it is not backed by a real share/follower contract.
- Flutter analyzer and tests are clean after fixing the premium-agency prompt-builder test type and making the v10 `server_id` story migration idempotent.

### 2026-06-26 News Agency Alignment Notes

- AYSTEK Core `news_agencies` detail page no longer has a separate "Abonelik" surface.
- That surface only stored plan/quota fields and was not used by the Personal Assistant scheduler, premium agency connector actions, mobile app, or notification flow.
- Premium agency access now stays on the enterprise path: provider catalog status + linked `integration_connectors` + action readiness + user credentials.
- Legacy news-agency-local connector/subscription/ingestion index and translation remnants were removed from the active module surface to avoid duplicate, idle work.

1. Hard-coded backend URL appears in multiple services.
   - `aystek_api_service.dart`
   - `news_notification_service.dart`
   - Should become shared environment/flavor/config.

2. Large screens carry too much responsibility.
   - `travel_detail_screen.dart`
   - `news_prefs_screen.dart`
   - `ai_chat_screen.dart`
   - `home_screen.dart`
   - Refactor into feature widgets/controllers/services as future work.

3. `DatabaseHelper` is a large mixed repository.
   - Consider domain repositories over time: persons, items, travels, reminders, chat/settings.

4. User-visible strings remain hard-coded in Dart.
   - Especially in chat modes, error fallback UI, settings/help content, and some home/offline banners.

5. Test suite is not reliable.
   - Analyzer has warnings.
   - `flutter test` fails.
   - Existing coverage is minimal.

6. README and platform branding are stale.
   - Project still carries `kt_kurum_portal` and `com.example` in multiple platform files.

7. AI provider/model defaults are hard-coded.
   - Acceptable short-term as technical defaults, but enterprise direction should make provider/model catalog dynamically manageable.

8. Agency credential flow may be duplicated.
   - Review whether `AgencyCredentialScreen` and `news_prefs_screen` agency sheets represent two active flows or one obsolete flow.

---

## Recommended Next Work Order

1. Make analyzer clean.
2. Fix or rewrite the failing `news_prefs_layout_test`.
3. Normalize l10n and move hard-coded user-visible strings into ARB files.
4. Centralize AYSTEK backend config and Host header in a single config object/service.
5. Clean stale platform branding and old package artifacts.
6. Review and unify agency credential flows.
7. Split the largest screens into smaller feature widgets/controllers.
8. Add meaningful tests around:
   - news preferences loading/saving
   - auth restore/login/logout
   - AYSTEK API parsing
   - database migrations
   - chat mode state/history
   - notification preference scheduling

---

## Operational Notes

ADB workflow from the main standards document still applies:

```bash
adb connect 192.168.1.2:PORT
adb -s 192.168.1.2:PORT install -r build/app/outputs/flutter-apk/app-debug.apk
```

Do not run `adb disconnect` unless explicitly requested.
