# 2026-03-07 - Integrations Scheduled Pull Phase-2 Cursor/Checkpoint Success-Only Hardening

## Completed
1. Pull worker cursor snapshot/commit semantiği harden edildi:
   - `internal/modules/integrations/service/pull_worker.go`
   - Run başında `cursor_before` snapshot alınır (`run.cursor_before`).
   - Cursor state action input’una `_cursor` alanı ile geçirilir.

2. Success-only checkpoint commit eklendi:
   - Sadece `execErr == nil` durumunda `cursor_after` hesaplanır ve `UpsertPullCursor` çağrılır.
   - `integration_pull_cursors.last_success_at` sadece başarıda güncellenir.

3. Failure/retry güvenliği:
   - Retry ve failed durumlarında `cursor_after` boş bırakılır.
   - Cursor koleksiyonu güncellenmez; önceki checkpoint korunur.

4. Strategy-config ile dinamik cursor çıkarımı:
   - `strategy_config.cursor_path` desteklenir (örn: `result.next_cursor`, `result.delta_token`).
   - `strategy_config.cursor_field` desteklenir.
   - Varsayılanlar strategy bazlıdır:
     - `cursor` -> `next_cursor`
     - `delta` -> `delta_token`
     - `file_watermark` -> `watermark`
     - `time_window` -> `window_end`

5. Run observability geliştirmesi:
   - `run.cursor_before` / `run.cursor_after` alanları doldurulur.
   - Cursor çıkarımı başarısızsa run’a warning (`cursor_commit_value_missing ...`) yazılır.

## Validation
- `go test ./internal/modules/integrations/... -run TestDoesNotExist` -> PASS
- `go test ./cmd/api/... -run TestDoesNotExist` -> PASS

## Next
1. Phase-3 başlangıcı: raw staging + canonical mapping/upsert pipeline.
