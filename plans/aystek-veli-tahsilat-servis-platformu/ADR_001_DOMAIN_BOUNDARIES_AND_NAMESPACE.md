# ADR 001 — Domain Sınırları ve Namespace Kararı

**Tarih:** 2026-05-09  
**Durum:** Kabul edildi  
**Karar türü:** Mimari sınır / bounded context  
**Kapsam:** KT kart/QR geçiş, cüzdan, çetele, GPS doğrulama ve veli/personel tahsilat platformu  

---

## 1. Bağlam

Proje iki büyük iş hattını aynı AYSTEK Core platformu içinde birleştirir:

1. **KT / Mobility Access:** Kartlı/QR’lı taşımacılık, kullanıcı bazlı cüzdan, Mifare/QR okutma, misafir bilet, abonelik, rota/ücret, araç hareket raporu, GPS doğrulama ve çetele.
2. **Veli/Personel Tahsilat:** Veli/personel portalı, kişi bağlama, borç hesabı, sözleşme gereksinimi, dinamik kişi alanları, sanal POS, havale bildirimi ve SAP dispatch.

İlk düşünce `transport_portal` altında ilerlemekti. Ancak kullanıcı anlatımı, bu ismin yalnızca ödeme/self-service tarafını karşıladığını; geçiş, offline cihaz, cüzdan, güzergâh ve çetele tarafının daha geniş bir domain olduğunu gösterdi.

Bu nedenle tek namespace altında her şeyi toplamak ileride:

- büyük paket karmaşası,
- mükerrer servisler,
- wallet/payment karışıklığı,
- ulaşım kartı ile banka kartı kavramlarının karışması,
- offline transaction ile online payment intent sınırlarının bulanması,
- admin menü ve permission şişmesi

üretecektir.

---

## 2. Karar

AYSTEK Core içinde dört ana namespace kullanılacaktır:

| Namespace | Ana sorumluluk |
|---|---|
| `transport_core` | Kurum, taşınan kişi, grup, güzergâh, sefer, çetele, araç hareket raporu ve ortak ulaşım master data |
| `transport_access` | Mifare/QR okuyucu, geçiş cihazı, access transaction, offline sync, bilet, abonelik ve geçiş kuralları |
| `transport_wallet` | Kullanıcı bazlı bakiye, yükleme, iade, kredi limiti, TOSLA/POS top-up ve fatura sorumluluğu |
| `transport_portal` | Veli/personel public self-service, borç hesabı, sözleşme gereksinimi, dinamik alan, havale ve duyuru |

Mevcut Core modülleri tekrar yazılmayacak; adapter/service dependency olarak kullanılacaktır:

| Mevcut modül | Kullanım |
|---|---|
| `payments` | Sanal POS, payment intent, gateway, installment, risk |
| `contracts` | Template/content ve versiyonlu sözleşme içerikleri |
| `ledger` / `billing` | Finansal hareket, faturalama, accounting bridge |
| `integrations` / `api_sharing` | SAP, legacy bridge, müşteri web servisleri, doğrulanmış iş paylaşımı |
| `telematics` / `fleet_tally` | GPS, cihaz, trip doğrulama, çetele desenleri |
| `sms` / `notifications` | OTP, yaklaşma bildirimi, işlem ve operasyon bildirimleri |
| `stock/vehicle-assets` | Araç varlığı ve zimmet/asset ilişkileriyle entegrasyon |

---

## 3. Namespace Sınırları

### 3.1 `transport_core`

Bu modül platformun ulaşım master data ve doğrulanmış iş çekirdeğidir.

Alt paket önerileri:

- `institutions`
- `persons`
- `groups`
- `routes`
- `trips`
- `movement_reports`
- `sapdispatch` veya ortak dispatch adapter
- `ops`

Sahip olduğu aggregate’lar:

- `Institution`
- `TransportedPerson`
- `PersonGroup`
- `Route`
- `RouteDirection`
- `FarePolicy`
- `Trip`
- `TripInference`
- `MovementReport`
- `VerifiedWorkUnit`

Sahip olmadığı alanlar:

- Kart/QR okutma ham transaction’ları: `transport_access`
- Kullanıcı bakiyesi: `transport_wallet`
- Veli borcu/sözleşmesi: `transport_portal`
- Banka sanal POS provider detayları: `payments`

### 3.2 `transport_access`

Bu modül geçiş noktası ve okutma domain’idir.

Alt paket önerileri:

- `devices`
- `instruments`
- `transactions`
- `tickets`
- `subscriptions`
- `offline_sync`
- `rules`

Sahip olduğu aggregate’lar:

- `AccessDevice`
- `WalletInstrument`
- `AccessTransaction`
- `OfflineTransactionBatch`
- `GuestTicket`
- `SubscriptionPass`
- `AccessRule`
- `RefundWindow`

Sahip olmadığı alanlar:

- Ana bakiye hesabı: `transport_wallet`
- Seferin nihai çetele raporu: `transport_core`
- POS ödemesi: `payments`

### 3.3 `transport_wallet`

Bu modül kullanıcı bazlı cüzdan ve değer hareketi domain’idir.

Alt paket önerileri:

- `accounts`
- `transactions`
- `topups`
- `refunds`
- `credit_policies`
- `invoice_responsibility`

Sahip olduğu aggregate’lar:

- `WalletAccount`
- `WalletTransaction`
- `TopUpIntent`
- `TopUpProviderBinding`
- `CreditPolicy`
- `RefundTransaction`
- `InvoiceResponsibility`

Sahip olmadığı alanlar:

- Ulaşım kartı fiziksel okuma protokolü: `transport_access`
- Banka gateway credential: `payments/gateways`
- Borç hesabı: `transport_portal/debts`

### 3.4 `transport_portal`

Bu modül veli/personel self-service ve borç tahsilat domain’idir.

Alt paket önerileri:

- `guardians`
- `attributes`
- `service_requests`
- `debts`
- `contracting`
- `payments`
- `bank_transfers`
- `announcements`

Sahip olduğu aggregate’lar:

- `GuardianLink`
- `LinkRequest`
- `DynamicFieldDefinition`
- `FieldRequirementPolicy`
- `PersonAttributeValue`
- `ServiceRequest`
- `DebtAccount`
- `DebtItem`
- `ContractRequirement`
- `ContractSignature`
- `PaymentQuote`
- `BankTransferNotice`
- `Announcement`

Sahip olmadığı alanlar:

- Genel wallet bakiyesi: `transport_wallet`
- Kart/QR access transaction: `transport_access`
- Kurum/person master data’nın nihai otoritesi: `transport_core`

---

## 4. Kritik Kavram Ayrımları

### 4.1 Banka Kartı vs Ulaşım Kartı

| Kavram | Domain | Açıklama |
|---|---|---|
| Banka/kredi kartı | `payments` | PAN/CVV AYSTEK’e girmemeli; hosted/provider component kullanılmalı |
| Mifare ulaşım kartı | `transport_access` | Wallet’a bağlı instrument; kart değişince bakiye taşınmaz, hesap kullanıcıdadır |
| Mobil QR kimlik | `transport_access` | Kullanıcı veya geçici bilet için access instrument |
| Misafir QR bilet | `transport_access/tickets` | Tek kullanımlık route bazlı bilet |

### 4.2 Wallet vs DebtAccount

| Kavram | Kullanım |
|---|---|
| `WalletAccount` | KT/üniversite geçişlerinde bakiye, yükleme, iade, eksi bakiye |
| `DebtAccount` | Veli/personel dönemsel borç, sözleşme ve tahsilat yükümlülüğü |

Bu iki model birleşmemelidir. Ancak ikisi de `payments`, `ledger`, `billing` ve `sapdispatch` ile entegre olur.

### 4.3 Trip vs AccessTransaction

| Kavram | Kullanım |
|---|---|
| `AccessTransaction` | Tek okutma/biniş/iniş/QR olayı |
| `Trip` | Araç ve zaman bağlamında sefer/tur |
| `MovementReport` | Sefer kapanınca çetele ve faturalama dayanağı |

Access transaction ham olaydır; çeteleye doğrudan eşit değildir.

### 4.4 Person vs PortalUser

| Kavram | Kullanım |
|---|---|
| `TransportedPerson` | Taşınan öğrenci/personel/öğretim görevlisi |
| `PortalUser` | Sisteme giriş yapan veli/personel/kullanıcı |
| `GuardianLink` | PortalUser ile TransportedPerson arasındaki doğrulanmış ilişki |

---

## 5. Admin Menü Kararı

Tek üst menü:

- Key: `transport-platform`
- Title key: `menu.operations.transport_platform`
- Path: `/admin/transport-platform/dashboard`

Alt bölümler:

| Bölüm | Sahip modül |
|---|---|
| Dashboard | `transport_core` |
| Kurumlar | `transport_core/institutions` |
| Taşınan Kişiler | `transport_core/persons` |
| Gruplar ve İndirimler | `transport_core/groups` |
| Güzergâh ve Ücretler | `transport_core/routes` |
| Sefer ve Çetele | `transport_core/trips` |
| Geçiş Cihazları | `transport_access/devices` |
| Geçiş İşlemleri | `transport_access/transactions` |
| Bilet ve Abonelikler | `transport_access/tickets` |
| Cüzdan ve Yüklemeler | `transport_wallet` |
| Veli/Personel Portalı | `transport_portal` |
| SAP Dispatch | `transport_core` veya ortak dispatch |
| Audit | `transport_core/ops` |

---

## 6. Permission Kararı

Temel permission module listesi:

- `transport_core`
- `transport_institutions`
- `transport_persons`
- `transport_groups`
- `transport_routes`
- `transport_trips`
- `transport_access`
- `transport_access_devices`
- `transport_access_transactions`
- `transport_tickets`
- `transport_subscriptions`
- `transport_wallet`
- `transport_wallet_topups`
- `transport_portal`
- `transport_guardians`
- `transport_attributes`
- `transport_service_requests`
- `transport_debts`
- `transport_contracting`
- `transport_payments`
- `transport_bank_transfers`
- `transport_announcements`
- `transport_sap_dispatch`
- `transport_ops`

Standart action seti:

- `read`
- `create`
- `update`
- `delete`
- `approve`
- `reject`
- `execute`
- `export`
- `replay`
- `unmask`

Hassas action’lar:

- `transport_attributes:unmask`
- `transport_ops:export`
- `transport_sap_dispatch:replay`
- `transport_bank_transfers:approve`
- `transport_wallet:adjust`

Bu action’lar Admin’e seed ile gelir; ancak gerekiyorsa `role_seed.go` içindeki super-admin-only kararları ayrı dokümante edilmelidir.

---

## 7. Seed Kararı

Yeni seed dosyaları önerisi:

| Seed | İçerik |
|---|---|
| `transport_core_seed.go` | Site bootstrap, kurum tipleri, temel route/service enums |
| `transport_access_seed.go` | Access rule defaults, ticket/subscription status labels |
| `transport_wallet_seed.go` | Wallet transaction types, credit policy defaults |
| `transport_portal_seed.go` | Portal status defaults, debt/contract state labels |
| `transport_event_catalog_seed.go` | Event key kataloğu |

Mevcut seed dosyaları güncellenecek:

- `menu_seed.go`
- `permission_seed.go`
- `translation_seed.go`
- `runner.go`
- gerekirse `gursel_turizm_seed.go`

---

## 8. Event Prefix Kararı

Event prefix standardı:

- `transport.core.*`
- `transport.access.*`
- `transport.wallet.*`
- `transport.portal.*`

Örnekler:

- `transport.core.trip.started`
- `transport.core.trip.closed`
- `transport.core.movement_report.verified`
- `transport.access.transaction.recorded`
- `transport.access.offline_batch.synced`
- `transport.access.ticket.consumed`
- `transport.wallet.topup.succeeded`
- `transport.wallet.balance.debited`
- `transport.wallet.refund.created`
- `transport.portal.guardian_link.approved`
- `transport.portal.contract.accepted`
- `transport.portal.payment.captured`

Her event:

- `tenant_id`
- `site_id`
- `institution_id` veya neden yok bilgisi
- `correlation_id`
- `event_version`
- `actor`
- `subject`
- `occurred_at`
- `idempotency_key`

taşımalıdır.

---

## 9. Sonuçlar

Bu kararın artıları:

- Ürün büyürken domain sınırları korunur.
- Veli ödeme ile KT geçiş sistemi aynı kişi/kurum/finans omurgasını paylaşır.
- Wallet, payment ve debt kavramları karışmaz.
- Offline geçiş ihtiyaçları online ödeme akışını kirletmez.
- Campus service genişlemesi ulaşım domain’ini bozmadan yapılabilir.
- AYSTEK Core’un multi-tenant, multi-site ve event-driven yapısına uyum korunur.

Eksileri:

- İlk foundation sprint’i biraz daha kapsamlıdır.
- `main.go` wiring ve seed sayısı artar.
- Başlangıçta daha fazla ADR ve test iskeleti gerekir.

Karar: Bu maliyet kabul edilir; çünkü uzun vadeli mükerrer geliştirme ve domain karışıklığı riski çok daha yüksektir.

