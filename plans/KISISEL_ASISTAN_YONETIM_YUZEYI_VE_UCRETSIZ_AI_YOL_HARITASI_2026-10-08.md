# Kişisel Asistan → AYSTEK Core Yönetim Yüzeyi + Ücretsiz AYSTEK AI — Analiz ve Yol Haritası

**Tarih:** 2026-10-08
**Durum:** Taslak — kullanıcı onayı bekliyor, henüz uygulamaya geçilmedi.
**İlişkili belgeler:** `KISISEL_ASISTAN_AJANS_AI_HABER_VE_SOHBET_GECMISI_PLAN_2026-06-24.md`,
`KISISEL_ASISTAN_ENTERPRISE_NEWS_AI_INTEGRATION_PLAN_2026-06-24.md` (bu belgedeki
"AYSTEK AI motor yükseltmesi" fazıyla KISMEN örtüşüyor — oradaki Faz 4 sadece haber
*filtreleme* pipeline'ı için harici bir AI çağrısı taslaklıyor, bu belge ise
**sohbet motorunun kendisini** ve **admin-aksiyon yürütmeyi** kapsıyor. İki plan
aynı `AystekAIEngine` soyutlamasına çıkabilir, uygulama sırasında birleştirilmeli.)

## 0. İstek

Kullanıcının kendi ifadesiyle (2026-10-08): *"Kişisel Asistanı AYSTEK Core'un
yönetim panelinden yapılan tüm işlemlerin sık kullanılan özelliklerini hem yapay
zekalara yaptırabileceğimiz hem de bizlerin de yapabileceği konuma
taşıyabilmemiz gerekiyor. Yapay zekada AYSTEK AI'nin geliştirilerek ücretsiz
halde de yönetilebilir olmasını kesinlikle istiyoruz."*

İki ayrı ama birbirine bağımlı hedef:
1. **Aksiyon yüzeyi genişletme** — Core admin panelindeki sık kullanılan
   işlemlerin, mobil uygulamadan hem AI'ya hem doğrudan kullanıcıya açılması.
2. **AYSTEK AI'nin güçlendirilmesi, ücretsiz kalarak** — server-taraflı "AYSTEK
   AI" sekmesinin gerçek bir yapay zeka motoruna dönüşmesi, ama kullanıcıya
   (ya da tenant'a) token/API maliyeti yansıtılmadan.

## 1. Mevcut Durum Analizi (bu oturumda 3 paralel araştırmayla doğrulandı)

### 1.1 Admin paneli — geniş ve zaten tekrarlayan bir "onay" deseni var

`internal/seeds/menu_seed.go` üzerinden çıkarılan tam envanter 14 üst kategoriye
ayrılıyor: dashboard, forum, commerce (ürün/envanter/sipariş/ödeme/fiyatlama),
content (haber/sayfa/yorum/medya/pazarlama), operations (PMS, sözleşmeler,
lisanslar), transport-platform-suite (Gürsel dönüşümü — kurum/kişi/cüzdan/erişim
noktası), integrations, telematics-suite (+ fleet_tally), stock-suite, pms-suite,
settings, security, ai-governance, analytics, system-ops.

Mobil için en doğal adaylar (düşük risk, yüksek "işe yarama" oranı) — 4 tekrarlayan
UI deseni etrafında kümeleniyor:

| Desen | Örnek modüller | Mobilde karşılığı |
|---|---|---|
| Tekli durum-geçişi | Orders, Vehicle Incidents, Service Reservations | Durum rozetli kart + "İlerlet" butonu |
| Gerekçe-zorunlu onay/red | Forum AI Governance Decisions | Onay kartı + zorunlu not alanı (Phase 6'daki `_ActionCard` deseniyle BİREBİR aynı) |
| Toplu onay/red | Content Bots / Moderation Queue | Çoklu-seçim + toplu aksiyon çubuğu |
| Tek-tık yayınlama | Classifieds, Comments | Tek buton, anında durum değişimi |

Bu dört desen, Phase 6'da inşa ettiğimiz `AiAction` / onay-kartı mekanizmasıyla
(`lib/models/ai_action.dart`, `_ActionCard`, `_executeAction`) **doğrudan uyumlu**
— sıfırdan bir UI paradigması gerekmiyor, var olanı genellemek yeterli.

### 1.2 İzin modeli — "admin işlemi" aslında admin'e özel değil

`internal/seeds/role_seed.go` incelendi: `editor` rolü (çoğu tenant kullanıcısının
gerçekte sahip olduğu rol) zaten `posts`, `pages`, `comments`, `media`, `events`,
`classifieds`, `stock_*`, `wallets`, `pms_*`, `institutions`, `receivables` gibi
onlarca modülde **tam CRUD** iznine sahip (`editorWrite` + `editorWriteActions`,
role_seed.go:200-247). Tek kesin, aşılamaz sınır `platformOnlyPermissionModules
= {"system", "system_ops"}` (role_seed.go:20-24) — bunlar AYSTEK platform
personeline özel, tenant admin'i bile dahil değil.

**Sonuç:** "Mobilden admin paneli işlemi yapma" çoğu sıradan kullanıcı için
zaten gerçekçi bir hedef, sadece "admin" rolüne özel bir ayrıcalık değil. İzin
kontrolü mekanizması da zaten hazır: `middleware.PermissionMiddleware` path'ten
modülü çözüp `"<modül>:*"` önekiyle kontrol ediyor (Phase 6'da `/admin/posts`
için doğrulandı) — yeni bir yetki sistemi kurmaya gerek yok, var olanı kullanmak
yeterli.

### 1.3 AYSTEK AI — en kritik bulgu: zaten ücretsiz, çünkü gerçek bir AI değil

`internal/modules/personal_assistant/service/engine.go` incelendi:

- `defaultEngine` (`knowledgeEngine` struct) **hiçbir dış LLM'e bağlı değil.**
  Mod bazlı (`assistant`/`newsDesk`/`coffeeFal`/`tarotFal`/`cooking`/`suggestions`)
  anahtar-kelime eşleştirmesi + sabit Türkçe şablon metinleri (`fmt.Sprintf`)
  üretiyor. Haber modunda bile RAG makalelerini şablona gömüyor, gerçek bir
  özetleme/üretim yapmıyor.
- Arayüzün kendi kod yorumu (engine.go:11-12) bunu zaten öngörmüş:
  *"AystekAIEngine, AYSTEK'in kendi AI yanıt motorudur. İleride Ollama veya
  başka bir yerleşik LLM ile değiştirilebilir."*
- `chat_service.go:188-243`'te gerçek bir LLM çağrısı beklentisiyle yazılmış
  ama **hiçbir yerde çağrılmayan** tam teşekküllü `buildSystemPrompt()` fonksiyonu
  duruyor — yarım bırakılmış, kullanılmayı bekleyen bir temel.
- Go repo'sunun tamamında (`grep -rli "openai\|anthropic\|gemini"`) **sıfır** dış
  LLM entegrasyonu var — bu alan tamamen yeşil saha.
- `internal/modules/chatbot/` (farklı, e-ticaret destek botu) zaten çalışan bir
  intent→action router'a sahip (`service/action_router.go`) — "AI bir backend
  işlemini gerçekten yürütür" için somut bir emsal, ama bizim amacımıza (admin
  CRUD aksiyonları) doğrudan uygulanabilir değil, sadece mimari ilham.
- Flutter tarafında "AYSTEK AI" sekmesi (`ServerChatSessionsScreen`) ile
  "Asistan" sekmesi (`LocalChatSessionsScreen`, `ai_service.dart` — kullanıcının
  KENDİ API anahtarı) **tamamen ayrı kod yolları.** Bu ayrım KORUNMALI (aşağıya
  bakınız).

**Sonuç:** Kullanıcının "ücretsiz ama gerçek AI" isteği, kod tabanının zaten
yöneldiği yönle birebir örtüşüyor. Eksik olan tek şey: `AystekAIEngine`
arayüzünün arkasına gerçek bir model bağlamak — ve bunu **platform ödemeli
tek bir yerde** yapmak, her kullanıcıdan ayrı API anahtarı istemeden.

## 2. Stratejik Karar Noktaları (uygulamaya geçmeden önce netleşmeli)

Bu üç karar, her birinin maliyet/risk/kalite dengesi farklı olduğu için
kullanıcı onayı gerektiriyor — bu belge sadece analiz, henüz seçim yapılmadı.

### Karar 1 — AYSTEK AI'yi nasıl güçlendirelim?

| Seçenek | Maliyet modeli | Artı | Eksi |
|---|---|---|---|
| **A. Self-hosted açık kaynak model (Ollama)** | Sabit altyapı (GPU/sunucu), kullanıcı başına artmaz | Öngörülebilir maliyet, veri Core'dan çıkmaz, kod yorumunun zaten işaret ettiği yön | Sunucu/GPU kurulumu + bakım yükü; kalite küçük modellerde sınırlı olabilir |
| **B. Merkezi ucuz bulut API (tek platform anahtarı, örn. DeepSeek/Groq)** | Değişken, token-bazlı, ama TEK fatura Core'a | Kurulum kolay, kalite yüksek | Kullanım arttıkça maliyet büyür → mutlaka kota/rate-limit şart (SMS modülündeki `QuotaPerDay` deseni doğrudan uyarlanabilir) |
| **C. Hibrit** | A varsayılan + kullanıcı isterse kendi anahtarıyla B/mevcut "Asistan" sekmesine geçsin | İki modelin de güçlü yanını alır, zaten "AYSTEK AI" + "Asistan" ayrımı bunu doğal olarak destekliyor | İki motoru paralel sürdürme yükü |

**Ön değerlendirme:** Kullanıcı "ücretsiz" diye net konuştuğu için **Seçenek C
(hibrit)** en güvenli başlangıç gibi görünüyor — "AYSTEK AI" sekmesi self-hosted/
ücretsiz bir modelle güçlendirilir (öngörülebilir maliyet), "Asistan" sekmesi
zaten kullanıcının kendi (isterse ücretli) anahtarıyla çalışmaya devam eder.
Kesin karar kullanıcıda.

### Karar 2 — Aksiyon kataloğu mimarisi

Phase 6'da öğrendiğimiz ders: **yeni backend yazmaya çoğu zaman gerek yok** —
var olan `/admin/posts` ucu mobil için de yeterliydi. Aynı prensip burada da
geçerli olmalı: her admin modülü için zaten bir `/admin/<modül>/...` ucu var,
mobil/AI bunu DOĞRUDAN (aynı JWT, aynı izin kontrolü) çağırabilir.

Ölçeklenme sorusu: her yeni aksiyon için Flutter'da elle yeni bir `AiActionType`
case'i + sistem promptu metni eklemek (Phase 6'nın yaptığı gibi) **az sayıda
aksiyon için yeterli ama 20-30 modüle büyüyünce sürdürülemez.** Bunun yerine
**bir "aksiyon kataloğu"** (JSON/Dart yapısı: `{name, module, method, path,
requiredFields, confirmationLabel, requiredPermission}`) tanımlanıp hem AI
promptuna hem insan arayüzüne BU TEK kataloğdan üretim yapılması önerilir —
"bir kere tanımla, ikisi de kullansın" prensibi. Bu, ayrı bir mimari tasarım
kararı; Faz 8'de detaylandırılacak.

### Karar 3 — Kapsam sırası

Önce düşük-riskli onay/moderasyon akışları mı (Notifications, Classifieds,
Comments — veri kaybı/yanlış işlem riski düşük), yoksa ticari etkisi olanlar mı
(Orders, Payments)? **Önerilen: düşük riskliden başla** — Phase 6'nın "haber
yayınlama" riskiyle (gerçek içerik oluşturma) zaten bir kez karşılaştık, aynı
dikkat ticari/finansal modüllerde daha da kritik.

## 3. Önerilen Fazlı Yol Haritası

Her faz kendi EnterPlanMode/onay turunu gerektirir — bu belge sadece üst
düzey sıralama ve kapsam.

### Faz 7 — Mobil "Hızlı Aksiyonlar" MVP (insan kullanımı, backend değişikliği yok)
- Flutter'da yeni bir "Hızlı Aksiyonlar" ekranı: 2-3 düşük riskli modülden
  (örn. Notifications mark-as-read, Classifieds publish, Comments approve/spam)
  bekleyen öğeleri listele, var olan admin uçlarını (`/admin/notifications/:id/read`
  vb.) doğrudan çağır.
  Site/tenant bağlamı zaten Phase 1-4'te kurulan `_TenantAwareClient` +
  site-switcher üzerinden gelir.
- Çıktı: kullanıcı değerini hızlı kanıtlayan, düşük riskli bir MVP; "aksiyon
  kataloğu" ihtiyacını gerçek kullanımla doğrulama fırsatı.

### Faz 8 — Aksiyon kataloğu + bu aksiyonları AI'ya açma
- Karar 2'deki kataloğu tasarla; Faz 7'nin 2-3 aksiyonunu bu kalıba taşı.
- AI tarafı: Phase 6'nın `<<<KA:create_news_post>>>` desenini genelleştir —
  katalogdan otomatik üretilen KA-komut dokümantasyonu + tek bir generic
  `_executeAction` dispatcher (modül bazlı switch yerine katalog-driven çağrı).
- Her aksiyon için onay kartı zorunlu (Phase 6'daki "production deploy"
  sınıflandırıcısı dersi: AI'nin gerçek içerik/işlem değiştiren adımları HER
  ZAMAN insan onayından geçmeli).

### Faz 9 — AYSTEK AI motor yükseltmesi (Karar 1'e göre)
- `AystekAIEngine` arayüzünün arkasına gerçek model bağlama (muhtemelen
  self-hosted/Ollama, Karar 1'e göre netleşecek).
- `chat_service.go`'daki kullanılmayan `buildSystemPrompt()` fonksiyonunu
  devreye alma.
- Eğer Seçenek B/C'nin bir parçası olarak merkezi ücretli API kullanılacaksa:
  SMS modülündeki `QuotaPerDay`/Redis sayaç desenini tenant-bazlı AI kullanım
  kotasına uyarlama (role_seed/sms deseninden doğrudan taşınabilir).
- `KISISEL_ASISTAN_ENTERPRISE_NEWS_AI_INTEGRATION_PLAN_2026-06-24.md`'deki
  Faz 4 (haber filtreleme AI'sı) ile bu fazın birleştirilip birleştirilmeyeceğine
  (aynı model/altyapıyı paylaşabilirler) burada karar verilmeli.

### Faz 10 — Kapsam genişletme + gözlemlenebilirlik
- Faz 7-8'in kataloğunu daha fazla modüle (Orders, PMS, Vehicle Incidents vb.)
  genişletme.
- Activity Logs modülüyle entegre "mobilden/AI'dan yapılan işlem" izlenebilirliği
  (zaten var olan `activity_logs` modülü — yeni bir audit sistemi gerekmez,
  sadece mobil/AI kaynaklı çağrıları etiketlemek yeterli).

## 4. Kapsam Dışı (bu analiz turunda bilinçli olarak ele alınmadı)
- `system`/`system_ops` modülleri — platform-only sınır, mobil/AI'ya ASLA
  açılmayacak.
- Telematics-suite / fleet_tally gibi çok-adımlı, donanıma bağlı operasyonel
  akışlar — "sık kullanılan" tanımına girmiyor, ayrı bir değerlendirme gerekir.
- Ödeme/finansal onaylar (Payments, Invoices) — Karar 3 gereği ilk fazlara dahil
  edilmedi, ayrı ve daha sıkı bir risk değerlendirmesi gerektirir.

## 5. Sıradaki Adım
Kullanıcının bu üç karar noktasında (Karar 1/2/3) tercihini belirtmesi
bekleniyor. Onay sonrası Faz 7 için ayrı bir EnterPlanMode oturumuyla somut
dosya/uç listesi çıkarılacak.
