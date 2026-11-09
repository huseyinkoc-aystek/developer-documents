---
Title: Integrations Platform Guide (End-to-End)
Version: 1.4.0
Status: Production
Last Updated: 2026-03-08
Module Path: /internal/modules/integrations
Audience: Product, Backend, Frontend, DevOps, QA
Related Documents:
  - ./modules/integrations.md
  - ./INTEGRATIONS_SHARED_SERVICE_UI_WALKTHROUGH.md
  - ../05-operations/INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ../05-operations/INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
  - ../05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md
  - ../05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md
  - ../05-operations/INTEGRATIONS_SHARED_SERVICE_ALERT_THRESHOLDS.md
  - ../05-operations/GITLAB_PROTECTED_ENV_SETUP.md
---

# Integrations Platform Guide (End-to-End)

Bu doküman Integrations modülünü uçtan uca açıklar: veri modeli, akışlar, UI ekranları, test sistemi, payload mapping/enrichment, response değerlendirme, RabbitMQ davranışı, güvenlik ve operasyon.

## 1. Modülün Amacı
Integrations modülü, AYSTEK içindeki domain event'lerini harici servislere güvenli ve izlenebilir şekilde bağlamak için kullanılır.

Temel hedefler:
1. Domain event -> hook eşleşmesi
2. Hook -> action yürütme (sync/async)
3. Payload dönüştürme (mapping + enrichment + defaults)
4. Harici isteğin güvenli şekilde gönderimi
5. Yanıtın dinamik olarak başarı/failure değerlendirmesi
6. Execution ve DLQ ile tam izlenebilirlik

## 2. Ana Bileşenler
1. Connector
- Harici sistem bağlantı bilgileri (base URL, test URL, auth, timeout, outbound allowlist).

2. Action
- HTTP method/path + request mapping + response schema/mapping.

3. Hook
- Hangi event_key dinlenecek, target tipi (`action|chain`), faz/mod/retry ayarları.

4. Chain
- Çok adımlı DAG akışı (`start/action/transform/branch/end`) ve action node üzerinden connector çağrısı.

5. Execution
- Her çalıştırmanın kayıtları: request/response/status/error/duration.

6. DLQ
- Retry limit sonrası düşen başarısız işler.

## 3. Çalışma Modları
1. Sync Hook
- Hook çağrısı istek kapsamında çalışır.
- Başarı/başarısızlık sonucu doğrudan alınır.

2. Async Hook (RabbitMQ)
- Hook queue'ya bırakılır (`integration_events`).
- Worker tüketir ve execution günceller.
- Retry ve DLQ kuralları uygulanır.

Not:
- RabbitMQ yoksa asenkron hook için kodda kontrollü fallback davranışı vardır; ancak production hedefinde queue-worker hattı esas alınmalıdır.

## 4. Event-Driven Akış
1. Bir modül domain event yayınlar (örn: `stock.vehicle_asset.created`).
2. Integrations event bridge, `event_key` ile hook'ları bulur.
3. Hook mode:
- `sync` ise anında yürütür.
- `async` ise RabbitMQ'ya publish eder.
4. Hook target:
- `target_type=action` ise mevcut action akışı.
- `target_type=chain` ise `hook.chain_id` ile chain runner akışı.
4. Action request payload hazırlanır.
5. HTTP çağrısı yapılır.
6. Response schema/mapping ile sonuç değerlendirilir.
7. Execution kaydı güncellenir.
8. Gerekirse retry veya DLQ çalışır.

## 5. Payload Builder (Mapping + Enrichment + Defaults)
Payload oluşturma pipeline'ı kırıcı olmayan (non-breaking) olacak şekilde tasarlanmıştır.

Sıra:
1. Field mapping
- `externalField -> internalPath` eşlemesi
- rename/nested path desteği

2. Static values
- Her isteğe sabit değer eklenmesi

3. Enrichment
- Primary payload'da olmayan alanları secondary collection/table lookup ile tamamlama
- Join key üzerinden lookup

4. Defaults
- Hala bulunamayan alanlar için `""` veya `null` fallback
- Trigger akışını durdurmadan warning üretimi

Bu yaklaşım sayesinde zorunlu alan eksikliği yüzünden event akışı gereksiz yere kesilmez.

## 6. Response Yakalama ve Başarı Kriteri
Response değerlendirmesi format bağımsızdır:
1. JSON
- Parse + flat path map

2. XML/SOAP
- XML flattening + SOAP fault yakalama

3. Text/Raw
- Ham içerik tutulur

Action içindeki `Yanıt Şeması (Response Schema)` runtime'da aktif başarı kriteridir.
Örnek:
- `gpsUpdateResult = true`
- `meta.status = ok`

Eşleşmezse execution başarısız işaretlenebilir ve hata kodu/mismatch detayı execution metadata içinde saklanır.

## 7. Test Kabiliyetleri
### 7.1 Action Test
Admin action ekranında:
1. Test input verisi hazırlanır
2. Payload preview görülür
3. Test URL ile çağrı atılır
4. Request preview + response preview + warnings gösterilir
5. Sonuç execution-like izlenebilir kayıt olarak saklanır

### 7.2 Hook Test
Hook ekranında:
1. Event key seçilir
2. Test payload verilir
3. Dry-run veya gerçek test çalıştırılır
4. Sonuçlar action bazında görülür

### 7.3 Dinamik Event Catalog
`/admin/integrations/hooks/event-catalog` endpoint'i:
1. Fallback event kataloğunu taşır
2. Tenant içindeki hook + action konfigürasyonunu tarar
3. `available_event_fields`, `field_mapping`, `test_values` kaynaklarından payload alanlarını derler

Bu sayede event key seçiminde kullanıcıya daha doğru örnek payload alanları sunulur.

### 7.4 Chain Test-Run
Chain edit ekranında:
1. Test payload girilir.
2. `dry_run` veya gerçek test çalıştırılır.
3. Step bazlı sonuçlar (status/duration/error) ve final payload önizlenir.
4. Hata satırından ilgili node'a odaklanılır.

### 7.5 Chain Graph Editor Usability
1. Pan/zoom + drag-drop node yerleşimi.
2. Kısayollar: undo/redo, klavye ile node taşıma, delete ile node kaldırma.
3. Snap-to-grid, komşu node hizalama ve auto-layout yardımcıları.

## 8. Güvenlik Modeli
1. SSRF koruması
- Outbound URL çözümlemesi private IP kontrolünden geçer.

2. Connector allowlist
- Domain/IP/CIDR/port/scheme bazlı kurallar
- `allow_private_ips_for_allowlist_match` ile kontrollü private IP erişimi

3. Kimlik doğrulama
- Connector auth ve action header politikaları desteklenir.

4. İzlenebilirlik
- Request/response ve error context execution kayıtlarında tutulur.

## 9. Admin API Yüzeyi (Özet)
Temel endpointler (`/admin/integrations`):
1. Connectors
- `GET /connectors`
- `POST /connectors`
- `GET /connectors/:id`
- `PUT /connectors/:id`
- `DELETE /connectors/:id`

2. Actions
- `GET /actions`
- `POST /actions`
- `GET /actions/:id`
- `PUT /actions/:id`
- `DELETE /actions/:id`
- `POST /actions/execute`

3. Hooks
- `GET /hooks`
- `GET /hooks/event-catalog`
- `POST /hooks`
- `GET /hooks/:id`
- `PUT /hooks/:id`
- `DELETE /hooks/:id`
- `POST /hooks/:id/execute`
- `POST /hooks/:id/test`
- `POST /hooks/bulk-status`

4. Chains
- `GET /chains`
- `POST /chains`
- `GET /chains/:id`
- `PUT /chains/:id`
- `DELETE /chains/:id`
- `POST /chains/:id/validate`
- `POST /chains/:id/test-run`

5. Executions / DLQ
- `GET /executions`
- `GET /executions/:id`
- `POST /executions/:id/retry`
- `GET /dlq`
- `POST /dlq/:id/requeue`
- `POST /dlq/requeue`
- `GET /dlq/analytics`

6. Mapping
- `POST /mapping/validate`

7. Scheduled Pull
- `GET /pull-schedules`
- `POST /pull-schedules`
- `GET /pull-schedules/:id`
- `PUT /pull-schedules/:id`
- `DELETE /pull-schedules/:id`
- `POST /pull-schedules/:id/test-run`
- `GET /pull-runs`
- `GET /pull-runs/:id`

8. Shared Service
- `GET /shared-services`
- `POST /shared-services`
- `GET /shared-services/:id`
- `GET /shared-services/:id/usage-summary`
- `PUT /shared-services/:id`
- `DELETE /shared-services/:id`

9. Shared Service Grant
- `GET /shared-service-grants`
- `POST /shared-service-grants`
- `GET /shared-service-grants/:id`
- `GET /shared-service-grants/:id/usage-summary`
- `PUT /shared-service-grants/:id`
- `DELETE /shared-service-grants/:id`

## 10. UI Ekran Akışı
1. Connectors
- Base URL/Test URL, allowlist, auth ayarları

2. Actions
- Field mapping, static values, test input, body template, response schema

3. Hooks
- Event key, mode/faz, retry, target selector (Action/Chain), test panel

4. Chains
- Görsel editor (drag-drop node), node tipine özel config paneli, validate ve test-run paneli

5. Executions
- Request/response, status, timing, hata detayları

6. DLQ
- Fail olmuş kayıtların analiz/requeue yönetimi

7. Scheduled Pulls / Pull Runs
- Pull schedule CRUD + test-run yönetimi
- Pull run list/detail ekranında ingestion summary + issues görünürlüğü

8. Shared Services / Service Grants
- Shared service/grant CRUD
- Policy alanları: scope/method/path/limits + secret binding + payload guard + mTLS
- Edit ekranında 7/14/30 gün usage summary ve warning paneli
- Quota baskısı/deny trendi için grant bazlı observability kartı

## 11. RabbitMQ Doğrulama (Operasyon)
RabbitMQ hattının canlı doğrulaması için:
1. `user-guide/05-operations/INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md`
2. `user-guide/05-operations/INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md`
3. `user-guide/05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`
4. `user-guide/05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md`
5. `user-guide/05-operations/GITLAB_PROTECTED_ENV_SETUP.md`
6. `go/aystek-core/tools/ops/integrations_eventbridge_smoke.sh`
7. `go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`
8. `user-guide/05-operations/scripts/integrations_pull_smoke.sh`
9. `user-guide/05-operations/scripts/integrations_shared_service_smoke.sh`
10. Pipeline job (`.gitlab-ci.yml`) ile post-deploy smoke

## 12. Sık Karşılaşılan Sorunlar
1. Event tetikleniyor ama execution oluşmuyor
- event_key mismatch
- hook disabled
- tenant/site context uyumsuz

2. Execution queued kalıyor
- worker process yok
- queue bağlantı sorunu

3. SOAP çağrısı 500 dönüyor
- `SOAPAction` yanlış format
- body template token değerleri eksik

4. SSRF private IP hatası
- hedef private range'e çözülüyor
- allowlist ve private-ip-match politikası kontrol edilmeli

## 13. Genişletme Rehberi
Yeni entegrasyon eklerken önerilen sıra:
1. Connector oluştur
2. Action request/response şemasını tanımla
3. Hook event key bağla
4. Test URL ile action/hook test et
5. Execution ve warning'leri doğrula
6. Retry/DLQ davranışını smoke test ile doğrula

## 14. Kapsam ve Sınırlar
Bu modül:
- Trigger mekanizmasını merkezileştirir
- Payload/response katmanını esnekleştirir
- Operasyonel izlenebilirlik sağlar

Bu modül dışında kalanlar:
- Domain event üretiminin iş kuralları (kaynak modüller)
- Harici servisin kendi SLA ve semantik kısıtları

## 15. Sonuç
Integrations modülü şu anda:
1. Çoklu modül event key bağlama
2. Hook target modeli ile action veya chain hedefleme
3. Mapping/enrichment/defaults ile dayanıklı payload üretimi
4. Format bağımsız response değerlendirme
5. Queue-worker ve DLQ ile operasyonel güvence
6. Admin UI üzerinden action/chain test + izleme

özellikleriyle enterprise düzeyde çalışacak mimariyi sağlar.

## 16. Kullanıcı İşletim Rehberi
Bu bölüm doğrudan admin kullanıcısının günlük kullanımını anlatır.

### 16.1 Event-Driven Entegrasyon Kurulumu
1. `Integrations > Connectors` ekranında connector oluşturun.
2. `Integrations > Actions` ekranında method/path, field mapping, static values, body template ve response schema tanımlayın.
3. `Integrations > Hooks` ekranında `event_key` seçip target belirleyin (`action` veya `chain`).
4. Hook mode’u seçin: `sync` veya `async`.
5. Hook test çalıştırın ve sonucu `Integrations > Executions` ekranından doğrulayın.

### 16.2 Action Test Kullanımı
1. Action detayında test input alanlarını doldurun.
2. Test çalıştırın.
3. Request preview (url/headers/body), response, warnings ve response schema eşleşmesini kontrol edin.

### 16.3 Chain Kullanımı
1. `Integrations > Chains` ekranında chain oluşturun.
2. Node’ları görsel editörde ekleyin (`start/action/transform/branch/end`).
3. `Validate` ile grafı doğrulayın.
4. `Test Run` ile node timeline çıktısını inceleyin.
5. Hook target’ını `chain` seçerek canlı akışa bağlayın.

### 16.4 Scheduled Pull Kullanımı
1. `Integrations > Scheduled Pulls` ekranında schedule oluşturun:
- connector + action
- trigger mode (`cron` / `interval`)
- incremental strategy
- retry/timeout limitleri
- strategy config (JSON)
2. Schedule detayında test-run yapın (`dry_run` opsiyonel).
3. `Integrations > Pull Runs` ekranında run geçmişini izleyin.
4. Run detail ekranında ingestion summary/issues, warnings ve request/response metadata alanlarını inceleyin.
5. Gerekirse run detail’den `Retry` aksiyonu ile manuel yeniden kuyruğa alın.

### 16.5 Hangi Sorunda Hangi Ekran
1. Mapping/body sorunu: `Actions`.
2. Event tetikleme sorunu: `Hooks`.
3. Akış node hatası: `Chains` + timeline.
4. Runtime çağrı hatası: `Executions`.
5. Retry limiti aşımı: `DLQ`.
6. Periyodik çekim sorunları: `Pull Runs` detail.

### 16.6 Operasyon Kısa Listesi
1. Konfigürasyon sonrası önce test-run yapın.
2. Async akışlarda worker/RabbitMQ sağlığını kontrol edin.
3. SSRF/allowlist bloklarında connector güvenlik ayarlarını güncelleyin.
4. Sık 429/5xx hatalarında retry/backoff ve schedule interval değerlerini birlikte optimize edin.
5. Post-deploy scheduled pull smoke:
- `user-guide/05-operations/scripts/integrations_pull_smoke.sh`
