---
Title: Fleet Tally Module Reference
Version: 1.0.0
Status: Production Reference
Last Updated: 2026-06-05
Classification: Internal - Module Reference
Related Documents:
  - ../FLEET_TALLY_USER_GUIDE_TR.md
  - ../INTEGRATIONS_PLATFORM_GUIDE.md
---

# GPS Çetele Kontrol — fleet_tally

`fleet_tally` modülü, dış sistemlerden alınan servis/çetele kayıtlarını GPS konum geçmişiyle doğrular ve sonucu dış sisteme geri gönderir.

## Kullanıcı Kılavuzu

Detaylı Türkçe kullanım ve operasyon rehberi:

- `../FLEET_TALLY_USER_GUIDE_TR.md`

## Alt Modüller

| Alt Modül | Amaç |
|---|---|
| `control_types` | Doğrulama algoritması ve şema tanımları |
| `gps_providers` | Trio, N2Mobil, Mobiliz, Arvento sağlayıcı konfigürasyonları |
| `gps_providers/vehicle_mappings` | Plaka ↔ sağlayıcı araç/cihaz ID eşlemeleri |
| `connectors` | SAP SOAP veya REST dış sistem bağlayıcıları |
| `schedules` | Cron tabanlı veya manuel doğrulama çalıştırmaları |
| `verifications` | Job/run kayıtları, sonuçlar ve retry işlemleri |

## Admin Ekranları

| Ekran | Yol |
|---|---|
| Kontrol Tipleri | `/admin/fleet-tally/control-types` |
| GPS Sağlayıcılar | `/admin/fleet-tally/gps-providers` |
| Veri Bağlayıcılar | `/admin/fleet-tally/connectors` |
| Zamanlayıcılar | `/admin/fleet-tally/schedules` |
| Doğrulama İşleri | `/admin/fleet-tally/verifications` |
| Çalışma Geçmişi | `/admin/fleet-tally/runs` |
| Araç Eşlemeleri | `/admin/fleet-tally/vehicle-mappings` |
| Kullanım Kılavuzu | `/admin/fleet-tally/user-guide` |

## Temel API Yüzeyi

| Metod | Yol | Açıklama |
|---|---|---|
| `GET` | `/admin/fleet-tally/control-types` | Kontrol tiplerini listeler |
| `GET` | `/admin/fleet-tally/gps-providers` | GPS sağlayıcılarını listeler |
| `GET` | `/admin/fleet-tally/connectors` | Dış sistem bağlayıcılarını listeler |
| `GET` | `/admin/fleet-tally/schedules` | Zamanlayıcıları listeler |
| `POST` | `/admin/fleet-tally/schedules/:id/run` | Zamanlayıcıyı manuel çalıştırır |
| `GET` | `/admin/fleet-tally/verifications/jobs` | Doğrulama işlerini listeler |
| `POST` | `/admin/fleet-tally/verifications/jobs/:id/retry` | Bir işi tekrar dener |
| `GET` | `/admin/fleet-tally/verifications/runs` | Çalıştırma geçmişini listeler |
| `GET` | `/admin/fleet-tally/verifications/runs/:id/stream` | Run durumunu SSE ile izler |

## Yetkiler

| Yetki Modülü | Kapsam |
|---|---|
| `fleet_tally_control_types` | Kontrol tipi CRUD |
| `fleet_tally_gps_providers` | GPS sağlayıcı CRUD ve konum sorgusu |
| `fleet_tally_connectors` | Dış sistem bağlayıcı CRUD ve iş çekme |
| `fleet_tally_schedules` | Zamanlayıcı CRUD ve manuel run |
| `fleet_tally_verifications` | Job/run görüntüleme ve retry |

## Çok-Tenant / Site Davranışı

Tüm kayıtlar `tenant_id` ve mümkün olduğunda `site_id` ile izole edilir. Repository sorguları tenant/site context olmadan çalışmamalıdır.

Gürsel Araç Takip sitesi (`aractakip.gurseltur.com.tr`) için görünüm profili seedleri fleet-tally ekranlarını görünür tutar.

