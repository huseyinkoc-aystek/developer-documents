# TASK_054 - Integrations Enterprise Work Plan

## Scope
Integrations modulu (Go backend + Next.js admin) icin:
1. Test URL tabanli test-run kabiliyetinin guclendirilmesi
2. Mapping + enrichment pipeline'in enterprise seviyede devreye alinmasi
3. Trigger mekanizmasinin aynen korunmasi
4. Gereksiz/tekrarlayan kisimlarin temizlenmesi

## Mandatory Constraints
- Hook trigger semantics degismeyecek.
- Enrichment fail oldugunda trigger durmayacak; field default/empty ile devam edecek.
- Runtime retry/DLQ akisina dokunulmadan payload-building iyilestirilecek.
- Test-run kayitlari traceable olacak ve test oldugu acik isaretlenecek.

## Current-State Findings (Code)
- `InsertAction` create akisi `request_mapping` kaydetmiyor.
- Frontend `test_values` gonderiyor ama backend modelinde yok.
- `ValidateMapping` mock donuyor, gercek evaluate yok.
- `TestHook` pratikte sadece `action_id` uzerinden calisiyor; `action_keys` kapsamli calismiyor.
- `DirectExecuteAction` test cagrilarini execution olarak kaydetmiyor.
- `response_body` tip uyumu frontend parser'da dar (string bekliyor).

## Target Architecture
- Shared `PayloadBuilder`:
  - applyMapping
  - applyTransforms
  - applyEnrichment
  - applyDefaults
- Runtime (real trigger) ve test-run ayni payload pipeline'ini kullanacak.
- Test-run sonucunda execution benzeri kayit olusturulacak (`execution_kind=test`).

## Data Mapping Spec (Target)
- required_fields
- mapping: external_field -> internal path | expression | constant
- defaults: missing field fallback
- enrichment_rules: secondary source lookup by join key
- transforms: trim/upper/date-format vb.

## Implementation Phases

### Phase 1 - Baseline Stabilization (P0)
- [x] `InsertAction` -> `request_mapping` persist fix
- [x] `RequestMapping` -> `test_values` field add
- [x] Frontend type/parser -> `test_values` parse fix
- [x] Auth type normalization (`apikey`/`api_key`)

### Phase 2 - Shared Payload Engine (P0)
- [x] `payload_builder` service eklendi (`executor/payload_builder.go`)
- [x] mapping + defaults + transform + enrichment steps devrede
- [x] repository-backed enrichment resolver eklendi (secondary source lookup)
- [x] `ValidateMapping` endpoint'i expression + payload preview uretecek sekilde guncellendi
- [x] missing-field warning'ler request/execution sonucuna tasindi (`payload_warnings`, `_payload_warnings`) ve metrics/tracing sinyali eklendi

### Phase 3 - Test-Run Traceability (P0)
- [x] Execution modeline test/runtime ayrimi
- [x] `DirectExecuteAction` sonucunu execution kaydi olarak persist
- [x] Admin execution listesinde test filter/etiket

### Phase 4 - Hook Multi-Action Alignment (P1)
- [x] `TestHook` icin `action_keys` cozumleme
- [x] legacy `action_id` fallback korunacak

### Phase 5 - UI/UX Consolidation (P1)
- [x] ActionForm'dan `Test Input Data` alanini kaldir
- [x] ActionTestPanel icine `Test Dataset` bolumu tasi
- [x] Payload preview + response preview netlestir

### Phase 6 - Validation & Cleanup (P1/P2)
- [x] Gercek mapping validation
- [x] Placeholder/gereksiz kod temizligi
- [x] Dokuman ve test guclendirme

## Checkpoints
1. Create/Edit action mapping verisi kaybolmadan saklaniyor.
2. Test paneli disardan degil kendi bolumunden test datasini yonetiyor.
3. Runtime trigger davranisi degismiyor.
4. Missing enrichment durumunda akim fail etmiyor.
5. Test-run sonuclari execution tarafinda izlenebiliyor.

## Risk Controls
- Tenant/site context propagation korunacak.
- Default value policy field-level uygulanacak (`""` veya `null`).
- Test URL ile production URL ayrimi UI'da gorunur olacak.
- DLQ davranisi sadece runtime cagrilar icin devam edecek.

## First Execution Batch (Now)
- [x] Plan dokumani kaydedildi.
- [x] Backend request_mapping persistence + test_values support
- [x] Next UI: Test Input Data panel relocation
- [x] Hedefli dogrulama ve ara rapor

## Phase-2 Progress Update
- [x] Action create `request_mapping` persistence fix tamamlandi.
- [x] `test_values` backend model + frontend type/parser zinciri tamamlandi.
- [x] Test Input Data ActionTestPanel icine tasindi.
- [x] Shared payload builder runtime ve test-run body build akisina baglandi.
- [x] Enrichment resolver repository-backed implemente edildi.
- [x] ValidateMapping artik gercek expression evaluation + payload preview donuyor.
- [x] Next client tarafina `validateMapping` API methodu eklendi (payload preview/warnings alinabilir).
- [x] ActionTestPanel canli `validateMapping` cagrisi ile payload preview + warning karti gosteriyor.
- [x] Execution listesi artik `execution_type` (runtime/test) filtre ve etiket gosteriyor.
- [x] Direct action test cagrilari execution kaydi olarak persist ediliyor (`execution_type=test`).
- [x] TestHook `action_keys` cozumleyip `action_id` fallback ile calisiyor.
- [x] Integrations service testleri guclendirildi (DLQ, test execution persist, auth normalize, action_keys kapsami) ve paket testleri yesil.
