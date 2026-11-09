# TASK 064 Özeti

Medya seçicinin boş görünmesine yol açan Core MIME alanı hatası düzeltildi; panel yeni yüklenen
dosyayı otomatik seçer. CMS'in `/static/...` görselleri ortak haber public uygulamasında güvenli,
site-kapsamlı same-origin proxy üzerinden sunulur. Haber, SEO, yazar, galeri ve reklam görselleri
aynı adaptörü kullanır.

Commitler GitHub'a aktarıldı: Core `21ef8c6`, panel `5c910c2`, ortak haber public
`f841f86`. Canlı dağıtım, çalışma ortamında VPN rotası bulunmadığı için bekliyor.

Detay: `TASK_064_MEDIA_LIBRARY_SELECTION_PUBLIC_DELIVERY.md`.
