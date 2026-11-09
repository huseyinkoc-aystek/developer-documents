# Giden E-posta (SMTP) Kurulumu — Gmail / Google Workspace

AYSTEK Core'un gönderdiği tüm e-postalar (hesap doğrulama, şifre sıfırlama, sipariş onayı, bildirimler)
tek bir SMTP hesabından çıkar. Bu rehber `info@aysteknolojileri.com` (Google Workspace) içindir.

## 1. Google tarafı (bir kez)

1. `info@aysteknolojileri.com` ile Google hesabına girin → **Güvenlik** → **2 Adımlı Doğrulama**'yı açın.
   (Workspace'te yönetici konsolunda *Güvenlik → Kimlik doğrulama → 2 Adımlı Doğrulama* kullanıcıya izinli olmalı.)
2. <https://myaccount.google.com/apppasswords> → ad: `AYSTEK Core` → **Oluştur**. 16 karakterli
   **uygulama şifresi** görünür (boşluklu yazılabilir). Normal Gmail şifresi SMTP'de çalışmaz.
   - "Uygulama şifreleri" seçeneği yoksa: 2 Adımlı Doğrulama kapalı ya da Workspace yöneticisi kapatmıştır.
3. Gönderilen postanın spam'e düşmemesi için (özellikle Fransa'daki müşteriler): Workspace yönetici
   konsolu → **Uygulamalar → Google Workspace → Gmail → E-posta kimlik doğrulama** ile **DKIM**'i
   `aysteknolojileri.com` için açın; DNS'te **SPF** (`v=spf1 include:_spf.google.com ~all`) ve **DMARC**
   (`_dmarc`) kayıtları bulunsun.

## 2. Yönetim paneli (aystek-core-next)

**Ayarlar → E-posta hesapları → Yeni hesap**

| Alan | Değer |
|---|---|
| Ad | `Gmail (info)` |
| Sunucu (host) | `smtp.gmail.com` |
| Port | `587` + **TLS** işaretli (STARTTLS) — alternatif: `465` + **SSL** işaretli |
| Kullanıcı adı | `info@aysteknolojileri.com` |
| Şifre | Google'dan aldığınız 16 haneli uygulama şifresi |
| Gönderen e-posta | `info@aysteknolojileri.com` (Gmail başka adresi kendi adresiyle değiştirir; farklı gönderen için Gmail'de "Başka ad olarak gönder" tanımlı olmalı) |
| Gönderen adı | `AYSTEK` (**tenant geneli**: tüm siteler — balci, AYSTEK Haber… — aynı hesaptan çıkar) |
| Amaç | `system` |
| Aktif / Varsayılan | ikisi de açık |

Kaydettikten sonra satırdaki **Test** düğmesi **giriş yapmış yöneticinin kendi e-posta adresine** bir mesaj yollar (kullanıcıda e-posta yoksa gönderilmez); hata olursa
SMTP'nin gerçek nedeni (ör. `535 Username and Password not accepted`) ekranda görünür.

- Şifre alanı düzenlemede **boş bırakılırsa mevcut şifre korunur**; şifre API'den hiç geri dönmez.
- **Şifre şifreli saklansın diye** Core ortamında `PAYMENT_GATEWAY_ENCRYPTION_KEY` tanımlı olmalı
  (ödeme ağ geçidi sırlarıyla aynı anahtar; yeni installer kurulumları `core.env`'e otomatik yazar, **mevcut sunucularda** `core.env`'e `PAYMENT_GATEWAY_ENCRYPTION_KEY=<ENCRYPTION_KEY ile aynı değer>` ekleyip Core'u yeniden başlatın). Anahtar yoksa şifre
  veritabanında düz metin kalır. Anahtarı sonradan değiştirirseniz kayıtlı şifre çözülemez: hesabı açıp şifreyi
  yeniden girin.
- Öncelik sırası: **varsayılan e-posta hesabı → sunucu ortamı (`SMTP_*` env, yalnız yedek)**. Hesap
  tanımlıyken env kullanılmaz; hesap yokken `.env`/`core.env` değerleri kullanılır.
- **Eski `smtp_*` ayarları kaldırıldı.** Bu anahtarlar hiçbir ekrandan düzenlenemiyordu. Seed
  (`SeedEmailAccounts`) etkin ve eksiksiz bir eski yapılandırmayı bir kez "SMTP (migrated)" adlı varsayılan
  hesaba taşır, eski `smtp_pass` değerini boşaltır ve `smtp_enabled`'ı kapatır; hesap zaten varsa hiçbir şey yapmaz,
  yer tutucu hesap oluşturmaz. Taşınan şifre olduğu gibi kopyalanır: hesabı panelde bir kez kaydedince şifrelenir.

## 3. Sınırlar ve notlar

- Gmail/Workspace günlük gönderim sınırı vardır (Workspace ≈ 2000 alıcı/gün, ücretsiz Gmail ≈ 500). Sipariş
  e-postaları için yeterlidir; toplu bülten için ayrı bir işlem sağlayıcısı (Brevo, Mailjet, SES…) düşünün —
  aynı ekrandan `marketing` amaçlı ikinci bir hesap olarak eklenebilir.
- Gönderen adı tenant geneli olduğundan site adı (örn. "Rucher du VELAY") e-postanın **içinde** ve konu
  satırında yer alır (`{{site_name}}`), gönderen adında değil.
- Şifreyi geliştirme ortamında `.env`'e yazmayın; `.env` git'te izlenir.
