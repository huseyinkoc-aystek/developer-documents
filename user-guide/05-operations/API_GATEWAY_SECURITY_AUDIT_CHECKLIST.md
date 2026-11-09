---
Title: API Gateway Security Audit Checklist
Version: 1.0.0
Status: Active
Last Updated: 2026-03-09
Related Documents:
  - ./API_GATEWAY_SECRET_ROTATION_RUNBOOK.md
  - ./API_GATEWAY_SLO_SLA_ALERTING.md
  - ../03-compliance/API_GATEWAY_KVKK_GDPR_COMPLIANCE_MATRIX.md
---

# API Gateway Security Audit Checklist

Bu checklist, Phase 5 kapsamindaki SSRF, injection ve replay riskleri icin operasyonel denetim adimlarini standartlastirir.

## 1. Preconditions
1. Production ortaminda `OAUTH2_SECRET_PROVIDER=aws_kms` set edilmeli.
2. `OAUTH2_ALLOW_INLINE_SECRET_MANAGER` production'da false/empty olmali.
3. API gateway metrics endpoint (`/metrics`) scrape ediliyor olmali.

## 2. SSRF Controls
1. Downstream cagrilar allowlist disina cikmiyor mu kontrol et.
2. IP literal, localhost, RFC1918 ve metadata endpoint (`169.254.169.254`) engellemeleri aktif mi dogrula.
3. URL parsing bypass testleri uygula:
   - mixed encoding
   - userinfo (`user@host`)
   - redirect chain

Pass criteria:
1. Dis hedeflere cikis red.
2. Denemeler audit/event loglarda izlenebilir.

## 3. Injection Controls
1. SQL/NoSQL benzeri payload'larda route auth + validation davranisini test et.
2. Header/query/body alanlarinda kontrol karakterleri normalize ediliyor mu denetle.
3. CSV export dosyalarinda formula injection korumasi dogrula (`=`, `+`, `-`, `@` prefixed values escaped).

Pass criteria:
1. Injection payloadlari execute olmamali.
2. Request reddi uygun error kodu ile donmeli.

## 4. Replay & Token Abuse Controls
1. Kisa omurlu access token sureleri politikaya uygun mu kontrol et.
2. Rotate-secret sonrasi eski secret ile token alma basarisiz olmali.
3. Repeated request patternlerinde rate-limit/quota deny metrikleri artmali.

Pass criteria:
1. Expired/invalid token reddedilmeli.
2. Replay denemeleri audit + metrics ile gorulebilmeli.

## 5. Audit and Evidence
1. Ilgili API client icin export alin:
   - `GET /admin/api-clients/:id/audit-logs/export?format=csv&limit=5000`
2. Asagidaki alanlarin kayitlandigini kontrol et:
   - `method`, `path`, `response_status`, `latency_ms`, `client_ip`, `created_at`
3. Security test penceresi icin hata/deny oranlarini dokumante et.

## 6. Sign-off Template
```text
Audit Window:
Environment:
Auditor:

SSRF: PASS/FAIL
Injection: PASS/FAIL
Replay: PASS/FAIL

Critical Findings:
Mitigations Applied:
Follow-up Date:
```
