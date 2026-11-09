---
Title: Enterprise Security Architecture Booklet (TR/EN)
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
  - ./TENANT_SITE_SECURITY_MATRIX.md
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md
  - ../archive/internal/MIDDLEWARE_CHAIN_INVARIANTS.md

---
# Table of Contents
1. [Türkçe](#türkçe)
2. [English](#english)

## Türkçe

### 1. Güvenlik Prensipleri
- Zero-trust: hiçbir bağlam varsayılan güvenli değildir.
- Defense-in-depth: tek katman değil çok katmanlı kontrol.
- Deterministik karar: ihlal koşullarında sabit hata modeli.

### 2. Güvenlik Katmanları
- HTTP sınırı: tenant/site çözümleme ve zorunluluk kontrolü.
- Repository sınırı: tenant dominance kanıtı.
- Cache sınırı: tenant/site partition.
- Event/worker sınırı: replay + state-machine guard.
- CI sınırı: güvenlik gate’leri.

### 3. Sonuç
Bu model, izolasyonu operasyonel süreçlere gömen ve denetlenebilir hale getiren bir kurumsal güvenlik mimarisi üretir.

## English

### 1. Security Principles
- Zero-trust posture for all incoming context.
- Defense-in-depth across runtime and delivery layers.
- Deterministic violation handling.

### 2. Security Layers
- HTTP boundary enforcement
- Repository dominance enforcement
- Cache partition isolation
- Event/worker replay and lifecycle guards
- CI security gates

### 3. Outcome
The architecture provides audit-ready, enterprise-grade isolation assurance with repeatable evidence.
