# AYSTEK Core Arac Takip ve Filo Telematics Entegrasyon Is Planlari

## Amac

Arac, filo ve nesne takip yeteneklerini AYSTEK Core icine tenant-aware, site-aware ve module-gated bir yapiyla entegre etmek; mevcut ortak modulleri dogru yerde kullanmak; yanlis ownership nedeniyle duplicate domain uretmemek.

Bu dokuman, AYSTEK Core genel mimarisi ve mevcut aktif moduller incelenerek revize edilmistir. Ozellikle:

- `integrations` platformun gercek entegrasyon orkestrasyon katmanidir.
- `stock/integrations` dar ve ikinci planda kalan bir yapidir; stratejik baz olmamalidir.
- `stock/locations` bugun aktif urun yuzeyidir, ancak domain owner olmamali; telematics projection/read-model katmanina donusmelidir.
- `stock/vehicles` ve `stock/vehicle_assets` aktif master-data ve operasyonel stock akislarina baglidir; aceleyle tasinmamali, kontrollu ownership matrisiyle korunmalidir.

## Incelenen AYSTEK Core Yuzeyleri

### Cekirdek platform ve gate zinciri

- Router middleware zinciri tenant, site, request context, module gate ve permission enforcement icin dogru kurulmus durumda.
- Module gate request path uzerinden modulu cozuyor; yeni moduller mevcut standarda uymali.
- Module enablement `module.{module}.enabled`, `tenant.{tenant}.module.{module}.enabled`, `site.{site}.module.{module}.enabled` key kalibiyla yonetiliyor.

### Entegrasyon platformu

- `integrations` modulu connector, action, hook, chain, execution, DLQ, pull schedule/run, shared service, grant ve internal worker akislariyla enterprise-level platform.
- `integrations/eventbridge` domain event -> hook orkestrasyonunu zaten sagliyor.
- `integrations/api_sharing` partner API exposure icin ayri ve zengin bir yetenek sagliyor.
- Frontend tarafinda `admin/integrations/*` yuzeyi dolu ve aktif.

### Stock ekosistemi

- `stock/vehicles` aktif arac kaydi ve plaka gecmisi moduludur.
- `stock/inventory` GPS, HAT, kamera ve benzeri fiziksel cihaz envanterini zaten tasiyor.
- `stock/vehicle_assets` arac-cihaz montaj ve yasam dongusunu tasiyor; event de uretiyor.
- `stock/sales`, `stock/returns`, `stock/reporting` gibi moduller vehicle-assets ile bagli; bu baglar domain tasima kararini etkiliyor.
- `stock/locations` backend, menu, admin UI ve portal UI tarafinda aktif kullaniliyor.

### PMS ve diger ortak moduller

- Alarm, ariza, bakim veya operasyonel gorev ihtiyacinda `pms/work_orders` yeniden kullanilabilir ve yeniden yazilmamalidir.
- Bildirim cikislari icin `notifications`, `sms`, `mailtemplates` ortak servis olarak korunmali.
- Audit ve takibin tek merkezi `activity_logs` olmaya devam etmeli.

## Dogru Lemma

AYSTEK Core icinde arac takip/filo takip cozumunun dogru bounded context dagilimi su olmali:

1. `integrations`:
   Telematics vendor, OEM, gateway, webhook, scheduled pull, partner API ve event-hook orkestrasyonunun tek stratejik giris/cikis katmani.
2. `stock`:
   Arac, fiziksel cihaz, montaj, satis, iade ve envanterin master-data ve ticari/lojistik operasyon katmani.
3. `telematics`:
   Canli telemetri, konum, trip, geofence, alert, driver behavior, command, CAN/OBD, replay ve fleet operasyonlarinin yeni bounded context'i.
4. `stock/locations`:
   Hedef durumda write-owner degil, telematics tarafindan beslenen projection/read-model veya compatibility facade.

Bu lemma ile hem bugunku urun gercegi korunur hem de yeni telematics domain'i stock icine gomulmeden kurulmus olur.

## Mimari Karar Ozeti

1. Yeni ana bounded context `telematics` olacaktir.
2. `integrations` telematics connector orkestrasyonu icin tek resmi platform olacak; `stock/integrations` ustune yeni stratejik kabiliyet eklenmeyecektir.
3. `stock/vehicles` ilk fazlarda master-data olarak korunacaktir.
4. `stock/inventory` fiziksel tracker, HAT, kamera, sensor, gateway, OBD dongle gibi ekipmanlarin envanter owner'i olarak korunacaktir.
5. `stock/vehicle_assets` arac-cihaz mount/unmount ve ownership resolution icin ana bridge olarak korunacaktir.
6. `stock/locations` aktif UI nedeniyle aniden kaldirilmamalidir; telematics write-modeli devreye alindiktan sonra projection veya facade'a donusturulmelidir.
7. Bakim, saha gorevi ve alarmdan is emri uretimi icin yeni bir gorev modulu yazilmayacak; mevcut `pms/work_orders` reuse edilecektir.
8. Tum yeni telematics modulleri tenant/site-aware, permissioned, seeded ve backend-driven menu standardina gore kurulacaktir.

## Ownership Matrisi

### `stock` tarafinda kalacak alanlar

- Arac kimligi
- Plaka gecmisi
- Fiziksel cihaz envanteri
- SIM/HAT/kamera/GPS/OBD dongle gibi envanter kalemleri
- Araca takili cihaz iliskisi
- Satis, iade, stok, assignment ve ticari operasyon akislari

### `telematics` tarafina alinacak alanlar

- Device provisioning durumu
- Protocol adapter metadata
- Heartbeat ve connection health
- Raw payload ve normalized telemetry
- Current position canonical state
- Route, trip, stop, idle, ignition state
- Geofence ve geofence eventleri
- Alert rule, alert instance, escalation, ack
- Driver session ve behavior score
- CAN/OBD normalized veri
- Remote command queue ve command result
- Replay, playback, fleet dashboard cache ve operational reporting

### `integrations` tarafinda kalacak alanlar

- Vendor connector governance
- Secret-aware auth config
- Webhook/hook delivery
- Scheduled pull orchestration
- Event-driven outbound integration
- Shared service ve partner API exposure
- Retry, DLQ, tracing, rate limiting, circuit breaker

## `stock/integrations` Icin Karar

`stock/integrations` bugun teknik olarak kayitli olsa da:

- urun yuzeyinde merkezi degildir,
- genel `integrations` kadar zengin governance ve worker altyapisina sahip degildir,
- frontend tarafinda birinci sinif admin deneyimiyle temsil edilmemektedir.

Bu nedenle:

1. Yeni telematics connector/yetki/webhook/schedule/DLQ kabiliyetleri `stock/integrations` altinda gelistirilmeyecektir.
2. Telematics ingestion ve vendor entegrasyonlari `integrations` modulu uzerinden tanimlanacaktir.
3. Ilerleyen fazda `stock/integrations` icin deprecation veya migration plani yazilacaktir.
4. Gerekirse gecis doneminde sadece mevcut veri veya backward compatibility nedeniyle okunacak; stratejik owner olmayacaktir.

## `stock/locations` Icin Karar

`stock/locations` su an aktif oldugu icin atil kabul edilmemelidir. Ancak domain owner olarak stock altinda kalmasi dogru degildir; cunku canli konum ve hareket verisi stock degil telematics problemidir.

Bu nedenle hedef model:

1. Canonical current location owner `telematics/telemetry` olacaktir.
2. Historical playback ve route sorgulari `telematics` tarafinda tutulacaktir.
3. `stock/locations` ilk gecis fazinda compatibility read-model olarak yasayacaktir.
4. Mevcut `/admin/stock/locations/*` ve `/portal/stock/locations/current` yuzeyleri gecici olarak korunacak, ama telematics projection'dan beslenecektir.
5. Uygun fazda yeni `admin/telematics/locations`, `admin/telematics/replay`, `admin/telematics/trips` UI'lari acilacak; eski stock location sayfalari yonlendirme veya compatibility facade'a indirgenecektir.

## `stock/vehicles` Icin Karar

Arac kaydini telematics'e tasimak teorik olarak mumkun olsa da mevcut urun baglantilari nedeniyle kisa ve orta vadede dogru karar degildir. Neden:

- `stock/vehicle_assets` araca bagli envanter mantigi kuruyor.
- `stock/sales` ve `stock/reporting` tarafinda arac-cihaz akislari var.
- mevcut menu ve operasyon algisi stock envanteri ile arac kaydini birlikte sunuyor.

Bu nedenle karar:

1. Faz 1-Faz 4 araliginda `stock/vehicles` master-data owner olarak korunur.
2. `telematics` modulu arac verisini repository/service kontrati ile referanslar.
3. Ancak arac domaini ileride stock'tan daha fazla operasyonel davranis alirsa, Faz 6+ sonrasinda ayrica bir `fleet/assets` veya `fleet/vehicles` re-segmentation ADR'i yazilabilir.

Yani bugunku dogru karar: hemen tasimamak.

## Hedef Domain Topolojisi

### Mevcut modullerden reuse edilecekler

- `stock/vehicles`
- `stock/inventory`
- `stock/vehicle_assets`
- `stock/catalog`
- `integrations`
- `integrations/api_sharing`
- `notifications`
- `sms`
- `mailtemplates`
- `activity_logs`
- `pms/work_orders`
- `settings`
- `permissions`
- `menus`
- `translations`

### Yeni acilacak telematics modulleri

- `telematics/devices`
- `telematics/ingestion`
- `telematics/telemetry`
- `telematics/trips`
- `telematics/geofences`
- `telematics/alerts`
- `telematics/drivers`
- `telematics/commands`
- `telematics/canbus`
- `telematics/dashboard`
- `telematics/reports`

Not:
`telematics/webhooks` diye ayri bir owner modul acmak yerine once `integrations/hooks` ve `integrations/api_sharing` reuse edilmelidir. Gerekirse telematics sadece bu platformlara domain-level adapter saglar.

## Hedef Veri Akisi

1. Vendor, cihaz veya gateway veriyi `integrations` uzerinden HTTP webhook, scheduled pull veya adapter pipeline'a yollar.
2. `telematics/ingestion` payload'i dogrular, normalize eder ve source metadata ile kaydeder.
3. Device resolution, `stock/inventory` ve `stock/vehicle_assets` uzerinden yapilir.
4. Device -> mounted asset -> vehicle zinciri cozulur.
5. Telemetry canonical modele yazilir.
6. Current location state `telematics/telemetry` tarafinda guncellenir.
7. Compatibility gerekiyorsa ayni veri `stock/locations` projection'una yazilir.
8. Trip assembler, geofence engine, alert engine ve driver scoring yan akislar halinde calisir.
9. Bildirimler `notifications/sms/mailtemplates` ile gonderilir.
10. Gerekli operasyonlar `pms/work_orders` ile goreve donusturulur.
11. Dis sistemlere acilim `integrations/hooks`, `integrations/shared-services` veya `integrations/api_sharing` ile yapilir.

## Fazlandirilmis Enterprise Is Plani

### Faz 0 - Kurumsal Kesinlestirme ve ADR Seti

- Tum aktif stock bagimliliklari cikarilacak.
- `integrations` ve `stock/integrations` arasinda resmi karar ADR'i yazilacak.
- `stock/locations` icin compatibility ve migration ADR'i yazilacak.
- `stock/vehicles` ownership karari resmi ADR ile sabitlenecek.
- Modul key, permission matrix, menu tree ve translation namespace tasarlanacak.
- Canonical event naming standard'i cikarilacak.

### Faz 1 - Telematics Foundation

- `telematics` kok modul ve alt paket yapisi acilacak.
- `Device`, `TelemetryPoint`, `CurrentState`, `Trip`, `Geofence`, `Alert`, `DriverSession`, `Command` aggregate'leri tanimlanacak.
- Tum repository ve modellerde tenant/site guard ve index standardi uygulanacak.
- Vehicle/device resolution domain service yazilacak.
- `stock_vehicle_assets` aktif mount kaydi owner kabul edilerek resolution zinciri sabitlenecek.

### Faz 2 - Integration-Driven Ingestion

- `integrations/connectors` icinde telematics vendor connector tipleri ve capability convention'lari tanimlanacak.
- Webhook ve scheduled pull ingestion adapter'lari gelistirilecek.
- Idempotency, replay protection, malformed payload handling ve DLQ akislari devreye alinacak.
- Raw payload audit ve correlation ID zinciri kurulacak.
- `stock/integrations` uzerinden yeni akış yazilmamali; varsa yalniz gecis adaptoru kalmali.

### Faz 3 - Device Registry ve Health

- `telematics/devices` modulunde logical device registry yazilacak.
- Device status, heartbeat, firmware, protocol, SIM relation, last seen, signal health alanlari eklenecek.
- Device ile stock envanteri arasinda canonical relation tablosu/aggregate'i tanimlanacak.
- Unmounted veya sold state'lerde davranis kurallari netlestirilecek.

### Faz 4 - Telemetry Core ve Konum Migration

- Raw telemetry ve normalized telemetry persistence yazilacak.
- Current location canonical state kurulacak.
- Historical telemetry query, replay ve basic playback altyapisi kurulacak.
- `stock/locations` projection writer eklenecek.
- Mevcut stock location endpointleri projection kaynagini yeni telematics write-modelden okuyacak sekilde evrilecek.

### Faz 5 - Fleet Intelligence

- Trip assembler
- Stop/idle/ignition state machine
- Geofence CRUD ve event generation
- Overspeed, offline, panic, tow, route deviation gibi alert kurallari
- Ack, suppression, escalation ve noise control mekanizmalari

### Faz 6 - Driver, CAN/OBD ve Operasyonel Derinlik

- Driver identity, driver-session ve assignment modeli
- RFID, iButton, mobil login veya QR tabanli tanimlama
- Driver behavior scoring
- CAN/OBD normalized schema
- Fuel, odometer, DTC, RPM, engine temp, harsh event alanlari
- Decoder strategy registry

### Faz 7 - Remote Command ve Enterprise Acilim

- Command queue
- ACK/NACK/result tracking
- Immobilizer ve kritik komutlar icin policy, approval ve audit
- Event-driven outbound integrations
- Partner ve bayi ekosistemi icin API sharing exposure
- Signed webhook delivery ve consumption contract'lari

### Faz 8 - UI, Raporlama ve Geçis Temizligi

- `admin/telematics/*` bilgi mimarisi ve menu yapisi
- Fleet dashboard
- Vehicle 360
- Device health center
- Alert center
- Replay/trip/geofence ekranlari
- `stock/locations` uyumluluk ekranlarinin sadeleştirilmesi veya yonlendirilmesi
- `stock/integrations` deprecation ve migration cleanup

Faz 8 kapanis referanslari:
- [FAZ8_STOCK_INTEGRATIONS_MIGRATION_INVENTORY.md](/workspace/aystek-core-dev-docs/plans/FAZ8_STOCK_INTEGRATIONS_MIGRATION_INVENTORY.md)
- [FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md](/workspace/aystek-core-dev-docs/plans/FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md)

## Sprint Bazli Backlog

### Sprint 1

- ADR seti
- module key ve permission matrisi
- `telematics` backend skeleton
- event naming convention
- resolution service prototype

### Sprint 2

- integration-driven ingestion
- raw payload persistence
- device registry temel modeli
- test doubles ve integration test skeleton

### Sprint 3

- normalized telemetry
- current state
- stock locations projection
- current/history API gecis katmani

### Sprint 4

- trips
- geofences
- temel alert engine
- notifications bridge

### Sprint 5

- drivers
- CAN/OBD
- behavior scoring
- fleet dashboard v1

### Sprint 6

- commands
- partner API/webhook
- operational hardening
- deprecation cleanup backlog

## Teknik Kurallar

- Her yeni koleksiyon `tenant_id` ve `site_id` ile indexlenecek.
- Tum admin endpointleri module gate ve permission middleware ile korunacak.
- Menu, permission, translation ve settings seed'leri eksiksiz gelecek.
- Cross-context erisim dogrudan koleksiyon okumayla degil service/repository kontratlariyla yapilacak.
- `stock/locations` artik write-owner gibi buyutulmeyecek.
- `integrations` disina ikinci bir connector orkestrasyon platformu acilmayacak.
- Event publishing ve event-hook orkestrasyonu mevcut `platform/events` ve `integrations/eventbridge` standartlariyla uyumlu olacak.
- Yuksek hacimli telemetry icin cold/hot retention ve aggregation stratejisi en bastan dusunulecek.

## Kritik Refactor Alanlari

1. `stock/locations`:
   Projection/facade modeline cekilecek.
2. `stock/inventory.product_type`:
   Telematics ekipman tiplerini daha net ifade edecek sekilde genisletilecek.
3. `stock/vehicle_assets`:
   Mount role gibi alanlarla zenginlestirilebilir.
4. `integrations`:
   Telematics vendor capability convention'lari ve reusable templates eklenebilir.
5. `stock/integrations`:
   Kullaniliyorsa migration matrisi yazilacak; kullanilmiyorsa deprecation adimi planlanacak.

## Riskler

- Stock ile telematics arasinda ownership bulanikligi
- Active UI varken `stock/locations` migration'inda kesinti
- Vendor payload cesitliligi nedeniyle normalization karmasasi
- Device -> inventory -> vehicle resolution bosluklari
- High-volume telemetry yazim maliyeti
- Alarm firtinasi ve bildirim patlamasi
- Multi-tenant veri sizmasi

## Risk Azaltma

- Tek bir resmi ownership matrisi
- ADR ile sabitlenmis migration kararleri
- Projection pattern
- Idempotency ve anti-replay
- DLQ ve retry policy
- Alert suppression ve escalation discipline
- Repository katmaninda tenant/site enforcement
- End-to-end integration test matrisi

## Basari Kriterleri

- Vendor entegrasyonlari yalniz `integrations` platformu uzerinden yonetilir.
- Arac ve cihaz kayitlari duplicate olmadan tek master-data zinciriyle yonetilir.
- `stock/locations` aktif urun yuzeyi kirilmadan telematics projection modeline gecirilir.
- Canli konum, gecmis, geofence, trip, alert ve replay yetenekleri `telematics` bounded context'inde tutarli calisir.
- UI ve permission yapisi backend-driven standardi bozmaz.
- Yeni moduller seed, translation, permission, audit ve test standartlarina tam uyar.
