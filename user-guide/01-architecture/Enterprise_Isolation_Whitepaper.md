---
Title: Enterprise Isolation Whitepaper (TR/EN)
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ../07-executive/Board_Level_Risk_Isolation_Report.md
  - ../02-security/Enterprise_Security_Architecture_Booklet.md
  - ./Distributed_Zero_Trust_Tenant_Architecture_Spec.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
  - ../04-api/ERROR_CODE_REFERENCE.md

---
# Table of Contents
1. [Türkçe](#türkçe)
2. [English](#english)

## Türkçe

### 1. Kapsam
Bu doküman AYSTEK Core izolasyon modelinin kurumsal seviyede teknik çerçevesini sunar.

### 2. Mimari Özet
- Middleware tabanlı tenant/site çözümleme
- Repository dominance (AST denetimli)
- Cache tenant/site partition
- Event replay ve signature doğrulama
- Tenant yaşam döngüsü (active/inactive/tombstoned)

### 3. Deterministik Güvence
İhlal koşullarında sistem tutarlı ve izlenebilir sonuç üretir. Bu yaklaşım denetim ve operasyon süreçlerini güçlendirir.

### 4. Operasyonel Kanıt
- Unit test
- Security gate
- Integration matris testleri

## English

### 1. Scope
This document defines the enterprise technical model for AYSTEK Core tenant/site isolation.

### 2. Architecture Summary
- Middleware-based tenant/site resolution
- AST-audited repository dominance
- Tenant/site cache partitioning
- Event replay and signature verification
- Tenant lifecycle state model (active/inactive/tombstoned)

### 3. Deterministic Assurance
Violation paths produce consistent, traceable outcomes suitable for enterprise governance.

### 4. Operational Evidence
- Unit-level control validation
- Security gate enforcement
- Integration matrix verification
