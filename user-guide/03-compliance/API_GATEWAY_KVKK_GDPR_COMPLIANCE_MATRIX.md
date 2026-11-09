---
Title: API Gateway KVKK & GDPR Compliance Matrix
Version: 1.0.0
Status: Active
Last Updated: 2026-03-09
Related Documents:
  - ./Regulatory_Compliance_Technical_Annex.md
  - ../04-api/API_CLIENT_ONBOARDING_BEST_PRACTICES.md
  - ../05-operations/API_GATEWAY_SECRET_ROTATION_RUNBOOK.md
  - ../05-operations/API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md
---

# API Gateway KVKK & GDPR Compliance Matrix

Bu dokuman, third-party API client erisimi icin KVKK/GDPR kontrol map'ini ve uygulama kanitlarini tanimlar.

## 1. Scope
1. OAuth2 `client_credentials` ile erisilen API gateway endpointleri.
2. `api_client` kimlikleri icin audit, retention ve least-privilege kontrolleri.
3. PII iceren endpointlerin siniflandirilmasi.

## 2. Data Classification Matrix
| Endpoint | Data Class | PII | KVKK/GDPR Note |
|---|---|---|---|
| `GET /products` | Product catalog | No | Genel urun verisi, low risk |
| `GET /products/:id` | Product detail | No | Genel urun verisi, low risk |
| `GET /orders` | Order summary | Potential | Musteri referansi/metadata varsa masked donus onerilir |
| `GET /orders/:id` | Order detail | Yes (conditional) | Kimlik/iletisim alanlari varsa data minimization zorunlu |
| `POST /orders` | Transactional write | Yes (conditional) | Purpose limitation + retention + audit zorunlu |

## 3. Control Mapping (KVKK/GDPR)
| Control Area | Requirement | Implementation |
|---|---|---|
| Identity & Access | Yetkisiz erisim engeli | OAuth2 client credentials + scope/route whitelist |
| Data Minimization | Gerekli olmayan PII paylasilmamali | Route bazli scope sinirlama + response filtering |
| Auditability | Islem izi zorunlu | `api_audit_logs` + export endpoint |
| Retention | Sinirli saklama suresi | TTL index, `API_GATEWAY_AUDIT_RETENTION_DAYS` |
| Key Management | Secret guvenligi | `aws_kms` provider + rotation runbook |
| Tenant Isolation | Veri ayrimi | Tenant/site context enforcement |
| Incident Response | Ihlal tespiti ve kayit | Security checklist + alerting + runbook |

## 4. Third-Party Data Processor Clause Template (Short Form)
Asagidaki metin, partner sozlesmesinde minimum clause olarak eklenmelidir.

```text
Data Processing Scope:
Processor receives API data only for documented service purpose and cannot reuse data outside controller instructions.

Security:
Processor enforces encryption in transit, credential protection, and role-based access controls.

Confidentiality:
Processor personnel with access to personal data are bound by confidentiality obligations.

Sub-processors:
Processor cannot appoint sub-processors without prior written approval.

Breach Notification:
Processor notifies controller without undue delay and no later than 24 hours after detection.

Retention & Deletion:
Processor deletes or returns personal data at contract termination unless legal retention applies.
```

## 5. Evidence Checklist
1. API client scope/route definitions exported from admin panel.
2. Rotation kayitlari (`POST /admin/api-clients/:id/rotate-secret`).
3. Audit export dosyalari (`/admin/api-clients/:id/audit-logs/export`).
4. Alert firing history (error-rate, deny burst, auth failures).
5. Partner sozlesmesinde data processor clause revizyon kaydi.

## 6. Review Cadence
1. Aylik: PII route/classification gozden gecirme.
2. Ceyreklik: Partner sozlesmesi + data processor clause audit.
3. Olay sonrasi: Ek rota/scope kisitlama ve zorunlu secret rotation.
