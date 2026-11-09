# AYSTEK Engage — Ar-Ge Raporu ve İş Planı (Çalışma Adı)

**Durum:** TASLAK v0.1 — onay bekliyor
**Tarih:** 2026-10-07
**Sahip:** Hüseyin Koç
**Doküman tipi:** Yaşayan iş planı (README_FOR_AI_ASSISTANTS.md → "Active/Planned Work Trackers" listesinde izlenir)
**Okuyucu:** Proje sahibi + bu işi sürdürecek mühendisler ve AI araçları. Kod yazmadan önce §5 (envanter) ve §12 (temizlik) okunmalıdır.

> "AYSTEK Engage" bir **çalışma adıdır**; ürün adı Faz D'den önce netleşir. Kod, menü anahtarı ve çeviri anahtarlarında ürün adı geçmez; modül adları fonksiyoneldir (`audience`, `messaging`, `campaigns`, `automation`, `partners` …).

---

## 0. Yönetici Özeti

**İstek:** SendPulse benzeri, çok kanallı (e-posta, SMS, WhatsApp/Telegram/… chatbot, web push, CRM, otomasyon) hizmet veren bir **public uygulama** sunmak; bunu AYSTEK Core üzerinde, **ortakların (partner) kendi markalarıyla kendi public uygulamalarını açabileceği** biçimde kurmak; ve rakipten daha iyi olmak.

**Ana bulgu:** AYSTEK Core bu işin %30-40'lık altyapısına zaten sahip (tenant/site izolasyonu, olay veri yolu, RabbitMQ işleyicileri, entegrasyon hook'ları, OAuth istemcileri + API ürünleri + kota, SMS outbox/politika/onay, abonelik/fatura/ödeme, defter (ledger), eklenti ve modül kapısı, çok dilli seed/menü sistemi). **Eksik olan çekirdek:** (1) pazarlama kitlesi (audience) + rıza defteri, (2) kanal soyutlaması ve gerçek kanal adaptörleri, (3) durum tutan otomasyon motoru, (4) partner/reseller katmanı, (5) plan + kullanım ölçümü + ön ödemeli kredi, (6) müşteriye dönük ürün-kalitesinde uygulama.

**Fark yaratma önerisi (özet, ayrıntı §4):**
1. **Olay-yerel otomasyon:** Her AYSTEK modülü (sipariş, aidat, araç/filo, okul servisi, forum…) otomasyon tetikleyicisi olur. SendPulse bunu yapamaz; biz dikey çözümleriyle birlikte satarız.
2. **Türkiye-önce uyumluluk:** İYS/ETK + KVKK + GDPR rıza defteri, çıkış (unsubscribe) ve baskılama (suppression) ürünün çekirdeğinde.
3. **Kendi sunucunda kur:** `aystek-core-installer` ile aynı ürün SaaS ya da müşterinin sunucusunda (veri yerleşimi) çalışır.
4. **Partner = birinci sınıf vatandaş:** Sadece yönlendirme linki (SendPulse'ın %25/%10 affiliate'i) değil; white-label alt hesap + kendi fiyatı + şablon/bağlayıcı pazar yeri payı + sertifikasyon.
5. **Yerel/BYO AI ajanı:** Kendi anahtarını getir ya da yerel model (Ollama emsali var); maliyet tavanı ve denetim kaydı zorunlu.
6. **Şeffaf fiyat:** Abonelik + ön ödemeli kredi (PAYG) + partner toptan fiyatı; hepsi admin panelinden yönetilen plan kataloğu (hard-code yok).

**Kritik dış bağımlılıklar (hemen başlatılmalı, kodla paralel):** Meta WhatsApp Tech Provider/BSP başvurusu; e-posta gönderim altyapısı kararı (kendi MTA yazılmaz); İYS firma kaydı/entegrasyon koşulları; SMS operatör sözleşmesi.

**Önerilen sıra (sizin verdiğiniz sıra korunarak):**
`Faz A` Core altyapısı → `Faz B` Figma şablonları → `Faz C` public uygulama entegrasyonu → `Faz D` ürün (product) modu.

---

## 1. Hedef ve Kapsam

### 1.1 Üç katman

| Katman | Ne | Kim kullanır |
|---|---|---|
| **Engage motoru** (Core modülleri) | Kitle, kanal, kampanya, otomasyon, sohbet kutusu, form/pop-up, CRM, ölçüm | Her şey bunun üstünde |
| **Partner programı** | Affiliate, ajans/reseller (white-label), geliştirici, sertifikalı uzman | AYSTEK dışı ortaklar |
| **Public uygulama** (müşteri uygulaması) | Kayıt, fiyat, panel, fatura; partner markasıyla yayınlanabilir | Son müşteri |

### 1.2 Kapsam dışı (şimdilik, bilinçli)
- Çevrimiçi kurs/LMS (SendPulse "Öğretin ve Kazanın"): `library`/`pages` ile karışmasın; Faz D sonrası değerlendirilir.
- TikTok chatbot, Viber ticari mesaj: kanal adaptör arayüzü hazır olacak, adaptör sonra.
- Kendi MTA/SMTP sunucusu yazmak: kapsam dışı (bkz. §9.2).

---

## 2. Ekran Görüntülerinden Çıkarılan Gereksinimler

| # | Gözlem (SendPulse) | Çıkan gereksinim | AYSTEK karşılığı |
|---|---|---|---|
| 1 | Chatbot akış editörü: Keywords (anahtar kelime tetikleyici, Türkçe eş anlamlılar), Welcome message, Standard reply, WhatsApp bölümü (Catalog order), Start/Filter/AI Agent/Flow düğümleri, düğüm başına kişi sayacı, "Edit the trigger", yapışkan not | Görsel akış editörü + tetikleyici türleri + akış içi koşul/AI/alt-akış düğümleri + düğüm bazlı istatistik | Yeni `automation`; editör için `integrations/chains/ChainGraphEditor.tsx` temel alınır (§5) |
| 2 | Ürün menüsü: Gönder (e-posta, SMS, Viber, Web Push, işlemsel e-posta), Etkileşim (açılış sayfası, akıllı pop-up), Chatbot (WhatsApp, Telegram, Instagram, Facebook, TikTok, Viber, e-posta botu, canlı sohbet), Yönet (CRM, E-ticaret CRM, sosyal CRM, omnichannel gelen kutusu, Otomasyon 360), Öğret, Entegre (AI asistanı, API, **MCP sunucusu**, entegrasyonlar, uzman bul) | Ürün haritası | §6 modül haritası |
| 3 | Ortaklar menüsü: Acenteler, Ortaklık Programı, Geliştiriciler, Ortak Kataloğu, Sertifikasyon | 4 ortak türü + katalog + sertifika | `partners` (§7) |
| 4 | Fiyatlandırma: ücretsiz katman (15.000 e-posta/ay), Abonelik vs Kullandığın Kadar Öde, liste büyüklüğü kaydırıcısı, yıllık %20 / 6 ay %10 indirim | Plan kataloğu + kullanım ölçümü + dönem indirimi | `plans` + `usage` (§8) |
| 5 | Affiliate paneli: ref linki (`/#<id>`), 30 günlük tık/gösterim/komisyon/iade, Promotion & Reports menüsü, plan satışında %25, bakiye yüklemede %10, yönetici kişi | Atıf (attribution) + komisyon defteri + iade geri alma + raporlar | `partners` + `ledger` (§7) |

> Not: Fiyat/komisyon oranları kullanıcının kendi ekran görüntülerinden alınmıştır; web araştırması ikincil kaynaklardır (§15). Rakip fiyatları yayın öncesi resmî sayfadan yeniden doğrulanmalıdır.

---

## 3. Pazar Araştırması (Özet)

| Oyuncu | Güçlü yan | Zayıf/boşluk | Bize dersi |
|---|---|---|---|
| **SendPulse** | Çok kanal tek hesapta, cömert ücretsiz katman, PAYG, Automation 360, CRM, açılış sayfası, kurs, API+MCP | Dikey veri yok; partner = çoğunlukla affiliate/ajans yönlendirme; SaaS-only | Ürün genişliği + ücretsiz katman şart; partner modelini genişletebiliriz |
| **GoHighLevel** | White-label "SaaS Mode": ajans platformu kendi markasıyla, kendi fiyatıyla yeniden satar | Karmaşık, ajans odaklı; yerel uyumluluk zayıf | Partner'ın **kendi fiyatını koyması** ana çekim; bunu birinci günden tasarla |
| **Brevo** | E-posta+SMS+CRM+işlemsel, CDP yönü, Avrupa uyumu | Chatbot/akış derinliği sınırlı | İşlemsel e-posta ayrı ürün gibi sunulmalı |
| **ManyChat** | Sosyal DM otomasyonunda UX lideri | Tek kanal odaklı, e-posta/CRM zayıf | Akış editörü UX çıtası |
| **WhatsApp Business Platform (Meta)** | Resmî kanal | 2025'ten itibaren şablon mesaj başına ücret (pazarlama/yardımcı/kimlik doğrulama); 24 saatlik hizmet penceresinde ücretsiz; BSP/Tech Provider statüsü gerekir | Maliyet **pass-through + marj** olarak ölçümlenmeli (§8) |

**Sonuç:** Rekabet "özellik sayısı" üzerinden kazanılmaz. Kazandıran: (a) dikey olay tetikleyicileri, (b) yerel uyumluluk, (c) partnerin kendi markası/fiyatı, (d) kurulum esnekliği.

---

## 4. Fark Yaratma Stratejisi

1. **Olay-yerel otomasyon.** `internal/platform/events` yayıncıları (orders, PMS aidat, fleet_tally, forum, vb.) `automation` tetikleyici kataloğuna otomatik kaydolur. Örnek akış: "Aidat vadesi 3 gün kaldı → WhatsApp şablonu → 24 saatte ödeme yoksa SMS". Tetikleyici kataloğu kodda değil, **olay kayıt defterinde** tutulur (README: dinamik davranış).
2. **Rıza defteri (consent ledger).** Kanal bazlı, kaynaklı, zaman damgalı, silinmez (append-only). Gönderim öncesi her zaman baskılama kontrolü; kampanya motoru rıza yoksa mesajı **kuyruğa bile almaz**. İYS senkron adaptörü (§9.1).
3. **Kendi sunucunda kurulum.** `aystek-core-installer` ile aynı ürün; tenant `data_strategy: shared|dedicated` (Canonical Architecture §6.2) ile uyumlu. Rakiplerde yok.
4. **Partner katmanı:** dört tür (§7), partnerin kendi plan fiyatı ve marjı, otomatik komisyon/ödeme defteri, şablon/bağlayıcı pazar yeri geliri payı.
5. **AI ajanı düğümü:** sağlayıcı soyutlaması (BYO anahtar / yerel Ollama / bulut), **çağrı başına ve günlük maliyet tavanı**, istem + yanıt denetim kaydı, insana devir (handoff). Not: mevcut `ai.Chat` kural tabanlıdır (README'ye göre `AystekAIEngine`), LLM değildir — yeni bir `llm` sağlayıcı soyutlaması gerekir ve ikisi karıştırılmamalıdır.
6. **MCP sunucusu + açık API:** SendPulse'ta var; bizde `integrations/api_sharing` (OAuth istemci, API ürünü, kota, kullanım, OpenAPI) üzerine ince bir `mcp` katmanı.
7. **Şeffaf, yönetilebilir fiyat.** Plan kataloğu DB'de; ücretsiz katman, dönem indirimleri, kanal bazlı kota hepsi admin ekranından.

---

## 5. Mevcut Altyapı Envanteri (Mükerrer İş Önleme)

Kod taramasıyla doğrulanmıştır (2026-10-07, `internal/modules/*`).

### 5.1 Yeniden kullan / genişlet

| Mevcut | Gerçek durum (doğrulandı) | Karar |
|---|---|---|
| `platform/events`, `platform/queue` (RabbitMQ), `platform/jobs`, `platform/idempotency`, `platform/ratelimit`, `platform/circuitbreaker` | Olay veri yolu + kuyruk + zamanlayıcı mevcut | **Aynen kullan**; otomasyonun bekleme/zamanlayıcı ihtiyacı için `jobs/scheduler` incelenecek |
| `integrations` (connector/action/hook/execution/DLQ, eventbridge, `chains` zincir grafiği, scheduled pull) | 89 kaynak/27 test dosyası; olay→hook→kuyruk→işçi hattı | **Otomasyonun "bağlayıcı eylemi" düğümü bunu çağırır.** Hook = durumsuz tek adım; `automation` = durumlu çok adım. **Birleştirme yok.** |
| `integrations/api_sharing` | OAuth istemci/jeton, API ürünü, kota, kullanım, denetim, `PartnerDPARecord`, OpenAPI | **Geliştirici portalı + MCP temeli** |
| `aystek-core-next .../integrations/chains/ChainGraphEditor.tsx` | 946 satır, özel grafik editörü: düğüm türleri start/action/transform/branch/end, geri al/yinele, ızgara, yakınlaştırma. `reactflow` bağımlılığı yok | **Kopyalama!** Ortak bir `components/graph-editor/` paketine çıkar; `chains` ve `automation` ikisi de onu kullanır (aksi halde iki editör = kusur) |
| `sms` (provider registry, outbox, policy/failover, quota, DLR webhook, şablon, toplu kampanya + onay, segment preset) | 13 kaynak dosya; **registry'de yalnızca `mock` sağlayıcı kayıtlı** (`service.go:47`); `RegisterProvider` uzantı noktası var | **SMS kanal adaptörü olarak sarmala.** Gerçek operatör adaptörleri henüz **yok** — Faz A'da yazılacak |
| `emailaccounts` + `platform/mail` + `mailtemplates` | Tenant bazlı SMTP hesabı (şifreli), şablon | E-posta kanalının **BYO SMTP** ayağı; gönderim domaini doğrulama/bounce eksik |
| `notifications` + `notifications/push` | 33 kaynak dosya, push tokenları | **Web/mobil push kanalı** |
| `subscriptions`, `subscription_contracts`, `billing/invoices`, `payments/gateways`, `licenses` | Dönem, sözleşme, fatura, ödeme geçidi, lisans havuzu | **Abonelik faturalaması bunlardan geçer**; yeniden yazma |
| `ledger` | Çift girişli defter | **Kredi bakiyesi + komisyon defteri** için tek doğruluk kaynağı |
| `plugins` + `platform/middleware/module_gate.go` | Modül kapısı `site.<id>.module.<m>.enabled` → tenant → global; eklenti özellik bayrakları | **Plan → modül/özellik hakkı** bu kapıya bağlanır |
| `view_profiles`, `menu_seed`, `permission_seed` (`messaging` menü ebeveyni zaten var) | RULE #2.2 | Yeni menü anahtarları `messaging`/`content-marketing` altına |
| `ai` (job/insight/recommendation/event) | Kural tabanlı | LLM için **yeni** soyutlama; `ai` içinde değil (§4.5) |
| `chatbot` (Dialogflow, intent schema, canned responses) | 17 kaynak dosya, **test yok** | **NLU sağlayıcısı + hazır yanıt** olarak otomasyonun "niyet" düğümüne bağlanır |
| `crm` (HubSpot/Salesforce senkron) | 9 dosya, **test yok**; yalnızca iki sağlayıcı | **Harici CRM konnektörü** olarak kalır; yerel CRM boşluğu ayrı (§6, `deals`) |
| `telematics/partner_onboarding` | Aşamalı (başvuru→UAT→canlı) partner kaydı; **yalnızca bellek-içi repo** | Telematik'e özgü; **aşama makinesi desenini** emsal al, tabloyu paylaşma |
| `pages`, `sliders`, `surveys` | İçerik/CMS | Açılış sayfası ve pop-up için temel; form/pop-up için yeni `forms` |
| `installer` | Kurulum, `public-site` profili, Ollama tavsiyecisi | Public uygulama dağıtımı + kendi-sunucu seçeneği |

### 5.2 İsim çakışması uyarısı
`contacts` modülü **iletişim formu gelen kutusu**dur (`Subject/Message/Status new|read|replied`); CRM kişisi değildir. Yeni kişi modeli için bu isim **kullanılmayacak** → `audience`.

### 5.3 Emekliye ayrılacak (ikili uygulama yasağı, §12)
`newsletter_subscribers`, `marketing` (Mailchimp sarmalı), `sms` toplu kampanya yüzeyi → ilgili yeni modül kendi yerini aldığında silinir.

---

## 6. Hedef Mimari

### 6.1 Modül haritası (DDD: model → dto → repository → service → handler → routes + test)

| Modül (yeni) | Sorumluluk | Temel varlıklar |
|---|---|---|
| `audience` | Pazarlama kitlesi; tek profil, çok kanal kimliği | `Profile`, `ChannelIdentity` (email/phone/telegram_id/wa_id/push_token), `CustomField`, `Tag`, `List`, `Segment` (kural ağacı), `ConsentRecord` (append-only), `SuppressionEntry` |
| `messaging` | Kanal soyutlaması + birleşik gönderim hattı | `Channel` kaydı, `ChannelAdapter` arayüzü, `MessageLog`, `Outbox` (mevcut `sms.OutboxJob` deseni genelleştirilir), kullanım olayı yayını |
| `campaigns` | Tek seferlik/zamanlı toplu gönderim, A/B, onay | `Campaign`, `CampaignVariant`, `CampaignStat` |
| `automation` | Durumlu akış motoru (e-posta serisi, chatbot, olay tetikli) | `Flow`, `FlowVersion`, `Trigger`, `NodeSpec`, `Run`, `RunStep`, `WaitToken` |
| `conversations` | Omnichannel gelen kutusu, insana devir | `Thread`, `Message`, `Assignment`, `Handoff` |
| `forms` | Abonelik formu, pop-up, açılış sayfası bağlama | `Form`, `Popup`, `Submission` (→ `audience`) |
| `deals` (Faz A+) | Hafif yerel CRM (pipeline/fırsat) | `Pipeline`, `Stage`, `Deal` |
| `llm` | AI sağlayıcı soyutlaması (BYO/yerel/bulut), maliyet tavanı, denetim | `LLMProvider`, `LLMUsage`, `PromptAudit` |
| `partners` | Ortak hesabı, atıf, komisyon, white-label, katalog, sertifika | §7 |
| `plans` | Plan kataloğu + hak/kota tanımı (admin yönetimli) | `Plan`, `Entitlement`, `PriceBook` |
| `usage` | Ölçüm → fiyatlandırma → kredi düşümü | `UsageEvent`, `Rate`, `CreditAccount` (ledger destekli) |
| `mcp` (ince) | MCP sunucusu; `api_sharing` üstünde | tool kataloğu |

**Bağımlılık yönü (döngü yok):** `forms/campaigns/automation/conversations → messaging, audience → platform`. `usage` hiçbirine bağımlı değildir; **olay dinler**. `partners → plans/usage/ledger`.

### 6.2 Kritik arayüzler (OCP: yeni kanal/düğüm = yeni dosya, mevcut kod değişmez)

```go
// messaging
type ChannelAdapter interface {
    Channel() ChannelKind                       // email, sms, whatsapp, telegram, push, instagram…
    Capabilities() Capabilities                // şablon zorunluluğu, pencere kuralı, maks uzunluk, medya
    ValidateConfig(map[string]any) error
    Send(ctx context.Context, req SendRequest) (SendResult, error)   // ctx: tenant+site taşır
}
type InboundAdapter interface {                // webhook → conversations
    VerifySignature(headers map[string]string, body []byte, cfg map[string]any) error
    ParseInbound(body []byte) ([]InboundEvent, error)
}

// automation
type NodeExecutor interface {
    Type() NodeType                            // send_message, wait, condition, ai_agent, call_action, set_field, handoff, sub_flow…
    Schema() NodeSchema                        // editör formu bu şemadan üretilir (UI'da hard-code yok)
    Execute(ctx context.Context, run *Run, step *RunStep) (Outcome, error)
}
```
Mevcut `sms/providers.Provider` zaten benzer bir sözleşme sunar (`Capabilities`, `FieldSpecs`, `Send`, `DLRProvider`); `messaging.ChannelAdapter` bunu **genelleştirir**, SMS adaptörü mevcut registry'yi sarar — SMS yeniden yazılmaz.

### 6.3 Akış motoru tasarım ilkeleri
- **Durumlu, kalıcı:** Her `Run` Mongo'da; bekleme düğümü `WaitToken` + zamanlayıcı (RabbitMQ gecikmeli kuyruk ya da `platform/jobs`); işçi çökse devam edebilir.
- **Idempotent adım:** `run_id + step_id + attempt` anahtarı (`platform/idempotency`).
- **Sürümleme:** Yayınlanmış akış değişmez (`FlowVersion`); çalışan Run kendi sürümünde biter.
- **İzolasyon:** Her sorgu `tenantfilter.WithTenantAndSite`; arka plan işi `contextkeys.ForTenantSite`/`Detached` (README §Context). `context.Background()` yasak.
- **Kaotik yük koruması:** Kanal bazlı hız limiti + `circuitbreaker`; tenant başına eşzamanlı Run tavanı (plan hakkı).
- **Gözlemlenebilirlik:** Düğüm başı sayaç (ekran görüntüsündeki "234") `RunStep` üzerinden toplulaştırılır; Prometheus metrikleri.

### 6.4 Olay anahtarları (README: EventKeysInfo zorunlu)
`audience.profile.created|updated|merged`, `audience.consent.granted|revoked`, `audience.suppression.added`, `messaging.message.queued|sent|delivered|failed|bounced|complained|opened|clicked|replied`, `campaigns.campaign.scheduled|started|completed`, `automation.run.started|completed|failed`, `conversations.thread.opened|assigned|handoff_requested`, `partners.referral.attributed`, `partners.commission.accrued|reversed|paid`, `usage.credit.low|exhausted`. Her biri ilgili liste sayfasında `EventKeysInfo` ile belgelenir ve `HookForm.tsx` `KNOWN_EVENT_KEYS`'e eklenir.

---

## 7. Partner Modeli

### 7.1 Hiyerarşi
```
Platform (AYSTEK)  →  Partner (reseller/ajans)  →  Tenant (müşteri)  →  Site (domain + amaç)
```
Canonical Architecture'daki tenant/site izolasyonu **değişmez**. Eklenen tek şey, tenant'ın bir `partner_id` ile ilişkilendirilmesi ve **ortak yönetici (partner admin)** kullanıcı türüdür: yalnızca kendi alt tenant'larında sınırlı yetkisi vardır, platform kullanıcıları gibi RBAC'ı bypass **etmez**. Bu, sızıntı riski en yüksek karardır → **Faz A0'da ADR + güvenlik incelemesi şart** (§13).

### 7.2 Ortak türleri

| Tür | SendPulse karşılığı | Bizde ek değer |
|---|---|---|
| **Affiliate** | Affiliate paneli (ref linki, %25/%10) | Atıf süresi/kural yönetimli; iade geri alma; sahtecilik koruması (kendi kendine yönlendirme engeli) |
| **Ajans/Reseller** | Acenteler | **White-label alt hesap, kendi fiyat defteri (toptan + marj), kendi alan adı/marka**, toplu müşteri yönetimi |
| **Geliştirici** | Geliştiriciler | Şablon / bağlayıcı / akış paketi yayını + gelir payı (`plugins` + `marketplaces` emsalleri) |
| **Sertifikalı uzman** | Sertifikasyon + Ortak Kataloğu | Katalog listesi, rozet, lead yönlendirme |

### 7.3 Varlıklar
`Partner` (tür, durum, aşama makinesi: başvuru→onay→aktif→askıda), `PartnerTier` (admin yönetimli oran/eşikler), `ReferralLink`/`Attribution` (çerez + kayıt anı kaydı), `CommissionRule`, `CommissionEntry` (**ledger'a yazılır**, `accrued→payable→paid|reversed`), `Payout`, `PartnerBranding` (logo, renk, alan adı, e-posta gönderen), `CatalogListing`, `Certification`.

### 7.4 Kurallar
- Komisyon oranları, atıf süresi (30 gün/ömür boyu/ilk 12 ay), asgari ödeme **koda gömülmez**; `PartnerTier` + `CommissionRule` ile admin'den yönetilir.
- İade/ödeme başarısızlığı → `commission.reversed` (defterde ters kayıt; silme yok).
- Ödeme manuel/dış süreç (Wallet emsali: otomatik banka API'si yok) — Faz D'ye kadar.
- Partner'ın müşterisi **son müşterisidir**; müşteri verisine partner erişimi yalnızca DPA kaydı (`PartnerDPARecord` emsali) ve açık yetkiyle.

---

## 8. Plan, Kullanım ve Fiyatlandırma

**İlke:** Hiçbir fiyat/kota koda gömülmez.
- `Plan` → `Entitlement[]` (modül kapısı anahtarları + sayısal kotalar: aylık e-posta, SMS kredisi, kişi sayısı, eşzamanlı akış, kanal sayısı, koltuk).
- `PriceBook` (para birimi + dönem: aylık/6 ay/yıllık; dönem indirimi yüzdesi **veri**). Partner için ayrı `PriceBook` (toptan) + partnerin müşteriye uyguladığı kendi fiyat defteri.
- **Ücretsiz katman** = fiyatı 0 bir plan; kötüye kullanım koruması (e-posta doğrulama, gönderim limitleri, ilk kampanya incelemesi).
- **Kullandığın kadar öde:** `CreditAccount` = `ledger` hesabı. `usage` servisi `messaging.message.sent` olaylarını dinler → `Rate` (kanal × ülke × kategori) ile fiyatlar → kredi düşer. WhatsApp için **Meta ücreti pass-through + marj** ayrı kalem.
- Abonelik faturalaması mevcut `subscriptions/billing/payments` hattından; kredi yükleme `payments/gateways` üzerinden; komisyon yükleme tutarı üzerinden de hesaplanır (SendPulse'taki %10 emsali).
- Düşük bakiye/aşım: `usage.credit.low|exhausted` olayı → bildirim + akış duraklatma (veri kaybı yok).

---

## 9. Uyumluluk ve Dış Bağımlılıklar

### 9.1 Rıza / mevzuat
- **Türkiye:** 6563 sayılı Kanun kapsamında ticari elektronik ileti (SMS, e-posta, arama) için önceden açık rıza, her iletide kolay ret yolu ve **İYS** kaydı/senkronu gerekir; ihlalde idari para cezaları vardır. KVKK ayrı olarak kişisel veri işlemeyi düzenler. *(İkincil kaynaklardan derlendi; İYS entegrasyon (API) koşulları ve tüzel kişi başvuru adımları resmî kaynaktan **Faz A0'da doğrulanacak**.)*
- **AB (ör. Balcı/Fransa sitesi):** GDPR — çift onay (double opt-in) seçeneği, silme/erişim hakkı, DPA.
- **Tasarım sonucu:** `ConsentRecord` append-only; kaynak (form/URL/IP/UA), kanal, amaç, zaman, kanıt; `SuppressionEntry` kampanya ve otomasyonda **zorunlu ön süzgeç**; tüm e-postalarda `List-Unsubscribe` başlığı + çıkış bağlantısı; `newsletter_subscribers.UnsubscribeToken` desenini devral.

### 9.2 E-posta teslim edilebilirliği (en büyük gizli maliyet)
- Kendi MTA yazılmaz. **Gönderim sağlayıcı adaptörü** (ör. bir transactional ESP/relay) + BYO SMTP (mevcut `emailaccounts`).
- Zorunlu: gönderen alan adı doğrulama (SPF/DKIM/DMARC kontrol ekranı), bounce/şikâyet webhook'u → otomatik baskılama, ısınma (warm-up) kotası, tenant başına itibar izleme, spam şikâyet eşiği aşılınca otomatik askı.
- Karar §14-K3.

### 9.3 WhatsApp / Meta
- Resmî yol: Meta **WhatsApp Business Platform (Cloud API)** + **Tech Provider/BSP** statüsü; gömülü kayıt (Embedded Signup) ile müşterinin kendi WABA'sı bağlanır. *(Kaynaklarda ISV ve BSP'nin gömülü kayıtı kendi sitesinde kullanabilme koşulları farklı anlatılıyor → başvuru sürecinde Meta'dan **doğrulanacak**.)*
- Faz A'da **BYO** (müşteri kendi Cloud API belirtecini girer) ile başla; Embedded Signup statü onayı sonrası.
- Şablon onay yaşam döngüsü + 24 saat hizmet penceresi kuralı `Capabilities`'te modellenir; pencere dışında serbest metin **gönderilemez** (motor engeller).
- Instagram/Messenger: Meta uygulama incelemesi gerekir → Faz A sonrası.
- Telegram Bot API: ücretsiz/kolay → ilk sohbet kanalı olarak önerilir.

### 9.4 SMS
Gerçek operatör/yerel sağlayıcı adaptörleri **yok** (yalnızca `mock`). Faz A'da en az bir yerel sağlayıcı + bir uluslararası sağlayıcı adaptörü; sözleşme/sender-ID (başlık) onay süreci iş tarafında paralel başlar.

---

## 10. Faz Planı

Büyüklük: **S** ≈ 1 küçük artım, **M** ≈ birkaç artım, **L** ≈ uzun seri. Takvim, ekip hızı bilinmediği için verilmemiştir; her artım README "Before Completing" kapısından geçer.

### FAZ A — AYSTEK Core altyapısı (şu an)

| Adım | İçerik | Büyüklük | Çıkış ölçütü |
|---|---|---|---|
| **A0** | **ADR'ler**: partner hiyerarşisi + yetki modeli, kitle modeli (`audience` ↔ `persons` sınırı), kanal/düğüm arayüzleri, ölçüm/kredi modeli, e-posta sağlayıcı kararı. **İYS/Meta/operatör dış başvurularını başlat.** Tenant/Site model kaynağı doğrulanır (`stores` modeli bulundu; "Site" varlığının konumu netleştirilecek) | S | ADR'ler onaylı; dış başvuru takip tablosu açık |
| **A1** | `audience` (profil, alan, etiket, liste, segment, rıza, baskılama) + `forms` (abonelik formu) + `newsletter_subscribers` **göçü ve silinmesi** | M | Eski aboneler yeni modele taşındı; eski modül/menü/seed/çeviri silindi |
| **A2** | `messaging` + e-posta kanalı (BYO SMTP + 1 sağlayıcı adaptörü, bounce/şikâyet, domain doğrulama) + SMS adaptörü sarmalayıcı + **gerçek SMS adaptörleri** + push adaptörü | L | Gerçek bir kanal üzerinden uçtan uca gönderim + DLR/bounce geri dönüşü; kullanım olayı yayınlanıyor |
| **A3** | `campaigns` (tek seferlik/zamanlı, A/B, onay) — `sms` toplu kampanya ve `marketing` Mailchimp yüzeyinin **yerini alır** | M | Kampanya, rıza süzgeciyle gönderiyor; eski toplu kampanya yüzeyi kaldırıldı |
| **A4** | `automation` motoru + editör (ortak `graph-editor` çıkarımı) + olay kayıt defterinden tetikleyici kataloğu + düğümler (mesaj, bekle, koşul, alan ata, bağlayıcı eylemi, alt akış) | L | Olay tetikli çok adımlı akış, bekleme dahil, çökme sonrası devam ediyor |
| **A5** | `conversations` + Telegram + WhatsApp (BYO) + `chatbot` (NLU/hazır yanıt) bağlanması + `llm` + AI ajanı düğümü | L | Anahtar kelime → akış → insana devir; AI maliyet tavanı çalışıyor |
| **A6** | `plans` + `usage` + kredi + ücretsiz katman + modül kapısı bağlantısı | M | Plan→hak→kapı zinciri; aşımda akış duraklıyor |
| **A7** | `partners` (affiliate → reseller → katalog) + white-label marka + komisyon defteri + partner paneli (aystek-core-next) | L | Ref linkiyle kayıt→satış→komisyon→iade ters kaydı uçtan uca |
| **A8** | `deals` (hafif CRM), `mcp`, geliştirici portalı cilası | M | Faz A kapanış denetimi |

**Her adım için zorunlu (README):** model+dto+repo+service+test+handler+routes+`main.go` bağlantısı; `translation_seed.go` **28 dil** + `fullyTranslatedPrefixes`; `menu_seed.go` + `allMenuKeys` + görünüm profilleri (RULE #2.2); `permission_seed.go`; `{domain}_seed.go`; `EventKeysInfo`; admin sayfası gradient-hero standardı; `go build`, `go vet ./...` (tüm repo), `go test ./internal/seeds`, `go run` ile sunucu başlatma (route çakışması), `npm run lint && npm run check-i18n`; `done/TASK_N_*.md` + `TASK_N_SUMMARY.md`; installer etkileniyorsa kurulum dokümanları.

### FAZ B — Figma şablonları
Önce **tasarım onayı** (renk/UX sübjektif; kod öncesi Design canvas mockup'ı ile onay — proje deneyimi). Şablon seti: pazarlama/ürün sitesi, fiyat sayfası (kaydırıcı + dönem anahtarı), kayıt/giriş, müşteri paneli (kampanya, akış editörü, gelen kutusu, kitle, raporlar), faturalama/kredi, partner portalı (affiliate paneli, reseller müşteri listesi, komisyon raporları, katalog), white-label tema belirteçleri (design tokens). Çıktı: bileşen envanteri + token seti + sayfa listesi → Faz C iş listesi.

### FAZ C — Public uygulama entegrasyonu
- Müşteri uygulaması **ayrı Next.js uygulaması** (öneri; §14-K1). Çok-kiracılı çalışma zamanı kimliği: `NEWS_PUBLIC_SITE_MULTIDOMAIN_WORKPLAN_TR.md` emsali (tek kod/imaj, dağıtım başına kimlik) ve BFF'in istemci IP'sini Core'a iletmesi şartı.
- Dil listesi `/public/languages/active`'ten (RULE #1.1); sayfa metinleri `pages` modülünden; tema `PartnerBranding`'den.
- Installer'a `engage-app` dağıtım profili + operasyon dokümanları.

### FAZ D — Ürün (product) modu
Ürün adı/marka, fiyat yayını, hukuki metinler (sözleşme, DPA, KVKK aydınlatma), SLA/destek (24/7 yardım merkezi), onboarding, ödeme canlı geçişi, partner programı duyurusu, otomatik komisyon ödemesi, operasyon runbook'ları, canlıya çıkış kontrol listesi. (Kayıtlı hatırlatma: canlıya çıkışta `.env` sızıntısı/FCM anahtarı rotasyonu önce ele alınmalı.)

---

## 11. Frontend (aystek-core-next) İş Başlıkları
`/admin/messaging/*` (kanallar, günlük), `/admin/audience/*` (profil, liste, segment, rıza, baskılama), `/admin/campaigns/*`, `/admin/automation/*` (akış listesi + editör), `/admin/conversations`, `/admin/forms`, `/admin/plans`, `/admin/usage`, `/admin/partners/*`. Hepsi: gradient-hero başlığı, `AdminDataTable`, `NoticeModal`, `rounded-2xl border-slate-100 shadow-sm` (referans: `posts`, `integrations/dashboard`); metinler çeviri anahtarı; `lib/api/*.ts` zarf (envelope) açma kuralları (çift sarma / liste düzleştirme dersleri).

---

## 12. Temizlik ve Emeklilik Planı (README: ikili uygulama = kusur)

| Eski | Yeni sahibi | Silme tetikleyicisi | Silinecekler |
|---|---|---|---|
| `newsletter_subscribers` (model/repo/service/handler/route/menü/seed/çeviri/`/admin/newsletter-subscribers`) | `audience` + `forms` | A1: veri göçü doğrulandı | Tüm modül, menü `newsletter-subscribers`, çeviri anahtarları, public abonelik ucu yeni forma yönlendirilir (kırılma yok: geçiş süresi boyunca eski `/public/...` uç **kabul eder**, sonra kalkar — kaldırma tarihi belgelenir) |
| `marketing` (Mailchimp sarmalı) | `messaging` e-posta kanalı + `integrations` konnektörü | A3 | Modül + rotalar |
| `sms` toplu kampanya (`BulkCampaign`, segment preset, onay, 13 checkpoint dokümanı) | `campaigns` | A3: eşdeğerlik testi geçti | Toplu kampanya/segment preset yüzeyi; `sms` yalnızca kanal adaptörü + politika/outbox olarak kalır. **Geçmiş kampanya verisi için arşiv/göç planı gerekir** (açık karar §14-K5) |
| `integrations/chains/ChainGraphEditor` kopyası | Ortak `graph-editor` | A4 | Editör içi kopya mantığı |
| `ai.Chat` ile "AI ajanı" karışıklığı | `llm` | A5 | Dokümanlarda/arayüzde "AI" ifadesi netleştirilir (README'deki AystekAIEngine uyarısı) |

---

## 13. Riskler

| Risk | Etki | Önlem |
|---|---|---|
| Partner yetki modeli tenant sızıntısı üretir | Kritik | A0 ADR + güvenlik incelemesi; partner admin RBAC bypass etmez; tenantfilter denetim testleri; negatif testler |
| E-posta itibarı yanma (spam, bounce) | Yüksek | Domain doğrulama, ısınma, otomatik askı, ücretsiz katman sınırları |
| Mevzuat (İYS/KVKK/GDPR) ihlali | Yüksek | Rıza defteri + ön süzgeç; hukuk incelemesi; resmî kaynakla doğrulama |
| Meta/operatör onay süreleri | Takvim | Dış başvurular A0'da başlar; BYO ile bağımsız ilerle |
| Kapsam şişmesi (tüm SendPulse'ı klonlama) | Yüksek | §1.2 kapsam dışı listesi; her faz çıkış ölçütü |
| Otomasyon motoru kararlılığı (çift gönderim, kayıp bekleme) | Yüksek | Idempotency, sürümleme, kaos testleri, DLQ |
| AI maliyet patlaması | Orta | Çağrı/gün tavanı, denetim kaydı, plan hakkı |
| Mevcut modüllerde test yok (`chatbot`, `crm`, `marketing`, `contacts`) | Orta | Dokunulan her modül için test borcu aynı işte kapatılır |
| Admin paneli "müşteri ürünü" olmaya zorlanır | Orta | Müşteri uygulaması ayrı (K1) |

---

## 14. Açık Kararlar (Varsayılan önerimle)

Aksi söylenmedikçe **varsayılanla** ilerlenir; itiraz edilen madde ADR'de yeniden açılır.

| # | Karar | Varsayılan öneri | Gerekçe |
|---|---|---|---|
| K1 | Müşteri arayüzü nerede? | Faz A'da admin panelinde mühendislik/yönetim ekranları; **müşteri uygulaması ayrı Next.js** (Faz C) | Canonical: admin panel çatallanmaz/özelleştirilmez; ürün UX çıtası farklı |
| K2 | Tenant–partner ilişkisi | Tenant'a `partner_id` + yeni "partner admin" kullanıcı türü | Mevcut izolasyon korunur |
| K3 | E-posta gönderimi | BYO SMTP + 1 transactional sağlayıcı adaptörü; sağlayıcı seçimi A0'da maliyet/teslim edilebilirlik/veri yerleşimi karşılaştırmasıyla | Kendi MTA yazmak kapsam dışı |
| K4 | İlk sohbet kanalları | Telegram → WhatsApp (BYO) → Instagram/Messenger | Onay yükü sırası |
| K5 | `sms` toplu kampanya geçmişi | Salt-okunur arşiv olarak `campaigns`'a göç, sonra silme | Denetim izi kaybolmaz |
| K6 | Kitle modeli | `audience` bağımsız; `persons` (kurum/Gürsel odaklı) ile **köprü** `person_id` opsiyonel alanı | `persons` site-izoleli ve kurum odaklı |
| K7 | AI sağlayıcıları | Arayüz + BYO OpenAI/Anthropic/Gemini uyumlu uç + yerel Ollama; varsayılan kapalı | Gizlilik + maliyet |
| K8 | Ödeme/komisyon ödemesi | Faz A–C manuel/dış, Faz D'de otomasyon | Wallet emsali + mevzuat |

**Kullanıcıdan gereken (dış):** ürün adı (Faz D'ye kadar çalışma adı), WhatsApp/Meta işletme doğrulaması için tüzel bilgiler, İYS kaydı, SMS operatör sözleşmesi, e-posta sağlayıcı hesabı.

---

## 15. Kaynaklar (ikincil; yayın öncesi birincil kaynaktan doğrulanmalı)
- SendPulse genel bakış / inceleme: https://firebearstudio.com/blog/what-is-sendpulse.html
- SendPulse fiyatlandırma (Capterra, PricingSaaS): https://www.capterra.com/p/146369/SendPulse/pricing/ · https://pricingsaas.com/companies/sendpulse
- SendPulse Automation 360 akış elemanları: https://sendpulse.com/knowledge-base/email-service/automation360/flow-elements
- SendPulse chatbot akış editörü: https://sendpulse.com/knowledge-base/chatbot/flow-builder
- GoHighLevel vs ManyChat / white-label: https://www.ghlexperts.com/comparisons/gohighlevel-vs-manychat · https://suitedash.com/gohighlevel-alternative-white-label.md
- GoHighLevel vs Brevo: https://postmake.io/compare/gohighlevel-vs-brevo
- WhatsApp fiyatlandırma 2026: https://www.authgear.com/post/whatsapp-api-pricing · https://www.unipile.com/whatsapp-business-api-pricing/
- Meta Embedded Signup: https://developers.facebook.com/docs/whatsapp/embedded-signup/
- Meta Tech Provider: https://developers.facebook.com/documentation/business-messaging/whatsapp/solution-providers/get-started-for-tech-providers
- BSP karşılaştırması: https://www.ycloud.com/blog/whatsapp-bsp-providers-official-options-2026
- İYS/ETK: https://netmera.com/iys/ · https://www.ideasoft.com.tr/ileti-yonetim-sistemi-iys-nedir/
- Kullanıcı ekran görüntüleri (SendPulse chatbot, ürün menüsü, ortaklar menüsü, fiyat sayfası, affiliate paneli) — 2026-10-07

---

## 16. İlerleme Takibi (her artımda güncellenir)

| Adım | Durum | Görev dokümanı | Not |
|---|---|---|---|
| A0 ADR + dış başvurular | ⬜ Başlamadı | — | Onay bekliyor |
| A1 audience + forms + newsletter göçü | ⬜ | — | |
| A2 messaging + kanallar | ⬜ | — | |
| A3 campaigns | ⬜ | — | |
| A4 automation | ⬜ | — | |
| A5 conversations + llm | ⬜ | — | |
| A6 plans + usage | ⬜ | — | |
| A7 partners | ⬜ | — | |
| A8 deals + mcp | ⬜ | — | |
| B Figma şablonları | ⬜ | — | |
| C Public uygulama | ⬜ | — | |
| D Ürün modu | ⬜ | — | |
