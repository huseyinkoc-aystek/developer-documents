# Telematics Wave C Sprint Backlog

Date: 2026-03-29
Owner: AI Assistant (Codex)
Status: Active
Depends on:
- [TELEMATICS_WAVE_B_SPRINT_BACKLOG_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_B_SPRINT_BACKLOG_2026-03-29.md)
- [TELEMATICS_WAVE_C_IOT_VIDEO_AND_SENSOR_STACK_EXECUTION_PLAN_2026-03-29.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_WAVE_C_IOT_VIDEO_AND_SENSOR_STACK_EXECUTION_PLAN_2026-03-29.md)

## Sprint Cadence

- Sprint length: 2 weeks
- Definition of done: code + tests + lint + docs + audit trace
- Environments: local dev, controlled pilot

## Sprint C1 (Sensor Core)

Target:
1. sensor registry for temperature, humidity, door, damper, panic, BLE
2. sensor threshold/rule policy templates
3. alert notice linkage for sensor exceedance

Primary deliverables:
1. sensor inventory model and repository contracts (C1-S1-01) ✅ completed 2026-03-29
2. sensor threshold policy model and rule templates (C1-S1-02) ✅ completed 2026-03-29
3. admin sensor CRUD and reading ingest endpoint (C1-S1-03) ✅ completed 2026-03-29

## Sprint C2 (Cold-Chain Operations)

Target:
1. route-level temperature compliance tracking
2. excursion timeline with acknowledgement workflow
3. compliance export and audit

Primary deliverables:
1. cold-chain excursion model and repository (C2-S1-01) ✅ completed 2026-03-29
2. excursion acknowledgement and status lifecycle (C2-S1-02) ✅ completed 2026-03-29
3. compliance summary and export endpoint (C2-S1-03) ✅ completed 2026-03-29

## Sprint C3 (Video Operations)

Target:
1. video event ingest contract
2. incident timeline and clip references
3. retention policy and access-control

Primary deliverables:
1. video event model and ingest endpoint (C3-S1-01) ✅ completed 2026-03-29
2. video incident timeline and clip reference handler (C3-S1-02) ✅ completed 2026-03-29
3. retention policy enforcement model (C3-S1-03) ✅ completed 2026-03-29

## Wave C Checkpoint (2026-03-29)

All C1/C2/C3 deliverables implemented and validated:

### New Packages
- `telematics/sensors/model` — Sensor, SensorReading, ThresholdPolicy models (temp/humidity/door/damper/panic/BLE)
- `telematics/sensors/repository` — interface + memory implementation (CRUD + readings + policies)
- `telematics/cold_chain/model` — Excursion, ComplianceSummary models
- `telematics/cold_chain/repository` — interface + memory implementation
- `telematics/video_events/model` — VideoEvent, ClipReference, RetentionPolicy models (ADAS/DMS/panic/overspeed/etc.)
- `telematics/video_events/repository` — interface + memory implementation

### New Handler Functions (in `handler/`)
- `sensors.go`: RegisterSensor, ListSensors, GetSensor, DecommissionSensor, IngestSensorReading, ListSensorReadings, UpsertSensorPolicy, ListSensorPolicies
- `cold_chain.go`: RecordExcursion, ListExcursions, GetExcursion, AcknowledgeExcursion, GetColdChainCompliance
- `video_events.go`: IngestVideoEvent, ListVideoEvents, GetVideoEvent, UpsertVideoRetentionPolicy, GetVideoRetentionPolicy

### New Routes (registered in `routes.go`)
```
POST   /iot/sensors
GET    /iot/sensors
GET    /iot/sensors/:id
DELETE /iot/sensors/:id
POST   /iot/sensors/:id/readings
GET    /iot/sensors/:id/readings
POST   /iot/sensor-policies
GET    /iot/sensor-policies
POST   /iot/cold-chain/excursions
GET    /iot/cold-chain/excursions
GET    /iot/cold-chain/excursions/:id
POST   /iot/cold-chain/excursions/:id/ack
GET    /iot/cold-chain/compliance
POST   /iot/video-events
GET    /iot/video-events
GET    /iot/video-events/:id
POST   /iot/video-retention
GET    /iot/video-retention
```

### Tests
- `handler/wave_c_test.go`: 18 tests, all PASS
- Full telematics suite: 11 packages OK, zero failures
