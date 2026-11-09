# GurselMobile - Current State Notes

**Date:** 2026-07-09
**Project path:** `/workspace/php/GurselMobile`
**Purpose:** Legacy Yii2 microservice for Gürsel Turizm (school/personnel transport). Serves: (1) a customer-facing balance top-up website ("kart bakiye yükleme"), (2) an internal back-office ("yonetim") covering institutions/vehicles/routes/students/staff/contracts/discounts, (3) a large internal+external JSON service layer (`modules/services`) consumed by SAP and by the parent-facing portal "GurselVB22", and (4) an in-vehicle QR-code safety-confirmation flow (`modules/qr_alert`) that is part of a "child left behind in vehicle" (mahsur kalma) alarm system.

This is a survey/inventory document, not a migration plan. Uncertain points are marked explicitly.

**Cross-check:** independently re-verified against source in `LEGACY_GURSEL_APPS_CODEX_REVIEW_2026-07-09.md`, which confirms the QR flow is a vehicle-level safety check (not student attendance — attendance is the separate RFID/card-tap system) and adds the AYSTEK Core bounded-context mapping / migration sequence / do-not-migrate list. Read that file alongside this one before planning the migration.

---

## 1. Overview

- Framework: Yii2 (`yiisoft/yii2-app-basic` template), `yiisoft/yii2: ~2.0.13`, PHP `>=5.4` declared in composer.json (code style and `\yii\base\Object` usage in `models/User.php` confirm an old Yii2 minor version; PHP 5/7-era codebase, not PHP 8 native — see tech debt on `utf8_encode`).
- Key extensions (from `composer.json`):
  - `yiisoft/yii2-swiftmailer` - mail (SMTP via Gmail, see below).
  - `yiisoft/yii2-mongodb` (`^2.0`) - primary datastore for most business collections is MongoDB, not the relational `db` (see Data model).
  - `kartik-v/*` (yii2-widgets, yii2-grid, yii2-dynagrid, yii2-detail-view, yii2-mpdf, yii2-export, yii2-tree-manager, yii2-money, yii2-editable, etc.) - large set of Kartik admin-grid/form widgets used throughout `modules/yonetim` CRUD screens.
  - `mihaildev/yii2-ckeditor` - rich text editing (likely for contract/announcement text).
  - `codemix/yii2-excelexport` - Excel export (PHPOffice/PhpSpreadsheet under the hood via composer.lock, not directly required at top level — no direct `phpoffice/*` line in composer.json; excel export goes through `codemix/yii2-excelexport`).
  - No Facebook SDK dependency found in `composer.json` (contrary to what might be expected from similar apps) — not present in this repo.
- Two datastores configured:
  - Relational MySQL (`config/db.php`): `mysql:host=94.73.170.201;dbname=kartlitasimacilik`, user `huseyinkoc`, password in plaintext in the file (see Tech debt).
  - MongoDB (`config/mongodb.php`): DSN with embedded username/password (value redacted here) pointing at `192.168.1.89:58987/GurselMobileMDB` — this is where almost all domain models actually live (`yii\mongodb\ActiveRecord`).
- Identity/user model (`models/Kullanicilar.php`) is itself a MongoDB `ActiveRecord` implementing `IdentityInterface` — auth state lives in Mongo, not MySQL.
- `web/index.php` hardcodes `YII_DEBUG = true` and `YII_ENV = 'dev'` with only a comment ("comment out when deployed to production") — see Tech debt/security.
- Git repo (unusual for this codebase family — see Recent activity): only 2 commits, both recent, "İlk commit" (first commit) style — this looks like the repo was only recently placed under git version control, not evidence of the app itself being new.

---

## 2. Domain modules (`modules/`)

### `modules/qr_alert` — in-vehicle QR safety confirmation

- `Module.php`: isolated module, own layout (`views/layouts/main-qr.php`), no shared app menu/assets.
- `controllers/DefaultController.php`: single real action, `actionIndex($id)`.
  - GET `/qr-alert/<id>` (24-hex-char Mongo ObjectId, enforced via URL rule in `config/web.php`: `'qr-alert/<id:[a-f0-9]{24}>' => 'qr-alert/default/index'`) renders a small Vue-in-PHP-view password prompt (`views/default/index.php`), scoped to one vehicle (`KurumBazliAracListesi`).
  - POST to the same URL (CSRF-protected, JSON in/out) validates a password (`PasswordForm`, min 4 chars) against the vehicle's registered "responsible" person(s) (`QrResponsible`), or falls back to a `registered_for => 'general'` catch-all responsible.
  - Every attempt (invalid vehicle id, vehicle not found, validation failure, password mismatch, success) is logged to `QrLog`.
- No dedicated attendance/roll-call action beyond this password check exists in `qr_alert`; see section 4 for what this flow is really for.

### `modules/odeme` — payment (Akbank virtual POS / card balance top-up)

- Customer-facing (not an API for GurselVB22): `DefaultController` under RBAC (`Yetkiler::getYetkiGruplari('kisi_sistem_girisi')`) except `akbanksanalpos` callback, which is intentionally `allow => true` (bank redirects here without a session).
- Flow: `actionIndex` builds an order (`SanalPosForm`) with a synthetic order number `'AB'.$kurumSisNum.time().'N'.$kartNum.'KTC'`; `actionKartbilgileri` builds Akbank 3D-Secure POST parameters (hash = `sha1($clientId.$oid.$amount.$okUrl.$failUrl.$rnd.$storekey)`); bank redirects back to `actionAkbanksanalpos`, which re-validates the HASH, builds a raw XML `CC5Request` provisioning call over cURL (`CURLOPT_SSL_VERIFYPEER => 0`), parses the XML response with `strpos`/`substr` (no XML parser), and on `Approved` calls `KartYuklemeHareketleri::gurselEYukleme(...)` to credit the transport card balance.
- Bank merchant credentials (`bilgi_1..bilgi_6`) are pulled from the `Sanalposlar` Mongo collection, not hardcoded in this controller (good) — but SSL verification is disabled on the provisioning cURL call (bad).

### `modules/services` — JSON service layer (SAP + GurselVB22 API surface)

~30 controllers, ~9k LOC. Notable ones:
- `VeliservisodemeservisleriController.php` ("Parent Payment/Service Services") — the clearest evidence of the **GurselVB22 (parent portal) integration**. Actions: `sozlesme`/`tumdonemsozlesmeleri`/`sozlesmeonaykaydi`/`sozlesmesorgula` (contract text, all-period contracts, contract-approval recording, contract lookup), `donemindirimoranlari` (period discount rates), `odemegecisi`/`odemealmadurumu`/`odemebilgisi`/`odetkinkurumlar` (payment-transition logging/status/active-institutions — this looks like GurselVB22 owns its own payment gateway and calls back here to record/query completed transactions, rather than this being a duplicate payment gateway), `velitckayitsorgula`/`velitckayit` (parent identity/TC-number lookup and registration, including blood type `ogr_kan_grubu`), `kardesindirimitalebi`/`kardesindiriminigonder` (sibling discount request + SAP submission), `banka-taksit-*` (bank installment lookups), `servis-on-kaydi` (service pre-registration), `kan-grubu-durumu`/`kisi-kan-grubu-zorunlulugu` (blood-type-requirement flags).
- `MusteriservisleriController.php` — personnel registration for two corporate accounts hardcoded by institution number (`0012010210`, `0012010269` "Mercedes-Benz Pazarlama"/"Mercedes Fabrika"), plus Trio Mobil RFID card-read ingestion (`triomobilKartOkumaBilgileri`, hardcoded Trio credentials as class constants).
- `SapentegrasyonislemleriController.php` — the largest controller (3253 lines). Two-way sync with SAP: institution/staff registration, vehicle/route registration, full school-person (student) record upsert (`kisibilgikayitislemi`, includes parent names/TC numbers/addresses/lat-long/blood group), route-contact and no-show notifications, "phone number login" (`actionSistemgirisi`) that returns full student + parent PII keyed by phone number, device registration for push (`sistemebaglikisicihazkaydi`), OneSignal-triggered notifications (`onesignalmesafebildirimi`, `onesignalyoklamabildirimi`, `onesignalduyurubildirimi`), and route/discount/contract queries. `actionSistemgirisi` is very likely the endpoint the GurselVB22 mobile/web client actually calls to pull a parent's student/route/attendance data.
- `AlertsController.php` — the cron-facing endpoint for the mahsur-kalma (child-left-behind) alarm; see section 6.
- `AlertsTempController.php` — near-duplicate of `AlertsController.php` (different bugs, different variable names, same hardcoded password value) left in the tree; appears to be an abandoned earlier draft. Dead code (see Tech debt).
- `BaglantianahtariController.php` — issues a bearer-style token (`WebServisBaglantiListesi::yeniBaglantiOlustur`) after checking `WebServisKullanicilari` username/password (stored and compared as **plaintext**, see `servisKullaniciDogrulama`: `->andWhere(['sifre'=>$sifre])`).
- `GuvenliiplistesiController.php` / `models/GuvenliIpListesi.php` — CRUD for an IP allow-list collection (`guvenli_ip_listesi`), but no controller in `modules/services` was found actually enforcing it via `beforeAction` (see Tech debt/security — the whitelist model exists but does not appear to gate any of the service endpoints reviewed).
- `WebservisbaglantilistesiController`, `WebserviskullanicilariController`, `WebservisislevleriController`, `BkmbakiyeyuklemeController` — CRUD/back-office screens for the service-user/connection/function-registry collections backing the above; not reviewed line-by-line (time-boxed), but they are standard Gii-generated CRUD controllers over the `WebServis*` models.

### `modules/yonetim` — back-office ("management")

By far the largest module: 71 models, 62 controllers (Gii-generated CRUD pattern: `IndexController` + `ModelSearch` classes). Represents essentially the whole business domain: institutions (`Kurumlar`), vehicles (`KurumBazliAracListesi` — also see qr_alert dependency), routes (`KurumBazliGuzergahlar`), drivers (`Soforler`), students/school-people (`OkulKisiListesi`, `OkulDonemBazliKisiListesi`, `OkulGuzergahBazliKisiListesi`), staff/customers (`KurumBazliMusteriListesi`, `KurumBazliPersonelListesi`, `PersonelListesi`), parent identity (`Velikimlikbilgileri`, `VeliKimlikBilgileriArama`), contracts (`Sozlesmeler`, `VoSozlesmeOnaylari`, `SozlesmeTurleri`), discounts (`KardesIndirimTalebi`, `DonemBazliIndirimOranlari`), payments/POS config (`Sanalposlar`, `KurumBankaTaksitleri`, `OdemeyeGecisler`), QR responsible-person registry (`QrResponsible`, `QrResponsibleSearch`), QR scan log search (`QrLogSearch`), alarm log (`AlarmLog`, `AlarmLogController`), attendance/roll-call notifications (`OkulKisiYoklamaBildirimleri`, `OkulKisiGuzergahOgrGelmemeBildirimleri`), RFID card taps (`AracGirisCikisKartOkutmalari`), device/push registry (`SistemeBagliKisiCihazlari`), RBAC (`Yetkiler`, `YetkiGruplari`, `YetkiAdlari`, `KullaniciYetkiGruplari`, `KullaniciYetkiDurumlari`), reservations (`Reservations`, `MusteriRezervasyonlari`, `Shifts`), driver-time-record import (`TrioPersonelDtr`, `ImportedFile` CSV import), Google Maps API key registry (`GoogleMapApiKeys` — plaintext API key stored in Mongo, editable via a plain CRUD form).

### `modules/raporlar` — reports

Very thin: `Module.php` + one `DefaultController` gated by `raporlar_modulu` permission, whose `actionIndex` just renders a list of report links from `ModulGorselBaglantilar` (a "module visual links" registry model). No report-generation logic found inside `modules/raporlar` itself — actual report data/rendering presumably lives behind the links this controller renders (not traced further; time-boxed).

---

## 3. Data model (selected core entities, all MongoDB `ActiveRecord` unless noted)

| Entity | Collection | Purpose | Key fields |
|---|---|---|---|
| `models\Kullanicilar` | `kullanicilar` | Back-office/portal login identity, implements `IdentityInterface` | `kullanici_eposta`, `kullanici_sifresi` (hashed via `Yii::$app->security`), `auth_key`, `access_token`, `kullanici_durumu` |
| `yonetim\models\KurumBazliAracListesi` | `kurum_bazli_arac_listesi` | Vehicle + assigned driver/hostess (per institution) | `plaka` (plate), `cam_numarasi`, `guzergah_kodu`, `sofor_*` (driver name/phone/docs), `hostes_*`, document expiry fields (`arac_muayene`, `arac_trafik_sigortasi`, etc.) |
| `yonetim\models\QrResponsible` | `qr_responsibles` | Person authorized to confirm the in-vehicle QR check for a given vehicle (or `registered_for=general`) | `name`, `email`, `password_hash` (bcrypt via `generatePasswordHash`), `arac_id`, `status` |
| `qr_alert\models\QrLog` | `qr_logs` | Every QR-page load/submit attempt | `arac_id`, `responsible_id`, `status` (`success`/`fail`/`invalid_id`), `ip`, `ua`, `created_at`, `note` |
| `yonetim\models\AlarmLog` | `alarm_log` | Dispatched "no scan" alarm record | `kurum_id`, `arac_id`, `plaka`, `emails[]`, `phones[]`, `reason`, `window`, `sent_at` |
| `yonetim\models\OkulKisiListesi` | (Mongo, name not directly confirmed but implied `okul_kisi_listesi`) | Student/school-person master record incl. parent PII | `ok_musteri_numarasi`, `ok_tc_numarasi`, `ok_anne_*`/`ok_baba_*` (mother/father name, phone, address, email), `ok_alinacagi_*`/`ok_birakilacagi_*` (pickup/drop route, address, lat/long, time) |
| `yonetim\models\OkulKisiYoklamaBildirimleri` | `okul_kisi_yoklama_bildirimleri` | Attendance/roll-call notification record | `ok_musteri_numarasi`, `cam_numarasi`, `guzergah_kodu`, `plaka`, `yoklama_durumu`, `yoklama_gonderimi` (sent flag) |
| `yonetim\models\AracGirisCikisKartOkutmalari` | `arac_giris_cikis_kart_okutmalari` | RFID card tap on vehicle entry/exit (the actual per-student attendance signal, distinct from the QR flow) | `kart_num`, `kart_num_trio`, `plaka`, `lat`, `lon` |
| `yonetim\models\Sanalposlar` | (Mongo) | Virtual POS / bank merchant config | `bilgi_1..bilgi_6` (client id, store key, credentials, success/fail URLs, provisioning URL) |
| `services\models\WebServisKullanicilari` | `web_servis_kullanicilari` | External/internal service-API credential | `kullanici_adi`, `sifre` (**plaintext**, compared with `==` in Mongo query), `islevi`, `ip_dk_guvenli_is_sayisi` |
| `services\models\GuvenliIpListesi` | `guvenli_ip_listesi` | IP allow-list entries | `ip_adresi` |
| `yonetim\models\Velikimlikbilgileri` | (Mongo) | Parent identity record used by `velitckayit*` API | mother/father TC, name, birth year; student blood group |

Relational MySQL (`kartlitasimacilik` DB) is configured but no model in the reviewed directories was found using `yii\db\ActiveRecord` against it — every domain model reviewed extends `yii\mongodb\ActiveRecord`. The MySQL connection's actual current usage inside this app is **unclear**; it may be vestigial/shared with a sibling app.

---

## 4. Business logic & workflows

### QR-based "child left behind" (mahsur kalma) safety flow, end to end

This is **not** a per-student attendance QR scan. It is a vehicle-level safety-compliance confirmation:

1. Each vehicle (`KurumBazliAracListesi`) has a physical QR code inside it (not generated dynamically in code reviewed — likely a static label encoding `/qr-alert/<vehicle _id>`; no QR-image-generation code was found in `modules/qr_alert`, so **QR code creation itself may happen out-of-band or in a script not reviewed**).
2. A driver or hostess scans the QR at the end of a route, lands on `GET /qr-alert/<id>` (`qr_alert\controllers\DefaultController::actionIndex`), sees a small Vue password page (`views/default/index.php`).
3. They submit a password. It's checked against any `QrResponsible` registered for that specific vehicle, falling back to a `registered_for=general` responsible. This is meant to attest "I walked through the vehicle and checked for a child left behind."
4. Every attempt (good or bad) is written to `QrLog` with vehicle id, responsible id (if matched), status, IP, user-agent, timestamp.
5. A separate cron-triggered endpoint, `services\AlertsController::actionIndex` (see section 6), scans, per institution with `mahsur_kalma_alarmi = true`, whether a *successful* `QrLog` exists inside that institution's configured time windows (`alarm_saat_araliklari`, e.g. `"09:00-10:00"`) on active alarm days (`alarm_aktif_gunler`). If a window closes (plus a 5-minute grace period) with no successful scan, it fires an email (via a **separate Gürsel Tur email-relay service**, `kt.gurseltur.com.tr/.../epostagonder`, called with hardcoded basic-auth-style credentials in the URL body) and logs an `AlarmLog` entry; SMS sending to `phones` is stubbed as a comment (`// SMS gönderimi için $phones kullanılabilir`) — **SMS alarm delivery for this specific flow is not implemented**, despite the model already collecting phone numbers.
6. Recipients are the vehicle's registered `QrResponsible` emails plus institution-level `alarm_gonderilecek_epostalar`/`alarm_gonderilecek_numaralar`.

Separately, actual **per-student attendance/roll-call** is driven by RFID card taps at vehicle entry/exit (`AracGirisCikisKartOkutmalari`, ingested via Trio Mobil GPS hardware and `MusteriservisleriController::actionGunlukokunankartlar`) and recorded/notified through `OkulKisiYoklamaBildirimleri` and `OkulKisiGuzergahOgrGelmemeBildirimleri` (student-did-not-board notifications), with parent push notifications dispatched via OneSignal from `SapentegrasyonislemleriController` (`onesignalyoklamabildirimi` etc.) and `BilgiGonderimi::onesignalIleAndroidPushMesajGonder`/`...IOSPushMesajGonder`. These two systems (QR vehicle-safety-check vs. RFID student-attendance) are independent but share the same vehicle/route/institution data model.

### Payment: two distinct paths, not obviously duplicated

- `modules/odeme` is this app's **own** customer-facing top-up UI: Akbank 3D-Secure virtual POS, credits a transport smart-card balance (`KartYuklemeHareketleri::gurselEYukleme`). This is GurselMobile's own gateway integration, used when a logged-in "kisi_sistem_girisi" user tops up their card from this website.
- `services\VeliservisodemeservisleriController` (`odemegecisi`/`odemealmadurumu`/`odemebilgisi`) does **not** initiate a card payment; it records/queries payment-transition status (order id, institution, amount, installment, bank) against `OdemeyeGecisler`. This strongly suggests GurselVB22 (the parent portal) runs its own payment/POS flow and calls back into GurselMobile purely to log the transition and let other systems (SAP?) query "did this order pay." **Not fully confirmed** — no code in this repo shows GurselVB22 actually calling these endpoints; the inference is based on the shape of the API (accepts a completed order + bank + amount, doesn't drive a bank redirect).

---

## 5. External integrations

- **SMS**: `BilgiGonderimi::smsGonder` posts XML to `https://processor.smsorigin.com/xml/process.aspx` using hardcoded credentials (username/password pair, values redacted here) (`modules/services/models/BilgiGonderimi.php:24-26`).
- **Push notifications**: OneSignal REST API, two separate apps (Android/iOS) with hardcoded `app_id` + base64 REST API keys in `onesignalIleAndroidPushMesajGonder`/`onesignalIleIOSPushMesajGonder` (`modules/services/models/BilgiGonderimi.php:258-259,302-303`).
- **Payment gateway**: Akbank virtual POS, 3D-Secure, raw XML `CC5Request` provisioning API (`modules/odeme/controllers/DefaultController.php`).
- **SAP**: bidirectional sync via `SapentegrasyonislemleriController` (student/staff/vehicle/route/contract/discount data in both directions) — this is the largest single integration surface in the app.
- **Trio Mobil** (GPS/RFID vendor): `MusteriservisleriController` polls `http://takip.triomobil.com/soap/GetRFIDLogs` with 3 different hardcoded credential pairs as class constants (`TRIOMOBIL_SVC_USER/PASS`, `..._MP_USER/PASS`, `..._MF_USER/PASS`).
- **A separate Gürsel Tur email-relay microservice**: `https://kt.gurseltur.com.tr/web/services/kartbakiyeislemleri/epostagonder`, called with hardcoded creds (username/password pair, values redacted here) (`AlertsController.php:273-274`, and again in `AlertsTempController.php:202`).
- **"AYSTEK" personnel transfer endpoint** (of historical interest for this migration effort): `BilgiGonderimi::aystekPersonelTransferi` posts to `http://asp.gurseltur.com.tr/web/services/staffregister/ucstaff` with hardcoded creds (username/password, values redacted here) (`modules/services/models/BilgiGonderimi.php:349-353`). This shows an **earlier, informal AYSTEK integration attempt already exists** in this legacy codebase (endpoint literally named "aystek..."), separate from the current AYSTEK Core platform — likely dead/superseded, not verified live.
- **GurselVB22 (parent portal) integration** — no direct code from GurselVB22 is in this repo, but the following are the concrete integration points a parent-facing client would use:
  - `services\VeliservisodemeservisleriController` — contracts, sibling discount, blood-group requirement, service pre-registration, payment-status logging. Auth is per-request `kullanici`/`sifre` fields checked via `LoginForm::login()` (i.e., reuses the same Mongo `kullanicilar` identity used by the web back-office — **no separate API key/JWT for this "parent service" surface**).
  - `services\SapentegrasyonislemleriController::actionSistemgirisi` — phone-number-keyed lookup returning full student + parent PII + route/pickup/drop/notification-preference data; same `kullanici`/`sifre` auth pattern.
  - `services\BaglantianahtariController` — issues a `WebServisBaglantiListesi` token after `WebServisKullanicilari` username/plaintext-password check; unclear if any other controller actually requires/validates this token (not observed in the controllers reviewed — most check `kullanici`/`sifre` per-request instead of a bearer token).
  - Shared MongoDB (`GurselMobileMDB`) is the most concrete "shared state" with GurselVB22 if that app connects to the same cluster — **not confirmed from this repo alone**; flagged for cross-check against the GurselVB22 survey.

---

## 6. Scheduled jobs (`crons/`)

Only one file: `crons/gursel_mobile_mahsur_alarm.php`.

- Not a Yii console command — a standalone PHP script (comment: `// /opt/cron/mahsur_alarm.php`) meant to be invoked directly by system cron.
- Does a single `curl` GET to `https://ats.gurseltur.com.tr/web/services/alerts/index?token=...` (shared-secret token value redacted here; note: this hostname is `ats.gurseltur.com.tr`, not the same host serving this codebase — **the actual alarm-check controller `services\AlertsController` in this repo may not be the one this cron calls**, or `ats.gurseltur.com.tr` is a production alias for this same app; not confirmed).
  - The `alertsToken` value hardcoded in the cron script matches `config/params.php`'s `alertsToken` (value redacted here), which **is** the token `AlertsController::beforeAction` checks — so it is very likely this cron drives that controller, just via its production hostname.
- Dead code in the script itself: after `print_r($res); exit();` on line 17-18, the following error-logging block (lines 20-24, writing to `/var/log/mahsur_alarm.log`) is unreachable.
- No cadence/schedule is defined in the repo (no crontab file present); cadence must live on the server's actual crontab, which is outside this repo. `AlertsController`'s own logic (`firstDivisible5MinuteAfter`, "±30s tolerance" comment) implies it's expected to run roughly **every 5 minutes**, but this is inferred from the code, not configured here.

---

## 7. Auth & permissions

- Human/web session auth: standard Yii2 `user` component, `identityClass = app\models\Kullanicilar` (Mongo-backed), `enableAutoLogin = true`, login via `LoginForm` (username = email, password checked with `Yii::$app->security->validatePassword` against a bcrypt hash).
- Separate, informal "customer"/"personnel" logins exist outside the Yii user component: `SiteController::actionPersonelgirisi` implements phone + SMS one-time-code login stored in **PHP session** (`$session['personel_login']`), not the Yii identity system; the OTP itself is generated via `BilgiGonderimi::kodOlustur` (weak 4-digit numeric code from `rand()`, not `random_int`) and delivered via the SMS integration above.
- RBAC: custom, not Yii2 RBAC (`yii\rbac`). `Yetkiler` (permissions) × `YetkiGruplari`/`KullaniciYetkiGruplari` (permission groups × user-group membership) is a hand-rolled Mongo aggregation check, `Yetkiler::getYetkiGruplari($yetkiAdi)`, used as `AccessControl`'s `allow` value throughout `modules/odeme`, `modules/raporlar`, and presumably `modules/yonetim` controllers.
- QR/attendance API auth specifically:
  - `qr_alert\DefaultController` — CSRF-token-protected form POST, password checked against `QrResponsible` bcrypt hash. No rate limiting, no lockout, no CAPTCHA on repeated password guesses against a 24-char-hex-guessable-if-leaked vehicle id; `QrLog` records every attempt but nothing appears to consume it for throttling (see Tech debt).
  - `services\AlertsController` (the cron-facing alarm endpoint) — single shared static token in a query string (value redacted here), checked with `!==` against `params['alertsToken']`; no HMAC, no per-caller identity, no expiry, no HTTPS enforcement inside the app itself (relies on `https://ats.gurseltur.com.tr` externally).
  - The `services\*` JSON API surface consumed by GurselVB22/SAP (`VeliservisodemeservisleriController`, `SapentegrasyonislemleriController`, `MusteriservisleriController`) — authenticates per-request with a `kullanici`/`sifre` pair run through the same `LoginForm`/`Kullanicilar` identity check used for the human web login (i.e., a *service account* shares the exact same login mechanism as a human user, no scoping, no token, no distinguishing claims). `BaglantianahtariController` implements a separate token-issuance flow against a different plaintext-password-checked collection (`WebServisKullanicilari`), but no controller reviewed actually requires that issued token on subsequent calls — it appears to be a partially-built or abandoned separate auth scheme.
  - `GuvenliIpListesi` (IP allow-list) exists as a data model/CRUD screen but was not found wired into any `beforeAction` in the controllers reviewed — likely intended as a gate for the services layer but **not enforced in code reviewed**.

---

## 8. Known issues / tech debt

Concrete, file:line-referenced findings (grep- and read-verified against this repo; not from memory):

1. **Debug/dev mode hardcoded in the production entrypoint.** `web/index.php:4-5`: `defined('YII_DEBUG') or define('YII_DEBUG', true); defined('YII_ENV') or define('YII_ENV', 'dev');` with only a comment saying to change it — if the real deployment doesn't override these via a server-level define, this exposes full stack traces, DB/Mongo DSNs, and file paths on any error.
2. **Hardcoded credentials committed in plaintext**, at minimum:
   - MySQL DB password: `config/db.php:13` (value redacted here), plus a second commented-out credential pair at lines 5-8.
   - MongoDB DSN with embedded credentials (value redacted here): `config/mongodb.php:6`.
   - SMTP password (value redacted here): `config/web.php:34`.
   - Cron shared secret (value redacted here): `config/params.php:13` and `crons/gursel_mobile_mahsur_alarm.php:3`.
   - SMS provider credentials: `modules/services/models/BilgiGonderimi.php:25-26`.
   - OneSignal REST API keys (2 apps): `modules/services/models/BilgiGonderimi.php:258-259,302-303`.
   - Trio Mobil GPS vendor credentials (3 pairs): `modules/services/controllers/MusteriservisleriController.php:31-35`.
   - Email-relay and "AYSTEK" staff-transfer service creds: `modules/services/controllers/AlertsController.php:273-274`, `modules/services/models/BilgiGonderimi.php:352-353`.
3. **Plaintext password storage/comparison for the internal service-API user table.** `modules/services/models/WebServisKullanicilari.php:112-121` (`servisKullaniciDogrulama`) compares `sifre` with a direct Mongo equality match — no hashing, unlike the human `Kullanicilar` model which correctly uses `Yii::$app->security`.
4. **Dead/duplicate controller with the same hardcoded secret.** `modules/services/controllers/AlertsTempController.php` is a near-identical, buggier earlier draft of `AlertsController.php` (confirmed via diff), left fully in place, still containing its own copy of a hardcoded password (line 202, value redacted here) and its own token check. Should be deleted, not migrated.
5. **`utf8_encode()` used on all inbound service credentials/params** (e.g. `modules/services/controllers/MusteriservisleriController.php:81`, `SapentegrasyonislemleriController.php` throughout `actionKurumvepersonelkayitislemi`/`actionSistemgirisi`/etc.) — `utf8_encode` assumes ISO-8859-1 input and is deprecated/removed behavior changes in PHP 8.2+; this alone will need remediation on any PHP upgrade, and likely already mangles non-Latin1 Turkish characters (ğ, ş, ı, etc.) in names/addresses.
6. **No rate limiting or lockout on the QR password endpoint.** `modules/qr_alert/controllers/DefaultController.php::actionIndex` logs every attempt to `QrLog` but nothing throttles repeated POSTs against a guessable/leaked vehicle id; a 4-character-minimum password (`PasswordForm::rules()`, `modules/qr_alert/models/PasswordForm.php:14`) is brute-forceable.
7. **IP allow-list model unused.** `services\models\GuvenliIpListesi` / `GuvenliiplistesiController` provide CRUD for an IP allow-list but no `beforeAction` in the reviewed `services` controllers actually calls `ipGuvenliMi()` — looks like partially-implemented defense-in-depth that was never wired up.
8. **Unfinished/placeholder second auth scheme.** `BaglantianahtariController` issues a connection token (`WebServisBaglantiListesi::yeniBaglantiOlustur`) but no other reviewed controller validates it on subsequent requests — either dead code or an integration point that was never finished.
9. **SMS channel for the mahsur-kalma alarm is unimplemented.** `services\AlertsController::gonderAlarm` collects `$phones` and comments `// SMS gönderimi için $phones kullanılabilir` (`AlertsController.php:264`) but never calls `BilgiGonderimi::smsGonder` — only email is actually sent, despite phone numbers being gathered from both institution config and `QrResponsible`.
10. **Weak OTP generation.** `BilgiGonderimi::kodOlustur` (`modules/services/models/BilgiGonderimi.php:29-48`) builds the personnel-login SMS code with `rand()`, not a CSPRNG (`random_int`/`Yii::$app->security->generateRandomString`); combined with a 4-digit numeric-only code space this is a low-entropy OTP.
11. **Hardcoded test/demo data reachable via a real model method.** `KurumBazliAracListesi::refreshFromJsonText()` (`modules/yonetim/models/KurumBazliAracListesi.php:407-onward`) embeds ~2000 lines of one school's actual fleet data (plates, driver names, personal phone numbers, chassis/engine numbers) directly in source as a "refresh" migration helper. This is real personal data (driver names/phones) committed to source control and should not ship as-is; also dead-code-shaped (looks like a one-off script that should have been run once and removed).
12. **Disabled TLS verification on multiple outbound integrations**: `curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false)` (and `VERIFYHOST, false`) appears in `crons/gursel_mobile_mahsur_alarm.php:9`, `services/controllers/AlertsController.php` (`sendEmail`), `services/models/BilgiGonderimi.php` (`smsGonder`, `aystekPersonelTransferi`), and `odeme/controllers/DefaultController.php` (Akbank provisioning call) — MITM risk on every one of these outbound calls.
13. **Yii2-scaffold dead code left in place**: `commands/HelloController.php` is the unmodified Yii2 example console command; `models/User.php` is the unmodified Yii2 tutorial in-memory user model with hardcoded `admin`/`admin` and `demo`/`demo` credentials (`models/User.php:16-27`) — not wired to the real `Kullanicilar` identity class, but still present and could be confusing/risky if ever accidentally referenced as `identityClass`.
14. **`web/index-test.php`** guards itself with an IP check (`REMOTE_ADDR` must be `127.0.0.1`/`::1`, `web/index-test.php:4-6`) but still hardcodes `YII_DEBUG = true` and loads a separate acceptance-test config — low risk if the IP guard is trustworthy behind the real proxy setup, worth confirming `REMOTE_ADDR` isn't spoofable through the front proxy.
15. **Raw XML built via string concatenation/`str_replace` templating** for the Akbank provisioning request (`modules/odeme/controllers/DefaultController.php`, `actionAkbanksanalpos`) instead of a proper XML writer — fragile if any field ever contains `<`, `>`, or `&`.

---

## 9. Recent activity (`git log --oneline -20`)

```
837a467 İlk commit: vendor hariç, config/web.php dahil
53a94a9 first commit
```

Only two commits exist, both generic "initial commit" messages (Turkish: "first commit, excluding vendor, including config/web.php"). This tells us the git history was only just established for this app — it provides **no signal about what was recently worked on feature-wise**; the actual development history (if any existed before) is not preserved. File modification timestamps on disk (from directory listings during this survey) show `modules/qr_alert`, `modules/odeme`, and the root `.gitignore` most recently touched (`Sep 23-24 2025`), and `crons/` most recently touched (`Oct 8 2025`) — consistent with the QR/mahsur-alarm feature and the cron script being the most recently active area of work, which lines up with the amount of in-progress-looking code found there (the abandoned `AlertsTempController.php` draft, the "İyi çalışmalar dileriz" alarm email copy, comments referencing 2026 dates in `Sozlesmeler`/`Velikimlikbilgileri`-adjacent code).

---

## 10. Open questions for cross-check against the GurselVB22 survey

- Does GurselVB22 connect to the same MongoDB cluster/DB (`GurselMobileMDB` at `192.168.1.89:58987`) directly, or only through the `services` HTTP API documented in section 5?
- Is `services\SapentegrasyonislemleriController::actionSistemgirisi` actually the endpoint GurselVB22's client calls for "find my child's route/pickup info by phone," or does GurselVB22 have its own copy of this logic?
- Does GurselVB22 ever call `BaglantianahtariController` for a token, and if so, which endpoint(s) actually check that token? (Not found in this repo.)
- Confirm whether `https://ats.gurseltur.com.tr` (used by the cron) and this GurselMobile codebase are the same production deployment, or whether `ats.gurseltur.com.tr` is a different/renamed app also worth surveying.
