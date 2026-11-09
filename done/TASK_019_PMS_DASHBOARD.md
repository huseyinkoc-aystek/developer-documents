# TASK_019 — PMS Dashboard

## Kapsam
- Admin KPI endpointi eklendi:
  - `GET /admin/pms/dashboard?period=YYYY-MM`
- KPI seti:
  - toplam acik borc
  - tahsilat orani
  - acik is emri sayisi
  - mutabakat kuyrugu sayisi
  - top 5 gecikmeli birim

## Backend
- Yeni paket: `internal/modules/pms/dashboard`
  - repository/service/handler/routes
- PMS module wiring'e dahil edildi.
- Permission guard: `pms_dashboard.read`.

## Frontend
- Yeni sayfa: `/admin/pms/dashboard`
- Posts/PMS tasarim cizgisine uygun:
  - modern header
  - period filtresi
  - KPI kartlari
  - top overdue units tablosu

## Seed/I18n
- Menu:
  - `menu.operations.pms.dashboard`
  - `pms-dashboard`, `pms-suite-dashboard`
- Permission:
  - `pms_dashboard:*` temel set
- Translation:
  - `pms.dashboard.*` anahtarlari
