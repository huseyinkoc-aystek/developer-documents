# 2026-03-07 - Integrations Scheduled Pull Phase-2 Idempotency + Retry/DLQ

## Completed
1. Pull retry due finder repository seviyesine eklendi:
   - `internal/modules/integrations/repository/pull_interface.go`
   - `internal/modules/integrations/repository/interface.go`
   - `internal/modules/integrations/repository/mongo.go` (`ListDuePullRuns`)

2. Pull worker semantiği enterprise pattern’e taşındı:
   - `internal/modules/integrations/service/pull_worker.go`
   - Deterministic run key: `pull-schedule:<schedule_id>:<slot_epoch>`
   - Duplicate run key koruması (find-or-create + duplicate key race handling)
   - Due schedule için pencere ilerletme (`next_run_at`) ve ayrı retry işlem hattı

3. Retry/backoff davranışı eklendi:
   - Pull run başarısızlığında `retry_count` artışı
   - `next_retry_at = now + backoff.CalculateBackoffDuration(attempt)`
   - `ProcessDuePullRetries` ile queued retry run’larının işlenmesi

4. DLQ semantiği eklendi:
   - Retry limit aşımında pull run failed + DLQ insert
   - DLQ reason: `pull_retry_limit_exceeded`
   - Pull run bağlamı metadata/request_payload içinde taşınıyor (`pull_run_id`, `run_key`, `schedule_id`)

5. Internal/scheduler entegrasyonu tamamlandı:
   - `internal/modules/integrations/handler/internal.go` (`RetryDuePull`)
   - `internal/modules/integrations/routes.go` (`POST /internal/integrations/pull/retry-due`)
   - `internal/modules/integrations/jobs/pull_scheduler.go` (execute-due sonrası retry-due çağrısı)

6. Test/mock uyumu:
   - `internal/modules/integrations/service/service_test.go` mock repo `ListDuePullRuns` eklendi.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

## Notes
- Hook tetikleme semantiğine dokunulmadı; değişiklikler sadece pull schedule runtime hattında.
- Retry/DLQ akışı mevcut integrations backoff yaklaşımıyla uyumlu çalışır.

## Next
1. Cursor/checkpoint success-only commit hardening.
2. Raw staging + canonical mapping/upsert pipeline (Phase-3 başlangıcı).
