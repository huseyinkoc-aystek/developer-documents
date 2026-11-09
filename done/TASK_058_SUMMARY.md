# TASK 058 Summary - AYSTEK Haber

AYSTEK Haber is now a Core-backed public news site for `aystekhaber.com` under
tenant `aystek`. The implementation uses existing Core CMS ownership for posts,
categories, menus, placements, translations and view profiles, and adds only the
missing sports snapshot domain.

Completed items:

- Core store/site seed for `aystekhaber.com`
- 70 Turkish/English launch posts from the Figma template
- Site-scoped public navigation menu
- Post placements for hero, breaking ticker, sidebar and trending blocks
- Sports module for standings and latest matches
- `/admin/sports` management page
- Posts page event key documentation for placement events
- Public Next.js header/topbar/ticker alignment with the template
- Standalone runtime fix for CSS/JS/static assets and `sharp`
- Complete Figma section composition including video, local news, authors,
  mobile app, ads, newsletter and footer
- Core author API as the single source for writer cards and profiles
- Public author active/published policy hardening
- Sports inactive admin-read fix and stronger domain validation
- Robots, sitemap, RSS and localized author detail routes
- Installer profiles and documentation updates
- Public-site TLS primary/alias SAN handling and deferred HTTPS vhost upgrade

Key validation passed:

- Core sports and seed tests
- Core API build
- AYSTEK Haber production build
- Sports admin targeted lint
- Full admin panel production build
- Installer config/TLS/installer package tests
- Live API checks for posts, nav, placements and sports
- Browser checks for desktop/mobile render, no horizontal overflow and no failed
  static/image requests

Operational caveat: the dev container has no Docker CLI, so real installer
container deployment was not executed locally. The installer config/runbook and
local standalone runtime were validated.
