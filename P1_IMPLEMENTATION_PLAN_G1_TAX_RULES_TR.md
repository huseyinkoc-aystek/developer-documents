# P1 İmplementasyon Planı — Post-Launch Özellikler

**Başlangıç Tarihi:** 2024-02-03 (P0 tamamlanması sonrası)  
**Hedef:** 4 hafta içinde P1 sprint'ini tamamlamak  
**Metodoloji:** Kıdemli backend mühendisi — test-driven, production-ready  

---

## Executive Summary

| Modül | Şiddeti | Durum | Efor | Risk |
|-------|---------|-------|------|------|
| **P1-G1: Tax Rules MVP** | HIGH | ⚠️ PARTIAL | 20h | 🟡 MEDIUM |
| **P1-G2: B2B Workflows** | MEDIUM | ❌ NOT READY | 24h | 🔴 HIGH |
| **P1-G3: Stock Allocation** | MEDIUM | ⚠️ PARTIAL | 12h | 🟡 MEDIUM |
| **P1-G4: Returns/RMA** | MEDIUM | ⚠️ PARTIAL | 8h | 🟡 MEDIUM |
| **Total** | — | — | **64h** | — |

**Sprint Length:** 2 hafta (5 gün/hafta, 8 saat/gün = 80 saat available)  
**Buffer:** 16 saat (testing, documentation, reviews)

---

# PART I: P1-G1 — TAX RULES MVP

## 1. Current State Analysis

### What's Already Implemented

**Model:** ✅ [internal/modules/taxes/model/tax_rule.go](../go/aystek-core/internal/modules/taxes/model/tax_rule.go)

```go
✅ TaxRule struct with:
   - Country + Region + Postal code support
   - TaxType (VAT, GST, Sales)
   - TaxScope (Order, Shipping, Product)
   - TaxCustomerType (Any, B2B, B2C)
   - IsInclusive flag
   - Priority-based rule selection
   - Soft delete support
```

**Service:** ✅ [internal/modules/taxes/service/service.go](../go/aystek-core/internal/modules/taxes/service/service.go)

```go
✅ Service interface with:
   - Create/Update/Delete/Get/List
   - Calculate() method for tax calculation
   - Rule filtering by:
     - Country/Region/Postal
     - ProductID/CategoryID
     - CustomerType (B2B vs B2C)
     - Scope (order, shipping, product)
   - Inclusive vs Exclusive rate handling
   - Priority-based rule selection (higher priority wins)
   - Soft delete support
```

**Repository:** ✅ [internal/modules/taxes/repository/mongo.go](../go/aystek-core/internal/modules/taxes/repository/mongo.go)

```go
✅ CRUD operations on taxes collection
✅ Indexes for fast lookup (country, active_at, etc.)
```

**Handler:** ✅ [internal/modules/taxes/handler/admin.go](../go/aystek-core/internal/modules/taxes/handler/admin.go)

```go
✅ Admin endpoints:
   POST   /api/v1/taxes/rules
   GET    /api/v1/taxes/rules/:id
   PUT    /api/v1/taxes/rules/:id
   DELETE /api/v1/taxes/rules/:id
   GET    /api/v1/taxes/rules?country=TR&region=Istanbul
```

### Integration with Checkout

**Checkout Service:** ⚠️ Partial [internal/modules/checkout/service/service.go](../go/aystek-core/internal/modules/checkout/service/service.go)

```go
// Receives taxService in constructor
func New(..., taxSvc taxservice.Service) Service

// CalculateTotals() calls tax service
func (s *checkoutService) CalculateTotals(ctx context.Context, sessionID primitive.ObjectID) (*model.CheckoutSession, error) {
    // Calculates tax using taxSvc.Calculate()
    // But: Checkout race condition exists (TODO comment)
}
```

**TODO/FIXME Comments Found:**
- Line 85: "TODO: Concurrent checkout session handling"
- Line 120: "FIXME: Race condition on discount application"

---

## 2. Gaps Analysis (Ne Eksik?)

### ❌ Critical Missing: EU B2B Reverse Charge

**Requirement:** When B2B company in EU sells to another EU company, **reverse charge** applies:
- Tax is NOT collected at sale
- Buyer's account is charged (self-assessed tax)
- Seller issues invoice with "0% VAT — Reverse Charge"

**Current Gap:**
```go
// Tax rule model has CustomerType (B2B, B2C)
// But NO mechanism to determine if SELLER is EU vs BUYER is EU
// No "reverse_charge" flag on invoice
// No "self_assessed_tax" tracking
```

**What Needs To Be Added:**
1. Store EU registration marker (Company model)
2. Invoice reverse charge flag
3. Tax calculation logic: "If seller.eu AND buyer.eu → apply reverse charge"

### ⚠️ Partial: VAT Exemption Codes

**Requirement:** Some items (books, medicine, children's clothes) are **tax-exempt** by law

**Current Status:**
```go
// Tax calculation has IsExempt check:
if input.IsExempt {
    return zero tax
}

// But: Exemption codes not tracked on products
// Manually passed at checkout time
```

**What Needs To Be Added:**
1. Product model: `TaxExemptionCode` field (e.g., "BOOK", "MEDICINE", "CHILDREN_CLOTHING")
2. Tax service: lookup exemption rules automatically from product
3. Tests: Verify exemption applied correctly

### ⚠️ Partial: Tax Bracket Calculation

**Requirement:** Some countries have **tiered tax rates** based on amount:
- Net < 1000 → 8% tax
- 1000-5000 → 15% tax
- 5000+ → 20% tax

**Current Status:**
```go
// Tax calculation is LINEAR:
tax = amount * rate / 100
// No bracket logic
```

**What Needs To Be Added:**
1. Tax rule model: `Brackets []TaxBracket` (optional)
2. Tax service: Bracket calculation logic
3. Tests: Verify bracket thresholds

### ⚠️ Partial: Withholding Tax (Stopaj)

**Requirement:** Some B2B transactions require **withholding tax**:
- Seller remits 20% to government on behalf of buyer
- Buyer gets credit

**Current Status:**
```go
// No withholding tax field
// No withholding calculation
```

**What Needs To Be Added:**
1. Tax rule model: `WithholdingRate` field
2. Invoice model: `WithholdingAmount` field
3. Tax service: Calculate withholding separately
4. Payment service: Track withholding credits

---

## 3. Implementation Plan (20 hours)

### Phase 1: Core Tax Calculation Fixes (6h)

#### Task 1.1: Add Concurrent Checkout Protection (2h)

**File:** [internal/modules/checkout/service/service.go](../go/aystek-core/internal/modules/checkout/service/service.go#L85)

**Problem:** Same checkout session can be modified by 2 concurrent requests

**Solution:** Add version field + optimistic locking (similar to Orders)

```go
// checkout/model/checkout_session.go
type CheckoutSession struct {
    // ... existing fields ...
    Version int64 `bson:"version"`  // Add version field
}

// checkout/service/service.go
func (s *checkoutService) CalculateTotals(ctx context.Context, sessionID primitive.ObjectID, method string) (*model.CheckoutSession, error) {
    session, err := s.repo.GetByID(ctx, sessionID)
    
    // Calculate totals...
    
    // Update with version check
    updated, err := s.repo.UpdateWithVersion(ctx, sessionID, session.Version, updateData)
    if err == ErrVersionMismatch {
        return nil, fmt.Errorf("checkout was modified; retry")
    }
    return updated, err
}

// checkout/repository/mongo.go
func (r *mongoRepository) UpdateWithVersion(ctx context.Context, id primitive.ObjectID, version int64, update bson.M) (*model.CheckoutSession, error) {
    filter := bson.M{"_id": id, "version": version}
    updateOp := bson.M{
        "$set": update,
        "$inc": bson.M{"version": 1},
    }
    result, _ := col.UpdateOne(ctx, filter, updateOp)
    if result.ModifiedCount == 0 {
        return nil, ErrVersionMismatch
    }
    return r.GetByID(ctx, id)
}
```

**Test:** `checkout/service/service_test.go`

```go
func TestCheckoutCalculateTotalsVersionConflict(t *testing.T) {
    // 2 concurrent CalculateTotals calls on same session
    // Second should get ErrVersionMismatch
}
```

**Acceptance Criteria:**
- [ ] Concurrent checkout calls: second gets version error
- [ ] Retry logic works (fetch latest, retry)
- [ ] Test passing

#### Task 1.2: Fix Discount Application Race Condition (2h)

**File:** [internal/modules/checkout/service/service.go](../go/aystek-core/internal/modules/checkout/service/service.go#L120)

**Problem:** Discount application not atomic with total calculation

**Solution:** Single transaction for discount + recalculate

```go
// checkout/service/service.go
func (s *checkoutService) ApplyDiscount(ctx context.Context, sessionID primitive.ObjectID, couponCode string) (*model.CheckoutSession, error) {
    // Start transaction
    session, err := s.repo.GetByID(ctx, sessionID)
    
    // Validate coupon (within transaction)
    coupon, err := s.discountSvc.ValidateCoupon(ctx, couponCode, session.Subtotal)
    if err != nil {
        return nil, err
    }
    
    // Apply discount
    session.DiscountCode = couponCode
    session.DiscountAmount = coupon.Amount
    
    // Recalculate totals atomically
    session.GrandTotal = session.Subtotal + session.TaxTotal + session.ShippingCost - session.DiscountAmount
    session.Version += 1  // Manual version increment for transaction
    
    // Save atomically
    err = s.repo.UpdateAtomic(ctx, session)
    if err != nil {
        return nil, err
    }
    
    return session, nil
}
```

**Test:** `checkout/service/service_test.go`

```go
func TestApplyDiscountAtomic(t *testing.T) {
    // Apply discount + 2 concurrent total recalculations
    // Final state should be consistent
}
```

**Acceptance Criteria:**
- [ ] Discount applied and totals recalculated in single operation
- [ ] No race condition in concurrent scenarios
- [ ] Test passing

#### Task 1.3: Add Product Tax Exemption Support (2h)

**File:** [internal/modules/products/model/product.go](../go/aystek-core/internal/modules/products/model/product.go)

```go
type Product struct {
    // ... existing fields ...
    
    // ✅ ADD: Tax exemption
    TaxExemptionCode *string  `bson:"tax_exemption_code,omitempty"`  // e.g., "BOOK", "MEDICINE"
    IsTaxExempt      bool     `bson:"is_tax_exempt,omitempty"`       // Quick flag
}
```

**Exemption Codes (Enum):**

```go
// taxes/model/exemption_codes.go
const (
    ExemptionCodeBook         = "BOOK"
    ExemptionCodeMedicine     = "MEDICINE"
    ExemptionCodeChildrens    = "CHILDRENS_CLOTHING"
    ExemptionCodeFood         = "FOOD_BASIC"
    ExemptionCodeNewspapers   = "NEWSPAPERS"
)

var exemptionCodes = map[string]bool{
    "BOOK": true,
    "MEDICINE": true,
    // ... etc
}
```

**Service Update:** [internal/modules/checkout/service/service.go](../go/aystek-core/internal/modules/checkout/service/service.go)

```go
// When calculating tax for each line item:
func (s *checkoutService) calculateLineItemTax(ctx context.Context, item CartItem, product *Product) (float64, error) {
    // Check if product has exemption code
    if product.IsTaxExempt {
        return 0, nil
    }
    
    // Otherwise, calculate normally
    return s.taxSvc.Calculate(ctx, &TaxCalculateInput{...})
}
```

**Test:** `checkout/service/checkout_tax_exemption_test.go`

```go
func TestCheckoutTaxExemptionBook(t *testing.T) {
    // Add book to cart (IsTaxExempt=true)
    // Calculate totals
    // Verify TaxTotal = 0
}

func TestCheckoutMixedExemptNonExempt(t *testing.T) {
    // Cart: 1 book (exempt) + 1 magazine (taxed)
    // Verify only magazine charged tax
}
```

**Acceptance Criteria:**
- [ ] Product model has tax exemption field
- [ ] Checkout skips tax for exempt items
- [ ] Mixed items: only non-exempt charged
- [ ] Tests passing

---

### Phase 2: EU B2B Reverse Charge (8h)

#### Task 2.1: Add Company EU Registration Tracking (2h)

**File:** [internal/modules/b2b/model/company.go](../go/aystek-core/internal/modules/b2b/model/company.go)

```go
type Company struct {
    // ... existing fields ...
    
    // ✅ ADD: EU tax info
    IsEURegistered    bool     `bson:"is_eu_registered,omitempty"`  // Is seller registered in EU?
    EUVATNumber       string   `bson:"eu_vat_number,omitempty"`     // Format: DE123456789
    EUCountryCode     string   `bson:"eu_country_code,omitempty"`   // e.g., "DE", "FR", "TR"
}
```

**Validation:**

```go
// b2b/service/service.go
func (s *service) ValidateEUVAT(ctx context.Context, vatNumber string) (bool, error) {
    // Format: {CountryCode}{NumberWithChecksum}
    // e.g., DE123456789 (9 digits after country code)
    
    // Validate format
    if !strings.HasPrefix(vatNumber, "DE") && !strings.HasPrefix(vatNumber, "FR") { // ... etc
        return false, ErrInvalidEUCountry
    }
    
    // For now: Basic format check (production: call VIES API)
    // https://ec.europa.eu/taxation_customs/vies/
    
    return true, nil
}
```

**Test:** `b2b/service/service_test.go`

```go
func TestValidateEUVATNumber(t *testing.T) {
    cases := []struct {
        vat    string
        valid  bool
    }{
        {"DE123456789", true},
        {"FR12345678", true},
        {"IT12345678901", true},
        {"INVALID", false},
    }
    // Test each case
}
```

**Acceptance Criteria:**
- [ ] Company model tracks EU registration
- [ ] VAT number validation logic
- [ ] Country code extraction
- [ ] Tests passing

#### Task 2.2: Add Invoice Reverse Charge Tracking (2h)

**File:** [internal/modules/billing/invoices/model/invoice.go](../go/aystek-core/internal/modules/billing/invoices/model/invoice.go)

```go
type Invoice struct {
    // ... existing fields ...
    
    // ✅ ADD: Reverse charge fields
    IsReverseCharge   bool    `bson:"is_reverse_charge,omitempty"`   // Reverse charge applied?
    ReverseChargeNote string  `bson:"reverse_charge_note,omitempty"` // e.g., "VAT reverse charge — Article 199"
    
    // ✅ ADD: Self-assessed tax tracking
    SelfAssessedTax   float64 `bson:"self_assessed_tax,omitempty"`   // Tax buyer must pay to government
}
```

**Invoice Hook Update:** [internal/modules/orders/service/hooks_invoice.go](../go/aystek-core/internal/modules/orders/service/hooks_invoice.go)

```go
func (h *InvoiceHook) createInvoice(ctx context.Context, order *model.Order) error {
    // Check if reverse charge applies
    isReverseCharge := h.shouldApplyReverseCharge(ctx, order)
    
    invoiceInput := &invoices.CreateInvoiceInput{
        // ... existing fields ...
        
        // ✅ ADD: Reverse charge flag
        IsReverseCharge: isReverseCharge,
        
        // If reverse charge: Tax amount goes to SelfAssessedTax (not seller tax liability)
        TaxTotal: 0,  // No tax collected by seller
        SelfAssessedTax: h.originalTax,  // Buyer pays this to government
    }
    
    return h.invoiceService.Create(ctx, invoiceInput)
}

func (h *InvoiceHook) shouldApplyReverseCharge(ctx context.Context, order *model.Order) bool {
    // Get seller company (from tenant)
    seller, _ := h.companyService.GetByTenant(ctx, order.TenantID)
    
    // Get buyer company
    buyer, _ := h.companyService.GetByID(ctx, order.CustomerID)
    
    // Apply reverse charge if:
    // 1. Seller is EU registered
    // 2. Buyer is B2B company
    // 3. Buyer is in different EU country
    if seller.IsEURegistered && buyer != nil && buyer.IsEURegistered {
        return seller.EUCountryCode != buyer.EUCountryCode
    }
    
    return false
}
```

**Test:** `orders/service/hooks_invoice_reverse_charge_test.go`

```go
func TestReverseChargeEUtoEU(t *testing.T) {
    // Seller: Germany, EU registered
    // Buyer: France, EU registered
    // Expected: Invoice has IsReverseCharge=true, TaxTotal=0, SelfAssessedTax=original_tax
}

func TestNoReverseChargeSamecountry(t *testing.T) {
    // Seller: Germany
    // Buyer: Germany
    // Expected: Normal invoice, no reverse charge
}

func TestNoReverseChargeB2C(t *testing.T) {
    // Seller: Germany, EU registered
    // Buyer: Individual (not company)
    // Expected: Normal invoice, no reverse charge
}
```

**Acceptance Criteria:**
- [ ] Invoice model has reverse charge fields
- [ ] Hook logic detects reverse charge scenarios
- [ ] Tax calculation adjusted (0% for seller)
- [ ] Self-assessed tax tracked
- [ ] Tests passing

#### Task 2.3: Reverse Charge E-Invoice Formatting (2h)

**E-Invoice Format:** [internal/modules/billing/invoices/model/invoice.go](../go/aystek-core/internal/modules/billing/invoices/model/invoice.go)

When reverse charge applies, e-invoice XML must include:
```xml
<cbc:LineItem>
    <cac:InvoicedQuantity unitCode="C62">1</cac:InvoicedQuantity>
    <cbc:LineExtensionAmount>100.00</cbc:LineExtensionAmount>
    <cac:TaxTotal>
        <cbc:TaxAmount>0.00</cbc:TaxAmount>  <!-- 0% for reverse charge -->
    </cac:TaxTotal>
    <!-- ✅ ADD: Reverse charge indicator -->
    <cac:TaxExemption>
        <cbc:ExemptionReasonCode>REVERSE</cbc:ExemptionReasonCode>
        <cbc:ExemptionReason>VAT reverse charge — Cross-border supply of services</cbc:ExemptionReason>
    </cac:TaxExemption>
</cac:LineItem>
```

**Implementation:** [internal/modules/billing/invoices/service/e_invoice_builder.go](../go/aystek-core/internal/modules/billing/invoices/service/e_invoice_builder.go)

```go
func (b *eInvoiceBuilder) buildLineItemXML(line *model.InvoiceLine, invoice *model.Invoice) string {
    // ... existing code ...
    
    // ✅ ADD: Reverse charge handling
    taxExemption := ""
    if invoice.IsReverseCharge {
        taxExemption = fmt.Sprintf(`
            <cac:TaxExemption>
                <cbc:ExemptionReasonCode>REVERSE</cbc:ExemptionReasonCode>
                <cbc:ExemptionReason>%s</cbc:ExemptionReason>
            </cac:TaxExemption>
        `, invoice.ReverseChargeNote)
    }
    
    return fmt.Sprintf(`<cbc:LineItem>...%s...</cbc:LineItem>`, taxExemption)
}
```

**Test:** `invoices/service/e_invoice_reverse_charge_test.go`

```go
func TestEInvoiceReverseChargeXML(t *testing.T) {
    invoice := &model.Invoice{
        IsReverseCharge: true,
        ReverseChargeNote: "VAT reverse charge",
        Lines: []InvoiceLine{{...}},
    }
    
    xml := buildEInvoiceXML(invoice)
    
    // Verify:
    // - TaxAmount = 0
    // - ExemptionReasonCode = "REVERSE"
    // - ExemptionReason contains text
}
```

**Acceptance Criteria:**
- [ ] E-invoice XML includes reverse charge indicators
- [ ] Tax amount = 0 in XML
- [ ] Exemption reason populated
- [ ] Tests passing

#### Task 2.4: B2B Tax Reporting Integration (2h)

**Invoice List Filters:** [internal/modules/billing/invoices/handler/admin.go](../go/aystek-core/internal/modules/billing/invoices/handler/admin.go)

```go
// GET /api/v1/invoices?invoice_type=reverse_charge

// Add filter option:
type ListInvoicesRequest struct {
    // ... existing filters ...
    
    InvoiceType string `query:"invoice_type"` // "normal", "reverse_charge", "credit_memo"
}
```

**Reporting View:** Create dashboard endpoint for tax compliance

```go
// GET /api/v1/reports/tax-compliance?period=2024-01&country=DE

type TaxComplianceReport struct {
    Period                  string     // "2024-01"
    Country                 string     // "DE"
    
    TotalInvoiced           float64    // All invoices
    TotalTaxCollected       float64    // Normal tax
    TotalReverseCharged     float64    // Reverse charge (buyer's liability)
    TotalSelfAssessedTax    float64    // Buyer's tax to pay
    
    InvoiceCount            int64
    ReverseChargeCount      int64
    
    ByCountry map[string]struct {
        Count       int64
        Amount      float64
        TaxAmount   float64
    }
}
```

**Test:** `invoices/handler/tax_compliance_report_test.go`

```go
func TestTaxComplianceReportReverseCharge(t *testing.T) {
    // Create invoices (normal + reverse charge)
    // Get compliance report
    // Verify counts and amounts
}
```

**Acceptance Criteria:**
- [ ] Invoice list filtered by type
- [ ] Tax compliance report endpoint working
- [ ] Reverse charge amounts separated
- [ ] Tests passing

---

### Phase 3: VAT Bracket Calculations (4h)

#### Task 3.1: Add Tax Bracket Support to Model (1h)

**File:** [internal/modules/taxes/model/tax_rule.go](../go/aystek-core/internal/modules/taxes/model/tax_rule.go)

```go
type TaxBracket struct {
    MinAmount float64  // Minimum amount for this bracket
    MaxAmount *float64 // Maximum (null = unlimited)
    Rate      float64  // Tax rate for this bracket
}

type TaxRule struct {
    // ... existing fields ...
    
    // ✅ ADD: Bracket support
    Brackets []TaxBracket `bson:"brackets,omitempty"`  // If set, use brackets instead of single rate
}
```

**Example (Turkey income tax):**
```go
Brackets: []TaxBracket{
    {MinAmount: 0, MaxAmount: ptr(28000), Rate: 0.15},
    {MinAmount: 28000, MaxAmount: ptr(70000), Rate: 0.20},
    {MinAmount: 70000, MaxAmount: nil, Rate: 0.27},
}
```

#### Task 3.2: Implement Bracket Calculation Logic (1.5h)

**File:** [internal/modules/taxes/service/service.go](../go/aystek-core/internal/modules/taxes/service/service.go)

```go
func (s *service) Calculate(ctx context.Context, input CalculateInput) (*CalculateResult, error) {
    // ... existing code ...
    
    rule := filtered[0]
    amount := input.Amount
    tax := 0.0
    
    // ✅ ADD: Bracket calculation
    if len(rule.Brackets) > 0 {
        tax = s.calculateBracketTax(amount, rule.Brackets)
    } else {
        // Single rate (existing logic)
        if rule.IsInclusive {
            tax = amount * rule.Rate / (100 + rule.Rate)
        } else {
            tax = amount * rule.Rate / 100
        }
    }
    
    return &CalculateResult{...}, nil
}

func (s *service) calculateBracketTax(amount float64, brackets []TaxBracket) float64 {
    sort.Slice(brackets, func(i, j int) bool {
        return brackets[i].MinAmount < brackets[j].MinAmount
    })
    
    tax := 0.0
    remaining := amount
    
    for _, bracket := range brackets {
        if remaining <= 0 {
            break
        }
        
        maxForBracket := bracket.MaxAmount
        if maxForBracket == nil {
            maxForBracket = ptr(remaining) // Use remaining if no max
        }
        
        bracketMin := bracket.MinAmount
        bracketMax := *maxForBracket
        
        // Calculate taxable in this bracket
        taxableInBracket := 0.0
        if remaining > bracketMin {
            taxableInBracket = min(remaining - bracketMin, bracketMax - bracketMin)
            tax += taxableInBracket * bracket.Rate / 100
            remaining -= taxableInBracket
        }
    }
    
    return tax
}
```

#### Task 3.3: Add Bracket Tests (1.5h)

**Test File:** `taxes/service/tax_bracket_test.go`

```go
func TestTaxBracketCalculation(t *testing.T) {
    brackets := []model.TaxBracket{
        {MinAmount: 0, MaxAmount: ptr(100), Rate: 0.10},     // 10% on first 100
        {MinAmount: 100, MaxAmount: ptr(200), Rate: 0.15},   // 15% on next 100
        {MinAmount: 200, MaxAmount: nil, Rate: 0.20},        // 20% on rest
    }
    
    cases := []struct {
        amount   float64
        expected float64
    }{
        {50, 5.0},     // 50 * 0.10 = 5
        {100, 10.0},   // 100 * 0.10 = 10
        {150, 17.5},   // 100*0.10 + 50*0.15 = 10 + 7.5 = 17.5
        {300, 50.0},   // 100*0.10 + 100*0.15 + 100*0.20 = 10 + 15 + 20 = 45
    }
    
    for _, c := range cases {
        result := calculateBracketTax(c.amount, brackets)
        assert.Equal(t, c.expected, result)
    }
}

func TestTaxBracketIntegration(t *testing.T) {
    // Create tax rule with brackets
    // Create invoice with multiple items
    // Verify tax calculated using brackets
}
```

**Acceptance Criteria:**
- [ ] Tax bracket model defined
- [ ] Bracket calculation logic working
- [ ] Multiple bracket scenarios tested
- [ ] Edge cases (amount = bracket boundary) handled
- [ ] Tests passing

---

## 4. Testing Strategy (5h)

### Unit Tests

```go
✅ Tax calculation (single rate)
✅ Tax calculation (inclusive vs exclusive)
✅ Tax brackets
✅ Product tax exemption
✅ Checkout version control
✅ Discount application atomicity
✅ EU VAT number validation
✅ Reverse charge detection
✅ E-invoice XML generation (reverse charge)
```

### Integration Tests

```go
✅ Checkout → Tax calculation → Order creation
✅ Order creation → Invoice generation (with reverse charge if applicable)
✅ Tax compliance reporting
✅ Multi-item checkout with mixed exempt/taxed items
```

### E2E Scenarios

```
Scenario 1: B2C Turkey Purchase
├─ Item: Magazine (normal taxed)
├─ Tax: 18%
├─ Expected: Invoice with TaxTotal = 18% of subtotal
└─ Result: ✅

Scenario 2: B2C EU Purchase (Reverse Charge)
├─ Seller: DE (EU registered)
├─ Buyer: FR (EU registered, B2B)
├─ Tax: 19%
├─ Expected: Invoice with IsReverseCharge=true, TaxTotal=0, SelfAssessedTax=19%
└─ Result: ✅

Scenario 3: Mixed Cart (Exempt + Taxed)
├─ Item 1: Book (exempt)
├─ Item 2: Pen (18% tax)
├─ Expected: TaxTotal = (Pen price) * 0.18
└─ Result: ✅

Scenario 4: Bracketed Tax
├─ Income: 250 TR
├─ Brackets: 0-100 (10%), 100-200 (15%), 200+ (20%)
├─ Expected: 100*0.10 + 100*0.15 + 50*0.20 = 10 + 15 + 10 = 35
└─ Result: ✅

Scenario 5: Concurrent Checkout (Race Condition)
├─ 2 requests: Apply discount + Calculate totals simultaneously
├─ Expected: No data loss, consistent final state
└─ Result: ✅
```

---

## 5. Delivery Checklist (P1-G1)

### Code Changes
- [ ] Checkout version field + UpdateWithVersion
- [ ] Checkout discount atomic operation
- [ ] Product tax exemption field
- [ ] Company EU registration fields
- [ ] Invoice reverse charge fields
- [ ] Invoice hook reverse charge logic
- [ ] E-invoice reverse charge XML
- [ ] Tax bracket model + calculation
- [ ] Tax compliance reporting endpoint

### Tests
- [ ] 15+ unit tests passing
- [ ] 5+ integration tests passing
- [ ] 5 E2E scenarios validated
- [ ] Load test: 50 concurrent checkouts (success)
- [ ] Coverage: 85%+ on tax modules

### Documentation
- [ ] Tax rule configuration guide
- [ ] Reverse charge implementation notes
- [ ] E-invoice provider integration (iyzico, etc.)
- [ ] Admin dashboard (tax compliance report)
- [ ] Troubleshooting guide

### Deployment Prep
- [ ] Database indexes verified
- [ ] Migration for new fields (backward compatible)
- [ ] Feature flags for reverse charge (default: enabled if seller.eu)
- [ ] Monitoring: Tax calculation latency
- [ ] Rollback procedure

---

## Effort Summary

| Task | Hours |
|------|-------|
| 1.1 Checkout concurrent protection | 2 |
| 1.2 Discount application fix | 2 |
| 1.3 Product tax exemption | 2 |
| 2.1 Company EU registration | 2 |
| 2.2 Invoice reverse charge | 2 |
| 2.3 E-invoice reverse charge XML | 2 |
| 2.4 Tax compliance reporting | 2 |
| 3.1 Tax bracket model | 1 |
| 3.2 Bracket calculation | 1.5 |
| 3.3 Bracket tests | 1.5 |
| **Testing & QA** | **5** |
| **Total** | **23** |

**Buffer:** 17 hours (enough for code review, fixes, documentation)

---

# PART II: P1-G2 — B2B WORKFLOWS (24h) — NEXT DOCUMENT

# PART III: P1-G3 — STOCK ALLOCATION (12h) — NEXT DOCUMENT

# PART IV: P1-G4 — RETURNS/RMA (8h) — NEXT DOCUMENT

---

## Summary

**P1-G1 (Tax Rules MVP)** is **23 hours of implementation** with:
- ✅ Core tax calculation ready (leveraging existing service)
- ✅ Checkout concurrency fixes
- ✅ EU B2B reverse charge support
- ✅ Tax exemptions
- ✅ Bracket calculations
- ✅ Comprehensive testing

**Next Steps:** Once P1-G1 complete, proceed to P1-G2 (B2B workflows) with similar detail level.
