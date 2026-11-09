# Device Vendor Command Orchestration Analysis

Date: 2026-04-17
Owner: Codex
Status: Analysis Ready

## Problem

Stock cihazlarına bağlı SIM hatlarına operasyonel komut gönderilmesi gerekiyor.

Bu ihtiyaç:
- serbest telefon numarasına mesaj atma
- basit şablon yönetimi

seviyesinden daha ileri bir ihtiyaçtır.

Gerçek ihtiyaç:
- doğru SIM hedefini bulmak
- cihaz markası / entegrasyon profiline göre doğru komut formatını seçmek
- cihaz durumuna göre doğru komutu üretmek
- sonucu delivery + callback + audit akışıyla izlemek

## Core Clarification

Bu alan `Trio adapter` olarak konumlandırılmamalıdır.

Doğru model:
- çoklu cihaz markası
- çoklu komut ailesi
- çoklu delivery gateway

Örnek vendor / device owners:
- Trio
- Arvento
- N2 Mobil
- Mobiliz

Bu yüzden gerekli şey tek bir partner adapter değil, bir `device vendor command registry` ve bunun üstünde çalışan `command adapter` katmanıdır.

## Domain Anchors Already Available

Mevcut sistemde komut orkestrasyonu için yeterli veri kökleri var:

- stock inventory:
  - `product_type`
  - `msisdn`
  - `imei`
  - `operator`
  - `brand_id`
  - `model_id`
- stock catalog:
  - `integration_key`
- stock vehicle assets:
  - araca bağlı aktif HAT/SIM ilişkisi
- telematics lifecycle / provisioning:
  - SIM readiness
  - inventory binding
  - onboarding state

## Required Enterprise Model

### 1. Target Resolution

Bir komut hedefi şu mantıkla çözülmeli:

1. seçilen araç / cihaz / identifier
2. aktif bağlı HAT/SIM asset
3. authoritative `msisdn`
4. device vendor / adapter identity

`imei` ve araç plakası lookup içindir.
Final delivery target her zaman çözümlenmiş `msisdn` olmalıdır.

### 2. Command Intent Layer

Sistem ham “message body” ile çalışmamalı.
Öncelikli kavram `command_intent` olmalıdır.

Örnek intent’ler:
- `query_status`
- `enable_tracking`
- `disable_tracking`
- `reset_device`
- `switch_profile`
- `sync_apn`
- `wake_up`

### 3. Device Vendor Command Registry

Sistem şu yapıyla çalışmalıdır:

- `command_intent`
- `vendor_key`
- `product_type`
- `template_code`
- `variable_schema`
- `rendering_rules`
- `delivery_strategy`

Örnek:
- `reset_device` + `vendor_key=trio`
- `reset_device` + `vendor_key=mobiliz`
- `reset_device` + `vendor_key=arvento`

Hepsi aynı intent’i temsil eder ama farklı komut string’i, farklı variable set’i ve farklı delivery kuralı kullanabilir.

### 4. Delivery Adapter Layer

Device vendor ile SMS gateway aynı şey değildir.

Bu yüzden iki ayrı katman gerekir:

- device command adapter
  - intent -> vendor command mapping
  - variable rendering
  - state guard
- SMS delivery adapter
  - provider selection
  - callback normalization
  - provider message correlation

## Proposed Registry Shape

### Device Command Definition

Suggested fields:
- `vendor_key`
- `command_intent`
- `title`
- `description`
- `message_template`
- `required_variables`
- `allowed_states`
- `cooldown_sec`
- `active`

### Vendor Resolution

Vendor key çözümü öncelik sırası:

1. explicit device/vendor profile
2. catalog brand `integration_key`
3. catalog model `integration_key`
4. fallback manual override

## Operational Guardrails

Aşağıdakiler olmadan bu alan enterprise-ready sayılmamalı:

- template approval
- tenant/site scope enforcement
- vendor allowlist
- cooldown policy
- duplicate suppression
- delivery callback reconciliation
- bridge error visibility
- audit trail with actor + target + intent + rendered body snapshot

## UI Implication

UI serbest telefon kutusu merkezli olmamalı.

Doğru akış:
- target lookup
- resolved SIM preview
- vendor resolution preview
- command intent selection
- rendered command preview
- guarded send
- execution stream

## Recommended Next Implementation Slice

1. `device_vendor_key` çözümleme kuralını aç
2. `command_intent` alanı tanımla
3. `device vendor command registry` read modelini ekle
4. `device command adapter` contract aç
5. ilk vendor adapter setlerini sırayla ekle:
   - Trio
   - Mobiliz
   - Arvento
   - N2 Mobil
6. SMS gateway layer ile callback/audit zincirini bağla

## Conclusion

Bugünkü taşıma işi doğru şekilde tamamlandı:
- UI artık integrations altında konumlandı
- stock backend domain sınırı korunuyor

Sonraki gerçek enterprise adımı:
- `device vendor command registry` açmak
- marka bazlı komut tanımlarını standartlaştırmak
- serbest SMS gönderimini operasyonel intent + vendor adapter modeline dönüştürmek
