# 2026-03-08 - Integrations Scheduled Pull Phase-3 Foundation (Raw Staging + Canonical Upsert)

## Completed
1. Canonical model/repository foundation eklendi:
   - `internal/modules/integrations/model/pull.go`
   - `internal/modules/integrations/repository/pull_interface.go`
   - `internal/modules/integrations/repository/interface.go`
   - `internal/modules/integrations/repository/mongo.go`
   - Yeni model: `PullCanonicalRecord`
   - Yeni repo metodu: `UpsertPullCanonicalRecord`

2. Mongo collection/index foundation eklendi:
   - `internal/modules/integrations/indexes.go`
   - Yeni collection: `integration_pull_canonical`
   - Unique index: `(tenant_id, site_id, schedule_id, record_key)`
   - Query index: `(tenant_id, site_id, schedule_id, observed_at desc)`

3. Pull worker başarı akışında ingestion pipeline eklendi:
   - `internal/modules/integrations/service/pull_worker.go`
   - `ingestPullSuccess(...)` fonksiyonu ile:
     - Raw response `integration_incoming_raw` içine TTL ile stage edilir (`raw_ttl_hours`, default 72h)
     - `records_path` ile record extraction yapılır (default `result.data`)
     - Canonical normalize adımı çalışır (mapping/default/required rules)
     - Her record için dedup key üretilip canonical upsert yapılır

4. Strategy-config destekleri eklendi (Phase-3 foundation):
   - `records_path`
   - `canonical_mapping_expr` (JMESPath)
   - `canonical_field_mapping`
   - `canonical_defaults`
   - `canonical_required_fields`
   - `canonical_dedup_path`
   - `raw_ttl_hours`

5. Non-breaking davranış korundu:
   - Ingestion hataları run’ı fail etmez; warnings içine yazılır.
   - Hook/event tetikleme semantiğine müdahale edilmedi.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

## Next
1. Canonical schema validation + enrichment rules + per-record error taxonomy.
2. Domain event publish flag (`publish_events`) + ingestion observability metrics.
