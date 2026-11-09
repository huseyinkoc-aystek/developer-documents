---
Title: GitLab Protected Environment Setup (Integrations Smoke)
Version: 1.1.0
Status: Production
Last Updated: 2026-03-06
Scope: CI/CD Security Hardening
Related Documents:
  - ./INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ./INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
---

# GitLab Protected Environment Setup

Bu doküman, `integrations_smoke_production` job'ı için GitLab tarafında zorunlu güvenlik ayarlarını tanımlar.

## 1. Neden Gerekli?
`manual + production` job'ının sadece yetkili kişilerce ve protected branch üzerinde çalışması gerekir.

## 2. Protected Environment Ayarı
1. GitLab projesine girin.
2. `Settings -> CI/CD` bölümünü açın.
3. `Protected environments` alanına gidin.
4. `Protect an environment` seçin.
5. `Environment` olarak `production` yazın/seçin.
6. `Allowed to deploy` için sadece yetkili rol/grupları seçin (örn. Maintainer veya Release Managers).
7. Kaydedin.

## 3. Protected + Masked Variable Ayarı
1. `Settings -> CI/CD -> Variables` alanına gidin.
2. `INTEGRATIONS_SMOKE_TOKEN` değişkenini oluşturun veya düzenleyin.
3. Değeri production admin bearer token olacak şekilde girin.
4. `Masked` seçeneğini aktif edin.
5. `Protected` seçeneğini aktif edin.
6. Ortam kapsamı gerekiyorsa `production` olarak sınırlandırın.
7. Kaydedin.

## 4. Önerilen Ek Değişkenler
1. `INTEGRATIONS_SMOKE_ENABLED=true` (Protected)
2. `INTEGRATIONS_SMOKE_API_BASE_URL=https://<production-api>` (Protected)
3. `INTEGRATIONS_SMOKE_EVENT_KEY=<kritik_event_key>` (opsiyonel)

## 5. Kod Seviyesi Koruma (Uygulandı)
Pipeline kuralında production smoke job sadece aşağıdaki koşulda görünür/çalışır:

- `CI_COMMIT_REF_PROTECTED == "true"`
- `CI_COMMIT_BRANCH == CI_DEFAULT_BRANCH`
- `INTEGRATIONS_SMOKE_ENABLED == "true"`

Bu koşullar `.gitlab-ci.yml` içinde uygulanmıştır.

## 6. Doğrulama Checklist
1. Default branch protected değilse job görünmemeli.
2. Protected olmayan branch pipeline'ında production smoke job görünmemeli.
3. Protected branch + manual tetikleme ile job görünmeli.
4. Token eksikse job güvenli şekilde fail etmeli.
5. Job loglarında token değeri maskelenmiş görünmeli.

## 7. API ile Otomatik Doğrulama (Ops Script)
UI ayarlarının doğru yapıldığını API üzerinden doğrulamak için script:

`/workspace/go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`

### Gereken Değişkenler
1. `GITLAB_TOKEN` (api/read_api yetkili token)
2. `GITLAB_PROJECT_ID` (numeric project id)

### Opsiyonel
1. `GITLAB_API_URL` (default: `https://gitlab.com/api/v4`)
2. `REQUIRED_ENVIRONMENT` (default: `production`)
3. `REQUIRED_VARIABLE` (default: `INTEGRATIONS_SMOKE_TOKEN`)
4. `REQUIRED_SCOPE` (default: `production`)

### Örnek
```bash
cd /workspace/go/aystek-core
chmod +x tools/ops/gitlab_verify_integrations_smoke_security.sh

GITLAB_TOKEN="<token>" \
GITLAB_PROJECT_ID="12345678" \
REQUIRED_ENVIRONMENT="production" \
REQUIRED_VARIABLE="INTEGRATIONS_SMOKE_TOKEN" \
REQUIRED_SCOPE="production" \
tools/ops/gitlab_verify_integrations_smoke_security.sh
```

### Beklenen Sonuç
1. `production` protected environment mevcut olmalı.
2. Deploy access levels en az bir kayıt içermeli.
3. `INTEGRATIONS_SMOKE_TOKEN` değişkeni `masked=true` ve `protected=true` olmalı.
4. Scope `production` (veya en azından `*`) olmalı.
