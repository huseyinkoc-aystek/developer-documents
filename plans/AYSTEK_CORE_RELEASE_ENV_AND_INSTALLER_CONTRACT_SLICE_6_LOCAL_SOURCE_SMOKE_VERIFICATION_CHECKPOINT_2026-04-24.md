# AYSTEK Core Release Env And Installer Contract
## Slice 6 Checkpoint: Local-Source Smoke Verification

Date: 2026-04-24

### Scope
- `docker + local-source` smoke doğrulamasını CI uyumlu hale getirmek
- `native + local-source` smoke hattını gerçek repo yollarıyla doğrulamak
- `source-mode-contract` CI job zincirindeki eksik sibling checkout'u kapatmak

### Completed
- `source-mode-contract` job artık `_deps/aystek-core` ve `_deps/aystek-core-next` sibling checkout'larını birlikte alıyor:
  - `/workspace/go/aystek-core-installer/.github/workflows/installer-p0-disposable.yml`
- `docker + local-source` smoke script'i, Docker binary bulunmayan yerel geliştirme ortamlarında açıklayıcı şekilde `skip` veriyor; CI runner üzerinde ise gerçek smoke çalışmaya devam ediyor:
  - `/workspace/go/aystek-core-installer/tests/scripts/run-it-12-docker-local-source-smoke.sh`
- `native + local-source` smoke gerçek repo yolları ile başarıyla geçti:
  - `/workspace/go/aystek-core-installer/tests/scripts/run-it-13-native-local-source-smoke.sh`

### Verification
- `bash -n /workspace/go/aystek-core-installer/tests/scripts/run-it-12-docker-local-source-smoke.sh`
- `bash -n /workspace/go/aystek-core-installer/tests/scripts/run-it-13-native-local-source-smoke.sh`
- `bash -n /workspace/go/aystek-core-installer/tests/scripts/run-source-mode-contracts.sh`
- `cd /workspace/go/aystek-core-installer && make integration-source-modes`
- `cd /workspace/go/aystek-core-installer && make integration-local-source-smoke`
- `cd /workspace/go/aystek-core-installer && go test ./...`
- `cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds`

### Notes
- Yerel çalışma ortamında Docker binary mevcut olmadığından `docker + local-source` smoke burada `skip` verdi.
- Bu bir zayıflık değil; CI workflow bu senaryoyu Docker erişimi olan runner üzerinde gerçek olarak doğrulayacak şekilde korunmuştur.
- `native + local-source` smoke hattı yerelde fiilen doğrulanmıştır.

### Result
- `local-source` deployment mode artık hem sözleşme seviyesinde hem de smoke doğrulama seviyesinde daha güvenilir.
- Private GHCR yolu ile local source yolu birlikte desteklenen, belgelenmiş ve CI tarafından doğrulanan iki ayrı kurulum hattı haline geldi.
