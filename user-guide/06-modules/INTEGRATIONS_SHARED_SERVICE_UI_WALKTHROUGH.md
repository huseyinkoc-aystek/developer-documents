---
Title: Integrations Shared Service UI Walkthrough
Version: 1.0.0
Status: Active
Last Updated: 2026-03-08
Related Documents:
  - ./INTEGRATIONS_PLATFORM_GUIDE.md
  - ./modules/integrations.md
  - ../05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md
---

# Integrations Shared Service UI Walkthrough

Bu doküman, Shared Service ve Service Grant ekranlarının kullanıcı tarafında adım adım nasıl kullanılacağını anlatır.

## 1. Shared Service Oluşturma
Ekran:
1. `Admin > Integrations > Shared Services > Oluştur`

![Shared Service Create Form](/workspace/aystek-core-dev-docs/user-guide/assets/integrations/shared-service-walkthrough/01-shared-service-create-form.svg)

Adımlar:
1. `Key`, `Name`, `Protocol`, `Status` alanlarını doldurun.
2. `Connector` seçin.
3. `Allowed Actions` listesinden izinli aksiyonları işaretleyin.
4. `Default Scopes`, `Allowed Methods`, `Allowed Path Rules` alanlarını girin.
5. `Default Limits` içinde `rpm` ve `daily_quota` tanımlayın.
6. Gerekliyse `Secret Bindings`, `Secret Values`, `Payload Guard`, `mTLS Policy` alanlarını doldurun.
7. Kaydedin.

Kontrol:
1. Form içindeki `Policy Hardening Checks` panelinde kritik uyarı kalmamalı.

![Policy Hardening Panel](/workspace/aystek-core-dev-docs/user-guide/assets/integrations/shared-service-walkthrough/02-policy-hardening-panel.svg)

## 2. Service Grant Oluşturma
Ekran:
1. `Admin > Integrations > Service Grants > Oluştur`

![Service Grant Create Form](/workspace/aystek-core-dev-docs/user-guide/assets/integrations/shared-service-walkthrough/03-service-grant-create-form.svg)

Adımlar:
1. `Shared Service` seçin.
2. `Grant Key`, `Name`, `Status` alanlarını doldurun.
3. `Allowed Actions` listesini servis kapsamına göre işaretleyin.
4. `Scopes` girin (`sms.send`, `sms.*` gibi).
5. Gerekirse grant seviyesinde `allowed_methods`, `allowed_path_rules`, `limits` override edin.
6. Gerekirse `secret_binding_overrides`, `payload_guard`, `mtls_policy` girin.
7. Kaydedin.

## 3. Hook ile Grant Bağlama
Ekran:
1. `Admin > Integrations > Hooks > [Hook Edit]`

Adımlar:
1. Hook config içinde `shared_service_grant_id` alanına grant id yazın.
2. Hook’u `async` veya `sync` modda ihtiyaca göre kaydedin.
3. Test çalıştırın.

Beklenen:
1. Execution sonucunda shared bağlam alanları görünür:
   - `shared_service_id`
   - `shared_service_grant_id`

## 4. Usage & Limit Observability Kullanımı
Ekran:
1. `Admin > Integrations > Shared Services > [Edit]`
2. `Admin > Integrations > Service Grants > [Edit]`

![Service Usage Observability](/workspace/aystek-core-dev-docs/user-guide/assets/integrations/shared-service-walkthrough/04-service-usage-observability.svg)
![Grant Usage Observability](/workspace/aystek-core-dev-docs/user-guide/assets/integrations/shared-service-walkthrough/05-grant-usage-observability.svg)

Adımlar:
1. `7d / 14d / 30d` pencere seçin.
2. Toplam istek, bugünkü istek, denied ve failure kartlarını kontrol edin.
3. Warning panelini inceleyin.

Yorumlama:
1. `Denied` yükseliyorsa policy drift veya kötü payload olabilir.
2. `Daily Usage %` yüksekse kota baskısı vardır.
3. `secret/mTLS` kaynaklı warning varsa güvenlik config’i düzeltilmelidir.

## 5. Hızlı Sorun Giderme
1. `scope_denied`: grant scopes ve action capability/action_key uyumunu kontrol edin.
2. `method_denied`: action method ile allowlist uyuşmasını kontrol edin.
3. `path_denied`: action path ile path rules uyuşmasını kontrol edin.
4. `secret_binding_denied`: binding map ve secret values tutarlılığını kontrol edin.
5. `mtls_denied`: cert/key path ve mTLS required alanını kontrol edin.

## 6. Ekran Görseli Checklist (Dokümana Eklenecek)
Bu başlık, dokümana gerçek ekran görselleri ekleneceği zaman referans checklist olarak kullanılır.

1. Shared service create form (policy alanları açık).
2. Hardening warning panel örneği.
3. Service grant create form (scope/limit override).
4. Service usage summary (7/14/30 selector + cards).
5. Grant usage summary (daily usage ratio + trend table).

Not:
1. Görsellerin standart klasörü: `user-guide/assets/integrations/shared-service-walkthrough/`
2. Şu an `.svg` placeholder dosyaları bağlıdır.
3. Gerçek capture alındığında aynı base adlarla `.png` dosya koyup linkleri `.png` olarak güncelleyin.
