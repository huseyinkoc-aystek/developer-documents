ROLÜN:
Sen AYSTEK Core üzerinde çalışan, prod deneyimi yüksek kıdemli yazılım mühendisisin. Görevin: mevcut “kanıtlı gap” raporuna dayanarak eksikleri kapatacak, uygulanabilir bir **Master Closure Plan** üretmek. Çıktı, ekip tarafından doğrudan sprint/PR’lara bölünerek uygulanabilir olmalı.

GİRDİLER (mutlaka kullan):

1. /workspace/aystek-core-dev-docs/04-AYSTEK_GOLIVE_READINESS_REPORT_TR.md  (kanıtlı P0/P1 gap raporu)
2. /workspace/aystek-core-dev-docs/AYSTEK_MASTER_WORKPLAN_TR.md (master fazlar)
3. Repo: /workspace/go/aystek-core
   NOT: Next.js tarafında düzenleme yapılmadı. UI işleri sadece backlog olarak listelensin; uygulama planı backend odaklı olsun.

ALTIN KURAL:

* Varsayım yok. Raporda “⚠️ PARTIAL / ❌ MISSING” olan maddeleri tek tek kapat.
* Her iş maddesi için: dosya yolu(lar) + hedeflenen fonksiyon/struct + test planı + kabul kriteri yaz.
* “Modül var” demek yetmez: P0 maddeleri prod-safe olacak (idempotency, concurrency, replay, metrics, alerting, reconciliation).

ÇIKTI FORMAT:
Tek bir Markdown doküman üret: **05-AYSTEK_CLOSURE_MASTER_PLAN_TR.md**

DOKÜMAN BÖLÜMLERİ (zorunlu):

# 1) Executive Summary

* P0 skoru, P1 skoru (rapordan al)
* Kapanması gereken kritik gap listesi (en fazla 10 madde)
* Tahmini efor aralığı (saat/gün) ve risk

# 2) P0 Closure Plan (Go-Live)

P0’daki “⚠️ PARTIAL / ❌” olan her madde için ayrı epik aç:

* P0-G1: Order optimistic locking / concurrency hardening
* P0-G2: Webhook replay protection (timestamp/nonce window) + kanıt
* P0-G3: FX snapshot persistence (Order/Payment/Invoice)
* P0-G4: Shipping gerçek carrier adapter MVP (1 taşıyıcı)
* P0-G5: Observability metrics + alerts (Prometheus/OTel) + dashboard kartları
* P0-G6: Invoice refund → credit memo senaryosu + test

Her epikte şu şablonu doldur:

## [EPIC KODU] Başlık

### Problem Kanıtı

* Rapor referansı: ilgili bölüm ve bulgu (kısa özet)

### Hedef Davranış (Done)

* Sistem davranışı maddeler halinde (prod-safe)

### Teknik Tasarım

* Değişecek/eklenecek ana yapı taşları
* DB değişikliği var mı? (model alanı, index, migration)
* Event/queue etkisi var mı? (outbox, DLQ)

### PR Planı (parçala)

Her PR için:

* PR-1: Başlık

  * Değişiklik kapsamı
  * Dokunulacak dosyalar (path listesi)
  * Yeni/Değişen fonksiyonlar (adıyla)
  * Testler (unit/integration) + hangi dosyalarda
  * Kabul kriterleri (Given/When/Then)
  * Geriye dönük uyumluluk notu

### Riskler ve Mitigasyon

* Data migration riski, concurrency riski, geri alma planı

# 3) P1 Closure Plan (Platform Derinliği)

Raporda “❌/⚠️” olan P1 maddeleri için aynı şablonla ama daha düşük öncelik:

* Tax rules implementation (service/handler + hesap doğrulama)
* B2B: company profile + credit limit + approval + payment terms (MVP)
* Inventory: bulk import/export job tracking + low stock notifications (kanıtlı eksikse)
* Returns: RMA lifecycle + partial refund kuralları (kanıtlı eksikse)
  Not: Burada “raporda kanıt yoksa” önce keşif PR’ı (spike) aç.

# 4) Test Stratejisi (Go-Live Gate)

* P0 için minimum E2E senaryo listesi (50+ gibi)
* Contract test kapsamı: webhook signature + idempotency + retry
* Load test hedefleri: 100 concurrent order/payment, webhook burst
* “Release Gate”: hangi testler geçmeden deploy yok

# 5) Next.js UI Backlog (Sadece liste)

Next düzenlenmediği için aşağıdaki ekranları “MVP go-live admin” olarak listele:

* Order timeline + state machine görselleştirme
* Payment status + webhook log görüntüleme
* Refund/Credit memo ekranı
* FX snapshot görüntüleme (order/payment/invoice)
* Shipping timeline + tracking
  Sadece backlog yaz; planlama backend’den sonra.

# 6) Uygulama Takvimi (2 Haftalık Sprint Planı)

* Gün gün: P0 PR sırası (bağımlılıkları gözet)
* Her günün çıktısı (deliverable) net olsun
* “Stop-the-line” kriteri: kırmızı risk varsa ne durdurur?

UYGULAMA TALİMATI:

* Önce 04-AYSTEK_GOLIVE_READINESS_REPORT_TR.md içinden P0 “⚠️/❌” maddelerini çıkar.
* Sonra repo içinde ilgili dosyaları aç, gerçek isimleri kullan.
* Planı “PR’lara böl” ve her PR’ı atomik, review edilebilir yap.
* Çıktıyı tek seferde, tam dolu üret.

BAŞLA.
