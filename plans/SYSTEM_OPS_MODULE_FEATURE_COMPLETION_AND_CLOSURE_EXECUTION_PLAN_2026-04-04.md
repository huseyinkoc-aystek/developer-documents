# System Ops Module Feature Completion and Closure Execution Plan

Date: 2026-04-04
Status: Completed
Owner: Codex

## Objective

System Ops modülünün hardening sonrası kalan gerçek capability boşluklarını kapatıp modülü “UI-polished but backend-stubbed” seviyesinden çıkararak operasyonel olarak tamamlanmış bir modül haline getirmek.

Bu planın amacı:
- hardening backlog ile feature backlog’u ayırmak
- kapanış için gerekli minimum gerçek işlevselliği netleştirmek
- Go backend ile Next admin yüzeyini aynı teslim kriterinde hizalamak

## Verified Current-State Assessment

Hardening tarafı büyük ölçüde tamamlandı:
- shared component stack
- ortak tasarım dili
- i18n coverage
- helper consolidation
- legacy wrapper cleanup

Buna karşılık aşağıdaki capability boşlukları hâlâ açık:

### 1. Update orchestration gerçek değil

Backend:
- `/admin/system/updates/check` hâlâ stub döndürüyor
- `StartUpdateRequest` içindeki `tarball_url` ve `confirm_token` handler/service içinde kullanılmıyor
- Next tarafındaki `options` sözleşmesi backend’e karşılık gelmiyor
- update manifest / release channel / health gate / preflight doğrulama zinciri yok

Frontend:
- update start ekranı sadece `version + channel` ile sınırlı
- dry-run / pre-backup / maintenance gate / auto rollback gibi opsiyonlar gerçek yüzey değil

### 2. Rollback akışı operasyonel olarak eksik

Backend:
- rollback listesi update listesinden page sonrası filtreleniyor; `total` ve pagination doğru değil
- `target_version` boş ise “previous version” çözümü yapılmıyor
- `confirm_token` tanımlı ama işlenmiyor

Frontend:
- “previous version” / önerilen target / rollback safety preview yok

### 3. Migration flow stub seviyesinde

Backend:
- `TriggerMigrations` yalnızca audit entry yazıyor; gerçek migration execution yok
- migration progress, output, failure detail ve rollbackability akışı yok

Frontend:
- migration ekranı liste odaklı; gerçek run/progress/result operasyon yüzeyi değil

### 4. Backup orchestration eksik

Backend:
- `TriggerBackupRequest.note` kullanılmıyor
- backup trigger her zaman `manual` tipinde başlatılıyor
- çalışan backup’ı tamamlayan worker/pipeline yok
- restore / download / retention yönetimi yok

Frontend:
- backup ekranı sadece liste + trigger seviyesinde
- backup type, retention, component scope, restore intent yüzeyi yok

### 5. Maintenance durability eksik

Backend:
- maintenance state in-memory tutuluyor
- Redis flag yazılıyor ama kalıcı repository/state modeli yok
- restart sonrası state continuity garanti değil
- bakım penceresi, scheduled activation, expiry, operator note yok

Frontend:
- çok dilli mesaj yüzeyi iyi durumda ama scheduling / expiry / preview / public impact alanları yok

### 6. License capability stub

Backend:
- license info ve activation hâlâ stub
- gerçek lisans sağlayıcısı / local signed license parse / entitlement validation yok

Frontend:
- lisans sayfası iyi paketlenmiş ama gerçek backend capability ile beslenmiyor

### 7. Test coverage modül seviyesinde eksik

Backend:
- `internal/modules/system_ops` altında `_test.go` yok
- worker recovery, audit filtering, update/rollback persistence, maintenance state, DTO normalization gibi alanlar testlenmiyor

Frontend:
- admin yüzeyinin davranışına dair page/component testleri yok

## Closure Definition

System Ops modülü ancak aşağıdaki noktalar tamamlanınca “feature-complete / closable” sayılmalıdır:

- update, rollback, backup, migration, maintenance ve license yüzeyleri stub olmaktan çıkmalı
- Go ve Next sözleşmeleri drift üretmeden aynı capability setini taşımalı
- operasyonel yan etkiler yalnız audit log’a değil, gerçek persistence ve lifecycle state’e bağlanmalı
- en az çekirdek backend test coverage açılmalı

## Execution Slices

### Slice 1: Update and Rollback Realization

Backend:
- update check için gerçek manifest/provider adapter aç
- `tarball_url`, `confirm_token`, update options ve preflight validation uygula
- rollback listesinde repository-level filtering ve doğru total hesabı yap
- empty target için “previous known version” çözümleyicisi ekle

Frontend:
- update start ekranına opsiyonel preflight seçenekleri ekle
- rollback start ekranına recommended target / previous version akışı ekle

Acceptance:
- update check stub olmaktan çıkar
- rollback pagination doğru çalışır
- request/response contract drift kapanır

### Slice 2: Migration and Backup Operationalization

Backend:
- migration run için gerçek execution pipeline ekle
- backup trigger için type, note, scope ve completion lifecycle ekle
- backup record’ları running’den terminal state’e taşıyan gerçek iş akışı kur

Frontend:
- migration ekranında run result / progress / terminal feedback yüzeyi aç
- backup ekranında type/note/scope alanları ve richer record detail göster

Acceptance:
- migration trigger artık yalnız audit log yazmaz
- backup record’ları kalıcı ve tamamlanabilir lifecycle taşır

### Slice 3: Maintenance Persistence and Scheduling

Backend:
- maintenance state için kalıcı repository/model aç
- restart sonrası recovery ekle
- optional start/end window ve expiry desteği ekle

Frontend:
- scheduled maintenance, expiry ve active window görünürlüğü ekle
- public impact preview veya operator summary alanı aç

Acceptance:
- maintenance state process restart ile kaybolmaz
- schedule/expiry destekli gerçek operational mode olur

### Slice 4: License Reality Layer

Backend:
- gerçek lisans provider veya local signed entitlement parse akışı ekle
- activate endpoint’i stub olmaktan çıkar
- entitlement/result/error modelini netleştir

Frontend:
- license ekranında real validation state, entitlement source ve failure reason göster

Acceptance:
- license endpoints static fake response dönmez

### Slice 5: Test and Verification Closure

Backend:
- service tests
- handler tests
- worker recovery tests
- repository filtering tests

Frontend:
- en az smoke-level component/page tests
- API contract normalization tests

Acceptance:
- System Ops için minimum regression safety açılmış olur

## Recommended Delivery Order

1. Update + rollback contract and backend reality
2. Backup + migration execution
3. Maintenance persistence
4. License reality layer
5. Test closure

## Out of Scope

Bu planın dışında tutulur:
- genel platform governance/export hattı
- tenant bazlı gelişmiş permission matrix
- public operator mobile tooling
- release engineering dışı deployment stack refactor

## Final Recommendation

System Ops bugün hardening açısından stabilize edilmiş durumda, ancak feature-level kapanış için hâlâ açık backlog var.

Bu yüzden doğru yönetim kararı:
- hardening kapandı
- module closure henüz tamamlanmadı
- kapanış için bu dokümandaki execution slice’lar izlenmeli
