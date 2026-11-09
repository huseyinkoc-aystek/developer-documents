---
Title: Regulatory & Compliance Technical Annex (TR/EN)
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
  - ../01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md
  - ../04-api/ERROR_CODE_REFERENCE.md
---

# Table of Contents
1. [Türkçe](#türkçe)
2. [English](#english)

## Türkçe

### 1. Uyum Kapsamı
Bu ek, içerik botları ve operasyon ayarları kapsamındaki kontrol modelinin kurumsal uyum çerçeveleriyle teknik eşlemesini sunar.

### 2. Kontrol Eşlemesi
- SOC 2: moderasyon aksiyonları, policy değişiklikleri ve kod enjeksiyon güncellemeleri için izlenebilir kontrol izi.
- ISO 27001: görev ayrılığı, RBAC ve değişiklik yönetimi.
- GDPR/KVKK: tenant/site bazlı veri ayrımı ve kaynak izlenebilirliği.

### 3. Denetim Kanıtları
- Moderasyon aksiyonları için audit kaydı (approve/reject/edit/schedule).
- Site policy, feature flags, limits ve ad policy değişiklikleri için denetim izi.
- Injection code create/update/delete işlemleri için sürüm ve aktör bilgisi.
- Outbox event kayıtları ile işlem sonrası izlenebilirlik.

### 4. İçerik İzlenebilirliği
Yayınlanan içerik, kaynak atfı (`source attribution`) ve moderation kararı ile izlenebilir olmalıdır.

### 5. Yönetişim Kontrolleri
- Editoryal sorumluluk eşlemesi (reviewer/approver actor mapping)
- Kritik değişikliklerde çoklu onay modeli (opsiyonel)
- Site bazlı operational governance ile merkezi kontrol + yerel özerklik dengesi

## English

### 1. Compliance Scope
This annex maps the content ingestion and operations settings control model to enterprise compliance expectations.

### 2. Control Mapping
- SOC 2: traceable controls for moderation actions, policy updates, and injection code changes.
- ISO 27001: segregation of duties, RBAC, and controlled change management.
- GDPR/KVKK: tenant/site data isolation and source-level traceability.

### 3. Audit Evidence
- Audit logs for moderation actions (approve/reject/edit/schedule).
- Audit trails for policy, flags, limits, and ad policy updates.
- Versioned audit records for injection code lifecycle operations.
- Outbox records for post-commit event traceability.

### 4. Content Traceability
Published content must retain source attribution and moderation decision lineage.

### 5. Governance Controls
- Editorial responsibility mapping (reviewer/approver identity)
- Optional multi-step approval for high-risk categories
- Site-scoped governance under central enterprise policy
