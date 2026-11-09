# WP-1 Kaynak Metinler — `languages.*` (62 anahtar)

> Bu dosya `I18N_TRANSLATION_DEBT_WORKPLAN_TR.md` içindeki **WP-1** için çeviri kaynağıdır.
> Ölçüm tarihi: 2026-09-21. Kaynak: `aystek-core/internal/seeds/translation_seed.go`
> (`languages.*` anahtarları, 19 dilde Türkçe metin tekrar ediyor).
> Güncel liste her zaman şu komutla alınır (bu tablo referanstır, komut çıktısı geçerlidir):
> `cd /workspace/go/aystek-core && TRANSLATION_REPORT_PREFIX=languages. go test ./internal/seeds -run TestTranslationCoverageReport -v`

**Yer tutucular** (aynen korunmalı, çevrilmemeli): `{{name}}`, `{{code}}`, `{{count}}`.
Tabloda yalnızca `languages.delete.message` `{{name}}` içerir.

| # | Anahtar | tr (kaynak) | en | de |
|---|---|---|---|---|
| 1 | `languages.actions.delete` | Sil | Delete | Löschen |
| 2 | `languages.actions.edit` | Düzenle | Edit | Bearbeiten |
| 3 | `languages.button.new` | Yeni Dil | New Language | Neue Sprache |
| 4 | `languages.default.badge` | Varsayılan | Default | Standard |
| 5 | `languages.default.makeDefault` | Varsayılan yap | Make default | Als Standard setzen |
| 6 | `languages.delete.confirm` | Dili Sil | Delete Language | Sprache löschen |
| 7 | `languages.delete.deleting` | Siliniyor... | Deleting... | Wird gelöscht... |
| 8 | `languages.delete.message` | “{{name}}” dilini silmek istediğinize emin misiniz? | Are you sure you want to delete the “{{name}}” language? | “{{name}}” dilini silmek istediğinize emin misiniz? |
| 9 | `languages.delete.title` | Dili Sil | Delete Language | Sprache löschen |
| 10 | `languages.empty` | Dil bulunamadı. | No language found. | Keine Sprache gefunden. |
| 11 | `languages.form.addEntry` | Yeni giriş | Add entry | Eintrag hinzufügen |
| 12 | `languages.form.cancel` | Vazgeç | Cancel | Abbrechen |
| 13 | `languages.form.code` | Dil Kodu | Language Code | Sprachcode |
| 14 | `languages.form.codePlaceholder` | ör. tr | e.g. tr | z.B. de |
| 15 | `languages.form.enabled` | Aktif | Active | Aktiv |
| 16 | `languages.form.entryCodePlaceholder` | Dil kodu | Language code | Sprachcode |
| 17 | `languages.form.entryLabelPlaceholder` | Dil adı | Language name | Sprachname |
| 18 | `languages.form.localizedTitle` | Localized İsimler | Localized Names | Lokalisierte Namen |
| 19 | `languages.form.makeDefault` | Varsayılan yap | Set as default | Als Standard setzen |
| 20 | `languages.form.removeEntry` | Sil | Remove | Entfernen |
| 21 | `languages.form.save` | Kaydet | Save | Speichern |
| 22 | `languages.loading` | Diller yükleniyor... | Languages are loading... | Sprachen werden geladen... |
| 23 | `languages.messages.createSuccess` | Dil oluşturuldu. | Language created. | Sprache erstellt. |
| 24 | `languages.messages.defaultChangeSuccess` | Varsayılan dil güncellendi. | Default language updated. | Standardsprache aktualisiert. |
| 25 | `languages.messages.deleteSuccess` | Dil silindi. | Language deleted. | Sprache gelöscht. |
| 26 | `languages.messages.updateSuccess` | Dil güncellendi. | Language updated. | Sprache aktualisiert. |
| 27 | `languages.modal.close` | Kapat | Close | Schließen |
| 28 | `languages.modal.createTitle` | Yeni Dil | New Language | Neue Sprache |
| 29 | `languages.modal.editTitle` | Dili Düzenle | Edit Language | Sprache bearbeiten |
| 30 | `languages.modal.localizedNames` | Localized İsimler | Localized Names | Lokalisierte Namen |
| 31 | `languages.modal.subtitle` | Dil kodu benzersiz olmalı ve localized isimler sanitize edilir. | Language code must be unique; localized names are sanitized. | Sprachcode muss eindeutig sein; lokalisierte Namen werden bereinigt. |
| 32 | `languages.modal.titleEdit` | Dili Düzenle | Edit Language | Sprache bearbeiten |
| 33 | `languages.modal.titleNew` | Yeni Dil | New Language | Neue Sprache |
| 34 | `languages.notice.codeRequired` | Dil kodu zorunludur. | Language code is required. | Sprachcode ist erforderlich. |
| 35 | `languages.notice.created` | Dil oluşturuldu. | Language created. | Sprache erstellt. |
| 36 | `languages.notice.defaultDeleteFail` | Varsayılan dil silinemez. | Default language cannot be deleted. | Standardsprache kann nicht gelöscht werden. |
| 37 | `languages.notice.defaultFail` | Varsayılan dil güncellenemedi. | Default language could not be updated. | Standardsprache konnte nicht aktualisiert werden. |
| 38 | `languages.notice.deleteFail` | Dil silinemedi. | Language could not be deleted. | Sprache konnte nicht gelöscht werden. |
| 39 | `languages.notice.deleted` | Dil silindi. | Language deleted. | Sprache gelöscht. |
| 40 | `languages.notice.loadError` | Diller yüklenemedi. | Languages could not be loaded. | Sprachen konnten nicht geladen werden. |
| 41 | `languages.notice.nameRequired` | En az bir localized isim girilmelidir. | At least one localized name is required. | Mindestens ein lokalisierter Name ist erforderlich. |
| 42 | `languages.notice.saveFail` | Dil kaydedilemedi. | Language could not be saved. | Sprache konnte nicht gespeichert werden. |
| 43 | `languages.notice.toggleFail` | Durum güncellenemedi. | Status could not be updated. | Status konnte nicht aktualisiert werden. |
| 44 | `languages.notice.updated` | Dil güncellendi. | Language updated. | Sprache aktualisiert. |
| 45 | `languages.statistics.active` | Aktif Dil | Active Languages | Aktive Sprachen |
| 46 | `languages.statistics.default` | Varsayılan Dil | Default Language | Standardsprache |
| 47 | `languages.statistics.total` | Toplam Dil | Total Languages | Gesamte Sprachen |
| 48 | `languages.statistics.translationRate` | Çeviri Oranı | Translation Rate | Übersetzungsrate |
| 49 | `languages.status.active` | Aktif | Active | Aktiv |
| 50 | `languages.status.inactive` | Pasif | Inactive | Inaktiv |
| 51 | `languages.subtitle` | Go mikro servisindeki language controller ile hizalı olarak dil kodlarını ve localized isimleri yönetin. | Manage language codes and localized names in sync with the Go language controller. | Verwalten Sie Sprachcodes und lokalisierte Namen synchron mit dem Go-Sprachcontroller. |
| 52 | `languages.table.actions` | İşlemler | Actions | Aktionen |
| 53 | `languages.table.code` | Kod | Code | Code |
| 54 | `languages.table.default` | Varsayılan | Default | Standard |
| 55 | `languages.table.headers.actions` | İşlemler | Actions | Aktionen |
| 56 | `languages.table.headers.code` | Kod | Code | Code |
| 57 | `languages.table.headers.default` | Varsayılan | Default | Standard |
| 58 | `languages.table.headers.localizedNames` | Localized İsimler | Localized Names | Lokalisierte Namen |
| 59 | `languages.table.headers.status` | Durum | Status | Status |
| 60 | `languages.table.localizedNames` | Localized İsimler | Localized Names | Lokalisierte Namen |
| 61 | `languages.table.status` | Durum | Status | Status |
| 62 | `languages.title` | Dil Yönetimi | Language Management | Sprachverwaltung |

