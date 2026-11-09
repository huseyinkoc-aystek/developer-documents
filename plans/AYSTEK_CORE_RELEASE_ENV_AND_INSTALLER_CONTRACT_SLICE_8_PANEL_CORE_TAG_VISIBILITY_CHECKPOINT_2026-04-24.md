# AYSTEK Core Release Env And Installer Contract
## Slice 8 Checkpoint: Panel-Core Tag Visibility Coordination

Date: 2026-04-24

### Scope
- `aystek-core-next` release workflow içinde matching core tag görünürlüğünü raporlamak
- dual-repo release runbook'u iki yönlü ama deadlock üretmeyen koordinasyon modeliyle güncellemek

### Completed
- Panel release workflow artık matching core image tag'in GHCR görünürlüğünü kontrol ediyor:
  - `/workspace/node/aystek-core-next/.github/workflows/release.yml`
- Panel GitHub release body artık matching core tag görünürlüğünü raporluyor.
- Runbook, panel tarafının görünürlük raporu verdiğini; core tarafının ise asıl tag lock enforcement yaptığını açıkça tanımlıyor:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Result
- `aystek-core-next` artık release koordinasyonunda göz ardı edilen taraf değil.
- İki repo birlikte çalışan bir operasyon modeline daha net bağlandı:
  - panel release: visibility/reporting
  - core release: hard lock enforcement
