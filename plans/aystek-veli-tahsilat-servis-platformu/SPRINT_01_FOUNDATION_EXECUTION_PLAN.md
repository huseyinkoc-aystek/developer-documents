# Sprint 01 — Transport Platform Foundation Execution Plan

**Tarih:** 2026-05-09  
**Durum:** Uygulamaya hazır sprint planı  
**Bağımlı ADR:** `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md`  
**Amaç:** Kod geliştirmeye başlamadan önce transport platformun Core içindeki seed, permission, menu, route, dashboard ve domain sınır temelini kurmak.

---

## 1. Sprint Hedefi

Sprint 01 sonunda AYSTEK Core içinde aşağıdaki temel görünür olmalıdır:

- `transport_core` backend foundation.
- `transport_access`, `transport_wallet`, `transport_portal` için module constant, doc ve ADR sınırları.
- `gursel-turizm` tenant altında mevcut `aractakip` yapısıyla çakışmadan `veli.gurseltur.com.tr` site seed’i.
- Transport Platform admin menüsü.
- Transport permission seed seti.
- Transport translation seed minimum key seti, 28 dil ile.
- Admin dashboard shell.
- Health/status endpoint.
- Event prefix ve event catalog temel dokümantasyonu.

Bu sprintte gerçek ödeme, gerçek SAP gönderimi, gerçek Mifare cihaz entegrasyonu veya mobil uygulama yoktur.

---

## 2. Kapsam

### 2.1 Backend Kapsamı

Oluşturulacak paketler:

```text
internal/modules/transport_core/
├── doc.go
├── module.go
├── dto/
│   └── dashboard.go
├── handler/
│   └── admin.go
├── model/
│   └── dashboard.go
├── repository/
│   ├── interface.go
│   └── mongo.go
├── service/
│   ├── service.go
│   └── service_test.go
└── routes.go

internal/modules/transport_access/
├── doc.go
└── module.go

internal/modules/transport_wallet/
├── doc.go
└── module.go

internal/modules/transport_portal/
├── doc.go
└── module.go
```

Sprint 01’de `transport_access`, `transport_wallet`, `transport_portal` için sadece boundary constants açılır. DDD tam katmanları sonraki sprintlerde, gerçek aggregate açıldığında oluşturulur.

### 2.2 Seed Kapsamı

Yeni dosya:

```text
internal/seeds/transport_platform_seed.go
```

Güncellenecek dosyalar:

```text
internal/seeds/runner.go
internal/seeds/gursel_turizm_seed.go
internal/seeds/menu_seed.go
internal/seeds/permission_seed.go
internal/seeds/translation_seed.go
```

Seed davranışı:

- `gursel-turizm` tenant zaten varsa idempotent update.
- `aractakip.gurseltur.com.tr` mevcut seed korunur.
- Yeni site/store:
  - code: `veli`
  - domain: `veli.gurseltur.com.tr`
  - type: `transport_parent_portal`
  - default language: `tr`
  - default currency: `TRY`
- Transport module feature flags/settings seed edilir.
- Transport menu ve permission seed’i tenant/site context ile uyumlu çalışır.

### 2.3 Frontend Kapsamı

Oluşturulacak önerilen dosyalar:

```text
app/(dashboard)/admin/transport-platform/dashboard/page.tsx
app/(dashboard)/admin/transport-platform/page.tsx
components/transport-platform/TransportPlatformDashboard.tsx
lib/api/transport-platform.ts
lib/types/transport-platform.ts
```

Dashboard ilk sürümde:

- tenant/site bağlamı,
- health bilgisi,
- modül boundary kartları,
- ileriki hatlara bağlantı placeholder’ları,
- permission denied / loading / error state

göstermelidir.

---

## 3. Backend Tasarım Detayı

### 3.1 Module Constants

Her namespace için module constant:

```go
const ModuleName = "transport_core"
const ModuleNameAccess = "transport_access"
const ModuleNameWallet = "transport_wallet"
const ModuleNamePortal = "transport_portal"
```

Paketler kendi `module.go` içinde sadece kendi adını bilmelidir. Cross-module string tekrarları minimumda tutulmalıdır.

### 3.2 Dashboard Endpoint

Route:

```text
GET /admin/transport-platform/dashboard
```

Permission:

```text
transport_core:read
```

Response veri şekli:

```json
{
  "module": "transport_core",
  "tenant_id": "gursel-turizm",
  "site_id": "veli",
  "status": "ready",
  "boundaries": [
    {"key": "transport_core", "status": "foundation"},
    {"key": "transport_access", "status": "planned"},
    {"key": "transport_wallet", "status": "planned"},
    {"key": "transport_portal", "status": "planned"}
  ],
  "next_actions": [
    "create_institution_registry",
    "create_person_registry",
    "create_route_policy"
  ]
}
```

Response envelope middleware zaten saracağı için handler raw object dönebilir.

### 3.3 Service

Sprint 01 service sadece context-aware dashboard summary üretir.

Zorunlu davranışlar:

- Context’ten tenant ID okur.
- Context’ten site/store ID okunabiliyorsa döndürür.
- Context yoksa güvenli fallback değil hata döndürür; tenant leak olmamalıdır.
- Status değerleri typed constant olmalıdır.

### 3.4 Repository

Sprint 01’de repository iki şekilde yapılabilir:

1. Dashboard repository yalnızca future metrics için boş interface olarak kalır.
2. Settings/store collection’dan site varlığını okuyarak dashboard health üretir.

Öneri: İlk sprintte minimal repository + DB ping/site exists check yapılır. Büyük business logic eklenmez.

---

## 4. Seed Detayı

### 4.1 Permission Modules

`permission_seed.go` içine eklenecek module seti:

```text
transport_core
transport_institutions
transport_persons
transport_groups
transport_routes
transport_trips
transport_access
transport_access_devices
transport_access_transactions
transport_tickets
transport_subscriptions
transport_wallet
transport_wallet_topups
transport_portal
transport_guardians
transport_attributes
transport_service_requests
transport_debts
transport_contracting
transport_payments
transport_bank_transfers
transport_announcements
transport_sap_dispatch
transport_ops
```

Varsayılan actions:

```text
read, create, update, delete, approve, reject, execute, export
```

Ek hassas actions:

```text
unmask, replay, adjust
```

### 4.2 Menu Entries

Üst menü:

```text
transport-platform
```

Title key:

```text
menu.operations.transport_platform
```

Path:

```text
/admin/transport-platform/dashboard
```

Alt menüler ilk sprintte görünür olabilir ama “planned/foundation” ekranına gidebilir:

- Dashboard
- Kurumlar
- Taşınan Kişiler
- Gruplar ve İndirimler
- Güzergâh ve Ücretler
- Sefer ve Çetele
- Geçiş Cihazları
- Geçiş İşlemleri
- Cüzdan ve Yüklemeler
- Veli/Personel Portalı
- SAP Dispatch
- Audit

### 4.3 Translation Keys

Minimum key set:

```text
menu.operations.transport_platform
menu.transport_platform.dashboard
menu.transport_platform.institutions
menu.transport_platform.persons
menu.transport_platform.groups
menu.transport_platform.routes
menu.transport_platform.trips
menu.transport_platform.access_devices
menu.transport_platform.access_transactions
menu.transport_platform.wallet
menu.transport_platform.portal
menu.transport_platform.sap_dispatch
menu.transport_platform.audit

transport_platform.dashboard.title
transport_platform.dashboard.description
transport_platform.boundary.core
transport_platform.boundary.access
transport_platform.boundary.wallet
transport_platform.boundary.portal
transport_platform.status.foundation
transport_platform.status.planned
transport_platform.status.ready
transport_platform.error.context_required
```

Bu key’ler 28 dil ile seed edilmelidir. Türkçe ve İngilizce kaliteli; diğer dillerde güvenli fallback kabul edilebilir ama boş değer kabul edilmez.

### 4.4 Transport Platform Seed

`transport_platform_seed.go` içinde:

- Transport platform default module settings.
- Event prefix catalog başlangıç kayıtları.
- Varsayılan status/type enum referansları.
- `veli` site bootstrap çağrısı veya `gursel_turizm_seed.go` içinden çağrılacak helper.

Idempotency zorunludur.

---

## 5. Frontend Tasarım Detayı

### 5.1 API Client

`lib/api/transport-platform.ts`:

- Envelope normalize eder.
- `getTransportPlatformDashboard()` fonksiyonu sunar.
- Hata mesajlarını typed normalize eder.

### 5.2 Types

`lib/types/transport-platform.ts`:

```ts
export type TransportBoundaryKey =
  | "transport_core"
  | "transport_access"
  | "transport_wallet"
  | "transport_portal";

export interface TransportBoundaryStatus {
  key: TransportBoundaryKey;
  status: "foundation" | "planned" | "ready" | "disabled";
}

export interface TransportPlatformDashboard {
  module: string;
  tenant_id: string;
  site_id?: string;
  status: "foundation" | "planned" | "ready";
  boundaries: TransportBoundaryStatus[];
  next_actions: string[];
}
```

### 5.3 Page

Dashboard:

- Server/client pattern mevcut repo standardına göre seçilir.
- Hard-coded text yok; key kullanılır.
- Skeleton/loading vardır.
- Empty state vardır.
- Error state vardır.

---

## 6. Test Planı

### 6.1 Backend

Komutlar:

```bash
cd /workspace/go/aystek-core
go test ./internal/modules/transport_core/...
go test ./internal/seeds/... ./internal/platform/middleware/... ./internal/platform/router/...
go build ./cmd/api
```

Testler:

- `service.GetDashboard` tenant context yoksa hata döndürür.
- `service.GetDashboard` context varsa boundary listesi döndürür.
- Seed idempotency: aynı seed iki kez çalışınca duplicate oluşmaz.
- Route permission middleware doğru module/action ile bağlanır.

### 6.2 Frontend

Komutlar:

```bash
cd /workspace/node/aystek-core-next
npm run check-i18n
npm run build
```

Testler:

- API client envelope normalize eder.
- Dashboard component loading/error/ready state render eder.
- Translation key eksikleri oluşmaz.

---

## 7. Definition of Done

Sprint 01 tamam sayılmak için:

- Backend compile başarılı.
- Targeted backend tests başarılı.
- Frontend build veya en azından yeni dosyalardan kaynaklı type/build hatası yok.
- `translation_seed.go` 28 dil kuralı sağlanmış.
- `menu_seed.go` transport platform menüsünü içeriyor.
- `permission_seed.go` transport modules/actions setini içeriyor.
- Seed runner transport seed’i çağırıyor.
- `gursel-turizm` altında `veli.gurseltur.com.tr` site seed’i idempotent.
- Admin dashboard route erişilebilir.
- Bu klasöre sprint checkpoint dokümanı eklenmiş.

---

## 8. Sprint Dışı Bırakılanlar

Bu sprintte yapılmayacaklar:

- Gerçek Institution CRUD.
- Gerçek Person import.
- Gerçek kart/QR okutma.
- Gerçek wallet balance.
- Gerçek payment capture.
- Gerçek SAP dispatch.
- Mifare/reader native entegrasyonu.
- Offline sync implementation.
- Public veli ödeme akışı.

Bu işleri foundation olmadan yapmak teknik borç üretir.

---

## 9. İlk Kodlama Sırası

1. Backend module constants ve dashboard service.
2. Backend routes ve main.go wiring.
3. Seed dosyaları: permission → translation → menu → transport platform → runner.
4. `gursel_turizm_seed.go` içinde `veli` site helper.
5. Frontend API/types/dashboard page.
6. Tests.
7. Build/check-i18n.
8. Sprint checkpoint.

---

## 10. Uygulama Notları

- `transport_core` ilk gerçek modüldür; diğerleri boundary placeholder olarak açılır.
- `transport_portal` önceki planın adı olarak değil, veli/personel portal subdomain’i olarak düşünülmelidir.
- `WalletInstrument` kavramı ulaşım kartı/QR içindir; banka kartı ile karıştırılmamalıdır.
- Hassas veri ve offline logic Sprint 01’de tasarım notu olarak kalır; veri toplamaya başlanmaz.
- Admin route path `transport-platform` olmalıdır; backend module adı `transport_core` olabilir.

