# Integrations Scheduled Pull Enterprise Execution Plan

Date: 2026-03-07
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / user-guide
Source:
- `/workspace/aystek-core-dev-docs/plans/AYSTEK CORE Integrations Modülüne Enterprise Düzeyde Zamanlanmış Dış Veri Çekme Yeteneği Eklenmesi.pdf`
- `/workspace/aystek-core-dev-docs/plans/AYSTEK CORE Integrations Modülüne Enterprise Düzeyde Zamanlanmış Dış Veri Çekme Yeteneği Eklenmesi..md`

## Objective
Integrations modülüne mevcut event-driven mimariyi bozmadan schedule-driven external data pull (periodic fetch + cursor/backfill + run observability + retry/DLQ) yeteneği eklemek.

## Guardrails
1. Hook tetikleme semantiği değişmeyecek (geri uyumluluk zorunlu).
2. Multi-tenant/site izolasyonu tüm yeni kaynaklarda zorunlu.
3. Async execution hattı RabbitMQ + worker temelli kalacak.
4. Güvenlik katmanı (SSRF/allowlist/auth/secrets) yeni pull akışında da zorunlu.
5. DDD katmanları + seed + translation + menu standartları korunacak.

## Delivery Phases

### Phase-1: Domain Foundation (Backend contracts)
1. `PullSchedule`, `PullRun`, `PullCursor`, `IncomingRaw` model kontratları.
2. DTO sözleşmeleri (create/update/list/test-run).
3. Repository contract extension (pull schedule/run/cursor CRUD + due finder).
4. Mongo index foundation (schedule due scanning, run idempotency, cursor lookup).
5. Validation error key seti ve normalize kuralları.

### Phase-2: Scheduler + Worker Runtime
1. Scheduler service: due schedule bulma + queue publish.
2. Pull worker: dequeue + connector/action fetch + retry/backoff + DLQ.
3. Run-level idempotency key ve duplicate run koruması.
4. Cursor/checkpoint update semantiği (success-only commit).

### Phase-3: Ingestion Pipeline
1. Raw staging (TTL) + canonical upsert pipeline.
2. Mapping/normalization + JSON schema validation + enrichment/defaults.
3. Record-level dedup key stratejisi.
4. Domain event publish (opsiyonel flag ile).

### Phase-4: Admin API + Next UI
1. Pull schedules list/create/edit/test-run ekranları.
2. Pull run list/detail + step/metrics + warnings paneli.
3. Retry/requeue operasyonları ve DLQ görünürlüğü.
4. Translation + menu seed + permission seed güncellemeleri.

### Phase-5: Ops/Hardening
1. Smoke runbook (schedule -> queue -> worker -> run/dlq).
2. CI post-deploy smoke job.
3. SLO dashboard metrikleri (due lag, success ratio, rate-limit hit).
4. Security/compliance checklist (KVKK/GDPR processing evidence).

## Immediate Sprint Scope (Started)
1. Phase-1 model/dto/repository/index foundation coding.
2. Compile checks (`go test` scoped packages).
3. Plan/progress dokümantasyonu.

## Progress Status (Updated: 2026-03-07)

### Completed
1. Plan dokümanı oluşturuldu ve roadmap netleştirildi.
2. Phase-1 başladı ve ilk foundation teslimi yapıldı:
   - `model/pull.go` ile PullSchedule/PullRun/PullCursor/IncomingRaw sözleşmeleri eklendi.
   - `dto/pull.go` ile pull schedule/run query-request DTO sözleşmeleri eklendi.
   - `repository/pull_interface.go` ile phased pull repository contract tanımlandı.
   - `indexes.go` içine pull schedule/run/cursor/raw koleksiyon indexleri ve raw TTL index foundation eklendi.
3. Phase-1 backend API foundation genişletildi:
   - `repository/interface.go` pull schedule/run/cursor/raw kontratları ana Repository interface’ine entegre edildi.
   - `repository/mongo.go` içinde pull schedule/run/cursor/raw için Mongo implementasyonları eklendi.
   - `service/pull_schedule_service.go` ile CRUD + list + test-run servis akışı eklendi.
   - `handler/admin.go` ve `routes.go` ile admin REST endpointleri eklendi:
     - `GET/POST /admin/integrations/pull-schedules`
     - `GET/PUT/DELETE /admin/integrations/pull-schedules/:id`
     - `POST /admin/integrations/pull-schedules/:id/test-run`
     - `GET /admin/integrations/pull-runs`
     - `GET /admin/integrations/pull-runs/:id`
   - Integrations modülü derleme doğrulaması tekrar çalıştırıldı: `go test ./internal/modules/integrations/... -run TestDoesNotExist` (başarılı).
4. Phase-1 seed standardı tamamlandı:
   - `menu_seed.go` içine pull schedules/runs menü kayıtları eklendi.
   - `translation_seed.go` içine pull schedule notice/validation ve menü anahtarları eklendi.
5. Phase-2 başlangıcı yapıldı:
   - `ProcessDuePullSchedules` servis metodu eklendi (due schedule -> run -> action execute -> run update -> next_run_at hesaplama).
   - Internal worker endpoint eklendi: `POST /internal/integrations/pull/execute-due`.
   - Phase-2 baseline compile doğrulaması geçti.
6. Phase-2 dispatcher başlangıcı eklendi:
   - `internal/modules/integrations/jobs/pull_scheduler.go` ile pull scheduler job eklendi.
   - Scheduler, `integration_pull_schedules` koleksiyonundan aktif tenant/site scope’larını toplayıp scope bazlı `ProcessDuePullSchedules` çağırır.
   - `cmd/api/main.go` içine scheduler bootstrap bağlantısı eklendi.
7. Phase-2 idempotency + retry/backoff + DLQ semantiği eklendi:
   - `ListDuePullRuns` repository contract/implementasyonu eklendi.
   - `ProcessDuePullSchedules` deterministic run key (`pull-schedule:<schedule_id>:<slot_epoch>`) + duplicate koruması ile güncellendi.
   - `ProcessDuePullRetries` ile queued retry run’ları backoff zamanına göre işlendi.
   - Pull run başarısızlıkları retry limit aşımında DLQ’ye yazılır hale getirildi (`reason: pull_retry_limit_exceeded`).
   - Internal endpoint eklendi: `POST /internal/integrations/pull/retry-due`.
   - Scheduler tick akışına retry processing de eklendi.
8. Phase-2 cursor/checkpoint success-only commit hardening tamamlandı:
   - Pull attempt başlangıcında `cursor_before` snapshot (run üzerinde) alınır.
   - Cursor state run input’una `_cursor` altında taşınır.
   - Sadece başarılı run sonrası `cursor_after` üretilir ve `integration_pull_cursors` upsert edilir.
   - Retry/failed durumlarında cursor commit yapılmaz (checkpoint korunur).
   - `strategy_config.cursor_path` / `strategy_config.cursor_field` ile cursor çıkarımı konfigüre edilebilir.
9. Phase-3 foundation başlatıldı (raw staging + canonical upsert):
   - `integration_pull_canonical` koleksiyonu için model/repository/index foundation eklendi.
   - Pull success akışında raw response `integration_incoming_raw` içine TTL ile stage edilir.
   - Strategy-config destekli canonical record extraction + normalize + defaults + required-field warning pipeline eklendi.
   - Canonical records `schedule_id + record_key` dedup anahtarıyla upsert edilir.
10. Phase-3 hardening tamamlandı (schema validation + enrichment + error taxonomy):
   - Canonical normalize akışına `canonical_schema` doğrulama eklendi (`required`, `types`).
   - `canonical_enrichment_rules` ile record/cursor/schedule/const/lookup_map kaynaklı zenginleştirme eklendi.
   - Per-record issue taxonomy, `response_meta.ingestion.issues` alanında structured olarak yazılır.
   - Ingestion hataları non-breaking kalır; run başarısını kırmadan warning/issue ile raporlanır.
11. Phase-3 hardening tamamlandı (`publish_events` + observability metrics):
   - `strategy_config.publish_events` flag’i ile canonical upsert sonrası domain event publish opsiyonel hale getirildi.
   - Event tipi: `integrations.pull.canonical.upserted` (scope: site, entity: integration_pull_canonical).
   - Pull ingestion metrikleri eklendi:
     - `pull_ingestion_raw_stage_total`
     - `pull_ingestion_canonical_records_total`
     - `pull_ingestion_issues_total`
     - `pull_ingestion_events_published_total`
   - Metric label setleri tenant/site/schedule ekseninde çalışacak şekilde bağlandı.
12. Phase-4 API/Next görünürlük adımı tamamlandı (pull run detail ingestion summary + issues):
   - Next tarafında `integrationsApi` genişletildi:
     - `listPullRuns(params)` ve `getPullRun(id)` metotları eklendi.
     - Pull run parser’ına `response_meta.ingestion` / `response_meta.ingestion_summary` parse desteği eklendi.
   - Next UI ekranları eklendi:
     - `/admin/integrations/pull-runs` liste ekranı (status/trigger/schedule filtreleri + detail linki)
     - `/admin/integrations/pull-runs/[id]` detay ekranı
   - Detay ekranına ingestion görünürlüğü eklendi:
     - Ingestion Summary kartları (raw/canonical/upsert/failed/issues/events/publish flag)
     - Ingestion Issues tablosu (code/severity/record/field/message)
     - Warnings paneli + request/response meta JSON panelleri.
   - Translation seed anahtarları pull-runs alanına genişletildi (liste, durum, trigger, ingestion, issues, warnings, meta başlıkları).
   - Scoped doğrulama:
     - `npx eslint` (yalnızca değişen dosyalar) başarılı.
     - `go test ./internal/seeds -run TestDoesNotExist` başarılı.
13. Phase-4 kalanları tamamlandı (pull schedules UI):
   - Next tarafında pull schedule tipleri ve API metotları eklendi:
     - `list/get/create/update/delete/test` pull schedule işlemleri
   - UI ekranları eklendi:
     - `/admin/integrations/pull-schedules`
     - `/admin/integrations/pull-schedules/new`
     - `/admin/integrations/pull-schedules/:id`
   - Formda connector/action seçimi, trigger mode (cron/interval), incremental strategy, retry/timeout limitleri ve `strategy_config` JSON yönetimi eklendi.
   - Edit ekranına schedule test-run paneli (payload + dry-run + sonuç preview) eklendi.
14. Phase-5 ops dokümantasyon finalizasyonu tamamlandı:
   - Yeni runbook:
     - `user-guide/05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`
   - `user-guide/README.md`, `OPERATIONS_GUIDE.md`, `INTEGRATIONS_PLATFORM_GUIDE.md` scheduled pull runbook referanslarıyla güncellendi.
15. Gap closure (kaynak doküman + plan karşılaştırma sonrası):
   - Scheduled pull için manuel retry operasyonu eklendi:
     - Backend: `POST /admin/integrations/pull-runs/:id/retry`
     - Service: `RetryPullRun` (queued/running için blok, failed/cancelled/success sonrası manuel requeue)
     - Next: pull-run detail ekranına `Retry` butonu + API client metodu
   - Translation seed’e ilgili notice/validation/action anahtarları eklendi.
16. Phase-5 CI smoke hardening tamamlandı:
   - `user-guide/05-operations/scripts/integrations_pull_smoke.sh` eklendi.
   - Scheduled pull runbook içine script kullanım + GitLab job örneği eklendi.

### Next
1. Plan kapsamındaki ana maddeler tamamlandı (kaynak dokümanla hizalı).
2. Sonraki iterasyon (opsiyonel): pull ve eventbridge smoke job'larını tek birleşik pipeline template altında birleştirme.
