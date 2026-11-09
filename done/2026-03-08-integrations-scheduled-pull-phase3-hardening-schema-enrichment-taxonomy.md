# 2026-03-08 - Integrations Scheduled Pull Phase-3 Hardening (Schema + Enrichment + Error Taxonomy)

## Context
- VSCode kesintisi sonrası kod tabanı yeniden doğrulandı.
- Derleme kontrolü tekrar çalıştırıldı ve temiz geçti.

## Completed
1. Re-validation (post-interruption):
   - `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
   - `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

2. Canonical schema validation eklendi (`pull_worker.go`):
   - `strategy_config.canonical_schema.required` alanları doğrulanır.
   - `strategy_config.canonical_schema.types` ile tip doğrulaması yapılır (`string|number|boolean|object|array`).
   - Uyumsuzluklar non-breaking warning + structured issue olarak kaydedilir.

3. Enrichment rules eklendi (`pull_worker.go`):
   - `strategy_config.canonical_enrichment_rules` desteklenir.
   - Kaynak tipleri: `const`, `record`, `cursor`, `schedule`, `lookup_map`.
   - Eksik kaynak veya kural hataları run’ı kırmadan warning/issue üretir.

4. Per-record error taxonomy eklendi:
   - Ingestion issues artık structured biçimde tutulur.
   - `response_meta.ingestion.issues` içinde code/severity/record_index/field/message alanları yazılır.
   - Örnek kodlar: `RAW_STAGE_INSERT_FAILED`, `RECORDS_PATH_NOT_FOUND`, `MAPPING_PARSE_ERROR`, `FIELD_MISSING`, `SCHEMA_TYPE_MISMATCH`, `CANONICAL_UPSERT_FAILED`.

5. Existing Phase-3 foundation ile birleşik davranış:
   - Raw staging + canonical upsert akışı korunur.
   - Ingestion hataları non-breaking olarak warning/issue ile raporlanır.

## File Changes
- `internal/modules/integrations/service/pull_worker.go`
- `plans/INTEGRATIONS_SCHEDULED_PULL_ENTERPRISE_EXECUTION_PLAN_2026-03-07.md`

## Next
1. `publish_events` flag ile opsiyonel domain event publish.
2. Pull ingestion observability metriklerinin genişletilmesi.
