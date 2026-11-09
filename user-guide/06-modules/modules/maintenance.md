# Bakım Modu (Maintenance Mode)

**Sürüm:** 2.0.0  
**Son Güncelleme:** 2026-03-31  
**Modül Yolu:** `internal/platform/middleware/maintenance.go` + `internal/modules/system_ops/service/service.go`

---

## Genel Bakış

Bakım modu, planlanmış bakım, güncelleme veya kritik müdahale sırasında genel kullanıcı trafiğini geçici olarak durdurmak için kullanılır. Yönetim paneli ve CLI araçları bu mod etkinken bile her zaman erişilebilir kalır.

### Mimari

Bakım modu **iki katmanlı** bir kontrol mekanizması kullanır:

```
HTTP İsteği
    │
    ▼
/admin/* veya /health* → Tam geçiş (yönetici + probe'lar her zaman çalışır)
    │
    ▼
Redis: GET system:maintenance_mode == "true"?
    ├── Evet → 503 Service Unavailable + Retry-After: 120
    │
    ▼
Settings: maintenance.enabled == true?
    ├── Evet → 503 Service Unavailable + Retry-After: 120
    │
    ▼
Normal akış devam eder
```

**Neden çift katman?**
- **Redis katmanı (birincil):** Bakım modu anında etkilidir. Servis yeniden başlatılmasına gerek yoktur.
- **Settings katmanı (yedek):** Redis erişilemeyen durumlarda (örn. Redis arızası) yedek olarak devreye girer.

---

## İzinler

| İzin | Rol |
|------|-----|
| `system_ops:maintenance` | Super Admin, Admin, Operator |
| `system_ops:read` | Tüm roller (görüntüleme) |

---

## Engellenen / Engellenmeyenler

Bakım modu etkinken **ENGELLENİR:**
- Tüm genel API istekleri

Bakım modu etkinken **HER ZAMAN ERİŞİLEBİLİR:**
- `/admin/*` — Yöneticiler modu kapatabilsin
- `/health` — Canlılık kontrolü
- `/health/live` — Kubernetes liveness probe
- `/health/ready` — Kubernetes readiness probe

---

## API Endpoint'leri

### Bakım Modunu Aç

```http
POST /admin/system/maintenance/enable
Authorization: Bearer {token}
Content-Type: application/json

{
  "message": {
    "tr": "Sistem güncellemesi yapılmaktadır.",
    "en": "System maintenance in progress."
  }
}
```

Yanıt:
```json
{ "success": true, "data": { "enabled": true } }
```

### Bakım Modunu Kapat

```http
POST /admin/system/maintenance/disable
Authorization: Bearer {token}
```

### Bakım Modu Durumu

```http
GET /admin/system/maintenance
Authorization: Bearer {token}
```

Yanıt:
```json
{
  "success": true,
  "data": {
    "enabled": true,
    "message": { "tr": "Güncelleme devam ediyor." },
    "enabled_by": "admin@sirket.com",
    "enabled_at": "2026-03-31T10:00:00Z"
  }
}
```

---

## Bakım Modu Yanıtı (Kullanıcıya)

Bakım modu etkinken engellenen isteklere dönen yanıt:

```http
HTTP/1.1 503 Service Unavailable
Retry-After: 120
Content-Type: application/json

{
  "success": false,
  "error": {
    "code": "maintenance_mode",
    "message": "Sistem bakım modunda. Lütfen daha sonra tekrar deneyin."
  }
}
```

---

## CLI ile Kullanım

```bash
aystek-cli maintenance on
aystek-cli maintenance off
aystek-cli maintenance status
```

---

## Güncelleme Sırasında Otomatik Bakım Modu

`upgrade.sh` ve `StartUpdate()` akışı bakım modunu **otomatik** yönetir:
- Adım 6: Bakım modu açılır
- Adım 13: Güncelleme başarıyla tamamlanınca bakım modu kapatılır

---

## Sorun Giderme

**Bakım modu kapanmıyor:**
```bash
redis-cli DEL system:maintenance_mode
curl -X POST http://127.0.0.1:9090/admin/system/maintenance/disable \
  -H "Authorization: Bearer ADMIN_TOKEN"
```
