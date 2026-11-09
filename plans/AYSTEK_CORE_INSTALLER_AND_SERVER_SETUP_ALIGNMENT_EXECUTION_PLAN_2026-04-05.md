# AYSTEK Core Installer and Server Setup Alignment Execution Plan

**Tarih:** 2026-04-05  
**Durum:** In Progress  
**Kapsam:** `aystek-core-installer` + server setup / deployment / operations documentation alignment

## Amaç

`aystek-core-installer`, `aystek-core` ve operasyon dokümantasyonu arasında kalan sözleşme farklarını kapatmak; kurulum rehberlerini güncel kurulum gerçekliğine hizalamak; eski `aystek-updater` / `14 adımlı atomik akış` anlatılarından doğan operatör riskini temizlemek.

## Tespit Edilen Başlıklar

1. Binary adı ve bootstrap akışı eski `aystek-installer` isminde kalmış.
2. `SUNUCU_KURULUM_REHBERI.md` içinde installer davranışı bazı yerlerde fazla geniş, bazı yerlerde tarihsel olarak eski.
3. `DEPLOYMENT_SCENARIOS.md`, `OVERVIEW.md` ve `system_ops.md` içinde eski updater/rollback modeli anlatılıyor.
4. `README_FOR_AI_ASSISTANTS.md` içinde installer-operasyon dokümanı güncelleme zorunluluğu yeterince açık değil.
5. Fresh-DB kurulumlarda ilk admin bootstrap akışı eksikti; `admin` ve `tenant-site` scope için operatör kontrollü bootstrap kullanıcı oluşturma gereksinimi açık değildi.
6. Installer repo içinde `tests/integration`, `tests/fixtures`, `tests/scripts` tabanlı gerçek bir integration harness iskeleti yoktu.

## Çalışma Dilimleri

### Slice 1 — Installer Canonical Naming and Bootstrap
- `Makefile` binary/module path uyumu
- bootstrap script binary adı ve release artifact adı uyumu

### Slice 2 — Server Guide Alignment
- `SUNUCU_KURULUM_REHBERI.md` metadata, sürüm ve kapsam güncellemesi
- installer-native/manual ayrımının bugünkü koda göre netleştirilmesi
- eski service / log / updater referanslarının temizlenmesi

### Slice 3 — Related Operations Documentation
- `DEPLOYMENT_SCENARIOS.md`
- `OVERVIEW.md`
- `system_ops.md`

### Slice 4 — Standards Hardening
- `README_FOR_AI_ASSISTANTS.md` içine installer-doküman güncelleme zorunluluğu eklenmesi

### Slice 5 — Release Artifact Contract and Integration Test Matrix
- canonical release artifact isimlendirmesini netleştirme
- native vs docker input sözleşmesini yazılı hale getirme
- bootstrap, install, update, rollback, backup, tenant-site ve mongodb-only akışları için integration test matrisi çıkarma
- hangi alanların yalnız unit test ile, hangilerinin VM/container smoke test ile doğrulanacağını ayırma

### Slice 6 — Fresh DB Bootstrap Admin and P0 Harness Scaffold
- `admin` scope fresh-DB kurulumlarda ilk admin kullanıcı bilgisini toplama
- `tenant-site` scope fresh-DB kurulumlarda ilk tenant admin bilgisini toplama
- Mongo seed katmanına bootstrap admin create/reconcile ekleme
- `tests/integration`, `tests/fixtures`, `tests/scripts` altında P0 harness iskeleti açma
- `IT-03`, `IT-05`, `IT-06`, `IT-07`, `IT-11` akışlarını scaffold etme

## Başarı Kriterleri

- Operatör bir dokümandan okuduğu komutu bugünkü installer ile gerçekten uygulayabiliyor olmalı.
- Binary adı, repo adı, bootstrap akışı ve doküman adı birbiriyle çelişmemeli.
- System Ops / deployment anlatıları, artık gerçek kodda olmayan updater-daemon veya migration-reversal varsayımlarını taşımamalı.
- Gelecekte installer değiştiğinde ilgili dokümanların da aynı işte güncellenmesi standart hale gelmeli.
- Release artifact beklentileri ile installer giriş alanları (`CoreBinaryPath`, `PanelBuildPath`, image tag'leri, bootstrap archive adı) tek bir sözleşmede toplanmalı.
- Integration test kapsamı belge seviyesinden çıkarılıp uygulanabilir bir matrise dönmelidir.
- Fresh-DB kurulumlarda operatör, ilk giriş için gerekli admin hesabını kurulum sırasında güvenli biçimde tanımlayabilmelidir.
