# SYSTEM OPS Module Hardening and I18N Execution Plan

Date: 2026-04-03
Status: Stabilized
Owner: Codex

## Objective

System Ops modülünü Copilot üretimi temel halinden çıkarıp AYSTEK standartlarına uygun, daha sürdürülebilir ve çoklu dil dostu bir operasyon yüzeyine dönüştürmek.

## Current-State Review

Backend tarafında `system_ops` modülü çalışır durumda:
- update
- rollback
- migration
- backup
- maintenance
- node
- audit
- license
- version

Next admin tarafında geniş bir yüzey var, ancak üç kritik boşluk tespit edildi:

1. API contract drift
- Go `json` alanları ile Next tipleri ve sayfa kullanımları tam hizalı değil.
- Aynı modül için iki ayrı tip sözleşmesi oluşmuş.
- `record_id` ve SSE log stream akışı sessiz şekilde yanlış bağlanabiliyor.

2. Çoklu dil uyumsuzluğu
- Maintenance mesajı veri modelinde çok dilli `map[string]string` taşımasına rağmen UI sadece `tr/en` ile sınırlıydı.
- Bu yüzden modül platformun aktif dil modeliyle uyumlu değildi.

3. Ürünleşme eksikliği
- Sayfalar tek tek vardı ama ortak System Ops operasyon rail’i yoktu.
- UX bütünlüğü düşük, discoverability zayıftı.

## Execution Slices

### Slice 1: Contract Hardening
- `system-ops.ts` API client ile `lib/types/system-ops.ts` sözleşmesini hizala
- snake_case ve camelCase geçişlerinde normalize edici katman ekle
- update / rollback start cevaplarına `record_id` ekle
- update detail log stream’i `job_id` veya `record_id` ile çözülebilir hale getir

### Slice 2: I18N Remediation
- maintenance ekranını tüm aktif dillere uyumlu hale getir
- System Ops rail ve yeni maintenance alanları için translation seed ekle
- fallback ile çalışan ama key-first ilerleyen yüzey standardını uygula

### Slice 3: UX Productization
- shared `SystemOpsRail` ekle
- ana System Ops sayfalarına ortak rail bağla
- modülü operasyon merkezi hissi verecek şekilde toparla

### Slice 4: Verification
- `go test ./internal/modules/system_ops/... ./internal/seeds ./cmd/api`
- System Ops page lint seti

## Done in This Pass

- API/type normalization katmanı eklendi
- update/rollback start response `record_id` ile güçlendirildi
- SSE log stream çözümü `job_id` ve `record_id` için dayanıklı hale getirildi
- maintenance UI aktif dillere göre dinamik mesaj editörüne dönüştürüldü
- shared System Ops rail eklendi ve ana sayfalara bağlandı
- translation seed yeni System Ops rail ve multilingual maintenance anahtarlarıyla genişletildi

## Remaining Improvement Backlog

- rollbacks, migrations ve alt yüzeylerin tamamına aynı rail/shell standardını taşımak
- health, backup ve audit tablolarında locale-aware tarih/sayı format helper’ı ortaklaştırmak
- license / node sayfalarındaki local inline type’ları kaldırıp tam shared type kullanmak
- System Ops için activity-log / stock reporting benzeri governance/export hattı gerekiyorsa ayrı plan açmak

## Closure Note

Bu plan kapsamındaki hardening backlog’u 2026-04-04 itibarıyla stabilize edildi.

Takip dokümanı:
- `SYSTEM_OPS_MODULE_STABILIZED_REPORT_2026-04-04.md`

Bu tarihten sonraki işler varsayılan olarak hardening değil, feature-level genişleme olarak ele alınmalıdır.

## Acceptance Criteria

- System Ops maintenance flow tüm aktif dillere mesaj girebilmeli
- Update start sonrası detail sayfası doğru record’a gitmeli
- Update detail log akışı record id ile de çalışmalı
- System Ops ana sayfaları ortak operasyon rail’i taşımalı
- translation seed ve code surface birlikte compile/lint/test geçmeli
