# AYShuttleGursel (Kartlı Taşımacılık) - Current State Notes

> **2026-10-29 documentation update:** This July inventory is retained as a historical snapshot. For the newly source-checked Turkish analysis, measured controller/collection inventory, and corrections to the PHP, PDF dependency, ledger immutability, token-expiry, SUL, and scheduler claims below, start with [Proje analizi](./AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md) and [Teknik envanter](./AYSHUTTLEGURSEL_TEKNIK_ENVANTER_TR.md). The new review is static analysis, not a production or migration certification.

**Date:** 2026-07-09
**Project path:** `/workspace/php/AYShuttleGursel`
**Purpose:** Card-based/RFID transit ticketing and fleet-billing system for Gürsel Turizm ("Kartlı Taşımacılık" = "Card-based Transportation"), branded internally as "K-Ulaşım". Manages RFID transit cards, per-ride fare deduction on vehicle-mounted readers, balance top-ups (cash/bank/e-wallet), discount rules for institution-affiliated riders (university/company shuttle contracts), vehicle/route/driver scheduling, online reservation, invoicing, and reporting. Multiple partner institutions (universities, companies) are onboarded as tenants ("kurum") with per-tenant feature flags.

This document is a factual inventory for planning a rewrite onto the AYSTEK Core platform (Go 1.24 + Gin + MongoDB at `/workspace/go/aystek-core`, Next.js 14 frontend at `/workspace/node/aystek-core-next`). Findings were gathered by direct reading of `controllers/`, `models/`, `modules/{odeme,raporlar,services,ticket_sales,yonetim}/`, `config/`, `commands/`, `mail/`, `widgets/`, and view-directory structure (view HTML bodies were not dumped). `vendor/`, `runtime/`, `.git/`, `nbproject/`, `tests/_output/`, and compiled assets were excluded per scope.

**Cross-check:** independently re-verified against source (grep-based counts + spot re-reads) in `LEGACY_GURSEL_APPS_CODEX_REVIEW_2026-07-09.md`, which found this inventory's payment-gateway coverage (Akbank + Yapı Kredi live, BKM Express stub-only) accurate and adds the AYSTEK Core bounded-context mapping / migration sequence / do-not-migrate list. Read that file alongside this one before planning the migration.

---

## 1. Overview

### 1.1 What the app does

A full back-office + rider-facing platform for a shuttle/transport operator:

- Issues and tracks RFID cards (and paper tickets/coupons/credit-linked cards) tied to persons ("kişi") who belong to a partner institution ("kurum" - a university, company, etc.).
- Vehicle-mounted card readers call into a webservice API (`modules/services`) to deduct fare balance per ride, apply institution/person-group discount rules, and log the transaction.
- Riders/administrators top up card balance via cash, bank transfer, invoiced batch loads, or online card payment (Akbank/Yapı Kredi 3D Secure virtual POS, BKM Express nominally).
- A newer, semi-independent `ticket_sales` module sells session-based paper/QR tickets to guest customers (phone/email OTP login, Yapı Kredi payment, QR code redemption).
- A `raporlar` (reports) module provides operational and financial reporting (vehicle earnings, card movement history, occupancy, debtor cards) with Excel export.
- A `yonetim` (management) module is the primary admin back office: ~85 CRUD-style controllers/models covering cards, persons, vehicles, routes, schedules, drivers, pricing, discounts, reservations, RBAC, and integrations with partner-institution systems (university card/SIS systems, SAP, Trio Mobil GPS/telematics).

### 1.2 Tech stack

| Layer | Technology |
|---|---|
| Framework | Yii2 (`yiisoft/yii2 ~2.0.45`), PHP >= 7.4, `yii2-app-basic` template origin |
| Primary datastore | **MongoDB** (`yiisoft/yii2-mongodb ^3.0`) - almost the entire domain model is `yii\mongodb\ActiveRecord` over MongoDB collections, not SQL tables |
| Secondary datastore | MySQL (`config/db.php`) - present in config and used only by 3 legacy external-integration models (`IsikUniversitesiKisiListesi`, `OzyeginKisiListesi`, `OzyeginKartBakiyeleri`); `config/db.php`/`config/test_db.php` still point at the unmodified Yii2-scaffold default `mysql:host=localhost;dbname=yii2basic` |
| UI framework | Bootstrap 5 (`yiisoft/yii2-bootstrap5`), extensive **kartik-v widget suite** (`yii2-widgets`, `yii2-grid`, `yii2-dynagrid`, `yii2-detail-view`, `yii2-nav-x`, `yii2-export`, `yii2-tree-manager`, `yii2-builder`, `yii2-editors`, `yii2-datecontrol`, `yii2-money`, `yii2-popover-x`, `yii2-tabs-x`, etc. - all `dev-master`, unpinned) |
| Excel/PDF export | `codemix/yii2-excelexport ^2.7` and `kartik-v/yii2-export` (wraps PhpSpreadsheet). **No mpdf/PhpOffice PDF library present** - PDF export is explicitly disabled everywhere it's wired (`FORMAT_PDF => false`) |
| Mailer | `yiisoft/yii2-symfonymailer ~2.0.3` (SMTP via Office365, credentials hardcoded - see §8) |
| Captcha | `pctux/yii2-invisible-recaptcha` (Google Invisible reCAPTCHA) |
| Phone input | `borales/yii2-phone-input` |
| JWT | `firebase/php-jwt ^6.11` declared but **not referenced anywhere in the codebase** - unused/aspirational dependency |
| Payment SDKs | No Akbank/Yapı Kredi SDK in composer.json - integrations are hand-rolled XML/cURL against bank virtual-POS APIs. BKM Express PHP SDK is vendored manually (not via Composer) under `web/bkmekspress/` and is **not actually wired up** (stub action only) |
| SMS | No SMS Composer package; SMS sending is hand-rolled XML POST to a third-party gateway ("SMS Origin", see §5) |
| Maps/GPS | Google Maps JS API (hardcoded API key, see §8), Google Directions API, and a third-party fleet-telematics integration ("Trio Mobil") |
| Dev/test | `yii2-debug`, `yii2-gii` (incl. a MongoDB Gii model generator), Codeception + PHPUnit scaffolding present but no evidence of maintained business-logic test suites was reviewed in depth (out of scope: `tests/`) |

### 1.3 High-level module map

```text
AYShuttleGursel/
  controllers/            root site controllers (login, contact, help, password reset)
  models/                 root-level auth/contact models (Kullanicilar = User, LoginForm, ...)
  modules/
    yonetim/              "management" - ~85 controllers/models, the core admin back office
    odeme/                "payment" - rider-facing balance top-up, payment gateways, reservation UI, QR
    raporlar/              "reports" - operational/financial reporting, Excel export
    services/              "services" - webservice API consumed by vehicle readers, partner systems, SAP
    ticket_sales/           newer semi-independent guest ticket-purchase module (own components/services dirs)
  commands/               console commands (stock Yii2 scaffold only, see §6)
  mail/                   mailer layouts only (no message templates, see §5)
  widgets/                Alert.php (flash message widget, stock Yii2 scaffold)
  assets/                 custom asset bundles (BKM Express JS, Google Maps tracking, stop/route map editor)
  web/bkmekspress/        vendored (non-Composer) BKM Express PHP SDK - unused
```

---

## 2. Domain modules

### 2.1 `modules/yonetim/` (Management) - the core admin back office

Largest module: ~85 controllers, ~85+ models (plus `*Arama` search-model twins), ~480 view files. No `modules/yonetim/components/` directory exists; `Module.php` (`modules/yonetim/Module.php`) is a bare scaffold (no bootstrap, no module-wide access filter). Every controller implements its own `AccessControl` behavior, checking permissions ad hoc via `Yetkiler::getYetkiGruplari('<permission-key>')` (see §7).

Controllers group into these business areas (only controllers with **non-CRUD custom actions** are called out; most others are plain Gii CRUD index/view/create/update/delete):

| Area | Controllers | Notable custom actions |
|---|---|---|
| Cards & balance | `KartlarController`, `KartyuklemehareketleriController`, `KartcekimhareketleriController`, `KartiadehareketleriController`, `KartdevirhareketleriController`, `KartdurumhareketleriController`, `CevrimdisikarthareketleriController`, `MobilcekimhareketleriController` | `KartlarController::actionOkunankart/actionOkunankartislembitir` (kiosk "swipe card" polling UI); `KartcekimhareketleriController::actionBakiyecek` (manual kiosk fare deduction, calls the same `bakiyeCekmeIslemi()` engine used by vehicle readers); `CevrimdisikarthareketleriController::actionCdisibakiyeleriaktar` (batch-replays offline reader transactions) |
| Persons & groups | `KisilerController`, `KisigruplariController`, `KisikartlariController`, `KisigrubugecerliliksureleriController`, `KisigrupguzergahindirimiController`, `KisisiniflariController`, `DigerKisiler` (model only) | `KisilerController::actionUpdate` cascades email changes across 7 collections via `Kisiler::kisiEpostaDegis()` |
| Vehicles, drivers, routes | `AraclarController`, `SoforlerController`, `GuzergahlarController`, `DuraklarController`, `AracsoforulasimatamasiController`, `AracisatamasiController`, `AracisdevirleriController`, `AracisformbilgileriController`, `ArackapasiteturleriController`, `AracguvenlikbilgileriController`, `AracabaglicihazlarController`, `ExternalshuttleController` | `AracisatamasiController` + `AraclarController::actionCreate` implement a session-mediated two-controller wizard (create vehicle-job, then auto-create vehicle if missing); `GuzergahlarController::actionDuraklar` cascading stop-list dropdown |
| Scheduling | `SaatliulasimlistesiController` (SUL - "Scheduled Transport List", the central route×time×day entity), `SaatlerController`, `SaataraliklariController`, `GunlerController`, `GuzergahvardiyatipleriController`, `GuzergahvardiyakapasitelistesiController` | `SaatliulasimlistesiController::actionUpdate` cascades SUL-name renames across 5 collections; `GuzergahvardiyakapasitelistesiController::actionVerileriguncelleme` pulls capacity/pricing from an external SAP endpoint over cURL with hardcoded creds and disabled TLS (see §8) |
| Pricing & discounts | `UlasimucretleriController`, `KartsatisfiyatlariController`, `GrupindirimlerikumesiController`, `AbonelikucretleriController`, `AbonelikyuklemeleriController`, `YuklemekampanyalariController` | discount-policy configuration consumed by the fare engine (§4) |
| Reservations | `RezervasyonisatamasiController`, `SeyahatrezervasyonlariController`, `YedekrezervasyonlarController` | `RezervasyonisatamasiController::actionCreate` bulk-creates a cross-product of `sul × date × time` reservation rows; `SeyahatrezervasyonlariController::actionUetdsreport` / `YedekrezervasyonlarController::actionUetdsreport` - **UETDS** compliance report (Turkey's mandatory Ministry-of-Transport bus-passenger manifest reporting system) |
| Cash/invoicing | `GunlukkasadekontlariController`, `FaturaliyuklemelerController`, `YazarkasaposlarController`, `YazarkasaposraporlariController` | `FaturaliyuklemelerController::actionTarihduzeltme` is a one-off maintenance script (hardcoded 2020-01-01..2021-01-20 date range) left live in the controller |
| RBAC / users | `KullanicilarController`, `KullaniciyetkigruplariController`, `KullaniciyetkidurumlariController`, `YetkilerController`, `YetkiadlariController`, `YetkigruplariController`, `KullanicidurumlariController`, `KullanicikurumayarlariController`, `KullanicisifreleriController`, `SifresureleriController` | see §7 |
| Partner-institution integrations | `IsikuniversitesikisilistesiController`, `OzyeginuniservicekisilistesiController`, `SabanciuniservicekisilistesiController`, `VeriaktarimController` | `VeriaktarimController` is a data-sync dashboard triggering batch imports/balance top-ups per institution (Sabancı, Özyeğin flat-file, Özyeğin web-service) |
| Devices / app versions | `AracabaglicihazlarController`, `UygulamacihazlariController`, `UygsurumleriController`, `UygsurumgruplariController`, `TasinabilircihazturleriController` | device connectivity dashboards |
| Misc | `DuyurularController` (announcements), `BaglantilarController` (connections - includes a live QR-decrypt debug harness, see §8), `LanguageController`, `ModulgorselbaglantilarController` (module tile/icon links), `BiletlerController` |

### 2.2 `modules/odeme/` (Payment) - rider-facing payment/reservation/tracking

`Module.php` redirects errors to `odeme/error/index` and sets layout `main`.

| Controller | Key actions | Purpose |
|---|---|---|
| `DefaultController` | `index`, `bakiyeyukle`, `servisnerede`, `haritadaizle`, `directions`, `locations`, `kartbilgileri`, `kvkkonayi`, `tcentry`, `akbanksanalpos`, `yapikredi`, `yapikrediodemeonay`, `bkmekspress`, `sartname`, `gizlilik`, `iptaliade`, `musterisozlesmesi`, `musterihizmetsozlesmesi` | Balance top-up entry point, live vehicle map tracking ("servisnerede" = "where's my shuttle"), Google Directions proxy, Trio Mobil vehicle-location proxy, Akbank/Yapı Kredi 3D Secure payment flows (see §4.3), legal/consent pages. `enableCsrfValidation = false` for the whole controller (bank callback endpoints require this) |
| `HaftalikPlanController` | `data` | JSON proxy (60s cached) for "weekly plan" - upcoming scheduled departures for the rider's institution |
| `HesapozetiController` | `karthareketozeti` | Account statement: aggregates card load/withdrawal/refund history for a date range |
| `KartsizgecisController` | `bakiyecek` | "Cardless passage" - POSTs to the internal services fare-deduction endpoint with a **hardcoded password** (value redacted, see §8) |
| `QrCodeController` | `generate` | Issues a 5-minute QR code (`01TCK` prefix + random hex) tied to person/card; imports but never calls `Endroid\QrCode\QrCode` (dead/unfinished import) |
| `RezervasyonController` | `index` | Seat reservation with waitlist fallback; if the institution requires balance deduction on reservation, calls the same internal fare-deduction endpoint (again with a hardcoded password, value redacted) |

Models: `QrCode` (Mongo AR, `qr_codes` collection - cross-module QR validation, delegates `02TCK`-prefixed codes to `ticket_sales\Tickets::validateQrCode()`), `RezervasyonForm` (plain form model), `SanalPosForm` (plain form model; fixed top-up amount menu 50₺-6000₺, scenario-based rules for balance/subscription/period top-up types).

### 2.3 `modules/raporlar/` (Reports)

`Module.php` is a trivial scaffold (sets layout only). Controllers:

- **AracbazlisultopkazanclariController** - per-vehicle/SUL total-earnings report; `actionCreate` bulk-inserts rows parsed from a delimited POST string; `actionTekrargonder` resets a "sent to SAP" flag for re-export.
- **ArachareketraporuController** - vehicle-trip movement reports: boarding lists, occupancy (`doluluk`), average occupancy, vehicle work-order forms, per-SUL total earnings, instant summary report, soft-delete/undelete of trip records.
- **DefaultController** - landing page of module tiles.
- **HarekettoplamraporuController** - aggregate totals report; builds its filter model directly from `$_REQUEST` rather than `$model->load()` (tech debt, see §8).
- **KarthareketraporuController** - card movement report/rebuild, debtor-card report, daily scan-count, card-movement summary, person/SUL usage info.
- **MusteriarackazancraporuController** - customer/vehicle earnings report; has an undefined-variable bug (see §8).

None of the report "models" are SQL views - they are either real Mongo collections populated by app-level "rebuild" batch jobs (`KartHareketRaporu`, `AracHareketRaporu`, `AracBazliSulTopKazanclari`, `MusteriAracKazancRaporu`), a raw-`batchInsert` staging duplicate (`KartHareketRaporuTmp`), or plain `yii\base\Model` DTOs that compute an in-memory report on every request by iterating another collection's records in PHP and wrapping the result in an `ArrayDataProvider` (`BorcluKartlarRaporu`, `CihazIletisimHareketleri`, `HareketToplamRaporu`, `KartHareketOzeti`).

### 2.4 `modules/services/` (Webservices) - the API consumed by hardware/partner systems

`Module.php` is a bare scaffold. This module is the machine-to-machine surface: vehicle card readers, mobile apps, partner-institution systems, and SAP all call into it. **CSRF validation is disabled on nearly every controller** (`public $enableCsrfValidation = false`) since these are not browser form posts.

Two (three) inconsistent auth patterns coexist:
1. `guvenlikKontrolu($kullanici, $sifre)` - tries system-user login (`LoginForm::login()`), falling back to a vehicle-plate + plate-password check (`AracGuvenlikBilgileri::plakaSifreKontrolu()`). Used by device/vehicle-facing endpoints. Plate-password check is **bypassed entirely** for plates starting with `'P'` (pseudo-plates generated for reservations).
2. Connection-key/token pattern: client exchanges `username`/`pass` for a random token via `BaglantianahtariController`/`EntegrebakiyeislemleriController::actionBaglantianahtari`, then presents the token on later calls; per-IP call-rate limiting is applied unless the caller's IP is in the `GuvenliIpListesi` whitelist (whitelisting bypasses only the *rate limit*, not authentication itself).
3. `AystekvalidatorController` - the only controller using Yii's standard `HttpBasicAuth` filter, against `Kullanicilar::validatePassword()`.

| Controller | Purpose |
|---|---|
| `KartbakiyeislemleriController` | **The central device-facing API** (23 POST actions) - card info/balance lookup, the core "pull balance for a ride" endpoints (`bakiyecekmeislemi`, `kartbakiyecekmeislemi`, `qrbakiyecekmeislemi`, `temazsizbakiyecekmeislemi`), per-institution QR decode branches (Özyeğin external lookup, Sabancı DES-decrypt with hardcoded key - see §8), offline-transaction transfer, schedule export, generic email/SMS dispatch proxies |
| `ArachareketislemleriController` | Vehicle-trip earnings reporting from device JSON, reservation SMS reminders, waitlist promotion |
| `AracisatamaislemleriController` | Auto-creates daily scheduled vehicle-job assignments, remote/manual assignment, close-out |
| `CihazbaglantiislemleriController` | Device registration (SIM ICCID, device↔vehicle link), app-login check, user-permission lookup |
| `EntegrebakiyeislemleriController` | Institution-driven balance top-up/refund/cancel, new-card provisioning (only Sabancı branch implemented; Özyeğin/Koç are empty stubs) |
| `FaturaislemleriController` | Invoiced/un-invoiced transaction listing, invoice-sent notification, invoice backfill |
| `GuvenliiplistesiController` | Admin CRUD for the IP whitelist |
| `HataliisatamasiduzeltmeleriController` | Replays balance deductions from an erroneous vehicle assignment onto a corrected one; contains a large dead comment block with hardcoded personal phone numbers/emails from a real incident (see §8) |
| `KarthareketleriController` | Last-balance lookup, paginated/date-ranged movement history (rebuilds a temp report collection per call), "person graduated" flag, institution list with **inconsistent CORS** (per-action allow-list vs a wildcard `Access-Control-Allow-Origin: *` helper) |
| `SapentegrasyonislemleriController` | **SAP "Teacher/Student" SOAP integration** - pushes per-vehicle/route earnings (`EDolumGuncelle`) and pulls route/shift/capacity pricing (`EDolumFiyatGuncelleme`) via native `SoapClient` against `Yii::$app->params['gurselTSSoapServiceUrlCanli']`/`...Test` (`config/params.php:16-17`), with **hardcoded SOAP credentials in source** (see §8) |
| `Webservisbaglantilistesi` / `Webservisislevleri` / `Webserviskullanicilari` Controllers | Admin CRUD over webservice connection logs, defined "functions", and webservice user accounts (whose passwords are stored/compared **in plaintext**, see §8) |
| `BaglantianahtariController` | Issues connection-key tokens |
| `AystekvalidatorController` | HTTP-Basic-Auth-gated: license-plate list export, purge of cards stuck "in stock" >30 days |

### 2.5 `modules/ticket_sales/` (guest ticket purchase - newest/most self-contained module)

Has its own `components/` and `services/` subdirectories (`BaseController`, `AuthService`) - a different, more modern code style than the rest of the app, but still Yii2/Mongo. Layout: `modules/ticket_sales/views/layouts/main`.

- **Auth**: plain PHP-session-based (`Yii::$app->session->has('customer_id')`), gated in `components/BaseController.php`; login is by institution-number lookup + 6-digit SMS/email OTP (`AuthController`, `services/AuthService.php`). The declared `firebase/php-jwt` dependency is **not used anywhere** in this module.
- **CustomerController** - profile management; guest customers who don't yet exist in the core `Kisiler` table are auto-provisioned with a `Kisiler` record + a "paper ticket" `Kartlar` card + `KisiKartlari` link, bridging ticket-sales guests into the core card system.
- **PaymentController** - a second, independently-implemented Yapı Kredi POSNET 3D Secure integration (parallel to `odeme\DefaultController`'s), persisting payment-flow state server-side in a `PaymentTransaction.order_note` JSON blob (to survive the bank redirect) rather than in PHP session.
- **TicketController** - lists a customer's tickets as JSON; renders each ticket's QR via a **public third-party API** (`api.qrserver.com`) rather than generating the QR image locally.
- Models: `TicketBuyers` (Mongo AR, `ticket_buyers`), `Tickets` (Mongo AR, `tickets` - status `active/used/cancelled/expired`, 24h expiry; `validateQrCode()` is called cross-module from `odeme\QrCode` for `02TCK`-prefixed codes, and on redemption triggers a real balance-credit call into `yonetim\KartYuklemeHareketleri::gurselEYukleme()`, i.e. a used paper/QR ticket is modeled as a card top-up+immediate-use).

---

## 3. Data model

**Architecture note:** the domain model is almost entirely **MongoDB**, not relational. Every core entity is a `yii\mongodb\ActiveRecord` implementing `collectionName()` (not `tableName()`). There are **no native Yii `hasOne()`/`hasMany()` relations** anywhere in the ~85+ yonetim models - every "relation" is a hand-written lookup method (`Model::find()->where(...)`) or a raw `Yii::$app->mongodb->createCommand()->aggregate()` pipeline, keyed by **denormalized string fields** (`kart_num`, `kurum` name, `kisi_eposta` email, route/SUL name) rather than ObjectId foreign keys. Renaming any of these "key" values (card number, person email, vehicle plate, route/SUL name, group name) requires the application to manually cascade updates across 5-8 collections - this happens in several controllers/models (see §4) and is a significant data-integrity risk that a relational or FK-based MongoDB schema in AYSTEK Core should eliminate.

Only 3 models use `yii\db\ActiveRecord` (MySQL), all external-integration read tables: `IsikUniversitesiKisiListesi` (`isik_universitesi_kisi_listesi`), `OzyeginKisiListesi` (`kisi_listesi`), `OzyeginKartBakiyeleri` (`ozyegin_kart_bakiyeleri`).

### 3.1 Core entities

| Entity (model) | Collection | Key fields | Notes |
|---|---|---|---|
| `Kartlar` (yonetim) | `kartlar` | `kart_num` (unique), `kurum`, `kart_turu`, `kart_durumu`, audit fields | Card types: `KartTurleri::NORMAL_KART`, `KAGIT_BILET` (paper ticket), `KREDILI_KART` (credit-linked), `KUPON` (coupon). Statuses: `KISIDE`/"Mezun Kişide" (with a graduated person), `STOKTA` (in stock), `KULLANIMDA` (in use), `KAYIP` (lost). `beforeSave()` auto-writes a status-history row to `KartDurumHareketleri` on every status change - implicit audit trail. Balance is **not stored on the card** - see §3.2. |
| `Kisiler` (yonetim) | `kisiler` | `kurum_kayit_num`, `kurum_numarasi`, `kimlik_numarasi` (Turkish national ID, unique), `kurum`, `kisi_sinifi`, `kisi_adi/soyadi`, `uyruk`, `cinsiyet`, `kisi_telefon`/`kisi_eposta` (both unique), `kisi_durumu`, `kisi_gruplari`, `online_odeme_kapali` | Person/cardholder. Custom validator implements the Turkish TCKN checksum algorithm plus a ban-list of known-fake IDs. `beforeSave()`: on status change away from initial status (e.g., "graduated"), **cascades a forced card close** (`Kartlar`+`KisiKartlari`). |
| `KisiKartlari` (yonetim) | `kisi_kartlari` | `kisi` (email FK), `kart_num` (FK), `kisi_kart_durumu` (Etkin/Kapalı) | Person↔card association. Validators enforce no duplicate active links and no two simultaneously-active cards of the same type per person. `beforeSave()` side-effects `Kartlar` status (activate → card "Kullanımda", close → card "Stokta"). |
| `Kurumlar` (yonetim) | `kurumlar` | `kurum_adi(+_tr/_en)`, `kart_bakiye_alt_siniri` (min balance floor), `kurum_azami_yukleme_limiti` (max top-up), `kurum_sistem_numarasi`, `abonelik_sistemi`, `rezervasyon_sistemi`, `rezervasyonda_bakiye_cek`, `hizli_kartsiz_gecis`, `hatali_is_duzenleme`, `qr_crypt_mode`, `misafir_kaydi`, `online_yayin`, `kampanyali_tutar_yukleme` | **Multi-tenant config/feature-flag entity** - nearly every business rule elsewhere branches on these institution-level booleans rather than a pluggable strategy. `kurumServiceBilgisi()` is a hardcoded `switch` on institution-name constants dispatching to per-university integration classes (tight coupling). |
| `Araclar` (yonetim) | `araclar` | `plaka` (unique, 7-10 chars), `arac_durumu` (usable bool), `arac_koc_num` (in-vehicle reader device number), `yolcu_kapasitesi` | Vehicles. |
| `Soforler` (yonetim) | `soforler` | `sofor_adi/soyadi`, `sofor_telefon`, `sofor_eposta` | Drivers. |
| `Guzergahlar` (yonetim) | `guzergahlar` | `kurum`, `guzergah_adi` (unique), `baslangic_xy`/`bitis_xy` (GPS coordinates), `duraklar` (array of stop refs) | Routes; `guzergahMesafeHesabi()` computes great-circle distance between a stop and the route endpoint, capped at `EN_COK_UZAKLIK = 50.0`. |
| `Duraklar` (yonetim) | `duraklar` | `durak_adi`, `durak_xy` | Stops (GPS-tagged). |
| `SaatliUlasimListesi` ("SUL") (yonetim) | `saatli_ulasim_listesi` | `saatli_ulasim_listesi_adi`, `guzergah`, `guzergah_kodu`, `baslangic_tarihi`/`bitis_tarihi`, `sefer_turu` (1=giriş/2=çıkış), `sefer_suresi`, per-weekday booleans (`pazartesi`..`pazar`), `mobil_yayinlama`, `rezervasyon_yayinla`, `kimlik_bilgisi_zorunlu` | **The central scheduling entity** - a specific route + time + set of active weekdays; everything (fares, job assignments, reservations, reports) keys off this "SUL" id/name. |
| `AracIsAtama` (vehicle job/trip assignment) (yonetim) | `arac_is_atama` | `saatli_ulasim_listesi`, `plaka`, `sofor`, `hareket_tarihi`/`hareket_saati`, `arac_is_durumu` (AÇIK/KAPALI), `guvenlik_sorgu_kodu` | The operational unit tying vehicle+driver+scheduled-trip together for a specific date/time; fare deductions and trip reports key off this record's `is_numarasi` (job id). |
| `UlasimUcretleri` (fares) (yonetim) | `ulasim_ucretleri` | `saatli_ulasim_listesi`, `kartli_ucret` (card fare), `kartsiz_ucret` (no-card/paper fare), `kupon_ucret` (coupon fare), `gecerlilik_baslangic/bitis_tarihi` (validity window) | Time-bounded, per-SUL, 3-tier pricing. |
| `GrupIndirimleriKumesi` (discount rule sets) (yonetim) | `grup_indirimleri_kumesi` | `grup_indirimi_adi`, `indirim_orani` (0-100%), `zaman_asimi` (timeout minutes), `gunler[]`, `limit_kontrol_durumu` (G=daily/H=weekly/A=subscription), `limit` (max uses), `saat_araliklari[]` (allowed hour ranges) | Discount policy engine config - consumed by the fare-deduction engine (§4.1). |
| `KisiGrupGuzergahIndirimi` (yonetim) | `kisi_grup_guzergah_indirimi` | person-group ↔ route ↔ discount-set mapping | Links a person-group's discount eligibility to specific routes. |
| `KisiGrubuGecerlilikSureleri` (yonetim) | `kisi_grubu_gecerlilik_sureleri` | group validity windows per person | Determines a person's *currently valid* group memberships at ride time. |
| `KartYuklemeHareketleri` (top-ups) (yonetim) | `kart_yukleme_hareketleri` | `kart_num`, `yukleme_turu` (`DEVIR`/transfer, `GURSEL_E_YUKLEME`/e-wallet, `NAKIT`/cash, `KREDI_KARTI`, `BKM`, `Kontur`, `Tosla`, `Kampanyali_Tutar_Yukleme`), `islem_kodu`, `dogrulama_kodu`, `referens_numarasi`, `gelen_yanit_metni`, `islem_tutari`, `fatura_tablosu` | Ledger of balance credits. **Insert-only** - no card balance field is ever updated; balance is computed by summing this collection (see §3.2). |
| `KartCekimHareketleri` (fare deductions) (yonetim) | `kart_cekim_hareketleri` | `is_numarasi` (job/trip id FK), `kisi_adi/soyadi`, `kart_num`, `cekim_turu`, `saatli_ulasim_listesi`, `hareket_tarihi/saati`, `plaka`, `sofor`, `islem_xy` (GPS), `islem_tutari`, `indirim_orani`, `hesaplanan_tutar`, `iade_durumu` | Ledger of balance debits - the core fare-deduction engine's output table; also insert-only. |
| `KartIadeHareketleri` (refunds) (yonetim) | `kart_iade_hareketleri` | `kart_num`, `iade_turu` (`NAKITODEME`, `KKIODEME`, `KARTAYUKLEME`, `BKMIADE`, `TOSLAIADE`), `islem_kodu`, `islem_tutari`, `aciklama` | Ledger of refunds/reversals; also feeds the balance calculation. |
| `KartSatislari` (yonetim) | `kart_satislari` | `kart_num`, `kart_fiyati` | Card-sale price records. |
| `Sanalposlar` (virtual POS config) (yonetim) | `sanalposlar` | `sanalpos_adi`, `sanalpos_durumu`, `yetkili_kurumlar`, generic `bilgi_1`..`bilgi_10` fields, `basarili/basarisiz_baglanti_adresi` | Bank merchant credentials stored in **10 unlabeled generic fields** (e.g. `bilgi_2`=storekey, `bilgi_3`=Name, `bilgi_4`=Password per the Akbank integration) - no semantic field names, a maintainability smell. |
| `PaymentTransaction` (yonetim) | `payment_transaction` | `user_id`, `order_id`, `vp_order_id`, `transaction_id`, `payment_method`, `card_num`, `ip`, `notification`, `order_note`, audit fields | Payment gateway call log/audit table. |
| `SeyahatRezervasyonlari` (yonetim) | `seyahat_rezervasyonlari` | `kurum`, `kisi_id`, `kimlik_numarasi`, `sul`, `binecegi_durak`, `sefer_turu/suresi`, `hareket_tarihi/saati`, `hareket_kodu`, `rezervasyon_durumu`, `bakiye_cekme_durumu`, `is_numarasi` | Seat reservations; feeds the UETDS government compliance report. |
| `Yetkiler` / `KullaniciYetkiGruplari` / `YetkiGruplari` (RBAC) (yonetim) | `yetkiler`, `kullanici_yetki_gruplari`, `yetki_gruplari` | pure string pairs (`{yetki_grubu, yetki}`, `{kullanici, yetki_grubu, yetki_durumu}`, `{yetki_grup_adi}`) | See §7. |
| `Kullanicilar` (root, `app\models`) | `kullanicilar` | `kullanici_eposta` (unique, used as username), `kullanici_sifresi` (hashed via `Yii::$app->security`), `kullanici_fotograf`, `kullanici_sifre_suresi`, `kullanici_durumu`, `kvkk_onay` (GDPR/KVKK consent), `auth_key`, `access_token` | System user / `IdentityInterface` implementation (Mongo AR). Note: a **second, unused** `models/User.php` (stock Yii2-scaffold demo user with hardcoded `admin`/`demo` passwords and `100-token`/`101-token` access tokens) still exists in the codebase but is not the configured `identityClass` (dead code, see §8). |
| `TicketBuyers` (ticket_sales) | `ticket_buyers` | `name`, `surname`, `id_number`, `email` (unique), `phone` (unique), `i_number`, `i_name`, `verified`, `created_at` | Guest ticket-purchase customer. |
| `Tickets` (ticket_sales) | `tickets` | `customer_id`, `route`, `order_id`, `auth_code`, `hostlog_key`, `client_id`, `response_txt`, `price`, `qr_code`, `status` (`active/used/cancelled/expired`), `expires_at` (24h), `created_at` | Purchased ticket/QR. |
| `QrCode` (odeme) | `qr_codes` | `person_id`, `code`, `card_num`, `created_at`, `expires_at` (5 min) | Cardless-passage QR token; `01TCK`-prefixed (odeme) vs `02TCK`-prefixed (ticket_sales) codes interoperate via `validateQrCode()`. |
| `WebServisKullanicilari` / `WebServisBaglantiListesi` / `GuvenliIpListesi` (services) | `web_servis_kullanicilari`, `web_servis_baglanti_listesi`, `guvenli_ip_listesi` | webservice accounts (**plaintext password field/comparison**), connection tokens with expiry/rate-limit, IP whitelist | See §7/§8. |

### 3.2 Balance model (important for migration planning)

**Card balance is not a stored field anywhere.** It is computed on demand as:

```
balance = SUM(kart_yukleme_hareketleri.islem_tutari for kart_num)
        − SUM(kart_cekim_hareketleri.hesaplanan_tutar for kart_num)
        + net(kart_iade_hareketleri for kart_num)
```

via MongoDB aggregation pipelines (e.g. `KartCekimHareketleri::toplamCekilenTutar()`, `KartYuklemeHareketleri::toplamYuklenenTutar()`, both in `modules/yonetim/models/`). This is an **event-sourced ledger design** - every top-up, deduction, and refund is an append-only document; nothing is ever updated in place to reflect "current balance". This is a valuable design property to preserve (or deliberately improve with a materialized/cached balance + ledger) in the AYSTEK Core rewrite, since it already matches an event-sourcing-friendly pattern.

There is **no visible transaction/locking mechanism** around the read-balance-then-decide-then-insert-deduction sequence in `KartCekimHareketleri::bakiyeIslemcisi()` (`modules/yonetim/models/KartCekimHareketleri.php:491-524`) - concurrent taps on the same card (e.g., from two devices, or offline-batch-replay racing a live tap) could both pass the balance-floor check before either insert lands, risking an overdrawn balance. Worth explicit handling (optimistic concurrency, or a Mongo transaction) in the rewrite.

---

## 4. Business logic & workflows

### 4.1 Fare deduction ("tap card to ride") - the core workflow

Entry point: `KartCekimHareketleri::bakiyeCekmeIslemi($cihazVerileri, $rezervasyonKontrolu = true)` (`modules/yonetim/models/KartCekimHareketleri.php:208-488`), invoked from the vehicle-reader-facing webservice (`modules/services/controllers/KartbakiyeislemleriController.php` - `actionBakiyecekmeislemi` and siblings) and from the admin kiosk UI (`yonetim/KartcekimhareketleriController::actionBakiyecek`).

Steps:
1. Resolve the scheduled-trip ("SUL") the reader is reporting for, and its 3-tier fare (`UlasimUcretleri::ulasimUcreti`).
2. Resolve the card via `OkunanKartlar::kartBilgisiOlustur()`. If the card is unknown, attempt to auto-register it by querying the institution's external webservice (`Kurumlar::kurumServiceBilgisi()`) and, if found, kick off a card-transfer (`KartDevirHareketleri::kartDevrineBasla`).
3. **Reservation short-circuit**: if the institution has `rezervasyonda_bakiye_cek` enabled, check whether the rider already holds a matching seat reservation (`SeyahatRezervasyonlari`) for this vehicle-job/date/time; if so, **skip the fare deduction entirely** (fare is considered pre-paid via the reservation).
4. Verify the card's institution matches the route's institution (else "KAPALI KART").
5. Branch on card type: `NORMAL_KART` (RFID card) goes through discount calculation (see below); `KAGIT_BILET`/`KREDILI_KART`/`KUPON` use a fixed price tier from `UlasimUcretleri` with no discount.
6. **Discount calculation** (`enUygunIndirimOrani()` / `zamanAraligiKontrolu()`, lines 526-606): resolves the rider's currently-valid person groups (`KisiGrubuGecerlilikSureleri::gecerliKisiGruplari`), looks up applicable group-route discount sets (`KisiGrupGuzergahIndirimi`), and for each candidate discount checks (a) a daily/weekly/subscription usage-count limit against past `kart_cekim_hareketleri`/`kart_iade_hareketleri` records, and (b) an allowed time-of-day window (`saat_araliklari`) with a timeout (`zaman_asimi`) since the rider's last matching ride, picking the highest applicable discount rate.
7. **Balance-floor check** (`bakiyeSinirKontrolu()`, line 608): `kalanBakiye = bakiye − (ucret − ucret*indirimOrani/100)` must be ≥ the institution's configured minimum balance (`kart_bakiye_alt_siniri`); if not, deduction is refused with "Kart Bakiyesi Yetersiz!" (insufficient balance) and the caller is told to retry as a credit-card payment (`krediKartiIslemiDene`).
8. If allowed, `bakiyeCek()` inserts the deduction record.

A near-duplicate pipeline, `cevrimDisiBakiyeCekmeIslemi()`, handles **offline reader replay** (devices buffer taps while offline and upload a batch later via `CevrimdisikarthareketleriController::actionCdisibakiyeleriaktar`), skipping the reservation-conflict step.

Numerous undocumented Turkish-language "error codes" are embedded as literal strings for support triage (e.g. `58EU001`..`58EU7`, `modules/yonetim/models/KartCekimHareketleri.php` throughout `bakiyeCekmeIslemi()`).

### 4.2 Balance top-up workflow

Entry points: admin manual entry (`yonetim/KartyuklemehareketleriController::actionCreate`), online payment gateway callback (`odeme/DefaultController::actionAkbanksanalpos`/`actionYapikrediodemeonay` → `KartYuklemeHareketleri::gurselEYukleme()`), institution-driven bulk top-up (`services/EntegrebakiyeislemleriController::actionBakiyeyukle`, `yonetim/VeriaktarimController::actionOzyegintutaryukleme`/`actionSabancitutaryukleme`).

Validation rules (`modules/yonetim/models/KartYuklemeHareketleri.php:108-134`):
- Card must have a current owner (`kartSahibiKontrolu`) - no anonymous top-ups.
- "Kontur" top-up type is restricted to certain card types (`konturYuklemeKontrolu`).
- Top-up amount is capped at the institution's configured maximum (`Kurumlar::kurumDNAzamiYuklemeLimiti`), **unless** the acting user holds the `limitsiz_bakiye_yukleme` permission (`islemTutariAzamiLimitKontrolu`).

`gurselEYukleme()` (e-wallet/online top-up, lines 427-468): records the load, and if the institution has `kampanyali_tutar_yukleme` (bonus campaigns) enabled, looks up a matching bonus tier (`YuklemeKampanyalari::getKampanyaByKurumAndTutar`) and **automatically credits an additional bonus/gift-point top-up**. On save failure, emails the admin the raw model errors and then calls `print_r($model->getErrors());exit();` (`modules/yonetim/models/KartYuklemeHareketleri.php:464`) - i.e., a failed top-up terminates the request with a raw PHP dump rather than a user-facing error page.

### 4.3 Online card payment (bank virtual POS)

Two 3D-Secure bank integrations exist, **duplicated** between `modules/odeme/controllers/DefaultController.php` and `modules/ticket_sales/controllers/PaymentController.php` (separate implementations of essentially the same POSNET protocol):

- **Akbank** (`odeme/DefaultController::akbankSanalPos()` + `actionAkbanksanalpos()`): builds a 3D-secure POST form with a SHA1 hash (`hashstr = clientId+oid+amount+okUrl+failUrl+rnd+storekey`, bank `storekey` pulled from the `Sanalposlar` config collection). The callback handler recomputes and compares the bank's hash but **does not abort processing on a mismatch** - it only appends a warning message and continues to the provizyon (authorization) step regardless (`modules/odeme/controllers/DefaultController.php:519-546`). Builds a raw `CC5Request` XML and posts it via cURL with **TLS verification disabled** (`CURLOPT_SSL_VERIFYPEER/HOST = false`), then parses the bank's XML response with manual `strpos` substring extraction rather than an XML parser.
- **Yapı Kredi** (POSNET OOS): `odeme/DefaultController::actionYapikredi()`/`actionYapikrediodemeonay()` and the parallel `ticket_sales/PaymentController::actionInitiatePayment()`/`actionPaymentResponse()`. Two-step OOS resolve→confirm flow (`oosResolveMerchantData` then `oosTranData`), also with TLS verification disabled on the outbound cURL calls.
- **BKM Express**: referenced in access-control rules and views, but `actionBkmekspress()` is an unimplemented stub (commented-out SDK calls only) - **not a working integration**.

On successful payment, both flows credit the card balance (`gurselEYukleme`) or a subscription (`AbonelikYuklemeleri::yuklemeYap`), and record an invoice (`FaturaliYuklemeler`).

### 4.4 Reservation workflow

`odeme/RezervasyonController` and `yonetim/RezervasyonisatamasiController` handle seat booking against a SUL's capacity; a waitlist (`YedekRezervasyonlar`) is used when full. If the institution's `rezervasyonda_bakiye_cek` flag is set, reservation creation itself triggers an immediate fare deduction via the same internal services endpoint used for card taps (again over cURL with the hardcoded password, value redacted, and disabled TLS - see §8). Reservation data also feeds the UETDS government passenger-manifest report.

### 4.5 Reports

Reports are generated either by re-scanning the ledger collections in a manual "rebuild" batch action (`KartHareketRaporu::raporuYenidenOlustur()`, `AracHareketRaporu`, invoked from controller actions like `actionRaporuyenidenolustur`/`actionKhrolustur` when a discrepancy is spotted), or computed live in PHP on each request by joining collections in memory (`BorcluKartlarRaporu`, `CihazIletisimHareketleri`, `HareketToplamRaporu`, `KartHareketOzeti`) and wrapped in an `ArrayDataProvider` for the kartik grid views. There is no scheduled job that keeps these report tables in sync automatically (see §6) - staleness is corrected manually by an admin clicking a "rebuild" action.

---

## 5. External integrations

| Integration | Where | Details |
|---|---|---|
| **Akbank virtual POS** | `modules/odeme/controllers/DefaultController.php` | 3D Secure, `CC5Request` XML API, SHA1 request hash, credentials in `Sanalposlar` Mongo config collection (`bilgi_2`=storekey, `bilgi_3`/`bilgi_4`=Name/Password) |
| **Yapı Kredi (POSNET) virtual POS** | `modules/odeme/controllers/DefaultController.php` AND `modules/ticket_sales/controllers/PaymentController.php` (duplicated) | OOS 3D Secure two-step resolve/confirm flow |
| **BKM Express** | `web/bkmekspress/` (vendored SDK, not Composer-managed), `assets/BKMAsset.php`, `odeme/DefaultController::actionBkmekspress()` | **Not implemented** - stub only |
| **SAP "Teacher/Student" SOAP webservice** | `modules/services/controllers/SapentegrasyonislemleriController.php`, `config/params.php:16-17` (`gurselTSSoapServiceUrlTest`/`...Canli`) | Native PHP `SoapClient` against `TeacherStudentService.asmx`; pushes per-vehicle/route earnings (`EDolumGuncelle`), pulls route/shift/capacity pricing (`EDolumFiyatGuncelleme`); hardcoded SOAP credentials (username/password pair, values redacted here) |
| **University partner systems** | `modules/yonetim/models/{SabanciuniServiceKisiListesi,OzyeginuniServiceKisiListesi,IstinyeuniServiceKisiListesi,KocuniServiceKisiListesi}` and matching `modules/yonetim/controllers/*` | Per-institution REST/SOAP integrations to resolve card↔person identity and sync rosters (Sabancı, Özyeğin, Koç, Işık, Istinye universities); dispatch is a hardcoded name-based `switch` in `Kurumlar::kurumServiceBilgisi()`. Large amounts of commented-out ad hoc cURL test code for these integrations remain in `controllers/SiteController.php` (see §8) |
| **Google Maps / Directions** | `assets/TrackingAsset.php`, `assets/LocalateAsset.php`, `modules/services/models/Google.php`, `odeme/DefaultController::actionDirections/actionHaritadaizle` | Live-vehicle-tracking map (multi-vehicle marker clustering), route/stop map editor, Directions API proxy - all with a **hardcoded Google API key** (see §8) |
| **Trio Mobil (fleet telematics/GPS)** | `modules/services/models/Triomobil.php`, `odeme/DefaultController::actionLocations` | Third-party vehicle GPS/location provider; token-based REST API, hardcoded default credentials |
| **SMS gateway ("SMS Origin")** | `modules/services/models/BilgiGonderimi.php` | Hand-rolled XML POST to `processor.smsorigin.com`, hardcoded gateway credentials |
| **QR code rendering (ticket_sales)** | `modules/ticket_sales/controllers/TicketController.php` | Uses the **public third-party API** `api.qrserver.com` to render ticket QR images (ticket content sent to an external service) rather than generating QR images locally |
| **reCAPTCHA (Invisible)** | `config/web.php` | Google Invisible reCAPTCHA, site key + secret hardcoded in config |
| **Mail** | `config/web.php` (Symfony Mailer via Office365 SMTP) | Only stock layout templates exist under `mail/layouts/` (`html.php`, `text.php`) - **no per-notification message templates**; all emails are composed inline (subject/body strings built ad hoc in controller/model code, e.g. admin-alert emails on save failures) |

---

## 6. Scheduled jobs / commands

`commands/HelloController.php` is the **unmodified stock Yii2-basic scaffold example command** (`actionIndex($message = 'hello world')` just echoes the message). **There are no custom console commands, cron jobs, or scheduled maintenance tasks anywhere in this codebase.** Any recurring reconciliation work observed in the code (report "rebuild" actions, invoice backfills, offline-transaction replay, institution data-sync) is triggered **manually by an admin clicking a button in the back office**, not by a scheduler. This is a gap worth deciding on explicitly for the rewrite (e.g., which of these manual "rebuild" flows should become real background jobs in AYSTEK Core).

---

## 7. Auth & permissions

Two separate identity systems coexist:

1. **System/admin users** (`app\models\Kullanicilar`, Mongo collection `kullanicilar`) - the configured `identityClass` in `config/web.php`. Standard Yii2 `IdentityInterface` (`findIdentity`, `findIdentityByAccessToken`, `validateAuthKey`), password hashing via `Yii::$app->security`. A second, **unused** stock-scaffold `models/User.php` with hardcoded demo credentials still ships in the codebase (dead code, not wired to `identityClass` - see §8).
2. **Institution customers/riders** log in via `SiteController::actionMusterigirisi()` against the same `Kullanicilar` collection but scoped by institution registration number (`kurumSisNum`), with a phone/email OTP-based self-registration flow (`actionMusterionayi`/`actionMusterikayit`) that provisions a new `Kullanicilar` + `Kisiler` record. The `ticket_sales` module implements a **third, independent** auth mechanism (plain PHP session + OTP, no relation to `Kullanicilar`).

**RBAC is entirely custom, not the Yii2 RBAC component (`yii\rbac`).** Three flat Mongo collections model it:
- `yetkiler`: `{yetki_grubu, yetki}` - pairs a permission-group name with a permission-key string (free text, no enum/registry).
- `kullanici_yetki_gruplari`: `{kullanici (email), yetki_grubu, yetki_durumu}` - assigns a user to one or more permission groups (`MUSTERI_YETKISI = "Kişi"` is the default group auto-assigned on customer self-registration).
- `yetki_gruplari`: `{yetki_grup_adi}` - the list of group names (pure lookup).

The single authorization primitive, `Yetkiler::getYetkiGruplari($yetkiAdi)`, is called from **every controller's `behaviors()` `AccessControl` rule** and from many action bodies directly, doing a fresh Mongo lookup of the current user's groups and a match against the `yetkiler` collection **on every single request/action, with no caching**. Permission keys are untyped strings (e.g. `kisi_sistem_girisi`, `yonetim_modulu`, `raporlar_modulu`, `limitsiz_bakiye_yukleme`, `sr_uetds_liste_goruntuleme`, `hatali_is_duzeltme_*`, `kullanici_okuma/ekleme/guncelleme/silme`) matched by convention, not by a central enum/constant registry.

Webservice-facing auth (`modules/services`) uses separate, weaker mechanisms - see §2.4/§8 (plaintext password comparison for `WebServisKullanicilari`, vehicle-plate+password fallback, IP whitelist that only bypasses rate-limiting rather than gating access).

CSRF protection is globally enabled by default but explicitly disabled on the payment callback controller and on almost all `modules/services` controllers (expected for machine-to-machine endpoints, but each such endpoint needs its own auth to compensate - see §8 for gaps).

---

## 8. Known issues / tech debt

### 8.1 Hardcoded credentials and secrets (highest priority to remediate before/during migration)

| Secret | Location |
|---|---|
| SMTP password (live) | `config/web.php:43` - literal SMTP DSN with plaintext password (value redacted here; see file) |
| SMTP password (dead, commented) | `config/web.php:42` |
| reCAPTCHA site key + secret | `config/web.php:105-106` |
| Cookie validation key | `config/web.php:20` |
| MongoDB test + production credentials | `config/params.php:11-15`, duplicated (commented) in `config/mongodb.php:6-10` |
| SAP SOAP credentials (username/password pair, plaintext, value redacted here) | `modules/services/controllers/SapentegrasyonislemleriController.php:156-158, 316-318` |
| Internal service-to-service creds (username/password pair, plaintext, value redacted here) | `modules/services/controllers/SapentegrasyonislemleriController.php:374` |
| Internal balance-pull "password" (hardcoded low-entropy shared secret, value redacted here) | `modules/odeme/controllers/KartsizgecisController.php:119`, `modules/odeme/controllers/RezervasyonController.php:227` |
| SAP capacity-pull creds (username/password pair, plaintext, value redacted here) | `modules/yonetim/controllers/GuzergahvardiyakapasitelistesiController.php` (`actionVerileriguncelleme`) |
| Sabancı university integration creds (username/password pair, plaintext, value redacted here) | `modules/yonetim/models/SabanciuniServiceKisiListesi.php:172` |
| DES-EDE-CBC key + IV for Sabancı QR decrypt (crypto key material, values redacted here) | `modules/services/controllers/KartbakiyeislemleriController.php:320,322` and duplicated in a debug harness `modules/yonetim/controllers/BaglantilarController.php` (`actionQrtest`) |
| SMS gateway creds (username/password pair, plaintext, value redacted here) | `modules/services/models/BilgiGonderimi.php:28-29` |
| Google Directions API key (value redacted here) | `modules/services/models/Google.php:21` |
| Google Maps JS API key (value redacted here) | `assets/TrackingAsset.php` (embedded in a script URL) |
| Trio Mobil GPS creds (multiple username/password pairs, plaintext, values redacted here) | `modules/services/models/Triomobil.php:21-22,53-56` |
| Hardcoded BCC email on outbound mail | `modules/services/controllers/KartbakiyeislemleriController.php:1078` (`huseyinkoc@aysteknolojileri.com`) |
| Stock-scaffold demo credentials (dead code) | `models/User.php:17,24` (`admin`/`admin`, `demo`/`demo`, `100-token`/`101-token`) |
| Plaintext webservice user passwords (stored + compared) | `modules/services/models/WebServisKullanicilari.php` (`servisKullaniciDogrulama`) |

Personal phone numbers and staff emails from a real incident left in a large dead-comment block: `modules/services/controllers/HataliisatamasiduzeltmeleriController.php:84-145`.

### 8.2 Security smells beyond hardcoded secrets

- **TLS verification disabled** (`CURLOPT_SSL_VERIFYPEER`/`CURLOPT_SSL_VERIFYHOST = false`) on essentially every outbound cURL call across `odeme` and `services` modules - bank endpoints, internal service-to-service calls, Trio Mobil, SAP-adjacent calls. Systemic, not a one-off.
- Akbank 3D Secure callback **does not abort on hash mismatch** - continues to authorization regardless (`modules/odeme/controllers/DefaultController.php:519-546`).
- Debug mode hardcoded on in the shipped entry point: `web/index.php:3-5` defines `YII_DEBUG=true`/`YII_ENV='dev'` unconditionally (only a comment tells the deployer to remove it).
- PII written to a plaintext debug log file on every Özyeğin/Sabancı QR balance-pull call: `runtime/logs/cekimTesti.txt`, opened in truncate mode (`"w"`), containing person id/email/card number (`modules/services/controllers/KartbakiyeislemleriController.php:235-266,438-468`).
- Inconsistent CORS policy within the same controller: per-action allow-list (`http://localhost:8096` only) vs. a shared header helper emitting `Access-Control-Allow-Origin: *` (`modules/services/controllers/KarthareketleriController.php`).
- IP whitelist (`GuvenliIpListesi`) only bypasses per-IP **rate limiting** for webservice token issuance, not an access gate on the endpoints themselves - worth confirming this is the intended security model, not an oversight.
- Vehicle-plate password check is fully bypassed for any plate starting with `'P'` (reservation-generated pseudo-plates) - `modules/services/controllers/KartbakiyeislemleriController.php` (~line 93-97).

### 8.3 Debug/dead-code leftovers

- Live `print_r($model->getErrors()); exit();` (or `die()`) error handling instead of proper exceptions/flash messages, in numerous files, e.g. `modules/odeme/controllers/DefaultController.php:425`, `modules/yonetim/models/KartYuklemeHareketleri.php:464`, `modules/yonetim/models/FaturaliYuklemeler.php:178`, `modules/yonetim/controllers/CevrimdisikarthareketleriController.php:182`, `modules/services/controllers/KartbakiyeislemleriController.php:235` (`or die(...)`). Dozens more `print_r($model->getErrors())` without `exit()` across `modules/yonetim/models/*`.
- A live QR-decryption **debug/test harness** left reachable in production routes: `modules/yonetim/controllers/BaglantilarController.php::actionQrtest()`.
- One-off maintenance script wired into a live controller: `modules/yonetim/controllers/FaturaliyuklemelerController.php::actionTarihduzeltme()` (hardcoded 2020-01-01..2021-01-20 date range, `exit()`s).
- Dead maintenance code in a model: `Kartlar::kartKapatmaIslemi()` contains a hardcoded array of ~54 literal card numbers and references an undefined `$arrayNin` variable (`modules/yonetim/models/Kartlar.php:527`).
- Large blocks of commented-out ad hoc integration-test cURL code (with hardcoded creds/keys) left in `controllers/SiteController.php` (roughly lines 79-193, `actionIndex()`).
- `modules/odeme/models/QrCode.php` imports `Endroid\QrCode\QrCode` but never calls it; the package isn't even in `composer.json` - unfinished/dead code.
- `config/db.php`/`config/test_db.php` still point at the never-used Yii2-scaffold default `mysql:host=localhost;dbname=yii2basic` - misleading leftover suggesting a MySQL dependency that (for business data) doesn't exist.
- A second, unused `models/User.php` (stock Yii2 demo identity class) ships alongside the real `Kullanicilar` - see §8.1.
- `commands/HelloController.php` is the unmodified Yii2 scaffold example (§6).
- `mail/` contains only the two stock layout files, no actual message templates (§5).

### 8.4 Likely functional bugs found while reading

- `GunlukKasaDekontlari::beforeSave()` uses `if (self::NAKIT)` / `if (self::KAGIT_BILET)` - since these are non-empty string constants, both conditions are always truthy, so the `NAKIT` branch always runs and the `KAGIT_BILET` branch is unreachable, instead of comparing against `$this->dekont_turu` as clearly intended.
- `UlasimUcretleri::uyusanTarihliUlasimUcretiVarMi()` filters on a misspelled field `gecerlilik_bit_tarihi` instead of `gecerlilik_bitis_tarihi`, likely silently disabling its own overlapping-fare-window guard.
- `modules/raporlar/controllers/MusteriarackazancraporuController.php` (`actionCreate` line ~112, `actionUpdate` line ~308) references `$saatlerObj->millisecondToDate(...)` without ever instantiating `$saatlerObj` in that method scope - an undefined-variable bug on that code path.
- `modules/raporlar/controllers/HarekettoplamraporuController.php` builds its filter model directly from `$_REQUEST` rather than via `$model->load()` - bypasses Yii2's normal mass-assignment/safe-attribute protections (not classic SQL injection since the backend is Mongo, but an unvalidated-input pattern).
- `AracbazlisultopkazanclariController`/`MusteriarackazancraporuController::actionCreatelist` parse a raw `;`/`,`-delimited string POST body into fields by manual `explode()`/index access - brittle, unvalidated "protocol" rather than structured input.

### 8.5 Architecture-level debt to weigh in the rewrite

- No relational integrity: cross-collection "foreign keys" are denormalized strings (card number, person email, vehicle plate, route/SUL name, group name); at least 6 places in the codebase manually cascade a rename across 5-8 collections (email change in `Kisiler::kisiEpostaDegis()`, plate rename in `AraclarController::actionUpdate`, SUL rename in `SaatliulasimlistesiController::actionUpdate`, card-number change, group-name rename, etc.) - a missed cascade site would silently corrupt data.
- RBAC permission checks re-query Mongo on every action dispatch with no caching (§7).
- No transactional guard around the fare-deduction read-balance→decide→insert sequence (§3.2) - a race-condition risk under concurrent taps.
- Two independent, duplicated bank-payment integrations (Akbank/Yapı Kredi logic exists once in `odeme` and again, separately, in `ticket_sales`) - consolidation candidate.
- Heavy reliance on raw `Yii::$app->mongodb->createCommand()->aggregate()` calls scattered across 70+ call sites in `modules/yonetim/models/*` (bypassing the ActiveRecord layer) means most real query/business logic lives outside any shared, testable query-builder abstraction.
- `firebase/php-jwt` is a declared but entirely unused dependency - if a stateless/mobile-friendly API is a goal for AYSTEK Core, no JWT wiring exists yet to carry over.
- `kartik-v/*` Composer dependencies are pinned to `dev-master` (unpinned/unreproducible builds) throughout `composer.json`.

---

## 9. Reporting / export features

- **Excel export**: exclusively via `kartik-v/yii2-export`'s `ExportMenu` widget (itself wrapping PhpSpreadsheet under the hood) in `raporlar` module views: `musteriarackazancraporu/index.php`, `karthareketraporu/index.php`, `karthareketraporu/borclu_kart_listesi.php`, `karthareketraporu/kisi_sul_kullanim_bilgisi.php`, `arachareketraporu/sul_arac_toplam_kazanc.php`, `arachareketraporu/doluluk.php`, `arachareketraporu/index.php`, `arachareketraporu/arac_is_form.php`. Two views (`arachareketraporu/sul_doluluk.php`, `arachareketraporu/ortalama_sul_doluluk.php`) have the same export widget **present but commented out** (feature disabled there). `karthareketraporu/kart_hareket_ozeti.php` imports `ExportMenu` but never calls it (unused import).
- **PDF export is explicitly disabled everywhere** (`ExportMenu::FORMAT_PDF => false` at every call site) - **no `mpdf` or `PhpOffice`/PhpSpreadsheet-direct dependency exists in `composer.json`**, confirming PDF generation is not a live feature of this app despite the task brief's expectation of checking for one.
- `codemix/yii2-excelexport ^2.7` is also declared in `composer.json` but no direct usage was found during this survey scope (may be a superseded/legacy dependency from before the `kartik-v/yii2-export` adoption, or used in a view not sampled).

---

## Appendix: file/directory reference index

- Config: `config/web.php`, `config/params.php`, `config/db.php`, `config/mongodb.php`, `config/console.php`, `config/test.php`, `config/test_db.php`
- Root auth: `controllers/SiteController.php`, `models/{Kullanicilar,LoginForm,MusteriLoginForm,User}.php`
- Core admin module: `modules/yonetim/` (85 controllers, 85+ models under `models/`)
- Payment/rider module: `modules/odeme/`
- Reports module: `modules/raporlar/`
- Webservice API module: `modules/services/`
- Guest ticket-sales module: `modules/ticket_sales/`
- Console commands: `commands/HelloController.php` (stock scaffold, no real jobs)
- Mail layouts (no templates): `mail/layouts/{html,text}.php`
- Vendored/unused BKM Express SDK: `web/bkmekspress/`
