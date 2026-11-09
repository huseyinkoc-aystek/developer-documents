# AYSTEK Taşımacılık, Geçiş, Çetele ve Veli/Personel Tahsilat Platformu

Bu klasör, `gursel-turizm` tenant’ı altında hem KT kartlı/QR’lı taşımacılık-geçiş-çetele sisteminin hem de `veli.gurseltur.com.tr` veli/personel tahsilat self-service kanalının AYSTEK Core içine taşınması için kalıcı proje iş planı alanıdır.

## Dokümanlar

### Planlama

- `MASTER_WORKPLAN.md` — Enterprise düzeyde ana iş planı, mimari kararlar, tam koleksiyon tablosu (4 domain, 45 koleksiyon), fazlar, kalite kapıları ve devam checklist'i.
- `SCOPE_RECONCILIATION_AND_TARGET_ARCHITECTURE.md` — Kullanıcı anlatımı ile birleşik rapor arasındaki farklar, temizlenmiş hedef kapsam ve genişletilmiş ürün mimarisi.

### Mimari Kararlar (ADR)

- `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md` — `transport_core`, `transport_access`, `transport_wallet`, `transport_portal` domain sınırları ve namespace kararı. SAP dispatch = `transport_core/sapdispatch`.
- `ADR_002_OFFLINE_FIRST_SYNC.md` — Offline-first cihaz geçiş senkronizasyon protokolü; cihaz policy snapshot, append-only log, idempotent batch sync, conflict resolution ve risk limitleri.
- `ADR_003_MULTI_POS_BANK_ROUTING.md` — Çoklu banka POS routing; BIN/IIN tabanlı acquirer seçimi, iki aşamalı provisional/final quote modeli ve TOSLA top-up ayrımı.

### Uygulama Planları

- `FAZ_0_DISCOVERY_EXECUTION_PLAN.md` — Keşif fazı uygulama planı; blocker soru listesi, legacy sistem analizi ve 5 keşif dokümanı için üretim kılavuzu.
- `SPRINT_01_FOUNDATION_EXECUTION_PLAN.md` — İlk uygulanabilir foundation sprint'inin dosya, seed, route, frontend ve test planı.

### Doğrulama

- `VALIDATION_AND_GAP_REPORT.md` — Codex çıktılarının doğrulanması; onaylanan kararlar, kritik boşluklar ve yapılan düzeltmelerin özeti.

## Faz 0'da Oluşturulacak Dokümanlar (henüz yok)

- `DISCOVERY_FINDINGS.md` — Belirsizlik sınıflandırma tablosu (BLOCKER/ASSUMPTION/ACCEPTED RISK)
- `LEGACY_MAPPING_MATRIX.md` — PHP/Yii2 alan → domain model eşleme tablosu
- `BANK_AND_POS_DECISION_RECORD.md` — Aktif banka POS profilleri ve BIN routing kararı
- `SAP_ADAPTER_DECISION_RECORD.md` — SAP sürümü, bağlantı tipi ve dispatch payload şeması
- `KVKK_SPECIAL_CATEGORY_DATA_POLICY.md` — Kan grubu ve özel nitelikli veri politikası

## Kaynak Referanslar

- `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`
- `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`
- `/workspace/php/GurselVB22` — Legacy veli borç ödeme uygulaması (Yii2)
- `/workspace/php/AYShuttleGursel` — Legacy araç takip / KT uygulaması (Yii2)
- `/workspace/go/aystek-core/internal/seeds/gursel_turizm_seed.go` — Mevcut Gürsel tenant ve aractakip site seed'i
- `/workspace/aystek-core-dev-docs/user-guide`
