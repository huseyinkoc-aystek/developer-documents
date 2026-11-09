# AYSTEK Core Installer Seed Chain And Disposable P0 Harness Checkpoint

Tarih: 2026-04-05

## Tamamlananlar

- Installer `db_mongo_seed` adımı gerçek `aystek-core` seed zincirini çalıştıracak şekilde bağlandı.
- Bu sayede `translation_seed`, `menu_seed` ve diğer çekirdek seed içerikleri kurulum ve update sırasında veritabanına reconcile edilir.
- `IT-05` ve `IT-11` disposable runtime harness script'leri artık gerçek Mongo/Redis container + `mongosh` assertion akışını hedefliyor.
- Test ortamları için tenant registry path override desteği açıldı.
- Seed closure, integration harness tarafından doğrudan çağrılabilir bir installer operasyonu olarak açıldı.

## Açılan Teknik Yapılar

- `installer.RunSeedClosure()`
- `tenant.ResolvedRegistryPath()`
- `tests/integration/cmd/seed-runtime-harness`

## Sonraki Doğru Adım

- Docker erişimi olan CI veya disposable host ortamında `IT-05` ve `IT-11` script'lerini gerçekten çalıştırıp pipeline'a almak.
