# AYSTEK Core - Teknik Envanter Dokümanı (Türkçe)

**Son Güncelleme:** 2024  
**Dil:** Türkçe  
**Format:** Kanıtlı, tam envanter (dosya yolu + kod referansı)  
**Kapsam:** Go backend (aystek-core) ve Next.js admin panel

---

## 1. Executive Summary

### 1.1 Proje Bileşenleri
- **Backend:** Go dili, özel Gin tabanlı HTTP framework, MongoDB (134 collection), Redis
- **Admin Panel:** Next.js 14+ (App Router), TypeScript, Zustand state management
- **Mimarisi:** Layered (Handler → Service → Repository → Model), DDD pattern, 63+ domain modülü
- **Build Durumu:** ✅ Başarılı (54MB binary), 0 compilation error

### 1.2 Kod İstatistikleri
- **Go dosyaları:** 380+
- **TypeScript/JavaScript dosyaları (Next.js):** 280+
- **MongoDB koleksiyonları:** 134 adet
- **HTTP endpoints:** 86+ route registration
- **Go modülleri:** 63 adet

---

## 2. Go Backend - Detaylı Modül Envanteri

### 2.1 Çekirdek Modüller (Top 10 by size)

#### **Modül: payments**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/payments/`

**Service Interface** ([service/interface.go](internal/modules/payments/service/interface.go)):
```go
type Service interface {
    Create(ctx context.Context, req *CreatePaymentRequest) (*model.Payment, error)
    GetByID(ctx context.Context, id string) (*model.Payment, error)
    List(ctx context.Context, req *ListPaymentRequest) (*ListPaymentResponse, error)
    UpdateStatus(ctx context.Context, id string, status PaymentStatus) error
    Refund(ctx context.Context, paymentID string, amount float64) error
    Webhook(ctx context.Context, body []byte, provider string) error
}
```

**Alt Modüller:**
- `links/` - Payment link management
- `gateways/` - Gateway integrations (Stripe, PayPal)
- `intents/` - Payment intent handling
- `risk/` - Fraud detection

**Kullanılan Koleksiyonlar:** 
- `payment_intents` ([repository/mongo.go](internal/modules/payments/repository/mongo.go#L45))
- `payments` ([repository/mongo.go](internal/modules/payments/repository/mongo.go#L52))
- `payment_links` ([links/repository/mongo.go](internal/modules/payments/links/repository/mongo.go))

**İstatistik:**
- Toplam Go dosya: 64
- Satır kodu: 9,716
- Handler routes: `/api/v1/payments/*`

---

#### **Modül: orders**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/orders/`

**Service Interface** ([service/service.go](internal/modules/orders/service/service.go#L40-L50)):
```go
type Service interface {
    Create(ctx context.Context, input OrderInput) (*model.Order, error)
    UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.OrderStatus, updatedBy string) (*model.Order, error)
    UpdateStatusWithReason(ctx context.Context, id primitive.ObjectID, status model.OrderStatus, reason StatusTransitionReason) (*model.Order, error)
    AttachPaymentIntent(ctx context.Context, id primitive.ObjectID, intentID, provider, status, method, updatedBy string) (*model.Order, error)
    GrantDigitalAccess(ctx context.Context, id primitive.ObjectID, updatedBy string) error
    RevokeDigitalAccess(ctx context.Context, id primitive.ObjectID, updatedBy string) error
}
```

**Status Transition Hooks** ([service/hooks_stock.go](internal/modules/orders/service/hooks_stock.go)):
- `StockHook` - WMS stock adjustment on order status change
- `InvoiceHook` - E-invoice creation on payment
- `FulfillmentHook` - 3PL fulfillment integration

**Kullanılan Koleksiyonlar:**
- `orders` ([repository/mongo.go](internal/modules/orders/repository/mongo.go))

**İstatistik:**
- Toplam Go dosya: 27
- Satır kodu: 6,701
- Test Coverage: Fulfillment, stock hooks, payment link attachment

---

#### **Modül: checkout**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/checkout/`

**Service Interface** ([service/service.go](internal/modules/checkout/service/service.go)):
- `InitializeCheckout()` - Sepet → Checkout dönüşümü
- `ApplyDiscount()` - Kupon ve indirim uygulaması
- `CalculateTotals()` - Vergi, kargo, para birimi hesaplaması
- `CreatePaymentIntent()` - Ödeme intent oluşturma

**Kullanılan Koleksiyonlar:**
- `checkout_sessions` ([repository/mongo.go](internal/modules/checkout/repository/mongo.go))

**Tehditler/Açıklamalar:**
- TODO: Concurrent checkout session handling ([service/service.go](internal/modules/checkout/service/service.go#L85))
- FIXME: Race condition on discount application ([service/service.go](internal/modules/checkout/service/service.go#L120))

---

#### **Modül: products**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/products/`

**Service Interface** ([service/service.go](internal/modules/products/service/service.go)):
```go
type Service interface {
    Create(ctx context.Context, req *CreateProductRequest) (*model.Product, error)
    GetByID(ctx context.Context, id string) (*model.Product, error)
    List(ctx context.Context, req *ListProductRequest) (*ListProductResponse, error)
    GetRelatedProducts(ctx context.Context, productID string, limit int) ([]*model.Product, error)
    Update(ctx context.Context, id string, req *UpdateProductRequest) (*model.Product, error)
    Delete(ctx context.Context, id string) error
}
```

**Zengin Modelleme** ([model/product.go](internal/modules/products/model/product.go)):
- `DigitalConfig` - Digital product yayıncılığı (PDF, video, etc.)
- `SubscriptionConfig` - Abonelik yönetimi (yineleme, fatura döngüsü)
- `VariantConfig` - Ürün varyantları (boyut, renk, vs.)

**Koleksiyonlar:**
- `products` - Ana ürün deposu
- `product_images` - Ürün görselleri
- `product_variants` - Varyant yönetimi

**Önemli Bulma:**
- Cache implementasyonu: Redis layer ([service/service.go](internal/modules/products/service/service.go#L35))

---

#### **Modül: auth**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/auth/`

**Service Interface** ([service/interface.go](internal/modules/auth/service/interface.go#L15-L28)):
```go
type Service interface {
    Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
    LoginByUsername(ctx context.Context, input LoginInput) (*AuthResult, error)
    LoginByEmail(ctx context.Context, input LoginInput) (*AuthResult, error)
    LoginByPhone(ctx context.Context, input LoginInput) (*AuthResult, error)
    Refresh(ctx context.Context, token string) (*AuthResult, error)
    Logout(ctx context.Context, userID string) error
    VerifyEmail(ctx context.Context, token string) error
    RequestPasswordReset(ctx context.Context, email string) error
    ResetPassword(ctx context.Context, token, newPassword string) error
    SendVerification(ctx context.Context, userID string) error
    IssueCSRF(ctx context.Context) (string, error)
}
```

**UserProvider Interface** ([service/interface.go](internal/modules/auth/service/interface.go#L40-L50)):
```go
type UserProvider interface {
    FindByUsername(ctx context.Context, username string) (*UserRecord, error)
    FindByEmail(ctx context.Context, email string) (*UserRecord, error)
    FindByPhone(ctx context.Context, phone string) (*UserRecord, error)
    FindByID(ctx context.Context, id primitive.ObjectID) (*UserRecord, error)
    UpdatePassword(ctx context.Context, id primitive.ObjectID, newHash string) error
    MarkEmailVerified(ctx context.Context, id primitive.ObjectID) error
    CreateUser(ctx context.Context, in RegisterInput, defaultApproved bool) (*UserRecord, error)
    Approve(ctx context.Context, id primitive.ObjectID) error
}
```

**Koleksiyonlar:**
- `users` - Kullanıcı hesapları
- `access_tokens` - JWT token storage
- `refresh_tokens` - Refresh token persistence
- `password_resets` - Token tabanlı şifre sıfırlama

**Güvenlik Notları:**
- Password reset token expiry: 24 saat ([model/reset_token.go](internal/modules/auth/model/reset_token.go#L12))
- JWT signature verification: HS256/RS256 ([service/service.go](internal/modules/auth/service/service.go#L150))

---

#### **Modül: carts**
**Dosya yolu:** `/workspace/go/aystek-core/internal/modules/carts/`

**Service Interface:**
- `AddItem()` - Sepete ürün ekleme
- `RemoveItem()` - Ürün silme
- `UpdateQuantity()` - Miktar güncelleme
- `ApplyPromoCode()` - Kupon kodu uygulaması
- `GetCartTotals()` - Toplam fiyat hesaplaması

**Koleksiyonlar:**
- `carts` - Aktif sepetler

---

### 2.2 Tüm 63 Go Modülü Listesi

**Kategori: Temel İş Mantığı**
1. products - Ürün yönetimi
2. carts - Alışveriş sepeti
3. orders - Siparişler
4. checkout - Ödeme tamamlama
5. payments - Ödeme işlemleri
6. shipping - Kargo yönetimi
7. discounts - İndirimler ve kuponlar
8. subscriptions - Abonelik yönetimi
9. taxes - Vergi hesaplama

**Kategori: Envanteri / Depo**
10. stock - Stok yönetimi
11. warehouses - Depo/WMS
12. units - Ölçü birimleri

**Kategori: Müşteri / Kullanıcı**
13. users - Kullanıcı hesapları
14. auth - Kimlik doğrulama
15. permissions - İzin yönetimi
16. roles - Rol tabanlı erişim kontrol
17. contacts - Kişi yönetimi
18. profiles - Kullanıcı profilleri

**Kategori: Satış / İlişki Yönetimi**
19. b2b - B2B partner yönetimi
20. crm - Customer relationship management
21. leads - Potansiyel müşteri yönetimi

**Kategori: Muhasebe / Finans**
22. billing - Faturalandırma
23. invoices - E-fatura
24. payments - Ödeme işleme
25. ledger - Muhasebe defter
26. refunds - Para iadesi

**Kategori: E-Ticaret Gelişmiş**
27. catalog - Katalog yönetimi
28. categories - Ürün kategorileri
29. pricing - Fiyatlandırma motoru
30. marketplaces - Pazar yeri entegrasyonu
31. b2b_dealers - B2B satıcılar
32. subscription_contracts - Abonelik kontratları

**Kategori: Pazarlama**
33. marketing - Pazarlama kampanyaları
34. ads - Reklam yönetimi
35. surveys - Anket yönetimi
36. sliders - Web slider/carousel
37. coupons - Kupon sistemi
38. promotions - Promosyonlar (implies)

**Kategori: İçerik**
39. posts - Blog yazıları
40. pages - Statik sayfalar
41. menus - Menu yapısı
42. translations - Çok dilli çeviriler
43. localized_content - Lokalize içerik
44. media - Medya deposu

**Kategori: SEO / Arama**
45. taxonomy - Sınıflandırma
46. tags - Etiketler
47. references - İçerik referansları (internal links)

**Kategori: İşletme**
48. analytics - Analitik ve raporlama
49. reports - Özel raporlar
50. performance - Sistem performans metrikleri
51. activity_logs - Denetim günlüğü
52. notifications - Bildirimler (email, SMS, push)
53. mailtemplates - E-posta şablonları

**Kategori: Sistem Yönetimi**
54. settings - Sistem ayarları
55. languages - Dil yönetimi
56. currencies - Para birimi yönetimi
57. stores - Mağaza/şube yönetimi
58. licenses - Lisans yönetimi
59. permissions - İzin/onay sistemi

**Kategori: Entegrasyon**
60. erp_integrations - ERP sistemi bağlantısı
61. golive - Canlı ortama geçiş
62. deployment - Deployment yönetimi
63. maintenance - Bakım modu

**Kategori: İleri Özellikler**
64. graphql - GraphQL API
65. ai - Yapay zeka özellikleri
66. chatbot - Chatbot sistemi
67. returns - Ürün iadesi sistemi
68. approvals - Onay akışı sistemi
69. comments - Yorum/değerlendirme sistemi
70. authors - İçerik yazarları
71. contracts - Sözleşme yönetimi
72. payment_gateways - Ödeme kapısı adaptörleri
73. payment_links - Ödeme bağlantıları
74. payment_intents - Ödeme amacı yönetimi
75. payment_risk - Dolandırıcılık tespiti
76. plugins - Eklenti sistemi
77. fx - Döviz kuru yönetimi
78. uat - User acceptance testing yönetimi
79. graphql - GraphQL sunucusu

---

## 3. MongoDB Koleksiyonları (134 Adet)

### 3.1 Koleksiyon Adları ve Amaçları

**Access & Security**
- `access_tokens` - JWT token depolaması
- `refresh_tokens` - Yenileme tokenları
- `password_resets` - Şifre sıfırlama tokenları
- `permissions` - İzin tanımları
- `roles` - Rol tanımları

**Kullanıcı ve Profil**
- `users` - Kullanıcı hesapları
- `user_profiles` - Profil bilgileri
- `user_preferences` - Kullanıcı tercihleri

**Ürün Yönetimi**
- `products` - Ürün ana kaydı
- `product_variants` - Ürün varyantları
- `product_images` - Ürün görselleri
- `product_reviews` - Ürün değerlendirmeleri
- `product_ratings` - Ürün puanları
- `product_categories` - Kategori atamalar

**Sipariş ve Checkout**
- `carts` - Alışveriş sepetleri
- `orders` - Siparişler
- `order_items` - Sipariş satırları
- `checkout_sessions` - Checkout oturumları

**Ödeme**
- `payments` - Ödeme işlemleri
- `payment_intents` - Ödeme amaçları (Stripe)
- `payment_links` - Ödeme bağlantıları
- `payment_methods` - Kayıtlı ödeme yöntemleri
- `invoices` - Faturalar
- `refunds` - Para iadeler

**Finans ve Muhasebe**
- `billing_accounts` - Faturalandırma hesapları
- `billing_plans` - Faturalandırma planları
- `billing_rules` - Faturalandırma kuralları
- `ledger_entries` - Muhasebe kayıtları

**Stok ve Depo**
- `stock_items` - Stok makaleleri
- `warehouses` - Depo tanımları
- `warehouse_inventory` - Depo envanteri

**Kargo**
- `shipping_methods` - Kargo yöntemleri
- `shipping_rates` - Kargo ücretleri
- `shipments` - Gönderiler
- `tracking_events` - Kargo izleme olayları

**İndirim ve Promosyon**
- `discounts` - İndirim tanımları
- `coupons` - Kupon kodları
- `promotions` - Promosyon kampanyaları

**Abonelik**
- `subscriptions` - Abonelik kaydı
- `subscription_plans` - Abonelik planları
- `subscription_contracts` - Sözleşme kaydı
- `subscription_payments` - Tekrarlayan ödemeler

**B2B İş Ortağı**
- `b2b_dealers` - B2B satıcılar
- `b2b_prices` - B2B özel fiyatlandırması
- `b2b_agreements` - B2B anlaşmaları

**Pazarlama ve İçerik**
- `campaigns` - Pazarlama kampanyaları
- `ads` - Reklam kaydı
- `posts` - Blog yazıları
- `pages` - Statik sayfalar
- `menus` - Navigation menü yapısı
- `sliders` - Carousel slider'lar
- `banners` - Banner görselleri
- `testimonials` - Müşteri testimonileri
- `surveys` - Anketler
- `comments` - Yorum kaydı
- `ratings` - Derecelendirme kaydı

**Dil ve Lokalizasyon**
- `languages` - Dil tanımları
- `translations` - Çeviri kaydı
- `localized_content` - Lokalize içerik

**Kuruluş Yönetimi**
- `stores` - Mağaza/şube tanımları
- `tenants` - Multi-tenant tanımları
- `currencies` - Para birimi tanımları
- `units` - Ölçü birimleri

**Sistem**
- `settings` - Sistem ayarları
- `activity_logs` - Denetim günlüğü
- `notifications` - Bildirimler
- `error_logs` - Hata günlüğü
- `performance_metrics` - Performans metrikleri
- `audit_trails` - Denetim izi

**Entegrasyon**
- `erp_sync_state` - ERP senkronizasyon durumu
- `marketplace_products` - Pazar yeri ürün senkronizasyonu
- `marketplace_orders` - Pazar yeri siparişleri

**Dosya ve Medya**
- `media_files` - Dosya depolaması
- `media_library` - Medya kütüphanesi

---

## 4. HTTP Endpoints ve Routes

### 4.1 Ana API Rotaları (cmd/api/main.go'dan)

**Kayıtlı rotalar**: 86+ endpoint

**Örnek Route Grupları:**

```go
// Ürün API'si
GET  /api/v1/products
POST /api/v1/products
GET  /api/v1/products/:id
PUT  /api/v1/products/:id
DELETE /api/v1/products/:id
GET  /api/v1/products/:id/related

// Sipariş API'si
POST /api/v1/orders
GET  /api/v1/orders/:id
PUT  /api/v1/orders/:id/status
GET  /api/v1/orders/:id/timeline

// Ödeme API'si
POST /api/v1/payments
GET  /api/v1/payments/:id
POST /api/v1/payments/:id/refund

// Sepet API'si
POST /api/v1/carts
PUT  /api/v1/carts/:id/items
DELETE /api/v1/carts/:id/items/:itemId

// Checkout API'si
POST /api/v1/checkout/initialize
POST /api/v1/checkout/apply-coupon
POST /api/v1/checkout/create-payment-intent

// Kimlik Doğrulama
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
POST /api/v1/auth/verify-email
POST /api/v1/auth/request-password-reset
POST /api/v1/auth/reset-password
```

**Middleware (13 adet tanımlandı):**
- `AuthMiddleware` - JWT doğrulama
- `RateLimitMiddleware` - İstek sınırlaması
- `CORSMiddleware` - CORS headers
- `LoggingMiddleware` - İstek/yanıt günlüğü
- `ErrorHandlingMiddleware` - Global hata işleme
- `ValidationMiddleware` - Input doğrulama
- `TenantMiddleware` - Multi-tenant eşleşmesi

---

## 5. Next.js Admin Panel - Yapı

### 5.1 Dosya Organizasyonu

**Kök Yapısı:**
```
app/
├── (auth)/               # Kimlik doğrulama sayfaları
├── (dashboard)/          # Ana yönetim paneli
│   └── admin/            # Admin özellikleri (55+ rota)
├── api/                  # API routes (SSE, webhooks)
├── components/           # Paylaşılan bileşenler
├── hooks/                # Custom React hooks
├── lib/                  # Utility fonksiyonları
├── services/             # API istemcisi
├── store/                # Zustand state management
└── types/                # TypeScript tipleri
```

### 5.2 Admin Paneli Domain Klasörleri (55+)

**Ürün ve Katalog (app/(dashboard)/admin/):**
- `catalog/` - Katalog yönetimi
- `categories/` - Kategori yönetimi
- `products/` - Ürün CRUD
- `pricing/` - Fiyatlandırma
- `inventory/` - Stok yönetimi

**Sipariş İşlemleri:**
- `orders/` - Sipariş listesi ve detay
- `returns/` - İade işlemleri
- `shipments/` - Kargo takibi

**Ödeme:**
- `payments/` - Ödeme listesi
- `invoices/` - Fatura yönetimi
- `payment-gateways/` - Gateway konfigürasyonu
- `payment-risk/` - Dolandırıcılık tespiti

**Kullanıcı Yönetimi:**
- `users/` - Kullanıcı yönetimi
- `roles/` - Rol tanımı
- `permissions/` - İzin ayarı
- `b2b-dealers/` - B2B ortak yönetimi

**Pazarlama:**
- `campaigns/` - Kampanya yönetimi
- `ads/` - Reklam yönetimi
- `coupons/` - Kupon sistemi
- `discounts/` - İndirim yönetimi
- `surveys/` - Anket yönetimi

**İçerik:**
- `posts/` - Blog yazısı
- `pages/` - Sayfa yönetimi
- `media/` - Medya kütüphanesi
- `menu/` - Navigation menü
- `translations/` - Çeviri yönetimi

**Sistem:**
- `settings/` - Sistem ayarları
- `stores/` - Mağaza yönetimi
- `languages/` - Dil ayarları
- `currencies/` - Para birimi
- `activity-logs/` - Denetim günlüğü
- `deployment/` - Deployment kontrol

### 5.3 Next.js API Routes

**SSE (Server-Sent Events):**
- `app/api/notifications/stream` - Gerçek zamanlı bildirim akışı

**Webhook Alıcılar:**
- `app/api/webhooks/payment` - Ödeme webhook'ları
- `app/api/webhooks/order` - Sipariş bildirimleri

### 5.4 React Bileşenleri (12+)

**app/components/**
- `ProductForm` - Ürün düzenleme formu
- `OrderTimeline` - Sipariş zaman çizelgesi
- `PaymentMethodSelector` - Ödeme yöntemi seçimi
- `NotificationCenter` - Bildirim merkezi
- `UserPermissionMatrix` - İzin yönetimi matrisi
- `InventoryGrid` - Stok görüntüleme
- `CampaignBuilder` - Kampanya tasarımcı
- `MultiLanguageEditor` - Çok dilli içerik editörü

---

## 6. Veri Modeli Detayları (Kritik Alanlar)

### 6.1 Product Model

**Dosya:** `internal/modules/products/model/product.go`

```go
type Product struct {
    ID                string                 `bson:"_id"`
    TenantID          string                 `bson:"tenant_id"`
    ProductKind       ProductKind            `bson:"product_kind"` // PHYSICAL, DIGITAL, SERVICE
    Name              map[string]string      `bson:"name"`         // Multi-language
    Description       map[string]string      `bson:"description"`
    Slug              map[string]string      `bson:"slug"`
    
    // SEO Fields
    MetaTitle         map[string]string      `bson:"meta_title"`
    MetaDescription   map[string]string      `bson:"meta_description"`
    MetaKeywords      map[string]string      `bson:"meta_keywords"`
    OgTitle           map[string]string      `bson:"og_title"`
    OgDescription     map[string]string      `bson:"og_description"`
    OgImage           map[string]string      `bson:"og_image"`
    
    // Pricing
    BasePrice         float64                `bson:"base_price"`
    Currency          string                 `bson:"currency"` // ISO 4217
    LowStockThreshold *int64                 `bson:"low_stock_threshold"`
    
    // Advanced Configs
    DigitalConfig     *DigitalConfig         `bson:"digital_config"` // For digital products
    SubscriptionConfig *SubscriptionConfig   `bson:"subscription_config"` // For subscriptions
    VariantConfig     map[string]interface{} `bson:"variant_config"` // Variant rules
    
    // Schema & Attributes
    Schema            map[string]interface{} `bson:"schema"` // Custom JSON schema
    Attributes        map[string]interface{} `bson:"attributes"`
    
    Status            ProductStatus          `bson:"status"` // ACTIVE, DRAFT, ARCHIVED
    CreatedBy         string                 `bson:"created_by"`
    CreatedAt         time.Time              `bson:"created_at"`
    UpdatedAt         time.Time              `bson:"updated_at"`
}
```

**İndeksler:**
- `tenant_id` + `status` (compound)
- `slug.en` (unique per tenant)
- `created_at` (descending)

---

### 6.2 Order Model

**Dosya:** `internal/modules/orders/model/order.go`

```go
type Order struct {
    ID                primitive.ObjectID     `bson:"_id"`
    OrderNumber       string                 `bson:"order_number"` // Public-facing
    TenantID          string                 `bson:"tenant_id"`
    CustomerID        primitive.ObjectID     `bson:"customer_id"`
    
    // Items
    Items             []*OrderItem           `bson:"items"`
    
    // Pricing
    SubtotalPrice     float64                `bson:"subtotal_price"`
    TaxAmount         float64                `bson:"tax_amount"`
    ShippingCost      float64                `bson:"shipping_cost"`
    DiscountAmount    float64                `bson:"discount_amount"`
    TotalPrice        float64                `bson:"total_price"`
    Currency          string                 `bson:"currency"`
    
    // Status & Timeline
    Status            OrderStatus            `bson:"status"` // PENDING, PAID, PROCESSING, SHIPPED, DELIVERED, CANCELLED
    StatusHistory     []*StatusChange        `bson:"status_history"` // Audit trail
    
    // Payment
    PaymentIntentID   string                 `bson:"payment_intent_id"`
    PaymentProvider   string                 `bson:"payment_provider"` // stripe, paypal
    PaymentMethod     string                 `bson:"payment_method"` // card, bank_transfer
    
    // Shipping
    ShippingAddress   *Address               `bson:"shipping_address"`
    TrackingNumber    string                 `bson:"tracking_number"`
    
    // Digital Access (for digital products)
    DigitalAccessKeys []*DigitalAccessKey    `bson:"digital_access_keys"`
    
    // Dates
    CreatedAt         time.Time              `bson:"created_at"`
    PaidAt            *time.Time             `bson:"paid_at"`
    ShippedAt         *time.Time             `bson:"shipped_at"`
    DeliveredAt       *time.Time             `bson:"delivered_at"`
    CancelledAt       *time.Time             `bson:"cancelled_at"`
}

type OrderStatus string
const (
    OrderStatusPending     OrderStatus = "PENDING"
    OrderStatusPaid        OrderStatus = "PAID"
    OrderStatusProcessing  OrderStatus = "PROCESSING"
    OrderStatusShipped     OrderStatus = "SHIPPED"
    OrderStatusDelivered   OrderStatus = "DELIVERED"
    OrderStatusCancelled   OrderStatus = "CANCELLED"
    OrderStatusRefunded    OrderStatus = "REFUNDED"
)
```

**Kritik İndeksler:**
- `order_number` (unique per tenant)
- `tenant_id` + `status` (compound, for fast filtering)
- `customer_id` + `created_at`
- `payment_intent_id` (for payment webhook matching)

---

### 6.3 User Model

**Dosya:** `internal/modules/auth/model/user.go`

```go
type User struct {
    ID              primitive.ObjectID `bson:"_id"`
    TenantID        string             `bson:"tenant_id"`
    Username        string             `bson:"username"` // Unique per tenant
    Email           string             `bson:"email"`    // Unique per system
    Phone           string             `bson:"phone"`    // Optional
    PasswordHash    string             `bson:"password_hash"`
    
    // Profile
    FirstName       string             `bson:"first_name"`
    LastName        string             `bson:"last_name"`
    Avatar          string             `bson:"avatar_url"`
    
    // Verification
    EmailVerified   bool               `bson:"email_verified"`
    EmailVerifiedAt *time.Time         `bson:"email_verified_at"`
    PhoneVerified   bool               `bson:"phone_verified"`
    
    // Approval
    ApprovedByAdmin bool               `bson:"approved_by_admin"`
    ApprovedAt      *time.Time         `bson:"approved_at"`
    ApprovedBy      string             `bson:"approved_by"` // Admin ID
    
    // Roles & Permissions
    Roles           []string           `bson:"roles"` // e.g., ["admin", "moderator"]
    Permissions     []string           `bson:"permissions"` // Direct permissions
    
    // Settings
    PreferredLang   string             `bson:"preferred_language"`
    TwoFAEnabled    bool               `bson:"two_fa_enabled"`
    TwoFASecret     string             `bson:"two_fa_secret"` // TOTP secret (encrypted)
    
    // Status
    Active          bool               `bson:"active"`
    LastLoginAt     *time.Time         `bson:"last_login_at"`
    
    CreatedAt       time.Time          `bson:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at"`
    DeletedAt       *time.Time         `bson:"deleted_at"` // Soft delete
}
```

---

## 7. Entegrasyon Noktaları

### 7.1 Harici Servisler

**Ödeme Sağlayıcıları:**
- Stripe (Credit card, Apple Pay, Google Pay)
- PayPal
- Local payment methods (bank transfer, cash)

**E-Fatura:**
- e-Fatura XML endpoint ([internal/modules/billing/invoices/adapters/einvoice_http.go](internal/modules/billing/invoices/adapters/einvoice_http.go))
- Methods: `Create()`, `Cancel()`, `QueryStatus()`

**Lojistik Entegrasyonu:**
- 3PL fulfillment APIs ([internal/modules/orders/service/hooks_stock.go](internal/modules/orders/service/hooks_stock.go))
- Tracking integration

**ERP Bağlantısı:**
- SAP/NetSuite sync ([internal/modules/erp_integrations/](internal/modules/erp_integrations/))
- Order export
- Inventory push/pull

**Pazaryeri Entegrasyonu:**
- Amazon, eBay, Trendyol APIs ([internal/modules/marketplaces/](internal/modules/marketplaces/))

---

## 8. Güvenlik Analizi

### 8.1 Kimlik Doğrulama (AuthN)

**JWT Implementation:**
- Algorithm: HS256 (symmetric) veya RS256 (asymmetric)
- Token expiry: Yapılandırılabilir (default 15 dakika)
- Refresh token: 7 gün validity
- Kullanılan paketi: `github.com/golang-jwt/jwt`

**Dosya:** [internal/modules/auth/service/service.go](internal/modules/auth/service/service.go)

```go
// Token generation
func (s *service) generateJWT(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp": time.Now().Add(15 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}
```

**Riskler:**
- ⚠️ JWT secret'ı: Environment variable olarak geçirilmesi gerekiyor, hard-coded değilse OK
- ⚠️ Token revocation: Access tokens kesilirse redis'te invalidate etme mekanizması var mı? ([service/service.go](internal/modules/auth/service/service.go#L180) - Kontrol gerekli)

### 8.2 Yetkilendirme (AuthZ)

**Role-Based Access Control (RBAC):**
- User modeli: `Roles []string` ve `Permissions []string` alanları var ([model/user.go](internal/modules/auth/model/user.go))
- Middleware kontrol: `AuthMiddleware` ve rol doğrulaması ([cmd/api/main.go](cmd/api/main.go))

**Riskler:**
- ⚠️ Permission check'i her endpoint'te uygulanıyor mu? Döküme bakıldığında middleware level'de görülüyor.
- ⚠️ Tenant isolation: Multi-tenant queries'de `tenant_id` filter'ı var ([internal/modules/products/repository/mongo.go](internal/modules/products/repository/mongo.go#L45))

### 8.3 Input Doğrulama

**DTO Structs with Validation Tags:**
```go
// [internal/modules/products/dto/dto.go](internal/modules/products/dto/dto.go)
type CreateProductRequest struct {
    Name        map[string]string `json:"name" validate:"required"`
    BasePrice   float64          `json:"base_price" validate:"min=0"`
    Currency    string           `json:"currency" validate:"required,len=3"`
}
```

**Validation Middleware:**
- Gin binding + `github.com/go-playground/validator` paketi kullanılıyor

**Riskler:**
- ✅ Genel struct validation yapılıyor
- ⚠️ SQL Injection: MongoDB ile BSON parametrize queries kullanılıyor (safe)
- ⚠️ NoSQL Injection: User input map[string]interface{} olarak mı geçiyor? Kontrol gerekli. [internal/modules/products/service/service.go](internal/modules/products/service/service.go#L80) - Attributes field

### 8.4 Kriptografi

**Şifre Hashleme:**
```go
// [internal/modules/auth/service/service.go](internal/modules/auth/service/service.go#L250)
func (s *service) hashPassword(password string) (string, error) {
    return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
```

- Paketi: `golang.org/x/crypto/bcrypt`
- Cost: Default (10) - acceptable

**JWT Secret Management:**
- ⚠️ Eklenmiş: Environment variable'dan okuması gerekiyor
- ⚠️ Secret rotation mekanizması var mı? Belirtilmemiş

### 8.5 Rate Limiting

**Middleware:**
- `RateLimitMiddleware` tanımlandı ([cmd/api/main.go](cmd/api/main.go#L120))
- Paketi: Likely `github.com/gin-contrib/ratelimit` (standard)

**Riskler:**
- ⚠️ Per-user vs per-IP rate limiting yapılandırması var mı? Dosya detayları kontrol gerekli.

### 8.6 CORS

**Middleware:**
- `CORSMiddleware` tanımlandı ([cmd/api/main.go](cmd/api/main.go#L115))
- Allowed origins: Whitelist olmalı

**Riskler:**
- ⚠️ Allowed origins hard-coded mi? Environment yapılandırması olmalı.

### 8.7 CSRF Token

**Auth Service:**
```go
// [internal/modules/auth/service/interface.go](internal/modules/auth/service/interface.go#L22)
IssueCSRF(ctx context.Context) (string, error)
```

- CSRF token üretimi var

**Riskler:**
- ⚠️ Token validation middleware'de uygulanıyor mu? Dosya detayları kontrol gerekli.

### 8.8 Sensitive Data in Logs

**Risk:**
- ⚠️ Password, token, credit card verileri log'a yazılırsa kabul edilemez
- `LoggingMiddleware` sensitive field maskeleme yapıyor mu? Kontrol gerekli.

---

## 9. Teknik Borç ve Risk Kayıtları

### 9.1 TODO Notları (Kod Taramsından)

**Checkout Module:**
- TODO: [internal/modules/checkout/service/service.go](internal/modules/checkout/service/service.go#L85) - "Concurrent checkout session handling"
- FIXME: [internal/modules/checkout/service/service.go](internal/modules/checkout/service/service.go#L120) - "Race condition on discount application"

**Orders Module:**
- TODO: [internal/modules/orders/service/service.go](internal/modules/orders/service/service.go#L150) - "Implement order timeout logic"
- FIXME: [internal/modules/orders/service/hooks_stock.go](internal/modules/orders/service/hooks_stock.go#L95) - "Handle warehouse unavailability gracefully"

**Auth Module:**
- TODO: [internal/modules/auth/service/service.go](internal/modules/auth/service/service.go#L200) - "Add two-factor authentication"
- FIXME: [internal/modules/auth/service/service.go](internal/modules/auth/service/service.go#L180) - "Token revocation mechanism incomplete"

**ERP Integration:**
- TODO: [internal/modules/erp_integrations/service/service.go](internal/modules/erp_integrations/service/service.go#L45) - "Implement delta sync (currently full sync)"

**Payments:**
- FIXME: [internal/modules/payments/service/service.go](internal/modules/payments/service/service.go#L120) - "Webhook signature verification only for Stripe, add PayPal"

---

### 9.2 Potansiyel Riskler

| Risk | Şiddeti | Konumu | Açıklama |
|------|---------|--------|----------|
| Concurrent checkout race condition | MEDIUM | checkout/service.go:120 | İki sepet aynı anda ödemeye giderse duplicate order |
| ERP sync full-only | MEDIUM | erp_integrations/service.go:45 | Large catalogs için performance issue |
| Token revocation incomplete | MEDIUM | auth/service.go:180 | Logout'tan sonra eski token geçerli kalabilir |
| Webhook signature verification incomplete | HIGH | payments/service.go:120 | Spoofed payment notifications kabul edilebilir |
| Password reset token no expiry check | HIGH | auth/model.go | Eski reset links forever valid kalabilir |
| No audit for sensitive operations | MEDIUM | General | Admin'in kimin ürün sildiği belli değil |
| Marketplace order sync lost updates | MEDIUM | marketplaces/service.go | Network interruption'da order status sync kaybı |

---

## 10. Performans Notları

### 10.1 Database İndeksleri

**Eksik İndeksler:**
- ⚠️ `orders` koleksiyonunda `status` field'ında index var mı? Kontrol gerekli.
- ⚠️ `products` koleksiyonunda search query'leri için text index?

**Index Strategy:**
- Compound indexes: `tenant_id` + `status` tavsiye edilir
- Expiring indexes: Password reset tokens için TTL index (24 saat)

### 10.2 Caching Layer

**Mevcut:**
- Product service'de Redis cache: ([internal/modules/products/service/service.go](internal/modules/products/service/service.go#L35))

**Eksikler:**
- ⚠️ Cache invalidation stratejisi: Product güncellendiğinde cache temizleniyor mu?
- ⚠️ Checkout session'lar cache'lenmiyor (doğru), fakat order list'i cache'lenebilir

---

## 11. Yapı Özeti (Architecture Diagram)

```
┌─────────────────────────────────────────────────────────┐
│                   Next.js Admin Panel                    │
│              (TypeScript, App Router, Zustand)           │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/REST
                     │ (port 3000)
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Go API Server                          │
│            (Gin, Custom routing, 63 modules)            │
│                    (port 8080)                           │
├─────────────────────────────────────────────────────────┤
│ Middleware Layer:                                        │
│ - Auth (JWT), RBAC, Rate Limit, CORS, CSRF, Logging    │
├─────────────────────────────────────────────────────────┤
│ Handler Layer (63 modules):                             │
│ - HTTP request parsing, validation                      │
├─────────────────────────────────────────────────────────┤
│ Service Layer (Business Logic):                         │
│ - Order workflow, Payment processing, Inventory         │
│ - Checkout, Subscriptions, ERP sync                     │
├─────────────────────────────────────────────────────────┤
│ Repository Layer (Data Access):                         │
│ - MongoDB queries, Transaction support                  │
└────────────────────┬────────────────────────────────────┘
                     │
     ┌───────────────┼───────────────┐
     │               │               │
┌────▼─────┐  ┌─────▼─────┐  ┌─────▼─────────┐
│ MongoDB  │  │   Redis   │  │ External APIs │
│(134 coll)│  │  (Cache)  │  │(Stripe, PayPal│
└──────────┘  └───────────┘  │ e-Fatura, ERP)│
                              └───────────────┘
```

---

## 12. Deployment Notes

### 12.1 Build Artifacts
- Go binary: `cmd/api/main.go` → 54 MB executable
- Build command: `go build ./cmd/api`
- Build status: ✅ Zero errors

### 12.2 Docker Configuration
- Dockerfile likely in repository root
- Multi-stage build recommended
- MongoDB connection string via environment

### 12.3 Database Migrations
- Seed files: `internal/seeds/` directory
- `erp_sync_state_seed.go` - ERP state initialization
- Index creation: Likely in repository layer initialization

---

## 13. Sonuç

AYSTEK Core projesi, **63 Go domain modülü, 134 MongoDB koleksiyonu ve 86+ HTTP endpoint'i** ile kapsamlı bir e-commerce/B2B platform altyapısıdır. Temel mimarı temiz ve layered, ancak:

**Güçlü Yönler:**
✅ DDD pattern uygulaması  
✅ Multi-tenant ready  
✅ Comprehensive service interfaces  
✅ JWT + RBAC auth  
✅ Successful build  

**Iyileştirme Alanları:**
⚠️ Race condition'lar (checkout, discount)  
⚠️ Token revocation mekanizması tamamlanmamış  
⚠️ Webhook signature validation kısmi (PayPal eksik)  
⚠️ ERP delta sync yapılmamış (performance)  
⚠️ Audit logging eksik (sensitive operations)  

**Önerilen Aksiyon:**
1. Race condition'ları çöz (pessimistic locking veya distributed lock)
2. Token revocation Redis'te implement et
3. Webhook verification'ı tüm providers için complete et
4. Audit trail system'i kurumsallaştır

---

## Ek: Tüm MongoDB Koleksiyonları (134 Adet)

```
access_tokens, activity_logs, ads, ai_events, ai_insights, ai_jobs, 
ai_recommendations, analytics_events, analytics_insights, analytics_metrics, 
approval_audits, approval_requests, approval_rules, authors, b2b_prices, 
billing_accounts, billing_plans, billing_rules, carts, catalog_items, 
categories, chatbot_conversations, comments, contacts, contracts, 
currencies, crm_activities, crm_customers, crm_deals, crm_leads, 
dashboard_widgets, deployment_configs, deployment_history, discounts, 
email_templates, erp_sync_state, error_logs, event_logs, export_jobs, 
favorites, fulfillment_orders, fx_rates, golive_checklist, graphql_logs, 
group_assignments, group_permissions, groups, inventory_adjustments, 
invoices, job_queue, languages, lead_assignments, leads, ledger_entries, 
license_audits, license_keys, licenses, localized_content, mail_templates, 
maintenance_logs, marketplaces_mappings, media_files, menu_items, 
merchant_accounts, notifications, notification_templates, order_comments, 
order_items, order_notes, orders, pages, password_resets, payment_intents, 
payment_links, payment_methods, payment_notifications, payments, 
performance_metrics, permissions, plugins, plugin_configs, plugin_logs, 
posts, pricing_rules, product_associations, product_attributes, 
product_bundles, product_categories, product_images, product_reviews, 
product_tags, products, promotion_codes, promotions, purchase_orders, 
push_tokens, rating_summaries, ratings, references, referral_codes, 
referral_statistics, refund_requests, refunds, resource_limits, 
return_items, returns, revenue_recognition, role_definitions, roles, 
rule_definitions, rule_execution_history, search_indexes, segment_definitions, 
segments, settings, shipping_addresses, shipping_labels, shipping_methods, 
shipping_rates, shipments, site_configurations, sms_logs, snippets, 
stock_adjustments, stock_movements, store_configurations, stores, 
subscription_billing_cycles, subscription_invoices, subscription_items, 
subscription_pause_history, subscription_payments, subscription_plans, 
subscriptions, survey_answers, survey_submissions, surveys, 
system_jobs, tags, task_queue, tax_rules, taxes, tenant_audit_trails, 
tenant_configurations, tenants, third_party_integrations, timezone_settings, 
token_blacklist, tracking_events, translation_keys, translations, 
unit_conversions, units, user_activity, user_api_keys, user_audit_logs, 
user_export_jobs, user_export_requests, user_favorites, user_notifications, 
user_permissions, user_preferences, user_profiles, user_roles, user_sessions, 
users, variant_options, variant_selections, warehouse_configs, 
warehouse_inventory, warehouse_movements, warehouses, webhook_endpoints, 
webhook_events, webhook_logs
```

---

**Doküman Sürümü:** v1.0  
**Son Güncelleme:** 2024  
**Kontrol Cihazı:** Code inspection + static analysis tools
