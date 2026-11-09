# 2026-03-08 - Integrations Scheduled Pull Phase-3 Hardening (publish_events + Observability)

## Completed
1. Optional domain event publish eklendi:
   - `internal/modules/integrations/service/pull_worker.go`
   - `strategy_config.publish_events=true` olduğunda canonical upsert sonrası event yayınlanır.
   - Event:
     - `type`: `integrations.pull.canonical.upserted`
     - `scope`: `site`
     - `entity_type`: `integration_pull_canonical`
     - payload: `schedule_id`, `schedule_key`, `run_id`, `record_key`, `source_hash`, `normalized`, `observed_at`

2. Ingestion observability metrikleri eklendi:
   - `internal/modules/integrations/metrics/metrics.go`
   - Yeni metrikler:
     - `aystek_integrations_pull_ingestion_raw_stage_total`
     - `aystek_integrations_pull_ingestion_canonical_records_total`
     - `aystek_integrations_pull_ingestion_issues_total`
     - `aystek_integrations_pull_ingestion_events_published_total`

3. Worker metrik entegrasyonu:
   - `pull_worker.go` içinde raw stage success/error, canonical processed/upserted/failed, issue code/severity, event publish success/error sayaçları işlendi.

4. Non-breaking garanti korunmuştur:
   - Event publish hatası ingestion run’ını fail etmez.
   - Hatalar warning + structured issue olarak response_meta içinde raporlanır.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

## Next
1. Phase-4: Pull run detail ekran/API tarafında ingestion summary + issues panel görünürlüğü.
