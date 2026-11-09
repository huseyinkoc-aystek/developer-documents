# TASK 063 — Balcı Next.js 15 Güvenlik Yükseltmesi

**Tarih:** 2026-10-12 · **Durum:** Kod, bağımlılık taraması, regresyon doğrulamaları ve canlı dağıtım tamamlandı.

## Amaç

`balci-aysteknolojileri-com` üretim bağımlılıklarında raporlanan kritik/yüksek Next.js
güvenlik borcunu, tenant/site runtime kimliğini ve mevcut mağaza davranışını bozmadan kapatmak.

## Uygulanan değişiklikler

- Next.js `14.2.33` → `15.5.27`; React ve React DOM → `19.3.0`.
- Next.js 15 App Router sözleşmesine göre `params`, `searchParams` ve `headers()` çağrıları
  asenkron hale getirildi.
- ESLint 9 flat config'e geçildi; eski `.eslintrc.json` kaldırıldı.
- Next.js 15'te kaldırılan `swcMinify` seçeneği kaldırıldı.
- Next.js'in sabitlediği eski transitif PostCSS sürümü, kilit dosyasında da etkili olacak şekilde
  `8.5.29` ile override edildi.
- CI'a yalnız üretim ağacını denetleyen `npm run audit:prod` kapısı eklendi.
- Haftalık npm ve GitHub Actions Dependabot yapılandırması eklendi.

## Güvenlik sonucu

`npm audit --omit=dev --audit-level=moderate` sonucu: **0 bulgu** (info/low/moderate/high/critical).

Tam geliştirme ağacında kalan bildirimler çalışma zamanı paketlerini etkilemez; ESLint/Tailwind
araç zincirindedir. Bunları zorla çözmek Next.js 16 ve/veya Tailwind 4 gibi daha geniş bir geçiş
gerektirdiğinden bu güvenlik yamasına karıştırılmadı.

## Doğrulama

- Temiz kurulum: `npm ci`.
- Ortam sözleşmesi: `npm run check-env`.
- Üretim bağımlılık taraması: `npm run audit:prod`.
- ESLint ve TypeScript (`tsc --noEmit`) temiz.
- 12/12 birim testi başarılı.
- Runtime env ile production build başarılı.
- Hiçbir `.env*` olmadan production build başarılı; reusable image sözleşmesi korundu.
- Standalone sunucuda ana mağaza, katalog, içerik, hesap, sepet, robots, sitemap ve health
  rotaları 200; tanımsız Host 421.

## Dağıtım ve geri dönüş

Dağıtım installer'ın siteye özel `local-source` profiliyle, benzersiz `git-<sha>` image etiketi
üzerinden yapılır. Güncelleme öncesi installer config, compose/env ve Core upload manifesti ayrı
bir geri dönüş dizinine alınır. Sorunda önceki image etiketi ve config ile public-site update
yeniden çalıştırılır. Public-site güncellemesi MongoDB'ye veya Core upload volume'üne dokunmaz.

Canlı dağıtım `e3f52f0` commit'i ve `git-e3f52f0` image etiketiyle tamamlandı. Konteyner içinden
Next.js sürümü `15.5.27` olarak doğrulandı. Ana sayfa, kataloglar, hikâye, blog, iletişim, hesap,
sepet, robots, sitemap ve health rotaları dış HTTPS üzerinden 200 döndü; tanımsız Host 421 döndü.
TLS sertifikası geçerli, konteyner logları temiz ve Core upload manifesti dağıtım öncesiyle aynı.
Geri dönüş paketi: `/var/backups/aystek-core/pre-balci-next15-e3f52f0`.

## Sonraki kontrollü yükseltme

Next.js 16 geçişi ayrı görev olmalıdır: Node sürümü, Turbopack, ESLint entegrasyonu, async API
değişiklikleri ve tüm mağaza rotaları yeniden regresyon testine alınmalıdır.
