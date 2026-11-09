# AYSTEK Core Installer Apache Decision And Supported Nginx Migration Checkpoint

Date: 2026-04-28

## Scope

This slice closes the operational gap around Apache/httpd being present on the target server while the installer is configured to enable Nginx on ports `80/443`.

## Completed

- Added installer config support for Apache migration decisions:
  - `apache_detected`
  - `apache_service_name`
  - `apache_migration_mode`
- Added wizard-side Apache detection and operator decision flow:
  - migrate supported Apache config to Nginx and stop Apache
  - keep Apache and disable installer-managed Nginx/TLS
  - abort for manual review
- Added supported Apache migration engine:
  - detects common Apache config roots
  - backs up discovered Apache config files into installer backup storage
  - parses supported `VirtualHost` patterns
  - migrates supported `ServerName`, `ServerAlias`, `DocumentRoot`, `ProxyPass`, `SSLCertificateFile`, `SSLCertificateKeyFile`
  - writes migrated Nginx files into `/etc/nginx/conf.d`
  - generates `migration-report.txt`
  - stops/disables Apache only after successful Nginx config generation
- Added safe-stop behavior for unsupported directives:
  - `RewriteRule`
  - `RewriteCond`
  - `ProxyPassMatch`
  - `Alias`
  - `ScriptAlias`
  - `Include`
  - `IncludeOptional`
- Strengthened secret entry flow by ensuring tenant-site Mongo local password also uses hidden confirm-entry behavior.

## Files

- `/workspace/go/aystek-core-installer/internal/config/schema.go`
- `/workspace/go/aystek-core-installer/internal/ui/wizard.go`
- `/workspace/go/aystek-core-installer/internal/installer/installer.go`
- `/workspace/go/aystek-core-installer/internal/apache/migration.go`
- `/workspace/go/aystek-core-installer/internal/i18n/i18n.go`
- `/workspace/go/aystek-core-installer/internal/i18n/tr.go`
- `/workspace/go/aystek-core-installer/internal/i18n/en.go`
- `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`

## Validation

- `cd /workspace/go/aystek-core-installer && go test ./...`
- `cd /workspace/go/aystek-core-installer && go build ./cmd/aystek-core-installer`

## Honest limitation

This is not a claim of arbitrary Apache-to-Nginx parity. It is a supported migration path for the common AYSTEK-compatible Apache patterns above. Unsupported directives intentionally stop the migration and force operator review instead of silently producing a broken Nginx config.
