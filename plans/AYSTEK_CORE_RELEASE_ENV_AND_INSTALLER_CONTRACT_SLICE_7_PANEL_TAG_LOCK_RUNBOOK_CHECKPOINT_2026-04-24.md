# AYSTEK Core Release Env And Installer Contract
## Slice 7 Checkpoint: Panel Tag Lock And Dual-Repo Runbook Hardening

Date: 2026-04-24

### Scope
- `aystek-core` release workflow içinde panel image tag lock doğrulaması eklemek
- canonical GHCR namespace örneklerini release body içinde düzeltmek
- dual-repo release runbook'u checklist seviyesinden tag-lock operasyon akışına taşımak

### Completed
- `aystek-core` release workflow artık matching panel tag olmadan devam etmiyor:
  - `/workspace/go/aystek-core/.github/workflows/release.yml`
- Core release body içindeki eski namespace örnekleri canonical hale getirildi:
  - `ghcr.io/aystek/aystek-core:<version>`
  - `ghcr.io/aystek/aystek-panel:<version>`
- Dual-repo runbook tag parity ve release order kurallarını daha görünür hale getirdi:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Result
- Core release artık panel image publish edilmeden operatöre yarım release çifti sunmayacak.
- Installer ve operasyon dokümantasyonu canonical GHCR sözleşmesiyle hizalı kaldı.
