# TASK 063 Özeti

Balcı mağazası Next.js `15.5.27` ve React `19.3.0` sürümlerine geçirildi. App Router'ın yeni
asenkron API sözleşmeleri uyarlandı; ESLint flat config, CI üretim audit kapısı ve haftalık
Dependabot kontrolleri eklendi. Üretim bağımlılık ağacında güvenlik bulgusu kalmadı.

Doğrulama: temiz kurulum, env sözleşmesi, ESLint, TypeScript, 12 birim testi, env'li/env'siz
production build ve standalone rota smoke testi başarılı. `e3f52f0` sürümü installer ile canlıya
alındı; public rotalar/TLS/Host admission doğrulandı ve Core upload dosyaları değişmedi.

Detay: `TASK_063_BALCI_NEXT15_SECURITY_UPGRADE.md`.
