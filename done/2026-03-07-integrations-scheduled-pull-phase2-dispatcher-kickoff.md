# 2026-03-07 - Integrations Scheduled Pull Phase-2 Dispatcher Kickoff

## Completed
1. Pull scheduler job eklendi:
   - `internal/modules/integrations/jobs/pull_scheduler.go`
   - Aktif pull schedule scope’larını (`tenant_id`,`site_id`) aggregate eder.
   - Her scope için `ProcessDuePullSchedules` çağırır.

2. API bootstrap entegrasyonu:
   - `cmd/api/main.go` içine `integrationsjobs.StartPullScheduler(...)` eklendi.

3. Phase-2 worker/internal endpoint:
   - `POST /internal/integrations/pull/execute-due` endpointi mevcut hale getirildi.

4. Seed standardı:
   - Menü ve çeviri anahtarları pull schedules/runs için eklendi.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS
- `go test ./internal/seeds -run TestDoesNotExist` -> PASS

## Next
1. Pull run idempotency + retry/backoff + DLQ semantiği.
2. Cursor/checkpoint success-only commit hardening.
3. Raw staging + canonical upsert pipeline (Phase-3 başlangıcı).
