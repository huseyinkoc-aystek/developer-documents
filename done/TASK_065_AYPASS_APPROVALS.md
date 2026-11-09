# TASK 065 — AYPASS Onay Mekanizması Teknik Referansı

**Tarih:** 2026-10-28 · **Durum:** Backend, seed ve panel işi tamamlandı ve origin'e gönderilmiş dallara aktarıldı (`74a5453`, `7b03fe7`). Yerel derleme ve test sonuçları Bölüm 12'de. Canlı ortam doğrulanmadı. Rev.06 kabul kriterleri henüz karşılanmıyor.

**Backend:** `/workspace/go/aystek-core`, dal `aypass-approval-ui`. Taban commit `f48748f`, birleşik ağaç `f286fd0`, Gürsel GPS kararı `c27c3b5`, iade anahtarı zorunluluğu `a00f69d`, B2B yol düzeltmesi `de3c4fa`, wallets:approve tohumu `18e27f8`, mock adaptör kapsamı `33467c8` (origin'e gönderilmiş; dal başı `74a5453`, birleşme commit'i).
**Panel:** `/workspace/node/aystek-core-next`, dal `approval-policy-ui`. Commit'ler `e5217b2`, `ab6fd9b`, `e89592d`, `dcf9d38`, `7b03fe7` (origin'e gönderilmiş; dal başı `7b03fe7`). Taban `main` (`56907eb`).
**Installer belgeleri:** `/workspace/go/aystek-core-installer` (commit edilmemiş).

Bu belge `README_FOR_AI_ASSISTANTS.md` içindeki görev belgesi kuralına göre yazıldı. Yönetici özeti için `TASK_065_SUMMARY.md` dosyasına bakın.

---

## 1. Kapsam ve kaynaklar

### 1.1 Kapsam

Bu iş, AYPASS geçiş sisteminin Rev.06 gereksinimlerinden onay (maker-checker) katmanını kapsar. Kapsamdaki on üç iş kalemi şunlardır:

1. Onay motoru (`internal/modules/approvals`): kapı kataloğu, kiracı ve site bazlı dinamik politikalar, onay talepleri, tek kullanımlık talep tüketimi ve karar kaydı.
2. Ödeme iadesi kapısı (`internal/modules/payments/intents`): eşik üstü iadeler onay talebi ister.
3. B2B satın alma siparişleri (`internal/modules/b2b/dealers`): politika tabanlı, sıralı onay zinciri.
4. Banka havalesi maker-checker (`internal/modules/wallets`): `wallets:approve` izni ve beyan eden ile onaylayanın ayrı olması.
5. Manuel sefer onayı (`internal/modules/fleet_tally/trip_assignments`): `approve` izni ve bağımsız onaylayan.
6. Erişim noktaları: pasif erişim noktası olayları ve yaya binişlerini reddeder.
7. Roller ve izinler: rol anahtarı doğrulaması, değişmez anahtar, dil bazlı etiketler.
8. Faturalama: GİB hatalarında sahte UUID yerine hata döner; zaten kesilmiş fatura yeniden gönderilmez.
9. Menüler: onay menüleri rol adı yerine izinle kapılanır.
10. Seed'ler: `approvals_seed`, modül hata çevirileri, onay arayüzü çevirileri ve rol çevirileri.
11. Testler: yeni testler ve `Test{Method}_{Scenario}` biçimine geçirilen test fonksiyonları.
12. Panel: onay politikası ve onay talebi ekranları, iade bildirimi, rol etiketi düzenleyicisi.
13. Operasyon belgeleri: kurulum rehberi, işletim kılavuzu, dağıtım senaryoları ve installer belgeleri.

Bu belgenin kapsamına bir ek daha girer: Gürsel GPS sağlayıcı yapılandırmalarının yalnızca ilk oluşturmada yazılması (`c27c3b5`, `f286fd0` üzerine eklendi). Ayrıntı Bölüm 9.4'te. İade anahtarının zorunlu hale getirilmesi (`a00f69d`, `f286fd0` üzerine eklendi) da bu kapsamdadır. Ayrıntı Bölüm 4.3.1'de. Bu turda ayrıca B2B yol düzeltmesi (`de3c4fa`), `wallets:approve` tohumu (`18e27f8`) ve mock adaptörlerin yalnızca production dışında kaydedilmesi (`33467c8`) yapıldı. Ayrıntı Bölüm 1.2'de.

### 1.2 Değişiklik özeti

`f48748f..f286fd0` aralığında 25 commit vardır. `git diff --stat` 111 dosyanın değiştiğini gösterir: 31 dosya eklendi, 73 dosya değiştirildi, 7 dosya silindi. Toplam 9614 satır eklendi, 2431 satır silindi.

Silinen dosyalar, AD-04 kararı gereği kullanılmayan eski onay modülüne aittir:

- `internal/modules/approvals/handlers/approvals.go`
- `internal/modules/approvals/handlers/workflows.go`
- `internal/modules/approvals/model/approval.go`
- `internal/modules/approvals/repository/repository.go`
- `internal/modules/approvals/service/approvals.go`
- `internal/modules/approvals/service/interface.go`
- `internal/modules/integrations/governance/approval.go`

Eski modül silinmiş, yerine yeni bir onay modülü yazılmıştır. Silme commit'i `eb581c7`dir. Yeniden tasarım sonraki commit'lerde yapılmıştır.

`f286fd0` sonrasında origin/aypass-approval-ui üzerinde şu commit'ler vardır (en yeni başta). Birleşme commit'i `74a5453`, P0 dalı `aypass-faz0-p0-fixes`'i ana dala bağlar:

- `74a5453` Merge branch 'aypass-faz0-p0-fixes' into aypass-approval-ui
- `25abcc1` Update the vehicle incidents test fake to the plate reason signatures
- `93b4232` Require a reason for every stock vehicle plate change
- `7eb7bfa` Make system account creation and journal confirmation race-safe
- `33467c8` Register mock e-invoice and ERP adapters only outside production
- `18e27f8` Grant wallets:approve in the permission seed
- `de3c4fa` Mount B2B routes under /admin and resolve them to the dealers module
- `a00f69d` Require an idempotency key for admin refunds
- `c27c3b5` Write Gürsel GPS provider configurations only when they are created

Aralık içinde ayrıca iki commit daha vardır: `ccb13e2` (Keep panel-managed role permissions across seed runs) ve `f31c571` (Name new tests Test{Method}_{Scenario} and assert with testify). Sahibin test sonuç listesi `f286fd0` ağacına aittir. Yukarıdaki commit'lerin test sonuçları bu belgede yoktur (Bölüm 12).

### 1.3 Kaynak belgeler

| Belge | Konum | Bu belgedeki kullanımı |
|---|---|---|
| Rev.06 analiz ve iş planı | `/workspace/aystek-core-dev-docs/AYPASS_GECIS_SISTEMI_REV06_ANALIZ_VE_IS_PLANI_TR.md` | Gereksinim kimlikleri, Faz 0 durumu, AD kararları, rol matrisi |
| Rev.06 Excel eşlemesi | `/workspace/aystek-core-dev-docs/AYPASS_REV06_EXCEL_ESLEME_TR.md` | 104 kabul testinin kod durumu (`f48748f` üzerinde hesaplandı) |
| Kurulum ve işletim belgeleri | `user-guide/05-operations/` altında `SUNUCU_KURULUM_REHBERI.md`, `OPERATIONS_GUIDE.md`, `DEPLOYMENT_SCENARIOS.md` | İşletim notları |
| Installer belgeleri | `/workspace/go/aystek-core-installer/README.md`, `docs/KULLANIM_KILAVUZU.md` | Seed kapsayıcısı ve kurulum adımları (commit edilmemiş) |

### 1.4 Sahip kararları

Aşağıdaki kararlar sahibin (Hüseyin Koç) kararlarıdır. Anlamları değiştirilmeden kaydedilmiştir.

- AD-04 onay kayıtları: kullanılmayan modül silinir ve yeniden tasarlanır. Onaylandı.
- Ödeme iadesi eşiği ve B2B zinciri kod değil, politika verisidir. Rol anahtarları İngilizcedir ve değişmezdir. Etiketler dile göredir. Sahibin yönlendirmesi.
- Banka havalesi ve manuel sefer onayları politika değil, izin artı kendi kendine onay kontrolü olarak kalır. Sahip bunu kabul etti.
- Seed test kimlik bilgileri kodda kalır ve üretimde değiştirilecektir. Sahip kararı.
- Editör rolünün menü görünürlüğü değişikliği sahip tarafından onaylandı.
- Ana dil incelemesi, bir ülke projesi kapsama girene ya da iş tamamlanana kadar ertelenir. Sahip kararı.
- Sistem rolleri `editor`, `reviewer` ve `fleet_viewer` oluşturulduktan sonra panelden yönetilir. Seed'ler bunların izinlerinin üzerine yazmaz. Karar 2026-10-28 tarihinde, sahibin "doğru olanı yap" yönergesi üzerine asistan tarafından verildi ve bu işte uygulandı.
- Gürsel GPS sağlayıcı yapılandırması (2026-10-28): yalnızca ilk oluşturmada yazılır. Uygulandı: `c27c3b5`. Ayrıntı Bölüm 9.4.

---

## 2. Mimari genel bakış

### 2.1 Bileşenler

Onay sistemi üç parçadan oluşur: kapı kataloğu (kod), politika verisi (veritabanı, yönetici API'si üzerinden yönetilir) ve onay talepleri (veritabanı). Tüketen modüller onay servisini bir arayüz üzerinden çağırır. Kendi onay mantıklarını yazmazlar.

```
 Panel (aystek-core-next)                       Core API (aystek-core, Go)
 ─────────────────────────                      ──────────────────────────────────────────────
 /admin/approval-policies ──REST──►  GET  /admin/approval-gates
 /admin/approval-requests ──REST──►  GET  /admin/approval-policies/:gate
 Sipariş sayfası iade bildirimi ──►  PUT  /admin/approval-policies/:gate
                                     GET  /admin/approval-requests
                                     POST /admin/approval-requests/:id/approve
                                     POST /admin/approval-requests/:id/reject
                                                    │
                       ┌────────────────────────────┴──────────────────────────┐
                       │  approvals modülü                                      │
                       │  Authorize · Claim · Release · ApprovalChain · Decide  │
                       └──────┬───────────────────┬───────────────────┬─────────┘
                              │                   │                   │
              payments/intents │       b2b/dealers │        returns (RefundWithApproval
              (RefundApprover) │      (ApprovalChain)│         üzerinden payments/intents)
```

Bu şemada kapılar kendi yetkilerini de taşır. Diyagram yalnızca mantıksal bağlantıları gösterir. Yetki kontrolleri Bölüm 4'te rota bazında verilmiştir.

### 2.2 Modül yapısı

`internal/modules/approvals` altındaki dosyalar:

| Dosya | Görev |
|---|---|
| `model/gates.go` | Kapı kataloğu, `DefaultPolicy`, izin anahtarı biçim kontrolü |
| `model/model.go` | `Policy`, `Request`, `Step`, durum sabitleri, hata değişkenleri, sayfalama |
| `dto/dto.go` | İstek ve yanıt gövdeleri (`PolicyRequest`, `DecisionRequest`, `GateResponse`) |
| `repository/interface.go`, `mongo.go`, `indexes.go` | Veritabanı erişimi, koşullu güncellemeler, dizinler |
| `service/service.go` | Politika, `Authorize`, `Claim`, `Release`, `Approve` / `Reject`, `ApprovalChain`, karar günlüğü |
| `handler/admin.go` | Yönetici uçları, hata-HTTP eşlemesi (`statusFor`) |
| `routes.go` | Uç noktaların yetki kapıları |

### 2.3 Tasarım ilkeleri

- Kapı anahtarları kodda sabittir (`payments.refund`, `b2b.purchase_order`). Eşik, etkinlik, para birimi, denetçi izni ve zincir adımları veridir ve yönetici API'si ile değiştirilir.
- Politikalar kiracı (tenant) ve site bazındadır.
- Onay talebi oluşturulurken gereken izin (`required_permission`) talebe kopyalanır (snapshot). Sonradan yapılan politika değişikliği bekleyen bir kararın denetçi kümesini değiştirmez.
- Bekleyen talepler kapı, konu (`subject_ref`) ve tutar bazında tekildir.
- Onay talebi tek kullanımlıktır. Tüketim koşullu bir güncellemedir.
- Kendi talebini onaylamak reddedilir.
- Zincir gerektiren kapılar kapalı başarısız (fail closed) çalışır: politika yoksa işlem durdurulur.
- Kararlar ve politika değişiklikleri etkinlik günlüğüne yazılır.

### 2.4 Entegrasyon özeti

| Tüketen modül | Arayüz veya kontrol | Kapı | Ayrıntı |
|---|---|---|---|
| `payments/intents` | `RefundApprover` (`Authorize`, `Release`) | `payments.refund` | Bölüm 7.1 |
| `returns` | `RefundWithApproval` üzerinden `payments/intents` | `payments.refund` (dolaylı) | Bölüm 7.2 |
| `service_reservations` | İade çağrısı (onay kimliği yok) | `payments.refund` (dolaylı) | Bölüm 7.3 |
| `b2b/dealers` | `ApprovalChain` | `b2b.purchase_order` | Bölüm 7.4 |
| `wallets` | Doğrudan izin kapısı (`wallets:approve`) | Politika dışı | Bölüm 7.5 |
| `fleet_tally/trip_assignments` | Doğrudan izin kapısı (`approve`) | Politika dışı | Bölüm 7.6 |
| `access_points` | Etkin durum kontrolü | Politika dışı | Bölüm 7.7 |
| `billing/invoices` | Hata sözleşmesi | Politika dışı | Bölüm 7.8 |
| `stock/vehicles` | Plaka geçmişi gerekçesi (yalnızca plaka değişince zorunlu) | Politika dışı | Bölüm 4.3.9 |
| `permissions`, `menus`, `seeds` | Rol, menü ve seed kuralları | Politika dışı | Bölüm 10 |

---

## 3. Onay kapıları

### 3.1 Kapı kataloğu

`internal/modules/approvals/model/gates.go` dosyasında iki kapı tanımlıdır. Katalog kod içindedir ve veritabanında saklanmaz.

| Anahtar | Modül | Etiket anahtarı | Zincir | Varsayılan politika |
|---|---|---|---|---|
| `payments.refund` | `payments` | `approvals.gate.payments_refund` | Hayır | Kapalı, izin `payments:approve` |
| `b2b.purchase_order` | `b2b_dealers` | `approvals.gate.b2b_purchase_order` | Evet | Kapalı, adım yok |

Varsayılan politika kodu (`gates.go` 46-55):

```go
// DefaultPolicy is the behaviour of a gate that has no stored Policy: the gate
// is off, so existing operations are unchanged until an administrator enables
// it. The checker permission follows the gate's module convention.
func DefaultPolicy(g Gate) Policy {
	return Policy{
		Gate:               g.Key,
		Enabled:            false,
		RequiredPermission: g.Module + ":approve",
	}
}
```

Kod varsayılanında zincir yoktur. Bu nedenle politikası olmayan kiracı ve site için zincir kapalı başarısız çalışır. Seed'in varsayılan kiracı için eklediği zincir Bölüm 3.5'te ayrıca açıklanmıştır.

### 3.2 Politika alanları

`approval_policies` koleksiyonundaki politika belgesi (`model.Policy`):

| Alan | Tür | Açıklama |
|---|---|---|
| `gate` | metin | Kapı anahtarı |
| `enabled` | bool | Kapı etkin mi |
| `threshold` | sayı, isteğe bağlı | Tek eşikli kapılarda eşik. Eşik üstü işlemler onay ister |
| `currency` | metin, isteğe bağlı | Eşiğin para birimi. Karşılaştırma büyük/küçük harf duyarsızdır |
| `required_permission` | metin, isteğe bağlı | Denetçinin sahip olması gereken izin (`modül:eylem`) |
| `steps` | dizi, isteğe bağlı | Zincir adımları (`approver_role`, `threshold_source`, `threshold`) |
| `updated_by` | metin, isteğe bağlı | Son güncelleyen |

Adım yapısı (`model.Step`, `model.go` 98-102):

| Alan | Değerler | Anlamı |
|---|---|---|
| `approver_role` | rol anahtarı | Adımı karar veren rol |
| `threshold_source` | `amount` veya `dealer_credit_limit` | Adımın hangi değere bakacağı |
| `threshold` | sayı | `amount` için tutar eşiği. `dealer_credit_limit` için 0 olarak gönderilir |

### 3.3 Doğrulama kuralları

`PutPolicy` ve zincir çözümlemesi aşağıdaki kuralları uygular. Hata değerleri çeviri anahtarıdır.

| Koşul | Hata anahtarı | HTTP |
|---|---|---|
| Kapı kataloğunda yok | `approvals.error.gate_unknown` | 404 |
| İşlemi yapan kimlik boş | `approvals.error.actor_required` | 400 |
| Tek eşikli kapıya adım gönderildi | `approvals.error.step_invalid` | 400 |
| İzin boş değil ve `modül:eylem` biçiminde değil | `approvals.error.permission_unknown` | 400 |
| İzin biçimi doğru ama izin kataloğunda yok | `approvals.error.permission_unknown` | 400 |
| Negatif eşik | `approvals.error.invalid_threshold` | 400 |
| Eşik var, para birimi yok | `approvals.error.currency_required` | 400 |
| Zincir kapısında adım yok | `approvals.error.steps_required` | 400 |
| Adımda boş rol | `approvals.error.step_invalid` | 400 |
| Adımda tanımsız rol | `approvals.error.role_unknown` | 400 |
| Adım kaynağı bilinmiyor | `approvals.error.step_invalid` | 400 |
| Tutar kaynağında negatif eşik | `approvals.error.invalid_threshold` | 400 |

Notlar:

- Tek eşikli kapıda izin boş gönderilirse `modül:approve` kullanılır (`fillSingleGate`).
- Zincir kapısında `required_permission` saklanmaz. `PutPolicy` bu alanı yok sayar.
- Tek eşikli kapı etkin ama eşiksizse her işlem onay gerektirir. Eşik boşken koşulsuz bir atlama yoktur (`Authorize`, `service.go` 286).

### 3.4 Zincir çalışma zamanı

`ApprovalChain`, zincir kapısı için çalışma anında adım listesi üretir:

- Politika saklı değilse, devre dışıysa veya adımı yoksa `approvals.error.chain_not_configured` döner. Bu, kapalı başarısız davranıştır. Admin uçlarında HTTP 409 olarak eşlenir. B2B gönderiminde ise 400 döner (`po_admin.go` 130-131).
- Adım seçimi: `amount` adımı, tutar eşiği aşıldığında uygulanır. `dealer_credit_limit` adımı, bayinin kredi limiti 0'dan büyükse ve tutar kredi limitini aşıyorsa uygulanır. Bu adımın eşiği bayinin kredi limitidir.
- Hiçbir adım uygulanmazsa zincir boş döner. Satın alma siparişi bu durumda gönderildiği anda onaylanır (Bölüm 6.2).

### 3.5 Varsayılan zincir durumu

Bu bölüm B2B zincirinin varsayılan durumunu açıklar. Sahibin "b2b.purchase_order için varsayılan zincir yok" ifadesi bu kiracı için düzeltilmiştir: `approvals_seed` varsayılan kiracının (`aystek`) varsayılan sitesi için etkin üç adımlı B2B zinciri ekler. Politikası kayıtlı olmayan veya devre dışı tenant ve site kapalı başarısız olur.

- Kod varsayılanında (`DefaultPolicy`) zincir yoktur. Politikası kayıtlı olmayan veya devre dışı kalan tenant ve site için zincir kapalı başarısız çalışır.
- `internal/seeds/approvals_seed.go` dosyası, varsayılan kiracının (`aystek`, `internal/seeds/util.go` 22) varsayılan sitesi için etkin bir üç adımlı zincir ekler:

```go
{
	Gate:    approvalsmodel.GateB2BPurchaseOrder,
	Enabled: true,
	Steps: []approvalsmodel.Step{
		{ApproverRole: "store_manager", ThresholdSource: approvalsmodel.StepSourceAmount, Threshold: 50000},
		{ApproverRole: "credit_dept", ThresholdSource: approvalsmodel.StepSourceDealerCreditLimit},
		{ApproverRole: "executive", ThresholdSource: approvalsmodel.StepSourceAmount, Threshold: 100000},
	},
},
```

Kaynak: `approvals_seed.go` 26-35. Sonuç: "varsayılan zincir yok" ifadesi aystek kiracısı için geçerli değildir. Politikası kayıtlı olmayan veya devre dışı tenant ve site kapalı başarısız olur. Ayrıca bu üç rolü hiçbir seed oluşturmaz. Roller yönetici tarafından oluşturulmadan adımlar onaylanamaz ve zincirli siparişler bekleyen kalır.

### 3.6 Ödeme iadesi kapısı varsayılanı

`payments.refund` seed'de kapalı gelir (`approvals_seed.go` 36-40). İzin `payments:approve` olarak tanımlanır. Eşik ve para birimi yoktur. Kapı açıldığında eşiksiz olarak her iade onay ister. Bu, Bölüm 3.2'de tanımlanan kuralın doğal sonucudur. Yönetici bunu bilmelidir.

---

## 4. API referansı

### 4.1 Genel kurallar

- Yönetici uçları `/admin` yoludu altındadır (`internal/platform/router/groups.go` 33). Onay uçları `approvals.RegisterRoutes` ile bu gruba bağlanır (`cmd/api/main.go` 2502).
- Her onay ucu kendi `ModulePermissionMiddleware(modül, eylem)` kapısını taşır (`routes.go` 18-23).
- Platform kullanıcıları modül kapısını atlar. Kapı yalnızca platform dışı kullanıcılara uygulanır (`module_permission.go` 26-28).
- Başarılı yanıtlar zarf içindedir: `{"success": true, "data": ..., "pagination": ...}`. Liste yanıtlarında `pagination` alanı `page`, `limit`, `total`, `total_pages`, `has_next` ve `has_prev` alanlarını taşır (`pagination.go`).
- Hata yanıtlarında handler `gin.H{"error": ...}` yazar. Zarf katmanı `error` değerini köke taşır: `{"success": false, "error": "<anahtar>"}`. Handler'ın `error` dışına yazdığı ek alanlar (ör. `approval_request_id`) köke korunur.

Kaynak: `internal/platform/middleware/response_envelope.go` 120-170 (hata gövdesinin köke taşınması ve ek alanların korunması) ve 183-184 (ek alanların birleştirilmesi).

Bu bölümdeki JSON örnekleri gösterim amaçlıdır. Kimlik, tarih ve tutar değerleri uydurulmuştur. Alan adları ve yapılar kod okumasından alınmıştır.

### 4.2 Onay yönetimi uçları

| Yöntem | Yol | Yetki | İşleyici | Başarı | Hatalar |
|---|---|---|---|---|---|
| GET | `/admin/approval-gates` | `approval_policies:read` | `ListGates` | 200 | 500 (ham metin) |
| GET | `/admin/approval-policies/:gate` | `approval_policies:read` | `GetPolicy` | 200 | 404 `gate_unknown`; 400 ham metin |
| PUT | `/admin/approval-policies/:gate` | `approval_policies:update` | `PutPolicy` | 200 | 404 `gate_unknown`; 400 doğrulama anahtarları |
| GET | `/admin/approval-requests` | `approval_requests:read` | `ListRequests` | 200 | 500 (ham metin) |
| POST | `/admin/approval-requests/:id/approve` | `approval_requests:read` (rota) | `Approve` | 200 | 404, 403, 409, 400 |
| POST | `/admin/approval-requests/:id/reject` | `approval_requests:read` (rota) | `Reject` | 200 | 404, 403, 409, 400 |

Önemli not: approve ve reject uçlarının rota kapısı `approval_requests:read` iznidir (`routes.go` 22-23). Asıl denetim servis katmanındadır. Talebe kaydedilmiş denetçi izni (`required_permission`) karar verenin oturumunda bulunmalıdır. Bu tasarım, yalnızca okuma izni olan bir kullanıcının karar ucuna ulaşabilmesi anlamına gelir. Denetçi izni olmayan biri servis katmanında `approvals.error.permission_missing` alır.

#### 4.2.1 GET /admin/approval-gates

İstek: gövde yok.

Yanıt (zarfın `data` kısmı; kısaltılmış):

```json
{
  "success": true,
  "data": [
    {
      "key": "payments.refund",
      "module": "payments",
      "label_key": "approvals.gate.payments_refund",
      "chain": false,
      "stored": false,
      "enabled": false,
      "required_permission": "payments:approve"
    },
    {
      "key": "b2b.purchase_order",
      "module": "b2b_dealers",
      "label_key": "approvals.gate.b2b_purchase_order",
      "chain": true,
      "stored": true,
      "enabled": true,
      "steps": [
        {"approver_role": "store_manager", "threshold_source": "amount", "threshold": 50000}
      ],
      "updated_by": "seed",
      "updated_at": "2026-10-28T09:00:00Z"
    }
  ],
  "pagination": {"page": 1, "limit": 2, "total": 2, "total_pages": 1, "has_next": false, "has_prev": false}
}
```

Liste sayfalı değildir. Katalogdaki tüm kapılar döner. `steps` dizisi örnekte kısaltılmıştır. Gerçek zincirde üç adım bulunur (Bölüm 3.5).

#### 4.2.2 GET /admin/approval-policies/:gate

Yol parametresi: `gate` (ör. `payments.refund`).

Yanıt: `dto.GateResponse`. Alanlar: `key`, `module`, `label_key`, `chain`, `stored`, `enabled`, `threshold`, `currency`, `required_permission`, `steps`, `updated_by`, `updated_at`. Politika saklı değilse `stored` değeri `false` döner ve diğer alanlar kapı varsayılanıdır.

Hata: bilinmeyen kapı için 404 ve gövde `{"success": false, "error": "approvals.error.gate_unknown"}`.

#### 4.2.3 PUT /admin/approval-policies/:gate

İstek gövdesi (`dto.PolicyRequest`, `dto.go` 11-17):

| Alan | Tür | Not |
|---|---|---|
| `enabled` | bool | `omitempty` yoktur. Gönderilmezse `false` olur ve kapı kapanır |
| `threshold` | sayı, isteğe bağlı | Tek eşikli kapılar için |
| `currency` | metin, isteğe bağlı | Eşik varsa zorunludur |
| `required_permission` | metin, isteğe bağlı | Tek eşikli kapılar için |
| `steps` | dizi, isteğe bağlı | Zincir kapıları için zorunludur |

Tek eşikli kapı örneği (`payments.refund`):

```json
{
  "enabled": true,
  "threshold": 500,
  "currency": "TRY",
  "required_permission": "payments:approve"
}
```

Zincir örneği (`b2b.purchase_order`):

```json
{
  "enabled": true,
  "steps": [
    {"approver_role": "store_manager", "threshold_source": "amount", "threshold": 50000},
    {"approver_role": "credit_dept", "threshold_source": "dealer_credit_limit", "threshold": 0},
    {"approver_role": "executive", "threshold_source": "amount", "threshold": 100000}
  ]
}
```

Başarı yanıtı: `dto.GateResponse`, zarf içinde `"success": true`. Her başarılı güncelleme etkinlik günlüğüne `approvals.activity.policy_updated` olarak yazılır (`service.go` 157-160).

PUT, politikayı tamamen değiştirir. Gönderilmeyen alan varsayılan değerine döner. Bu bir işletim notudur (Bölüm 13.3).

#### 4.2.4 GET /admin/approval-requests

Sorgu parametreleri (`handler/admin.go` 75-78, `model.go` 126-141):

| Parametre | Varsayılan | Sınır | Not |
|---|---|---|---|
| `status` | boş | Serbest metin | Değer izin listesiyle doğrulanmaz |
| `page` | 1 | 1'den küçükse 1 | Ayrıştırma hatası sessizce varsayılana düşer |
| `limit` | 50 | 1'den küçükse 50, 200'den büyükse 200 | Ayrıştırma hatası sessizce varsayılana düşer |

Sıralama: `created_at` azalan. Sonuç yoksa `data` boş dizidir.

Örnek yanıt (zarfın tamamı, kısaltılmış):

```json
{
  "success": true,
  "data": [
    {
      "id": "6720c0a1f0e4b2a1d3c4e5f6",
      "tenant_id": "aystek",
      "site_id": "site-1",
      "gate": "payments.refund",
      "subject_ref": "66f0a1b2c3d4e5f607080910",
      "amount": 750,
      "currency": "TRY",
      "requested_by": "66ee00000000000000000001",
      "required_permission": "payments:approve",
      "status": "pending",
      "created_at": "2026-10-28T10:15:00Z",
      "updated_at": "2026-10-28T10:15:00Z"
    }
  ],
  "pagination": {"page": 1, "limit": 50, "total": 1, "total_pages": 1, "has_next": false, "has_prev": false}
}
```

#### 4.2.5 POST /admin/approval-requests/:id/approve ve /reject

İstek gövdesi (`dto.DecisionRequest`, `dto.go` 20-22): isteğe bağlı `note` alanı. Gövde boş olabilir (`admin.go` 104).

```json
{
  "note": "Tutar ve gerekçe uygun."
}
```

Servis kontrolleri, sırayla (`service.go` 351-395):

1. İşlemi yapan kimlik boş değilse devam edilir. Boşsa `approvals.error.actor_required`.
2. Talep bulunur. Bulunamazsa `approvals.error.request_not_found` (404).
3. Talep `pending` durumunda olmalıdır. Değilse `approvals.error.request_not_pending` (409).
4. Karar veren, talebi açan kişi olamaz. Ise `approvals.error.self_approval` (403).
5. Karar veren, talebe kaydedilmiş `required_permission` iznine sahip olmalıdır. Yoksa `approvals.error.permission_missing` (403).
6. Koşullu güncelleme yalnızca `pending` durumdaki talebi değiştirir. Eşzamanlı karar durumunda ikinci karar `request_not_pending` alır.

Başarı yanıtı: güncellenmiş talep nesnesi. Karar, `approvals.request.approve` veya `approvals.request.reject` eylemi ile `approvals.activity.request_approved` ya da `approvals.activity.request_rejected` açıklamasıyla etkinlik günlüğüne yazılır (`service.go` 384-393).

Başarılı onay örneği (zarfın `data` kısmı):

```json
{
  "id": "6720c0a1f0e4b2a1d3c4e5f6",
  "gate": "payments.refund",
  "subject_ref": "66f0a1b2c3d4e5f607080910",
  "amount": 750,
  "currency": "TRY",
  "requested_by": "66ee00000000000000000001",
  "required_permission": "payments:approve",
  "status": "approved",
  "decided_by": "66ee00000000000000000002",
  "decided_at": "2026-10-28T10:20:00Z",
  "decision_note": "Tutar ve gerekçe uygun.",
  "created_at": "2026-10-28T10:15:00Z",
  "updated_at": "2026-10-28T10:20:00Z"
}
```

Hata örneği (kendi talebini onaylama):

```json
{"success": false, "error": "approvals.error.self_approval"}
```

`invalid id` metni bir çeviri anahtarı değildir. Geçersiz kimlikte handler bu düz metni döndürür. Bu bilinen bir istisnadır (Bölüm 14).

### 4.3 Değişen uçlar

Aşağıdaki uçların çoğunun yolu bu iş öncesinde de vardı. B2B uçlarının taban commit'teki `/api` önekli yolu `de3c4fa` ile `/admin/b2b` altına taşındı. Tabloda "Taban" sütunu, `git show f48748f:<dosya>` ile karşılaştırılan başlangıç durumunu verir. "Yeni" ifadesi yalnızca taban commit'te hiç bulunmayan uçlar için kullanılır.

| Yöntem | Yol | Taban durumu | Son durum | Yetki | Değişiklik |
|---|---|---|---|---|---|
| POST | `/admin/payment-intents/:id/refunds` | Yetki kapısı yok (`routes.go` 20) | Kapı `payments:execute` | `payments:execute` | `approval_request_id` alanı; `idempotency_key` zorunlu; 400 `payment_intents.idempotency_key_required` |
| PATCH | `/admin/returns/:id/refund` | Rota kapısı yok | Rota kapısı yok; genel `returns:*` kontrolü | `returns:*` | 409 `approvals.error.required` ve `not_granted` |
| POST | `/admin/b2b/purchase-orders/:id/approve` | Kapı yok (`routes.go` 70) | `b2b_dealers:approve` | `b2b_dealers:approve` | Sıralı zincir; rol oturumdan okunur |
| POST | `/admin/b2b/purchase-orders/:id/reject` | Kapı yok (`routes.go` 71) | `b2b_dealers:approve` | `b2b_dealers:approve` | Mevcut adımın rolü gerekir |
| POST | `/admin/b2b/purchase-orders/:id/submit` | Mevcut | Mevcut | Zincir koruması | Zincir politikadan oluşur; yoksa 400 |
| POST | `/admin/wallets/:id/bank-transfers/:journalId/confirm` | Kapı yok (`routes.go` 25) | `wallets:approve` | `wallets:approve` | 403 `bank_transfer_self_confirm` |
| POST | `/me/wallet/bank-transfer-declarations` | Mevcut (`routes.go` 39) | Mevcut | Oturum (modül anahtarı yok) | Bekleyen defter kaydı; bakiye değişmez |
| POST | `/admin/fleet-tally/trip-assignments/:id/review` | `execute` (`routes.go` 19) | `approve` | `fleet_tally_trip_assignments:approve` | Bağımsız onaylayan |
| GET | `/admin/fleet-tally/trip-assignments/pending-approval` | `read` (`routes.go` 18) | Değişmedi | `fleet_tally_trip_assignments:read` | Değişiklik yok |
| POST | `/admin/invoices/:id/issue` | Mevcut (`routes.go` 16) | Mevcut | `billing:*` (genel) | 409 `invoices.already_issued`; hata sözleşmesi |
| GET | `/admin/invoices/:id/eissue/status` | Mevcut (`routes.go` 21) | Mevcut | `billing:*` (genel) | Hata 400 `{error, invoice}` |
| POST | `/admin/invoices/:id/eissue/retry` | Mevcut (`routes.go` 22) | Mevcut | `billing:*` (genel) | Aynı hata eşlemesi |
| POST | `/admin/invoices/bulk/eissue` | Mevcut (`routes.go` 24) | Mevcut | `billing:*` (genel) | Başarısız öğeler `status: failed` alır |
| POST, PUT, DELETE | `/admin/roles`, `/admin/roles/:id` | Mevcut | Mevcut | `roles:*` (genel) | Anahtar doğrulaması, değişmezlik, etiketler |
| PUT | `/admin/stock/vehicles/:id` | Taban: gerekçe alanı yoktu (`93b4232` ile eklendi) | `stock_vehicles:update` | `stock_vehicles:update` | `plate_reason`: plaka değişince zorunlu (Bölüm 4.3.9) |
| POST | `/admin/stock/vehicles/:id/plates` | Taban: gerekçe alanı yoktu (`93b4232` ile eklendi) | `stock_vehicles:update` | `stock_vehicles:update` | `plate_reason`: plaka değişince zorunlu (Bölüm 4.3.9) |
| Menü yönetici uçları (yol bu belgede ayrıca doğrulanmadı) | Mevcut; `routes.go` değişmedi | Mevcut | Menü yetkisi | Görünürlük izne göre; `permissions` yönetici API'sinden yazılamaz |
| Erişim noktası olay uçları | Mevcut | Mevcut | Mevcut | Mevcut | Pasif erişim noktası reddedilir |

Notlar:

- Onay uçlarının tamamı (`/admin/approval-gates`, `/admin/approval-policies/:gate`, `/admin/approval-requests`, `.../approve`, `.../reject`) bu işte yeni eklenmiştir. Taban commit'te karşılıkları yoktur.
- B2B uçlarının taban commit'teki `/api` öneki `de3c4fa` ile kaldırıldı. Bu önekin yetki sonucu (modül `api`, platform dışı kullanıcılara 403) düzeltildi; Bölüm 14, birinci risk maddesine bakın.
- Erişim noktası olay uçlarının yolu bu belgede ayrıca listelenmemiştir. Davranış Bölüm 7.7'de anlatılır.

#### 4.3.1 Ödeme iadesi: POST /admin/payment-intents/:id/refunds

İstek (`payments/intents/dto/dto.go` 62-69):

| Alan | Tür | Not |
|---|---|---|
| `amount` | sayı | Zorunlu, sıfırdan büyük |
| `reason` | metin | İsteğe bağlı |
| `idempotency_key` | metin | Zorunlu. Sırasıyla gövde, `Idempotency-Key` başlığı, `X-Idempotency-Key` başlığı okunur; üçünde de yoksa 400 |
| `approval_request_id` | metin, isteğe bağlı | Onaylanmış `payments.refund` talebinin kimliği (hex) |

Anahtar öncelik sırası: gövde, `Idempotency-Key` başlığı, `X-Idempotency-Key` başlığı. Üçü de boşsa istek reddedilir. Servis, para hareketinden ve onay kontrolünden önce `payment_intents.idempotency_key_required` ile 400 döner. Yönetici işleyicisi bu durumda yeni bir anahtar üretmez (`a00f69d`). Bu nedenle anahtarsız istek artık kabul edilmez ve anahtarsız ağ tekrarı ikinci bir iade yaratamaz.

Örnek istek:

```json
{
  "amount": 750,
  "reason": "Hizmet iptali",
  "idempotency_key": "refund-2026-10-28-0001",
  "approval_request_id": "6720c0a1f0e4b2a1d3c4e5f6"
}
```

Yanıtlar:

| Durum | Gövde | Kaynak |
|---|---|---|
| 200 | Ödeme niyeti (`refunded_amount`, `status`). `approval_request_id` yoktur | `admin.go` 46 |
| 409 | `{"success": false, "error": "approvals.error.required", "approval_request_id": "<hex>"}` | `admin.go` 181 |
| 409 | `{"success": false, "error": "approvals.error.not_granted"}` | `admin.go` 184-186 |
| 400 | `payment_intents.idempotency_key_required` | `service.go` 537; `admin.go` (`idempotencyErrorStatus` varsayılanı 400) |
| 422 | `payment_intents.idempotency_conflict` | `service.go` 38 |
| 409 | `payment_intents.idempotency_in_progress` | `service.go` 40 |
| 400 | `invalid approval_request_id` (düz metin), `payment_intents.refund_not_allowed`, `payment_intents.refund_exceeds_capture` | `admin.go` 172; `service.go` 32, 34 |

Eşik üstü bir iade ilk istekte 409 `approvals.error.required` döndürür ve `approval_request_id` verir. Denetçi onayladıktan sonra aynı anahtar ve aynı `approval_request_id` ile yapılan istek talebi tüketerek iadeyi gerçekleştirir.

Tutar eşik altındaysa veya kapı kapalıysa talep oluşturulmaz ve iade doğrudan yapılır (`Authorize`, `service.go` 281-288).

#### 4.3.2 Returns iadesi: PATCH /admin/returns/:id/refund

İstek (`returns/dto/dto.go` 81-89): `amount`, `currency`, `reason`, `payment_intent_id` ve isteğe bağlı `approval_request_id`.

Yanıtlar:

| Durum | Gövde | Kaynak |
|---|---|---|
| 200 | İade yanıtı | `returns/handler/admin.go` 118-156 |
| 409 | `approvals.error.required` ve `approval_request_id` | `admin.go` 147-150 |
| 409 | `approvals.error.not_granted` | `admin.go` 152-155 |
| 400 | `returns.error.refund_not_allowed`, `returns.error.refund_invalid_amount`, `returns.error.invalid_payload` | `service.go` 21-25 |

Servis davranışı (`returns/service/service.go` 293-377): Ödeme niyeti iadesi `return:<return_id>` anahtarıyla çağrılır. Onay bekleyen veya kullanılamayan onay durumunda (`awaitsApproval`) iadenin durumu yazılmaz ve iade `failed` olarak işaretlenmez (`service.go` 353-355). Diğer hatalarda `failed` yazılır. Başarıda `refunded` yazılır.

Bu, returns iadesinin de `payments.refund` kapısından geçtiği anlamına gelir.

#### 4.3.3 B2B satın alma siparişleri

Tüm uçların yolu `/admin/b2b/...` biçimindedir. Taban commit'teki `/api` önekli yollar `de3c4fa` ile değiştirildi. Yetki kapıları aşağıdadır (`b2b/dealers/routes.go` 36-80).

| Yöntem | Yol | Kapı | İşleyici | Başarı |
|---|---|---|---|---|
| POST | `/admin/b2b/purchase-orders` | Zincir koruması | `CreatePurchaseOrder` | 201 |
| GET | `/admin/b2b/purchase-orders` | Zincir koruması | `ListPurchaseOrders` | 200 |
| GET | `/admin/b2b/purchase-orders/:id` | Zincir koruması | `GetPurchaseOrder` | 200 |
| PUT | `/admin/b2b/purchase-orders/:id` | Zincir koruması | `UpdatePurchaseOrder` | 200 |
| POST | `/admin/b2b/purchase-orders/:id/submit` | Zincir koruması | `SubmitPurchaseOrder` | 200 |
| POST | `/admin/b2b/purchase-orders/:id/approve` | `b2b_dealers:approve` | `ApprovePurchaseOrder` | 200 |
| POST | `/admin/b2b/purchase-orders/:id/reject` | `b2b_dealers:approve` | `RejectPurchaseOrder` | 200 |
| POST | `/admin/b2b/purchase-orders/:id/cancel` | Zincir koruması | `CancelPurchaseOrder` | 200 |
| POST | `/admin/b2b/purchase-orders/:id/convert` | Zincir koruması | `ConvertPurchaseOrder` | 200 |

"Zincir koruması" ifadesi, rotaya `protectedChain` geçildiği anlamına gelir (`cmd/api/main.go` 2519). Bu zincir `PermissionMiddleware` içerir. Modül yoldan `b2b` olarak çıkarılır ve `b2b_dealers` takma adına çözülür (`module_gate.go` 143). Önceki `/api` önekinde modül `api` çözülüyordu ve bu düzeltildi. Approve ve reject uçlarındaki `b2b_dealers:approve` kapısı, bu genel kontrolden sonra çalışır.

Onay ve red gövdeleri:

- Onay (`dto.PurchaseOrderApprovalRequest`, `dto.go` 411-413): `notes`. Kod okumasına göre gövde boş olamaz (`ShouldBindJSON`).
- Red (`dto.PurchaseOrderRejectRequest`, `dto.go` 416-418): `reason`, zorunlu.

Örnekler (gösterim amaçlı):

```json
{"notes": "Bayi kredisi kontrol edildi."}
```

```json
{"reason": "Teslim tarihi uygun değil."}
```

Rol ve kimlik kaynağı: Onay veya red için rol istek gövdesinden okunmaz. `callerRoles(c)` yalnızca oturumun `roles` bağlamını okur (`po_admin.go` 323-327). Taban commit'te (`f48748f`) rol gövdeden okunuyordu. `60cd404` bunu kaldırdı.

Sipariş yanıtında `submitted_by` alanı dönmez. Alan modelde saklanır, ancak yanıt modelinde yer almaz.

Hata anahtarları (`b2b/dealers/service/po_service.go` 19-24): `b2b.po.invalid`, `b2b.po.not_found`, `b2b.po.invalid_status`, `b2b.po.approval_required`, `b2b.po.self_approval`, `b2b.po.order_service_not_configured`. Gönderimde `approvals.error.chain_not_configured` olduğu gibi geçer ve 400 döner (`po_admin.go` 130-131).

#### 4.3.4 Banka havalesi

Beyan: `POST /me/wallet/bank-transfer-declarations` (`wallets/routes.go` 39-40). Uç taban commit'te de vardır. Modül anahtarı yoktur. Oturum ve bağlı kişi gerekir.

İstek (`wallets/dto/dto.go` 75-83): `amount`, `reference` (boşsa üretilir), `currency` (varsayılan TRY), `person_id` (isteğe bağlı, vasi hedefi).

```json
{
  "amount": 1500,
  "reference": "EFT-2026-0001",
  "currency": "TRY"
}
```

Yanıt: 201, `{wallet, journal}`. Hatalar: 404 "no person linked to this account", 403 "not your wallet", 404 "wallet not found", 400 (`public.go` 121-159). Beyan anahtarı `bank_transfer:<walletID>:<reference>` biçimindedir; aynı anahtarla yapılan tekrar aynı kaydı döndürür.

Onay: `POST /admin/wallets/:id/bank-transfers/:journalId/confirm` (`routes.go` 26), yetki `wallets:approve`. Gövde yoktur.

Yanıtlar:

| Durum | Gövde | Kaynak |
|---|---|---|
| 200 | Cüzdan nesnesi | `handler/admin.go` 144-165 |
| 400 | `invalid id` veya `invalid journal id` (düz metin) | `admin.go` 147, 152 |
| 403 | `{"success": false, "error": "wallets.error.bank_transfer_self_confirm"}` | `admin.go` 158-160 |
| 400 | Diğer hatalar (İngilizce düz metin) | `admin.go` 157-161 |

Cüzdan nesnesinin alanları (kod okumasına göre): `id`, `tenant_id`, `site_id`, `person_id`, `currency`, `ledger_account_id`, `cached_balance`, `cached_balance_updated_at`, `last_applied_journal_id`, `status`, `closed_at`, `balance_floor_override`, `floor_exempt`, `created_at`, `updated_at`. Onaylayanın kimliği `c.GetString("userID")` ile alınır (`admin.go` 155).

#### 4.3.5 Sefer onayı

Onay: `POST /admin/fleet-tally/trip-assignments/:id/review`, yetki `fleet_tally_trip_assignments:approve`.

İstek (`trip_assignments/dto/dto.go` 72-79):

| Alan | Değerler | Not |
|---|---|---|
| `action` | `approve`, `reject_refund`, `reject_transfer` | İşlem türü |
| `access_point_id` | kimlik | `reject_refund` için zorunlu |
| `transfer_plate` | metin | `reject_transfer` için |
| `transfer_route_id` | kimlik | Boş değilse ObjectID olmalıdır |
| `memo` | metin | Okunmaz |

```json
{
  "action": "reject_refund",
  "access_point_id": "6710a0b1c2d3e4f506070809"
}
```

Bekleyen liste: `GET /admin/fleet-tally/trip-assignments/pending-approval`, yetki `fleet_tally_trip_assignments:read`. Filtre: `approval_status` pending ve `status` completed. Sıralama `started_at` azalan. Varsayılan limit 25 (`admin.go` 206-207, `repository/mongo.go` 274-305).

Bağımsızlık kuralı: oluşturan kişi (`created_by`) boş değilse ve işlemi yapan aynı kişiyse 403 `fleet_tally.trip_assignments.error.self_approval` döner (`admin.go` 275-279). Oluşturan boşsa engel yoktur.

Diğer hatalar (400, İngilizce metin): `ErrTripNotManual`, `ErrTripNotInProgress`, `ErrApprovalActionInvalid` (`service.go` 20-23).

#### 4.3.6 Faturalama

| Yöntem | Yol | Davranış |
|---|---|---|
| POST | `/admin/invoices/:id/issue` | Gövde `{provider, doc_type}`. Zaten kesilmişse 409 `invoices.already_issued`. Kayıtlı sağlayıcı yoksa `invoices.provider_missing` (kayıt yapılmaz). Sağlayıcı boş sonuç dönerse `invoices.provider_empty_result` (400). |
| POST | `/admin/invoices/:id/eissue/retry` | Aynı eşleme. Başarısız fatura yeniden denenebilir. |
| POST | `/admin/invoices/bulk/eissue` | Gövde `invoice_ids`. Başarısız öğeler `status: failed` ve hata metni alır. |
| GET | `/admin/invoices/:id/eissue/status` | Hata 400 `{error, invoice}` biçimindedir. |

Örnek istek (değer alanları gösterim amaçlıdır):

```json
{
  "provider": "<sağlayıcı anahtarı>",
  "doc_type": "<belge tipi>"
}
```

Örnek hata yanıtı:

```json
{"success": false, "error": "invoices.already_issued"}
```

GİB adaptörü hataları: taşıma hatası boş UUID dönmez; `gib.submission_failed` ile sarılır. UUID içermeyen yanıt `gib.invalid_response` olur. Durum sorgusunda `Status` öğesi yoksa `gib.invalid_response` döner. Eski davranış `ACCEPTED` dönüyordu. `IssuedAt` zorunludur (`gib.invalid_invoice`). Belge UUID'si fatura numarası ve `IssuedAt` Unix saniyesinden deterministik üretilir. Tekrar gönderimde değişmez (`gib_einvoice.go` 411-418).

Güncelleme (`33467c8`): `mock` adlı e-invoice ve ERP sağlayıcıları yalnızca `cfg.Env` production değilken kaydedilir (`cmd/api/main.go` 1448-1468). Production ortamında `mock` adıyla yapılan fatura isteği kayıtlı bir sağlayıcı bulamaz ve reddedilir (kod okumasına göre). Bu sağlayıcının `Issue` çağrısı `Status: "sent"` ve `ExternalID: "EINV-<numara>"` ile hatasız döner (`adapters/mock.go` 39-47). Production dışındaki ortamlarda sahte başarı yolu olarak kalır (Bölüm 14, risk 3).

#### 4.3.7 Roller ve izinler

Rol oluşturma ve güncelleme davranışı (`permissions/service/service.go` 99-163, `role_keys.go`):

- Rol anahtarı `^[a-z][a-z0-9_]{1,62}$` desenine uymalıdır (`model/role.go` 31, `role_keys.go` 23-28). Oluşturmada uymazsa `permissions.error.role_key_invalid`.
- Güncellemede ad değişirse `permissions.error.role_key_immutable` (`service.go` 138-140).
- Dil kodu `^[a-z]{2,3}(-[a-z0-9]{2,8})?$` desenine uymalıdır (`model/role.go` 32). Etiket boş olamaz ve 100 karakteri geçemez. Hata anahtarı `permissions.error.role_label_invalid`.
- Sistem rolleri güncellenemez ve silinemez. Döndürülen metin düz metindir: `system role immutable` (`service.go` 135-137, 170-172).
- Etiketler `labels` alanında dil kodu → metin eşlemesi olarak saklanır. Boş `labels` nesnesi etiketleri temizler. Alan hiç gönderilmezse etiketler korunur.
- `permissions` dizisi gönderilirse izinler değiştirilir. Boş dizi izinleri temizler.

Örnek istek (`POST /admin/roles`, değerler gösterim amaçlıdır):

```json
{
  "name": "muhasebe_uzmani",
  "description": "Muhasebe kayıtlarını izler.",
  "labels": {"tr": "Muhasebe uzmanı", "en": "Accounting specialist"},
  "permissions": []
}
```

Örnek hata yanıtı (yanıt yardımcısının biçimi):

```json
{
  "success": false,
  "message": "role_create_failed",
  "error": "permissions.error.role_key_invalid"
}
```

Etiket geri dönüşü (`model/role.go` 50-58): dil kodu, sonra `en`, sonra rol adı. Backend'de bu fonksiyonun çağıranı yoktur. Panel geri dönüşü kendisi uygular (Bölüm 11.5).

#### 4.3.8 Menüler ve erişim noktaları

Menü uçlarının yolu değişmemiştir. Görünürlük filtresi değişmiştir (Bölüm 10.3). `MenuRequest` içinde `permissions` alanı yoktur (`menus/dto/dto.go` 5-16). İzinler yalnızca seed'ler aracılığıyla yazılır.

Erişim noktaları: pasif (`inactive`) bir erişim noktası olay kaydını (`RecordEvent`), taşıma nakit kapanışını (`SettleTripCash`) ve yaya binişi kaydını (`RecordWalkOnCashEntry`) reddeder. Hata değeri `access_points.error.inactive`dir (`access_points/service/service.go` 52). Olay kaydında bu hata 400 olarak döner. Aynı anahtarla tekrar gönderilen olay, durum kontrolünden önce idempotensi denetiminde bulunur ve saklanan olayı döndürür (`service.go` 944-949).

#### 4.3.9 Stok araçları: plaka değişikliği gerekçesi

Plaka değişikliği artık gerekçe ister (`93b4232`). Gerekçe `plate_reason` alanında gönderilir. Taban commit'te bu alan yoktu. Uçlar `stock/vehicles/routes.go` dosyasındadır.

| Yöntem | Yol | Yetki | Gerekçe kuralı |
|---|---|---|---|
| PUT | `/admin/stock/vehicles/:id` | `stock_vehicles:update` | `current_plate` mevcut plakadan farklıysa `plate_reason` zorunlu |
| POST | `/admin/stock/vehicles/:id/plates` | `stock_vehicles:update` | `plate` mevcut plakadan farklıysa `plate_reason` zorunlu |

İstek alanları:

- PUT (`dto.UpdateVehicleRequest`, `dto.go` 30-37): `current_plate`, `name`, `supplier_id` (isteğe bağlı), `status`, `plate_reason`.
- POST (`dto.AddPlateRequest`, `dto.go` 39-43): `plate`, `started_at` (RFC3339, isteğe bağlı), `plate_reason`.

Örnek PUT gövdesi (değerler gösterim amaçlıdır):

```json
{
  "current_plate": "34ABC123",
  "name": "Servis aracı 3",
  "status": "active",
  "plate_reason": "Araç yeniden plakalandı"
}
```

Başarılı yanıt: PUT için `stock.vehicles.notice.saved`, POST için `stock.vehicles.notice.plate_added` mesajı ve zarfın `data` kısmında `{"id": "<hex>"}`.

Hata yanıtları (`handler/admin.go` `writeError`, AddPlate işleyicisi):

| HTTP | Anahtar | Koşul |
|---|---|---|
| 400 | `stock.vehicles.validation.plate_reason_required` | Plaka değişti; gerekçe boş veya yalnızca boşluk |
| 400 | `stock.vehicles.validation.plate_reason_too_long` | Kırpılmış gerekçe 500 karakteri aşıyor |
| 400 | `stock.vehicles.validation.plate_required` | Plaka boş |
| 400 | `stock.vehicles.validation.duplicate_plate` | Plaka başka bir araçta kayıtlı |
| 400 | `stock.vehicles.validation.status_invalid` | Durum geçersiz |
| 400 | `stock.vehicles.validation.invalid_payload` | Gövde bağlanamadı |
| 404 | `stock.vehicles.error.not_found` | Araç bulunamadı |
| 403 | `stock.vehicles.error.forbidden` | Plaka ekleme: araç çağıranın kurumuna ait değil (AddPlate işleyicisi) |
| 500 | err.Error() | Diğer hatalar |

Kural ayrıntıları (`service/service.go`):

- Gerekçe `strings.TrimSpace` ile kırpılır. Boşsa `ErrPlateReasonRequired`. `utf8.RuneCountInString` değeri `maxPlateReasonLength` (500) değerini aşarsa `ErrPlateReasonTooLong` (`service.go` 271, 275-284).
- Plaka değişikliği, plaka normalize edildikten sonra mevcut plakayla karşılaştırılır. Plaka değişmemişse gerekçe istenmez.
- Gerekçe `stock_vehicle_plate_history.reason` alanına yazılır (`model.PlateHistory.Reason`, BSON `reason,omitempty`). Dizin: `stock_vehicle_plate_history__tenant_site_vehicle_started` (`indexes.go` 23-25).
- Bilinen boşluk: geçmiş yazma hataları yutuluyor (Bölüm 14.3, risk 25).

Çeviri anahtarları: `stock.vehicles.validation.plate_reason_required` ve `stock.vehicles.validation.plate_reason_too_long` (`module_error_translations_seed.go`), `stock.vehicles.fields.plate_reason` (`translation_seed.go`). Koordinatör bildirimine göre 28 dilde; ana dil incelemesi bekleyen diller değişmedi.

Panel tarafı: Bölüm 11.9.

### 4.4 Onay hata anahtarları

Kaynak: `internal/modules/approvals/model/model.go` 66-86 ve 111-122. HTTP eşlemesi `handler/admin.go` 135-146 (`statusFor`).

| Anahtar | HTTP | Üreten yer |
|---|---|---|
| `approvals.error.gate_unknown` | 404 | Politika, PUT, zincir, Authorize |
| `approvals.error.request_not_found` | 404 | Karar |
| `approvals.error.self_approval` | 403 | Karar |
| `approvals.error.permission_missing` | 403 | Karar |
| `approvals.error.request_not_pending` | 409 | Karar |
| `approvals.error.not_granted` | 409 | Authorize (talep tüketilemedi) |
| `approvals.error.chain_not_configured` | 409 | Zincir çözümlemesi (B2B gönderiminde 400) |
| `approvals.error.required` | `statusFor` içinde yok; varsayılan 400 | Authorize; iade ve B2B çağıranlarına döner. Refund uçları bunu ayrıca 409 olarak yazar |
| `approvals.error.actor_required` | 400 | PUT, karar |
| `approvals.error.permission_unknown` | 400 | Tek eşikli kapı izni |
| `approvals.error.currency_required` | 400 | Tek eşikli kapı |
| `approvals.error.invalid_threshold` | 400 | Eşik doğrulaması |
| `approvals.error.steps_required` | 400 | Zincir doğrulaması |
| `approvals.error.step_invalid` | 400 | Adım doğrulaması |
| `approvals.error.role_unknown` | 400 | Adım rolü |

Yanıtlarda çeviri anahtarı olmayan metinler de vardır: `invalid id` (düz metin), Go bağlama hata metinleri ve `statusFor` varsayılanındaki ham hata metinleri. Bu metinler panelde olduğu gibi gösterilir (Bölüm 11.7).

`RequiredError` tipi (`model.go` 115-122) iade ve B2B çağıranlarında 409 yanıtı üretmek için kullanılır. Yönetici uçları bu hatayı kendisi üretmez.

### 4.5 Diğer modüllerin hata anahtarları

| Modül | Anahtarlar |
|---|---|
| Ödeme niyetleri | `payment_intents.refund_exceeds_capture`, `payment_intents.idempotency_key_required`, `payment_intents.idempotency_conflict`, `payment_intents.idempotency_in_progress`, `payment_intents.callback_credit_pending` |
| Returns | `returns.error.refund_not_allowed`, `returns.error.refund_invalid_amount`, `returns.error.invalid_payload` |
| B2B | `b2b.po.invalid`, `b2b.po.not_found`, `b2b.po.invalid_status`, `b2b.po.approval_required`, `b2b.po.self_approval`, `b2b.po.order_service_not_configured` |
| Cüzdanlar | `wallets.error.bank_transfer_self_confirm`, `wallets.error.balance_condition_not_met`, `wallets.error.journal_already_applied` |
| Sefer atamaları | `fleet_tally.trip_assignments.error.self_approval` |
| Erişim noktaları | `access_points.error.inactive` |
| Faturalama | `invoices.already_issued`, `invoices.provider_empty_result`, `invoices.provider_missing`, `gib.invalid_invoice`, `gib.invalid_response`, `gib.submission_failed` |
| Stok araçları | `stock.vehicles.validation.plate_reason_required`, `stock.vehicles.validation.plate_reason_too_long`, `stock.vehicles.validation.plate_required`, `stock.vehicles.validation.duplicate_plate`, `stock.vehicles.validation.status_invalid`, `stock.vehicles.validation.invalid_payload`, `stock.vehicles.error.not_found` |
| Roller | `permissions.error.role_key_invalid`, `permissions.error.role_key_immutable`, `permissions.error.role_label_invalid` |

Bu anahtarların çoğu seed çevirilerinde 28 dilde bulunur (Bölüm 9.3). Stok araçlarının eski doğrulama anahtarları (`plate_required`, `duplicate_plate`, `status_invalid`, `invalid_payload`, `error.not_found`) hiçbir seed'de yoktur; panel bunları ham anahtar olarak gösterebilir (Bölüm 14.3, risk 26).

### 4.6 Değişen uçların karar sırası

Yönetici, yeni onay kontrollerinin hangi sırada çalıştığını bilmelidir. Bir isteğin hangi noktada reddedildiği, hata anahtarını ve yanıt kodunu belirler.

Refund (`POST /admin/payment-intents/:id/refunds`):

1. Oturum ve kiracı bağlamı (zincir).
2. Rota kapısı `payments:execute`.
3. Gövde bağlama (`approval_request_id` biçimi).
4. Servis: tutar > 0; anahtar zorunlu (boşsa `payment_intents.idempotency_key_required`, para hareketinden ve onaydan önce); intent yükleme; idempotensi.
5. Intent durumu `succeeded` olmalı.
6. Onay kapısı (`Authorize`): eşik, claim veya yeni talep.
7. Rezervasyon (`ReserveRefund`).
8. Sağlayıcı çağrısı (`adapter.Refund`).
9. Kayıt güncellemesi.

B2B onayı (`POST /admin/b2b/purchase-orders/:id/approve`):

1. Oturum ve `PermissionMiddleware` (modül `b2b`, takma ad `b2b_dealers`; `de3c4fa`).
2. Rota kapısı `b2b_dealers:approve`.
3. Gövde bağlama (`notes`).
4. Servis: durum `submitted`; mevcut adım; çağıranın rolü oturumdan; talep sahibi değil.

Banka havalesi onayı (`POST /admin/wallets/:id/bank-transfers/:journalId/confirm`):

1. Oturum ve kiracı bağlamı.
2. Rota kapısı `wallets:approve`.
3. Servis: journal bu cüzdana ait mi; journal zaten `posted` mi; beyan eden onaylayan mı (403).
4. `ConfirmPendingJournal`, `applyJournal` ve bildirim.

Sefer onayı (`POST /admin/fleet-tally/trip-assignments/:id/review`):

1. Rota kapısı `fleet_tally_trip_assignments:approve`.
2. Gövde bağlama (`action`).
3. Servis: manuel sefer mi; `in_progress` mi; bağımsız onaylayan mı (403).
4. Eylemin uygulanması.

### 4.7 Örnek istek ve yanıtlar: diğer değişen uçlar

Aşağıdaki örnekler ilgili uçların gövde alanlarını ve hata biçimini gösterir. Değerler uydurulmuştur.

Banka havalesi onayı (gövdesiz istek, `wallets:approve` gerekir):

```
POST /admin/wallets/66f1c2d3e4f5a6b7c8d9e0f1/bank-transfers/66f1c2d3e4f5a6b7c8d9e0f2/confirm
```

Hata örneği (beyan eden ile onaylayan aynıysa):

```json
{"success": false, "error": "wallets.error.bank_transfer_self_confirm"}
```

Satın alma siparişi gönderimi (zincir politikası yoksa):

```json
{"success": false, "error": "approvals.error.chain_not_configured"}
```

Satın alma siparişi kendi kendine onay:

```json
{"success": false, "error": "b2b.po.self_approval"}
```

Sefer onayı, kendi oluşturduğu seferi onaylamaya çalışınca:

```json
{"success": false, "error": "fleet_tally.trip_assignments.error.self_approval"}
```

Fatura, zaten kesilmiş:

```json
{"success": false, "error": "invoices.already_issued"}
```

Erişim noktası, pasif:

```json
{"success": false, "error": "access_points.error.inactive"}
```

Ödeme iadesi, anahtar gönderilmemiş:

```json
{"success": false, "error": "payment_intents.idempotency_key_required"}
```

Politika güncelleme, bilinmeyen kapı:

```json
{"success": false, "error": "approvals.error.gate_unknown"}
```

---

## 5. Veritabanı şeması

### 5.1 approval_policies

Politika belgeleri `tenant_id` ve `site_id` ile kapsamlanır. Aynı kiracı, site ve kapı için en fazla bir belge vardır.

| Alan | BSON | JSON | Tür | Not |
|---|---|---|---|---|
| `_id` | `_id,omitempty` | `id` | ObjectID | |
| `tenant_id` | `tenant_id` | `tenant_id` | metin | Kapsam |
| `site_id` | `site_id` | `site_id` | metin | Kapsam |
| `gate` | `gate` | `gate` | metin | Kapı anahtarı |
| `enabled` | `enabled` | `enabled` | bool | |
| `threshold` | `threshold,omitempty` | `threshold,omitempty` | *float64 | Tek eşikli kapılar |
| `currency` | `currency,omitempty` | `currency,omitempty` | metin | |
| `required_permission` | `required_permission,omitempty` | `required_permission,omitempty` | metin | |
| `steps` | `steps,omitempty` | `steps,omitempty` | dizi | Adım yapısı Bölüm 3.2 |
| `updated_by` | `updated_by,omitempty` | `updated_by,omitempty` | metin | |
| `created_at` | `created_at` | `created_at` | zaman | |
| `updated_at` | `updated_at` | `updated_at` | zaman | |

Dizin (`repository/indexes.go` 11-38):

| Ad | Koleksiyon | Anahtarlar | Benzersiz | Kısmi filtre |
|---|---|---|---|---|
| `approval_policies__tenant_site_gate` | `approval_policies` | `tenant_id` 1, `site_id` 1, `gate` 1 | Evet | Yok |

Dizinler `init` içinde kayıt edilir (`indexes.go` 40-42).

### 5.2 approval_requests

| Alan | BSON | JSON | Tür | Not |
|---|---|---|---|---|
| `_id` | `_id,omitempty` | `id` | ObjectID | |
| `tenant_id` | `tenant_id` | `tenant_id` | metin | |
| `site_id` | `site_id` | `site_id` | metin | |
| `gate` | `gate` | `gate` | metin | |
| `subject_ref` | `subject_ref` | `subject_ref` | metin | Ödeme niyeti kimliği (hex) |
| `amount` | `amount` | `amount` | float64 | Para birimi ayrı alandadır |
| `currency` | `currency` | `currency` | metin | |
| `requested_by` | `requested_by` | `requested_by` | metin | Talebi açan kullanıcı |
| `required_permission` | `required_permission` | `required_permission` | metin | Oluşturma anındaki kopya |
| `status` | `status` | `status` | metin | `pending`, `approved`, `rejected`, `consumed` |
| `decided_by` | `decided_by,omitempty` | `decided_by,omitempty` | metin | |
| `decided_at` | `decided_at,omitempty` | `decided_at,omitempty` | zaman | |
| `decision_note` | `decision_note,omitempty` | `decision_note,omitempty` | metin | |
| `created_at` | `created_at` | `created_at` | zaman | |
| `updated_at` | `updated_at` | `updated_at` | zaman | |

Dizinler (`repository/indexes.go` 11-38):

| Ad | Koleksiyon | Anahtarlar | Benzersiz | Kısmi filtre |
|---|---|---|---|---|
| `approval_requests__tenant_site_status` | `approval_requests` | `tenant_id` 1, `site_id` 1, `status` 1, `created_at` -1 | Hayır | Yok |
| `approval_requests__one_pending_per_subject` | `approval_requests` | `tenant_id` 1, `site_id` 1, `gate` 1, `subject_ref` 1, `amount` 1 | Evet | `{status: "pending"}` |

Kısmi benzersiz dizin, aynı kiracı, site, kapı, konu ve tutar için aynı anda yalnızca bir bekleyen talep bırakır. Servis ayrıca `FindPendingRequest` ile önceden denetler (`service.go` 301-305). Ekleme hatasında mevcut bekleyen talep yeniden okunur (`service.go` 318-327).

Talep sayısı sınırsızdır ve süre dolumu yoktur (Bölüm 14).

### 5.3 Değişen diğer koleksiyonlar

| Koleksiyon | Değişiklik | Kaynak |
|---|---|---|
| `roles` | `labels` alanı: dil kodu → metin eşlemesi (`labels,omitempty`) | `permissions/model/role.go` 17-20 |
| `menus` | `permissions` alanı, `nil` değilse yazılır; `nil` ise saklı izinler korunur | `menus/repository/mongo.go` 124-128 |
| `wallets` | `applied_journal_ids` (JSON'da `-`); son 50 kayıt tutulur | `wallets/model/wallet.go` 82; `repository/mongo.go` 150 |
| `trip_assignments` | `approval_status`, `driver_profile_id`, `transferred_to_trip_id` güncellemede yazılır | `trip_assignments/repository/mongo.go` 57-72 |
| `purchase_orders` | `approval_chain`, `approval_status`, `submitted_by` | `b2b/dealers/model/purchase_order.go` 63-69 |
| `stock_vehicle_plate_history` | `reason` alanı (`reason,omitempty`); plaka değişikliği gerekçesi burada saklanır | `stock/vehicles/model/model.go` (PlateHistory), `repository/mongo.go` |
| `ledger` dizinleri | Dizin tanımları bu işte değişti. Değişikliğin içeriği bu belgede ayrıntılandırılmadı | `ledger/repository/indexes.go` |
| `payment_intents` ve işlemler | Değişiklik yok; `approval_request_id` bu koleksiyonlarda saklanmaz | Bölüm 6.5 |

Satın alma siparişi onay adımı alanları (`purchase_order.go` 32-39): `approver_id`, `approver_role`, `approval_threshold`, `status` (`pending`, `approved`, `rejected`), `decision_at` (isteğe bağlı), `notes` (isteğe bağlı).

Satın alma siparişi durumları (`purchase_order.go` 12-20): `draft`, `submitted`, `approved`, `rejected`, `partial`, `fulfilled`, `cancelled`. `partial` ve `fulfilled` hiçbir yerde atanmaz.

Onay durumu (`ApprovalStatus`, `purchase_order.go` 25-29): `pending`, `approved`, `rejected`.

Sefer onay durumları (`trip_assignments/model/trip_assignment.go` 37-46): boş (yok), `pending`, `approved`, `rejected_refunded`, `rejected_transferred`. Kaynaklar: `auto_matched`, `manual`, `planned`.

Defter (ledger) durumları: `pending_reconciliation` (beyanda), `posted` (onayda). Beyan kaydının kaynağı `wallets_bank_transfer_declaration` değeridir.

### 5.4 Yazma ve okuma yolları

- Politika yazımı: PUT ucu, kiracı, site ve kapı anahtarıyla eşleşen politika belgesini yazar. `updated_by` ve `updated_at` alanları güncellenir.
- Talep oluşturma: `InsertRequest`. Yalnızca `pending` durumunda.
- Karar: koşullu `UpdateOne` (`pending` → `approved` veya `rejected`). Filtre talebin `pending` durumda olmasını da içerir.
- Tüketim: koşullu `UpdateOne` (`approved` → `consumed`). Filtre: `_id`, `gate`, `subject_ref`, `amount` ve `status` = `approved`. Eşleşen belge sayısı 1 olmalıdır (`repository/mongo.go` 185-203).
- Serbest bırakma: koşullu `UpdateOne` (`consumed` → `approved`). Filtre yalnızca `_id` ve `status` = `consumed` alanlarını içerir. Kapı, konu ve tutar yeniden denetlenmez (`repository/mongo.go` 205-220).

---

## 6. Durum makineleri

Aşağıdaki diyagramlar metin tabanlıdır. Her geçiş, koşul sağlanmadığında reddedilir.

### 6.1 Onay talebi

Durumlar: `pending`, `approved`, `rejected`, `consumed`.

```
  [kayıt yok]
      │  Authorize: eşik üstü tutar ve onay kimliği yok
      ▼
  pending ──── Approve: denetçi snapshot iznine sahip, karar veren talep sahibi değil ───► approved
     │                                                                                       │
     └──── Reject: aynı kontroller ───────────────────────────────────────────────► rejected  │
                   (terminal)                                                                │
                                                                                             │
  approved ── Authorize(approval_request_id): kapı, konu ve tutar eşleşir ──────────► consumed
                                                                                             │
  consumed ── Release: iade sağlayıcı kabulünden önce başarısız oldu ──────────────► approved
```

| Geçiş | İşlev | Konum | Koşullar |
|---|---|---|---|
| kayıt yok → pending | `Authorize` | `service.go` 301-329 | Kapı bilinir; politika etkin; eşik var ve tutar eşiği aşıyor ya da para birimi uyuşmuyor; onay kimliği yok. Aynı kapı, konu ve tutar için bekleyen talep varsa o kullanılır |
| pending → approved | `Approve` (`decide`) | `service.go` 343-395 | Karar veren boş değil; talep bulunur; durum `pending`; karar veren talep sahibi değil; snapshot izni var; koşullu güncelleme |
| pending → rejected | `Reject` (`decide`) | `service.go` 347-349 | Aynı kontroller |
| approved → consumed | `Claim` (`Authorize` içinde) | `service.go` 290-298; `repository/mongo.go` 185-203 | Onay kimliği verilmiş; `_id`, `gate`, `subject_ref`, `amount` ve durum eşleşir; tek belge güncellenir |
| consumed → approved | `Release` | `service.go` 334-337; `repository/mongo.go` 205-220 | Yalnızca bu talebi tüketen iade, sağlayıcı kabul etmeden başarısız olduğunda |

Süre dolumu yoktur. Bekleyen ve onaylanmış talepler karar veya tüketim gelene kadar kalır. `rejected` terminaldir. `consumed` kayıt, serbest bırakılmadıkça tüketilmiş kalır. Serbest bırakma için yönetici API'si yoktur.

### 6.2 Satın alma siparişi zinciri

Durumlar: `draft`, `submitted`, `approved`, `rejected`, `cancelled`. `approval_status` ayrı bir alandır: `pending`, `approved`, `rejected`.

```
  draft ──Submit──► submitted ──┬── son adım onaylandı ───────► approved ──Convert──► approved (order_id dolu)
    │                           │
    │                           └── mevcut adım reddedildi ──► rejected (terminal)
    │
    ├── Cancel ───────────────────────────────────────────────► cancelled
    └── Update (yalnızca draft)

  submitted ── Cancel ──────────────────────────────────────────► cancelled
  submitted ── Approve / Reject: yalnızca mevcut adım (sıra kuralı)
```

Ayrıntılı kurallar:

| Eylem | Önkoşul | Hata | Sonuç |
|---|---|---|---|
| Submit | Durum `draft` | `b2b.po.invalid_status` | Zincir çözülür; durum `submitted`; `approval_status` `pending`; `submitted_at` ve `submitted_by` yazılır |
| Submit, zincir boş | Hiçbir adım uygulanmadı | — | Durum `approved`; `approval_status` `approved`; `approved_at` yazılır |
| Submit, politika yok | Kapı yapılandırılmamış | `approvals.error.chain_not_configured` (400) | Sipariş `draft` kalır |
| Approve | Durum `submitted`; çağıran mevcut adımın rolünü taşır | `b2b.po.invalid_status`, `b2b.po.approval_required` | Mevcut adım `approved`; başka adım kaldıysa sipariş `submitted` kalır; kalmadıysa `approved` |
| Approve, talep sahibi | `submitted_by` çağıranla aynı | `b2b.po.self_approval` | Değişiklik yok |
| Reject | Durum `submitted`; çağıran mevcut adımın rolünü taşır | `b2b.po.invalid_status`, `b2b.po.approval_required` | Adım `rejected`; sipariş `rejected`; `rejected_at` yazılır |
| Cancel | Durum `draft` veya `submitted` | `b2b.po.invalid_status` | `cancelled` |
| Convert | Durum `approved`; `order_id` boş; sipariş servisi tanımlı | `b2b.po.invalid_status`, `b2b.po.order_service_not_configured` | `order_id` yazılır; durum değişmez. Mevcut `order_id` varsa reddedilir |

Kaynak: `po_service.go` 240-272 (Submit), 300-339 (Approve), 343-372 (Reject), 379-411 (Cancel, Convert), 474 (`order_id`).

Sıra kuralı: `currentStep` ilk bekleyen adımı döndürür. Yalnızca o adım karara bağlanabilir (`po_service.go` 276-283). Rol eşleşmesi büyük/küçük harf duyarsızdır (`hasRole`, 287-294). Taban zincirin "önce adım, sonra adım" kuralını ihlal eden bir karar yolu yoktur. Sıra kodla zorunludur.

Talep sahibi kuralı: `submitted_by` boş değilse ve çağıranın `user_id` değerine eşitse `b2b.po.self_approval` döner (`po_service.go` 308-310, 351-353). `submitted_by` boşsa denetim atlanır.

### 6.3 Banka havalesi beyanı

```
  Beyan (declare)                                    Onay (confirm)
  POST /me/wallet/bank-transfer-declarations         POST /admin/wallets/:id/bank-transfers/:journalId/confirm
          │                                                   │   (yetki: wallets:approve)
          ▼                                                   ▼
  defter kaydı: pending_reconciliation                 journal.SourceRef == cüzdan kimliği mi?
  bakiye DEĞİŞMEZ                                       ├── hayır ──► hata: "journal does not belong to this wallet"
  anahtar: bank_transfer:<cüzdan>:<referans>            ▼
  aynı anahtar aynı kaydı döndürür                     journal zaten posted mu?
                                                        ├── evet ──► cüzdan döner, hata yok (erken dönüş)
                                                        ▼
                                                      journal.CreatedBy == actor mü?
                                                        ├── evet ──► 403 wallets.error.bank_transfer_self_confirm
                                                        ▼
                                                      ConfirmPendingJournal ──► posted
                                                        ▼
                                                      applyJournal: ApplyBalanceChange (atomik $inc)
                                                        ▼
                                                      bildirim (yalnızca uygulandıysa)
```

| Durum | Anlamı |
|---|---|
| `pending_reconciliation` | Beyan edildi; bakiye henüz etkilenmedi |
| `posted` | Onaylandı ve bakiye uygulandı |

Cüzdan durumları (`wallets/model/wallet.go`): `active`, `frozen`, `closed`. Dondurulmuş veya kapalı cüzdanda işlem reddedilir ("wallet is frozen", "wallet is closed").

Kaynak: `wallets/service/service.go` 644-686. Bu akış kod okumasına dayanır. Çalışma zamanında doğrulanmamıştır.

Atomik bakiye değişikliği (`repository/mongo.go` 123-174): tek bir `FindOneAndUpdate`. Filtre: kiracı, site, `_id`, silinmemiş; isteğe bağlı durum `$nin`; isteğe bağlı `$expr` ile bakiye artı delta en az taban; defter değişikliğinde `applied_journal_ids` içinde olmayan kayıt. Güncelleme: `$inc cached_balance`, zaman damgaları, `last_applied_journal_id`, `$push applied_journal_ids` ile `$slice: -50`. Eşleşme yoksa kayıt yeniden okunur ve `BalanceChange.Check` ile sınıflandırılır. En fazla 3 deneme yapılır (`applyBalanceChangeAttempts`).

Pencere sabiti: `AppliedJournalWindow = 50` (`repository/interface.go` 31).

### 6.4 Sefer atama onayı

```
  Manuel sefer tamamlanır ──► approval_status = pending
                                      │
                 GET .../pending-approval (read) listeler
                                      │
                 POST .../:id/review  (approve yetkisi)
                                      │
     ┌────────────────────────────────┼────────────────────────────────────┐
     ▼                                ▼                                    ▼
  ApproveTrip                 RejectAndRefund                     RejectAndTransfer
  approval_status=approved    approval_status=rejected_refunded   yeni sefer oluşturulur
  kimlik ve zaman saklanmaz   iade en iyi çabayla denenir          (CreatedBy = onaylayan)
                              (iade hatası yutulur)                approval_status=rejected_transferred
                                                                   transferred_to_trip_id yazılır
```

Kurallar:

- `review` için izin `fleet_tally_trip_assignments:approve`. Rota, `fleet_tally` modül adıyla da ayrıca kontrol edilir (Bölüm 14).
- Bağımsızlık: oluşturan boş değilse ve çağıranla aynıysa 403 `fleet_tally.trip_assignments.error.self_approval` (`service.go` 332-337).
- `reject_refund` için `access_point_id` zorunludur.
- `ApproveTrip` yalnızca `approval_status` değerini değiştirir. Onaylayanın kimliği ve onay zamanı saklanmaz (bilinen eksik, Bölüm 14).
- Manuel seferler bitince `pending` olur (`service.go` 292-297).

Kaynak: `trip_assignments/service/service.go` 332-429.

### 6.5 İade rezervasyonu ve onay tüketimi

Bu diyagram `payments/intents/service/service.go` içindeki `RefundWithApproval` akışını gösterir (satır 528-636).

```
 RefundWithApproval(id, amount, reason, key, actor, approval_request_id)
   │
   ├─ amount > 0 değilse ──────────────────────────────────────────────► hata
   ├─ key boşsa ──────────────────────────────────────────────────────► payment_intents.idempotency_key_required (400; onaydan ve para hareketinden önce)
   ├─ intent yüklenir; claimOperation (idempotensi)
   │     aynı anahtar, farklı tutar ─────────────────────────────────► payment_intents.idempotency_conflict (422)
   │     aynı anahtar, başarılı ──────────────────────────────────────► replay (200)
   │     aynı anahtar, işlemde ───────────────────────────────────────► payment_intents.idempotency_in_progress (409)
   ├─ intent durumu succeeded değilse ───────────────────────────────► payment_intents.refund_not_allowed
   ├─ adapter çözülemezse ────────────────────────────────────────────► ErrInvalidPayload
   │
   ├─ Authorize(payments.refund, intent.ID, amount, currency, actor, approval_request_id)
   │     politika kapalı veya eşik altı ─────────────────────────────► (false, nil): onay gerekmez
   │     onay kimliği yok ───────────────────────────────────────────► bekleyen talep oluşur → RequiredError (409)
   │     onay kimliği var ────────────────────────────────────────────► Claim: approved → consumed
   │                                                                    eşleşmezse ► approvals.error.not_granted (409)
   │     claimedApproval = kimlik (yalnızca tüketildiyse)
   │
   ├─ ReserveRefund: refunded_amount + tutar ≤ tahsil edilen tutar ? (atomik)
   │     hayır ──► withApprovalReleased ──► ErrRefundExceedsCapture ──► operation failed
   │
   ├─ adapter.Refund(...)
   │     hata ──► ReleaseRefund (rezervasyon geri alınır)
   │              ──► withApprovalReleased (talep approved'a döner)
   │              ──► operation failed ──► hata döner
   │     başarı ─► operation succeeded (anahtar her durumda kapanır)
   │              ──► intent.ProviderPayload["refund"] yazılır; repo.Update
   │              ──► createRefundPayment; sipariş tipinde dijital erişim iptali
   │
   └─ Son durum: talep consumed (başarılı iade) veya approved (sağlayıcı öncesi başarısız)
```

Kaynak: `payments/intents/service/service.go` 528-636 ve 641-650 (`withApprovalReleased`).

Durum değerleri (kesin metinler):

| Kavram | Değerler | Kaynak |
|---|---|---|
| Onay talebi | `pending`, `approved`, `rejected`, `consumed` | `approvals/model/model.go` 36-42 |
| Ödeme niyeti | `requires_payment_method`, `requires_confirmation`, `requires_action`, `requires_capture`, `succeeded`, `cancelled`, `failed` | `payments/intents/model/intent.go` 12-18 |
| İşlem türü | `capture`, `refund`, `callback`, `cancel` | `intents/model/operation.go` 12-15 |
| İşlem durumu | `succeeded`, `failed`, `pending` | `intents/model/operation.go` 21-26 |
| İade durumu (returns) | `pending`, `refunded`, `failed` | `returns/model/return.go` 48-50 |

"awaiting-approval" veya "awaiting_approval" adı kodda yoktur. Onay bekleyen iade, bekleyen bir onay talebi ve 409 yanıtıyla temsil edilir.

Bir iade ödeme niyetinin durumunu değiştirmez. Kısmi ve tam iadelerde durum `succeeded` kalır.

Bilinen boşluklar (Bölüm 14'te ayrıntılı):

- Talep tüketildikten sonra sağlayıcı çağrısından önce süreç çökerse talep `consumed` kalır. Serbest bırakma API'si yoktur.
- `adapter.Refund` herhangi bir hata döndürürse onay serbest bırakılır. Adaptörün yalnızca "para hareket etmedi" anlamında hata döndürüp döndürmediği bu belgede doğrulanmamıştır. Zaman aşımı durumu ve PENDING_EXTERNAL sorgusu karar bekliyor (Bölüm 14.2, madde 4).
- Sağlayıcı kabul ettikten sonra `repo.Update` başarısız olursa işlem hata döner. Para hareket etmiş olsa da onay `consumed` kalır.
- `approval_request_id` iade kaydında (ödeme niyeti, işlem veya iade ödemesi) saklanmaz. Bağlantı yalnızca onay talebinin `subject_ref`, `amount` ve `status` alanlarından izlenebilir.

---

## 7. Entegrasyon noktaları

### 7.1 Ödeme niyetleri (`payments/intents`)

- Arayüz `RefundApprover` (`service.go` 108-116): `Authorize(ctx, gate, subjectRef, amount, currency, actorID, approvalRequestID) (bool, error)` ve `Release(ctx, id) error`.
- `approvals/service` bu arayüzü uygular. Bağlantı `cmd/api/main.go` 2042-2043: `approvalsSvc := approvalsservice.New(...)` ve `intentSvc.WithRefundApprover(approvalsSvc)`.
- `RefundWithApproval` (`service.go` 528): iade için kapı çağrısı `gate = payments.refund`, `subjectRef = intent.ID.Hex()`, `currency = intent.Currency`.
- Yeni istek alanı: `approval_request_id`. Yalnızca `Authorize` tarafından kullanılır.
- Eski `Refund` (`service.go` 518) `RefundWithApproval`'ı `nil` onay kimliğiyle çağırır.

### 7.2 Returns (`returns`)

- `returns/service/service.go` 293-377 iadeyi `RefundWithApproval` üzerinden yapar.
- İdempotensi anahtarı `return:<return_id>` biçimindedir (`service.go` 350).
- `awaitsApproval` (`service.go` 543-545): `RequiredError` veya `ErrApprovalNotGranted` durumunda iadenin durumu yazılmaz ve hata döner.

### 7.3 Servis rezervasyonları (`service_reservations`)

- `POST /admin/service-reservations/:id/cancel` (`routes.go` 24; `cmd/api/main.go` 2516) iadeyi `approval_request_id` olmadan çağırır (`service_reservations/service/service.go` 197). Hata her zaman 400 döner (`handler/admin.go` 184-189).
- Kod okumasına göre bu yol onaylanmış bir talebi tüketemez. Çalışma zamanında doğrulanmamıştır.

### 7.4 B2B satın alma siparişleri (`b2b/dealers`)

- Arayüz `service.ApprovalChainPolicy`. Bağlantı `cmd/api/main.go` 2519: `b2b_dealers.RegisterRoutes(r, database, protectedChain, orderSvc, approvalsSvc)`.
- Çağrı (`po_service.go` 576): `ApprovalChain("b2b.purchase_order", GrandTotal, dealer.CreditLimit.Total)`.
- Dönen zincir sipariş modelinde `approval_chain` olarak kopyalanır (`po_service.go` 254).

### 7.5 Cüzdanlar (`wallets`)

- Onay kapısı politika dışıdır. İzin ve kendi kendine onay kontrolü kullanır (sahip kararı).
- Beyan ve onay servisleri defter servisini kullanır (`PostPendingJournal`, `ConfirmPendingJournal`).

### 7.6 Sefer atamaları (`fleet_tally/trip_assignments`)

- İzin kapısı politika dışıdır (sahip kararı).
- Bağımsızlık kontrolü servis içindedir (`requireIndependentApprover`, `service.go` 332-337).

### 7.7 Erişim noktaları (`access_points`)

- `requireActiveAccessPoint` (`service/service.go` 378-383) pasif erişim noktalarını reddeder.
- Çağıran işlemler: `RecordEvent` (`service.go` 943-957), `SettleTripCash` (1590-1597), `RecordWalkOnCashEntry` (1723-1730).
- `RefundEvent` (`service.go` 324) bu kontrolü çağırmaz. Kod okumasına göre.

### 7.8 Faturalama (`billing/invoices`)

- Sağlayıcı kaydı (`cmd/api/main.go` 1447-1470): mock e-invoice ve mock ERP yalnızca `cfg.Env` production değilken kaydedilir (`33467c8`). Mock PDF renderer (1477), mock kargo (1353) ve mock SMS (`internal/modules/sms/service/service.go` 47) her ortamda kayıtlıdır. GİB adaptörü yalnızca GİB ortam değişkenleri tanımlıysa eklenir. Logo adaptörü de yalnızca kendi ortam değişkenleri tanımlıysa eklenir.
- Hata sözleşmesi Bölüm 4.3.6'da.

### 7.9 Menüler ve izinler

- Menü görünürlüğü `handler/menu_gate.go` içindedir. Her menü öğesi `roles` ve `permissions` ile kapılanır (Bölüm 10.3).
- Yetki kontrolü `middleware.HasPermission` ile yapılır: tam anahtar, `modül:*` veya `*:*`.

### 7.10 Seed'ler

- `internal/seeds/runner.go` sırası: `SeedLocalization`, `SeedPermissions`, `SeedRoles`, `RunApprovalsSeed`, `SeedTranslations`, ardından `SeedMenus` ve diğerleri (Bölüm 9.1).

### 7.11 Panel

- Panel onay uçlarını `lib/api/approvals.ts` içinden çağırır (Bölüm 11.7).

---

## 8. Önemli fonksiyonlar

Aşağıdaki parçalar okunan koddan alınmıştır. Uzun fonksiyonlarda yalnızca ilgili satırlar gösterilmiş ve kısaltma `// ...` ile belirtilmiştir.

### 8.1 Authorize

`internal/modules/approvals/service/service.go` 272-308:

```go
func (s *service) Authorize(ctx context.Context, gate, subjectRef string, amount float64, currency, actorID string, approvalRequestID *primitive.ObjectID) (bool, error) {
	g, ok := model.FindGate(gate)
	if !ok {
		return false, model.ErrGateUnknown
	}
	policy, _, err := s.effectivePolicy(ctx, g)
	if err != nil {
		return false, err
	}
	if !policy.Enabled {
		return false, nil
	}
	// Below the threshold the operation proceeds without a checker. A currency
	// that does not match the policy's cannot be compared, so it needs approval.
	if policy.Threshold != nil && strings.EqualFold(policy.Currency, currency) && amount <= *policy.Threshold {
		return false, nil
	}

	if approvalRequestID != nil {
		claimed, err := s.repo.Claim(ctx, *approvalRequestID, gate, subjectRef, amount)
		if err != nil {
			return false, err
		}
		if !claimed {
			return false, model.ErrApprovalNotGranted
		}
		return true, nil
	}
	// ... bekleyen talep bulunur veya oluşturulur; RequiredError döner.
}
```

Eşitlik durumunda (`amount == threshold`) iade onaysız geçer. `Threshold` nil ise ve kapı etkinse, her iade onay ister.

### 8.2 Kapı çağrısı ve serbest bırakma (iade)

`internal/modules/payments/intents/service/service.go` 562-587 (kısaltılmış):

```go
var claimedApproval *primitive.ObjectID
if s.refundApprover != nil {
	claimed, err := s.refundApprover.Authorize(ctx, approvalsmodel.GateRefund, intent.ID.Hex(), amount, intent.Currency, updatedBy, approvalRequestID)
	if err != nil {
		s.finishOperation(ctx, claim, model.OperationFailed, map[string]any{"error": err.Error()})
		return nil, err
	}
	if claimed {
		claimedApproval = approvalRequestID
	}
}
```

Sağlayıcı çağrısı ve hata yolu (kısaltılmış, `service.go` 589-611):

```go
resp, err := adapter.Refund(ctx, AdapterRefundInput{ /* ... */ })
settleCtx := contextkeys.Detached(ctx)
if err != nil {
	if releaseErr := s.repo.ReleaseRefund(settleCtx, id, amount); releaseErr != nil {
		err = errors.Join(err, releaseErr)
	}
	err = s.withApprovalReleased(settleCtx, claimedApproval, err)
	s.finishOperation(settleCtx, claim, model.OperationFailed, map[string]any{"error": err.Error()})
	return nil, err
}
```

`withApprovalReleased` (`service.go` 641-650):

```go
func (s *service) withApprovalReleased(ctx context.Context, claimed *primitive.ObjectID, err error) error {
	if claimed == nil || s.refundApprover == nil {
		return err
	}
	if releaseErr := s.refundApprover.Release(contextkeys.Detached(ctx), *claimed); releaseErr != nil {
		return errors.Join(err, releaseErr)
	}
	return err
}
```

Bu kod, `adapter.Refund` bir hata döndürdüğünde paranın hareket etmediğini varsayar. Bu varsayım adaptör sözleşmesine bağlıdır ve bu belgede doğrulanmamıştır. Bu konu karar bekliyor (Bölüm 14.2, madde 4).

### 8.3 Returns: onay bekleyen hata ayrımı

`internal/modules/returns/service/service.go` 543-546:

```go
// awaitsApproval reports whether a refund error means a checker decision is
// still needed or the supplied approval was not usable. Neither is a failed refund.
func awaitsApproval(err error) bool {
	var required *approvalsmodel.RequiredError
	return errors.As(err, &required) || errors.Is(err, approvalsmodel.ErrApprovalNotGranted)
}
```

### 8.4 Kapı kataloğu ve izin biçimi

Kapı tanımları ve `DefaultPolicy` için Bölüm 3.1'e bakın. İzin biçimi kontrolü (`gates.go` 57-61):

```go
func PermissionKeyValid(key string) bool {
	module, action, ok := strings.Cut(key, ":")
	return ok && module != "" && action != ""
}
```

### 8.5 Yetki eşleşmesi

`internal/platform/middleware/permission.go` 111-122:

```go
func HasModulePermission(perms []string, module string) bool {
	if len(perms) == 0 {
		return false
	}
	prefix := module + ":"
	for _, p := range perms {
		if p == "*:*" || p == module+":*" || strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}
```

Bu fonksiyon `PermissionMiddleware` içinde modül izni için çağrılır. `module:*` ve `*:*` joker izinleri de kabul edilir.

### 8.6 Rota kapıları

`internal/modules/approvals/routes.go` 18-23:

```go
g.GET("/approval-gates", middleware.ModulePermissionMiddleware("approval_policies", "read"), h.ListGates)
g.GET("/approval-policies/:gate", middleware.ModulePermissionMiddleware("approval_policies", "read"), h.GetPolicy)
g.PUT("/approval-policies/:gate", middleware.ModulePermissionMiddleware("approval_policies", "update"), h.PutPolicy)
g.GET("/approval-requests", middleware.ModulePermissionMiddleware("approval_requests", "read"), h.ListRequests)
g.POST("/approval-requests/:id/approve", middleware.ModulePermissionMiddleware("approval_requests", "read"), h.Approve)
g.POST("/approval-requests/:id/reject", middleware.ModulePermissionMiddleware("approval_requests", "read"), h.Reject)
```

Banka havalesi onayı (`wallets/routes.go` 26):

```go
g.POST("/wallets/:id/bank-transfers/:journalId/confirm", middleware.ModulePermissionMiddleware("wallets", "approve"), admin.ConfirmBankTransfer)
```

Sefer onayı (`fleet_tally/trip_assignments/routes.go` 19):

```go
g.POST("/:id/review", middleware.ModulePermissionMiddleware("fleet_tally_trip_assignments", "approve"), admin.ReviewTrip)
```

### 8.7 Satın alma siparişi kapısı

`internal/modules/b2b/dealers/routes.go` (kısaltılmış):

```go
poApproveGuard := middleware.ModulePermissionMiddleware("b2b_dealers", "approve")
attachProtected("POST", "/admin/b2b/purchase-orders/:id/approve", poHandler.ApprovePurchaseOrder, poApproveGuard)
attachProtected("POST", "/admin/b2b/purchase-orders/:id/reject", poHandler.RejectPurchaseOrder, poApproveGuard)
```

`attachProtected`, verilen `authMiddleware` dilimini (B2B için `protectedChain`) ve ek korumaları birleştirir. Bu nedenle `PermissionMiddleware` her iki uçta da çalışır.

### 8.8 Boş zincirde otomatik onay

`internal/modules/b2b/dealers/service/po_service.go` 262-266:

```go
if len(po.ApprovalChain) == 0 {
	po.Status = model.POStatusApproved
	po.ApprovalStatus = model.ApprovalStatusApproved
	po.ApprovedAt = &now
}
```

Bu davranış belgelenmiştir. Bölüm 14'te risk olarak yer alır.

### 8.9 Tohum sırası

`internal/seeds/runner.go` 56-66 (kısaltılmış):

```go
permIDsByTenant, err := SeedPermissions(ctx, db)
if err != nil {
	return err
}
if err := SeedRoles(ctx, db, permIDsByTenant); err != nil {
	return err
}
if err := RunApprovalsSeed(ctx, db); err != nil {
	return err
}
if err := SeedTranslations(ctx, db); err != nil {
	return err
}
```

### 8.10 Rol anahtarı ve dil kodu denetimi

- Rol anahtarı deseni `^[a-z][a-z0-9_]{1,62}$` (`permissions/service/role_keys.go` 23-28). Oluşturmada `CreateRole` içinde denetlenir (`service.go` 100).
- Dil kodu deseni `^[a-z]{2,3}(-[a-z0-9]{2,8})?$` (`permissions/model/role.go` 32, 43-45). Etiket denetimi `CreateRole` (`service.go` 103) ve `UpdateRole` (`service.go` 141) içinde yapılır.
- Değişmezlik denetimi `UpdateRole` içindedir (`service.go` 138-140). Ad değişikliği bu hatayla reddedilir.

### 8.11 Gürsel GPS yapılandırması (c27c3b5)

`internal/seeds/gursel_turizm_seed.go` içindeki `seedGurselTurizimGpsProviders` artık sağlayıcı alanlarını yalnızca ekleme sırasında yazar (kısaltılmış):

```go
// Written once, when the configuration is created. Afterwards the panel owns it,
// so a later seed run never overwrites an admin's credentials, endpoint,
// activation or name.
update := bson.M{
	"$setOnInsert": bson.M{
		"_id":             primitive.NewObjectID(),
		"tenant_id":       s.TenantID,
		"site_id":         s.SiteID,
		"provider_type":   s.ProviderType,
		// ... credentials, timezone_offset, is_active, created_at ve updated_at aynı bloktadır
	},
},
```

Önceki sürümde bu alanlar `$set` içindeydi ve her tam seed zincirinde kod değerlerine dönüyordu.

---

### 8.12 Plaka değişikliği gerekçesi

Aşağıdaki parçalar `origin/aypass-approval-ui` üzerindeki `internal/modules/stock/vehicles/service/service.go` dosyasından alınmıştır (kısaltılmış).

```go
const maxPlateReasonLength = 500

// normalizePlateReason gerekçeyi kırpar; boş değer ve 500 karakteri aşan değer reddedilir.
func normalizePlateReason(v string) (string, error) {
	reason := strings.TrimSpace(v)
	if reason == "" {
		return "", ErrPlateReasonRequired
	}
	if utf8.RuneCountInString(reason) > maxPlateReasonLength {
		return "", ErrPlateReasonTooLong
	}
	return reason, nil
}
```

Güncelleme yolunda gerekçe yalnızca plaka değiştiğinde denetlenir ve geçmiş kaydına yazılır:

```go
plateChanged := current.CurrentPlate != item.CurrentPlate
reason := ""
if plateChanged {
	if reason, err = normalizePlateReason(plateReason); err != nil {
		return err
	}
}
// ... araç kaydı güncellenir ...
if plateChanged {
	now := time.Now().UTC()
	_ = s.repo.CloseActivePlateHistory(ctx, id, now)  // hata yok sayılıyor (Bölüm 14.3, risk 25)
	_ = s.repo.InsertPlateHistory(ctx, model.PlateHistory{VehicleID: id, Plate: item.CurrentPlate, StartedAt: now, CreatedBy: item.UpdatedBy, Reason: reason})
}
```

---

## 9. Seed davranışı

### 9.1 Çalışma sırası

`internal/seeds/runner.go`:

| Sıra | Fonksiyon | Satır | Not |
|---|---|---|---|
| 1 | `SeedLocalization` | 53 | Diller ve yerelleştirme önce gelir |
| 2 | `SeedPermissions` | 56 | Kiracı başına izin kimlikleri |
| 3 | `SeedRoles` | 60 | Rollerin politika öncesi olması gerekir (`21a8d61`) |
| 4 | `RunApprovalsSeed` | 63 | Yalnızca eksik politikaları ekler |
| 5 | `SeedTranslations` | 66 | Çeviri katmanları |
| sonra | `SeedMenus` (78), `SeedPublicNavMenus` (81) ve diğerleri | 78 ve sonrası | |

`cmd/seed/main.go` 168 satırında `RunApprovalsSeed` ayrıca çağrılır.

### 9.2 Onay politikası tohumu

`RunApprovalsSeed` (`internal/seeds/approvals_seed.go` 47-74):

- Yalnızca varsayılan kiracı ve onun çözülen varsayılan sitesi için çalışır.
- Eklemeler `$setOnInsert` ile yapılır. Upsert `tenant_id`, `site_id` ve `gate` üzerindendir. Mevcut belge asla üzerine yazılmaz. Yönetici tarafından düzenlenmiş bir politika korunur.
- `updated_by` değeri `seed` olur.
- Varsayılanlar: `b2b.purchase_order` etkin üç adımlı zincir (Bölüm 3.5). `payments.refund` kapalı, izin `payments:approve`, eşik yok (Bölüm 3.6).
- Diğer kiracılar ve siteler politikalarını yönetici ekranlarından oluşturur.

### 9.3 Çeviri tohumları

| Dosya | Anahtar sayısı | İçerik | Dil sayısı |
|---|---|---|---|
| `approvals_translations_seed.go` | 17 | 15 `approvals.error.*`, 2 `approvals.gate.*` | 28 |
| `approvals_ui_translations_seed.go` | 76 | `menu.approvals.policies`, `menu.approvals.requests`, 74 `approvals.ui.*` | 28 |
| `module_error_translations_seed.go` | 21 | 3 `approvals.activity.*`, 2 `stock.vehicles.validation.plate_reason_*` ve diğer modül hata anahtarları | 28 |
| `translation_seed.go` (`stock.vehicles.fields.plate_reason`) | 1 | Plaka gerekçesi alanı etiketi | 28 |
| `role_translations_seed.go` | 3 | `permissions.error.role_key_invalid`, `role_key_immutable`, `role_label_invalid` | 28 |

Kullanılan 28 dil kodu: `tr`, `en`, `de`, `fr`, `it`, `es`, `pt`, `pt-br`, `az`, `zh`, `ja`, `fa`, `ar`, `kk`, `uz`, `ru`, `uk`, `pl`, `hr`, `cs`, `bs`, `sr`, `sq`, `hy`, `az-cyrl`, `ky`, `tk`, `mn`.

Bu çeviriler `translation_seed.go` içinde çağrılan katmanlarla uygulanır (yaklaşık 12532-12535 satırları). Boş değer yoktur.

Ana dil incelemesi: `ky`, `kk`, `mn`, `hy`, `az-cyrl`, `uz` ve `tk` dilleri ana dil incelemesi bekler. Sahip kararıyla ertelenmiştir (Bölüm 14.2). Latince "idempotency" kelimesi `fa`, `ar` ve `hy` dillerinde bilinçli olarak bırakılmıştır.

Panelin kullandığı 98 sabit anahtarın tümü `internal/seeds` içinde bulunur (Bölüm 11.6).

### 9.4 Gürsel GPS sağlayıcı yapılandırmaları

Karar: yapılandırma yalnızca oluşturulurken yazılır. Commit `c27c3b5` ("Write Gürsel GPS provider configurations only when they are created"), backend dalı `aypass-approval-ui`, `f286fd0` üzerine eklendi.

- Değişen dosyalar: `internal/seeds/gursel_turizm_seed.go` (11 satır) ve `internal/seeds/gursel_turizm_gps_seed_test.go` (yeni, 46 satır).
- Tam seed zinciri, güncelleme ve konteyner başlatma artık `credentials`, `base_url`, `is_active`, `display_name` ve `timezone_offset` alanlarını mevcut kayıtlarda değiştirmez. Panelden yapılan değişiklik korunur.
- Regresyon testi: `TestSeedGurselTurizimGpsProviders_KeepsAdminChangesOnRerun`. Test, yönetici tarafından değiştirilen bir alanın ikinci seed çalıştırmasında korunduğunu doğrular. Bu testin çalıştırma sonucu sahibin doğrulama listesinde yer almaz. Bu belgede sonucu iddia edilmez.
- Taban commit'te bu alanlar her tam seed çalıştırmasında kod değerlerine dönüyordu.

### 9.5 Rol tohumları

`internal/seeds/role_seed.go`:

- `admin`: kiracının izinlerinin tümünü alır. Platforma özel `system` ve `system_ops` modülleri hariçtir (`role_seed.go` 21-24, 85-100). Rol `is_system: true` ile oluşturulur (102-115). Rol mevcutsa izinleri her çalıştırmada zorla eşitlenir (117-130). Yani `admin`, platform dışı tüm izinlere her zaman eşitlenir.
- `editor`, `reviewer` ve `fleet_viewer`: yalnızca eksikse oluşturulur (`role_seed.go` 383, 416, 442). Var olan rollerin izinlerine rutin çalıştırmada dokunulmaz. `is_system` değeri `false`'tur.
- `reviewer` varsayılanı: sınıflandırılmış ilanlar için `read`, `review`, `publish` veya `update` izinleri (400-414). `fleet_viewer` varsayılanı: `fleetViewerModules` listesi ve `read` (436-501).
- `syncPublicSiteManagerRoles` (462) her iki modda çalışır. Bu tasarım gereğidir.
- `GrantRoleDefaults` (`role_seed.go` 56-61) yalnızca `aystek-seed --roles` ile çalışır (`cmd/seed/main.go` 38, 158-173). Var olan `editor`, `reviewer` ve `fleet_viewer` rollerine eksik varsayılan izinleri ekler (`grantRoleDefaults`, 523-542). Hiçbir izin kaldırılmaz. Etiketler gönderilmediği için değişmez.
- `store_manager`, `credit_dept` ve `executive` rollerini hiçbir seed oluşturmaz. Bu roller yönetici tarafından oluşturulmalıdır. Aksi halde B2B zincirindeki bu adımlar onaylanamaz.

### 9.6 Kurulum ve kapsayıcı

- `deployments/production/docker-compose.yml` satır 63-64: `SEED_ON_START: "false"`. Yanındaki yorum, seed'lerin installer'ın tek seferlik seed kapsayıcısında çalıştığını belirtir. Bu dosyada seed ile ilgili başka ayar yoktur.
- Commit `ccb13e2` değeri `"true"` iken `"false"` olarak değiştirdi.
- `cmd/api/main.go` içindeki `maybeRunSeeds` (3051): seed'ler yalnızca `cfg.Env` `dev` ise veya `SEED_ON_START` doğruysa çalışır (3052). `SEED_TIMEOUT_SECONDS` varsayılanı 600 saniyedir. Seed hatası API başlatmayı durdurmaz; uyarı olarak kaydedilir (3082).
- Installer tek seferlik kapsayıcıda tam zinciri çalıştırır. `SEED_ON_START=true`, `SEED_EXIT_AFTER_RUN=true` ve `SEED_FULL_RESET=false` kullanılır.

### 9.7 Test seed kimlik bilgileri

Test için seed kimlik bilgileri kodda kalır. Üretimde değiştirilecektir (sahip kararı).

---

## 10. Yetkiler ve menüler

### 10.1 Onay izinleri

`internal/seeds/permission_seed.go` içinde `approval_policies` (satır 104) ve `approval_requests` (satır 105) modülleri listelenir. Varsayılan eylem listesi (satır 215-231) `approve` içerir. Özel listesi olmayan modüller bu varsayılan listeyi alır. Bu nedenle `approval_policies:approve` ve `approval_requests:approve` belgeleri de oluşur (kod okumasına göre). Uç yetkileri bu eylemleri kullanmaz.

Onay uçlarında ve yeni kapılarda kullanılan izinler:

| İzin | Kullanım |
|---|---|
| `approval_policies:read` | Kapı listesi, politika okuma, politika menüsü |
| `approval_policies:update` | Politika yazma |
| `approval_requests:read` | Talep listesi ve karar uçları (rota kapısı) |
| `payments:execute` | İade uçları |
| `payments:approve` | Refund kapısının varsayılan denetçi izni |
| `b2b_dealers:approve` | B2B satın alma siparişi onayı ve reddi |
| `wallets:approve` | Banka havalesi onayı (tohumlanır, `18e27f8`; Bölüm 10.2) |
| `fleet_tally_trip_assignments:approve` | Sefer onayı |
| `fleet_tally_trip_assignments:read` | Bekleyen sefer listesi |

### 10.2 wallets:approve tohumu (düzeltildi)

`permission_seed.go` 270 satırında `wallets` modülünün eylem listesi `create`, `read`, `update`, `delete`, `execute` ve `approve` olarak yazılmıştır (`18e27f8`). Taban commit'te `approve` bu listede yoktu; bu nedenle `wallets:approve` belgesi hiçbir role verilemiyordu.

Etkisi: Platform kullanıcıları (`is_platform_user`, `impersonated` olmadığı sürece) bu kapıyı geçer. Platform dışı yöneticiler, `wallets:approve` izni verilmiş bir role sahipse banka havalesini onaylayabilir. Kurulumdaki canlı rol atamaları bu belgede doğrulanmamıştır.

### 10.3 Menü kapısı

`internal/modules/menus/handler/menu_gate.go`:

- `filterMenuTreeByAccess` (64-76): Platform kullanıcısı (impersonate değilse) tüm ağacı görür. Rolü ve izni olmayan öğe herkese görünür. Diğer durumda rolü (`hasAnyRole`) veya izni (`hasAnyPermission`) eşleşen öğe görünür.
- `hasAnyPermission` `middleware.HasPermission` kullanır (tam anahtar, `modül:*` veya `*:*`). Yalnızca `approval_policies:update` iznine sahip bir kullanıcı, `approval_policies:read` gerektiren politika menüsünü görmez.
- Grup koruma kuralı (`filterMenuTree`, 84-101): Alt öğeler önce süzülür. Hiç alt öğesi kalmayan grup (boş URL) düşer. `allow` tarafından reddedilen düğüm düşer. Yalnızca `filterMenuTreeByAccess` `keepRejectedGroups` ayarını açar. Modül kapısı ve görünüm profili filtreleri bu ayarı kullanmaz.
- Taban commit'te filtre yalnızca rolleri kullanıyordu. İzinler yoktu.

Editör rolünün görünürlük değişikliği (sahip tarafından onaylandı): Editör artık görmeye izinli olduğu alt öğeleri olan üç grubu görür. Bu gruplar ticari siparişler, ticari ödemeler ve operasyondur. Bu değişikliğin çalışan bir kiracıdaki etkisi bu belgede doğrulanmamıştır.

### 10.4 Menü tohumu

`internal/seeds/menu_seed.go`:

| Anahtar | Yol | Başlık anahtarı | İzin | Rol | Sıra | İkon |
|---|---|---|---|---|---|---|
| `approval-policies` | `/admin/approval-policies` | `menu.approvals.policies` | `approval_policies:read` | yok | 5 | scale |
| `approval-requests` | `/admin/approval-requests` | `menu.approvals.requests` | `approval_requests:read` | yok | 6 | list-checks |

Her iki öğe `security` grubu altında yer alır (satır 362; üst grubun rolü `admin`dir). Görünür olmaları için ilgili izin gerekir.

Menü yazımı (`menus/repository/mongo.go` 124-128): `menu.Permissions` `nil` değilse izinler yazılır. `nil` ise saklı izinler korunur. Yönetici API'si `permissions` alanını kabul etmez (`MenuRequest`, `dto.go` 5-16). İzinler yalnızca seed'ler aracılığıyla yazılır.

### 10.5 Yetki ara katmanı

`ModulePermissionMiddleware(modül, eylem)` (`internal/platform/middleware/module_permission.go` 13-51):

- Kullanıcı kimliği yoksa 401.
- Platform kullanıcısı (impersonate değilse) geçer (27-30).
- `permissions` bağlamı yoksa veya türü `[]string` değilse 403.
- `HasPermission(list, "modül:eylem")` geçerse devam eder. Geçmezse 403 `missing permission modül:eylem`.

`PermissionMiddleware()` (`permission.go` 11-76): `/public`, `/pa`, `/me` ve `/notifications` yollarını muaf tutar. Modülü yol segmentinden çıkarır (`resolveModuleFromPath`). `/admin/translation/*` yolu `posts` modülüne eşlenir. Modül izni yoksa platform kullanıcısı bypass alır. Diğerleri 403 `missing module permission` alır.

Test: `module_permission_approve_test.go` (`TestModulePermissionMiddleware_ApproveGateIsExactAction`). `fleet_tally_trip_assignments:approve` için yalnızca `execute` iznine sahip kullanıcı 403 alır. Açık `approve` iznine sahip kullanıcı 200 alır. Modül joker izni de 200 alır.

---

## 11. Panel ekranları

Panel deposu: `/workspace/node/aystek-core-next`, dal `approval-policy-ui`. Bu iş panelde "admin-panel-next" olarak da anılmıştır. Depo dizini `aystek-core-next`tir.

### 11.1 Ekran listesi

| Yol | Erişim kontrolü | Bileşenler |
|---|---|---|
| `/admin/approval-policies` | `approval_policies:read` ile görünür. Düzenleme yalnızca `approval_policies:update` ile etkin | `ApprovalsHero`, `AdminDataTable`, `PolicyEditorDialog` (`ChainStepsEditor`, `PermissionSelect`), `useApprovalGates` |
| `/admin/approval-requests` | `approval_requests:read` ile görünür | `ApprovalsHero`, `AdminDataTable`, `ReasonModal` (not isteğe bağlı), `useApprovalRequests`, `useApprovalGates` |
| `/admin/orders/[id]` (değişti) | Mevcut sayfa | `usePaymentIntentRefund`, `RefundApprovalNotice` |
| `/admin/roles/[id]` (değişti) | Mevcut sayfa | `RoleLabelsEditor`, `useAdminLanguages` |
| Rol oluşturma (`RoleForm.tsx`, değişti) | Mevcut sayfa | `RoleLabelsEditor`, rol anahtarı doğrulaması |

Gezinme: `ApprovalsHero` iki sayfaya bağlantı verir. Yan menü girişleri bu farkın parçası değildir. Menü tohumu ayrıca yapılır (Bölüm 10.4).

### 11.2 Onay politikası ekranı

- Kapı listesi `GET /admin/approval-gates` ile alınır. Her kapı için `summarizeGate` özet üretir: zincirse adım sayısı, eşiksizse "all", aksi halde eşik ve para birimi (`lib/approvals/policy.ts` 13-17).
- Düzenleme iletişim kutusu (`PolicyEditorDialog`): tek eşikli kapılarda eşik, para birimi ve izin. Zincir kapılarında adım düzenleyicisi.
- Zincir düzenleyicisi (`ChainStepsEditor.tsx`):
  - Ekle: boş rol, kaynak `amount`, boş eşik ile yeni adım (`policy.ts` 73-78).
  - Sil: onay istemeden kaldırır (115-124).
  - Sırala: komşu ile yer değiştirir. İlk adımda yukarı, son adımda aşağı devre dışıdır (94-114; `moveItem` 113-119).
  - Kaynak `dealer_credit_limit` ise eşik alanı devre dışı kalır ve ipucu gösterilir (89, 127-129).
  - Rol seçimi rol listesinden yapılır. Listede olmayan kayıtlı rol anahtarı ham anahtarla seçili kalır (31-36).
  - Boş listede `approvals.ui.step.empty` gösterilir (40-42).
- Doğrulama (`lib/approvals/policy.ts` 89-110): boş liste `approvals.error.steps_required`; boş rol veya geçersiz kaynak `approvals.error.step_invalid`; tutar adımında negatif veya sayısal olmayan eşik `approvals.error.invalid_threshold`. `dealer_credit_limit` adımı 0 olarak gönderilir. Yinelenen roller reddedilmez.
- Adımlar dizideki sırayla gönderilir. İlk adım önce sorulur (`ChainStepsEditor.tsx` 24).

### 11.3 Onay talepleri ekranı

- Liste `GET /admin/approval-requests?status=&limit=100` ile alınır. Yalnızca ilk sayfa çekilir. Sayfalama kontrolü yoktur.
- Varsayılan filtre `pending`tir. `?request=<id>` sorgusu ilgili satırı vurgular ve filtreyi tüm durumlara genişletir (`approval-requests/page.tsx` 37, 42).
- Onayla ve reddet düğmeleri yalnızca şu koşullarda görünür (`page.tsx` 90-93): talep `pending`; `hasPermission(permissions, request.required_permission)` doğru; talebi açan kişi mevcut kullanıcı değil.
- Red ve onay notu `ReasonModal` ile alınır. Not isteğe bağlıdır (270-280).
- Uç çağrıları: `POST /admin/approval-requests/{id}/approve` ve `/reject`.

### 11.4 Sipariş sayfası iade bildirimi

- `usePaymentIntentRefund(intentId)`: anahtar tutar ve para birimi başına üretilir (`lib/approvals/refund.ts` `nextRefundAttempt`, 32-40). Başarıdan sonra sıfırlanır.
- Gönderim `POST /admin/payment-intents/{id}/refunds` çağrısını yapar. Her gönderim `idempotency_key` alanını içerir (`hooks/usePaymentIntentRefund.ts` 52). `approval_request_id` bekleyen kimlikse eklenir (`hooks/usePaymentIntentRefund.ts` 42-78).
- Sonuç sınıflandırması (`classifyRefundFailure`, `refund.ts` 51-56):
  - `refunded`: iade döndü. Bekleme ve `notGranted` temizlenir.
  - `approval_required`: üst düzey `approval_request_id` dolu. Bekleme durumu tutarla birlikte kaydedilir (61-66).
  - `not_granted`: hata `approvals.error.not_granted`. Bekleme korunur ve `notGranted` işaretlenir (67-71).
  - `failed`: diğer hatalar.
- Bekleme sırasında "İadeyi işle" düğmesi yerine `RefundApprovalNotice` gösterilir (`orders/[id]/page.tsx` 708-711). Bildirim tutarı, para birimini ve talep kimliğini gösterir. `/admin/approval-requests?request=<id>` bağlantısı verir. Yeniden dene aynı anahtar ve aynı onay kimliğiyle gönderir. İptal bekleme durumunu temizler ama anahtarı korur (`hooks` 81-84).

### 11.5 Rol sayfası ve etiket düzenleyicisi

- `dcf9d38` ile rol sayfasındaki `alert(...)` çağrısı sayfa içi mesaja çevrildi: `setMessage(t("roles.form.errors.nameRequired")); return;` (`roles/[id]/page.tsx` 49-53). Anahtar önceden vardı. Bu iki dosyada (`roles/[id]/page.tsx`, `RoleForm.tsx`) `alert(` çağrısı kalmamıştır.
- `saveRole` sırası (`roles/[id]/page.tsx` 49-97): ad denetimi (sayfa içi mesaj); sistem rolü koruması (`roles.form.notes.systemLimited`, istek yok); etiket dönüşümü (`roleLabelsFromRows`; geçersiz satır `permissions.error.role_label_invalid` verir ve istek yapılmaz); izinler; `PUT /admin/roles/{id}` (`name`, `description`, `permissions`, `is_system`, `labels`); başarıda `roles.edit.success`, hatada `roles.edit.fail`.
- Düzenleme sayfasında rol anahtarı salt okunurdur ve ipucu `approvals.ui.roles.key_locked` gösterilir (173, 179).
- Etiket düzenleyicisi (`RoleLabelsEditor.tsx`): her dil için bir satır. Bir dil diğer satırlarda kullanılıyorsa gizlenir. Metin en fazla 100 karakterdir. Boş satırlar `roles.ts` içinde atlanır (42).
- Oluşturma formunda anahtar hatası anahtarın altında gösterilir (`permissions.error.role_key_invalid`). Anahtar ve etiketler kaydetmeden önce denetlenir (`RoleForm.tsx` 180, 185).
- Sistem rolü `admin`dir (`is_system: true`). `editor`, `reviewer` ve `fleet_viewer` sistem rolü değildir. Panelin sistem kilidi yalnızca `admin` için geçerlidir.

### 11.6 Çeviri anahtarları ve sabit metin

- Değişen panel dosyalarında 98 farklı sabit anahtar vardır: 79 `approvals.*` (73 `approvals.ui.*`, 6 `approvals.error.*`), 2 `permissions.error.*`, 2 `menu.approvals.*`, ve önceden var olan `roles.form`, `roles.edit` ve `common` anahtarları.
- Bu 98 anahtarın tümü `internal/seeds` içinde dizge olarak bulunur (`git grep` ile doğrulandı; olumsuz kontrol de yapıldı).
- Dinamik anahtarlar: kapı etiketleri (`approvals.gate.payments_refund`, `approvals.gate.b2b_purchase_order`), sunucu hata anahtarları (`apiErrorMessage` ile), doğrulama anahtarları (`policy.ts`, `roles.ts`).
- `approvals.ui` ön ekleri: `action` 6, `field` 20, `filter` 1, `hint` 5, `policies` 9, `refund` 4, `requests` 6, `roles` 8 (`labels` dahil), `source` 2, `status` 4, `step` 5, `toast` 3. Toplam 73.
- Yeni JSX'te kullanıcıya görünen İngilizce sabit metin yoktur. Yalnızca tipografik metinler vardır: boş değer için em dash ve ayırıcı olarak orta nokta.
- Önceden var olan ve değişmeyen metinler: sipariş sayfasındaki yer tutucular (`0.00`, `Enter refund reason...`). İade bildirimlerinin varsayılan İngilizce metinleri `t(key, default)` içinde kalır.

### 11.7 Panel API istemcisi

`lib/api/approvals.ts` ve `lib/api/payment-intents.ts`:

| İşlev | Yöntem ve yol | Yanıt tipi | Backend ile eşleşme |
|---|---|---|---|
| `listGates` | GET `/admin/approval-gates` | `ApprovalGate[]` (`items`) | Eşleşir |
| `putPolicy` | PUT `/admin/approval-policies/{gate}` | `ApprovalGate` | Eşleşir |
| `listRequests` | GET `/admin/approval-requests?status=&limit=100` | `ApprovalRequest[]` | Eşleşir |
| `approveRequest` | POST `/admin/approval-requests/{id}/approve` | `ApprovalRequest` | Eşleşir |
| `rejectRequest` | POST `/admin/approval-requests/{id}/reject` | `ApprovalRequest` | Eşleşir |
| `paymentIntentsApi.refund` | POST `/admin/payment-intents/{id}/refunds` | `PaymentIntent` | Eşleşir |

Panel, `GET /admin/approval-policies/:gate` ucunu kullanmaz. Ekranlar tüm kapı listesini okur.

Hata gövdesindeki `error` alanı panelde çevrilir. `^(approvals|permissions)\.error\.[a-z0-9_]+$` deseniyle eşleşen anahtarlar `apiErrorMessage` ile çevrilir (`lib/approvals/errors.ts` 4, 22-28). Diğer metinler olduğu gibi gösterilir.

Panel tipleri (`lib/types/approvals.ts`) backend JSON adlarıyla eşleşir. Panel, `status` alanı yoksa `pending` değerini, `id` yoksa `_id` değerini kullanır. Bilinmeyen `threshold_source` değeri `amount` olarak okunur. Bu geri dönüşler eksik veriyi gizleyebilir.

### 11.8 Panel sınırları

1. Talep listesi en fazla 100 satır ve yalnızca ilk sayfadır.
2. Yetki `approval_policies:read` yoksa ham kapı anahtarı gösterilir.
3. Zincir düzenleyicisi yinelenen onaylayıcı rollerini kabul eder.
4. Tek eşikli politikada eşik boşsa para birimi gönderilmez (`policy.ts` 56-57).
5. Sistem rolü `admin` düzenleme sayfasından kaydedilemez (önceden var olan koruma).
6. Talebi açan kişi UI'de kendi talebini onaylayamaz. Backend de reddeder.
7. `lib/approvals/roles.ts` içindeki istemci denetimleri yalnızca erken geri bildirimdir. Yetkili olan sunucudur.

### 11.9 Araç formu: plaka değişikliği gerekçesi (`7b03fe7`)

- Araç düzenleme formunda plaka girdisi kayıtlı plakadan farklıysa gerekçe alanı görünür ve zorunludur (`app/(dashboard)/admin/stock/vehicles/VehicleForm.tsx`). Karşılaştırma büyük harfe çevrilmiş ve kırpılmış değerlerle yapılır.
- Kayıt düğmesi gerekçe girilene kadar beklenir.
- Güncelleme isteği `plate_reason` alanını gönderir (`lib/api/stock-ecosystem.ts`).
- Sunucu hataları çevrilerek gösterilir (`t(message, message)`). Koordinatör bildirimine göre gerekçe kuralları ham anahtar yerine metin olarak görünür. Çalışan panelde bu davranış ayrıca doğrulanmamıştır.

---

## 12. Test yaklaşımı ve sonuçları

### 12.1 Yaklaşım

- Test fonksiyonları `Test{Method}_{Scenario}` biçimine getirildi ve assertion'lar testify'a çevrildi (commit `f31c571`).
- Yeniden adlandırma sayısı: `f31c571` commit'i 20 dosyayı değiştirdi. Diff'te 70 `func Test` satırı kaldırıldı ve 70 yeni ad eklendi. Yani bu commit 70 test fonksiyonunu yeniden adlandırdı.
- Taban commit (`f48748f`) ile karşılaştırıldığında aynı 20 dosyada 95 test fonksiyonu vardı. Son durumda (`f286fd0`) 185 vardır. Bu karşılaştırmada 91 test adı tabana göre yenidir (yeni eklenen veya yeniden adlandırılan). Bir taban adı kalkmıştır: `TestFilterMenuTreeByRoles_ImpersonatedPlatformUserUsesTenantRoles`, yeni adı `TestFilterMenuTreeByAccess_ImpersonatedPlatformUserUsesTenantRoles`.
- Sahibin listesindeki "91 test fonksiyonu f31c571 ile yeniden adlandırıldı" ifadesi bu taban karşılaştırmasına karşılık gelir. Commit diff'inde yalnızca 70 yeniden adlandırma vardır. Bu fark burada kayda geçirilmiştir.
- Yeni testler: `internal/seeds/role_seed_test.go` içinde `TestSeedRoles_KeepsEditorPermissionsAfterAdminChange` ve `TestGrantRoleDefaults_AddsMissingDefaultsAndKeepsCustomOnes`. `internal/modules/b2b/dealers/service/po_approval_test.go` içinde `TestApprove_RequiresEarlierStepFirst` ve `TestReject_RequiresCurrentStepRole`.
- Mutasyon kontrolü: rol tohumlamasının eski üzerine yazma davranışı geri konduğunda rol tohumlama testi başarısız olur (sahibin bildirimi).
- Mongo bağlantılı testler yerel bir replica set'e karşı çalıştırıldı (sahibin bildirimi). Ortam değişkeni tanımlı değilse bu testler atlanır.
- Gürsel GPS regresyon testi: `TestSeedGurselTurizimGpsProviders_KeepsAdminChangesOnRerun` (`c27c3b5`). Bu testin sonucu sahibin doğrulama listesinde yer almaz. Bu belgede sonucu iddia edilmez.

### 12.2 Sahibin doğrulama sonuçları

Bu sonuçlar sahibin bildirdiği şekilde, `f286fd0` ağacı üzerinde alınmıştır. Bu belgeleme aşamasında yeniden çalıştırılmamıştır. Sonraki commit'ler (`c27c3b5`, `a00f69d`, `de3c4fa`, `18e27f8`, `33467c8`) için sahibin listesinde sonuç yoktur. `a00f69d` ile değişen `handler_mongo_test.go` testinin sonucu bu belgede yoktur.

Backend:

- `go build ./...` geçer. Bu iş kapsamındaki 103 mevcut değiştirilmiş Go dosyası `gofmt` ile temizdir.
- `go vet ./...` iki tanı verir. İkisi de bu işin dokunmadığı paketlerdedir: `personal_assistant/service` test dosyasında derleme hatası ve `stock/subscription_pricings/service` içinde kendine atama uyarısı.
- `go test` geçer: approvals (repository, service); b2b/dealers/service; wallets (repository, service, topup_sources/service); fleet_tally/trip_assignments (repository, service); permissions (handler, service); returns/service; menus (handler, repository); access_points/service; platform/middleware; payments/intents (handler, service); billing/invoices (adapters, handler, service); orders/handler; service_reservations/service; ve seeds içinde rol, provisioning, çeviri kapsamı, menü ve profil testleri.
- Yedi GİB sandbox adaptör testi atlanır. Bu atlamalar taban commit'te de vardır. Düzeltme: atlama mesajlarından yalnızca biri `GIB_ENVIRONMENT` adını içerir ("GIB_ENVIRONMENT not set to 'test', skipping sandbox test", `gib_sandbox_test.go` 25). Diğer altı atlama, GİB adaptörünün yapılandırılmadığını (`GIB adapter not configured`) veya gerçek GİB kimlik bilgilerinin gerektiğini (`requires actual GIB credentials`) belirtir. Sahibin ifadesi bu ayrıma göre düzeltilmiştir.
- Bu işten bağımsız, önceden var olan bir başarısızlık: `internal/seeds` içindeki `TestNewsSite*` testleri commit `7c6ae1e`'de de başarısızdır.
- Rota kontrolü: `cmd/api` son ağaçtan derlendi ve `APP_ENV=test` ile, seed'ler kapalı, geçici bir veritabanına karşı 40 saniye çalıştırıldı. Panik yoktu, rota çakışması yoktu ve sunucu dinledi. Kaydedilen onay rotaları: `GET /admin/approval-gates`; `GET` ve `PUT /admin/approval-policies/:gate`; `GET /admin/approval-requests`; `POST /admin/approval-requests/:id/approve` ve `/reject`; `GET /admin/fleet-tally/trip-assignments/pending-approval`.

Panel:

- `npm test` 41 testin 39'unu geçer. İki başarısızlık, bu işin dokunmadığı dosyalardaki kiracı kimliği testleridir ve `main` üzerinde de başarısızdır.
- `tsc --noEmit` 328 hata verir. Bu sayı bu işten önceki sayıyla aynıdır.
- `npm run check-i18n` 1 kodu ile çıkar. Eksik anahtarlar `main` ile aynıdır ve önceden vardır.
- Değişen dosyalarda ESLint üç `no-explicit-any` hatası verir. Bu satırlar `main` üzerinde zaten vardı.

### 12.3 Bu belgeleme aşamasında yapılan kontroller

Bu belgenin yazımı sırasında kod okuması yapıldı. Aşağıdakiler doğrudan kaynak dosyada doğrulandı:

- Kapı kataloğu, varsayılan politika, rota kapıları (approvals, wallets, sefer incelemesi, B2B).
- Onay durum değerleri ve hata anahtarları; `statusFor` eşlemesi.
- Onay kaydı ve politika JSON alanları; PolicyRequest ve DecisionRequest gövdeleri.
- Tohum sırası, onay politikası tohumu, `SEED_ON_START` değeri.
- İade için claim ve release akışı; returns için bekleyen hata ayrımı.
- Yanıt zarfının ek alanları köke koruduğu.
- `api` modül çıkarımı ve B2B rota zinciri.
- `wallets:approve` tohum eksikliği.
- Mock e-invoice sağlayıcısının kaydı.
- Commit sayıları, test adı farkları ve GİB atlama mesajları.
- Gürsel GPS seed değişikliği (`c27c3b5`).

Çalıştırılarak doğrulanmamış olanlar: canlı istekler, Mongo testleri, panel testleri ve adaptörlerin gerçek sağlayıcı davranışı.

---

### 12.4 Son doğrulama (74a5453; koordinatör bildirimi)

Aşağıdaki sonuçlar koordinatör tarafından `origin/aypass-approval-ui` HEAD `74a5453` üzerinde bildirilmiştir. Bu belgede yeniden çalıştırılmamıştır. Sahibin önceki sonuçları `f286fd0` ağacına aittir.

- `go build ./...` geçer.
- `go vet ./...` yalnızca iki önceden var olan bulgu verir.
- Hedefli testler geçer: ledger, wallets, stock vehicles, access points, payments handler, b2b, middleware, invoices, seeds (çeviri, rol, GPS, izin).
- Ledger ve wallets eşzamanlılık testleri `-race` ile geçer.
- Route probe: `/admin/stock/vehicles` ve `/admin/b2b/purchase-orders` 401 döndü (rota var, kimlik doğrulama istiyor). Eski `/api/admin/b2b` 200 döndü; bu yolda rota kaydı yoktur. Bu 200 yanıtının kaynağı doğrulanmadı. Kod tarafında `NoRoute` işleyicisi `cmd` ve `internal/platform` altında bulunmadı.
- Panik ve rota çakışması yok.
- Panel (`origin/approval-policy-ui` HEAD `7b03fe7`): `npm test` 39/41 (2 önceden var olan başarısızlık); değişen dosyalarda ESLint 0 sorun; `tsc` 328 (taban sayısı); değişen dosyalarda tip hatası yok; i18n için yeni eksik anahtar yok.
- `gofmt`: `internal/platform/middleware/module_gate.go` taban commit `f48748f`'de de biçimlendirilmemiş (`gofmt -l` bu belgede doğrulandı). Bu dal dosyaya yalnızca alias satırını ekledi.

---

## 13. Operasyon notları

### 13.1 Belgeler

`README_FOR_AI_ASSISTANTS.md` kuralı gereği operasyon ve installer belgelerinin bu iş kapsamında gözden geçirilmesi gerekir. Belgelerin durumu:

| Belge | Onay ve seed içeriği |
|---|---|
| `user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md` | Onay politikası ve seed bölümleri 10.4 (995-1049 satırları), 15.3 ve 15.7 bölümlerindedir |
| `user-guide/05-operations/OPERATIONS_GUIDE.md` | 3.15 bölümü (141-150 satırları) |
| `user-guide/05-operations/DEPLOYMENT_SCENARIOS.md` | 8. bölüm (271-280 satırları) |
| `/workspace/go/aystek-core-installer/README.md` | Commit edilmemiş. Seed davranışı 87-91, operatör adımları 96-99 |
| `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md` | Commit edilmemiş. 20.6 (2220-2262), 20.8 (2287-2306) ve 396. satır |

Eski uyarı: `OPERATIONS_GUIDE.md` 3.15 ve `KULLANIM_KILAVUZU.md` 20.8 bölümleri, Gürsel GPS yapılandırmalarının her tam seed zincirinde sıfırlandığını yazar. Commit `c27c3b5` bu belgeleri değiştirmemiştir. Bu uyarılar artık eski davranışı anlatır ve güncellenmelidir. Bu belge kapsamında düzeltilmemiştir.

### 13.2 Kurulum sonrası yapılacaklar

1. `store_manager`, `credit_dept` ve `executive` rollerini oluşturun. Bunlar olmadan B2B zincirindeki adımlar onaylanamaz.
2. Onay politikaları ekranında `b2b.purchase_order` zincirini ve `payments.refund` kapısını gözden geçirin.
3. İade kapısını açmadan önce bir eşik belirleyin. Eşiksiz etkin kapı her iadeyi onaya bağlar.
4. Mevcut kurulumlarda `aystek-seed --roles` komutunu çalıştırın. Bu komut eksik varsayılan izinleri ekler ve hiçbir izni kaldırmaz.
5. `approvals.error.chain_not_configured` (HTTP 400, B2B gönderiminde) ve `approvals.error.required` (HTTP 409, iade) yanıtlarını izleyin. İkincisi beklenen davranıştır.
6. `SEED_ON_START` değerini kurulumdan sonra `false` tutun.
7. İade isteği gönderen her istemci `idempotency_key` göndermelidir. Anahtarsız istek 400 `payment_intents.idempotency_key_required` ile döner. Panel anahtarı gönderir. Diğer istemciler bu belgede denetlenmedi.
8. B2B uçları `/admin/b2b` altındadır (`de3c4fa`). Eski `/api` önekli yolları kullanan istemciler güncellenmelidir.
9. Production ortamında mock e-invoice ve ERP adaptörleri kayıtlı değildir (`33467c8`). Mock PDF renderer, mock kargo ve mock SMS sağlayıcıları tüm ortamlarda kayıtlı kalır (Bölüm 14.3, risk 3).
10. Araç plakası değiştirilirken gerekçe zorunludur (kırpılmış, en fazla 500 karakter). Eksik gerekçe `stock.vehicles.validation.plate_reason_required`, uzun gerekçe `stock.vehicles.validation.plate_reason_too_long` ile döner (`93b4232`).

### 13.3 Seed ve kapsayıcı

- Tam seed zinciri (`SEED_ON_START=true`, `SEED_EXIT_AFTER_RUN=true`, `SEED_FULL_RESET=false`) yalnızca installer'ın tek seferlik kapsayıcısında çalıştırılır (`docker compose run --rm`).
- `--roles` yalnızca eksik rollere ve izinlere dokunur.
- `PUT /admin/approval-policies/:gate` politikayı tamamen değiştirir. Gönderilmeyen alanlar varsayılana döner. Yönetici, değişiklik öncesi mevcut politikayı okuyup kopyalamalıdır.

### 13.4 Üretim durumu

- Üretim sunucusu (`aypass.gurseltur.com.tr`) bu ortamdan erişilemedi. Canlı durum doğrulanmamıştır.
- Flutter uygulamaları `aypass_erisim_app` ve `aypass_kullanici_app` bu çalışmada derlenmedi ve doğrulanmadı.
- Installer belgeleri commit edilmemiştir.

### 13.5 Gizli bilgiler

- Depo `.env` dosyasını izlemektedir. Proje notlarına göre bu dosya SMTP ve Firebase gizli bilgileri içermekte ve bu bilgiler gönderilmiş (push) durumdadır. Bu belgede hiçbir değer yazılmadı ve dosya okunmadı. Üretimden önce bu bilgiler döndürülmeli ve git geçmişi temizlenmelidir.
- `cmd/migrate_gursel/main.go` taban commit'te sabit kodlanmış bir kaynak bağlantı bilgisi içeriyordu. `f286fd0` bu değeri `MIGRATE_SOURCE_URI` ortam değişkenine taşır. Değer geçmişte kaldığı için döndürme listesinde yer alır. Değer bu belgede yazılmadı.
- Sunucu bağlantı bilgisi içeren dosyalar bu belgenin kapsamında okunmadı.

---

## 14. Açık konular, riskler ve kararlar

### 14.1 Sahip kararları ve uygulanmış kararlar

Bölüm 1.4'te yer alan kararlar:

- AD-04 onay kayıtları: silinip yeniden tasarlandı. Onaylandı.
- Eşik ve B2B zinciri politika verisidir. Rol anahtarları İngilizce ve değişmezdir. Etiketler dile göredir.
- Banka havalesi ve sefer onayları politika değil, izin artı kendi kendine onay kontrolüdür.
- Seed test kimlik bilgileri kodda kalır. Üretimde değiştirilecektir.
- Editör menü görünürlüğü değişikliği onaylandı.
- Ana dil incelemesi ertelendi.
- `editor`, `reviewer`, `fleet_viewer` panelden yönetilir; seed'ler üzerine yazmaz.

Uygulanmış karar (2026-10-28, `c27c3b5`):

- Gürsel GPS sağlayıcı yapılandırması yalnızca ilk oluşturmada yazılır. Tam seed zinciri, güncelleme ve konteyner başlatma mevcut kayıtların `credentials`, `base_url`, `is_active`, `display_name` ve `timezone_offset` alanlarını değiştirmez. Panelden yapılan değişiklik korunur. Regresyon testi: `TestSeedGurselTurizimGpsProviders_KeepsAdminChangesOnRerun`. Bu konu artık açık risk değildir.

Uygulanmış değişiklik (2026-10-28, `a00f69d`): admin iade işleyicisi anahtarsız istekte yeni anahtar üretmez. Anahtar gövdeden ya da `Idempotency-Key` / `X-Idempotency-Key` başlığından gelir ve zorunludur. Boş anahtar servis tarafından `payment_intents.idempotency_key_required` ile 400 olarak reddedilir. Anahtarsız istek artık kabul edilmediği için, anahtarsız ağ tekrarının ikinci bir iade yaratması bu yolda mümkün değildir.

Bu turda uygulanan diğer değişiklikler (sahip kararı değildir): `de3c4fa` (B2B rotaları `/admin/b2b` altına taşındı ve `b2b` takma adı `b2b_dealers` modülüne çözülüyor), `18e27f8` (`wallets:approve` tohumlanıyor), `33467c8` (mock e-invoice ve ERP adaptörleri yalnızca production dışında kaydediliyor). Birleşme `74a5453` ile P0 dalı (`7eb7bfa`: sistem hesabı oluşturma ve defter onayı yarış güvenli; `93b4232`: plaka değişikliğinde gerekçe zorunlu; `25abcc1`: test fake güncellemesi) origin dalına girdi. Panel `7b03fe7` plaka gerekçesi alanını ekledi.

### 14.2 Açık kararlar

Aşağıdaki dört konu açık karar olarak kalır:

1. `cmd/migrate_gursel` ikili dosyasının (binary) git'ten kaldırılması. Sahibin kararı bekleniyor.
2. Kimlik bilgisi rotasyonu ve git geçmişinin temizlenmesi. Üretimden (go-live) önce yapılmalıdır. Kapsam: `.env` dosyası ve taban commit'teki `cmd/migrate_gursel` içindeki kaynak bağlantı bilgisi. Değerler bu belgede yazılmadı.
3. Ana dil incelemesi: `ky`, `kk`, `mn`, `hy`, `az-cyrl`, `uz`, `tk`. Bir ülke projesi kapsama girene ya da iş tamamlanana kadar ertelenmiştir (sahip kararı).
4. İade serbest bırakma kararı (karar bekliyor): sağlayıcı hatasında rezervasyon ve onay serbest bırakılıyor. Zaman aşımında sağlayıcı iadeyi yapmış olabilir. Rev.06 analiz belgesinin MK-07 satırı (PDF sayfa 27) ödeme durumu olarak PENDING_EXTERNAL ister; zaman aşımını ret saymamayı ve aynı intent ile sağlayıcı referansıyla sorgulamayı öngörür. Karar, bu sorgulamanın tasarlanıp tasarlanmayacağına ve serbest bırakma kuralının buna göre daraltılıp daraltılmayacağına ilişkindir.

### 14.3 Riskler

Önem sırasına göre. Her madde kod okumasına dayanır ya da bu belgede doğrulanmamıştır. Belirtilen yerlerde bu açıkça yazılmıştır.

1. B2B rota yolu ve yetki (düzeltildi, doğrulanmadı): B2B uçları `de3c4fa` ile `/api` önekinden `/admin/b2b` altına taşındı; `b2b` takma adı `b2b_dealers` modülüne çözülüyor (`module_gate.go` 143). Önceki yapıda `/api/admin/b2b` yolu modül `api` olarak çözülüyor ve platform dışı çağıranlar `missing module permission` ile 403 alıyordu (kod okumasına göre). Bu durum taban commit'te de vardı. Koordinatör probe (`74a5453`): `/admin/b2b/purchase-orders` 401 döndü (rota var, kimlik istiyor); eski `/api/admin/b2b` 200 döndü, bu yolda rota kaydı yok ve yanıtın kaynağı doğrulanmadı. Yetkili bir istekle doğrulama bu belgede yapılmadı. Eski yolu kullanan istemciler güncellenmelidir.
2. `wallets:approve` izni artık tohumlanıyor (`18e27f8`, Bölüm 10.2). Platform dışı yöneticiler bu izni içeren bir rolle banka havalesini onaylayabilir. Uçtan uca doğrulama bu belgede yapılmadı.
3. Sahte sağlayıcılar (açık risk, düzeltilmedi): mock e-invoice ve mock ERP adaptörleri `33467c8` ile yalnızca production dışında kaydediliyor (`cmd/api/main.go` 1448-1468). Mock PDF renderer (`cmd/api/main.go` 1477), mock kargo sağlayıcısı (`cmd/api/main.go` 1353) ve mock SMS sağlayıcısı (`internal/modules/sms/service/service.go` 47) her ortamda kayıtlı kalır. Adlarından anlaşıldığı üzere sahte sağlayıcılardır; davranışları bu belgede ayrıca doğrulanmamıştır. Production dışında `mock` e-invoice sağlayıcısı hatasız `sent` sonucu döner (Bölüm 4.3.6).
4. İade serbest bırakma kararı bekliyor: `adapter.Refund` hata döndürdüğünde rezervasyon ve onay serbest bırakılır (`service.go` 581, 586, 605; Bölüm 6.5, 8.2). Zaman aşımında sağlayıcı iadeyi yapmış olabilir; sorgulama yapılmadan serbest bırakma ikinci bir iade riski doğurabilir. Rev.06 analiz belgesinin MK-07 satırı (PDF sayfa 27) PENDING_EXTERNAL durumunu ve aynı intent ile sağlayıcı referansıyla sorgulamayı ister; zaman aşımı ret sayılmamalıdır. Sağlayıcıya iade anahtarı gönderilmiyor (`AdapterRefundInput` alanları: Amount, Currency, ProviderIntentID, Reason, Metadata; `adapters.go` 60-66). Karar bekliyor (Bölüm 14.2, madde 4).
5. Sağlayıcı kabul ettikten sonra `repo.Update` başarısız olursa iade 400 döner. Para hareket etmiştir, ama onay `consumed` kalır. Serbest bırakma API'si yoktur.
6. Çökme riski: talep tüketildikten sonra sağlayıcı çağrısından önce süreç çökerse talep `consumed` kalır. Kurtarma için yeni talep gerekir.
7. Bekleyen ve onaylanmış talepler süresiz kalır. Süre dolumu, iş veya zamanlayıcı yoktur.
8. Onay talebi iade kaydına bağlanmaz. `approval_request_id` ödeme niyetinde, işlemde veya iade ödemesinde saklanmaz. İzleme yalnızca onay talebinden yapılabilir.
9. Varsayılan B2B zincirindeki `store_manager`, `credit_dept` ve `executive` rolleri hiçbir seed tarafından oluşturulmaz. Oluşturuluncaya kadar eşiği aşan siparişler bekleyen kalır.
10. Varsayılan zincir: `approvals_seed.go`, varsayılan kiracının (`aystek`) varsayılan sitesi için etkin üç adımlı B2B zinciri ekler (store_manager 50.000 tutar eşiği, credit_dept bayi kredi limiti, executive 100.000). Politikası kayıtlı olmayan veya devre dışı tenant ve site kapalı başarısız olur. Bu nedenle "varsayılan zincir yok" ifadesi aystek kiracısı için geçerli değildir (Bölüm 3.5).
11. Boş zincirde otomatik onay: Hiçbir adım uygulanmıyorsa sipariş gönderildiği anda onaylanır (`po_service.go` 262-266). Davranış belgelenmiştir.
12. Satın alma siparişi gönderiminde bayi okuma hatası yok sayılır (`po_service.go` 248: `dealer, _ := ...`). Hata durumunda kredi limiti 0 kabul edilir ve `credit_dept` adımı düşer (kod okumasına göre).
13. Banka havalesi onayında kısmi hata: defter kaydı `posted` olduktan sonra `applyJournal` başarısız olursa yeniden onay bakiyeyi uygulamadan başarılı döner (kod okumasına göre; `service.go` 656-658 ve 678). Tutarsızlık yeniden mutabakatta görünür.
14. Rol adları benzersiz değildir. `CreateRole` tekillik denetimi yapmaz. Roles koleksiyonunda dizin bulunmadı (kod deposunda arama yapıldı, canlı veritabanı kontrol edilmedi).
15. Ödeme iadesi kapısı eşiksiz açılırsa her iade onay ister (Bölüm 3.6).
16. Etiket geri dönüşü backend'de çağrılmaz (`LabelFor`). Panel geri dönüşü kendisi uygular.
17. Rol güncelleme anlamı: boş `labels` nesnesi ve boş `permissions` dizisi alanları temizler. Panel bunları yanlışlıkla göndermemelidir.
18. Analiz belgesindeki `ensureSystemAccount` yarış durumu bu işte doğrulanmadı.
19. `go vet` iki tanı verir (Bölüm 12.2). Bu işin dokunmadığı paketlerdedir.
20. Test yardımcı fonksiyonları hâlâ `t.Fatalf` kullanır. Yalnızca `Test` fonksiyonları dönüştürüldü.
21. Flutter uygulamaları (`aypass_erisim_app`, `aypass_kullanici_app`) derlenmedi ve doğrulanmadı.
22. Panelde `tsc` hataları (328) ve `check-i18n` eksikleri önceden vardır ve bu iş tarafından çözülmemiştir.
23. Admin rolü düzenleme sayfasından kaydedilemez (önceden var olan koruma).
24. Anahtar zorunluluğu (`a00f69d`): iade için `idempotency_key` göndermeyen istemciler artık 400 alır. Panel gönderir; diğer istemciler bu belgede denetlenmedi.
25. Plaka geçmişi yazım hataları yutuluyor (Z02 açığı): `Create` (service.go 75) ve `UpdateVehicle` (service.go 128-129) geçmiş yazma hatalarını `_ =` ile yok sayıyor; geçmiş kaydı sessizce eksik kalabilir. `AddPlate` yazma hatalarını döndürür; ancak geçmişi kapatma, geçmiş ekleme ve araç güncelleme ayrı yazmalar olarak, işlem sınırı olmadan yapılır. Ara bir hata geçmiş kaydı ile araç kaydını ayrıştırabilir.
26. `stock.vehicles.validation.plate_required`, `duplicate_plate`, `status_invalid`, `invalid_payload` ve `stock.vehicles.error.not_found` anahtarlarının hiçbir seed'de yoktur. Panel bu hataları ham anahtar olarak gösterebilir (koordinatör bildirimi; çalışan panelde doğrulanmadı).
27. `internal/platform/middleware/module_gate.go` taban commit `f48748f`'de de biçimlendirilmemiştir (`gofmt -l` çıktısı). Bu dal dosyaya yalnızca alias satırını ekledi (`de3c4fa`); dosya biçimlendirilmemiş olarak kalır.
### 14.4 Rev.06 kabul kriterleri

Rev.06 kabul kriterleri karşılanmamaktadır. Z01 ile Z14 arasında tamamen karşılanan gereksinim yoktur. 104 kabul testinin kod durumu taban commit `f48748f` üzerinde hesaplanmıştır: 3 karşılanıyor, 34 kısmi, 45 yok, 22 çelişki. Bu sayılar bu işten sonra yeniden hesaplanmamıştır (Bölüm 15.1).

Model 1 (DV01-DV20, DT01-DT26) sağlayıcı transferi ve Model 4 (M01-M24, MK01-MK30, D01-D14, F01-F12) Küresel Pazar uygulanmamıştır.

### 14.5 Rev.06 analiz belgesindeki eskimiş ifadeler

`AYPASS_GECIS_SISTEMI_REV06_ANALIZ_VE_IS_PLANI_TR.md` bu işten önce yazılmıştır. Aşağıdaki maddeler artık güncel değildir. Belge bu iş kapsamında düzeltilmemiştir:

- §12.7 "Politika için seed dosyası yoktur": `approvals_seed.go` vardır.
- §12.7 "tek kapı": kataloğda iki kapı vardır.
- §12.7 "11 hata anahtarı ve 1 etiket": 15 hata anahtarı ve 2 kapı etiketi, 28 dilde.
- §12.7 "panel ekranları yapılmadı": `approval-policy-ui` dalında yapıldı.
- §8 ve §11 "AD-04 tasarım açık": yeniden tasarım uygulandı.

---

## 15. Rev.06 eşlemesi ve sonraki adımlar

### 15.1 Gereksinim kimlikleri ve kabul testleri

Kaynak: `AYPASS_REV06_EXCEL_ESLEME_TR.md` (tarih 2026-10-28, kod referansı `f48748f`). Durum sözcükleri: KARŞILANIYOR (tamamen), KISMİ, YOK, ÇELİŞKİ.

Gereksinimler:

| Grup | Toplam | KARŞILANIYOR | KISMİ | YOK | ÇELİŞKİ |
|---|---|---|---|---|---|
| Z01-Z14 | 14 | 0 | 5 | 2 | 7 |
| M01-M24 | 24 | 0 | 10 | 10 | 4 |
| DV01-DV20 | 20 | 0 | 4 | 16 | 0 |

Kabul testleri (104):

| Grup | Toplam | KARŞILANIYOR | KISMİ | YOK | ÇELİŞKİ |
|---|---|---|---|---|---|
| K01-K24 | 24 | 1 | 7 | 8 | 8 |
| MK01-MK30 | 30 | 0 | 12 | 12 | 6 |
| DT01-DT26 | 26 | 0 | 3 | 22 | 1 |
| T01-T15 | 15 | 1 | 6 | 2 | 6 |
| FX1-FX9 | 9 | 1 | 6 | 1 | 1 |
| Toplam | 104 | 3 | 34 | 45 | 22 |

Bu sayılar yeniden hesaplanmamıştır. Bu işin yeni kodu bu sayıları değiştirmiş olabilir. Yeniden hesaplama yapılmalıdır.

### 15.2 Faz 0 durumu (Rev.06 analiz §8, 2026-10-28)

| Madde | Durum | Not |
|---|---|---|
| 0.1 RefundEvent | Tamam | Yalnızca `LedgerEntryID` olan kabul edilmiş olaylar iade edilir |
| 0.2 GİB | Tamam | Hatalar döner, üretimde sahte UUID yok, `ErrAlreadyIssued` (409). Açık: iki eşzamanlı gönderim için atomik koruma |
| 0.3 İade idempotensi | Kısmi | `refunded_amount` üzerinde koşullu güncelleme; sağlayıcı hatasında serbest bırakma; `ErrIdempotencyConflict` (422). Kapatıldı (`a00f69d`): anahtarsız admin iadesi yeni anahtar üretmez, 400 `payment_intents.idempotency_key_required` döner. Açık: sağlayıcı tarafında iade anahtarı yok; çökmede kayıt bekler kalır |
| 0.4 Cüzdan eşzamanlılığı | Tamam | `$set` kaldırıldı; `ApplyBalanceChange` ile `$inc` ve taban filtresi; son 50 defter kaydı penceresi |
| 0.5 Callback tekrarı | Tamam | İşlem kaydı kancadan sonra yazılır; kanca hatası 503 |
| 0.6 Defter dizini | Tamam | Koordinatör bildirimi (`74a5453`): sistem hesabı oluşturma ve defter onayı yarış güvenli hale getirildi (`7eb7bfa`). Kalan açıklar: `ledger_accounts` üzerindeki unique indeks veri denetimi bekliyor (indeks henüz tanımlı değil); soft-delete sonrası yeniden açma ve eski ile yeni binary aynı anda çalışırken mükerrer hesap riski; canlı replica set doğrulaması; `ConfirmBankTransfer` yeniden denemesi |
| 0.7 approvals (AD-04) | Kısmi | Kullanılmayan modül silindi (`eb581c7`); yeniden tasarım yapıldı. Açık: kiracı ve site kapsamının kanıtı, talep eden ile onaylayanın farklılığı ve kalıcı karar kaydı |
| 0.8 Sefer atamaları | Tamam | `TestUpdate_PersistsApprovalAndDriverFields` gerçek MongoDB'de |
| 0.9 Pasif erişim noktası | Tamam | `ErrAccessPointInactive`; `TestRecordEvent_RejectsInactiveAccessPoint` |
| 0.10 Sabit kimlik bilgileri | Kısmi | Göç kaynağı `MIGRATE_SOURCE_URI` veya `-src-uri`. Sahip kararı: seed test değerleri kodda kalır. Açık: `cmd/migrate_gursel` ikili dosyası ve git geçmişinde rotasyon |
| 0.11 Stok plakası aktörü | Tamam | Aktör oturumdan alınır. Plaka değişikliğinde gerekçe zorunludur (kırpılmış, 1–500 karakter; yalnızca plaka gerçekten değiştiğinde). Kayıt `stock_vehicle_plate_history.reason` alanına yazılır (`93b4232`; panel `7b03fe7`) |

0.6 ve 0.11 kapatıldı (koordinatör bildirimi; `origin/aypass-approval-ui` HEAD `74a5453`). Analiz belgesindeki durum sütunu hâlâ Kısmi olarak kayıtlıdır; o belge bu işte güncellenmedi.

Taahhütler (`origin/aypass-approval-ui` içinde): `b2398bd`, `884b0a7`, `c79cc9b`, `eb581c7`. 0.3 için ek commit: `a00f69d`. 0.6 ve 0.11 için: `7eb7bfa`, `93b4232`, `25abcc1` (birleşme `74a5453`).

### 15.3 Açık kararlar

Dört karar açıktır (Bölüm 14.2):

1. `cmd/migrate_gursel` ikili dosyasının kaldırılması.
2. Üretimden önce kimlik bilgisi rotasyonu ve git geçmişinin temizlenmesi.
3. Ana dil incelemesi (ertelendi).
4. İade serbest bırakma kararı (karar bekliyor; Bölüm 14.3, risk 4).

### 15.4 Teknik sonraki adımlar

Aşağıdaki maddeler teknik işlerdir. Karar gerektirmezler.

1. B2B rotalarının (`/admin/b2b`, `de3c4fa`) çalışan bir istekle doğrulanması (Bölüm 14.3, risk 1). En öncelikli adım.
2. Banka havalesi onayının (`wallets:approve`, `18e27f8`) uçtan uca doğrulanması (Bölüm 10.2).
3. Mock PDF renderer, mock kargo ve mock SMS sağlayıcılarının production için kapatılması veya kayıttan çıkarılması (Bölüm 14.3, risk 3).
4. Adaptör sözleşmesinin netleştirilmesi: `adapter.Refund` hatasının "para hareket etmedi" anlamına geldiğinin doğrulanması, ya da serbest bırakmanın buna göre daraltılması.
5. Sağlayıcı kabulünden sonraki `repo.Update` hatası için telafi mekanizması.
6. Onay talebi süre dolumu ve süresi dolmuş talepler için bir zamanlayıcı.
7. Refund kayıtlarında `approval_request_id` saklanması (izlenebilirlik).
8. Rev.06 kabul matrisinin yeniden hesaplanması.
9. Model 1 (DV, DT) ve Model 4 (M, MK, D, F) uygulamaları.
10. Operasyon belgelerindeki eski GPS uyarılarının (`OPERATIONS_GUIDE.md` 3.15, `KULLANIM_KILAVUZU.md` 20.8) düzeltilmesi.
11. Rev.06 analiz belgesi §12.7 ve ilgili maddelerinin güncellenmesi (Bölüm 14.5).
12. Gürsel GPS regresyon testinin bir MongoDB veritabanına karşı çalıştırılması.
13. Flutter uygulamalarının derlenmesi ve doğrulanması.
14. Üretim sunucusuna erişildiğinde canlı durum kontrolü.
15. Faz 0 0.6 kalan açıkları: `ledger_accounts` unique indeksi (veri denetiminden sonra), soft-delete sonrası yeniden açma ve eski+yeni binary riski, canlı replica set doğrulaması, `ConfirmBankTransfer` yeniden denemesi (Bölüm 15.2).
16. İade serbest bırakma kararının verilmesi ve PENDING_EXTERNAL sorgulamasının tasarlanması (Bölüm 14.3, risk 4).
17. Plaka geçmişi yazım hatalarının işlenmesi: `Create` (service.go 75) ve `UpdateVehicle` (service.go 128-129) (Bölüm 14.3, risk 25). `AddPlate` için geçmiş kapatma, geçmiş ekleme ve araç güncellemesinin tek işlemde yapılıp yapılamayacağının incelenmesi.
18. `stock.vehicles.validation.*` ve `stock.vehicles.error.not_found` anahtarlarının seed çevirilerine eklenmesi (Bölüm 14.3, risk 26).
19. `internal/platform/middleware/module_gate.go` dosyasının biçimlendirilmesi (Bölüm 14.3, risk 27).
20. Eski `/api/admin/b2b` yolunun 200 yanıtının kaynağının bulunması (Bölüm 14.3, risk 1).
