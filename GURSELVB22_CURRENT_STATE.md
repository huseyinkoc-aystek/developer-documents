# GurselVB22 ("Veli" Parent Portal) - Current State Notes

**Date:** 2026-07-09
**Project path:** `/workspace/php/GurselVB22`
**Purpose:** Yii2 PHP web app used by parents ("Veli") of school students (K-12, not universities) to log in with the student's national ID number, enroll/register the child for the shuttle/bus service, request withdrawal from the service, view sibling ("kardeş") discounts, sign digital service contracts, and pay for the service by credit card (multiple bank 3D-Secure gateways) or by manually-reconciled bank transfer ("havale").

This document only covers GurselVB22. A companion product, **GurselMobile**, is surveyed separately by another agent and is not read here in depth — but GurselVB22 is a thin client over GurselMobile's HTTP API (see "External Integrations" below), so understanding that boundary matters a lot for migration planning.

**Cross-check:** independently re-verified against source in `LEGACY_GURSEL_APPS_CODEX_REVIEW_2026-07-09.md`, which confirms the "no local domain DB, thin orchestration shell over SAP/ATS/KT" finding and adds the AYSTEK Core bounded-context mapping / migration sequence / do-not-migrate list. Read that file alongside this one before planning the migration.

---

## 1. Overview

GurselVB22 is a **Yii 2 Basic Project Template** app (`yiisoft/yii2-app-basic`, framework `~2.0.13`, PHP `>=5.4`). It is thin: only **2 controllers** (`SiteController.php`, ~182 KB / ~4,020 lines, and `HavaleController.php`), **15 model files** (mostly plain `yii\base\Model` form objects, not `ActiveRecord`), and a `views/` tree of ~40 view files. There is no admin backend in this app; admin/back-office functions apparently live elsewhere (SAP-based ERP + the `ats.gurseltur.com.tr` service, + possibly GurselMobile).

Key composer dependencies (`composer.json`):
- `yiisoft/yii2: ~2.0.13`, `yiisoft/yii2-bootstrap`, `yiisoft/yii2-swiftmailer` (email), `yiisoft/yii2-mongodb: ^2.1` (declared but **no MongoDB usage was found** anywhere in `models/`, `controllers/`, or `config/` — looks like unused/vestigial dependency).
- Large set of **kartik-v** widget packages (`yii2-widgets`, `yii2-grid`, `yii2-dynagrid`, `yii2-detail-view`, `yii2-dropdown-x`, `yii2-nav-x`, `yii2-context-menu`, `yii2-slider`, `yii2-field-range`, `yii2-sortable`, `yii2-sortable-input`, `yii2-social`, `yii2-checkbox-x`, `yii2-date-range`, `yii2-datecontrol`, `yii2-money`, `yii2-popover-x`, `yii2-tabs-x`, `yii2-editable`, `yii2-mpdf`, `yii2-export`, `yii2-tree-manager`, `yii2-ipinfo`, `yii2-dialog`, `yii2-widget-datepicker`, `yii2-builder`) — heavier widget footprint than the actual number of views would suggest is used; likely leftover from a shared boilerplate/other Gursel apps.
- `codemix/yii2-excelexport`, `pctux/yii2-invisible-recaptcha`.
- `kartik-v/yii2-social` includes Facebook-login-capable widgets as a dependency, but **no Facebook SDK usage, no `facebook` string, and no social-login controller code were found** anywhere in `controllers/` or `models/` — the dependency appears unused in this app's actual code paths.
- Dev-only: `yii2-codeception`, `yii2-debug`, `yii2-gii`, `yii2-faker` (fzaninotto faker, dev-only as expected).
- `auth.json` at repo root contains a **live-looking GitHub OAuth token** (`github-oauth` for `github.com`) committed in plaintext — this is a secret-in-repo finding (see Tech Debt).

The app has **no local persistent database of its own for domain data**. `config/db.php` configures a MySQL connection (`huseyinkoc.com.tr`, hardcoded credentials — see Tech Debt), but no `ActiveRecord` model or `Yii::$app->db` query was found anywhere in `models/`, `controllers/`, or `commands/`. All student/parent/enrollment/payment data is fetched and written through:
1. A **SOAP web service** (`sapservice.gurseltur.com.tr/WebServices/TeacherStudentService.asmx`) — the actual system of record, evidently backed by an SAP-adjacent ERP ("TeacherStudentService", "PlateToSupplierService").
2. A **REST/HTTP API on `ats.gurseltur.com.tr/web/services/...`** — this is almost certainly the deployed identity of **GurselMobile** (see section 5).
3. A separate **REST/HTTP API on `kt.gurseltur.com.tr/web/services/kartbakiyeislemleri/...`** used purely for sending email/SMS notifications (name suggests "card balance operations" service — likely another Gursel microservice, distinct from GurselMobile/ats).

So functionally, GurselVB22 is a **UI/orchestration shell** — session-driven, mostly stateless itself — that authenticates parents against the SOAP service, drives multi-step business logic in PHP, and pushes state changes back out to `ats.gurseltur.com.tr` and `kt.gurseltur.com.tr` over plain HTTP with credentials in the query string.

---

## 2. Domain Concepts

- **Veli (parent/guardian)**: not a persisted local entity. The "logged in user" is a `yii\web\IdentityInterface` implementation (`app\models\User`) whose identity data is fetched live from SOAP (`OgrenciBilgileri::ogrenciBilgisi`) keyed by student TC (national ID) number + `projeKodu` (school/project code), and cached only in PHP session (`$session["ogrBilgi_info"]`). There is no `users` table.
- **Öğrenci (student/child)**: represented by SOAP response objects (`TeacherStudentInfoResult`), with fields such as `TeacherStudentIdentityNo` (TC no), `TeacherStudentNameSurnameText`, `StudentNo` (SAP student number), `Project` (school code), `TeacherStudentPeriod` (academic term/donem), `AllCare`/`PaidCare`/`RevertTeacherStudentNumber` (balance fields used to compute amount due), `RoadTeacherStudent` (assigned route/guzergah), `Plaka` (vehicle plate), `FatherEmail`/`MotherEmail`/`OtherPersonEmail`, mobile phone fields.
- **Okul/Proje (school)**: identified by a `projeKodu`/`kurum` code (e.g. `"0012020014"`), resolved to a name via `OgrenciBilgileri::projeListesi()` which now delegates to `projelistesiATS()` (curl call to `ats.gurseltur.com.tr .../odetkinkurumlar`); an older SAP SOAP-based project list (`ProjeListInfo`) is present but commented out.
- **Servis/Route enrollment**: pre-registration ("Servis Ön Kayıt", `ServisOnKayitForm`) is a **lead-capture form**, not an actual enrollment transaction — it emails the school's contact list and pushes the form data to `ats.gurseltur.com.tr/.../servis-on-kaydi`. It does not go through SAP and does not create a login-capable student record itself.
- **Servisten Ayrılma (withdrawal)**: `ServistenAyrilmaTalebiForm` — a request form (reason, start/leave dates) that similarly emails staff and does not appear to update SAP/enrollment state directly; it is a "request", presumably actioned by staff off-system.
- **Kardeş İndirimi (sibling discount)**: `KardesIndirimi` model — parent enters one or more siblings' TC numbers plus an uploaded proof document; validated with a Turkish TC checksum algorithm, submitted to `ats.gurseltur.com.tr/.../kardesindirimitalebi`, and the uploaded file is additionally pushed via **FTP** to `10.34.18.12` (plaintext credentials, see Tech Debt) in addition to being saved locally under `web/customerFiles/`.
- **Sözleşme (service contract) flows**: `VSozlesmesi` model manages "soz_1"/"soz_2" (customer-service and cancellation/refund contracts) — checked via `sozlesmeOnayKontrolu()`, rendered with placeholder substitution (`[FIYAT_BILGISI]`, `[GUZERGAH_BILGISI]`, `[TC_KIMLIK_BILGISI]`, `[OGRENCI_AD_SOYAD]`, `[OGRENCI_ADRES]`, etc.), and confirmed back to `ats.gurseltur.com.tr/.../sozlesmeonaykaydi`. Login and several payment-flow actions are gated on the parent having open, unsigned contracts (`actionIndex`, `actionKardesindirimi`, `actionServistenayrilmatalebi`, `actionAkbanktaksit`, etc. all call `tcKayitKontrolu()`/`sozlesmeOnayKontrolu()` first).
- **TC Bilgi Alma (parent identity capture)**: `TCBilgiAlmaForm` — collects mother/father TC, name, birth year (at least one parent required) plus student blood group; posts to `ats.gurseltur.com.tr/.../velitckayit`. Includes dead/unused validators (`anneKimlikDogrulama`/`babaKimlikDogrulama`) that call the real Turkish government identity-verification SOAP service `tckimlik.nvi.gov.tr/Service/KPSPublic.asmx` — these methods exist but are **not wired into `rules()`**, so they are dead code today.
- **Ödeme (payment)**: two paths —
  - **Card payment**: bank-hosted 3D-Secure flow across **6 bank integrations**: Akbank, Garanti, İşbank (plus a newer "isbankyeni" variant), Yapı Kredi, QNB Finansbank (plus a newer "finansyeni" variant). Each has a `<bank>taksit` (installment selection) action, a `<bank>odeme3d` (submit to bank's 3D gateway) action, and a `<bank>odeme3donay` (3D callback/confirmation, hash validation, then server-to-server "Auth" XML request to the bank) action.
  - **Havale (bank transfer)**: `HavaleController::actionBildirim` + `HavaleForm` — a "notify us you've transferred money" flow; no bank API integration, no automatic reconciliation. It records a `"HAVALE"` payment intent via `OgrenciBilgileri::ogrenciOdemeSonucu(...)` back to SAP, which appears to be **manually reconciled by school/finance staff** off-system (there is no evidence of a bank statement import or matching job anywhere in this codebase).

---

## 3. Data Model

There is effectively **no ActiveRecord layer**; the term "model" here means Yii2 form/DTO objects (`yii\base\Model`) that validate input and orchestrate SOAP/HTTP calls. None declare `tableName()` or `rules()`-backed DB persistence, and no `extends ActiveRecord` class exists in the codebase (confirmed via grep).

| File | Role | Key fields | Persistence |
|---|---|---|---|
| `models/User.php` | Parent identity (`IdentityInterface`) | `tcNum`, `proje`, `ogrSapNo`, `username`, `fatherEmail`/`motherEmail`/`otherEmail`, `authKey`, `accessToken` | Session only (`$session["ogrBilgi_info"]`); source of truth is SOAP `TeacherStudentInfo` |
| `models/User_1.php` | **Dead/legacy** copy of the Yii2 scaffold `User` (hardcoded `admin`/`demo` test users, `extends \yii\base\Object` — pre-2.0.4 API) | n/a | Hardcoded array |
| `models/LoginForm.php` | Login form (TC number + school + KVKK consent) | `tcNum`, `projeKodu`, `donem`, `kvkkOnayi` | n/a (delegates to `User::findByTCNum`) |
| `models/LoginForm_1.php` | **Dead/legacy** copy of scaffold login form (username/password) | n/a | n/a |
| `models/OgrenciBilgileri.php` | Not really a "model" — a service/gateway class wrapping SOAP calls to SAP and REST calls to `ats.gurseltur.com.tr` | n/a | External systems |
| `models/ServisOnKayitForm.php` | Pre-registration lead form | ~20 fields: student name/TC/school/class/address/birth date/blood group, mother/father/other contact details, service start date, usage type (`CIFT_YON`/`TEK_YON_SABAH`/`TEK_YON_AKSAM`) | Emails + curl POST to `ats.gurseltur.com.tr/.../servis-on-kaydi` |
| `models/ServistenAyrilmaTalebiForm.php` | Withdrawal request form | `okulAdi`, `projeKodu`, `ogrenciAdSoyad`, `ogrenciKimlikNum`, `servisBasTarihi`, `servisAyrilmaTarihi`, `ayrilmaNedeni` | Email only (no `ats` call found in this model) |
| `models/KardesIndirimi.php` | Sibling discount request | `ogr_tc` (array of TC numbers), `file` (upload) | curl POST to `ats.gurseltur.com.tr/.../kardesindirimitalebi` + FTP upload of the file to `10.34.18.12` |
| `models/TCBilgiAlmaForm.php` | Parent TC/identity capture | `ogr_tc`, `ogr_dogum_tarihi`, `ogr_kan_grubu`, `proje_kodu`, `anne_*`/`baba_*` (TC, ad, soyad, dogum_yili), `kvkk` | curl POST to `ats.gurseltur.com.tr/.../velitckayit` |
| `models/VSozlesmesi.php` | Contract text + approval | `sozlesme_kodu`, `sozlesme_basligi`, `sozlesme_metni`, `sozlesme_onay_tarihi` | curl POST to `ats.gurseltur.com.tr/.../sozlesmeonaykaydi` |
| `models/HavaleForm.php` | Bank-transfer notification form | `tcNum`, `projeKodu`, `persId`, `donem`, `onBilgilendirme`, `musteriHizmetSozlesmesi`, `amount` | Delegates persistence to `OgrenciBilgileri::ogrenciOdemeSonucu(...)` (SOAP `PaymentsSet`) |
| `models/OgrenciSMSOnay.php` | SMS OTP confirmation (login second factor) | `ogrenciTel`, `smsOnayKodu` | Code generated and compared in **PHP session only**; SMS sent via curl POST to `kt.gurseltur.com.tr/.../smsgonder` |
| `models/SurucuHostesBilgileri.php` | Driver/hostess lookup by vehicle plate | n/a | SOAP `DriverAndHostessInfo`/`DriverAndHostessDetails` |
| `models/ContactForm.php` | Generic contact-us form | `name`, `email`, `subject`, `body`, `verifyCode` | Sends email via `Yii::$app->mailer` (Swiftmailer) |
| `models/OdemeModel.php` | Empty stub (`class OdemeModel extends Model {}`), unused | — | — |

**Core entities, restated for migration mapping:** Parent = `User`/SOAP `TeacherStudentInfo[0]`; Student = SOAP `TeacherStudentInfo[1]` (`StudentNo`); School = `projeKodu`/`kurum`; Route = `RoadTeacherStudent`/`guzergah`; Enrollment record = not owned by this app (lives in SAP + `ats`); Payment/transaction = SOAP `PaymentsSet` + `ats.../odemegecisi`/`odemebilgisi`; Bank transfer record = `HavaleForm` + SOAP `PaymentsSet` with `UYE_ISYERI_NO = "HAVALE"`.

---

## 4. Business Logic & Workflows

### 4.1 Login / session bootstrap (`SiteController::actionIndex`, `actionLogin`)
1. Parent submits TC number + school (`projeKodu`) + KVKK consent (`LoginForm`).
2. `LoginForm::login()` calls `User::findByTCNum()`, which first checks `OgrenciBilgileri::odemeAlmaDurumu($projeKodu)` (curl to `ats.../odemealmadurumu` — "is this school open for payments") and, if true, calls SOAP `TeacherStudentInfo`.
3. On success, a synthetic `User` identity is built (password = TC number itself — i.e., **the TC number is both username and password**) and logged in via `Yii::$app->user->login()`.
4. **SMS OTP second factor**: `actionLogin` then looks up the child's assigned driver (`SurucuHostesBilgileri`) purely to source phone numbers, generates a 5-digit code (`OgrenciSMSOnay::kodOlustur`), stores it in session, and sends it via curl to `kt.gurseltur.com.tr/.../smsgonder` to the parent's mobile numbers (mother/father/other, deduplicated). The `ogrenci_sms_onay` view collects the code; `OgrenciSMSOnay::smsKodKontrol()` compares it to the session value.
5. Once OTP is confirmed (`actionOgrencismsonay` sets `$session['smsOnay'] = 1`), `actionIndex` runs a chain of gates before showing the dashboard: (a) TC/parent-info completeness check (`tcKayitKontrolu`, redirects to `tctalebinigonder` if missing), (b) contract "soz_1"/"soz_2" approval check (redirects to `sozlesmetalebinigonder` if unsigned), each re-checked at most every 300 seconds via session timestamps (`tc_son_sorgu_zamani`, `soz1_son_sorgu_zamani`, `soz2_son_sorgu_zamani`).
6. A **hardcoded TC-number special case** exists: `if(Yii::$app->user->identity->tcNum == "68206285232") { return $this->redirect(['kardesindirimi']); }` (`SiteController.php:284`) — a one-off business rule baked directly into the controller for a specific parent.

### 4.2 Registration / withdrawal (lead-capture, not transactional)
- **Servis Ön Kayıt** (`actionServisonkaydi` → `ServisOnKayitForm::servisKayitIslemi()`): builds an HTML summary, emails school staff (hardcoded recipient list plus school's own contact email resolved via SOAP `SchoolEmailInfo`), and separately POSTs the raw form fields to `ats.gurseltur.com.tr/.../servis-on-kaydi` for logging. This is a **request**, not an enrollment — actual service assignment/route creation happens outside this app (presumably in the SAP-side ERP or GurselMobile/ats backend).
- **Servisten Ayrılma Talebi** (`actionServistenayrilmatalebi` → `ServistenAyrilmaTalebiForm::ayrilmaTalebiIslemi()`): requires an existing logged-in parent whose TC record is already complete (`tcKayitKontrolu` gate); emails the same staff list with the withdrawal request; the SAP-side unenrollment itself is presumably done manually by staff, since no SOAP write call for "remove from service" was found.

### 4.3 Kardeş İndirimi (sibling discount)
`actionKardesindirimi` gates on contract approval, then on POST validates sibling TC numbers (Turkish TC checksum, uniqueness), saves the uploaded proof file into `web/customerFiles/<parentTC>-<randomstring>.<ext>`, POSTs the claim to `ats.gurseltur.com.tr/.../kardesindirimitalebi`, and additionally FTP-uploads the same file to `10.34.18.12` using hardcoded plaintext creds (`KardesIndirimi::sendFile()`).

### 4.4 Contract signing (`sözleşme`)
`actionSozlesmetalebinigonder`/`actionSozlesmeonay` fetch a templated contract from `ats.../sozlesme`, substitute placeholders with live SOAP data (price, route, parent name, student address), render it for the parent to accept, and record acceptance back to `ats.../sozlesmeonaykaydi` including the request IP address, phone number, and the accepting parent's TC. `views/havale/index.php` also independently fires the same "soz_1"/"soz_2" acceptance via an inline jQuery `$.ajax` POST when the parent ticks a checkbox during the havale flow — i.e., contract-acceptance logging can be triggered from two different code paths.

### 4.5 Card payment (example: Akbank, representative of all 6 bank integrations)
1. `actionAkbanktaksit`: recomputes balance/discount via `OgrenciBilgileri::getCare("AKBANK")` (which itself calls `ats.../donemindirimoranlari` for bank-specific discount rates) and renders an installment-selection form.
2. `actionAkbankodeme3d`: builds a 3D-Secure init payload (`clientId`, `oid` = `studentNo + "AB" + timestamp`, `hash = base64(sha1(clientId+oid+amount+okUrl+failUrl+rnd+storekey))`), logs a `"KART GIRISI YAPILIYOR"` / `"KART DOGRULANDI"` breadcrumb to `ats.../odemegecisi` (which also embeds parent email addresses in a JSON blob for logging), then renders an auto-submitting form to the bank's hosted 3D gateway.
3. `actionAkbankodeme3donay`: bank redirects back with `mdStatus`, `HASHPARAMS`/`HASHPARAMSVAL`/`HASH`; the code recomputes the hash and compares (`if ($paramsval != $hashparamsval || $hashparam != $hash) { $hata .= "Güvenlik Uyarisi..." }`) — **but continues processing regardless of this check failing** (the `$hata` string is accumulated for display, not used to `return`/abort the transaction at this point in the code read). It then builds an XML `CC5Request` "Auth" request server-to-server to the bank using the merchant name/password from `config/params.php`, and on success presumably calls back into `ogrenciOdemeSonucu()`/SAP (not fully re-traced beyond this point given file size, but the pattern is consistent across all 6 banks based on the action list).
4. Store/Finansbank ("finansyeni"), İşbank ("isbankyeni") have newer duplicate action sets alongside older ones (`finanstaksit` vs `finansyenitaksit`, `isbankodeme3d` vs `isbankyeniodeme3d`), suggesting an in-place bank-integration migration that left the old code paths in place rather than removed.

### 4.6 Havale (bank transfer) payment
`HavaleController::actionBildirim` (see code walkthrough in section 2) computes the outstanding balance from SOAP data, applies the same contract/TC gating as the card flows, then simply records a `"HAVALE"` payment attempt via `ogrenciOdemeSonucu(..., $UYE_ISYERI_NO = "HAVALE", ...)`. There is **no bank statement reconciliation, no webhook, no scheduled job** touching this in the codebase — it is a notification mechanism only; actual reconciliation of the transfer against the bank account is implicitly a manual/off-system finance-team process.

### 4.7 Notifications
- **Email**: mix of direct Swiftmailer (`Yii::$app->mailer->compose()`, used in `ContactForm`, and commented-out in several other places) and an external "email-sending microservice" reached via curl at `kt.gurseltur.com.tr/.../epostagonder` (used in `ServisOnKayitForm`, `ServistenAyrilmaTalebiForm`, `OgrenciSMSOnay` for SMS-failure alerts, `OgrenciBilgileri` for SAP-write-failure alerts). Many direct-mailer call sites are commented out, suggesting a migration from Swiftmailer-direct to the `kt.gurseltur.com.tr` service that was done inconsistently (dead code left behind rather than removed).
- **SMS**: sent exclusively via curl to `kt.gurseltur.com.tr/.../smsgonder`; a much older direct integration to `processor.smsorigin.com` (channel code and username, values redacted here) is present but fully commented out (dead code, `models/OgrenciSMSOnay.php:96-174`).
- **Logging/monitoring email target**: `config/web.php` configures a Yii `EmailTarget` log component that emails `huseyinkoc@aysteknolojileri.com` (the AYSTEK developer doing this migration) on every `error`-level log message.

---

## 5. External Integrations

| System | Host/Endpoint | Protocol | Used for |
|---|---|---|---|
| SAP-adjacent ERP | `sapservice.gurseltur.com.tr/WebServices/TeacherStudentService.asmx` (also test/local variants in `params.php`) | SOAP/WSDL | System of record for student/parent/enrollment/route/payment data (`TeacherStudentInfo`, `PaymentsSet`, `SchoolEmailInfo`, `DriverAndHostessInfo`, `DriverAndHostessDetails`, `PesinIndirimOrani`) |
| **`ats.gurseltur.com.tr/web/services/...`** | REST-ish, `application/x-www-form-urlencoded` POST, plaintext `kullanici`/`sifre` (username/password) creds in body — two distinct credential pairs used across call sites, values redacted here | HTTP (curl, `CURLOPT_SSL_VERIFYPEER`/`VERIFYHOST` disabled everywhere) | **Strong evidence this is GurselMobile's deployed API.** Several call sites have a commented-out `localhost/GurselMobile/web/services/...` or `localhost/projects/GurselMobile/web/services/...` URL directly above the live `ats.gurseltur.com.tr` URL (e.g. `models/OgrenciBilgileri.php:246,319,352,386,418`), i.e. the same endpoint path (`veliservisodemeservisleri/...`, `sapentegrasyonislemleri/...`) was developed against a local GurselMobile checkout and then pointed at the `ats.gurseltur.com.tr` production host. Endpoints used: `odetkinkurumlar`, `kurumindirimdurumu`, `donemindirimoranlari`, `kurumsozlesmeleri`, `odemealmadurumu`, `tumdonemsozlesmeleri`, `sozlesme`, `odemebilgisi`, `odemegecisi`, `velitckayitsorgula`, `velitckayit`, `kan-grubu-durumu`, `kisi-kan-grubu-zorunlulugu`, `banka-bilgisi`, `banka-taksit-sayisi`, `sozlesmesorgula`, `sozlesmeonaykaydi`, `servis-on-kaydi`, `kardesindirimitalebi` |
| `kt.gurseltur.com.tr/web/services/kartbakiyeislemleri/...` | Same HTTP pattern as above, same hardcoded creds (values redacted here) | HTTP | Separate microservice ("card balance operations") used only for `epostagonder` (send email) and `smsgonder` (send SMS) — likely a shared Gursel-wide notification service, distinct from GurselMobile |
| Card payment gateways | Akbank (`sanalakpos.com`), Garanti (`sanalposprov.garanti.com.tr`), İşbank (`sanalpos.isbank.com.tr`), Yapı Kredi (`posnet.yapikredi.com.tr` / `3DSWebService`), QNB Finansbank old (`fbwebpos.com`) and new (`vpos.qnbfinansbank.com`); shared test gateway `entegrasyon.asseco-see.com.tr` for several banks | 3D-Secure HTML form POST + server-to-server XML (`CC5Request`) | Card payment processing; all merchant IDs/store keys/usernames/passwords for both **live and test** are hardcoded in `config/params.php` |
| Turkish government identity registry | `tckimlik.nvi.gov.tr/Service/KPSPublic.asmx` | SOAP | Present in `TCBilgiAlmaForm::anneKimlikDogrulama`/`babaKimlikDogrulama` but **not called from `rules()`** — dead/unused code today |
| FTP file store | `10.34.18.12` | FTP (`ftp_ssl_connect`) | Sibling-discount proof documents are duplicated here (`KardesIndirimi::sendFile`), hardcoded creds (values redacted here) |
| SMTP (Office 365) | `smtp.office365.com:587` | Swiftmailer/SMTP | Fallback/legacy direct email path; hardcoded creds in `config/web.php` |
| Facebook SDK | — | — | **Not used.** `kartik-v/yii2-social` is a composer dependency but no Facebook login code, no `facebook` string, and no OAuth controller action exist in this app. |
| MongoDB | — | — | **Not used.** `yiisoft/yii2-mongodb` is declared in composer but no Mongo connection/config/usage exists anywhere. |

**Conclusion on GurselMobile linkage:** GurselVB22 does not call GurselMobile via a documented internal client class or SDK — it calls a set of REST endpoints on `ats.gurseltur.com.tr` using raw `curl_init()` calls scattered across nearly every model file, with hardcoded plaintext basic-auth-style credentials in the POST body on every single call. The commented-out `localhost/GurselMobile/...` URLs strongly suggest `ats.gurseltur.com.tr` **is** (or fronts) the deployed GurselMobile instance, sharing the exact same `web/services/<controller>/<action>` URL scheme GurselMobile's own agent should be seeing on its side (e.g. `veliservisodemeservisleri`, `sapentegrasyonislemleri` controllers). This should be confirmed against the GurselMobile survey, but the evidence here is strong and consistent across 6+ files.

---

## 6. Scheduled Jobs / Commands

`commands/` contains only the stock Yii2 scaffold example: `commands/HelloController.php` (`app\commands\HelloController::actionIndex`, echoes a message — not a real business command). **No cron/scheduled console commands exist in this app.** `config/console.php` has a commented-out `yii\faker\FixtureController` mapping only. Any recurring reconciliation (e.g. for havale bank-transfer matching) is not implemented here — if it exists at all, it must live in GurselMobile, the SAP-side ERP, or a manual staff process.

---

## 7. Auth & Permissions

- **No RBAC.** There is a single logical role: an authenticated parent (`Yii::$app->user`, identity class `app\models\User`). There is no admin login, no `AccessControl` role beyond `'@'` (any authenticated user) and `'?'/'*'` (guest/any) in `SiteController::behaviors()`.
- Authentication is TC-number + school-code + SMS OTP, with **no real password** (`User::validatePassword()` compares to the TC number itself, since `password` is set to the TC number when the identity is built in `User::findByTCNum`).
- `AccessControl` in `SiteController::behaviors()` only lists a subset of actions in its `only` array (`logout`, `index`, `tctalebinigonder`, `sozlesmetalebinigonder`, `servistenayrilmatalebi`, and the various bank payment/taksit/result actions) all gated to `roles => ['@']` (must be logged in), with `tctalebinigonder`/`sozlesmetalebinigonder` also allowed for guests (`'*'`, `'?'`) since those can be reached with a `pTCNum`/`pProje` query param before full login (used from the havale flow). **Actions not listed in `only` (e.g. `servisonkaydi`, `kardesindirimi` is listed but `contact`, `changelang`, `ogrencismsonay`, `sartname`, `gizlilik`, `ogrencilistesi`, etc. are not) fall outside the `AccessControl` filter entirely** and are open to anyone by default under Yii2's `AccessControl` semantics (unlisted actions are not filtered by an `only`-scoped rule set) — this needs to be verified against Yii2's exact default-deny/allow behavior for `only`, but at minimum `actionOgrencilistesi` (dumps a SOAP-derived list, see Tech Debt) has no visible access gate of its own.
- `HavaleController` has **no `behaviors()`/`AccessControl` at all** — `actionBildirim` is reachable by anyone, gated only by whatever business-logic checks are inline (contract signature, TC lookup), not by Yii's access-control layer.
- CSRF: **disabled application-wide** — `SiteController.php:30`: `public $enableCsrfValidation = false;` (see Tech Debt).

---

## 8. Known Issues / Tech Debt

Concrete, file:line-referenced findings, ordered roughly by severity:

1. **CSRF validation is disabled for the entire main controller.** `controllers/SiteController.php:30`: `public $enableCsrfValidation = false;`. Combined with the app processing money (card payments, bank-transfer notifications) and personal data (TC numbers, blood types, addresses) via POST forms, this is a significant security gap. A commented-out `beforeAction` override right above it (lines 24-27) shows this was a deliberate, explicit choice, not an oversight — but the reasoning isn't documented.

2. **Every external HTTP integration disables TLS certificate validation.** Every single `curl_init()` call in the codebase (over 25 call sites across `models/OgrenciBilgileri.php`, `models/ServisOnKayitForm.php`, `models/ServistenAyrilmaTalebiForm.php`, `models/KardesIndirimi.php`, `models/VSozlesmesi.php`, `models/OgrenciSMSOnay.php`, `models/TCBilgiAlmaForm.php`, `controllers/SiteController.php`) sets `CURLOPT_SSL_VERIFYPEER, false` and `CURLOPT_SSL_VERIFYHOST, false`, e.g. `models/OgrenciBilgileri.php:228-229`, `:260-261`, `:297-298`, etc. This defeats HTTPS's core protection (MITM) for every call to `ats.gurseltur.com.tr`, `kt.gurseltur.com.tr`, and payment gateways, even though the URLs are `https://`.

3. **Plaintext credentials hardcoded and repeated dozens of times across the codebase**, in three main forms:
   - Application-service creds embedded directly in query strings on every call: a `kullanici=`/`sifre=` (username/password) pair appears identically in at least 15+ locations, values redacted here (e.g. `models/OgrenciBilgileri.php:221`, `:253`, `:290`, `:325`, `:354`, `:388`, `:420`, `:445`, `:473`; `models/TCBilgiAlmaForm.php:255`, `:284`, `:311`, `:373`, `:407`; `models/VSozlesmesi.php:66`, `:95`, `:137`; `controllers/SiteController.php:614`) and a second distinct username/password pair (e.g. `models/OgrenciBilgileri.php:253`, `:290`, `:325`).
   - FTP creds in `models/KardesIndirimi.php:155-157` (host, username, password — values redacted here).
   - Bank merchant credentials (live **and** test) for 6 payment providers, plus SMTP password, plus MySQL DB password, all hardcoded in `config/params.php` and `config/db.php` / `config/web.php` (see below) — these files are part of the app source tree, not `.env`/ignored config, so they ship wherever the code ships (and are present in this survey checkout).
   - `auth.json` (repo root) contains what looks like a live GitHub OAuth token in plaintext (value redacted here).

4. **Hardcoded database and SMTP credentials in versioned config files.**
   - `config/db.php:5-7`: MySQL DSN with embedded username/password (values redacted here) — note also the malformed DSN, missing `;` between `port=3308` and `dbname=...`, which would break the connection string as written, though the app doesn't appear to use `Yii::$app->db` anywhere so this may simply be inert/unused.
   - `config/web.php:58-61`: SMTP transport `host => smtp.office365.com`, `username => no-reply@gurseltur.com.tr`, password value redacted here.
   - `config/params.php`: 6 sets of live+test payment gateway merchant IDs, store keys, usernames and passwords in cleartext (lines 16-97), values redacted here.

5. **`web/customerFiles/` upload handling has no visible content-type/extension allow-list.** `controllers/SiteController.php:429-432` (`actionKardesindirimi`): `$model->file = UploadedFile::getInstance(...)`, filename becomes `<parentTC>-<randomstring>.<clientSuppliedExtension>`, saved directly under the public web root with `$model->file->saveAs(...)`. `KardesIndirimi::rules()` (`models/KardesIndirimi.php:29-37`) marks `file` merely `'required'`/`'safe'` — there is no MIME-type, extension allow-list, or file-size validator visible anywhere in `KardesIndirimi.php` or the controller. Because the file is saved under `web/` (publicly served, per `web/.htaccess`), and the filename is guessable in part (parent's own TC number as prefix, but only the random suffix is unknown), this is a file-upload hardening gap worth revisiting even though it's not fully unauthenticated (requires a logged-in session).

6. **Dead/legacy scaffold files left in `models/`:** `models/User_1.php` and `models/LoginForm_1.php` are untouched copies of the original Yii2 project-template `User`/`LoginForm` classes (hardcoded `admin`/`demo` users, `extends \yii\base\Object` — a pre-Yii-2.0.4 base class name that likely doesn't even exist in the installed framework version, i.e. these files would fatal-error if ever autoloaded/instantiated). They are not referenced anywhere else (`app\models\User` and `app\models\LoginForm` — without the `_1` suffix — are what's actually used), but they remain in the tree.

7. **A user-specific business rule hardcoded into the controller.** `controllers/SiteController.php:284`: `if(Yii::$app->user->identity->tcNum == "68206285232"){ return $this->redirect(['kardesindirimi']); }` — a one-off redirect rule tied to a literal national ID number, with no comment explaining why, embedded inside the general login/dashboard gating logic.

8. **Significant duplicate/parallel code paths from incomplete migrations.** Bank integrations have "old" and "new" versions living side by side (`finanstaksit`/`finansodeme3d`/`finansodeme3donay` vs `finansyenitaksit`/`finansyeniodeme3d`/`finansyeniodeme3donay`; `isbanktaksit`/`isbankodeme3d`/`isbankodeme3donay` vs `isbankyenitaksit`/`isbankyeniodeme3d`/`isbankyeniodeme3ddata`/`isbankyeniodeme3donay`) with no indication in code which is currently live in production routing/views — a modernization/migration risk if both are still reachable.

9. **Large amounts of commented-out dead code left in place rather than removed**, including: an entire legacy SMS gateway integration to `processor.smsorigin.com` (`models/OgrenciSMSOnay.php:96-174`, ~80 lines), a legacy SOAP-based project list (`models/OgrenciBilgileri.php:179-202`), multiple commented-out direct-Swiftmailer-send blocks alongside the `kt.gurseltur.com.tr` HTTP replacement (e.g. `models/ServisOnKayitForm.php:299-304`, `models/ServistenAyrilmaTalebiForm.php:128-133`), and a commented-out `beforeAction` CSRF override (`controllers/SiteController.php:24-27`).

10. **`TCBilgiAlmaForm::anneKimlikDogrulama`/`babaKimlikDogrulama` (real government ID verification via `tckimlik.nvi.gov.tr`) exist but are not wired into `rules()`**, so the actual government-identity cross-check the code appears designed to perform never runs; only the TC-checksum (`ctype_digit` + mod-10/11 algorithm) is enforced.

11. **Declared-but-unused dependencies**: `yiisoft/yii2-mongodb` (composer.json) — no Mongo code anywhere; `kartik-v/yii2-social` — no Facebook/social-login code anywhere. Both add install/maintenance surface with no corresponding functionality, worth dropping during migration scoping.

12. **A likely-broken/inert code path**: `models/TCBilgiAlmaForm.php:323`: `if(!empTy($model->ogr_kan_grubu)){` — `empTy` (mixed case, not a real PHP function) would raise a fatal `Call to undefined function` error if this branch of `tcKaydiGonder()` is ever reached with the wrong casing not being an existing shim; worth flagging for functional testing before migration (this may indicate the whole `tcKaydiGonder` method is currently dead/never called, since this typo would otherwise have been caught in production use — grep shows no call sites for `tcKaydiGonder` in `controllers/`, supporting that theory).

13. **No automated tests found relevant to business logic.** `tests/` exists but was out of scope per the survey instructions (Codeception scaffold dependency present in `composer.json`); given the payment and personal-data handling in this app, the absence of any request/functional test coverage visible from the app's own directories is a notable modernization risk to flag, not just a "known issue" to fix in place.

---

## 9. Views Structure (skim only, per scope)

- `views/site/` (~40 files): one view per controller action, including a full page per bank/payment step (`akbankodeme3d.php`, `akbankodeme3donay.php`, `akbanktaksit.php`, and equivalents for finans/finansyeni/garanti/isbank/isbankyeni/yapikredi), plus `login.php`, `index.php` (and a stray `index_1.php`), `servisonkaydi.php`, `servistenayrilmatalebi.php`, `kardesindirimi.php`, `ogrenci_sms_onay.php`, `sozlesme_talep_et.php` (plus a stray copy `sozlesme_talep_et kopyası.php` — Turkish for "copy", i.e. an accidental duplicate left in the tree), `tc_talep_et.php`, `hizmet_sozlesmesi.php`, `musteri_sozlesmesi.php`, `sartname.php`, `gizlilik.php`, `iptal_iade.php`, `alevmuvafakatname.php`, `contact.php`, `about.php`, `error.php`, `result.php`.
- `views/havale/`: `index.php` (the transfer notification form, with inline jQuery posting contract-acceptance to `site/sozlesmeonay`), `response.php` (simple confirmation alert view).
- `views/layouts/main.php`: sets `SameSite=Lax` and `SameSite=None; Secure` cookie headers manually via raw `header()` calls (lines 2-3) in addition to Yii's own session handling — an unusual pattern worth checking for interaction with Yii's session cookie config in `config/web.php`.
- `mail/layouts/html.php`: minimal Swiftmailer HTML layout, standard Yii2 scaffold style.
- `web/customerFiles/`: **currently empty directory** in this checkout; its purpose (confirmed via code, not contents) is storage of sibling-discount proof-of-relationship documents uploaded through `actionKardesindirimi`, publicly served from the web root per `web/.htaccess`'s rewrite rules (no dedicated deny rule for this subdirectory was found in `web/.htaccess`).

---

## 10. Open Questions For Migration Planning

- Confirm with the GurselMobile-side survey whether `ats.gurseltur.com.tr/web/services/*` is literally the GurselMobile deployment, or a separate "ATS" system that both GurselVB22 and GurselMobile talk to. The in-code evidence (matching `localhost/GurselMobile/...` comments directly above live `ats.gurseltur.com.tr` calls with identical paths) points strongly to the former.
- Determine which of the duplicate "old" vs "yeni" (new) bank integration action sets are actually routed to in production (e.g. via web server / reverse proxy config not present in this repo) before deciding what to port.
- Determine whether `config/db.php`'s MySQL connection is truly dead (no `ActiveRecord`/`Yii::$app->db` usage was found) or whether some vendor extension (e.g. an admin/gii tool, a kartik widget with its own persistence) uses it implicitly — worth a runtime check, not just a static grep, before assuming it can be dropped.
- Get clarity from the business owner on how havale (bank transfer) reconciliation actually happens today — this codebase only shows the "customer notifies us" half of that workflow; the matching-against-actual-bank-statement half is not present anywhere in GurselVB22.
