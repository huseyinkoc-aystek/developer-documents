# 2026-03-07 - Integrations Scheduled Pull Phase-1 API Foundation

## Scope
Phase-1 kapsamında scheduled pull capability için backend API foundation katmanı tamamlandı.

## Completed
1. Integrations repository ana interface’i pull schedule/run/cursor/raw metodları ile genişletildi.
2. Mongo repository implementasyonu eklendi:
   - pull schedule CRUD + soft delete + due listing
   - pull run CRUD/list/find
   - pull cursor upsert/find
   - incoming raw insert
3. Service katmanı eklendi:
   - `Create/Update/Delete/Get/ListPullSchedule`
   - `List/GetPullRun`
   - `TestPullSchedule` (action direct execute üzerinden test run kaydı üretir)
4. Admin handler + route katmanı eklendi:
   - `/admin/integrations/pull-schedules` CRUD + `/test-run`
   - `/admin/integrations/pull-runs` list/detail
5. Integrations service test mock repository, yeni interface metodlarını destekleyecek şekilde güncellendi.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

## Next
1. Pull schedule/runs için translation + permission + menu seed tamamlanacak.
2. Scheduler dispatcher + worker runtime (Phase-2) başlatılacak.
