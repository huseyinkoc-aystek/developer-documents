# 2026-03-07 - Integrations Scheduled Pull Phase-1 Seed Completion + Phase-2 Kickoff

## Completed
1. Seed standard updates
   - `internal/seeds/menu_seed.go`
     - Integrations altında `pull-schedules` ve `pull-runs` menü kayıtları eklendi.
   - `internal/seeds/translation_seed.go`
     - Pull schedule/runs menü anahtarları eklendi.
     - Pull schedule notice/validation anahtarları eklendi.

2. Phase-2 başlangıç implementasyonu
   - `internal/modules/integrations/service/pull_worker.go`
     - `ProcessDuePullSchedules` metodu eklendi.
     - Akış: due schedules -> pull run oluştur -> action execute -> run status update -> schedule next_run_at update.
   - `internal/modules/integrations/handler/internal.go`
     - `ExecuteDuePull` internal handler eklendi.
   - `internal/modules/integrations/routes.go`
     - `POST /internal/integrations/pull/execute-due` route eklendi.

3. Service interface genişletmesi
   - `internal/modules/integrations/service/service.go`
     - `ProcessDuePullSchedules` sözleşmesi eklendi.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./internal/seeds -run TestDoesNotExist` -> PASS

## Next
1. Dispatcher/scheduler katmanını tenant/site kapsamlı şekilde bağla.
2. Pull run idempotency + retry/backoff + DLQ semantiğini ekle.
3. Cursor/checkpoint commit davranışını success-only netleştir.
