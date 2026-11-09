---
Title: API Gateway Secret Rotation Runbook
Version: 1.0.0
Status: Active
Last Updated: 2026-03-09
Related Documents:
  - ./OPERATIONS_GUIDE.md
  - ../02-security/AUTHENTICATION_AND_AUTHORIZATION.md
  - ../../plans/AYSTEK_CORE_API_GATEWAY_THIRD_PARTY_ACCESS_CONTROL.md
---

# API Gateway Secret Rotation Runbook

Bu runbook, OAuth2 `api_client` secret rotation islemlerinin guvenli ve izlenebilir sekilde uygulanmasi icindir.

## 1. Kapsam
1. API client secret rotation sureci.
2. Rotation sonrasi token alma dogrulamasi.
3. Audit log export ile degisiklik sonrasi izleme.
4. Hata durumunda geri donus adimlari.

## 2. On Kosullar
1. `auth` modulunde admin yetkisi aktif olmali.
2. Production ortaminda `OAUTH2_SECRET_PROVIDER` acikca set edilmeli:
   - onerilen: `aws_kms`
3. Production ortaminda `inline` provider varsayilan olarak kapali olmali.
4. API gateway audit indeksleri olusmus olmali (`api_audit_logs`).

## 3. Rotation Is Akisi
### 3.1 Rotation Tetikleme
1. Hedef API client ID belirlenir.
2. Endpoint cagrilir:
   - `POST /admin/api-clients/:id/rotate-secret`
3. Donen `client_secret` yalnizca bir kez guvenli kanaldan partner sisteme aktarilir.

Beklenen:
1. `200 OK`
2. Response icinde `client_secret` ve `secret_expires_at` doner.

### 3.2 Partner Cutover Dogrulamasi
1. Partner tarafi yeni secret ile token endpoint'i cagirir:
   - `POST /auth/oauth2/token`
2. `grant_type=client_credentials` ile yeni access token alinabildigi dogrulanir.

Beklenen:
1. Token uretimi basarili.
2. Scope/route kontrolleri degismeden calisir.

### 3.3 Audit Dogrulamasi
1. Rotation sonrasi ilk 15 dakika boyunca API kullanimini izleyin.
2. Audit export alin:
   - `GET /admin/api-clients/:id/audit-logs/export?format=csv&limit=1000`
3. Gerekirse filtre ekleyin:
   - `from_date`, `to_date`, `status`, `method`, `path`

Beklenen:
1. 401/403 anomali pikleri yok.
2. Basarili trafik trendi stabil.

## 4. Retention Politikasi
1. Audit retention TTL env:
   - `API_GATEWAY_AUDIT_RETENTION_DAYS`
2. Varsayilan deger: `90` gun.
3. TTL indeksi: `idx_created_at_ttl` (`created_at`).

Not:
1. Regulator gereksinimine gore retention gunu ortam bazli ayarlanabilir.
2. Uzun sureli arsiv icin export dosyalari ayrica saklanmalidir.

## 5. Hata ve Geri Donus
1. Rotation sonrasi tum token talepleri basarisiz ise:
   - Partner konfigrasyonunu ve secret dagitim kanalini dogrulayin.
2. KMS kaynakli hata varsa:
   - KMS key policy, IAM rol, region ve key id degerlerini kontrol edin.
3. Gecici geri donus:
   - Yeni secret partnera ulasana kadar trafik penceresi kontrollu daraltilir.
4. Kalici cozum:
   - Yeni rotation islemi planli pencerede tekrar uygulanir.

## 6. Operasyon Siklik Onerisi
1. Standart rotation periyodu: 90 gun.
2. Yuksek riskli partnerler: 30-60 gun.
3. Olay sonrasi zorunlu rotation: aninda.
