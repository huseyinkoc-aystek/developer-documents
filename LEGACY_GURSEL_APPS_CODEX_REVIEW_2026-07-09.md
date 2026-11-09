# Legacy Gürsel / AYSTEK PHP Apps - Codex Review Notes

**Date:** 2026-07-09  
**Reviewer:** Codex  
**Scope:** Review Claude's four current-state inventories against the local source trees and AYSTEK Core standards.

Reviewed source inventories:

- `AYSHUTTLEGURSEL_CURRENT_STATE.md`
- `AYSRPAY_CURRENT_STATE.md`
- `GURSELVB22_CURRENT_STATE.md`
- `GURSELMOBILE_CURRENT_STATE.md`
- `README_FOR_AI_ASSISTANTS.md`
- `AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`

This file is the correction/synthesis layer. The four `*_CURRENT_STATE.md` files remain the detailed evidence base.

---

## 1. Executive verdict

Claude's detailed inventory is broadly useful and mostly consistent with the source code, but the high-level summary needs a few important corrections:

1. **Payment duplication is real, but not every app contains every bank integration.**
   - The correct finding is: multiple legacy apps independently implement overlapping 3D-Secure / virtual-POS flows, but their bank coverage differs by app.
   - Do not read the previous summary as "all four apps each implement Akbank + Garanti + Yapı Kredi + İşbank + QNB."

2. **GurselMobile's QR flow is not student QR attendance.**
   - It is a vehicle-level "child left behind" safety confirmation flow.
   - Actual per-student attendance/yoklama is RFID / card-tap based and tied to Trio Mobil / vehicle card-read data.

3. **GurselVB22 is genuinely a thin UI/orchestration shell.**
   - It declares some old dependencies, but no local domain DB usage was found in app code.
   - It reads/writes state through SAP SOAP, `ats.gurseltur.com.tr`, `kt.gurseltur.com.tr`, and session state.

4. **The security warning is valid and should be treated as production-impacting.**
   - The source trees contain hardcoded credentials, disabled TLS verification, disabled CSRF in broad areas, plaintext service passwords, and payment-card handling that puts the app in PCI scope.
   - The generated docs use redaction markers; that does not remediate the source repositories. Secrets should be considered exposed and rotated.

5. **Migration should consolidate domain ownership, not copy legacy controllers.**
   - AYSTEK Core already has overlapping bounded contexts: `payments`, `ledger`, `billing`, `contracts`, `subscription_contracts`, `subscriptions`, `fleet_tally`, `telematics`, `integrations`, `stock/vehicles`, and `stock/reservations`.
   - Legacy behavior should be mapped into these modules or into new DDD submodules only when the existing owner is not appropriate.

---

## 2. Source-code verification snapshot

Counts below are from app code only, excluding vendor/runtime/tests/assets/backups.

| App | PHP files | Mongo AR refs/classes | Yii DB AR refs/classes | `Yii::$app->db` refs | `Yii::$app->mongodb` / Mongo refs | CSRF disabled mentions | TLS verification disabled mentions |
|---|---:|---:|---:|---:|---:|---:|---:|
| AYShuttleGursel | 1040 | 136 | 5 | 0 | 256 | 13 | 44 |
| AYSRPAY | 234 | 19 | 0 | 0 | 22 | 3 | 11 |
| GurselVB22 | 75 | 0 | 0 | 0 | 0 | 2 | 74 |
| GurselMobile | 688 | 92 | 0 | 0 | 87 | 10 | 17 |

Interpretation:

- The "mostly MongoDB" claim is valid for AYShuttleGursel, AYSRPAY, and GurselMobile.
- GurselVB22 is the exception: no local Mongo/SQL domain persistence appears in the reviewed app code.
- Broad CSRF/TLS issues are real and cross-cutting.

---

## 3. Corrections to the previous high-level summary

### 3.1 Payment gateway coverage

Corrected app-by-app picture:

| App | Live / relevant payment behavior | Important correction |
|---|---|---|
| AYShuttleGursel | Card balance top-up / ticket-sales payment flows; Akbank and Yapı Kredi/Posnet are the relevant explicit implementations; BKM Express SDK exists but is stub/unused. | Do not count the vendored BKM bank enum as live Garanti/İşbank/QNB integration. |
| AYSRPAY | Payment-link debt collection; Garanti is the active path; Yapı Kredi exists but appears disconnected/dead. | No Akbank/İşbank/QNB implementation was found in this app's active business path. |
| GurselVB22 | Parent portal school-service payments; multiple bank action sets exist: Akbank, Garanti, Yapı Kredi, İşbank old/new, Finans/QNB old/new. | This is the app with the broadest duplicated bank surface. |
| GurselMobile | Own balance top-up UI uses Akbank virtual POS; parent-portal payment endpoints mostly log/query payment state for GurselVB22. | Its `Veliservisodemeservisleri` payment endpoints are not another bank gateway implementation. |

Migration implication:

- Build a single `payments` abstraction in AYSTEK Core:
  - `PaymentProvider`
  - `MerchantAccount`
  - `InstallmentPlan`
  - `PaymentIntent`
  - `PaymentAttempt`
  - `ThreeDSecureSession`
  - `PaymentCallback`
  - `PaymentLink`
- Gateway code must be provider-specific but centrally owned; no module should hand-roll bank XML/cURL independently.
- New implementation should use hosted/tokenized bank fields where available. AYSTEK-owned servers and DOM must not receive raw PAN/CVV.

### 3.2 "QR Yoklama" wording

The previous table label "Mobil/QR Yoklama Microservice" is misleading.

Correct wording:

- GurselMobile includes a **QR vehicle safety confirmation** flow for "mahsur kalma" checks.
- Student attendance/yoklama is a separate RFID/card-tap flow:
  - vehicle card reads,
  - Trio Mobil / GPS integrations,
  - `AracGirisCikisKartOkutmalari`,
  - `OkulKisiYoklamaBildirimleri`,
  - OneSignal parent notifications.

Migration implication:

- Model these as two separate domain concepts:
  - `VehicleSafetyCheck` / `SafetyCheckWindow` / `SafetyAlert`
  - `AttendanceEvent` / `BoardingEvent` / `NoShowNotification`
- Do not merge QR safety confirmation and student attendance into one "QR yoklama" module.

### 3.3 Balance model generalization

The append-only movement pattern is valid for the transport card systems:

- card top-up,
- fare debit,
- refund,
- transfer,
- invoice/batch load,
- reservation-related debit.

But AYSRPAY is not the same model:

- It manages customer debt deltas from Gürsel SAP and payment records.
- `DebtList` accumulation and single-row delete behavior need migration review.

Migration implication:

- Use `ledger` for money/account movements, but split account types:
  - `transport_card_account`
  - `customer_debt_account`
  - `subscription_contract_account`
  - optional `reservation_hold_account`
- Preserve append-only entries where they represent auditable financial movement.
- Do not blindly migrate legacy "current balance by summing rows" without idempotency, concurrency control, and reconciliation strategy.

---

## 4. AYSTEK Core bounded-context mapping

Recommended target ownership:

| Legacy concept | Target AYSTEK Core owner |
|---|---|
| Virtual POS providers, 3D Secure, payment callbacks, payment links | `payments` |
| Card balance, debt/payment movements, refunds, top-ups, fare debits | `ledger` |
| Invoices / billing notices / receivables | `billing` |
| Student/parent service contracts and approval snapshots | `contracts` or `subscription_contracts` depending on product semantics |
| Shuttle subscription/service entitlement | `subscriptions` + `subscription_contracts` |
| Vehicles, routes, GPS providers, trip verification | `fleet_tally` + `telematics` + `stock/vehicles` |
| UETDS and transport regulatory reporting | `telematics/uetds` if operational, otherwise dedicated reporting submodule |
| SAP SOAP, Trio Mobil, SMS Origin, OneSignal, kt relay, ATS endpoints | `integrations` / `erp_integrations` adapters |
| Parent pre-registration, withdrawal, sibling discount requests | New transport-school service module or `subscriptions` workflow, not raw controller ports |
| QR vehicle safety confirmation / missed-scan alarm | `fleet_tally` or `telematics/alerts` subdomain |
| Per-student RFID attendance | `telematics/telemetry` + dedicated attendance/read-model service |

Critical platform rules from `AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`:

- Every migrated module must be tenant-aware and site-aware.
- Module access must be controlled by ModuleGate / permissions server-side.
- Repositories must enforce `tenant_id` and `site_id` where applicable.
- Public UI and customer flows must consume backend authority; frontend cannot own business decisions.

---

## 5. Security remediation backlog

These are migration blockers, not cosmetic cleanup:

1. **Credential rotation**
   - Rotate SMTP, MongoDB, MySQL, SAP SOAP, bank VPOS, SMS, OneSignal, Trio Mobil, FTP, Google API, GitHub OAuth, cookie-validation, and shared service tokens.
   - Do this even if values are moved to env/secret manager, because the source history / filesystem exposure already exists.

2. **Payment PCI scope reduction**
   - Do not port server-rendered PAN/CVV hidden-input flows.
   - Prefer hosted payment page, hosted fields, iframe, or tokenization.
   - Persist only bank-safe tokens, references, auth codes, masked PAN where legally/contractually allowed, and explicit allow-listed callback fields.

3. **TLS verification**
   - Remove all `CURLOPT_SSL_VERIFYPEER=false` / `CURLOPT_SSL_VERIFYHOST=false` behavior.
   - Use CA bundles and explicit endpoint allow-lists where partner systems require special handling.

4. **CSRF and API auth**
   - Browser forms must keep CSRF protection.
   - Bank callbacks and machine APIs need narrow exemptions with signature/HMAC/token verification.
   - Replace shared human-login credentials for service APIs with scoped service accounts, client credentials, mTLS, signed requests, or token exchange.

5. **Plaintext service passwords**
   - Do not migrate plaintext password comparison models.
   - Hash at rest and scope permissions per integration/action.

6. **File uploads**
   - Sibling-discount proof files should not be written directly under public webroot.
   - Enforce extension/MIME/size validation, malware scanning where available, private object storage, and signed download URLs.

7. **Logging and PII**
   - Do not persist raw request bodies containing card/bank/personal fields.
   - Add log redaction and explicit field allow-lists.

---

## 6. Do-not-migrate list unless business confirms

Do not carry these forward by default:

- AYSRPAY generic e-commerce leftovers: Category/Menu, Ozon fields, unused push helpers, broken/dead `SentemailController::actionOrderinfo`, unused mail themes.
- AYSRPAY disconnected Yapı Kredi path unless production usage is proven.
- GurselVB22 old/new duplicate bank action sets without production route evidence.
- GurselVB22 `User_1.php`, `LoginForm_1.php`, and hardcoded TC-specific redirect rule.
- GurselMobile `AlertsTempController`.
- Stock Yii scaffold `HelloController` / demo `User` classes.
- Dead BKM Express stub paths unless the business still needs BKM.
- One-off maintenance/debug actions with real PII or card numbers embedded in source.

Rule: if no current production owner confirms a flow, document it as legacy/dead and exclude it from the first migration slice.

---

## 7. Recommended migration sequence

1. **Security freeze and credential rotation**
   - Inventory all external credentials and rotate before shared CI/GitHub exposure.

2. **Canonical payment service**
   - Implement provider abstraction, hosted/tokenized flow, callbacks, payment attempts, and idempotency.
   - Port only one provider first with tests; add others through the same interface.

3. **Ledger/account model**
   - Define append-only ledger entries and account types.
   - Add idempotency keys for vehicle/device reads, payment callbacks, and SAP confirmations.

4. **Transport/fleet core mapping**
   - Map vehicles/routes/cards/riders/institutions into tenant/site-aware Core models.
   - Reuse `fleet_tally`, `telematics`, `stock/vehicles`, and `stock/reservations` where they own the concept.

5. **Parent portal workflows**
   - Rebuild login, contract approval, TC info capture, sibling discount, withdrawal, and havale notification as typed workflows.
   - Avoid scattered cURL calls; all external communication must go through `integrations` adapters.

6. **QR safety and RFID attendance**
   - Keep safety confirmation and attendance as separate workflows.
   - Add alert rules, notification routing, and audit logs.

7. **Reporting and reconciliation**
   - Recreate only confirmed reports.
   - Use read models/materialized views instead of controller-triggered ad hoc rebuilds.

---

## 8. Open questions before implementation

- Which bank providers are actually active in production today per product/site?
- Which "old" vs "new" GurselVB22 İşbank/Finans paths are still routed?
- Is `ats.gurseltur.com.tr` definitively the deployed GurselMobile app, a reverse proxy to it, or a separate ATS service?
- Which SAP SOAP endpoints are source-of-truth versus side-effect logging?
- How is havale reconciliation actually completed today?
- Which reports are legally/operationally required versus historical convenience screens?
- Are legacy RFID card IDs/person identifiers allowed to be migrated as-is, or must they be re-tokenized/anonymized?
- What is the target bank integration contract for hosted/tokenized fields in 2026?

---

## 9. Final position

Use Claude's four `*_CURRENT_STATE.md` files as detailed source inventories, but use this review as the corrected migration lens:

- consolidate duplicate payment code,
- separate QR safety from RFID attendance,
- treat secrets as exposed,
- preserve append-only financial audit where valid,
- drop dead scaffolding,
- map behavior into AYSTEK Core's existing DDD modules instead of porting Yii controllers.
