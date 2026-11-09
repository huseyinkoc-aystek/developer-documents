#!/bin/bash
# OpenAPI Auto-Generation Setup
# Purpose: Configure automatic OpenAPI/Swagger documentation generation
# Date: February 2, 2026

set -e

WORKSPACE_ROOT="$(pwd)"
API_PKG_PATH="./cmd/api"
DOCS_OUTPUT="${WORKSPACE_ROOT}/docs/openapi.yaml"

echo "🔧 Setting up OpenAPI auto-generation..."

# Step 1: Install swag CLI if not present
echo "Step 1: Checking swag CLI installation..."
if ! command -v swag &> /dev/null; then
    echo "  Installing swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
    echo "  ✅ swag installed"
else
    echo "  ✅ swag already installed: $(swag -version)"
fi

# Step 2: Add swag annotations to main.go if not present
echo "Step 2: Checking for swag annotations in main.go..."
if ! grep -q "@title" "${API_PKG_PATH}/main.go"; then
    echo "  Adding swag package comment to main.go..."
    cat > /tmp/swag_header.txt << 'EOF'
// @title AYSTEK E-Commerce Platform API
// @version 2.0.0
// @description Enterprise e-commerce platform with multi-vendor support, marketplace integrations, and ERP sync
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.aystek.com/support
// @contact.email support@aystek.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host api.aystek.com
// @schemes https
// @produce json
// @consume json
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @x-google-audiences internal
EOF
    
    # Insert after package declaration
    sed -i '/^package main/a\\' "${API_PKG_PATH}/main.go"
    sed -i '/^package main/r /tmp/swag_header.txt' "${API_PKG_PATH}/main.go"
    echo "  ✅ Swag annotations added"
else
    echo "  ✅ Swag annotations already present"
fi

# Step 3: Generate OpenAPI documentation
echo "Step 3: Generating OpenAPI documentation..."
cd "${API_PKG_PATH}"
swag init -g main.go --output ../../docs --parseVendor
cd - > /dev/null

if [ -f "${DOCS_OUTPUT}" ]; then
    echo "  ✅ OpenAPI spec generated: ${DOCS_OUTPUT}"
    echo "  📊 Spec size: $(wc -l < ${DOCS_OUTPUT}) lines"
else
    echo "  ⚠️  OpenAPI spec not found at expected path"
fi

# Step 4: Setup GitHub Actions workflow for auto-generation
echo "Step 4: Creating CI/CD workflow for OpenAPI generation..."
mkdir -p .github/workflows

cat > .github/workflows/openapi-generate.yml << 'EOF'
name: OpenAPI Auto-Generation

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'cmd/api/**'
      - 'internal/modules/**/handler/**'
  pull_request:
    branches: [ main, develop ]
    paths:
      - 'cmd/api/**'
      - 'internal/modules/**/handler/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install swag
        run: go install github.com/swaggo/swag/cmd/swag@latest
      
      - name: Generate OpenAPI
        run: |
          cd cmd/api
          swag init -g main.go --output ../../docs --parseVendor
      
      - name: Validate OpenAPI
        run: |
          if [ -f docs/openapi.yaml ]; then
            echo "✅ OpenAPI spec generated successfully"
            wc -l docs/openapi.yaml
          else
            echo "❌ OpenAPI spec generation failed"
            exit 1
          fi
      
      - name: Upload OpenAPI artifact
        uses: actions/upload-artifact@v3
        with:
          name: openapi-spec
          path: docs/openapi.yaml
      
      - name: Comment PR with spec URL
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '📚 OpenAPI spec generated and available in artifacts'
            })
EOF

echo "  ✅ Workflow created: .github/workflows/openapi-generate.yml"

# Step 5: Create handler annotation helpers
echo "Step 5: Creating handler annotation documentation..."
cat > docs/OPENAPI_ANNOTATIONS.md << 'EOF'
# OpenAPI Annotations Guide

## How to Document Endpoints with Swag

### Basic Endpoint Example

```go
// GetOrder retrieves an order by ID
// @Summary Get order details
// @Description Get complete order information including items and payment status
// @ID get-order
// @Accept  json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order "Order details"
// @Failure 404 {object} error "Order not found"
// @Router /admin/orders/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) GetOrder(c *gin.Context) {
    // implementation
}
```

### Parameters

```go
// Query Parameter
// @Param limit query int false "Page size" default(10) maximum(100)
// @Param offset query int false "Page offset" default(0)
// @Param status query string false "Filter by status" Enums(pending,completed,cancelled)

// Path Parameter
// @Param id path string true "Resource ID"

// Form/Body Parameter
// @Param body body CreateOrderRequest true "Order creation request"
```

### Responses

```go
// Success Response
// @Success 200 {object} model.Order "Order created successfully"
// @Success 201 {object} model.Order "Resource created"

// Array Response
// @Success 200 {array} model.Order "List of orders"

// Paginated Response
// @Success 200 {object} dto.PaginatedResponse{data=[]model.Order} "Paginated orders"

// Error Response
// @Failure 400 {object} error "Validation error"
// @Failure 401 {object} error "Unauthorized"
// @Failure 500 {object} error "Server error"
```

### Security

```go
// @Security ApiKeyAuth
// @Security OAuth2Implicit[read, write]

// Multiple security schemes
// @Security ApiKeyAuth || OAuth2Implicit
```

### Generated Output

Run `swag init` in `cmd/api/` to generate:
- `docs/swagger.json` - Machine readable spec
- `docs/swagger.yaml` - Human readable spec
- Routes available at `/docs/swagger-ui` (if hosted)
EOF

echo "  ✅ Annotation documentation created"

# Step 6: Setup Swagger UI (optional)
echo "Step 6: Setting up Swagger UI endpoint..."
cat > docs/README.md << 'EOF'
# API Documentation

## OpenAPI Specification

- **Format:** OpenAPI 3.0.0
- **Location:** `./openapi.yaml`
- **Generated:** Automatically via `swag init`
- **Updated on:** Every push to main/develop branches

## Viewing Documentation

### Option 1: Online Editor
- Navigate to: https://editor.swagger.io
- Import: `raw` OpenAPI YAML from this repository

### Option 2: Local Swagger UI
```bash
docker run -p 8888:8080 \
  -e SWAGGER_JSON=/app/openapi.yaml \
  -v $(pwd)/docs:/app \
  swaggerapi/swagger-ui
# Open: http://localhost:8888
```

### Option 3: ReDoc (Alternative UI)
```bash
docker run -p 8888:8080 \
  -v $(pwd)/docs:/usr/share/nginx/html/data \
  -e SPEC_URL=/data/openapi.yaml \
  redocly/redoc
# Open: http://localhost:8888
```

## Adding New Endpoints

1. Add swag annotations to handler function
2. Run: `cd cmd/api && swag init -g main.go --output ../../docs`
3. Commit updated `openapi.yaml`
4. CI/CD automatically validates spec

## Common Annotations

See `OPENAPI_ANNOTATIONS.md` for detailed examples.

## Spec Validation

```bash
# Validate spec syntax
docker run -it --rm \
  -v $(pwd)/docs:/data \
  redocly/cli lint /data/openapi.yaml
```
EOF

echo "  ✅ Documentation created: docs/README.md"

# Step 7: Make script executable
chmod +x "${WORKSPACE_ROOT}/docs/setup-openapi.sh" 2>/dev/null || true

echo ""
echo "✅ OpenAPI Auto-Generation Setup Complete!"
echo ""
echo "📋 Next Steps:"
echo "  1. Review: docs/OPENAPI_ANNOTATIONS.md"
echo "  2. Generate: cd cmd/api && swag init -g main.go --output ../../docs"
echo "  3. Commit: git add docs/openapi.yaml"
echo "  4. CI/CD: Automatic generation on push"
echo ""
echo "📚 View documentation at: https://editor.swagger.io (import openapi.yaml)"
echo ""
