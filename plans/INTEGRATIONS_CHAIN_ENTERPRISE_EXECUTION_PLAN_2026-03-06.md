# Integrations Chain Enterprise Execution Plan

Date: 2026-03-06
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / user-guide
Source: Integrations modülü için çok adımlı entegrasyon akışları

## Objective
Integrations modülünü tek action tetikleme modelinden geriye uyumlu şekilde çok adımlı chain modeline genişletmek.

## Guardrails
1. Existing hook->action behavior must remain unchanged.
2. Trigger semantics must not change.
3. DSL must remain declarative (JMESPath/template), no arbitrary scripting.
4. Tenant/site isolation and existing module standards are mandatory.

## Delivery Phases

### Phase-1 (Start Now)
1. Backend domain model for `Chain` + `ChainNode` + `ChainEdge` (Integrations module).
2. CRUD APIs: `/admin/integrations/chains`.
3. Validate API: `/admin/integrations/chains/:id/validate` + payload shape validation helper.
4. DAG validation (single start, at least one end, no cycles, edge consistency).
5. Persistence: Mongo collection + tenant/site scoped repository methods.
6. i18n keys for validation/notice messages.

### Phase-2
1. Hook target expansion: `target_type=action|chain` and `target_ref` compatibility strategy.
2. Chain runner service with shared execution context (`event`, `vars`, `steps`).
3. Node execution for `action`, `transform`, `branch`.
4. Async worker integration (RabbitMQ path) + retry/DLQ semantics.

### Phase-3
1. Next.js admin pages for chains list/create/edit.
2. Graph editor integration (node/edge config panels).
3. Hook form target selector (Action/Chain).
4. Chain validate + test-run UI.

### Phase-4
1. Observability hardening (run-level + step-level traceability).
2. smoke/runbook updates and production playbook.
3. QA matrix: success/failure/branch/retry/DLQ/idempotency.

## Immediate Execution Checklist (Phase-1)
1. Add model + dto contracts.
2. Extend repository interface + mongo implementation.
3. Extend service interface + validation logic.
4. Add admin handler methods and routes.
5. Add translation seed keys.
6. Run `gofmt` and focused `go test`.

## Non-Goals for Phase-1
1. Full chain runtime execution.
2. Frontend graph editor.
3. Hook runtime binding to chain.

## Progress Status (Updated: 2026-03-06)

### Completed
1. Phase-1 tamamlandı:
   - Chain model + CRUD + validate endpoint + DAG validation + persistence + i18n seed.
2. Phase-2 büyük ölçüde tamamlandı:
   - Hook `target_type=action|chain` (default `action`, geriye uyumlu).
   - Chain runner (`start/action/transform/branch/end`) ve async worker entegrasyonu.
   - Execution/DLQ semantiği korunarak chain çalıştırma.
3. Phase-4/1 observability tamamlandı:
   - `step_execution` kayıt modeli.
   - Step execution list/detail API.
   - Next execution detail ekranında node timeline.
4. Phase-3 kısmi tamamlandı (bu iterasyon):
   - Next `chains` sayfaları eklendi: list/create/edit.
   - Hook form target selector (Action/Chain) canlı hale getirildi.
   - Hook listesine target type filtre/kolonları eklendi.
   - Görsel chain editor (node/edge panel + SVG flow preview) JSON ile çift yönlü senkron eklendi.
   - Enterprise editor iyileştirmesi: drag-drop node yerleşimi + node tipine özel config formları (action/transform/branch) + branch condition builder.
   - Backend ValidateChain iyileştirmesi: node tipine göre config schema doğrulaması + alan bazlı `validation_errors` çıktısı.
   - ValidateChain UI: `validation_errors` için filtrelenebilir tablo (node_id/field/code) ve satıra tıklayınca node odaklama.
   - Flow builder UX modernizasyonu: pan/zoom destekli canvas, sürükle-bırak taşınabilir düğümler, sağda seçili node konfigürasyon paneli, odaklanan node için otomatik input imleci.
   - Chain test-run UI: payload girişi + dry-run seçeneği + adım bazlı sonuç tablosu + final payload önizlemesi + node odak entegrasyonu.
   - Menü + çeviri seed kayıtları chains için güncellendi.
5. Phase-4/2-3 tamamlandı:
   - Runbook/smoke adımları `INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md` ile operasyonel olarak finalize edildi.
   - QA matrix (branch/retry/DLQ/idempotency) checklist senaryoları eklendi.
6. Phase-3/4 sonrası hardening tamamlandı:
   - Graph editor usability polish: undo/redo (Ctrl/Cmd+Z, Ctrl/Cmd+Shift+Z), klavye ile node taşıma, snap-to-grid, komşu node hizalama yardımcıları ve auto-layout eklendi.

### Remaining
1. Opsiyonel operasyon adımı:
   - GitLab kullanılacaksa protected environment + protected/masked variable (`INTEGRATIONS_SMOKE_TOKEN`) UI konfigürasyonunun canlı ortamda tamamlanması.
