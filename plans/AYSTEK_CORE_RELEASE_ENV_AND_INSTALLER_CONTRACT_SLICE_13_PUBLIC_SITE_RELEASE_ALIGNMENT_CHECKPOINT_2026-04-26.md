# AYSTEK Core Release Env And Installer Contract
## Slice 13 Checkpoint: Public Site Release Alignment

Date: 2026-04-26

### Scope
- `aysteknoloji-com` repo için eksik release hattını açmak
- public site env/Docker sözleşmesini canonical hale getirmek
- dört repo release runbook'unu uyumlu hale getirmek

### Completed
- `aysteknoloji-com` için Dockerfile eklendi:
  - `/workspace/node/aysteknoloji-com/Dockerfile`
- public site canonical env example eklendi:
  - `/workspace/node/aysteknoloji-com/.env.example`
- public site release workflow eklendi:
  - `/workspace/node/aysteknoloji-com/.github/workflows/release.yml`
- public site repo README release/env sözleşmesi ile güncellendi:
  - `/workspace/node/aysteknoloji-com/README.md`
- CI validation runbook ve dual-repo runbook dört repo modeline genişletildi:
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_CI_FIRST_RUN_VALIDATION_RUNBOOK_2026-04-24.md`
  - `/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_AND_PANEL_DUAL_REPO_RELEASE_RUNBOOK_2026-04-24.md`

### Result
- `aysteknoloji-com` artık release/install konusunun dışarıda kalan repo'su değil.
- Core, panel, installer ve public site için daha tutarlı bir release/env sözleşmesi oluştu.
