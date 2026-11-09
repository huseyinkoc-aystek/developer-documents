# Telematics Wave A Device And Protocol Productization Execution Plan

Date: 2026-03-25
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / aystek-core-dev-docs
Source: `TELEMATICS_SCOPE_COMPLETENESS_AND_NEXT_WAVE_PLAN_2026-03-25.md`

## Objective

Tamamlanmış çekirdek telematics platformunu gerçek cihaz ekosistemiyle ürün seviyesinde buluşturmak.

Wave A hedefi:
- vendor protocol adapter library
- protocol certification ve test disiplini
- device lifecycle admin surface
- object/asset tracking için düşük güç ve beacon odaklı model

## Scope

### In Scope

1. Vendor protocol adapter registry
2. Teltonika Codec8 ailesi için ilk gerçek adapter
3. Queclink sınıfı ikinci adapter için skeleton + certification contract
4. Device capability matrix
5. Protocol raw payload -> canonical telemetry mapping catalog
6. Device lifecycle admin center
7. Object tracking state modeli
8. Certification test pack ve replay fixtures

### Out of Scope

1. Mobil uygulamalar
2. Video telematics
3. IoT sensor full product stack
4. UTTS / U-ETDS / takograf modülleri
5. Ticari paketleme ve GTM araçları

## Guardrails

1. Mevcut generic telemetry ingest kırılmamalı.
2. Vendor adapter akışı tenant/site isolation standardını bozmamalı.
3. Canonical telemetry model tek kalmalı; vendor-specific alanlar adapter metadata altında tutulmalı.
4. Existing replay, alert, CAN, command lifecycle adapter-agnostic kalmalı.
5. Certification başarısız adapter production-default olmamalı.

## Target Vendors And Tracks

### Track A1

- Teltonika Codec8 / Codec8 Extended
- hedef: ilk production-grade adapter

### Track A2

- Queclink benzeri @Track sınıfı adapter
- hedef: ingestion skeleton + canonical mapping contract + certification harness

### Track A3

- low-power object trackers / BLE beacon gateway akışı
- hedef: object tracking state modeli

## Delivery Phases

Not:
- Bu planın ilk yürütme dalgası `Phase-1: Adapter Foundation And Certification Operations` olarak kapatıldı.
- Sonraki aktif dalga `Phase-2: Device Lifecycle Center` olarak ayrı execution plan ile yürütülür.

### Phase-1: Adapter Foundation And Certification Operations

1. `protocol_adapters` registry domain modeli
2. adapter capability schema:
   - gps
   - ignition
   - can
   - obd
   - beacon
   - battery
   - digital_io
3. canonical field mapping contract
4. vendor raw payload fixture repository
5. certification result modeli

### Phase-2: Device Lifecycle Center

1. device lifecycle read model
2. provisioning state machine
3. certification gate visibility
4. firmware/SIM/connectivity health surface
5. device timeline ve remediation notes
6. fleet-level lifecycle KPI
7. stale / uncertified / frozen segmentation

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md)

### Phase-3: Queclink Skeleton And Contract

1. parser skeleton
2. message family registry
3. canonical mapping draft
4. certification placeholders
5. unsupported field behavior policy

Active next execution:
- [TELEMATICS_WAVE_A_PHASE_3_TELTONIKA_ONBOARDING_AND_PROVISIONING_AUTOMATION_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_3_TELTONIKA_ONBOARDING_AND_PROVISIONING_AUTOMATION_EXECUTION_PLAN_2026-03-26.md)

### Phase-4: Teltonika Onboarding Hardening And Assignment Automation

1. IMEI to stock/device assignment automation
2. assignment mismatch and confidence model
3. SIM/operator readiness and first-signal timeout hardening
4. activation-ready transition automation
5. vendor-specific onboarding playbook expansion

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md)

Completion report:
- [TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_COMPLETION_REPORT_2026-03-27.md)

### Phase-5: Operations And Hardening

1. first-signal timeout governance
2. activation-ready automation
3. provisioning blocked escalation
4. vendor-specific provisioning hardening
5. provisioning runbook and closure trend

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_FIRST_SIGNAL_AND_PROVISIONING_HARDENING_EXECUTION_PLAN_2026-03-27.md)

Completion report:
- [TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_5_COMPLETION_REPORT_2026-03-27.md)

### Phase-6: Vendor Activation Automation And Field Exception Handling

1. vendor-specific activation automation
2. field exception classification and remediation
3. activation hold / retry / manual handoff lifecycle
4. exception governance and closure trend

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_VENDOR_ACTIVATION_AUTOMATION_AND_FIELD_EXCEPTION_HANDLING_EXECUTION_PLAN_2026-03-27.md)

Completion report:
- [TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_6_COMPLETION_REPORT_2026-03-28.md)

### Phase-7: Activation Orchestration Hardening And Field Exception Recurrence Escalation

1. activation orchestration policy matrix
2. retry budget and cooldown ladder
3. field exception recurrence escalation
4. freeze recommendation and controlled recovery visibility

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_7_ACTIVATION_ORCHESTRATION_HARDENING_AND_FIELD_EXCEPTION_RECURRENCE_ESCALATION_EXECUTION_PLAN_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_7_ACTIVATION_ORCHESTRATION_HARDENING_AND_FIELD_EXCEPTION_RECURRENCE_ESCALATION_EXECUTION_PLAN_2026-03-28.md)

Completion report:
- [TELEMATICS_WAVE_A_PHASE_7_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_7_COMPLETION_REPORT_2026-03-28.md)

### Phase-8: Freeze Recommendation Operationalization And Controlled Recovery Validation

1. freeze recommendation decision ladder
2. controlled recovery validation checklist
3. re-freeze escalation and rollback confidence
4. governance export and closure trend

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_8_FREEZE_RECOMMENDATION_OPERATIONALIZATION_AND_CONTROLLED_RECOVERY_VALIDATION_EXECUTION_PLAN_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_8_FREEZE_RECOMMENDATION_OPERATIONALIZATION_AND_CONTROLLED_RECOVERY_VALIDATION_EXECUTION_PLAN_2026-03-28.md)

Completion report:
- [TELEMATICS_WAVE_A_PHASE_8_COMPLETION_REPORT_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_8_COMPLETION_REPORT_2026-03-28.md)

### Phase-9: Freeze Closure And Recovery Effectiveness Scoring

1. freeze acted / not-acted closure quality
2. recovery effectiveness scoring
3. regression after action visibility
4. governance export and scheduled snapshot

Detailed plan:
- [TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_9_FREEZE_CLOSURE_AND_RECOVERY_EFFECTIVENESS_SCORING_EXECUTION_PLAN_2026-03-28.md)

## Concrete Backlog

## Backend

1. Create `telematics/protocol_adapters` module
2. Add adapter registry collection + repository + indexes
3. Add certification result collection or embedded audit model
4. Implement Teltonika Codec8 parser package
5. Add fixture-driven parser tests
6. Add canonical mapping service
7. Extend ingest pipeline with adapter resolution strategy
8. Add unsupported payload DLQ / audit marker
9. Add object tracking state aggregate
10. Add gateway-beacon relation service

## Frontend

1. Add `admin/telematics/protocol-adapters` page
2. Add adapter registry list/detail
3. Add certification status badges
4. Add dry-run payload tester
5. Add device lifecycle center
6. Add object tracking dashboard

## Docs / Ops

1. Create vendor certification checklist
2. Create Teltonika fixture curation guide
3. Create adapter onboarding runbook
4. Create unsupported payload triage guide

## Acceptance Criteria

1. En az 1 gerçek vendor adapter production-grade çalışmalı.
2. En az 1 ikinci vendor için certification-ready skeleton bulunmalı.
3. Device lifecycle ekranı adapter/capability/certification durumunu göstermeli.
4. Unsupported payload’lar audit ve triage akışına düşmeli.
5. Object tracking için low-power state modeli çalışmalı.
6. Replay/alert/command akışları vendor adapter katmanından bağımsız kalmalı.

## Immediate Execution Checklist

1. `protocol_adapters` bounded context aç
2. registry + capability schema yaz
3. Teltonika Codec8 fixture setini oluştur
4. parser test harness kur
5. ingest pipeline’a adapter resolution ekle
6. admin registry page ve dry-run tester aç
7. device lifecycle center için ilk read-modeli kur

## Risks

1. Vendor protocol nüansları canonical modele zor sığabilir.
2. Raw payload fixture eksikliği certification kalitesini düşürür.
3. Low-power object tracking ile vehicle tracking aynı cadence varsayımlarıyla yönetilemez.
4. Adapter-specific alanlar modele kontrolsüz sızabilir.

## Risk Mitigation

1. Mapping contract’ı versioned tut.
2. Fixture setini gerçek saha payload’larıyla büyüt.
3. Vehicle ve object tracking state machine’lerini ayrı tasarla.
4. Unsupported field ve unsupported message audit zinciri kur.

## Success Metrics

1. Certified adapter count
2. Unsupported payload rate
3. Parser failure rate
4. Canonical mapping coverage
5. Device certification freshness
6. Object tracking stale-rate accuracy

## Progress Status

### Completed

1. Çekirdek telematics ingest, replay, alert, command ve governance altyapısı hazır.
2. CAN decoder registry ve dry-run yaklaşımı hazır.
3. External partner/callback governance ve audit/export zinciri hazır.
4. Phase-1 adapter foundation, Teltonika ingress, certification lifecycle ve governance tamamlandı.

### Remaining

1. Teltonika onboarding hardening ve assignment automation
2. Queclink skeleton ve ikinci vendor contract
3. Object tracking özel lifecycle
4. Multi-vendor provisioning parity

## References

- [TELEMATICS_SCOPE_COMPLETENESS_AND_NEXT_WAVE_PLAN_2026-03-25.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_SCOPE_COMPLETENESS_AND_NEXT_WAVE_PLAN_2026-03-25.md)
- [PHASE_8_COMPLETION_REPORT.md](/workspace/aystek-core-dev-docs/plans/PHASE_8_COMPLETION_REPORT.md)
- [TELEMATICS_WAVE_A_PHASE_1_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_1_COMPLETION_REPORT_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_2_DEVICE_LIFECYCLE_CENTER_EXECUTION_PLAN_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_3_COMPLETION_REPORT_2026-03-26.md)
- [TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_A_PHASE_4_TELTONIKA_ONBOARDING_HARDENING_AND_ASSIGNMENT_AUTOMATION_EXECUTION_PLAN_2026-03-26.md)
