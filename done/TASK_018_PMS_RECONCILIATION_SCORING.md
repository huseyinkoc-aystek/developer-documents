# TASK_018 — PMS Reconciliation Scoring v2

## Kapsam
- Match scoring algoritmasi deterministic hale getirildi.
- Candidate uretiminde V2 kriterleri:
  - description normalization
  - amount tolerance band (percent + absolute)
  - weighted scoring (amount + party_ref + unit_hint)

## Teknik Notlar
- `ProposeMatch` servisi `ProposeInput` alir:
  - `tolerance_percent`
  - `tolerance_abs`
- Adaylar skor ve baglayici tie-break ile siralanir.
- Ayni bank islemi icin cift confirm state machine ile engellenir (409).

## UI
- `/admin/pms/reconciliation` ekraninda candidate skor gorunur.
- Matched satirlarda manual candidate secilip confirm edilir.

## Test
- `internal/modules/pms/reconciliation/service/service_test.go`:
  - deterministic scoring
  - duplicate/invalid confirm guard
