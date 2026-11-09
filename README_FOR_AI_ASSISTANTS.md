# 🤖 README FOR AI ASSISTANTS - MUST READ FIRST

**📌 CRITICAL: All AI Tools (ChatGPT, Claude, Copilot, etc.) MUST Read This Document Before Starting Any Work**

**Version:** 1.3.1  
**Last Updated:** September 29, 2026  
**Status:** IMMUTABLE REFERENCE - Mandatory Standards for ALL AI Tools

---

> **⚠️ ATTENTION AI ASSISTANTS:**  
> This document contains MANDATORY coding standards, architecture patterns, and operational procedures for the AYSTEK ecosystem.  
> **DO NOT skip this document.** Every rule here MUST be followed automatically without requiring reminders.  
> **Failure to follow these standards will result in rejected code and wasted development time.**

---

---

## 📋 Table of Contents

1. [Introduction & Purpose](#introduction--purpose)
2. [AI Role Definition](#ai-role-definition)
3. [Project Overview](#project-overview)
4. [Go Backend (aystek-core)](#go-backend-aystek-core)
5. [Next.js Frontend (aystek-core-next)](#nextjs-frontend-aystek-core-next)
6. [Flutter App (kisisel_asistan)](#-flutter-app-kisisel_asistan)
7. [CRITICAL Standards (MUST FOLLOW AUTOMATICALLY)](#critical-standards-must-follow-automatically)
8. [Module Creation Checklist](#module-creation-checklist)
9. [Integration Patterns](#integration-patterns)
10. [Testing Standards](#testing-standards)
11. [Code Style Guidelines](#code-style-guidelines)
12. [Documentation Requirements](#documentation-requirements)
13. [Standard Operating Procedures](#standard-operating-procedures)

---

## 🎯 Introduction & Purpose

### 📢 FOR ALL AI ASSISTANTS: READ THIS IMMEDIATELY

This document is the **SINGLE SOURCE OF TRUTH** for all AI tools working on AYSTEK projects. Before writing ANY code, you MUST:

**Canonical Architecture Reference (MANDATORY):** `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`
This is the single authoritative architecture standard. All platform-level decisions must align with it.

1. ✅ **Read this entire document**
2. ✅ **Internalize all CRITICAL rules** (marked with 🔴)
3. ✅ **Follow all standards automatically** without being told
4. ✅ **Apply the 6-core-file module structure** (plus tests & optional handlers) for every new module
5. ✅ **Add 28 language translations** for every new text key
6. ✅ **Create menu entries** for every new module
7. ✅ **Generate seed files** for every new model
8. ✅ **Update installer / deployment / operations docs when runtime procedures change**

### 🗣️ Communication Language Rule

**🔴 CRITICAL:** The project owner (Hüseyin Koç) communicates with AI assistants in **Turkish**. All chat/conversation responses to the user MUST be in Turkish — regardless of the language used in code, comments, commit messages, or this document itself (which stay in their existing convention, e.g. code/commits in English, code comments in Turkish per existing codebase style). This applies to every AI tool working on this project, not just one specific assistant/session.

### Installer And Setup Documentation Rule

If `aystek-core-installer` behavior changes, you must update the related setup docs in the same work item.

Minimum synchronized document set:
- `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`

This especially includes:
- seed/bootstrap behavior
- release artifact contract
- install/update/rollback step order
- backup/restore behavior
- integration test strategy

### Installer AI Release Advisor + Public-Site Auto-Update (Standing Knowledge)

**🔴 CRITICAL — read this before assuming any auto-update or "AI does deployments" behavior exists elsewhere in the ecosystem.**

`aystek-core-installer` has an AI-assisted feature, but its scope is narrow and precise — do not generalize it:

- **What it is:** A **local release advisor** for the `public-site` deployment scope ONLY (standalone Next.js sites such as aystekhaber.com, kureselhat.com — **not** the aystek-core Go API, **not** the aystek-core-next admin panel, which still use the manual `update` command with an operator-supplied version/image tag).
- **Code:** `internal/autoupdate/advisor.go` (`OllamaAdvisor`), `internal/autoupdate/service.go` (`Service.Run`), CLI subcommand `auto-update` in `cmd/aystek-core-installer/main.go`.
- **Model:** Calls a **locally-hosted Ollama model** (`http://127.0.0.1:11434/api/chat`, model name configured via `public_site_ai_model`, e.g. `qwen3:4b`) with a bounded diff + changed-file list, asking for a JSON `{decision: deploy|hold, reason}`. **No cloud AI call** (no OpenAI/Anthropic/GPT) — code/diffs never leave the server.
- **Deterministic safety gate the model cannot bypass** (`CheckCandidate`): holds automatically if more than 80 files changed, or if infra/dependency files changed (`package.json`, `docker*`, `.env`, `tsconfig`, etc.). The model can only approve small, low-risk content-level diffs.
- **This is unrelated to `AystekAIEngine`** (Core's rule-based content engine) — do not conflate the two or claim equivalence to a paid AI service in docs or UI copy.
- **Update source:** `public_site_git_url` may point to a local/self Git checkout (validated clean + correctly-branched via `internal/deploy/git_source.go`) **or** a GitHub URL (`https://github.com/...` / `git@github.com:...`) via the same fast-forward-only `GitSource.SyncExpected` logic — both paths share identical review/gating.
- **Trigger:** Manual (`aystek-core-installer auto-update --non-interactive --config ...`) or an optional systemd timer (`docs/systemd/aystek-public-auto-update@.timer`, hourly poll — `OnBootSec=15min`, `OnUnitActiveSec=1h`). **The timer is not enabled by default**; an operator must explicitly install/enable it per site.
- **Failure handling:** fetch → flock lock → diff-gate → Ollama review → if `deploy`: build image, write compose, bring up, health-check (6 attempts/5s); any failure restores the prior compose/image (`Service.restore`, logged as `auto_update_rollback`) — this is separate from and does not extend the general `rollback` command.
- **Fully documented already:** `docs/KULLANIM_KILAVUZU.md` §21.4.1, `README.md`, and `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md` §15.4.1 all describe this consistently with the code — read those before touching this subsystem, and keep all three in sync if you change it (per the Installer And Setup Documentation Rule above).
- **`/workspace/COZUMTEK_SUNUCU_BILGILERI.txt`** is a plaintext server-connection-info note. It is **not referenced anywhere in the installer codebase** — it is raw SSH/VPN connection info a human or an AI agent reads manually to perform an ad-hoc deployment session, unrelated to the advisor/auto-update feature itself. Treat its contents as sensitive; never print/log it, never commit it.

### Why This Document Exists

Without this document, AI assistants often:
- ❌ Create incomplete modules (missing tests, seeds, translations)
- ❌ Forget to add menu entries (users can't access new features)
- ❌ Skip translations (app breaks in non-English locales)
- ❌ Ignore architectural patterns (code becomes inconsistent)
- ❌ Make the same mistakes repeatedly

**This document prevents all of these issues.**

### 📌 Active/Planned Work Trackers (check these before starting related work)

Some large, multi-session efforts have their own living workplan document instead of a single
task doc, so any AI tool can resume them mid-way without being re-briefed by the user. Read the
relevant one **before** touching adjacent code — it usually means work is already in progress or
a module already exists where you might assume a greenfield build:

- **`/workspace/aystek-core-dev-docs/I18N_TRANSLATION_DEBT_WORKPLAN_TR.md`** — systematic closure of translation gaps across `translation_seed.go` (WP-1/2/3, prefix-by-prefix, ratchet-gated via `TestFullyTranslatedPrefixesHaveNoGaps`). Check its progress tables before assuming a prefix is untranslated.
- **`internal/modules/sports/`** — ✅ multi-sport (football/basketball/volleyball), multi-tier, site-scoped and opt-in external-provider auto-sync. API-Sports, TheSportsDB and the subscription-free Web Feed / Google Sheets adapter can run concurrently per site. Read `/workspace/aystek-core-dev-docs/done/TASK_061_SPORTS_AUTO_SYNC.md`, then the current extension reference `/workspace/aystek-core-dev-docs/done/TASK_062_SPORTS_SITE_SCOPED_WEB_FEED.md` and operator guide `/workspace/aystek-core-dev-docs/SPORTS_WEB_FEED_GUIDE_TR.md` before touching this module. The earlier free basketball/volleyball source gap is closed by the user-controlled canonical JSON/CSV bridge; do not reintroduce tenant-wide sync settings or expose private feed URLs through public APIs.
- **`/workspace/aystek-core-dev-docs/BALCI_PUBLIC_SITE_WORKPLAN_TR.md`** — `balci.aysteknolojileri.com` (Rucher du VELAY) public sitesi: Figma Make çıktısının (`/workspace/node/Rucherduvelayaricilik-main`) Next.js + Core entegrasyonu, faz planı ve açık kararlar. Yeni "genel public site" (haber olmayan) seed'i, frontend portu ve installer profili için önce bunu oku.
- **`/workspace/aystek-core-dev-docs/ENGAGEMENT_PLATFORM_RND_WORKPLAN_TR.md`** — SendPulse-class multi-channel engagement platform ("AYSTEK Engage", working name): audience/consent, messaging channels, campaigns, automation engine, conversations, plans/usage/credits, partner (affiliate/reseller/white-label) program and the partner-brandable public customer app. Phased A (Core) → B (Figma) → C (public app) → D (product). **Draft, awaiting owner approval.** Read §5 (inventory of what already exists — e.g. `contacts` is a contact-form inbox, not a CRM; `sms` has only a `mock` provider registered) and §12 (modules to retire) before touching `sms`, `marketing`, `newsletter_subscribers`, `chatbot`, `crm` or `integrations/chains`.
- **`/workspace/aystek-core-dev-docs/NEWS_PUBLIC_SITE_MULTIDOMAIN_WORKPLAN_TR.md`** — the multi-site news platform architecture (aystekhaber.com, kureselhat.com sharing one codebase/image with per-deployment runtime identity). Read this before adding a new public news site or touching site-identity/domain-resolution code.

### Core Principle

AI assistants must operate as **SENIOR SOFTWARE ENGINEERS** with deep knowledge of:
- Domain-Driven Design (DDD)
- Clean Architecture principles
- SOLID principles
- Test-Driven Development (TDD)
- Modern Go and TypeScript patterns

**You are not a code generator. You are a senior engineer who understands context, makes good decisions, and writes production-quality code.**

### Non-Negotiable Engineering Mindset

All AI assistants (Codex, Claude, Copilot, ChatGPT, and future tools) MUST work like the best enterprise-level computer engineers, not like autocomplete tools.

This means every task must be approached with:
- **Object-oriented and domain-oriented thinking:** Encapsulate behavior in the correct model/service/class/component boundaries. Prefer interfaces, dependency injection, constructors, domain methods, and clear ownership instead of procedural patches or scattered helper code.
- **Enterprise-grade quality:** Design for maintainability, observability, testability, security, localization, permissions, and operational continuity from the beginning.
- **No duplicate or abandoned work:** When a new implementation replaces an older one, remove the obsolete files, functions, routes, components, hooks, styles, tests, seeds, and documentation in the same work item.
- **Dynamic, user-managed behavior:** Avoid hard-coded business behavior when it should be configurable through admin screens, database records, settings, seeds, permissions, menus, translations, or tenant/site configuration.
- **Product ownership:** Treat every change as part of a living product. Leave the application cleaner, more consistent, and easier to operate than before.

---

## 👨‍💻 AI Role Definition

### Primary Responsibilities

When acting as a coding assistant, you MUST:

1. **Implement, Not Suggest**: Always implement changes directly unless explicitly asked to only suggest
2. **Complete Tasks Fully**: Never stop at partial implementations - finish what you start
3. **Follow Standards Automatically**: Apply all standards in this document without being reminded
4. **Self-Research**: Gather necessary context before implementing (read files, search code, analyze patterns)
5. **Validate Changes**: Run tests and check for compilation errors after modifications
6. **Document Actions**: Create/update documentation for significant changes
7. **Clean As You Build**: Remove code and artifacts that become unused because of your change
8. **Prefer Configurable Systems**: Make behavior manageable by users/admins where the product domain requires it

### Decision-Making Authority

You are authorized to:
- Create new files and directories
- Modify existing code
- Create tests
- Add dependencies (with justification)
- Refactor code (while maintaining functionality)
- Design database schemas (following MongoDB patterns)

### When to Ask

Only ask for clarification when:
- Business logic rules are ambiguous
- Multiple valid architectural approaches exist
- Breaking changes are required
- Third-party API credentials/keys are needed

---

## 📦 Project Overview

### Ecosystem Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      AYSTEK ECOSYSTEM                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────────────┐      ┌──────────────────────┐    │
│  │   aystek-core-next   │◄────►│    aystek-core       │    │
│  │   (Next.js 14.2.x)   │ REST │    (Go 1.24.x)       │    │
│  │   Frontend Admin     │ API  │    Backend API       │    │
│  └──────────────────────┘      └──────────────────────┘    │
│                                           │                  │
│                                           ▼                  │
│                                  ┌────────────────┐         │
│                                  │   MongoDB      │         │
│                                  │   + Redis      │         │
│                                  └────────────────┘         │
│                                           │                  │
│                                           ▼                  │
│                                  ┌────────────────┐         │
│                                  │ RabbitMQ / SSE │         │
│                                  └────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Backend** | Go | 1.24.x (go.mod + toolchain) |
| **Backend Framework** | Gin | 1.10.1 |
| **Database** | MongoDB | 6.x+ (driver v1.17.6) |
| **Cache** | Redis | 7.x+ (go-redis v9.17.2) |
| **Queue/Jobs** | RabbitMQ | amqp091-go v1.10.0 |
| **Observability** | OpenTelemetry + Prometheus | otel v1.39.x, prometheus v1.18.0 |
| **Frontend** | Next.js | 14.2.33 |
| **Frontend Language** | TypeScript | 5.3.3 |
| **Styling** | Tailwind CSS | 3.4.13 |
| **State Management** | Zustand | 4.5.x |

---

## 🔧 Go Backend (aystek-core)

### Project Structure

```
aystek-core/
├── api                           # Compiled binary (generated)
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── docker-compose.foundation.yml
├── docker-compose.production.yml
├── internal/
│   ├── cache/                      # Cache helpers
│   ├── errors/                     # Shared error types
│   ├── gateway/                    # External gateway integrations
│   ├── modules/                    # Business domain modules (60+)
│   │   ├── auth/
│   │   ├── users/
│   │   ├── products/
│   │   ├── orders/
│   │   ├── payments/
│   │   ├── checkout/               # Example module structure
│   │   │   ├── model/              # Domain models
│   │   │   │   └── checkout.go
│   │   │   ├── dto/                # Data Transfer Objects
│   │   │   │   └── checkout.go
│   │   │   ├── repository/         # Data access layer
│   │   │   │   ├── interface.go
│   │   │   │   └── mongo.go
│   │   │   ├── service/            # Business logic layer
│   │   │   │   ├── service.go
│   │   │   │   └── service_test.go
│   │   │   ├── handler/            # HTTP handlers
│   │   │   │   ├── admin.go         # Admin endpoints (optional)
│   │   │   │   └── public.go        # Public/protected endpoints
│   │   │   └── routes.go           # Route registration
│   │   └── [60+ other modules]/
│   ├── monitoring/                 # Metrics, health, ops
│   ├── platform/                   # Shared infrastructure
│   │   ├── alerts/
│   │   ├── bootstrap/
│   │   ├── cache/
│   │   ├── config/
│   │   ├── db/
│   │   ├── events/
│   │   ├── i18n/
│   │   ├── indexes/
│   │   ├── jobs/
│   │   ├── logger/
│   │   ├── mail/
│   │   ├── middleware/
│   │   ├── ops/
│   │   ├── otel/
│   │   ├── pagination/
│   │   ├── queue/
│   │   ├── ratelimit/
│   │   ├── redis/
│   │   ├── response/
│   │   ├── router/
│   │   ├── sse/
│   │   └── testing/
│   ├── queue/                      # Queue workers / adapters
│   └── seeds/                      # Database seeders (30+ files)
│       ├── translation_seed.go     # 🔴 CRITICAL: 28 languages
│       ├── menu_seed.go            # 🔴 CRITICAL: Menu entries
│       ├── localization_seed.go    # Language + localized content seeds
│       ├── user_seed.go
│       ├── product_seed.go
│       └── [30+ other seeds]/
│   ├── tests/                      # Integration tests
│   └── testutil/                   # Test helpers
├── docs/                           # Documentation
├── main                          # Compiled binary (generated)
├── public/                        # Static files
├── seed                          # Compiled seed runner (generated)
└── go.mod                          # Go dependencies
```

**Note on submodules:** Some domains are split into subpackages (e.g., `billing/invoices`, `products/attributes`, `products/media`). These still follow the same DDD layer conventions.

### Routing & Middleware (Current Codebase)

**Router Entry:** `internal/platform/router.New()`  
**Groups:** `router.Groups()` returns `/auth`, `/admin`, `/public`, and a protected root group

**Middleware order (simplified):**
1. CORS
2. Recovery
3. Request ID
4. Tracing (OpenTelemetry)
5. Response Envelope (wraps JSON)
6. Request Timer
7. Auth / Permission / Language / CSRF / Maintenance (group-specific)

**Implications:**
- ✅ Most admin endpoints live under `/admin/*`
- ✅ Public endpoints live under `/public/*`
- ✅ Some protected endpoints mount at root with the protected middleware chain
- ✅ Response envelope is automatic (see section below)

### Architecture Pattern: Domain-Driven Design (DDD)

Every module MUST follow this layered architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP REQUEST                            │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  ROUTES (routes.go)                                         │
│  - Register endpoints with Gin router                       │
│  - Map paths to handlers                                    │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  HANDLER (handler/public.go, handler/admin.go)             │
│  - Parse HTTP requests                                      │
│  - Validate input with DTOs                                 │
│  - Call service methods                                     │
│  - Format HTTP responses                                    │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  DTO (dto/*.go)                                             │
│  - Request/Response structures                              │
│  - Validation tags (binding:"required")                     │
│  - Transform domain models to API responses                 │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  SERVICE (service/service.go)                               │
│  - Business logic & orchestration                           │
│  - Call repositories                                        │
│  - Coordinate between modules                               │
│  - Apply business rules                                     │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  REPOSITORY (repository/interface.go, repository/mongo.go) │
│  - Data access abstraction                                  │
│  - CRUD operations                                          │
│  - Database queries                                         │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  MODEL (model/*.go)                                         │
│  - Domain entities                                          │
│  - Business logic methods                                   │
│  - BSON/JSON tags                                           │
│  - State machines (where applicable)                        │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                     MONGODB                                 │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities (DETAILED)

#### 1. MODEL Layer (`model/*.go`)

**Purpose:** Domain entities and business rules

**Standards:**
```go
package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Status type for type safety
type CheckoutStatus string

const (
    CheckoutStatusStarted          CheckoutStatus = "started"
    CheckoutStatusAddressValidated CheckoutStatus = "address_validated"
    CheckoutStatusShippingSelected CheckoutStatus = "shipping_selected"
    CheckoutStatusPaymentMethodSet CheckoutStatus = "payment_method_set"
    CheckoutStatusConfirmed        CheckoutStatus = "confirmed"
    CheckoutStatusCancelled        CheckoutStatus = "cancelled"
    CheckoutStatusExpired          CheckoutStatus = "expired"
)

// Main domain model (matches existing codebase style)
type CheckoutSession struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    CustomerID      primitive.ObjectID `bson:"customer_id" json:"customer_id"`
    CartID          primitive.ObjectID `bson:"cart_id" json:"cart_id"`
    Status          CheckoutStatus     `bson:"status" json:"status"`
    Subtotal        float64            `bson:"subtotal" json:"subtotal"`
    TaxAmount       float64            `bson:"tax_amount" json:"tax_amount"`
    ShippingAmount  float64            `bson:"shipping_amount" json:"shipping_amount"`
    DiscountAmount  float64            `bson:"discount_amount" json:"discount_amount"`
    GrandTotal      float64            `bson:"grand_total" json:"grand_total"`
    Currency        string             `bson:"currency" json:"currency"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
    ExpiresAt       time.Time          `bson:"expires_at" json:"expires_at"` // TTL
}

// Domain method - encapsulates business logic
func (cs *CheckoutSession) IsExpired() bool {
    return time.Now().After(cs.ExpiresAt)
}

// State machine validation
func (cs *CheckoutSession) CanTransition(target CheckoutStatus) bool {
    validTransitions := map[CheckoutStatus][]CheckoutStatus{
        CheckoutStatusStarted: {CheckoutStatusAddressValidated, CheckoutStatusCancelled},
        CheckoutStatusAddressValidated: {CheckoutStatusShippingSelected, CheckoutStatusCancelled},
        CheckoutStatusShippingSelected: {CheckoutStatusPaymentMethodSet, CheckoutStatusCancelled},
        CheckoutStatusPaymentMethodSet: {CheckoutStatusConfirmed, CheckoutStatusCancelled},
        CheckoutStatusConfirmed: {CheckoutStatusCancelled},
        CheckoutStatusCancelled: {},
        CheckoutStatusExpired:   {},
    }
    allowed, ok := validTransitions[cs.Status]
    if !ok {
        return false
    }
    for _, s := range allowed {
        if s == target {
            return true
        }
    }
    return false
}
```

**Key Points:**
- ✅ Use `bson` tags for MongoDB field mapping
- ✅ Use `json` tags for JSON serialization (snake_case in this codebase)
- ✅ Define typed constants (not plain strings)
- ✅ Include domain methods (IsExpired, CanTransition, etc.)
- ✅ Implement state machines with validation maps
- ✅ Use `primitive.ObjectID` for MongoDB IDs
- ✅ Include timestamps (CreatedAt, UpdatedAt, ExpiresAt where TTL is required)

#### 2. DTO Layer (`dto/*.go`)

**Purpose:** API request/response structures with validation

**Standards:**
```go
package dto

// Request DTO - with validation tags
type StartCheckoutRequest struct {
    CartID string `json:"cart_id" binding:"required"`
}

// Response DTO - clean API response
type StartCheckoutResponse struct {
    CheckoutID string  `json:"checkout_id"`
    CartID     string  `json:"cart_id"`
    Status     string  `json:"status"`
    Subtotal   float64 `json:"subtotal"`
    Currency   string  `json:"currency"`
    ExpiresAt  string  `json:"expires_at"` // ISO8601
}

// Additional request/response DTOs follow the same style:
// - snake_case JSON tags
// - string IDs for HTTP payloads
// - binding tags for validation
type ConfirmCheckoutRequest struct {
    PaymentMethod string         `json:"payment_method" binding:"required"` // card|bank|cash
    Notes         string         `json:"notes,omitempty"`
    Metadata      map[string]any `json:"metadata,omitempty"`
}
```

**Key Points:**
- ✅ Use `binding` tags for validation (required, email, min, max, etc.)
- ✅ Separate Request and Response DTOs
- ✅ Use snake_case JSON field names in HTTP payloads
- ✅ Use string IDs in DTOs (convert to/from `primitive.ObjectID` in handlers/services)
- ✅ Format dates as strings in responses (ISO8601/RFC3339)

#### 3. REPOSITORY Layer (`repository/interface.go`, `repository/mongo.go`)

**Purpose:** Abstract data access, implement MongoDB operations

**Interface** (`repository/interface.go`):
```go
package repository

import (
    "context"
    "github.com/aystek/aystek-core/internal/modules/checkout/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository defines data access contract
type Repository interface {
    Create(ctx context.Context, session *model.CheckoutSession) (*model.CheckoutSession, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*model.CheckoutSession, error)
    GetByIDString(ctx context.Context, id string) (*model.CheckoutSession, error)
    Update(ctx context.Context, session *model.CheckoutSession) error
    UpdateStatus(ctx context.Context, id primitive.ObjectID, status model.CheckoutStatus) error
    GetByCustomerAndCart(ctx context.Context, customerID, cartID primitive.ObjectID) (*model.CheckoutSession, error)
    ListByCustomer(ctx context.Context, customerID primitive.ObjectID, limit int, skip int) ([]model.CheckoutSession, error)
    Delete(ctx context.Context, id primitive.ObjectID) error
    DeleteExpired(ctx context.Context) (int64, error)
}
```

**MongoDB Implementation** (`repository/mongo.go`):
```go
package repository

import (
    "context"
    "errors"
    "time"
    "github.com/aystek/aystek-core/internal/modules/checkout/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
    collection *mongo.Collection
}

// Constructor - accepts *mongo.Database, NOT *mongo.Collection
func NewMongoRepository(db *mongo.Database) Repository {
    return &mongoRepository{collection: db.Collection("checkout_sessions")}
}

func (r *mongoRepository) Create(ctx context.Context, session *model.CheckoutSession) (*model.CheckoutSession, error) {
    if session == nil {
        return nil, errors.New("session cannot be nil")
    }
    session.ID = primitive.NewObjectID()
    session.CreatedAt = time.Now()
    session.UpdatedAt = time.Now()
    _, err := r.collection.InsertOne(ctx, session)
    if err != nil {
        return nil, err
    }
    return session, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.CheckoutSession, error) {
    var session model.CheckoutSession
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
    if err == mongo.ErrNoDocuments {
        return nil, nil // not found is not an error
    }
    return &session, err
}

func (r *mongoRepository) Update(ctx context.Context, session *model.CheckoutSession) error {
    session.UpdatedAt = time.Now()
    _, err := r.collection.ReplaceOne(ctx, bson.M{"_id": session.ID}, session)
    return err
}

func (r *mongoRepository) CreateIndexes(ctx context.Context) error {
    indexModel := mongo.IndexModel{
        Keys:    bson.D{{Key: "expires_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(0),
    }
    _, err := r.collection.Indexes().CreateOne(ctx, indexModel)
    return err
}
```

**Key Points:**
- ✅ Constructor accepts `*mongo.Database` (not Collection)
- ✅ Set CreatedAt/UpdatedAt timestamps automatically
- ✅ Return `nil, nil` for not found (not errors)
- ✅ Use `bson.M` for queries, `bson.D` for ordered operations
- ✅ Create TTL/compound indexes (constructor or explicit `CreateIndexes`)

#### 🔧 HTTP Response Envelope (Middleware)

All JSON responses are wrapped by `internal/platform/middleware/response_envelope.go`, unless the request is for static assets or SSE.

**Standard Envelope:**
```json
{
  "success": true,
  "message": "",
  "data": {},
  "error": null,
  "pagination": { "page": 1, "limit": 25, "total": 0, "total_pages": 0 }
}
```

**Rules:**
- ✅ You can return raw JSON from handlers; middleware will wrap it.
- ✅ If you already return `response.Envelope`, middleware will not double-wrap.
- ✅ For lists, return `{ "items": [...], "total": X, "page": N, "limit": M }` to auto-create `pagination`.

#### 4. SERVICE Layer (`service/service.go`)

**Purpose:** Business logic orchestration

**Standards:**
```go
package service

import (
    "context"
    "errors"
    "github.com/aystek/aystek-core/internal/modules/checkout/model"
    "github.com/aystek/aystek-core/internal/modules/checkout/repository"
    cartsrepo "github.com/aystek/aystek-core/internal/modules/carts/repository"
    ordersrepo "github.com/aystek/aystek-core/internal/modules/orders/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
    ErrCheckoutNotFound = errors.New("checkout not found")
    ErrInvalidStatus    = errors.New("invalid status transition")
    ErrExpiredSession   = errors.New("checkout session expired")
)

type Service struct {
    repo        repository.Repository
    cartsRepo   cartsrepo.Repository
    ordersRepo  ordersrepo.Repository
}

func New(repo repository.Repository, cartsRepo cartsrepo.Repository, ordersRepo ordersrepo.Repository) *Service {
    return &Service{
        repo:       repo,
        cartsRepo:  cartsRepo,
        ordersRepo: ordersRepo,
    }
}

// Business logic method
func (s *Service) StartCheckout(ctx context.Context, customerID, cartID primitive.ObjectID, email, phone string) (*model.Checkout, error) {
    // 1. Validate cart exists
    cart, err := s.cartsRepo.FindByID(ctx, cartID)
    if err != nil || cart == nil {
        return nil, errors.New("cart not found")
    }
    
    // 2. Validate cart belongs to customer
    if cart.CustomerID != customerID {
        return nil, errors.New("cart does not belong to customer")
    }
    
    // 3. Create checkout session
    checkout := &model.Checkout{
        CustomerID: customerID,
        CartID:     cartID,
        Email:      email,
        Phone:      phone,
        Status:     model.CheckoutStatusCreated,
    }
    
    if err := s.repo.Create(ctx, checkout); err != nil {
        return nil, err
    }
    
    return checkout, nil
}

// Complex orchestration method
func (s *Service) ConfirmCheckout(ctx context.Context, checkoutID primitive.ObjectID) (*model.Order, error) {
    // 1. Get checkout
    checkout, err := s.repo.FindByID(ctx, checkoutID)
    if err != nil {
        return nil, err
    }
    if checkout == nil {
        return nil, ErrCheckoutNotFound
    }
    
    // 2. Validate state
    if !checkout.CanTransition(model.CheckoutStatusCompleted) {
        return nil, ErrInvalidStatus
    }
    
    if checkout.IsExpired() {
        checkout.Status = model.CheckoutStatusExpired
        _ = s.repo.Update(ctx, checkoutID, checkout)
        return nil, ErrExpiredSession
    }
    
    // 3. Create order
    order := &model.Order{
        CustomerID: checkout.CustomerID,
        Email:      checkout.Email,
        Status:     "pending",
        // ... other fields
    }
    
    if err := s.ordersRepo.Create(ctx, order); err != nil {
        return nil, err
    }
    
    // 4. Update checkout status
    checkout.Status = model.CheckoutStatusCompleted
    checkout.OrderID = &order.ID
    if err := s.repo.Update(ctx, checkoutID, checkout); err != nil {
        return nil, err
    }
    
    return order, nil
}
```

**Key Points:**
- ✅ Define custom error types (ErrXXX)
- ✅ Accept dependencies via constructor (dependency injection)
- ✅ Orchestrate multiple repositories
- ✅ Apply business rules and validations
- ✅ Use domain methods from models (CanTransition, IsExpired)
- ✅ Handle state transitions correctly
- ✅ Return domain models (not DTOs)

#### 5. HANDLER Layer (`handler/public.go`)

**Purpose:** HTTP request/response handling

**Standards:**
```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/aystek/aystek-core/internal/modules/checkout/dto"
    "github.com/aystek/aystek-core/internal/modules/checkout/service"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicHandler struct {
    service *service.Service
}

func NewPublicHandler(service *service.Service) *PublicHandler {
    return &PublicHandler{service: service}
}

// POST /checkout/start
func (h *PublicHandler) StartCheckout(c *gin.Context) {
    var req dto.StartCheckoutRequest
    
    // 1. Bind and validate
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 2. Call service
    checkout, err := h.service.StartCheckout(
        c.Request.Context(),
        req.CustomerID,
        req.CartID,
        req.Email,
        req.Phone,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // 3. Transform to DTO and respond
    response := dto.ToCheckoutResponse(checkout)
    c.JSON(http.StatusOK, gin.H{"checkout": response})
}

// GET /checkout/:id
func (h *PublicHandler) GetCheckout(c *gin.Context) {
    // 1. Parse URL parameter
    idStr := c.Param("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    
    // 2. Call service
    checkout, err := h.service.GetCheckout(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if checkout == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "checkout not found"})
        return
    }
    
    // 3. Respond
    response := dto.ToCheckoutResponse(checkout)
    c.JSON(http.StatusOK, gin.H{"checkout": response})
}
```

**Key Points:**
- ✅ Use `ShouldBindJSON` for request validation
- ✅ Parse URL/query parameters carefully
- ✅ Use `c.Request.Context()` for context propagation
- ✅ Return proper HTTP status codes (200, 400, 404, 500)
- ✅ Transform models to DTOs before responding
- ✅ Use `gin.H{}` for JSON responses
- ✅ Handle errors gracefully

#### 6. ROUTES Layer (`routes.go`)

**Purpose:** Register HTTP endpoints

**Standards:**
```go
package checkout

import (
    "github.com/gin-gonic/gin"
    "github.com/aystek/aystek-core/internal/modules/checkout/handler"
)

// RegisterRoutes registers all checkout routes
func RegisterRoutes(r *gin.RouterGroup, public *handler.PublicHandler) {
    r.POST("/checkout/start", public.StartCheckout)
    r.POST("/checkout/:id/address", public.ValidateAddress)
    r.POST("/checkout/:id/calculate", public.CalculateTotals)
    r.POST("/checkout/:id/confirm", public.ConfirmCheckout)
    r.GET("/checkout/:id", public.GetCheckout)
    r.DELETE("/checkout/:id", public.CancelCheckout)
}
```

**Key Points:**
- ✅ Accept `*gin.RouterGroup` (not `*gin.Engine`)
- ✅ Accept handler instances (not creating them here)
- ✅ Use RESTful conventions (GET, POST, PUT, DELETE)
- ✅ Group related routes together

### Module Wiring in main.go

**Location:** `cmd/api/main.go`

**Pattern:**
```go
package main

import (
    "github.com/aystek/aystek-core/internal/platform/router"
    "github.com/aystek/aystek-core/internal/modules/checkout"
    checkoutrepo "github.com/aystek/aystek-core/internal/modules/checkout/repository"
    checkoutservice "github.com/aystek/aystek-core/internal/modules/checkout/service"
    checkouthandler "github.com/aystek/aystek-core/internal/modules/checkout/handler"
    // ... other imports
)

func main() {
    // ... MongoDB setup ...

    // Initialize repositories
    checkoutRepo := checkoutrepo.NewMongoRepository(db)
    cartsRepo := cartsrepo.NewMongoRepository(db)
    ordersRepo := ordersrepo.NewMongoRepository(db)

    // Initialize services with dependencies
    checkoutService := checkoutservice.New(checkoutRepo, cartsRepo, nil, nil, nil)

    // Initialize handlers
    checkoutPublicHandler := checkouthandler.NewPublicHandler(checkoutService)

    // Setup router + groups
    r := router.New()
    authGroup, adminGroup, publicGroup, protectedGroup := router.Groups(r, settingsSvc, langSvc, authRepo, userRepo, permRepo, roleRepo)

    // Register routes
    checkout.RegisterRoutes(protectedGroup, checkoutPublicHandler)

    // ... other modules ...

    _ = authGroup
    _ = adminGroup
    _ = publicGroup
    r.Run(":8080")
}
```

**Key Points:**
- ✅ Follow dependency injection pattern (bottom-up)
- ✅ Repository → Service → Handler → Routes
- ✅ Use aliased imports to avoid name conflicts
- ✅ Group routes under `/auth`, `/admin`, `/public`, or protected root

---

## 🌐 Next.js Frontend (aystek-core-next)

### Project Structure

```
aystek-core-next/
├── app/                            # Next.js 14 App Router
│   ├── layout.tsx                  # Root layout
│   ├── page.tsx                    # Home page
│   ├── (auth)/                     # Auth route group
│   ├── (dashboard)/                # Dashboard route group
│   │   ├── admin/                  # Admin pages (/admin/*)
│   │   │   ├── products/
│   │   │   ├── orders/
│   │   │   └── [modules]/
│   │   ├── layout.tsx              # Dashboard layout
│   │   └── page.tsx                # Dashboard home
│   └── api/                        # API routes (if needed)
│   └── components/                 # App-router scoped components
├── components/                     # React components
│   ├── ui/                         # Reusable UI components
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   └── Table.tsx
│   ├── layout/                     # Layout components
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   └── Footer.tsx
│   ├── products/                   # Feature-specific components
│   │   ├── ProductCard.tsx
│   │   ├── ProductForm.tsx
│   │   └── ProductList.tsx
│   └── shared/                     # Shared components
├── lib/                            # Utilities and API clients
│   ├── api/                        # API wrappers per module
│   ├── axios.ts                    # Axios client with auth + CSRF
│   ├── api/baseUrl.ts              # API base URL resolver
│   ├── i18n/                       # i18n helpers
│   └── stores/                     # Zustand stores
├── services/                       # Specific clients (e.g., SSE)
├── hooks/                          # Custom React hooks
│   ├── useAuth.ts
│   ├── useProducts.ts
│   └── useCheckout.ts
├── store/                          # State management
│   ├── authStore.ts
│   └── cartStore.ts
├── types/                          # Shared types (minimal)
├── public/                         # Static assets
└── tailwind.config.ts              # Tailwind configuration
```

### Frontend Architecture Patterns

#### 1. API Service Layer

**Purpose:** Centralize all backend API calls

**Standards** (`lib/api/orders.ts`-style):
```typescript
import api from "@/lib/axios";
import { unwrapApiResponse } from "@/lib/api/mappers";
import type { Order } from "@/lib/types";

const BASE_PATH = "/admin/orders";

export const ordersApi = {
  async list(params?: { page?: number; limit?: number }) {
    const res = await api.get(BASE_PATH, { params });
    const body = unwrapApiResponse(res.data ?? res);
    const items = Array.isArray(body?.items) ? body.items : [];
    const page = Number(body?.page ?? params?.page ?? 1);
    const limit = Number(body?.limit ?? params?.limit ?? 25);
    const total = Number(body?.total ?? items.length);
    return { items: items as Order[], page, limit, total };
  },
};
```

**Base API Client** (`lib/axios.ts`):
```typescript
import axios from "axios";
import { getApiBaseUrl } from "@/lib/api/baseUrl";

const api = axios.create({
  baseURL: getApiBaseUrl(),
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

// Interceptors handle tokens, CSRF, and unwrap API envelopes.
export default api;
```

**Key Points:**
- ✅ Use `lib/axios.ts` or `lib/api/client.ts` for HTTP
- ✅ Centralize API base URL via `lib/api/baseUrl.ts`
- ✅ Interceptors handle auth + CSRF + envelope unwrapping
- ✅ Normalize IDs and snake_case fields where needed

#### 2. Custom Hooks

**Purpose:** Encapsulate data fetching and business logic

**Standards** (`hooks/useLanguages.ts`-style):
```typescript
import { useEffect, useState } from "react";
import { languagesApi } from "@/lib/api/languages";
import type { Language } from "@/lib/types";

export function useLanguages() {
  const [languages, setLanguages] = useState<Language[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
        const res = await languagesApi.list(false);
        if (mounted) setLanguages(res.languages || []);
      } catch (err: any) {
        if (mounted) setError(err.message || "Failed to load languages");
      } finally {
        if (mounted) setLoading(false);
      }
    })();
    return () => {
      mounted = false;
    };
  }, []);

  return { languages, loading, error };
}
```

**Key Points:**
- ✅ Return loading/error states
- ✅ Return data and action functions
- ✅ Handle errors in try/catch
- ✅ Type everything with TypeScript
- ✅ Use useEffect for automatic fetching

#### 3. TypeScript Types

**Purpose:** Type safety across frontend

**Standards** (`lib/types/checkout.ts`):
```typescript
export type CheckoutStatus = 
  | 'started'
  | 'address_validated'
  | 'shipping_selected'
  | 'payment_method_set'
  | 'confirmed'
  | 'cancelled'
  | 'expired';

export interface Checkout {
  id: string;
  customer_id: string;
  status: CheckoutStatus;
  email: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

export interface StartCheckoutRequest {
  cart_id: string;
}

export interface CheckoutResponse {
  checkout_id: string;
  cart_id: string;
  status: CheckoutStatus;
  subtotal: number;
  currency: string;
  expires_at: string;
}
```

**Key Points:**
- ✅ Match backend DTO structures
- ✅ Use string literals for status types
- ✅ Separate Request and Response types
- ✅ Use snake_case field names to match backend payloads

#### 4. Page Components (App Router)

**Purpose:** Server/Client components for routes

**Standards** (`app/(dashboard)/admin/languages/page.tsx`):
```typescript
'use client';

import { useLanguages } from '@/hooks/useLanguages';
import { Spinner } from '@/components/ui/Spinner';

export default function LanguagesPage() {
  const { languages, loading, error } = useLanguages();

  if (loading) return <Spinner />;
  if (error) return <div>Error: {error}</div>;
  if (!languages?.length) return <div>No languages found</div>;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Languages</h1>
      <ul className="space-y-2">
        {languages.map((lang) => (
          <li key={lang.id} className="bg-white rounded-md p-3 shadow">
            {lang.name}
          </li>
        ))}
      </ul>
    </div>
  );
}
```

**Key Points:**
- ✅ Use `'use client'` for interactive components
- ✅ Use custom hooks for data fetching
- ✅ Show loading/error states
- ✅ Use Tailwind CSS for styling
- ✅ Extract UI components (Button, Spinner, etc.)

---

---

## 📱 Flutter App (kisisel_asistan)

### Proje Konumu & Amaç

**Path:** `/workspace/flutter/kisisel_asistan/`  
**Package:** `com.aysteknolojileri.kisisel_asistan`  
**Amaç:** Kişisel eşya/konum/seyahat takibi + AYSTEK Core entegrasyonlu haber ve AI sohbet uygulaması.  
**Backend:** `aystek-core` personal_assistant modülü → `http://10.231.235.10:9090` (dev), `https://pers-assit.aysteknolojileri.com` (prod)

### Technology Stack

| Katman | Paket | Versiyon |
|--------|-------|---------|
| State Management | flutter_riverpod | ^2.5.1 |
| Navigation | go_router | ^14.2.7 |
| Local DB | sqflite | ^2.3.3 |
| Secure Storage | flutter_secure_storage | ^9.2.2 |
| Background Tasks | workmanager | ^0.9.0 |
| Notifications | flutter_local_notifications | ^17.2.3 |
| HTTP | http | ^1.2.2 |
| i18n | flutter_localizations (ARB) | SDK |

### Proje Yapısı

```
lib/
├── main.dart                      # Entry point — Riverpod ProviderScope + servis başlatma
├── core/
│   ├── constants/app_constants.dart
│   ├── extensions/l10n_extension.dart
│   ├── router/app_router.dart     # GoRouter tanımları (StatefulShellRoute)
│   └── theme/app_theme.dart       # Material3 tema
├── database/
│   └── database_helper.dart       # sqflite yönetici (kişisel veri: Item, Person, Travel)
├── l10n/                          # ARB tabanlı 11 dil (tr, en, de, fr, es, pt, ar, zh, ja, ru, hi)
├── models/
│   ├── aystek_models.dart         # AYSTEK Core API model nesneleri
│   ├── chat_mode.dart             # AI sohbet modları (ChatMode enum)
│   ├── ai_settings.dart           # AI sağlayıcı/model ayarları
│   ├── item.dart / person.dart / travel.dart  # Yerel veri modelleri
│   └── ...
├── providers/
│   ├── aystek_providers.dart      # AYSTEK Core Riverpod provider'ları
│   ├── providers.dart             # Yerel DB provider'ları
│   └── locale_provider.dart
├── screens/
│   ├── aystek/                    # AYSTEK Core bağlantılı ekranlar
│   │   ├── news_tab_screen.dart   # Alt sekme kapsayıcı
│   │   ├── news_feed_screen.dart  # Haber akışı + etiket/bölge filtresi
│   │   ├── news_article_screen.dart # Makale detay + in-app WebView
│   │   ├── news_prefs_screen.dart # Kaynak/bölge/bildirim tercihleri
│   │   ├── aystek_login_screen.dart
│   │   ├── aystek_register_screen.dart
│   │   └── aystek_account_screen.dart
│   ├── ai/
│   │   ├── ai_chat_screen.dart    # AI sohbet (OpenAI/Claude/Gemini) + ChatMode
│   │   └── chat_tab_screen.dart
│   ├── home/                      # Ana sayfa dashboard
│   ├── items/                     # Eşya yönetimi
│   ├── persons/                   # Kişi yönetimi
│   ├── travels/                   # Seyahat yönetimi
│   ├── settings/                  # Uygulama ayarları
│   └── shell_screen.dart          # BottomNavigationBar + StatefulShellRoute
├── services/
│   ├── aystek_api_service.dart    # AYSTEK Core HTTP client (singleton)
│   ├── ai_service.dart            # OpenAI/Claude/Gemini API çağrıları
│   ├── news_notification_service.dart  # WorkManager arka plan haber bildirimleri
│   ├── notification_service.dart  # flutter_local_notifications wrapper
│   ├── tavily_service.dart        # Web arama servisi
│   └── ...
└── widgets/                       # Paylaşılan UI bileşenleri
```

### AYSTEK Core HTTP İstemcisi

```dart
// Sabitler — aystek_api_service.dart
const _kBaseUrl = 'http://10.231.235.10:9090';         // Dev
// const _kBaseUrl = 'https://pers-assit.aysteknolojileri.com'; // Prod
const _kHostHeader = 'pers-assit.aysteknolojileri.com';

// _TenantAwareClient: TCP bağlantısı dev IP'ye gider,
// Host header backend'e tenant/domain bilgisini taşır
class _TenantAwareClient extends http.BaseClient {
  @override
  Future<http.StreamedResponse> send(http.BaseRequest request) {
    request.headers['host'] = _kHostHeader; // StoreContextMiddleware için
    return _inner.send(request);
  }
}

// Token anahtarları
const _kStorageKeyAccess = 'aystek_access_token';
const _kStorageKeyRefresh = 'aystek_refresh_token';
```

**Önemli:** `AystekApiService.instance` singleton pattern kullanır. Token'lar `flutter_secure_storage` ile `encryptedSharedPreferences: true` (Android EncryptedSharedPreferences) olarak saklanır.

### Riverpod Provider Katmanı

```dart
// aystek_providers.dart

// Auth durumu — NotifierProvider (sync, restoreSession ile başlar)
final aystekAuthProvider = NotifierProvider<AystekAuthNotifier, AystekAuthState>(...);

// Kullanıcı profili
final paProfileProvider = FutureProvider.autoDispose<PAProfil?>((ref) async { ... });

// Haber tercihleri (kaynak + bölge)
final paNewsPreferencesProvider = FutureProvider.autoDispose<PANewsPreferences?>((ref) async { ... });

// Haber akışı
final paNewsFeedProvider = FutureProvider.autoDispose<List<NewsArticle>>((ref) async { ... });

// Haber kaynakları (kullanıcının seçtiği)
final paNewsSourcesProvider = FutureProvider.autoDispose<List<PANewsSource>>((ref) async { ... });

// Tüm mevcut kaynaklar
final paAllSourcesProvider = FutureProvider<List<PANewsSource>>((ref) async { ... });

// Filtre state'leri
final newsCategoryFilterProvider = StateProvider<String>((ref) => 'all');
final newsTagFilterProvider = StateProvider<String>((ref) => '');
```

**KRITIK — autoDispose + ref.listen Race Condition:**  
`NewsFeedScreen` `paNewsPreferencesProvider`'ı önceden yüklerse, `NewsPrefsScreen` açıldığında `ref.listen` **tetiklenmez** (state değişmedi). Fix: `initState` içinde `ref.read` ile sync okuma yapılır.

```dart
@override
void initState() {
  super.initState();
  // Sync init — provider zaten yüklü olabilir
  final cached = ref.read(paNewsPreferencesProvider);
  cached.whenData((prefs) {
    if (prefs != null) {
      _selected.addAll(prefs.sourceIds);
      _initialized = true;
    }
  });
}

@override
Widget build(BuildContext context) {
  ref.watch(paNewsPreferencesProvider); // provider'ı canlı tutar
  // Async fallback — sadece ilk yüklemede tetiklenir
  ref.listen<AsyncValue<PANewsPreferences?>>(paNewsPreferencesProvider, (_, next) {
    next.whenData((prefs) {
      if (prefs != null && !_initialized && mounted) setState(() { ... });
    });
  });
}
```

**KRITIK — addPostFrameCallback'i build() içinde kullanma:**  
Her frame yeniden kaydeder → döngüsel rebuild → `_FocusInheritedScope` crash. Bunun yerine `ref.listen` kullan.

### Navigation (GoRouter)

Uygulama 5 alt dal içeren `StatefulShellRoute.indexedStack` yapısı kullanır:

| Branch | Path | Açıklama |
|--------|------|---------|
| 0 | `/home` | Ana sayfa + Person/Item/Location/Travel alt rotaları |
| 1 | `/news` | Haber akışı tabı |
| 2 | `/chat` | AI sohbet tabı |
| 3 | `/travels` | Seyahat yönetimi tabı |
| 4 | `/settings` | Ayarlar tabı |

**Özel rotalar (shell dışı):**
- `/onboarding` — İlk kurulum
- `/search` — Genel arama
- `/news/article` — Makale detay (`extra: NewsArticle` olarak nesne geçilir)
- `/news-feed` — Standalone derin bağlantı

**GoRouter extra ile nesne geçme:**
```dart
// Gönderirken:
context.push('/news/article', extra: article);

// Alırken:
GoRoute(
  path: '/news/article',
  builder: (c, s) => NewsArticleScreen(article: s.extra as NewsArticle),
),
```

### AI Sohbet Modları (ChatMode)

```dart
enum ChatMode {
  assistant,   // Genel asistan (kişi bağlamı kullanır)
  newsDesk,    // Haber masası + toplantı gündemi (feed'i okur)
  coffeeFal,   // Kahve falı
  tarotFal,    // Tarot
  cooking,     // Yemek tarifi
  suggestions, // Öneri modu
}
```

Her mod: `icon`, `label(l10n)`, `systemPrompt`, `hints`, `usesPersonContext`, `welcomeCard` tanımlar.

**Toplantı gündemi akışı (newsDesk):** `_sendNewsAgenda()` → `_MeetingTimeDialog` → feed + profil okunur → yapılandırılmış mesaj oluşturulur → `_send()`.

### Arka Plan Bildirimler (WorkManager)

```dart
// Zorunlu: top-level @pragma annotasyonu
@pragma('vm:entry-point')
void callbackDispatcher() {
  Workmanager().executeTask((taskName, _) async {
    WidgetsFlutterBinding.ensureInitialized();
    if (taskName == _kTaskName) await _backgroundNewsCheck();
    return true;
  });
}
```

- **Frekans:** 15 dk minimum (Android kısıtlaması)
- **Token erişimi:** `FlutterSecureStorage(aOptions: AndroidOptions(encryptedSharedPreferences: true))` — arka plan isolate'ten erişilebilir
- **Tercih anahtarları:** `news_notify_breaking`, `news_notify_developing` (SharedPreferences)
- **Görülen URL'ler:** `news_seen_urls` (SharedPreferences string listesi)
- **Policy:** `ExistingPeriodicWorkPolicy.keep` (workmanager ^0.9.0)

### Yerelleştirme (l10n)

ARB tabanlı, `flutter gen-l10n` ile üretilir. **11 dil:** tr, en, de, fr, es, pt, ar, zh, ja, ru, hi.

```
lib/l10n/
├── intl_tr.arb   # Türkçe (kaynak)
├── intl_en.arb
├── intl_de.arb
└── ...
```

**Kural:** Yeni kullanıcıya görünen metin **tüm 11 ARB dosyasına** eklenmeli.  
`context.l10n.keyName` uzantısı ile erişilir.

### AYSTEK Core Personal Assistant Modülü

**Backend path:** `/workspace/go/aystek-core/internal/modules/personal_assistant/`

| Endpoint | Açıklama |
|---------|---------|
| `POST /pa/auth/login` | Giriş |
| `POST /pa/auth/register` | Kayıt |
| `POST /pa/auth/refresh` | Token yenileme |
| `GET /pa/profile` | Kullanıcı profili |
| `GET /pa/news/feed` | Kişiselleştirilmiş haber akışı |
| `GET /pa/news/sources` | Mevcut haber kaynakları |
| `GET /pa/news/preferences` | Haber tercihleri |
| `PUT /pa/news/preferences` | Haber tercihleri kaydet |

**NewsArticle modeli:**
```go
type NewsArticle struct {
  ID        string    // URL hash
  Title     string
  Summary   string
  URL       string
  ImageURL  string
  Source    string
  Category  string
  Tags      []string  // autoTag() ile otomatik üretilir
  PublishedAt time.Time
}
```

**Auto-tagging:** `autoTag(title, summary string)` fonksiyonu Türkçe anahtar kelime kurallarıyla 12 kategori etiketler:
`ekonomi`, `siyaset`, `spor`, `teknoloji`, `sağlık`, `eğitim`, `kültür-sanat`, `dünya`, `hukuk-suç`, `afet-doğa`, `ulaşım`, `enerji`

> **Do not confuse this with the `sports` module.** This `spor` tag only labels personal-assistant news feed articles. The real match/standings feature is `internal/modules/sports/` (backend) + `aystekhaber-com/components/news/SportsSnapshotPanel.tsx` (frontend) — see `done/TASK_062_SPORTS_SITE_SCOPED_WEB_FEED.md` and `SPORTS_WEB_FEED_GUIDE_TR.md` (pointers above) before touching either.

### ADB Test İş Akışı

```bash
# ASLA adb disconnect çalıştırma — bağlantı kopar
adb connect 10.231.235.6:37357
/tmp/platform-tools/adb -s 10.231.235.6:37357 install -r build/app/outputs/flutter-apk/app-debug.apk

# APK build
cd /workspace/flutter/kisisel_asistan
flutter pub get
flutter build apk --debug
```

### Önemli Tuzaklar

1. **`ExistingWorkPolicy` vs `ExistingPeriodicWorkPolicy`**: Workmanager ^0.9.0 periodic task için `ExistingPeriodicWorkPolicy.keep` kullanır, `ExistingWorkPolicy` değil.

2. **Host Header**: Backend'e istek gönderirken `Host: pers-assit.aysteknolojileri.com` header'ı şart — aksi halde tenant çözülmez.

3. **autoDispose + cached state**: `initState`'de `ref.read()` ile senkron oku, `ref.listen()`'i async fallback olarak bırak.

4. **`build()` içinde yan etki**: `addPostFrameCallback` veya `setState()` build fazında çağrılmamalı. `ref.listen` build dışı faz garantisi verir.

5. **WorkManager callbackDispatcher**: Mutlaka `@pragma('vm:entry-point')` annotasyonu olmalı ve top-level fonksiyon olmalı.

---

## 🔴 CRITICAL Standards (MUST FOLLOW AUTOMATICALLY)

### ⚠️ GLOBAL RULE: Object-Oriented, Enterprise, Clean Product Engineering

**THE RULE:** Every implementation MUST follow object-oriented/domain-oriented engineering principles and enterprise product quality. AI assistants MUST NOT create one-off procedural patches, duplicate implementations, or hard-coded behavior that should be managed by users/admins.

This rule applies to all repositories:
- `/workspace/go/aystek-core`
- `/workspace/node/aystek-core-next`
- `/workspace/flutter/kisisel_asistan`
- `/workspace/aysteknoloji-com`
- `/workspace/go/aystek-core-installer`
- all related documentation under `/workspace/aystek-core-dev-docs`

**Mandatory engineering behavior:**
- ✅ Place behavior in the correct owner: domain model, service, repository, handler, API client, hook, component, provider, widget, controller, or settings module.
- ✅ Use interfaces, constructors, dependency injection, typed DTOs/models, and domain methods where they fit the existing codebase.
- ✅ Keep modules cohesive and loosely coupled. Do not scatter business rules across unrelated UI handlers, scripts, or utility files.
- ✅ Reuse existing abstractions and patterns before creating new ones.
- ✅ Add a new abstraction only when it removes real duplication, clarifies ownership, or matches an established architecture pattern.
- ✅ Preserve tenant/site/user boundaries, permissions, localization, auditability, and operational safety.
- ✅ Treat every change as production work: testable, observable, maintainable, reversible, and documented when needed.

**Forbidden:**
- ❌ Copy-paste implementations that create duplicate business logic.
- ❌ Hard-coded text, permissions, statuses, menus, feature flags, provider lists, tenant behavior, or operational settings when they belong in seeds/admin/settings/database configuration.
- ❌ UI-only fixes that bypass backend/domain rules.
- ❌ Backend-only defaults that users/admins cannot manage when the feature is meant to be configurable.
- ❌ Temporary compatibility code left behind without an explicit migration/removal plan.

### ⚠️ GLOBAL RULE: Remove Obsolete, Unused, and Replaced Work

**THE RULE:** When a change makes existing code, files, routes, components, hooks, services, seeds, assets, styles, docs, or tests obsolete, the obsolete work MUST be removed in the same task.

AYSTEK applications must never become a code dump. Every AI assistant is responsible for preventing dead code, duplicate modules, unused screens, outdated docs, and abandoned experiments.

**Mandatory cleanup workflow (every task):**
- ✅ Search for old implementations before adding a new one.
- ✅ Prefer extending/refactoring the existing owner over creating parallel code.
- ✅ Two implementations of the same repository/service/handler are a defect: find which one is wired in `cmd/api/main.go`, keep that one, delete the other (and its tests) in the same task. Tests must exercise the real implementation, not a private mock of it.
- ✅ After implementation, search for unused/replaced symbols, files, routes, imports, components, hooks, styles, translations, permissions, menu entries, seeds, and docs.
- ✅ Delete obsolete code and artifacts that are no longer used.
- ✅ Update tests to target the active implementation only.
- ✅ Update documentation to describe the current behavior only.
- ✅ If obsolete data must remain temporarily for migration/backward compatibility, document why, where it is used, and when/how it can be removed.

**Examples of obsolete work that MUST be removed:**
- Replaced React components, Flutter widgets, hooks, providers, services, or API clients
- Old Go handlers, services, repositories, interfaces, models, routes, seeds, indexes, or migrations
- Duplicate CSS/Tailwind classes, theme files, icons, images, and generated assets
- Unused translation keys, menu entries, permissions, settings, feature flags, and mock data
- Outdated documentation that describes behavior no longer present in the product

### ⚠️ GLOBAL RULE: Dynamic and User/Admin-Managed Product Behavior

**THE RULE:** Product behavior that users, admins, tenants, or operators reasonably need to manage MUST be modeled dynamically through AYSTEK Core data/configuration instead of being hard-coded in application code.

**Default approach:**
- ✅ Store configurable domain data in MongoDB through proper models/repositories/services.
- ✅ Seed safe defaults through the appropriate seed files.
- ✅ Expose management through admin pages, settings screens, menus, permissions, translations, or tenant/site configuration where relevant.
- ✅ Keep frontend screens driven by backend data/contracts, not duplicated static lists.
- ✅ Keep Flutter/offline behavior synchronized through explicit local storage/sync contracts instead of hidden hard-coded assumptions.

**Hard-code only when:**
- The value is a true technical constant.
- The value is part of a protocol/library contract.
- The value is intentionally immutable and documented as such.
- Changing it through UI/configuration would be unsafe or meaningless.

### ⚠️ GLOBAL RULE: All UI Text Must Be Translation Keys + Seeded

**THE RULE:** Any user-visible text you add (Go or Next.js) MUST be implemented as a translation **key** and seeded in `translation_seed.go` with **all 28 languages**.  
This applies to: labels, buttons, headings, table columns, tooltips, errors, notices, empty states, and menu titles.

**Mandatory workflow (every task):**
- ✅ Add/verify translation **key(s)** for every new text
- ✅ Update `/workspace/go/aystek-core/internal/seeds/translation_seed.go` with **28 languages**
- ✅ Do not ship plain strings for UI text

### ⚠️ GLOBAL RULE: Event-Driven Modules Must Document Event Keys

**THE RULE:** Any module that publishes domain events MUST document its event keys with the `EventKeysInfo` component on the frontend list page.

**When to apply:**
- Module publishes events via `internal/platform/events`
- Backend service calls `events.Publish()` with domain events
- Integration hooks can subscribe to the module's events

**Required components:**
- ✅ Define `EventKeyDefinition[]` array in frontend list page with all event keys
- ✅ Document event description and all payload fields (name, type, description, required flag)
- ✅ Import and render `<EventKeysInfo event_keys={...} />` component in page
- ✅ Place component between header and main content for user visibility
- ✅ Translation keys already available (no additional translations needed for standard UI)

**Example reference:** `/workspace/node/ays tek-core-next/app/(dashboard)/admin/stock/vehicle-assets/page.tsx`

**Benefits:**
- Users discover integration capabilities without backend code access
- Integration developers understand exact payload structure
- Reduces documentation maintenance and support requests
- Encourages consistent event publishing patterns across modules

### ⚠️ RULE #1: Translation Seed (28 Languages)

**Location:** `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

**THE RULE:** Every new translation key MUST have translations in ALL 28 languages.

**28 Supported Languages:**
```
tr, en, de, fr, it, es, pt, pt-br, az, zh, ja, fa, ar, kk, uz, ru, uk, pl, hr, cs, bs, sr, sq, hy, az-cyrl, ky, tk, mn
```
**Source of truth:** `/workspace/go/aystek-core/internal/seeds/localization_seed.go` (language list and defaults)

**Note:** Some legacy keys in `translation_seed.go` may be partially translated. New keys MUST include all 28 languages.

**Pattern:**
```go
"menu.checkout": {
    "tr": "Ödeme",
    "en": "Checkout",
    "de": "Kasse",
    "fr": "Paiement",
    "it": "Pagamento",
    "es": "Pago",
    "pt": "Pagamento",
    "pt-br": "Pagamento",
    "az": "Ödəmə",
    "zh": "结账",
    "ja": "チェックアウト",
    "fa": "پرداخت",
    "ar": "الدفع",
    "kk": "Төлем",
    "uz": "To'lov",
    "ru": "Оплата",
    "uk": "Оплата",
    "pl": "Płatność",
    "hr": "Plaćanje",
    "cs": "Platba",
    "bs": "Plaćanje",
    "sr": "Плаћање",
    "sq": "Pagesa",
    "hy": "Վճարում",
    "az-cyrl": "Өдəмə",
    "ky": "Төлөм",
    "tk": "Töleg",
    "mn": "Төлбөр"
},
```

**When to Add Translations:**
- ✅ New menu items
- ✅ New status labels
- ✅ New error messages
- ✅ New form labels
- ✅ Any user-facing text

**Enforcement (automatic):**
- `cd /workspace/go/aystek-core && go test ./internal/seeds` runs `TestFullyTranslatedPrefixesHaveNoGaps`: every key under a prefix listed in `fullyTranslatedPrefixes` (`translation_coverage_test.go`) must have all 28 languages, no Turkish text reused for another language and no English text in a non-Latin-script language. **Add the prefix of every new feature to that list** and use an overlay file (`apply<Feature>Translations`, see `site_language_translations_seed.go`) instead of `allLangs(tr, en)`, which fills every language with English.
- `cd /workspace/node/aystek-core-next && npm run check-i18n` fails when the frontend uses a `t("key")` no seed defines and the key is not in the legacy baseline (`scripts/i18n-baseline.json`). The baseline may only shrink.
- Diagnostics: `TRANSLATION_REPORT_PREFIX=<prefix> go test ./internal/seeds -run TestTranslationCoverageReport -v`.
- Open legacy debt and its plan: `I18N_TRANSLATION_DEBT_WORKPLAN_TR.md`.

**Validation Checklist:**
- [ ] Key added to `translations` map
- [ ] All 28 languages have values
- [ ] Turkish (tr) translation is primary
- [ ] English (en) translation is complete
- [ ] Other languages have translations (even if placeholder)

**Example (Adding New Module):**
```go
// NEW MODULE: Invoices
"menu.invoices": {
    "tr": "Faturalar",
    "en": "Invoices",
    "de": "Rechnungen",
    "fr": "Factures",
    "it": "Fatture",
    "es": "Facturas",
    "pt": "Faturas",
    "pt-br": "Faturas",
    "az": "Fakturalar",
    "zh": "发票",
    "ja": "請求書",
    "fa": "فاکتورها",
    "ar": "الفواتير",
    "kk": "Шоттар",
    "uz": "Hisob-fakturalar",
    "ru": "Счета",
    "uk": "Рахунки",
    "pl": "Faktury",
    "hr": "Fakture",
    "cs": "Faktury",
    "bs": "Fakture",
    "sr": "Фактуре",
    "sq": "Faturat",
    "hy": "Հաշիվ-ապրանքագրեր",
    "az-cyrl": "Фактуралар",
    "ky": "Эсептер",
    "tk": "Fakturalar",
    "mn": "Нэхэмжлэхүүд"
},
```

### ⚠️ RULE #1.1: Languages Are Site-Scoped (Multi-Language Standard)

**THE RULE:** Which languages a public site offers is **data managed in the admin panel**, never a hard-coded list. The tenant owns the catalog; each site may narrow it.

**Model (authoritative):** `/workspace/go/aystek-core/internal/modules/languages`
- `languages` collection = tenant catalog (names, native names, tenant-wide `is_active` / `enabled` / `is_default`).
- `site_language_configs` collection = optional per-site override (`tenant_id + site_id`, ordered `codes`, `default_code`).
- **No override → the site uses the tenant-wide active languages** (previous behavior, safe default).
- Override codes must be active tenant languages; a language deactivated tenant-wide silently drops out of the site's list.

**API:** `GET/PUT/DELETE /admin/languages/site` (site = `X-Site-ID`), `GET /public/languages/active` and `/public/languages/default` are site-aware. `GET /public/languages` (paged list used by the admin panel) intentionally stays tenant-wide.

**Admin UI (aystek-core-next):** `/admin/languages` (catalog + "Bu Sitede" column + site banner) and `/admin/languages/site` (per-site selection). Menu key `site-languages` under Settings.

**Mandatory for every public website (`aysteknoloji-com`, `aystekhaber-com`, future sites):**
- ✅ Read the language list from `/public/languages/active`; never hard-code `['tr','en']`-style lists, language switchers, sitemap language loops or `lang === 'en' ? … : …` resolution.
- ✅ Resolve the requested `?lang=` against that list; an unavailable language falls back to the site default (see `aystekhaber-com/lib/api/languages.ts` → `resolveSiteLanguage`).
- ✅ Page copy / portal labels are admin-managed per language in the Pages module. A language with no localization falls back **English → Turkish** (see `aystekhaber-com/lib/api/pages.ts`); do not paste translations into code.
- ✅ Sitemaps list only the languages the site offers.

**Mandatory for backend code:**
- ✅ Per-request language caches are keyed by **tenant + site** (never one process-wide cache).
- ✅ When you read the site override, guard on `contextkeys.GetSiteID(ctx) != ""` first (`RequireSiteID` panics in dev without a site).
- ✅ New per-site configuration follows the same shape: optional override record, absent record = inherit tenant behavior, all repository calls through `tenantfilter.WithTenantAndSite`.

**Adding a language to a site (operator flow):** activate it on `/admin/languages` (tenant catalog) → select it on `/admin/languages/site` → add a localization for it to the site's home page in the Pages module so portal labels are translated (until then the site shows English labels).

### ⚠️ RULE #2: Menu Seed

**Location:** `/workspace/go/aystek-core/internal/seeds/menu_seed.go`

**THE RULE:** Every new module MUST have menu entries added to `menu_seed.go`.
**Also:** Any new **navigation link or admin page** added on the Next.js side MUST have a corresponding menu entry in `menu_seed.go`.

**Menu Structure:**
```go
type menuSeed struct {
    Key      string        // Unique identifier (e.g., "checkout")
    Parent   string        // Parent menu key (e.g., "operations")
    TitleKey string        // Translation key (e.g., "menu.checkout")
    TitleTR  string        // Turkish title (fallback)
    TitleEN  string        // English title (fallback)
    Path     string        // URL path (e.g., "/admin/checkout")
    Icon     string        // Icon name (e.g., "credit-card")
    Order    int           // Display order within parent
    Roles    []string      // Required roles (e.g., []string{"admin"})
    Visible  bool          // Show in menu?
}
```

**Example (Adding Checkout Module):**
```go
seeds := []menuSeed{
    // ... existing menus ...
    
    // NEW: Checkout under Operations
    {
        Key:      "checkout",
        Parent:   "operations",
        TitleKey: "menu.operations.checkout",
        TitleTR:  "Ödeme",
        TitleEN:  "Checkout",
        Path:     "/admin/checkout",
        Icon:     "credit-card",
        Order:    1,
        Roles:    []string{"admin"},
        Visible:  true,
    },
}
```

**Parent Groups (Standard):**
- `dashboard` - Overview + KPI landing
- `forum` - Forum administration (overview, moderation, governance tools)
- `commerce` - E-commerce (products, inventory, orders, payments, pricing, customers, settings)
- `content` - Content hub (posts, pages, comments, localized content, marketing assets)
- `operations` - Operational workflows (checkout, contracts, subscriptions, licenses, stores)
- `settings` - Platform configuration (menu, languages, translations, plugins, mail templates, GraphQL)
- `security` - Access control & audit (users, roles, permissions, logs, notifications)
- `ai-governance` - AI center & governance (jobs, insights, forum AI decisions/policies, chatbot tooling)
- `analytics` - Deployment/UAT/Go-live dashboards

**Key Naming Standard:**
- Use `menu.<domain>.<section>.<item>` for new items (e.g., `menu.commerce.products.items`).
- Legacy keys remain valid (e.g., `menu.forum_overview`) and must not be renamed without a documented migration map.

**Icon Names (Common):**
- `dashboard`, `shopping-cart`, `package`, `credit-card`, `truck`, `users`, `settings`, `bar-chart`, `inbox`, `bell`, `shield`, `percent`, `ticket`, `calculator`

**Validation Checklist:**
- [ ] Key is unique and descriptive
- [ ] Parent exists (or empty for top-level)
- [ ] TitleKey exists in translation_seed.go
- [ ] Path follows /admin/{module} pattern
- [ ] Icon is appropriate and exists
- [ ] Order is set correctly
- [ ] Roles are appropriate
- [ ] Visible is true (unless hidden)

### ⚠️ RULE #2.1: Permission Seed + Admin Auto-Assign

**Location:** `/workspace/go/aystek-core/internal/seeds/permission_seed.go` and `/workspace/go/aystek-core/internal/seeds/role_seed.go`

**Permission Naming Standard:**
- Format: `module:action` (e.g., `forum:read`, `orders:update`)
- Actions must come from the shared seed list (create/read/update/delete/publish/reorder/etc.).

**Admin Auto-Assign Rule (MANDATORY):**
- Admin role **always receives all permissions** at seed time.
- This is enforced in `role_seed.go` by listing all permissions from the DB and assigning them to Admin.
- **No manual admin updates are allowed** when new permissions are added.

**Super-Admin Only Exception (Rare):**
- If a permission must be excluded from Admin, it must be listed in the `superAdminOnly` map in `role_seed.go` and documented.
- **Platform-only modules (in force):** `system` and `system_ops` (updates, rollbacks, backups, migrations, licence, maintenance) are excluded from the tenant Admin role via `platformOnlyPermissionModules` in `role_seed.go`. AYSTEK platform users (`user_type=platform`, `@aystek.com`, or role `platform`) bypass module RBAC in `PermissionMiddleware`, so operators are unaffected; a customer tenant admin gets `403` on `/admin/system/*`. Any new installation-level module must be added to that map, not handed to tenant admins.

**Editor Baseline:**
- Editor role receives **read-only** permissions by default.
- Write permissions are only granted for explicitly allowed content modules in `role_seed.go`.
### ⚠️ RULE #2.2: View Profiles Must Know Every New Menu Key

**Location:** `/workspace/go/aystek-core/internal/seeds/view_profiles_seed.go` (plus `aysfilo_seed.go`, `gursel_turizm_seed.go` for their own menu lists)

View profiles (system → tenant → store → profile) control which **admin menu keys** each site's users see. A new menu entry is not finished until the profiles know it.

**Mandatory for every new `menu_seed.go` entry (and every per-tenant menu list such as `gursel_turizm_seed.go`):**
- ✅ Add the key to `allMenuKeys` (system default = visible everywhere). Least-privilege profiles built with `hiddenMenuItemsExcept(...)` then hide it by default.
- ✅ Decide, per least-privilege profile, whether the site's users need it and add the key to that profile's visible set (`aysTekHaberVisibleMenuKeys`, `gurselTurizmVehicleTrackingVisibleMenuKeys`, `aysfiloVisibleMenuKeys`, `gurselAypassVisibleMenuKeys`, …). Update the allowlist test (`TestAysTekHaberViewProfileAllowlist`) when an editorial menu is required.
- ✅ `filterMenuTree` drops a hidden parent's whole subtree: if you add a child, its parent must be visible in the profiles that should show the child.
- ✅ Re-run `go run ./cmd/seed` locally and confirm the key is not in the profile's hidden list.

### ⚠️ RULE #3: Seed Files for New Models

**Location:** `/workspace/go/aystek-core/internal/seeds/`

**THE RULE:** Every new data model MUST have a corresponding seed file.

**Seed File Naming:** `{domain}_seed.go` (e.g., `checkout_seed.go`, `invoice_seed.go`)

**Seed File Template:**
```go
package seeds

import (
    "context"
    "log"
    "github.com/aystek/aystek-core/internal/modules/checkout/model"
    checkoutrepo "github.com/aystek/aystek-core/internal/modules/checkout/repository"
    "go.mongodb.org/mongo-driver/mongo"
)

// RunCheckoutSeed initializes default checkout data
func RunCheckoutSeed(ctx context.Context, db *mongo.Database) error {
    log.Println("🚀 Seeding checkouts...")
    
    repo := checkoutrepo.NewMongoRepository(db)
    
    // Define seed data
    checkouts := []*model.Checkout{
        {
            // Sample checkout 1
            // ...
        },
        {
            // Sample checkout 2
            // ...
        },
    }
    
    // Insert seed data
    for _, checkout := range checkouts {
        // Check if exists (by unique field)
        existing, _ := repo.FindByEmail(ctx, checkout.Email)
        if existing != nil {
            log.Printf("⏭️  Checkout %s already exists, skipping...", checkout.Email)
            continue
        }
        
        // Create
        if err := repo.Create(ctx, checkout); err != nil {
            log.Printf("❌ Failed to seed checkout %s: %v", checkout.Email, err)
            continue
        }
        log.Printf("✅ Seeded checkout: %s", checkout.Email)
    }
    
    log.Println("✅ Checkout seeding completed!")
    return nil
}
```

**Seed Registration (main.go or seed runner):**
```go
// Run all seeds
func runSeeds(ctx context.Context, db *mongo.Database) {
    seeds.RunUserSeed(ctx, db)
    seeds.RunProductSeed(ctx, db)
    seeds.RunCheckoutSeed(ctx, db)  // ← Add new seed here
    // ... other seeds
}
```

**Validation Checklist:**
- [ ] Seed file created in `/internal/seeds/`
- [ ] Named as `{domain}_seed.go`
- [ ] Exports `Run{Domain}Seed()` function
- [ ] Checks for existing data before inserting
- [ ] Logs progress (🚀, ✅, ❌, ⏭️)
- [ ] Registered in seed runner

---

### ⚠️ GLOBAL RULE: Admin Page Visual Design Standard (aystek-core-next)

**THE RULE:** Every `/admin/*` page in `aystek-core-next` MUST follow the exact visual/structural standard below.

**History note (do not repeat this mistake):** a 2026-07-23 session incorrectly concluded that the gradient hero header was "forbidden" and rewrote `/admin/products/` and drafted rules around a flat, gradient-free header — based on a plausible-looking but wrong reading of which pages were canonical. The product owner corrected this on 2026-07-24: the gradient hero (as seen on `/admin/posts/`, `/admin/posts/create/`, and `/admin/integrations/dashboard/`) **is** the correct, intended standard — `/admin/products/` had been wrongly "flattened" and was restored to its original gradient template. **Do not assume an existing page (or a previous README revision) is correct just because it's old, large, or was "cleaned up" by a past session.** Verify against the concrete markup in this section and against the reference pages below, and if something looks inconsistent, ask rather than guess.

**Verified gold-standard reference pages (read these before writing any admin page):**
- `app/(dashboard)/admin/posts/page.tsx` (list) + `app/(dashboard)/admin/posts/create/page.tsx` (create/edit form) — the canonical gradient-hero header + card-based form body.
- `app/(dashboard)/admin/integrations/dashboard/page.tsx` → `app/(dashboard)/admin/integrations/components/IntegrationsAdminHeader.tsx` — the canonical shared-header-component version of the same pattern (breadcrumb + title + refresh/create/back pill actions), for pages nested under a module with sub-pages.
- `app/(dashboard)/admin/products/` (list+create+`[id]/edit`) — restored to this standard on 2026-07-24, safe to copy directly.

**Required — every admin page must have these:**
- ✅ Gradient hero header: `rounded-3xl bg-gradient-to-r from-{colorA} via-{colorB} to-{colorC} text-white p-6 shadow-xl`. Each module may pick its own 3-stop gradient identity (e.g. products: `from-indigo-600 via-fuchsia-500 to-orange-400`; posts: `from-indigo-600 via-sky-500 to-cyan-400`; integrations: `from-violet-600 via-purple-500 to-fuchsia-400`) — the gradient itself is not fixed to one palette, but every header must be a colored gradient, never a flat/white background.
- ✅ Header action buttons are pills: refresh = `rounded-full bg-white text-slate-800 px-4 py-2 text-sm font-semibold shadow hover:shadow-lg`; primary/create = `rounded-full bg-{module-dark-accent}/70-or-90 text-white px-4 py-2 text-sm font-semibold shadow`; back (on create/edit) = `rounded-full bg-white/20 text-white px-4 py-2 text-sm font-semibold shadow border border-white/30`.
- ✅ Optional stat tiles live **inside** the gradient header: `rounded-2xl bg-white/10 border border-white/20 px-4 py-3` in a `grid gap-3 sm:grid-cols-2 lg:grid-cols-6` under the title row — never as separate white cards below the header.
- ✅ Section/card containers (filters row, form sections, table wrapper) are `rounded-2xl border border-slate-100 bg-white p-5-or-6 shadow-sm` (not `rounded-xl`/`border-slate-200`/no-shadow).
- ✅ `AdminDataTable`/`TableToolbar` (`app/components/`) for tables/search/pagination — never hand-rolled `<table>` markup.
- ✅ `NoticeModal` with `onConfirm`/`confirmLabel`/`confirmLoading` for destructive-action confirmation — never native `window.confirm()`/`alert()`.
- ✅ Form section body content (labels, inputs, in-card buttons) stays understated: labels `text-sm font-medium text-slate-600 mb-1`, inputs `w-full rounded-lg border border-slate-200 px-3 py-2 text-sm` (or `border` shorthand), submit button `rounded-lg bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700` — the gradient/pill treatment is for the **header only**, not for in-card controls.
- ✅ Page width: wide multi-column forms (2-3 field columns per row, like `products/create`) can span the full dashboard width with no wrapper constraint — the grid naturally fills the space. A **simple single-column form** (like `institutions/create`) must NOT be left unconstrained at full width either — that stretches inputs edge-to-edge and, worse, without a wrapper it hugs the left with a huge dead zone on the right instead of looking centered. Wrap single-column forms in `mx-auto max-w-2xl` (or similar) so they read as a centered card, not a left-pinned strip.

**Forbidden:**
- ❌ Flat/white page headers with no gradient — every `/admin/*` page needs the gradient hero.
- ❌ `rounded-xl` + `border-slate-200` + no-shadow as the card standard — that was the incorrect 2026-07-23 convention; use `rounded-2xl` + `border-slate-100` + `shadow-sm`.
- ❌ Hand-rolled `<table>` markup or hand-rolled search/pagination bars when `AdminDataTable`/`TableToolbar` already cover the case.
- ❌ Native `window.confirm()`/`alert()` for destructive actions.

**Mandatory workflow before writing or editing ANY admin page:**
- ✅ Open `posts` (list+create) and/or `integrations/dashboard`+`IntegrationsAdminHeader` and copy their header/card structure — do not improvise a flat header or invent a different card radius/border/shadow combination.
- ✅ If you are fixing an old page that has drifted from this, restore the gradient-hero + `rounded-2xl`/`border-slate-100`/`shadow-sm` pattern in the same change — do not leave a flat header next to gradient siblings, and do not "clean up" a gradient header into a flat one on your own initiative.
- ✅ If you genuinely believe a page (including one of the reference pages above) is wrong, say so and ask — don't silently rewrite the standard based on your own read of "what looks more modern."

---

## ✅ Module Creation Checklist

Use this checklist when creating a new module:

### Backend (Go) - Core 6 Files + Tests + Optional Handlers

#### 1. Model (`model/{domain}.go`)
- [ ] Define domain entity struct
- [ ] Add BSON and JSON tags to all fields
- [ ] Define status/type constants
- [ ] Implement domain methods (validation, state machine, etc.)
- [ ] Include timestamps (CreatedAt, UpdatedAt)
- [ ] Use `primitive.ObjectID` for MongoDB IDs

#### 2. DTO (`dto/{domain}.go`)
- [ ] Define Request DTOs with validation tags
- [ ] Define Response DTOs
- [ ] Define List response with pagination
- [ ] Implement transformation functions (ToXXXResponse)
- [ ] Match backend field names (snake_case for JSON)

#### 3. Repository Interface (`repository/interface.go`)
- [ ] Define Repository interface
- [ ] Include CRUD methods (Create, FindByID, Update, Delete)
- [ ] Add domain-specific queries (FindByCustomerID, etc.)
- [ ] Add List method with pagination

#### 4. Repository Implementation (`repository/mongo.go`)
- [ ] Implement mongoRepository struct
- [ ] Constructor accepts `*mongo.Database`
- [ ] Create indexes (constructor or `CreateIndexes`)
- [ ] Implement all interface methods
- [ ] Set timestamps automatically
- [ ] Return `nil, nil` for not found (not errors)

#### 5. Service (`service/service.go`)
- [ ] Define Service struct with dependencies
- [ ] Define custom error types (ErrXXX)
- [ ] Implement constructor with dependency injection
- [ ] Implement business logic methods
- [ ] Orchestrate multiple repositories if needed
- [ ] Apply business rules and validations

#### 6. Service Tests (`service/service_test.go`)
- [ ] Test all business logic methods
- [ ] Use testify/assert for assertions
- [ ] Mock repository dependencies
- [ ] Test success cases
- [ ] Test error cases
- [ ] Test edge cases (expired sessions, state transitions, etc.)

#### 7. Handler (`handler/public.go` and/or `handler/admin.go`)
- [ ] Implement PublicHandler struct
- [ ] Implement handler methods for each endpoint
- [ ] Bind and validate requests
- [ ] Call service methods
- [ ] Transform to DTOs
- [ ] Return proper HTTP status codes
- [ ] Handle errors gracefully
  - ✅ Use `response` helpers or rely on response envelope middleware

#### 8. Routes (`routes.go`)
- [ ] Export RegisterRoutes function
- [ ] Accept `*gin.RouterGroup` and handlers
- [ ] Register all endpoints with proper HTTP methods
- [ ] Follow RESTful conventions

### Backend Integration

#### 9. Main.go Wiring
- [ ] Import all module packages (model, repo, service, handler)
- [ ] Initialize repository (pass `*mongo.Database`)
- [ ] Initialize service (pass repository dependencies)
- [ ] Initialize handler (pass service)
- [ ] Use `router.New()` and `router.Groups()` for route grouping
- [ ] Call RegisterRoutes with the correct group (`admin`, `public`, or `protected`)

#### 10. Seeds
- [ ] Create `{domain}_seed.go` in `/internal/seeds/`
- [ ] Implement `Run{Domain}Seed()` function
- [ ] Add sample data
- [ ] Register in seed runner

#### 11. Translations
- [ ] Add menu title key to `translation_seed.go`
- [ ] Include ALL 28 languages
- [ ] Add status labels if needed
- [ ] Add error message keys if needed

#### 12. Menu
- [ ] Add menu entry to `menu_seed.go`
- [ ] Set correct parent group
- [ ] Set appropriate icon
- [ ] Set order within parent
- [ ] Set correct roles
- [ ] Add the key to `allMenuKeys` and to the visible set of every least-privilege view profile that needs it (RULE #2.2)

### Frontend (Next.js)

#### 13. Types (`lib/types/{domain}.ts` or `types/{domain}.ts`)
- [ ] Define TypeScript interfaces matching backend DTOs (snake_case)
- [ ] Define status/type unions
- [ ] Export all types

#### 14. API Client (`lib/api/{domain}.ts`)
- [ ] Create typed API methods
- [ ] Use `apiClient` / `axios` helpers and unwrap envelope
- [ ] Normalize IDs and field naming if needed

#### 15. Hooks (`hooks/use{Domain}.ts`)
- [ ] Create custom hook for data fetching (optional)
- [ ] Return loading/error/data states
- [ ] Implement action functions

#### 16. Components (`components/{domain}/`)
- [ ] Create list/detail/form components as needed
- [ ] Use Tailwind CSS
- [ ] Use shared UI components

#### 17. Pages (`app/(dashboard)/admin/{domain}/`)
- [ ] Create list page (`page.tsx`)
- [ ] Create detail page (`[id]/page.tsx`) if needed
- [ ] Add loading and error states

#### 18. Event Keys Documentation (MANDATORY for modules that publish events)
- [ ] Define `EventKeyDefinition` array with all event keys published by module
- [ ] Document each event key with description and payload fields
- [ ] Import and use `EventKeysInfo` component in main list page
- [ ] Place component between header and main content (for visibility)
- [ ] Add translation keys for event documentation (if custom text needed)

**Example Pattern:**
```tsx
import { EventKeysInfo, type EventKeyDefinition } from "@/app/components/EventKeysInfo";

// Define event keys published by this module
const MODULE_EVENT_KEYS: EventKeyDefinition[] = [
  {
    key: "module.entity.created",
    description: "Triggered when entity is created",
    payload_fields: [
      { name: "entity_id", type: "string", description: "Entity ID", required: true },
      { name: "entity_name", type: "string", description: "Entity name", required: true },
      { name: "status", type: "string", description: "Entity status" },
      // ... other fields
    ],
  },
  {
    key: "module.entity.updated",
    description: "Triggered when entity is updated",
    payload_fields: [
      // ... fields
    ],
  },
  // ... other event keys
];

export default function ModuleListPage() {
  return (
    <div className="space-y-5">
      <ModuleHeader title="..." subtitle="..." />
      
      {/* Event keys documentation */}
      <EventKeysInfo event_keys={MODULE_EVENT_KEYS} />
      
      {/* Rest of page content */}
      <div className="grid gap-4">...</div>
    </div>
  );
}
```

**Why This is Important:**
- ✅ Users can discover available event keys without reading backend code
- ✅ Integration developers know exactly what fields are available in payloads
- ✅ Reduces support burden by making event structure self-documenting
- ✅ Encourages consistent event naming and documentation across modules
- ✅ Enables users to create integration hooks directly from the module page

**Translation Keys (Already Added):**
- `module.event_keys.title` - "Event Keys"
- `module.event_keys.description` - Expandable description text
- `module.event_keys.copy_key` - Copy button tooltip
- `module.event_keys.expand` / `collapse` - Toggle buttons
- `module.event_keys.payload_fields` - Section header
- `module.event_keys.integration_hint` - Usage hint text

---

## 🔗 Integration Patterns

### Cross-Module Communication

**Best Practice:** Prefer depending on other module repositories to reduce coupling, but the current codebase DOES use service-to-service dependencies in several modules (e.g., checkout uses carts/orders/taxes services). If you depend on another service:
- Use an interface type (not a concrete struct)
- Avoid circular dependencies
- Keep the dependency direction stable (core → peripheral)

**✅ CORRECT:**
```go
// Checkout service depends on repositories (preferred)
type CheckoutService struct {
    repo       checkoutrepo.Repository
    cartsRepo  cartsrepo.Repository      // ✅ OK - Repository
    ordersRepo ordersrepo.Repository     // ✅ OK - Repository
}
```

**❌ INCORRECT:**
```go
// Checkout service depends on Orders service via interface (allowed when needed)
type CheckoutService struct {
    repo          checkoutrepo.Repository
    ordersService ordersservice.Service   // ✅ OK IF interface + no cycles
}
```

**Reasoning:** This prevents circular dependencies and keeps service layers focused, while allowing real-world orchestration where necessary.

### Event-Driven Communication (Current Standard)

For complex workflows, use the platform event bus and RabbitMQ-backed async processing patterns.

```go
// Event publisher interface
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}

// Example: Checkout publishes OrderCreated event
func (s *Service) ConfirmCheckout(ctx context.Context, checkoutID primitive.ObjectID) error {
    // ... create order ...
    
    // Publish event
    event := Event{
        Type: "OrderCreated",
        Data: map[string]interface{}{
            "orderId": order.ID,
            "customerId": order.CustomerID,
        },
    }
    s.eventPublisher.Publish(ctx, event)
    
    return nil
}
```

### 🧩 Integrations Module Structure (aystek-core) — AI Implementation Guide

This section explains the canonical structure of the `integrations` module so AI assistants can implement new integrations consistently.

#### 1) Module Folder Map (Backend)

Primary path: `/workspace/go/aystek-core/internal/modules/integrations/`

Expected subpackages:
- `model/` → domain entities (`connector`, `action`, `hook`, `execution`, `dlq`)
- `dto/` → admin/public request-response contracts
- `repository/` → MongoDB read/write layer and filtering queries
- `service/` → orchestration, validation, execution triggering
- `handler/` → HTTP endpoints (usually admin-focused)
- `eventbridge/` → event listener + queue publish bridge
- `executor/` → outbound HTTP/SOAP request builder + response parsing
- `worker/` → async consumer processing and retry handling
- `routes.go` → route registration

#### 2) Core Domain Concepts (MUST use this mental model)

- **Connector**: Target system definition (base URL, auth config, status, tenant/site scope)
- **Action**: Concrete operation on connector (method/path/template/timeout/retry)
    - **Action Key Naming Standard**: All `action_key` values MUST use `UPPER_SNAKE_CASE` format (e.g., `GPS_UPDATE`, `CREATE_ORDER`, `UPDATE_INVENTORY`)
    - This ensures consistency and readability across the integration system
- **Hook**: Event subscription rule (`event_key` + conditions + action mapping + mode sync/async)
- **Execution**: Runtime attempt record (request/response/status/error/latency)
- **DLQ Record**: Exhausted-failure message for later replay/analysis

AI assistants should design every new integration change around these 5 objects, not ad-hoc endpoint calls.

#### 3) Runtime Flow (Authoritative Sequence)

1. Domain module publishes platform event (`internal/platform/events`)
2. Integrations eventbridge listener receives event
3. Matching hooks are found by `event_key` + tenant/site scope
4. `sync` hook executes inline OR `async` hook enqueues RabbitMQ job
5. Worker executes action via executor and records execution result
6. Failures apply retry policy; exhausted jobs are persisted to DLQ/execution logs

This is the default integration pipeline in `aystek-core`.

#### 4) Tenant/Site Isolation Rules (MANDATORY)

- Hook matching MUST respect tenant/site boundaries.
- Event payload context (`tenant_id`, `site_id`) must be propagated from publisher to queue worker.
- Do not implement cross-tenant forwarding unless explicitly requested and approved as a separate feature.
- Panel-managed connector/hook configuration is the default control plane; code changes should not be required for routine routing changes.

#### 5) Executor Rules (HTTP + SOAP/XML)

- Request builder must support JSON and XML payload generation from mapping/template definitions.
- Basic auth credentials must come from connector auth config and be encoded correctly.
- Content-Type should be configurable from action/request mapping.
- Timeout, retry, and idempotency behavior must be explicit in action/hook configuration.

#### 6) AI Coding Checklist for Integrations Changes

When adding or modifying integration behavior, AI assistants MUST verify:

- [ ] Event key is stable and documented (publisher + hook side)
- [ ] Hook seed/config includes correct `action_id`/`action_keys`
- [ ] Tenant/site filters are present in repository query path
- [ ] Async flows are RabbitMQ-backed and observable
- [ ] Execution records include enough context for audit/debug
- [ ] Error paths are retry-safe and do not drop failures silently
- [ ] Tests cover success + failure + matching logic

#### 7) Frontend Boundary (aystek-core-next)

- `aystek-core-next` manages connectors/hooks/executions through backend APIs.
- Frontend must not connect to AMQP directly.
- UI should present status/history from API/SSE, while execution happens in backend workers.

**User-Friendly Event Key Selection:**
- Hook form (`HookForm.tsx`) includes autocomplete dropdown for event keys
- `KNOWN_EVENT_KEYS` constant lists available events from various modules
- Users type and see filtered suggestions with module name and description
- Dropdown shows: event key (code), description, and source module badge
- Translation key: `integrations.event_key_hint` for input help text

**When adding new event keys to autocomplete:**
1. Add entry to `KNOWN_EVENT_KEYS` array in `HookForm.tsx`
2. Include: `key`, `module`, `description` fields
3. Group by module for better organization
4. Update this list when new modules publish events

If unsure, default to: **event publish in domain → hook match in integrations → async queue execution in backend**.

### 🐇 RabbitMQ Usage Standard (MANDATORY DECISION GUIDE)

This section defines when AI assistants MUST/SHOULD use RabbitMQ in `aystek-core` and how `aystek-core-next` should consume results.

#### 1) `aystek-core` (Go backend) — Producer/Worker Side

**Current implementation references (authoritative examples):**
- `internal/modules/integrations/eventbridge/listener.go` (domain event listener)
- `internal/modules/integrations/eventbridge/queue.go` (RabbitMQ publish)
- `internal/modules/integrations/worker/` (async execution workers)
- `cmd/api/main.go` (RabbitMQ connection + worker startup)
- `docker-compose.test.yml` (RabbitMQ test service)

**When RabbitMQ SHOULD be used:**
- External API/webhook/integration calls that must not block request-response flow
- Retryable operations with failure isolation (DLQ/retry logic)
- Bursty workloads requiring backpressure and controlled concurrency
- Multi-step async orchestration where eventual consistency is acceptable

**When RabbitMQ SHOULD NOT be used:**
- Simple synchronous CRUD calls that must return immediate business result
- Single-process in-memory notifications where loss is acceptable
- Very low-latency operations that require strict inline transaction guarantees

**Implementation rule for AI assistants:**
1. Publish domain event (`internal/platform/events`) from service layer
2. Match/filter hook/subscription by event key
3. For `async` mode, enqueue to RabbitMQ (durable queue)
4. Process with worker pool (ack/nack, retry, timeout, observability)
5. Persist execution result and error context for auditability

**Operational baseline:**
- Queue must be durable
- Worker count/prefetch must be configurable
- Shutdown must be graceful
- Failures must be traceable (execution logs / DLQ)

#### 2) `aystek-core-next` (Next.js frontend) — Consumer Side

`aystek-core-next` MUST NOT connect directly to RabbitMQ. Frontend integration model is:
- Call backend REST APIs
- Receive near-real-time updates via SSE endpoints (for supported domains)
- Poll backend status endpoints when SSE is not available

**AI assistants must prefer:**
- Backend queue + worker for async jobs
- Frontend only observes job state/result through API/SSE

**Do not implement:**
- Direct AMQP client usage in Next.js
- Browser-level queue credentials exposure
- Queue topology management in frontend code

#### 3) AI Decision Checklist (Use Before Coding)

Before adding a new integration/workflow, answer:
- Is this operation external/slow/failure-prone? → Use RabbitMQ async pipeline
- Does user need immediate final result in same request? → Keep synchronous
- Is retry/idempotency required? → RabbitMQ + worker + idempotency key
- Should frontend only display progress/result? → API/SSE only, no RabbitMQ in frontend

If at least one of the first or third conditions is true, AI assistants should default to RabbitMQ-backed async design in `aystek-core`.

---

## 🧪 Testing Standards

### Unit Tests (Go)

**Location:** `service/service_test.go`

**Pattern:**
```go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, checkout *model.Checkout) error {
    args := m.Called(ctx, checkout)
    return args.Error(0)
}

// Test case
func TestStartCheckout_Success(t *testing.T) {
    mockRepo := new(MockRepository)
    mockCartsRepo := new(MockCartsRepository)
    
    service := New(mockRepo, mockCartsRepo, nil)
    
    // Setup mocks
    mockCartsRepo.On("FindByID", mock.Anything, mock.Anything).Return(&model.Cart{}, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Execute
    checkout, err := service.StartCheckout(context.Background(), customerID, cartID, "test@example.com", "1234567890")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, checkout)
    assert.Equal(t, model.CheckoutStatusCreated, checkout.Status)
    
    // Verify mocks
    mockRepo.AssertExpectations(t)
    mockCartsRepo.AssertExpectations(t)
}

func TestStartCheckout_CartNotFound(t *testing.T) {
    mockRepo := new(MockRepository)
    mockCartsRepo := new(MockCartsRepository)
    
    service := New(mockRepo, mockCartsRepo, nil)
    
    // Setup mocks
    mockCartsRepo.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
    
    // Execute
    checkout, err := service.StartCheckout(context.Background(), customerID, cartID, "test@example.com", "1234567890")
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, checkout)
}
```

**Key Points:**
- ✅ Use testify/assert and testify/mock
- ✅ Mock all dependencies
- ✅ Test success and error cases
- ✅ Verify mock expectations
- ✅ Name tests clearly (Test{Method}_{Scenario})

### Integration Tests (Optional)

For critical flows, add integration tests:

```go
func TestCheckoutFlow_Integration(t *testing.T) {
    // Setup real MongoDB
    client := connectToTestMongoDB()
    db := client.Database("test_db")
    defer cleanupTestDB(db)
    
    // Initialize real repositories
    checkoutRepo := checkoutrepo.NewMongoRepository(db)
    cartsRepo := cartsrepo.NewMongoRepository(db)
    ordersRepo := ordersrepo.NewMongoRepository(db)
    
    // Initialize service
    service := New(checkoutRepo, cartsRepo, ordersRepo)
    
    // Test full workflow
    // 1. Create cart
    // 2. Start checkout
    // 3. Confirm checkout
    // 4. Verify order created
}
```

**Locations in repo:**
- `/workspace/go/aystek-core/internal/tests/integration`
- `/workspace/go/aystek-core/internal/testutil`

### Frontend Tests (Current State)

The frontend repo does **not** include a dedicated testing framework in `package.json`.

**Current expectations:**
- ✅ Run `npm run lint`
- ✅ Keep TypeScript types strict in the files you touch
- ✅ Add tests only if the task explicitly requests it (and include the needed tooling)

---

## 🎨 Code Style Guidelines

### Go Code Style

#### 1. Naming Conventions
- **Packages:** Lowercase, single word (e.g., `checkout`, `orders`)
- **Files:** Snake case (e.g., `checkout.go`, `service_test.go`)
- **Types:** PascalCase (e.g., `CheckoutService`, `OrderStatus`)
- **Functions/Methods:** PascalCase for exported, camelCase for private
- **Variables:** camelCase (e.g., `customerID`, `checkoutRepo`)
- **Constants:** PascalCase or SCREAMING_SNAKE_CASE (e.g., `CheckoutStatusCreated` or `MAX_RETRIES`)

#### 2. Error Handling
```go
// ✅ CORRECT: Named error types
var (
    ErrCheckoutNotFound = errors.New("checkout not found")
    ErrInvalidStatus    = errors.New("invalid status")
)

// ✅ CORRECT: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create checkout: %w", err)
}

// ❌ INCORRECT: Generic errors
if err != nil {
    return errors.New("error occurred")
}
```

#### 3. Context Propagation
```go
// ✅ CORRECT: Always pass context
func (s *Service) GetCheckout(ctx context.Context, id primitive.ObjectID) (*model.Checkout, error) {
    return s.repo.FindByID(ctx, id)
}

// ❌ INCORRECT: Missing context
func (s *Service) GetCheckout(id primitive.ObjectID) (*model.Checkout, error) {
    return s.repo.FindByID(context.Background(), id)
}
```

**Background work (audited — `TestNoContextBackgroundInRepositoryAndServiceLayers`):** a literal `context.Background()` is forbidden in `repository/` and `service/` layers so background work can never silently lose tenant/site identity. Use the helpers in `internal/platform/contextkeys`:

```go
// Fire-and-forget work started from a request: keeps tenant/site/trace values, drops cancellation + deadline
go s.processJob(contextkeys.Detached(ctx), job)

// Scheduled/startup job acting for one tenant (and optionally one site), no request to derive from
runCtx := contextkeys.ForTenantSite(sched.TenantID, sched.SiteID)

// No request and no tenant: index creation in a repository constructor, platform bootstrap hooks
col.Indexes().CreateOne(contextkeys.System(), model)
```

- ✅ Prefer passing the caller's `ctx` down (thread it through helpers instead of creating a new root).
- ✅ Add a timeout on top of the helper when the work must be bounded: `context.WithTimeout(contextkeys.Detached(ctx), 30*time.Second)`.
- ❌ Never call tenant-scoped repository methods with `contextkeys.System()`.

**Repository tenant audit exceptions:** a repository file without tenant scoping is only acceptable when its model genuinely has no tenant field (global catalog, platform lifecycle). It must be a dedicated file (never mixed with tenant data), listed in `repositoryAuditAllowlist` (`internal/platform/tenantfilter/tenantfilter_audit_test.go`) and documented in `user-guide/archive/internal/REPOSITORY_ALLOWLIST_RATIONALE.md`.

#### 4. Dependency Injection
```go
// ✅ CORRECT: Constructor with dependencies
func New(repo repository.Repository, cartsRepo cartsrepo.Repository) *Service {
    return &Service{
        repo:      repo,
        cartsRepo: cartsRepo,
    }
}

// ❌ INCORRECT: Global variables
var repo repository.Repository

func New() *Service {
    return &Service{}
}
```

### TypeScript/React Code Style

#### 1. Naming Conventions
- **Files:** Kebab case for pages (e.g., `checkout-detail.tsx`), PascalCase for components (e.g., `CheckoutCard.tsx`)
- **Components:** PascalCase (e.g., `CheckoutDetail`, `ProductCard`)
- **Functions:** camelCase (e.g., `startCheckout`, `handleSubmit`)
- **Hooks:** camelCase with `use` prefix (e.g., `useCheckout`, `useProducts`)
- **Types/Interfaces:** PascalCase (e.g., `Checkout`, `StartCheckoutRequest`)

#### 2. Component Structure
```typescript
// ✅ CORRECT: Props interface, typed component
interface CheckoutCardProps {
  checkout: Checkout;
  onConfirm: (id: string) => void;
}

export function CheckoutCard({ checkout, onConfirm }: CheckoutCardProps) {
  return (
    <div className="p-4 border rounded">
      <h3>{checkout.email}</h3>
      <button onClick={() => onConfirm(checkout.id)}>Confirm</button>
    </div>
  );
}

// ❌ INCORRECT: Untyped props
export function CheckoutCard({ checkout, onConfirm }) {
  // ...
}
```

#### 3. Async Handling
```typescript
// ✅ CORRECT: Try/catch with error handling
const handleSubmit = async () => {
  try {
    setLoading(true);
    await checkoutService.confirmCheckout(id);
    toast.success('Checkout confirmed!');
  } catch (error: any) {
    toast.error(error.message || 'Failed to confirm');
  } finally {
    setLoading(false);
  }
};

// ❌ INCORRECT: Unhandled promise
const handleSubmit = () => {
  checkoutService.confirmCheckout(id);
};
```

### General Principles (Both)

1. **DRY (Don't Repeat Yourself):** Extract common logic into utilities/helpers
2. **KISS (Keep It Simple, Stupid):** Avoid over-engineering
3. **YAGNI (You Aren't Gonna Need It):** Don't add features before they're needed
4. **SOLID:**
   - **S**ingle Responsibility
   - **O**pen/Closed Principle
   - **L**iskov Substitution
   - **I**nterface Segregation
   - **D**ependency Inversion
5. **Object-Oriented Ownership:** Put state and behavior in the object/module/layer that owns the domain concept. Avoid procedural code paths that bypass models, services, repositories, hooks, providers, widgets, or controllers.
6. **Enterprise Maintainability:** Prefer typed contracts, dependency injection, cohesive modules, clear boundaries, tests, observability, permissions, localization, and documentation over quick fixes.
7. **No Dead Work:** Remove obsolete code, unused files, duplicate implementations, outdated translations, stale menu entries, unused permissions, and outdated docs when a change replaces them.
8. **Dynamic Product Behavior:** Business behavior that should be managed by users/admins/tenants must live in settings, seeds, database records, permissions, menus, translations, or admin screens rather than hard-coded lists.

---

## 📚 Documentation Requirements

### Code Documentation

#### Go (GoDoc)
```go
// Package checkout provides checkout session management.
package checkout

// CheckoutService handles checkout business logic.
type Service struct {
    repo       repository.Repository
    cartsRepo  cartsrepo.Repository
    ordersRepo ordersrepo.Repository
}

// New creates a new checkout service with required dependencies.
func New(repo repository.Repository, cartsRepo cartsrepo.Repository, ordersRepo ordersrepo.Repository) *Service {
    return &Service{
        repo:       repo,
        cartsRepo:  cartsRepo,
        ordersRepo: ordersRepo,
    }
}

// StartCheckout creates a new checkout session for the given customer and cart.
// It validates the cart exists and belongs to the customer before creating the session.
//
// Parameters:
//   - ctx: Request context
//   - customerID: ID of the customer starting checkout
//   - cartID: ID of the cart to checkout
//   - email: Customer email address
//   - phone: Customer phone number
//
// Returns:
//   - *model.Checkout: The created checkout session
//   - error: Error if validation fails or creation fails
func (s *Service) StartCheckout(ctx context.Context, customerID, cartID primitive.ObjectID, email, phone string) (*model.Checkout, error) {
    // Implementation...
}
```

#### TypeScript (JSDoc)
```typescript
/**
 * Checkout service for managing checkout sessions.
 */
export const checkoutService = {
  /**
   * Start a new checkout session.
   * 
   * @param data - Checkout request data
   * @returns Promise resolving to checkout response
   * @throws {Error} If cart is invalid or customer is unauthorized
   */
  async startCheckout(data: StartCheckoutRequest): Promise<CheckoutResponse> {
    // Implementation...
  },
};

/**
 * Custom hook for managing checkout state and actions.
 * 
 * @param checkoutId - Optional checkout ID to fetch
 * @returns Checkout state, loading status, error, and action functions
 * 
 * @example
 * ```tsx
 * const { checkout, loading, startCheckout } = useCheckout();
 * 
 * const handleStart = async () => {
 *   await startCheckout({ customerId, cartId, email, phone });
 * };
 * ```
 */
export function useCheckout(checkoutId?: string) {
  // Implementation...
}
```

### Task Documentation

When completing a task, create two documentation files in `/workspace/aystek-core-dev-docs/done/`:

1. **`TASK_{N}_{MODULE_NAME}.md`** - Complete technical reference (1,000+ lines)
   - Architecture overview
   - All endpoints with examples
   - Database schema
   - State machine diagrams
   - Integration points
   - Code snippets

2. **`TASK_{N}_SUMMARY.md`** - Executive summary (300-500 lines)
   - What was built
   - Key decisions
   - Completion checklist
   - Testing results
   - Next steps

### Installer and Operations Documentation

Installer, deployment, update, rollback, backup, maintenance, recovery, and server setup behavior are part of the product contract. If these behaviors change, the related documentation must be reviewed and updated in the same work item.

At minimum, evaluate and update the relevant files:

1. `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`
2. `/workspace/aystek-core-dev-docs/user-guide/05-operations/OPERATIONS_GUIDE.md`
3. `/workspace/aystek-core-dev-docs/user-guide/05-operations/DEPLOYMENT_SCENARIOS.md`
4. `/workspace/go/aystek-core-installer/README.md`
5. `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`

AI assistants must treat installer documentation maintenance as a standing requirement. Shipping code without aligning these operator-facing documents is considered incomplete work.

---

## 🔄 Standard Operating Procedures

### When Starting a New Task

1. **Read Requirements Carefully**
   - Understand business requirements
   - Identify affected modules
   - Note integration points

2. **Research Existing Code**
   - Search for similar modules
   - Read existing patterns
   - Check for reusable code
   - Identify old/duplicate implementations that may need to be removed

3. **Plan Architecture**
   - Sketch data model
   - Identify dependencies
   - Design API endpoints
   - Decide which behavior must be configurable by users/admins instead of hard-coded

4. **Create Checklist**
   - Use Module Creation Checklist (above)
   - Break into sub-tasks
   - Estimate complexity

### During Implementation

1. **Follow Patterns**
   - Use existing module structure
   - Follow DDD layers
   - Apply SOLID principles
   - Preserve object-oriented ownership and dependency injection
   - Prefer dynamic/admin-managed configuration for product behavior

2. **Write Tests**
   - Test as you code (TDD preferred)
   - Cover success and error cases
   - Mock dependencies

3. **Validate Changes**
   - Run tests after each change
   - Check compilation
   - Verify no regressions

4. **Remove Obsolete Work**
   - Delete replaced code, files, components, widgets, routes, hooks, services, seeds, assets, and docs
   - Remove unused imports, translations, menu entries, permissions, settings, and tests
   - Document any temporary compatibility code with a removal plan

5. **Update Documentation**
   - Add GoDoc/JSDoc comments
   - Update README if needed
   - Note breaking changes

### Before Completing

1. **Final Validation**
   - [ ] All tests passing
   - [ ] Project compiles
   - [ ] No linter errors
   - [ ] Integration with main.go complete
   - [ ] No duplicate/replaced implementation remains active

2. **Seed Data**
   - [ ] translation_seed.go updated (28 languages)
   - [ ] menu_seed.go updated
   - [ ] Domain seed file created
   - [ ] All new UI text uses translation keys (no hard-coded strings)
   - [ ] All new admin pages/links have menu entries
   - [ ] New menu keys are registered in view profiles (RULE #2.2)
   - [ ] Public-site language behavior comes from site-scoped languages, not hard-coded lists (RULE #1.1)
   - [ ] Dynamic/admin-managed settings, permissions, menus, and defaults are seeded where needed

3. **Cleanup**
   - [ ] Obsolete files/functions/routes/components/hooks/widgets/styles/assets removed
   - [ ] Unused translations/menu entries/permissions/settings removed or migrated
   - [ ] Temporary compatibility code has a documented removal plan
   - [ ] Documentation reflects only the active implementation

4. **Documentation**
   - [ ] Code comments complete
   - [ ] Task documentation created
   - [ ] Summary document created
   - [ ] Files moved to `/done/` folder

5. **Communication**
   - [ ] Notify about completion
   - [ ] Highlight key decisions
   - [ ] Note any blockers or issues

---

## 🚨 Common Pitfalls to Avoid

### Backend (Go)

1. **❌ Passing wrong type to repository constructor**
   ```go
   // WRONG
   repo := checkoutrepo.NewMongoRepository(collection)
   
   // CORRECT
   repo := checkoutrepo.NewMongoRepository(db)
   ```

2. **❌ Not handling MongoDB not-found correctly**
   ```go
   // WRONG
   if err != nil {
       return nil, err  // ErrNoDocuments is treated as error
   }
   
   // CORRECT
   if err == mongo.ErrNoDocuments {
       return nil, nil  // Not found is not an error
   }
   ```

3. **❌ Forgetting to set timestamps**
   ```go
   // WRONG
   checkout := &model.Checkout{
       Email: "test@example.com",
   }
   
   // CORRECT
   checkout := &model.Checkout{
       Email:     "test@example.com",
       CreatedAt: time.Now(),
       UpdatedAt: time.Now(),
   }
   ```

4. **❌ Not validating state transitions**
   ```go
   // WRONG
   checkout.Status = model.CheckoutStatusCompleted
   
   // CORRECT
   if !checkout.CanTransition(model.CheckoutStatusCompleted) {
       return ErrInvalidStatus
   }
   checkout.Status = model.CheckoutStatusCompleted
   ```

### Frontend (TypeScript)

1. **❌ Not handling loading states**
   ```tsx
   // WRONG
   const { checkout } = useCheckout(id);
   return <div>{checkout.email}</div>;  // Crashes if checkout is null
   
   // CORRECT
   const { checkout, loading } = useCheckout(id);
   if (loading) return <Spinner />;
   if (!checkout) return <div>Not found</div>;
   return <div>{checkout.email}</div>;
   ```

2. **❌ Not typing API responses**
   ```typescript
   // WRONG
   const response = await api.get('/checkout/123');
   
   // CORRECT
   const response = await api.get<{ checkout: Checkout }>('/checkout/123');
   const checkout = response.data.checkout;
   ```

3. **❌ Not handling async errors**
   ```typescript
   // WRONG
   const handleClick = async () => {
       await checkoutService.confirm(id);
   };
   
   // CORRECT
   const handleClick = async () => {
       try {
           await checkoutService.confirm(id);
           toast.success('Success!');
       } catch (error) {
           toast.error('Failed!');
       }
   };
   ```

---

## 📖 Quick Reference

### Go Module File Structure
```
modules/{module}/
├── model/{module}.go          # Domain models
├── dto/{module}.go            # Request/Response DTOs
├── repository/
│   ├── interface.go           # Repository contract
│   └── mongo.go               # MongoDB implementation
├── service/
│   ├── service.go             # Business logic
│   └── service_test.go        # Unit tests
├── handler/
│   ├── public.go              # Public endpoints
│   └── admin.go               # Admin endpoints (optional)
└── routes.go                  # Route registration
```

### TypeScript Module File Structure
```
lib/
├── api/{module}.ts            # API client for module
├── types/{module}.ts          # Module types (when needed)
hooks/
├── use{Module}.ts             # Custom React hook (optional)
components/
├── {module}/                  # Feature components
│   ├── {Module}Card.tsx
│   ├── {Module}Form.tsx
│   └── {Module}List.tsx
app/(dashboard)/admin/{module}/
├── page.tsx                   # List page
└── [id]/page.tsx              # Detail page (optional)
```

### 🔑 Production Server Access (öğrenilenler — 2026-10-23)

AI asistanların bu sandbox'tan prod sunucusuna erişirken tekrar tekrar düştüğü
tuzaklar ve doğru yol:

1. **Hangi host'a bağlan**: `core.aysteknolojileri.com:22` (public SSH) bu
   sandbox'tan HER ZAMAN engelli (platform ağ politikası, host+port bazlı
   allowlist — VPN açıp kapatmak bunu etkilemez). Sunucunun bir İÇ IP'si var
   (konumu: `/workspace/COZUMTEK_SUNUCU_BILGILERI.txt`, bu dosya git-tracked
   DEĞİL) — o iç IP:22 bu sandbox'tan ERİŞİLEBİLİR.
2. **Şifreyi komut satırına yazma**: `sshpass`/env-var ile şifre geçiren bir
   SSH komutu, Claude Code'un auto-mode sınıflandırıcısı tarafından
   ("canlı bir credential'ı komutta göstermek") engellenir — bu ayarla
   aşılamayan sabit bir kural. Doğru yol: `~/.ssh/id_ed25519.pub`'ı sunucuda
   hedef kullanıcının `~/.ssh/authorized_keys`'ine ekletip anahtar tabanlı,
   şifresiz `ssh user@host "..."` kullan. Bu hem sınıflandırıcıya takılmaz
   hem de `~/.claude/settings.json`'daki `permissions.allow` içinde
   `"Bash(ssh user@host:*)"` gibi dar bir kuralla tam eşleşir (komut
   `sshpass` ile değil doğrudan `ssh` ile başlamalı — pattern eşleşmesi
   tam komut önekine bakıyor).
3. **Sunucu kaynak kodu nerede**: release dizinleri (`/opt/aystek-releases/
   <tarih>-<hash>/aystek-core`) otomatik deploy geçmişi — buraya DOKUNMA.
   Asıl git repo (buradan `git pull` yapılır) ayrı bir yerde, host'un
   `/home/<kullanıcı>/aystek-core` gibi bir dizininde.
4. **Go toolchain sunucuda yok, Docker var**: prod container'ları (`docker ps`
   ile görünür — `aystek-core`, `aystek-mongodb`, `aystek-redis`,
   `aystek-rabbitmq`, panel ve her public site'ın kendi container'ı) sadece
   derlenmiş binary çalıştırır, Go SDK içermez. Tek seferlik bir Go komutu
   (örn. bir seed migration) çalıştırmak için GEÇİCİ bir `golang:<go.mod'daki
   sürüm>` container'ı aç, kaynak dizini mount et, ve
   `--network container:aystek-core` ile ÇALIŞAN aystek-core container'ının
   network namespace'ini paylaş — böylece `mongodb`/`redis`/`rabbitmq` gibi
   container hostname'lerine, prod'daki gibi erişebilir:
   ```bash
   docker run --rm -v /path/to/aystek-core:/app -w /app \
     --network container:aystek-core --env-file <env-dosyası> \
     golang:1.26 go run ./cmd/seed --pa-news-per-site
   ```
5. **Doğru ortam değişkenlerini nereden al**: repodaki `.env`/`.env.production`
   dosyaları çalışan container'ın GERÇEK ortamıyla eşleşmeyebilir (eski/farklı
   olabilir). Garantili yol: `docker exec aystek-core env > /tmp/x.env` ile
   ÇALIŞAN container'ın kendi ortamını kopyala, SONRA `grep -v '^PATH=' /tmp/x.env
   > /tmp/x-nopath.env` ile PATH satırını çıkar (yoksa `golang` image'ının
   kendi `go` binary'sinin PATH'ini ezer → "go: command not found").
6. **Yeni bir seed fonksiyonu `seedNewsStore`/`upsertSiteSeedDefaults` gibi
   tenant-scoped yardımcıları çağırıyorsa**: çağrıdan ÖNCE
   `contextkeys.WithSiteID(contextkeys.WithTenantID(ctx, tenantID), "")` ile
   context'i hazırlamayı UNUTMA — `SeedNewsSite`'ın zaten yaptığı gibi. Aksi
   halde `contextkeys.GetTenantID(ctx)` boş döner ve "requires tenant
   context" hatasıyla sessizce skip edilir (prod'da canlı test edilirken
   bulundu — internal/seeds/pa_news_sources_seed.go, pa_agency_integrations_seed.go).
7. **SSH anahtarı sandbox sıfırlanınca/yeni bir sandbox'ta SIFIRLANIR**:
   `~/.ssh/id_ed25519.pub` her sandbox örneğinde YENİDEN üretilir, sunucudaki
   `authorized_keys`'de ESKİ bir sandbox'ın anahtarı kalmış olabilir ama YENİ
   sandbox'ınki orada YOKTUR. `ssh root@192.168.200.250 "echo OK"` ile ÖNCE
   test et; "Permission denied (publickey,password)" dönerse kullanıcıdan
   `cat ~/.ssh/id_ed25519.pub` çıktını sunucuda
   `echo "<satır>" >> ~/.ssh/authorized_keys` ile EKLEMESİNİ (üzerine
   yazmadan, `>>` ile) iste — kullanıcı dosyanın mevcut içeriğini
   `cat ~/.ssh/authorized_keys` ile gösterirse bunun ESKİ bir sandbox'ın
   anahtarı olup olmadığını (seninkiyle birebir aynı olup olmadığını)
   kontrol et.
8. **Bağlantı aniden `Connection refused`/`Connection timed out` verebilir**:
   iç IP'ye erişim kullanıcının VPN durumuna bağlı — oturum ortasında VPN
   kapanırsa SSH düşer. Bu durumda kullanıcıdan VPN'i tekrar açmasını iste,
   `echo OK` ile yeniden test et, sonra kaldığın yerden devam et (installer
   komutları idempotent, yarıda kesilen bir `update` tekrar çalıştırılabilir).

### 🚀 Backend/Panel'i Canlıya Deploy Etme (öğrenilenler — 2026-11-08)

`aystek-core` ve `aystek-core-next`'teki bir değişikliği production'a
yansıtmanın KANITLANMIŞ, uçtan uca yolu — SSH erişimi yukarıdaki bölümdeki
gibi kurulduktan SONRA:

1. **İki reponun GERÇEK git checkout'unu güncelle** (host'ta, `/opt/aystek-
   releases/...` DEĞİL — madde 3'teki uyarı):
   ```bash
   ssh root@192.168.200.250 "cd /home/aysadm/aystek-core && git pull origin main"
   ssh root@192.168.200.250 "cd /home/aysadm/aystek-core-next && git pull origin main"
   ```
2. **Deploy mekanizması `aystek-core-installer`'dır** (`/usr/local/bin/
   aystek-core-installer`), konfigürasyonu `/etc/aystek-core/installer.json`.
   Bu dosyada `source_mode: "local-source"` ve `core_source_path`/
   `panel_source_path` alanları, installer'ın docker image'ı HANGİ KAYNAK
   DİZİNİNDEN build edeceğini belirler (gereksinim: o dizinde docker modu
   için `deployments/Dockerfile`/`Dockerfile` bulunması — bkz. installer'ın
   kendi `docs/KULLANIM_KILAVUZU.md` §"Local-source prerequisite kuralları").
   **ÖNEMLİ KEŞİF**: Bu alanların `/opt/aystek-releases/<tarih>-<hash>/...`
   gibi DATED bir kopyaya işaret etmesi ZORUNLU DEĞİL — doğrudan madde 1'de
   güncellediğin GERÇEK git checkout'una (`/home/aysadm/aystek-core`,
   `/home/aysadm/aystek-core-next`) işaret edebilir, installer oradan
   build eder. (Bir önceki operatör/AI oturumunda `/opt/aystek-releases/...`
   tamamen silinmiş/temizlenmiş bulundu — installer.json o zaman hâlâ ESKİ,
   var olmayan bir yola işaret ediyordu ve `update` "core_source_path
   invalid: stat ... no such file or directory" ile patlıyordu; çözüm,
   alanı doğrudan git checkout'a yönlendirmekti, YENİ bir dated kopya
   oluşturmaya gerek YOK.)
3. **Deploy'dan ÖNCE `installer.json`'ı yedekle ve commit hash'lerini
   güncelle** (image tag'leri + `release_version`, sadece izlenebilirlik
   için — installer bunları gerçekte kaynak dizininden build ettiği için
   build'in kendisini etkilemiyor ama raporlamada/validate'te doğru
   görünmesi için güncellemek iyi pratik):
   ```bash
   ssh root@192.168.200.250 "cp /etc/aystek-core/installer.json /etc/aystek-core/installer.json.bak-\$(date +%Y%m%d%H%M%S)"
   ssh root@192.168.200.250 'python3 - <<"PYEOF"
   import json
   p = "/etc/aystek-core/installer.json"
   with open(p) as f: cfg = json.load(f)
   cfg["core_source_path"] = "/home/aysadm/aystek-core"
   cfg["panel_source_path"] = "/home/aysadm/aystek-core-next"
   cfg["core_image_tag"] = "aystek/aystek-core:git-<core-kisa-hash>"
   cfg["panel_image_tag"] = "aystek/aystek-panel:git-<panel-kisa-hash>"
   cfg["release_version"] = "git-<core-kisa-hash>"
   with open(p, "w") as f: json.dump(cfg, f, indent=2, ensure_ascii=False)
   print("UPDATED_OK")
   PYEOF'
   ```
   (Bu JSON güncellemesi zararsız, idempotent bir yapılandırma değişikliği
   — bir dosyadaki iki path alanını ve üç sürüm etiketini günceller,
   production'da kalıcı/geri alınamaz bir etki yaratmaz.)
4. **Önce `--dry-run` ile doğrula, sonra gerçek `update`i çalıştır**:
   ```bash
   ssh root@192.168.200.250 "/usr/local/bin/aystek-core-installer update --config /etc/aystek-core/installer.json --non-interactive --dry-run -v"
   ssh root@192.168.200.250 "/usr/local/bin/aystek-core-installer update --config /etc/aystek-core/installer.json --non-interactive -v"
   ```
   Gerçek `update` docker image'ı YENİDEN BUILD ettiği için 45 saniye ile
   150+ saniye arası sürebilir (build cache durumuna göre değişken) —
   komut çoğu zaman 600 saniyelik ön plan zaman aşımını aşıp arka plana
   taşınır, bu NORMALDİR; arka plan görev bildirimini bekle. `update`
   ADIMLARI idempotent'tir — bağlantı VPN nedeniyle koparsa (madde 8),
   bağlantı dönünce AYNI komutu tekrar çalıştırmak güvenlidir.
5. **`update` başarıyla biterse Validation Report'ta şunları gör**:
   `api_health`/`api_liveness`/`api_readiness`/`mongo_tcp`/`redis_tcp`/
   `rabbit_tcp`/`panel_tcp` ✓ olmalı. `nginx_http` kontrolü SIK SIK
   "context deadline exceeded" ile ⚠ (NON-CRITICAL) döner — bu sandbox'tan
   `https://core.aysteknolojileri.com` DNS/firewall nedeniyle bazen
   ulaşılamaz olabilir, deploy'un kendisinin başarısız olduğu anlamına
   GELMEZ; madde 6'daki curl yöntemiyle ayrıca doğrula.
6. **Deploy'un GERÇEKTEN yeni kodu çalıştırdığını, sadece bir proxy
   catch-all'dan "başarılı" görünmediğini DOĞRULA** — bu projede public
   internetten `https://core.aysteknolojileri.com/pa/<KAYITLI-OLMAYAN-HERHANGİ-BİR-YOL>`
   bile `200 {"success":true}` döner (genel bir reverse-proxy/nginx
   catch-all'ı), yani "200 başarılı" tek başına kanıt DEĞİLDİR. Doğru
   doğrulama: YENİ eklediğin route'u, zaten var olduğunu bildiğin GERÇEK,
   auth gerektiren bir route ile (ör. `/pa/web-search-settings`, 401 döner)
   VE kasıtlı UYDURMA bir yolla (`/pa/bu-yol-asla-var-olmayacak-xyz123`,
   her zaman `{"success":true}` döner) KARŞILAŞTIRARAK yap:
   ```bash
   curl -s -o /dev/null -w "%{http_code}\n" https://core.aysteknolojileri.com/pa/YENI-ROTAN     # deploy öncesi: 200 (catch-all) → deploy sonrası: 401 (gerçek auth kontrolü) bekle
   curl -s https://core.aysteknolojileri.com/pa/bu-yol-asla-var-olmayacak-xyz123                  # HER ZAMAN {"success":true} — kontrol grubu
   ```
   Deploy öncesi/sonrası bu ikisinin farklılaştığını görmeden "deploy
   edildi, çalışıyor" deme.
7. **`/opt/aystek-docker/docker-compose.yml`**: installer tarafından
   OTOMATİK ÜRETİLİR ("Generated by aystek-core-installer... DO NOT EDIT
   MANUALLY" başlığını taşır) — elle değiştirme, `installer.json` üzerinden
   yönet.
8. **Sadece admin panel (ya da sadece backend) değiştiyse**, sadece ilgili
   `*_image_tag`'i güncelleyip AYNI `update` komutunu çalıştırmak yeterli
   — ikisini birden deploy etmek zorunda değilsin, installer her ikisini de
   her `update`de yeniden build/restart eder ama bu genelde sorun
   yaratmaz (idempotent).

### Command Quick Reference

```bash
# Backend (Go)
cd /workspace/go/aystek-core
go test ./internal/modules/checkout/service/...  # Run tests
go build ./cmd/api                                # Build project
go run ./cmd/api/main.go                          # Run server

# Frontend (Next.js)
cd /workspace/node/aystek-core-next
npm run dev                                       # Start dev server
npm run build                                     # Build for production
npm run lint                                      # Lint
npm run check-i18n                                # fails on NEW t("key") usages missing from the backend seeds (baseline: scripts/i18n-baseline.json)
```

---

## 🎯 Final Note

**This document is your SINGLE SOURCE OF TRUTH.**

When in doubt:
1. ✅ Re-read this document
2. ✅ Search for similar existing code
3. ✅ Ask for clarification if business logic is unclear

**Never:**
- ❌ Deviate from the core module structure (model + dto + repository + service + handler + routes + tests)
- ❌ Skip the 28 language translations
- ❌ Forget to add menu entries
- ❌ Skip creating seed files
- ❌ Implement without tests

**Remember:** You are a senior engineer. Make good decisions. Write clean code. Test thoroughly. Document well.

---

---

## 📣 Final Message to AI Assistants

**Congratulations on reading this entire document!** 🎉

You now have all the knowledge needed to work effectively on AYSTEK projects. Here's what you should remember:

### The 3 Golden Rules (Never Forget)

1. **🔴 28 LANGUAGES** - Every translation key needs ALL 28 languages
2. **🔴 MENU ENTRIES** - Every module needs a menu entry in menu_seed.go
3. **🔴 SEED FILES** - Every model needs a seed file

### Your First Task Checklist

When you start working on a new task:
- [ ] I have read this document completely
- [ ] I understand the DDD architecture (Model → DTO → Repository → Service → Handler → Routes)
- [ ] I know the 3 Golden Rules
- [ ] I have the Module Creation Checklist (17 steps)
- [ ] I will create tests for everything I build
- [ ] I will update translations, menus, and seeds automatically

### Questions to Ask Yourself Before Submitting Code

- ✅ Did I create all core files (model, dto, repository, service, handler, routes) and tests?
- ✅ Did I add 28 language translations?
- ✅ Did I add the menu entry?
- ✅ Did I create the seed file?
- ✅ Did I write tests?
- ✅ Did I wire it in main.go?
- ✅ Does the project compile?
- ✅ Do all tests pass?

If you answered "NO" to any of these, **DO NOT submit the code.** Go back and complete it.

### Remember

> **You are a senior software engineer, not a code monkey.**  
> Think critically. Make good decisions. Write clean code. Test thoroughly. Document well.

---

**END OF AI ASSISTANT README**

*Version 1.2.0 - September 21, 2026*

---

**Document Stats:**
- Lines: 2,880+
- Words: 12,100+
- Code Examples: 50+
- Sections: 14
- Checklists: 60+ items
- Languages Supported: 28
- Estimated Read Time: 35-45 minutes

**Pro Tip for AI Assistants:**  
Bookmark this document's key sections. You'll refer to them constantly:
- CRITICAL Standards (global rules + Rules #1, #2, #3)
- Module Creation Checklist (17 steps)
- Common Pitfalls
- Quick Reference
