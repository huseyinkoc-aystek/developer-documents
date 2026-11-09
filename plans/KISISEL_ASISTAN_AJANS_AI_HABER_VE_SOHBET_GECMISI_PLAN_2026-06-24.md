# Kişisel Asistan — Ajans Haberleri, AI Filtreleme ve Sohbet Geçmişi Kalıcılığı

**Tarih:** 2026-06-24  
**Proje:** `aystek-pers-assit` (kisisel_asistan Flutter uygulaması)  
**Durum:** Tasarım / Uygulama Planı

---

## 1. Vizyon

Kullanıcı abone olduğu haber ajanslarına (AA, ANKA ve ileride diğerleri) kendi kullanıcı adı/şifresi ile bağlanır. Uygulama bu ajansların API'lerinden haberleri çeker, AI (Claude/GPT/Gemini) ile özetler/filtreler ve tercih edilen il, kategori ile etiketlere göre kişiselleştirilmiş bir akış sunar. Seçilen haberlerde güncellenme olduğunda bildirim gönderilir. Tüm AI sohbet geçmişi cihazda kalıcı olarak saklanır; kullanıcı istediğinde devam ettirebilir veya silebilir.

---

## 2. Kapsam

### 2.1 Ajans Kimlik Bilgisi Yönetimi (Ajanslar Sekmesi)
- Kullanıcı adı + şifre `FlutterSecureStorage`'a (EncryptedSharedPreferences) şifreli kaydedilir.
- Her ajans için ayrı anahtar: `agency_credential_<id>`
- Desteklenen ajanslar: AA (`aa`), ANKA (`ankahaber`); ileride DHA, İHA, Sabah, vb.
- Bağlantı durumu kartlarda gösterilir; güncelleme/bağlantı kesme desteklenir.

### 2.2 Ajans API Entegrasyonu
Her ajansın API'si farklı; bağlantı stratejileri:

| Ajans | Yöntem | Not |
|-------|--------|-----|
| Anadolu Ajansı (AA) | OAuth2 / Basic Auth + REST | `haber.aa.com.tr` — API dokümantasyonu talep edilmeli |
| ANKA Haber | Session cookie + scraping veya API | `ankahaber.net` |

**Backend mimarisi:**  
Ajans API çağrıları doğrudan mobil cihazdan değil **AYSTEK Core backend** üzerinden yapılır. Nedenler:
- Kimlik bilgileri sunucuya gönderilmez; token/cookie yönetimi backend'de.
- Rate limiting, cache ve retry mantığı merkezi tutulur.
- İleride çok kullanıcı desteklenebilir.

AYSTEK Core'da yeni endpoint önerileri:
```
POST /api/v1/news/agency/connect          — token/cookie al, cihaza session_token dön
GET  /api/v1/news/agency/feed             — filtrelenmiş haber listesi
GET  /api/v1/news/agency/article/{id}     — tam metin
POST /api/v1/news/agency/watch/{id}       — makale izlemeye al (bildirim için)
```

### 2.3 AI Destekli Filtreleme
- **Model:** Claude Haiku 4.5 (hız+maliyet dengesi) veya kullanıcı tercihine göre GPT-4o/Gemini
- **Pipeline:**
  1. Backend ajans API'den ham haber listesi çeker (başlık + özet + etiketler).
  2. AI şu filtreleri uygular:
     - **İl/Bölge:** Kullanıcı tercihindeki şehir/bölgeyle ilgili mi?
     - **Kategori:** Kullanıcının seçtiği kategoriler (siyaset, ekonomi, spor, teknik, kamu).
     - **Etiket:** Kullanıcı tanımlı anahtar kelimeler.
  3. Skor 0-100 arası; eşik üstü akışa alınır.
- **Prompt şablonu** (backend'de) kullanıcı tercihlerini dinamik enjekte eder.
- Filtreleme sonucu önbelleğe alınır (Redis, 15 dk).

### 2.4 Bildirim Sistemi
- **İzleme (watch):** Kullanıcı bir makaleyi "izle"ye alırsa AYSTEK Core dakika başı kontrol eder.
- **Yeni eşleşme:** Filtreleme pipeline'ı yeni puanı yüksek haber bulduğunda push bildirimi gönderir.
- Mobil tarafta mevcut `NewsNotificationService` + `WorkManager` genişletilir.
- Bildirime tıklandığında ilgili haber detayına yönlendirilir.

### 2.5 AI Sohbet Geçmişi Kalıcılığı (Kritik Hata Düzeltmesi)
**Mevcut sorun:** Uygulama her açıldığında sohbet geçmişi sıfırlanıyor.

**Çözüm mimarisi:**

```
lib/
  models/
    chat_session.dart       — ChatSession, ChatMessage modelleri
  services/
    chat_history_service.dart — Hive/SQLite üzerinde CRUD
  providers/
    chat_providers.dart     — chatSessionsProvider, activeChatProvider
  screens/
    chat/
      chat_screen.dart      — mevcut ekran, provider'dan okunur
      chat_sessions_screen.dart — geçmiş listesi, devam etme, silme
```

**Veri modeli:**
```dart
class ChatSession {
  final String id;            // UUID
  final String title;         // İlk mesajdan otomatik türetilir
  final String aiProvider;    // 'claude' | 'openai' | 'gemini'
  final String model;
  final DateTime createdAt;
  final DateTime updatedAt;
  final List<ChatMessage> messages;
}

class ChatMessage {
  final String role;          // 'user' | 'assistant' | 'system'
  final String content;
  final DateTime timestamp;
}
```

**Depolama:** `Hive` (yerel, hızlı, şifreleme desteği var) veya `sqflite`.  
Tercih: `Hive` — Flutter'da native, key-value, binary format, hızlı.

**UX akışı:**
1. Sohbet ekranı açıldığında mevcut aktif oturum gösterilir.
2. "Yeni Sohbet" butonu → yeni oturum oluşturulur.
3. "Geçmiş" butonu → `ChatSessionsScreen`: tüm oturumlar listesi.
4. Listeden oturuma tıklanınca o oturum aktif olur, mesajlar yüklenir.
5. Kullanıcı sola kaydırınca oturum silinir (confirm dialog).
6. Oturum başlığı: ilk kullanıcı mesajının ilk 50 karakteri (veya AI tarafından üretilmiş özet).

**Provider tasarımı:**
```dart
// Tüm oturumlar
final chatSessionsProvider = StateNotifierProvider<ChatSessionsNotifier, List<ChatSession>>(...);

// Aktif oturum ID
final activeChatIdProvider = StateProvider<String?>(...);

// Aktif oturum mesajları
final activeChatMessagesProvider = Provider<List<ChatMessage>>((ref) {
  final id = ref.watch(activeChatIdProvider);
  final sessions = ref.watch(chatSessionsProvider);
  return sessions.firstWhereOrNull((s) => s.id == id)?.messages ?? [];
});
```

---

## 3. Uygulama Sırası

### Faz 1 — Hızlı Kazanımlar (1-2 sprint)
1. **Sohbet geçmişi kalıcılığı** — Hive entegrasyonu, `ChatSessionsScreen`, provider refactor.
2. **Ajanslar sekmesi crash düzeltmesi** — ✅ Tamamlandı (SingleChildScrollView + parent state lift).

### Faz 2 — Backend Altyapısı (2-3 sprint)
3. AYSTEK Core'a ajans session yönetimi endpoint'leri.
4. Ajans API adaptörleri (AA ve ANKA için).
5. AI filtreleme pipeline (Haiku 4.5 + kullanıcı tercihleri).

### Faz 3 — Mobil Entegrasyon (1-2 sprint)
6. Flutter'da ajans feed ekranı (filtrelenmiş akış).
7. Bildirim sistemi genişletmesi (izleme + yeni eşleşme).
8. AI sohbetinde haber bağlamı enjeksiyonu (seçili haber → sohbete ekle).

### Faz 4 — Yeni Ajanslar ve Genişleme
9. DHA, İHA, NTV vb. adaptörler.
10. Makale tam metin okuma + AI özet.
11. Sohbet oturumlarını AYSTEK Core'a yedekleme (opsiyonel).

---

## 4. Teknik Kısıtlar ve Kararlar

| Konu | Karar | Gerekçe |
|------|-------|---------|
| Ajans kimlik bilgisi | Cihazda `FlutterSecureStorage` | Sunucuya şifre gönderilmez |
| Ajans API çağrısı | Backend proxy üzerinden | Rate limit, token yönetimi, cache |
| AI filtreleme | Backend'de (AYSTEK Core) | Model değiştirilebilirlik, maliyet kontrolü |
| Sohbet geçmişi | Hive, cihazda | Gizlilik, offline erişim |
| Bildirimler | FCM + WorkManager | Mevcut altyapı, düşük pil etkisi |

---

## 5. Açık Sorular

- AA ve ANKA'nın developer API erişimi mevcut mu? Sözleşme gerekiyor mu?
- AYSTEK Core'da ajans proxy modülü yeni bir microservice mi yoksa mevcut `integrations` modülüne mi eklenmeli?
- Sohbet geçmişi şifrelenmeli mi? (Hive şifreleme eklentisi var)
- Kullanıcı başına kaç oturum sınırı? Cihaz depolama yönetimi nasıl olacak?

---

## 6. Bağlantılı Dokümanlar

- `project_kisisel_asistan.md` — Genel uygulama vizyonu
- `project_aystek_core_vision.md` — AYSTEK Core entegrasyon hedefleri
- `AYSTEK CORE Integrations Modülüne Enterprise Düzeyde...` — Entegrasyon altyapısı
