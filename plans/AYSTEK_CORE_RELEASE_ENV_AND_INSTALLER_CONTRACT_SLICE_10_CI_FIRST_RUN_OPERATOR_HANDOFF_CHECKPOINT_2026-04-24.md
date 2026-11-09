# AYSTEK Core Release Env And Installer Contract
## Slice 10 Checkpoint: CI First-Run Operator Handoff

Date: 2026-04-24

### Scope
- GitHub UI üzerinden yapılacak ilk gerçek CI koşuları için operatör handoff dökümanı hazırlamak
- dual-repo release runbook içine `Validated in CI` alanını eklemek

### Completed
- İlk gerçek GitHub UI koşuları için ayrı operatör runbook hazırlandı:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_CI_FIRST_RUN_VALIDATION_RUNBOOK_2026-04-24.md`
- Dual-repo release runbook artık bu operatör runbook'a referans veriyor ve ilk gerçek CI başarılarından sonra doldurulacak `CI Validation Status` alanı içeriyor:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Constraints
- Mevcut çalışma ortamında GitHub Actions `workflow_dispatch` tetikleme aracı yok.
- Yerelde `gh` CLI kurulu değil.
- Bu yüzden ilk gerçek run operatör tarafından GitHub UI üstünden başlatılmalıdır.

### Handoff
Operator should run, in order:
1. `huseyinkoc-aystek/aystek-core-installer` → `Installer P0 Disposable Matrix`
2. `huseyinkoc-aystek/admin-panel-next` → `Release - Panel Docker Image`
3. `huseyinkoc-aystek/aystek-core` → `Release & Package`

After that, bring back:
- workflow URLs
- success/failure status
- final namespace used
- release image coordinates

Then the runbooks can be updated with `Validated in CI`.
