# AYSTEK Core Release Env And Installer Contract
## Slice 9 Checkpoint: Namespace-Aware Release Notes

Date: 2026-04-24

### Scope
- release workflow env sözleşmesini GHCR namespace esnekliği ile hizalamak
- GitHub release body içindeki örnek `docker pull` ve installer tag satırlarını gerçek workflow env ile uyumlu hale getirmek

### Completed
- `aystek-core` release body artık sabit `ghcr.io/aystek/...` yerine workflow env içindeki `CORE_IMAGE` ve `PANEL_IMAGE` değerlerini kullanıyor:
  - `/workspace/go/aystek-core/.github/workflows/release.yml`
- `aystek-core-next` release body de aynı şekilde namespace-aware hale geldi:
  - `/workspace/node/aystek-core-next/.github/workflows/release.yml`

### Result
- Private GHCR namespace repo owner veya `GHCR_NAMESPACE` override ile değiştiğinde release notları yanlış image yolu göstermeyecek.
- Workflow davranışı ile operatör dokümantasyonu aynı sözleşmeyi paylaşır hale geldi.
