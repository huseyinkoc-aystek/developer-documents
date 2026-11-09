# AYSTEK Mühendislik Public CMS Kurumsal Sitesi İçerik ve Çok Dilli Enterprise Yayın Stratejisi

## Hedefler, kitleler ve public içerik ilkeleri

AYSTEK Mühendislik’in yeni public CMS kurumsal sitesinin (AYSTEK Core ile entegre) ana hedefi, **kurumsal satın alma kararına etki eden paydaşlara** güven veren, kendi kendine anlaşılabilen (self-service) ve uzun satış döngülerini destekleyen bir bilgi merkezi kurmaktır. B2B satın alma süreçlerinde alıcıların dijital “self-service”e yöneldiğini gösteren güncel araştırmalar, sitenin yalnızca vitrin değil aynı zamanda “satış + onboarding + destek” bileşeni gibi kurgulanmasını destekler. citeturn6search2turn6search11

Bu nedenle public sitede içerikler; **(i)** karar verici yönetici (rektorat/IT/finans), **(ii)** operasyon yöneticisi (kampüs işletme/ulaşım/servis), **(iii)** teknik ekip (entegrasyon/siber güvenlik), **(iv)** kullanıcı destek ekipleri için ayrı niyet (intent) kümelerinde tasarlanmalıdır. Bu yaklaşım, HubSpot’un “buyer journey” (farkındalık–değerlendirme–karar) modelinde önerilen içerik eşleşmesiyle de uyumludur. citeturn6search1turn6search17

Public içerik için iki kritik kural önerilir:  
Birincisi, **iddiaları ispatlanabilir şekilde yazmak** (ör. “%X tasarruf” gibi ölçülemeyen rakamlardan kaçınmak, varsa case-study/kanıt ile vermek). İkincisi, **sözleşmesel/teknik hassas detayları (müşteri özel mimari, uç noktalar, iç tablo adları vb.) public’e taşımamak**; bunun yerine “entegrasyon prensipleri, güvenlik yaklaşımı, uyumluluk hedefleri” gibi genellenebilir trust-signal’lar sunmaktır. B2B web sitelerinin uzun ve karmaşık satın alma süreçlerini desteklemesi gerektiğine dair NN/g değerlendirmeleri bu tür “ikna + bilgi” dengesini doğrular. citeturn6search11

## AYSTEK Core ile entegre public CMS için referans mimari

AYSTEK Core ile entegre çalışan bir public CMS’de tipik iki seçenek vardır: **(A)** Core’un içerik servislerini “headless CMS” gibi kullanmak veya **(B)** CMS’i ayrı tutup Core’a yalnızca kimlik/varlık/veri API’leriyle bağlamak. Hangi yol seçilirse seçilsin, enterprise ölçekte sürdürülebilirlik için şu mimari ilkeler önerilir:

Public site performansı ve SEO açısından **Core Web Vitals** (yükleme, etkileşim, görsel stabilite) metrikleri önemlidir. Google, Core Web Vitals’ın gerçek kullanıcı deneyimini ölçen metrikler olduğunu ve iyileştirilmesini özellikle önerir. citeturn2search0 Bu da içerik servis katmanında CDN önbellekleme, görsel optimizasyonu ve sayfa render stratejilerinin (SSR/SSG/ISR) planlı olmasını gerektirir.

Güvenlik tarafında, public CMS + Core entegrasyonunu tasarlarken en azından **OWASP Top 10** risk sınıflarını (kimlik doğrulama/erişim kontrolü, enjeksiyon, XSS vb.) bir kontrol listesi olarak kullanmak iyi pratiktir. citeturn0search2 Public sitede özellikle form kanalları (demo talebi/iletişim) hedef olacağından, tarayıcı tarafı korumalar için **Content Security Policy (CSP)** gibi mekanizmalarla XSS ve içerik enjeksiyonu riskini azaltmak önerilir. citeturn3search2turn3search18

Görsel varlıkların “AYSTEK Core DB’sine eklenmesi” isteği için pratikte iki katmanlı model daha güvenli ve ölçeklenebilir olur: (1) DB’de varlık kaydı/metadata (dil, alt text, lisans, kullanım yeri, checksum), (2) binary dosya depolaması için DB içi BLOB veya DB ile tutarlı dosya sistemi nesne depolama. SQL Server ekosisteminde FILESTREAM yaklaşımı, BLOB verisini dosya sisteminde saklarken ilişkisel verilerle **transactional consistency** sağlayabilen bir model sunar. citeturn1search3 (Core’un hangi veritabanını kullandığı bilinmediği için bu bölüm “örnek referans” olarak değerlendirilmelidir.)

## Çok dilli içerik, URL stratejisi ve SEO-i18n gereksinimleri

Çok dilli kurumsal sitede en sık hata; dil seçimini yalnızca UI’da yapıp arama motoru tarafını eksik bırakmaktır. Google, her dil versiyonu için **ayrı URL** kullanılmasını; dilin cookie/Accept-Language ile dinamik gösterilmesine güvenilmemesini önerir. Çünkü crawler her varyasyonu göremeyebilir; ayrıca otomatik yönlendirmeler kullanıcıların ve arama motorlarının tüm versiyonlara erişimini engelleyebilir. citeturn9view0turn9view1

Bu nedenle önerilen standart:

Site dili için rota yapısı (örnek):  
`/tr/…`, `/en/…`, `/ar/…`, `/fr/…`, `/es/…`, `/ru/…`  
ve her sayfada karşılıklarını ifade eden **hreflang** işaretlemeleri. Google, localized sayfalar için hreflang anotasyonları + sitemap gibi yöntemlerle varyantları açıkça belirtmeyi anlatır; ayrıca hreflang setinin karşılıklı (bidirectional) olmasının önemini vurgular. citeturn9view1

Dil etiketleme standardında UI çevirileri kadar teknik doğruluk da gerekir. HTML’de `lang` attribute’u; ekran okuyucu/yardımcı teknolojiler ve doğru metin işleme için önemlidir ve BCP 47 dil etiketleriyle kullanılmalıdır. citeturn2search10turn1search2turn2search2

Dil seti önerisi (Türkçe hariç en az 5 dil) için, global erişilebilirlik ve kurumsal kullanım açısından **BM’nin 6 resmî dili** iyi bir başlangıç kümesi sağlar: Arapça, İngilizce, Fransızca, Rusça, İspanyolca (ve Çince opsiyonel). BM bu dillerin resmî diller olduğunu açıkça belirtir. citeturn5search2 Bu set, “en çok bilinen” ihtiyacını pratikte karşılar; ayrıca Avrupa müşterilerinde Almanca’nın değerli olması muhtemeldir (opsiyonel 7. dil olarak eklenebilir).

Görsel SEO ve erişilebilirlik için: Google, görsellerde en önemli meta verinin **alt text** olduğunu; bunun erişilebilirliği de iyileştirdiğini belirtir ve keyword-stuffing’den kaçınmayı önerir. citeturn9view2 Performans tarafında, web.dev; responsive görseller için `srcset/sizes` ile farklı boyutlarda görsel sunmayı ve modern formatların (WebP/AVIF vb.) dosya boyutlarını düşürerek yüklemeyi hızlandırabildiğini anlatır. citeturn4search0turn4search9

## Güven, uyumluluk ve kamuya açık metinler

AYSTEK Mühendislik çözümleri; kampüs geçiş, nakitsiz ödeme/top-up, taşımacılık (GPS/çetele), servis ödemeleri, çocuk güvenliği ve araç içi panik alarm gibi alanlarda kişisel veri ve bazen ödeme verisiyle temas edebilir. Bu nedenle public sitede “trust center” mantığında, aşağıdaki metinlerin kurumsal bir dille yayınlanması tavsiye edilir:

KVKK kapsamında aydınlatma yükümlülüğü, veri sorumlusunun ilgili kişiye; veri sorumlusunun kimliği, işleme amacı, aktarım alıcıları ve amaçları, toplama yöntemi ve hukuki sebep ve ilgili kişi hakları gibi başlıkları sağlamasını gerektirir. citeturn7view1turn1search12 Bu, public sitede en azından “Aydınlatma Metni / Gizlilik Politikası / Çerez Aydınlatma Metni / İletişim Formu Aydınlatması” gibi bölümlerin kurumsal dille hazırlanmasını gerektirir.

GDPR tarafında (özellikle AB’deki kullanıcılar hedefleniyorsa veya davranış izleme/çerezlerle AB ziyaretçileri etkileniyorsa), işleme ilkeleri (hukuka uygunluk-şeffaflık, amaç sınırlaması, veri minimizasyonu, saklama sınırlaması vb.) açıklandığından; site metinleri bu ilkelere uyumlu tasarlanmalıdır. citeturn8view2turn8view0turn8view1

Ödeme toplama/top-up gibi akışlar için, PCI DSS; kart/hesap verisini korumaya yönelik teknik ve operasyonel gereksinimler için global bir temel standarttır. PCI SSC, PCI DSS’in ödeme hesap verisini korumayı hedefleyen “baseline” gereksinimler sunduğunu belirtir. citeturn10search5turn0search3 Ayrıca PCI DSS v4.x’de “future-dated” gereksinimlerin 31 Mart 2025 itibarıyla değerlendirmelerin parçası olması gerektiği PCI SSC kaynaklarında vurgulanır. citeturn10search1turn10search10 (Public sitede “PCI compliant” gibi sert bir iddia ancak gerçekten belgelendirilebiliyorsa kullanılmalı; aksi halde “PCI DSS prensipleri dikkate alınarak tasarlanır / ödeme altyapısı PCI DSS uyumlu sağlayıcılarla entegre edilir” gibi daha doğru bir söylem seçilmelidir.)

Türkiye’de ödeme hizmetleri ve elektronik para alanında 6493 sayılı Kanun çerçevesi bulunduğundan; nakitsiz/online yükleme ve ödeme fonksiyonlarında rol dağılımı (AYSTEK mi ödeme kuruluşu, yoksa entegre olunan PSP mi?) public içerikte net ve ihtiyatlı ifade edilmelidir. TCMB’nin 6493 mevzuat sayfası ve kanun metnindeki amaç/kapsam maddeleri, bu alanın regüle olduğunu doğrular. citeturn4search3turn4search7

Çocuk unutulmasını önleyici çözümler ve panik alarm gibi güvenlik fonksiyonlarında, public içerik güven verici ama aşırı vaat içermeyen bir dil kullanmalıdır. Örneğin NHTSA ve CDC, araç içinde çocuk bırakmanın hızla ölümcül olabilecek sıcaklıklara yol açabildiğini vurgular. citeturn11search0turn11search1 Bu bulgular, ürününüzün “risk azaltma” değerini anlatmak için kullanılabilir; ancak “%100 önler” gibi mutlak iddialardan kaçınmak daha doğrudur.

## İçerik mimarisi ve AYSTEK Mühendislik hizmet makaleleri

Public CMS’de, enterprise alıcılar için en hızlı güven inşa eden yapı “çözüm odaklı” bilgi mimarisidir. Bu kapsamda önerilen üst düzey site haritası:

Kurumsal sayfalar: Ana Sayfa, Hakkımızda, Yaklaşım (Kurumsal teslimat modeli), Kariyer, İletişim.  
Platform: AYSTEK Core (mimari bakış, modüller, entegrasyon yaklaşımı).  
Çözümler (AYSTEK Mühendislik neler yapıyor):  
Kartlı/QR geçiş + nakitsiz/online yükleme + indirim ve kurumlar arası iletişim (kampüs/kurum ekosistemi)  
GPS verileriyle personel/öğrenci taşımacılığı için çetele analizi ve doğrulama  
Okul servis taşımacılığında online ödeme, servis kaydı ve ayrılma süreçleri  
Araç içi çocuk unutulmasını önleyici çözümler  
Araç içi panik alarm hizmeti  
Stok yazılımları  
Özel çözümler (entegrasyon ve proje bazlı geliştirme)

Bu sayfaların her biri “makale formatında” yayınlanmalı; sayfa yapısı tekrar kullanılabilir olmalıdır: problem tanımı → çözüm yaklaşımı → öne çıkan yetenekler → işletim/operasyon → entegrasyon → güvenlik/uyumluluk → raporlama/KPI → sık sorulan sorular → demo/teklif CTA.

Çok dilli üretim için pratik operasyon önerisi: Her makaleyi önce Türkçe “kaynak metin” olarak çıkarıp, aynı içerik modelini diğer dillere taşımak; URL’ler ve hreflang setini otomatik üretmek. Google, aynı sayfanın farklı dil varyantlarını açıkça belirtmeyi ve kullanıcıya dil değiştirme imkânı sunmayı önerir; otomatik redirect ile kullanıcıyı zorla başka dile göndermemeyi özellikle vurgular. citeturn9view0turn9view1

Forum/UGC (kullanıcı üretimi içerik) planlanıyorsa, Google dokümantasyonu; bazı sitelerin ana içeriği tek dilde tutup şablon öğelerini çevirerek forum gibi UGC alanlarında bu yaklaşımları kullandığını belirtir. citeturn9view1 Bu, “forum çok dilli olacak mı?” sorusuna iki seçenek sunar: (1) forumu tek dil (TR veya EN) tutmak, (2) kategorilere göre dil bazlı forum alanları açmak. Enterprise pratikte genellikle (1) ile başlar, talep oldukça genişler.

## Blog ve forum içerik üretim sistemi ve yayın takvimi önerisi

Blog, Google’da discoverability ve satın alma sürecinde “farkındalık” yaratmak için en etkili kanallardan biridir; ancak B2B’de blog yazıları tek başına trafik değil, aynı zamanda “self-service buyer journey” karar desteği üretmelidir. Gartner’ın bulgusu, alıcıların önemli bir kısmının temsilcisiz (rep-free) bir deneyime yöneldiğini söyler; bu da sitenin teknik/operasyonel sorulara yanıt veren içerik üretmesini kritik hale getirir. citeturn6search2

Bu çerçevede önerilen blog kümeleri (pillar–cluster):  
Kampüs teknolojileri (nakitsiz kampüs, kart/QR geçiş, indirim/entitlement yönetimi)  
Okul servis ve taşımacılık dijitalleşmesi (ödemeler, kayıt/ayrılma, doğrulama, güvenlik)  
Telemetri ve GPS ile operasyon doğrulama (çetele/rota/gerçekleşen hizmet analizi)  
Güvenlik ve uyumluluk (KVKK, veri minimizasyonu, erişim kontrol, denetim izleri)  
Platform ve entegrasyon (AYSTEK Core entegrasyon prensipleri, API/webhook yaklaşımı)

Forum tarafında, amaç “support deflection + entegrasyon hızlandırma + topluluk güveni” olmalıdır. Güvenlik açısından da forum giriş/üye yönetimi, OWASP riskleri ve temel web uygulama güvenliği ilkeleriyle tasarlanmalıdır. citeturn0search2 Ayrıca forum sayfaları da public olduğundan, performans ve görsel optimizasyon web standardı olmalı: görsellerin boyutlandırılması için `srcset/sizes` ve modern formatlar gibi yöntemler hem hız hem maliyet için önemlidir. citeturn4search0turn4search9

Ayrıca aksesibiliteyi “enterprise kalite” standardına çekmek için WCAG 2.2 hedeflenmelidir. W3C WCAG 2.2’nin güncel sürüm olduğunu ve erişilebilirlik çalışmalarında en güncel sürümün kullanılmasını önerir; 2.2, 2.1’i genişletir. citeturn0search1turn0search13

## Codex veya Copilot için .mdf prompt dosyası

Aşağıdaki dosya; Codex/Copilot’a verildiğinde, örnek şablonunuzdaki public siteye **çok dilli sayfaları, çözüm makalelerini, blog + forum seed içeriklerini** eklemesi ve görsel varlıkları AYSTEK Core DB varlık modeline uygun biçimde “asset kaydı” olarak yerleştirmesi için hazırlanmış “enterprise düzey” bir çalışma promptudur.

Bu prompt; çok dil için BCP 47 etiket mantığına (ör. `tr-TR`, `en-US`) dayanır. citeturn1search2turn2search10 Çok dilli SEO için “her dil için ayrı URL + hreflang + otomatik redirect’ten kaçınma + kullanıcıya dil seçimi” yaklaşımını temel alır. citeturn9view0turn9view1 Görsel SEO/erişilebilirlik için alt text’in önemini ve keyword stuffing’den kaçınmayı ekibe aktarır. citeturn9view2

Aşağıdaki içerik bir “metin prompt dosyasıdır”; uzantı olarak `.mdf` kullanmanız istenmiş olsa da, bu dosya SQL Server’ın `.mdf` veritabanı dosyası değildir (metin dosyasıdır). Eğer gerçekten SQL Server `.mdf` (database) üretimi gerekiyorsa, bunun için şema ve ortam bilgisi şarttır; bu prompt ise “repo/şablona içerik ekleme” içindir.

[Download the .mdf prompt file](sandbox:/mnt/data/AYSTEK_public_cms_prompt.mdf)

```text
AYSTEK_PUBLIC_CMS_ENTERPRISE_PROMPT (v1.0)
Role: You are a senior software engineer + senior content strategist.
Goal: Add enterprise-grade, multilingual public website content for AYSTEK Mühendislik into the existing "public CMS template" project.
Integration: The site works integrated with AYSTEK Core (use existing integration layer / APIs / DB seeds already in the repo).

CORE RULES (DO NOT VIOLATE)
1) Do not invent certifications, compliance badges, customer approvals, SLAs, or metrics unless there is an existing verified source in repo.
2) If customers like "Sabancı, Özyeğin, KOÇ, Işık, Bahçeşehir" are mentioned, place them behind a feature flag or a "references" component that requires explicit approval (comment explaining legal/logo permission).
3) SEO & i18n:
   - Use separate URLs per language (e.g., /tr/, /en/, /ar/, /fr/, /es/, /ru/).
   - Implement hreflang alternates for each localized page.
   - Do NOT force language redirects based on guessed user language.
   - Provide a visible language switcher UI.
4) Accessibility: ensure alt text for every image; ensure html lang is set per page; respect keyboard navigation in template if applicable.
5) Security: do not introduce unsafe HTML injection; sanitize rich text.
6) Performance: use responsive images and avoid huge payloads; provide multiple image sizes if template supports it.
7) Images: use "enterprise realistic" visuals. We will not embed external hotlinks in production. Add assets into AYSTEK Core DB as assets (binary or object references, depending on repo conventions).

LANGUAGES (BCP47)
- tr-TR (primary)
- en-US
- ar-SA
- fr-FR
- es-ES
- ru-RU

CONTENT TYPES TO CREATE
A) Pages (corporate website)
B) Solution Articles (detailed service pages)
C) Blog Posts (SEO + buyer journey aligned)
D) Forum Seed Content (categories + pinned topics + initial Q&A style threads)
E) Assets (images) records in AYSTEK Core (DB seeds)

DISCOVER TEMPLATE & INTEGRATION
Step 1: Inspect repository structure.
- Identify where pages live (e.g., /content/pages, /src/pages, /cms/pages).
- Identify content model (JSON/YAML/MDX) and how locales are handled.
- Identify DB seed / migration mechanism for AYSTEK Core (e.g., /db/seeds, /migrations, /sql, /prisma/seed).
Step 2: Do not change existing architecture. Use conventions.
Step 3: Add content using the existing sample schema; if missing, implement a minimal deterministic schema:
- Content entity: id, type, slug, locale, title, summary, body, seo(metaTitle, metaDescription, canonical), og(ogTitle, ogDescription, ogImageAssetId), breadcrumbs, tags, updatedAt
- Asset entity: id, locale(optional), name, altText, license, source, checksum, mimeType, width, height, variants, storageRef

SITE MAP (CREATE THESE PAGES IN ALL LANGUAGES)
1) Home
2) About AYSTEK Mühendislik
3) AYSTEK Core Platform
4) Solutions (Landing)
   4.1) Card/QR Access + Cashless Top-up + Discounts + Inter-Institution Communication (Enterprise)
   4.2) GPS / Fleet Ledger (Çetele) Analysis & Verification for Personnel and School Transportation
   4.3) School Shuttle Payments Online + Service Registration & Exit Records
   4.4) In-Vehicle Child Safety: Forgotten Child Prevention Solutions
   4.5) In-Vehicle Panic Alarm Service
   4.6) Inventory / Stock Software
   4.7) Custom Solutions (Integration + tailored development)
5) Industries
   5.1) Universities & Campus Ecosystems
   5.2) K-12 Schools & School Transportation
   5.3) Corporate Campuses & Personnel Transportation
6) Security, Privacy & Compliance (Trust Center style)
   - KVKK notice (high-level)
   - GDPR principles (high-level)
   - Payments and PCI DSS approach (no unverified compliance claims)
   - Data minimization and security-by-design statements
7) Resources
   - Blog
   - Forum
   - FAQ
8) Contact / Request Demo

COPY STYLE GUIDE (ENTERPRISE)
- Tone: confident, factual, solution-oriented. Avoid slang.
- Use "we" voice but do not overpromise.
- Provide clear sections: Problem → Approach → Key Capabilities → Integration → Security & Operations → Reporting → FAQs → CTA.
- Avoid deep internal technical detail; keep public-friendly.

GLOSSARY (KEEP CONSISTENT ACROSS LANGUAGES)
- "AYSTEK Core" (do not translate)
- "Çetele" -> keep as "Ledger / Trip Ledger (Çetele)" in EN; provide short explanation the first time.
- "Nakitsiz Kampüs" -> "Cashless Campus"
- "Kartlı Geçiş" -> "Card-Based Access"
- "QR Geçiş" -> "QR Access"
- "Servis" in K-12 context -> "School Shuttle / School Bus Service"
- "Panik Alarm" -> "Panic Alarm"

SEO METADATA RULES (PER PAGE)
- metaTitle: 50–60 chars target
- metaDescription: 150–160 chars target
- Slug: kebab-case; language prefix in URL, not in slug.
- Provide canonical per language URL.
- Provide hreflang alternates list for each page.

IMAGES (ASSET BRIEFS)
Provide 2–4 images per major page; enterprise realistic; no childish icons.
- Home hero: modern enterprise software dashboard (generic)
- Access control: campus gate/turnstile or QR scan at entrance (generic)
- Transportation analytics: fleet dashboards / map route analysis (generic)
- Payments: secure online payment UI abstract (generic)
- Child safety: school bus interior safety concept (generic, non-identifiable children)
- Panic alarm: emergency button concept (generic)
- Inventory: warehouse inventory scanning concept (generic)
For each image create:
- altText (localized)
- captions optional
- usage: hero | inline | ogImage
- variants: xs/s/m/l + webp/avif if supported

IMPLEMENTATION TASKS
Task A: Create content entries for all pages listed above in all locales.
Task B: Create 12 blog posts in tr-TR and en-US; for ar-SA/fr-FR/es-ES/ru-RU create at least 4 "summary translations" each (or full if time).
Task C: Create forum:
- Categories: Announcements, Getting Started, Integrations, Access & Payments, Transportation, Safety Features, Inventory, Troubleshooting, Feature Requests
- Pinned topics:
  - Welcome & Code of Conduct
  - How to ask for help (provide logs, steps, environment)
  - Integration Quickstart (public-friendly)
- Seed 20 Q&A threads (mix technical + operational) with accepted answers.
Task D: Add DB seed/migration so AYSTEK Core can store:
- Page content
- Blog posts
- Forum seed (optional)
- Asset metadata (and binary or storage refs per repo)
Task E: Ensure language switcher exists in navigation and in footer.
Task F: Add sitemap generation if template has it; otherwise add a simple sitemap file list for each locale.

DRAFT TR CONTENT (USE AS PRIMARY SOURCE TEXT)
[Home - tr-TR]
H1: Kampüs ve Taşımacılık Operasyonlarında Enterprise Yazılım Çözümleri
Kısa Özet: AYSTEK Mühendislik; üniversiteler, okullar ve kurumsal yapılarda kartlı/QR geçişten nakitsiz ödemeye, servis ödemelerinden GPS tabanlı operasyon doğrulamaya kadar kritik süreçleri AYSTEK Core ile entegre şekilde yönetmenizi sağlar.
CTA: Demo Talep Et | Çözümleri İncele

[AYSTEK Core - tr-TR]
H1: AYSTEK Core: Tüm Hizmetlerinizi Tek Platformda Birleştirin
Ana Mesaj: Kimlik, yetkilendirme, ödeme/indirim kuralları, cihaz/entegrasyon yönetimi, olay kayıtları ve raporlama gibi yetenekleri ortak bir çekirdekte toplayarak farklı hizmetleri uçtan uca işletmenizi sağlar.
Bölümler: Modüler Mimari, Entegrasyon Yaklaşımı, Operasyon ve Raporlama, Güvenlik Yaklaşımı (public-friendly), SSS.

[Solution pages - tr-TR] (For each solution: produce 800–1200 words; include FAQ and CTA)
- Kart/QR + nakitsiz/topup + indirim + kurumlar arası iletişim
- GPS çetele analizi & doğrulama
- Okul servis ödemeleri online + kayıt/ayrılma
- Çocuk unutma önleme (risk azaltma yaklaşımı, uyarı senaryoları, raporlama)
- Panik alarm (tetikleme, eskalasyon, kayıt altına alma)
- Stok/Envanter (depo, tüketim, rapor)
- Özel çözümler (entegrasyon, IoT, proje bazlı)

TRANSLATIONS
- Translate all pages into en-US, ar-SA, fr-FR, es-ES, ru-RU.
- Keep product names and AYSTEK Core unchanged.
- For ar-SA ensure RTL direction support if template supports; otherwise keep LTR but text Arabic.
- Validate that each localized page has unique metaTitle/metaDescription in that language.

OUTPUT EXPECTATIONS
- Provide PR-ready changes:
  - New content files OR DB seed scripts.
  - Navigation updates.
  - Locale routing & hreflang if missing.
  - Asset records & sample placeholders if binary not included.
- Provide a short README summarizing what was added and where.

END OF PROMPT
```

