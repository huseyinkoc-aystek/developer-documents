---
Title: API Client Onboarding & Best Practices
Version: 1.0.0
Status: Active
Last Updated: 2026-03-09
Related Documents:
  - ./INTEGRATION_GUIDE.md
  - ./ERROR_CODE_REFERENCE.md
  - ../02-security/AUTHENTICATION_AND_AUTHORIZATION.md
  - ../05-operations/API_GATEWAY_SECRET_ROTATION_RUNBOOK.md
---

# API Client Onboarding & Best Practices

Bu dokuman, third-party partnerlarin AYSTEK API Gateway'e guvenli sekilde onboarding surecini tanimlar.

## 1. Onboarding Steps
1. API client olustur:
   - `POST /admin/api-clients`
2. Scope ve route whitelist tanimla:
   - minimum yetki prensibi (least privilege)
3. Partner'a `client_id` ve ilk `client_secret` guvenli kanaldan ilet.
4. Token akisini dogrula:
   - `POST /auth/oauth2/token`
5. Pilot trafik ile izleme dogrulamasi yap:
   - audit logs + metrics

## 2. Recommended Initial Scope Set
1. Read-only entegrasyonlar icin:
   - `read:products`
   - `read:orders`
2. Yazma ihtiyaci olanlar icin ek:
   - `write:orders`

Not: Her partner icin gerekli olmayan scope verilmemelidir.

## 3. Token Usage Best Practices
1. Access token cache edin, gereksiz token cagrisi yapmayin.
2. Token expiry oncesi kontrollu refresh stratejisi kullanin.
3. Token ve secret degerlerini loglamayin.
4. Error handling icin 401/403/429 senaryolarini ayri ele alin.

## 4. Request Quality Best Practices
1. `X-Request-ID` gonderin (traceability).
2. Idempotent yazma operasyonlarinda unique islem anahtari stratejisi kullanin.
3. Retry yaparken exponential backoff uygulayin.
4. 429 durumunda cooldown uygulayip hizli tekrar denemeyin.

## 5. Secret Rotation Policy
1. Standart rotation periyodu: 90 gun.
2. Olay sonrasi zorunlu rotation: aninda.
3. Rotation endpoint:
   - `POST /admin/api-clients/:id/rotate-secret`
4. Rotation sonrasi token alma smoke testi zorunludur.

## 6. Audit and Support
1. API client audit izleme:
   - `GET /admin/api-clients/:id/audit-logs`
2. Export:
   - `GET /admin/api-clients/:id/audit-logs/export?format=csv|json`
3. Support eskalasyonunda `request_id` + zaman penceresi paylasin.

## 7. Go-Live Checklist
1. Scope/route whitelist onaylandi.
2. Rate-limit ve daily quota partner hacmine gore ayarlandi.
3. Secret rotation takvimi tanimlandi.
4. Monitoring alert abonelikleri aktif.
5. Incident contact listesi karsilikli dogrulandi.
