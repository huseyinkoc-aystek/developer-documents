---
Title: Integrations Chain Runbook & QA Matrix
Version: 1.0.0
Status: Production
Last Updated: 2026-03-06
Scope: All Modules (Hook target: chain)
Related Documents:
  - ./INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ../06-modules/INTEGRATIONS_PLATFORM_GUIDE.md
---

# Integrations Chain Runbook & QA Matrix

## 1. Binding Model (Hook -> Chain -> Action -> Connector)
Chain akışı aşağıdaki bağla çalışır:
1. Hook tetiklenir (`event_key` eşleşmesi).
2. Hook `target_type=chain` ise `hook.chain_id` yüklenir.
3. Chain graph node'ları sırayla yürür (`start -> ... -> end`).
4. `action` node içinde `config.action_id` ile Action bulunur.
5. Action içindeki `connector_id` ile Connector bulunur ve HTTP çağrısı yapılır.

Notlar:
1. Trigger semantiği değişmez; sadece hedef action yerine chain olur.
2. Bir chain içinde farklı action node'ları farklı connector'lara gidebilir.
3. Geriye uyumluluk için `target_type` verilmezse hook varsayılan `action` kabul edilir.

## 2. Phase-4 Smoke Runbook (Chain)

### 2.1 Preconditions
1. API healthy (`/healthz`, `/readyz`).
2. RabbitMQ healthy (async smoke için).
3. En az bir active chain ve chain target'lı active hook var.
4. Chain graph validate sonucu `valid=true`.
5. Chain içindeki action node'larda `config.action_id` doğru.

### 2.2 Functional Smoke (Admin UI)
1. Chain edit ekranı aç.
2. `Validate` çalıştır, `validation_errors` boş olmalı.
3. `Test Run` panelinde payload ile `dry_run=true` test et.
4. Step tablosunda node sırası ve `next_node` akışı doğrulanır.
5. `dry_run=false` test et, action node HTTP sonuçları doğrulanır.
6. Hata satırına tıklayıp node focus davranışı doğrulanır.

### 2.3 Async Runtime Smoke (EventBridge + Worker)
1. Hook `mode=async`, `target_type=chain` olacak şekilde kaydedilir.
2. Domain event tetiklenir.
3. `integration_events` kuyruğuna publish gözlenir.
4. Worker tüketimi log'da gözlenir.
5. Execution kaydı oluşmalı (`success` veya kontrollü `failed`).
6. Execution detail'de step timeline kayıtları görünmeli.

### 2.4 Retry / DLQ Drill
1. Chain içindeki bir action endpoint geçici bozulur.
2. Execution retry semantiği (`next_retry_at`) doğrulanır.
3. Retry limiti aşılınca DLQ kaydı doğrulanır.
4. Requeue sonrası tekrar çalıştırma gözlenir.

## 3. QA Matrix

| Case ID | Senaryo | Girdi | Beklenen |
|---|---|---|---|
| CH-001 | Validate başarılı | Geçerli DAG + node config | `valid=true`, `validation_errors=[]` |
| CH-002 | Action node config eksik | `action` node, `config.action_id` yok | `valid=false`, field=`config.action_id` |
| CH-003 | Transform expr hatalı | parse edilemeyen expression | `valid=false`, field=`config.expression` |
| CH-004 | Branch expr hatalı | parse edilemeyen predicate | `valid=false`, field=`config.predicate_expr` |
| CH-005 | Non-branch multi-edge | action/transform node'dan 2 edge | `valid=false`, outgoing edge hatası |
| CH-006 | Branch edge kuralları | duplicate true/default edge | `valid=false`, condition hatası |
| CH-007 | Dry run chain test | `dry_run=true` | Step akışı döner, DB execution yazımı beklenmez |
| CH-008 | Real chain test | `dry_run=false` | Step HTTP sonuçları ve final payload döner |
| CH-009 | Async event success | Hook async chain | Execution `success`, step timeline dolu |
| CH-010 | Async retry->DLQ | endpoint fail | retry ve ardından DLQ davranışı |
| CH-011 | Idempotency same key | aynı idempotency ile tekrar | duplicate çalışma engeli/uygun semantik |
| CH-012 | Multi-connector chain | farklı action node'lar farklı connector | Her step doğru connector ile çalışır |

## 4. Acceptance Criteria
1. Hook target chain akışında event tetiklemesi stabil.
2. Chain validation alan bazlı hata üretir.
3. Chain test-run sonucu adım bazlı izlenebilir.
4. Async execution + retry + DLQ semantiği bozulmamış.
5. UI focus/filtre ergonomisi kullanıcı dostu çalışır.

## 5. Operational Notes
1. Chain test-run endpoint preview amaçlıdır; persist etmez.
2. Canlı davranış için hook execute + worker hattı izlenmelidir.
3. Production smoke için protected token ve restricted environment kullanılmalıdır.
