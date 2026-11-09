# AYSTEK AI için Ollama Kurulum Rehberi

**Tarih:** 2026-10-10
**Durum:** Kod tarafı hazır ve commit edildi (aystek-core `3cef8e3`, Faz 9a).
Bu belge SADECE sunucu tarafı kurulumu anlatır — bu adımlar üretim
sunucusunda (SSH erişimi olan biri tarafından) çalıştırılmalı, bu belgeyi
yazan AI oturumu üretim sunucusuna erişemiyor (platform kısıtı).

**İlişkili:** `KISISEL_ASISTAN_YONETIM_YUZEYI_VE_UCRETSIZ_AI_YOL_HARITASI_2026-10-08.md`
(Faz 9), `SUNUCU_KURULUM_REHBERI.md`, `DEPLOYMENT_SCENARIOS.md`.

## Ön Koşul — Kod Zaten Hazır

`internal/modules/personal_assistant/service/engine.go` artık iki ortam
değişkeni okuyor:
- `OLLAMA_BASE_URL` — örn. `http://ollama:11434` (Docker) veya
  `http://127.0.0.1:11434` (systemd/native).
- `OLLAMA_MODEL` — örn. `qwen2.5:3b`.

**İkisi de boşsa** Core otomatik olarak eski kural-tabanlı motora döner —
yani bu kurulumu yapmadan önce veya sırasında Core'u yeniden başlatmak
GÜVENLİ, davranış bozulmaz.

## Model Seçimi

Sunucu profili: 2-4 CPU, 4-8GB RAM, **GPU yok**, bu RAM Mongo/Redis/
RabbitMQ/Core/Panel ile paylaşılıyor. Öneri sırası:

1. **`qwen2.5:3b-instruct-q4_K_M`** (önerilen başlangıç) — ~2GB RAM,
   Türkçe dahil çok dilli, CPU'da birkaç saniyede yanıt verir.
2. `llama3.2:3b` — alternatif, benzer boyut/performans.
3. Eğer 1 numara sunucuda gözle görülür yavaşsa (>15-20sn/yanıt) veya RAM
   baskısı yaratıyorsa: `qwen2.5:1.5b` ya da `llama3.2:1b` gibi daha küçük
   bir modele düşün — kalite biraz düşer ama hız/kaynak baskısı azalır.

Model adını DEĞİŞTİRMEK tek bir env değişkeni (`OLLAMA_MODEL`) + Core
restart'ı kadar kolay — ilk seçim kesin olmak zorunda değil.

## A) Docker Compose Modunda Kurulum

`deployments/production/docker-compose.yml` dosyasına yeni bir servis
eklenir, `aystek-core` servisine de iki env satırı eklenir:

```yaml
services:
  # ... mongodb, redis, rabbitmq (değişiklik yok) ...

  ollama:
    image: ollama/ollama:latest
    container_name: aystek-ollama
    restart: unless-stopped
    ports:
      - "127.0.0.1:11434:11434"   # sadece sunucu içinden erişilebilir
    volumes:
      - ollama_data:/root/.ollama
    networks:
      - aystek_internal
    deploy:
      resources:
        limits:
          memory: 3g   # 4-8GB'lık sunucuda diğer servisleri boğmasın

  aystek-core:
    # ... mevcut tanım aynı kalır ...
    depends_on:
      - mongodb
      - redis
      - rabbitmq
      - ollama               # YENİ
    environment:
      # ... mevcut satırlar aynı kalır, ALTINA ekle:
      OLLAMA_BASE_URL: http://ollama:11434
      OLLAMA_MODEL: ${OLLAMA_MODEL:-qwen2.5:3b-instruct-q4_K_M}

volumes:
  # ... mevcutlar aynı kalır, ALTINA ekle:
  ollama_data:
```

Uygulama adımları:
```bash
cd /opt/aystek-core/current   # ya da gerçek compose dizini
# 1) compose.yml'i yukarıdaki gibi düzenle (elle veya bir editörle)
# 2) Ollama servisini ayağa kaldır
docker compose up -d ollama
# 3) Modeli indir (ilk sefer birkaç dakika sürebilir, ~2GB indirme)
docker compose exec ollama ollama pull qwen2.5:3b-instruct-q4_K_M
# 4) Core'u env değişiklikleriyle yeniden başlat
docker compose up -d aystek-core
# 5) Logda motor seçimini doğrula
docker compose logs aystek-core --tail=50 | grep "AYSTEK AI"
# Beklenen satır: "✅ AYSTEK AI: Ollama motoru aktif (http://ollama:11434, model=qwen2.5:3b-instruct-q4_K_M)"
```

## B) Native/systemd Modunda Kurulum

```bash
# 1) Ollama'yı kur (resmi kurulum betiği)
curl -fsSL https://ollama.com/install.sh | sh

# 2) Servisin ayakta ve localhost:11434'te dinlediğini doğrula
systemctl status ollama
curl -s http://127.0.0.1:11434/api/tags

# 3) Modeli indir
ollama pull qwen2.5:3b-instruct-q4_K_M

# 4) Core'un env dosyasına ekle
echo 'OLLAMA_BASE_URL=http://127.0.0.1:11434' | sudo tee -a /etc/aystek-core/core.env
echo 'OLLAMA_MODEL=qwen2.5:3b-instruct-q4_K_M' | sudo tee -a /etc/aystek-core/core.env

# 5) Core'u yeniden başlat
sudo systemctl restart aystek-core

# 6) Logda motor seçimini doğrula
sudo journalctl -u aystek-core -n 50 | grep "AYSTEK AI"
```

## Doğrulama (her iki mod için ortak)

1. Log satırı `✅ AYSTEK AI: Ollama motoru aktif (...)` görünmeli (görünmüyorsa
   env okunmamış demektir — env dosyası/compose satırları + servis
   restart'ı kontrol edilmeli).
2. Mobil uygulamada "AYSTEK AI" sekmesinden ("Asistan" DEĞİL — o ayrı,
   kullanıcının kendi API anahtarıyla çalışır) bir mesaj gönderip gerçek,
   şablon-olmayan bir yanıt geldiğini doğrula.
3. `docker stats` (Docker modu) ya da `free -h` + `top` (systemd modu) ile
   birkaç mesaj sonrası RAM/CPU baskısını izle — sunucudaki diğer
   servisleri (Mongo/Redis/RabbitMQ) aç kalıyorlar mı kontrol et.

## Geri Alma

Herhangi bir sorunda SADECE iki env değişkenini kaldırıp Core'u yeniden
başlatmak yeterli — kod otomatik olarak eski kural-tabanlı motora döner,
başka hiçbir değişiklik gerekmez:
```bash
# Docker: compose.yml'den OLLAMA_* satırlarını sil, `docker compose up -d aystek-core`
# systemd: core.env'den OLLAMA_* satırlarını sil, `sudo systemctl restart aystek-core`
```

## Sıradaki Adımlar (Faz 9b/9c — henüz uygulanmadı)

Bu kurulum tamamlanıp doğrulandıktan sonra:
- `management` modu desteği (bildirim/ilan/yorum aksiyonları) Go tarafına
  eklenecek.
- Kişisel bağlam (eşya/konum/seyahat) ve yönetim bağlamının (bekleyen
  öğeler) sohbete enjekte edilmesi — ayrı bir plan turu.
