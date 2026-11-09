# P1-G2 — B2B Workflows Implementation Plan

**Kapsam:** Company management, credit limits, purchase order (PO) workflow, approval process  
**Efor:** 24 saat  
**Bağlantılar:** P0-G1 Order state machine + P1-G1 Tax rules (reverse charge)

---

## 1. Current State Analysis

### ✅ What Exists

**Model:** [internal/modules/b2b/model/company.go](../go/aystek-core/internal/modules/b2b/model/company.go)

```go
type Company struct {
    ID                primitive.ObjectID
    Name              string
    TaxID             string           // VAT, EIN, etc.
    Industry          string
    Website           string
    
    PrimaryContact   *Contact
    BillingAddress   *Address
    ShippingAddress  *Address
    
    CreditLineID     *primitive.ObjectID  // Reference to credit line
    
    Status           string               // active, inactive, suspended
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

**Model:** [internal/modules/b2b/model/credit_line.go](../go/aystek-core/internal/modules/b2b/model/credit_line.go)

```go
type CreditLine struct {
    ID              primitive.ObjectID
    CompanyID       primitive.ObjectID
    
    LimitAmount     float64      // Total credit allowed
    UsedAmount      float64      // Currently used
    AvailableAmount float64      // LimitAmount - UsedAmount
    
    Status          string       // active, suspended
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Service:** [internal/modules/b2b/service/service.go](../go/aystek-core/internal/modules/b2b/service/service.go)

```go
✅ Basic CRUD operations on company + credit line
✅ Credit balance checking
✅ Company status management
```

---

### ❌ Critical Missing: Purchase Order (PO) Workflow

**Requirement:** B2B customers should use **Purchase Order** instead of immediate checkout:

1. **PO Creation:** Customer creates PO (item list + quantities + delivery date)
2. **PO Approval:** Seller reviews + approves/rejects
3. **PO to Order:** Approved PO auto-converts to Order (with special B2B terms)
4. **Payment Terms:** Net 30 / Net 60 (not immediate payment)
5. **Invoice Sequence:** Multiple invoices (partial delivery) possible

**Current Status:** ❌ NO MODELS, NO SERVICE, NO WORKFLOW

### ❌ Critical Missing: Approval Workflows

**Requirement:** For high-value POs (> credit threshold), **approval chain** needed:

1. Store manager approval (if PO > store threshold)
2. Credit department approval (if PO > company credit limit)
3. Executive approval (if PO > executive threshold)

**Current Status:** ❌ NO APPROVAL SYSTEM

### ❌ Critical Missing: Payment Terms Management

**Requirement:** B2B orders can have **payment terms** (not immediate payment):

- **Net 30:** Pay within 30 days of invoice
- **Net 60:** Pay within 60 days of invoice
- **2/10 Net 30:** 2% discount if paid within 10 days
- **Deposits:** 50% upfront, 50% on delivery

**Current Status:** ❌ NO PAYMENT TERMS MODEL

---

## 2. Implementation Plan (24 hours)

### Phase 1: Purchase Order Model & Service (8h)

#### Task 1.1: Create PO Model (2h)

**File:** `internal/modules/b2b/model/purchase_order.go` (NEW)

```go
package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type POStatus string

const (
    POStatusDraft       POStatus = "draft"       // Created, not submitted
    POStatusSubmitted   POStatus = "submitted"   // Awaiting approval
    POStatusApproved    POStatus = "approved"    // All approvals passed
    POStatusRejected    POStatus = "rejected"    // Rejected in approval
    POStatusPartial     POStatus = "partial"     // Partially fulfilled
    POStatusFulfilled   POStatus = "fulfilled"   // All items delivered
    POStatusCancelled   POStatus = "cancelled"   // Cancelled by buyer/seller
)

type PurchaseOrder struct {
    ID                  primitive.ObjectID   `bson:"_id"`
    PONumber            string               `bson:"po_number"`  // Public-facing: PO-2024-001
    
    // Reference data
    TenantID            string               `bson:"tenant_id"`  // Seller
    CompanyID           primitive.ObjectID   `bson:"company_id"` // Buyer
    ContactPersonID     *primitive.ObjectID  `bson:"contact_person_id,omitempty"`  // Who created PO
    
    // Order info
    Status              POStatus             `bson:"status"`
    Items               []POItem             `bson:"items"`
    
    // Pricing
    Subtotal            float64              `bson:"subtotal"`
    TaxTotal            float64              `bson:"tax_total,omitempty"`
    ShippingCost        float64              `bson:"shipping_cost,omitempty"`
    DiscountTotal       float64              `bson:"discount_total,omitempty"`
    GrandTotal          float64              `bson:"grand_total"`
    Currency            string               `bson:"currency"`
    
    // Payment terms
    PaymentTermsID      *primitive.ObjectID  `bson:"payment_terms_id,omitempty"`
    PaymentTermsDesc    string               `bson:"payment_terms_desc,omitempty"`  // e.g., "Net 30"
    DepositRequired     float64              `bson:"deposit_required,omitempty"`    // 50% upfront?
    DepositPaid         float64              `bson:"deposit_paid,omitempty"`
    
    // Delivery
    DeliveryDate        time.Time            `bson:"delivery_date"`
    DeliveryAddress     *Address             `bson:"delivery_address"`
    ShippingMethod      string               `bson:"shipping_method,omitempty"`
    
    // Approval workflow
    ApprovalStatus      ApprovalStatus       `bson:"approval_status"`
    ApprovalChain       []Approval           `bson:"approval_chain,omitempty"`  // Who needs to approve + status
    
    // Conversion to Order
    OrderID             *primitive.ObjectID  `bson:"order_id,omitempty"`  // Link to converted order
    
    // Metadata
    Notes               string               `bson:"notes,omitempty"`
    Tags                []string             `bson:"tags,omitempty"`
    
    CreatedAt           time.Time            `bson:"created_at"`
    UpdatedAt           time.Time            `bson:"updated_at"`
    SubmittedAt         *time.Time           `bson:"submitted_at,omitempty"`
    ApprovedAt          *time.Time           `bson:"approved_at,omitempty"`
    RejectedAt          *time.Time           `bson:"rejected_at,omitempty"`
    
    CreatedBy           string               `bson:"created_by"`  // Company user ID
    UpdatedBy           string               `bson:"updated_by"`
}

type POItem struct {
    ProductID           primitive.ObjectID   `bson:"product_id"`
    VariantID           *primitive.ObjectID  `bson:"variant_id,omitempty"`
    SKU                 string               `bson:"sku,omitempty"`
    Name                string               `bson:"name"`
    Quantity            int                  `bson:"quantity"`
    UnitPrice           float64              `bson:"unit_price"`
    LineTotal           float64              `bson:"line_total"`
    
    FulfilledQuantity   int                  `bson:"fulfilled_quantity,omitempty"`  // Partial delivery
    DeliverySchedule    *DeliverySchedule   `bson:"delivery_schedule,omitempty"`   // Phased delivery?
}

type DeliverySchedule struct {
    Quantity            int       `bson:"quantity"`
    PlannedDate         time.Time `bson:"planned_date"`
    ActualDate          *time.Time `bson:"actual_date,omitempty"`
}

type ApprovalStatus string

const (
    ApprovalStatusPending  ApprovalStatus = "pending"   // Waiting for approvals
    ApprovalStatusApproved ApprovalStatus = "approved"  // All approved
    ApprovalStatusRejected ApprovalStatus = "rejected"  // At least one rejected
)

type Approval struct {
    ApproverID          string                `bson:"approver_id"`
    ApproverRole        string                `bson:"approver_role"`  // "store_manager", "credit_dept", "executive"
    ApprovalThreshold   float64               `bson:"approval_threshold"`  // Why this person needed to approve
    Status              string                `bson:"status"`  // "pending", "approved", "rejected"
    DecisionAt          *time.Time            `bson:"decision_at,omitempty"`
    Notes               string                `bson:"notes,omitempty"`
}
```

**Index:** [internal/modules/b2b/repository/po_indexes.go](../go/aystek-core/internal/modules/b2b/repository/po_indexes.go)

```go
// Po number per tenant (unique)
{Key: bson.D{{Key: "tenant_id", Value: 1}, {Key: "po_number", Value: 1}}, Options: options.Index().SetUnique(true)}

// Company POs
{Key: bson.D{{Key: "company_id", Value: 1}, {Key: "created_at", Value: -1}}}

// Status queries
{Key: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}}

// Approval chain
{Key: bson.D{{Key: "approval_chain.approver_id", Value: 1}, {Key: "approval_status", Value: 1}}}
```

#### Task 1.2: PO Service Interface & CRUD (2h)

**File:** `internal/modules/b2b/service/po_service.go` (NEW)

```go
package service

import (
    "context"
    "github.com/aystek/aystek-core/internal/modules/b2b/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type POService interface {
    // CRUD
    CreatePO(ctx context.Context, input CreatePOInput) (*model.PurchaseOrder, error)
    GetPO(ctx context.Context, id primitive.ObjectID) (*model.PurchaseOrder, error)
    ListPOs(ctx context.Context, filter POListFilter) ([]model.PurchaseOrder, int64, error)
    UpdatePO(ctx context.Context, id primitive.ObjectID, input UpdatePOInput) (*model.PurchaseOrder, error)
    
    // Workflow
    SubmitPO(ctx context.Context, id primitive.ObjectID, submittedBy string) (*model.PurchaseOrder, error)
    CancelPO(ctx context.Context, id primitive.ObjectID, reason string, cancelledBy string) (*model.PurchaseOrder, error)
    
    // Conversion
    ConvertPOToOrder(ctx context.Context, id primitive.ObjectID) (*Order, error)
    
    // Validation
    ValidatePOForApproval(ctx context.Context, po *model.PurchaseOrder) error
}

type CreatePOInput struct {
    CompanyID           primitive.ObjectID
    Items               []POItemInput
    DeliveryDate        time.Time
    DeliveryAddress     *model.Address
    PaymentTermsID      *primitive.ObjectID
    DepositRequired     *float64  // e.g., 0.50 for 50%
    ShippingMethod      string    // e.g., "aras_kargo"
    Notes               string
    CreatedBy           string
}

type POItemInput struct {
    ProductID           primitive.ObjectID
    VariantID           *primitive.ObjectID
    Quantity            int
    UnitPrice           *float64  // If not provided, use current product price
}

type UpdatePOInput struct {
    Status              *POStatus
    DeliveryDate        *time.Time
    ShippingMethod      *string
    Notes               *string
    UpdatedBy           string
}

type POListFilter struct {
    CompanyID           *primitive.ObjectID
    Status              *POStatus
    SubmittedAfter      *time.Time
    CreatedBy           *string
    Page                int
    Limit               int
}
```

**Implementation:** [internal/modules/b2b/service/po_service_impl.go](../go/aystek-core/internal/modules/b2b/service/po_service_impl.go)

```go
func (s *b2bService) CreatePO(ctx context.Context, input CreatePOInput) (*model.PurchaseOrder, error) {
    // 1. Validate company exists + is active
    company, _ := s.companyRepo.GetByID(ctx, input.CompanyID)
    if company.Status != "active" {
        return nil, ErrCompanyInactive
    }
    
    // 2. Build line items (fetch product prices)
    lines := make([]model.POItem, 0, len(input.Items))
    subtotal := 0.0
    
    for _, itemInput := range input.Items {
        product, _ := s.productSvc.Get(ctx, itemInput.ProductID.Hex())
        
        unitPrice := itemInput.UnitPrice
        if unitPrice == nil {
            unitPrice = &product.BasePrice
        }
        
        lineTotal := float64(itemInput.Quantity) * *unitPrice
        subtotal += lineTotal
        
        lines = append(lines, model.POItem{
            ProductID: itemInput.ProductID,
            VariantID: itemInput.VariantID,
            Name:      product.Name["en"],  // TODO: use company's language
            Quantity:  itemInput.Quantity,
            UnitPrice: *unitPrice,
            LineTotal: lineTotal,
        })
    }
    
    // 3. Calculate tax (using tax service)
    taxTotal, _ := s.taxSvc.CalculateTaxForPO(ctx, &TaxForPOInput{
        Items:             lines,
        CompanyID:         input.CompanyID,
        DeliveryAddress:   input.DeliveryAddress,
    })
    
    // 4. Calculate shipping
    shippingCost := 0.0  // TODO: lookup shipping
    
    // 5. Generate PO number (e.g., PO-2024-00001)
    poNumber := s.generatePONumber(ctx)
    
    // 6. Build PO
    po := model.PurchaseOrder{
        ID:               primitive.NewObjectID(),
        PONumber:         poNumber,
        TenantID:         s.tenantID,
        CompanyID:        input.CompanyID,
        Status:           model.POStatusDraft,
        Items:            lines,
        Subtotal:         subtotal,
        TaxTotal:         taxTotal,
        ShippingCost:     shippingCost,
        GrandTotal:       subtotal + taxTotal + shippingCost,
        Currency:         "TRY",  // TODO: from company/tenant
        DeliveryDate:     input.DeliveryDate,
        DeliveryAddress:  input.DeliveryAddress,
        ShippingMethod:   input.ShippingMethod,
        Notes:            input.Notes,
        ApprovalStatus:   model.ApprovalStatusPending,
        CreatedAt:        time.Now(),
        UpdatedAt:        time.Now(),
        CreatedBy:        input.CreatedBy,
        UpdatedBy:        input.CreatedBy,
    }
    
    // 7. Save to DB
    id, err := s.poRepo.Insert(ctx, po)
    if err != nil {
        return nil, err
    }
    
    return s.poRepo.GetByID(ctx, id)
}

func (s *b2bService) SubmitPO(ctx context.Context, id primitive.ObjectID, submittedBy string) (*model.PurchaseOrder, error) {
    po, _ := s.poRepo.GetByID(ctx, id)
    
    if po.Status != model.POStatusDraft {
        return nil, ErrPOInvalidStatus
    }
    
    // 1. Validate PO
    if err := s.ValidatePOForApproval(ctx, po); err != nil {
        return nil, err
    }
    
    // 2. Build approval chain (who needs to approve?)
    approvalChain := s.buildApprovalChain(ctx, po)
    
    // 3. Update status + submission time
    now := time.Now()
    po.Status = model.POStatusSubmitted
    po.SubmittedAt = &now
    po.ApprovalChain = approvalChain
    po.UpdatedBy = submittedBy
    po.UpdatedAt = now
    
    // 4. Trigger notifications (approvers notified)
    s.notifyApprovers(ctx, po)
    
    return s.poRepo.Update(ctx, id, po)
}

func (s *b2bService) buildApprovalChain(ctx context.Context, po *model.PurchaseOrder) []model.Approval {
    approvals := []model.Approval{}
    
    // Threshold 1: Store manager (PO > 50,000 TRY)
    if po.GrandTotal > 50000 {
        approvals = append(approvals, model.Approval{
            ApproverRole:       "store_manager",
            ApprovalThreshold:  50000,
            Status:             "pending",
        })
    }
    
    // Threshold 2: Credit dept (PO > company credit limit)
    company, _ := s.companyRepo.GetByID(ctx, po.CompanyID)
    creditLine, _ := s.creditLineRepo.GetByCompany(ctx, company.ID)
    
    if po.GrandTotal > creditLine.AvailableAmount {
        approvals = append(approvals, model.Approval{
            ApproverRole:       "credit_department",
            ApprovalThreshold:  creditLine.LimitAmount,
            Status:             "pending",
        })
    }
    
    // Threshold 3: Executive (PO > 500,000 TRY)
    if po.GrandTotal > 500000 {
        approvals = append(approvals, model.Approval{
            ApproverRole:       "executive",
            ApprovalThreshold:  500000,
            Status:             "pending",
        })
    }
    
    return approvals
}
```

#### Task 1.3: PO Conversion to Order (2h)

**File:** [internal/modules/b2b/service/po_to_order.go](../go/aystek-core/internal/modules/b2b/service/po_to_order.go) (NEW)

```go
func (s *b2bService) ConvertPOToOrder(ctx context.Context, poID primitive.ObjectID) (*Order, error) {
    po, _ := s.poRepo.GetByID(ctx, poID)
    
    // 1. PO must be approved
    if po.Status != model.POStatusApproved {
        return nil, ErrPONotApproved
    }
    
    // 2. Check credit availability
    company, _ := s.companyRepo.GetByID(ctx, po.CompanyID)
    creditLine, _ := s.creditLineRepo.GetByCompany(ctx, company.ID)
    
    if po.GrandTotal > creditLine.AvailableAmount {
        return nil, ErrInsufficientCredit
    }
    
    // 3. If deposit required: check if paid
    if po.DepositRequired > 0 {
        if po.DepositPaid < po.DepositRequired {
            return nil, ErrDepositNotPaid
        }
    }
    
    // 4. Build Order from PO
    orderItems := make([]order.OrderItem, 0, len(po.Items))
    for _, poItem := range po.Items {
        orderItems = append(orderItems, order.OrderItem{
            ProductID:  poItem.ProductID,
            VariantID:  poItem.VariantID,
            Quantity:   poItem.Quantity,
            UnitPrice:  poItem.UnitPrice,
            TotalPrice: poItem.LineTotal,
            Name:       poItem.Name,
        })
    }
    
    // 5. Create Order
    orderInput := order.OrderInput{
        CustomerID:        po.CompanyID,
        Items:             orderItems,
        Subtotal:          po.Subtotal,
        TaxTotal:          po.TaxTotal,
        ShippingTotal:     po.ShippingCost,
        Currency:          po.Currency,
        ShippingAddress:   po.DeliveryAddress,
        PaymentTermsID:    po.PaymentTermsID,  // ← Pass B2B payment terms
        PONumber:          po.PONumber,         // ← Link to original PO
        Status:            order.StatusPending,
        CreatedBy:         "po-conversion",
    }
    
    createdOrder, err := s.orderSvc.Create(ctx, orderInput)
    if err != nil {
        return nil, err
    }
    
    // 6. Update PO with order reference
    po.OrderID = &createdOrder.ID
    po.Status = model.POStatusPartial  // Will update to fulfilled as items ship
    s.poRepo.Update(ctx, poID, po)
    
    // 7. Reserve credit
    creditLine.UsedAmount += po.GrandTotal
    creditLine.AvailableAmount -= po.GrandTotal
    s.creditLineRepo.Update(ctx, creditLine.ID, creditLine)
    
    return createdOrder, nil
}
```

#### Task 1.4: PO Repository (2h)

**File:** `internal/modules/b2b/repository/po_repository.go` (NEW)

```go
package repository

import (
    "context"
    "github.com/aystek/aystek-core/internal/modules/b2b/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type PORepository interface {
    Insert(ctx context.Context, po model.PurchaseOrder) (primitive.ObjectID, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*model.PurchaseOrder, error)
    GetByNumber(ctx context.Context, tenantID, poNumber string) (*model.PurchaseOrder, error)
    List(ctx context.Context, filter POListFilter) ([]model.PurchaseOrder, int64, error)
    Update(ctx context.Context, id primitive.ObjectID, po model.PurchaseOrder) (*model.PurchaseOrder, error)
    Delete(ctx context.Context, id primitive.ObjectID) error
}

type mongoPORepository struct {
    db *mongo.Database
}

func (r *mongoPORepository) Insert(ctx context.Context, po model.PurchaseOrder) (primitive.ObjectID, error) {
    col := r.db.Collection("purchase_orders")
    result, err := col.InsertOne(ctx, po)
    if err != nil {
        return primitive.NilObjectID, err
    }
    return result.InsertedID.(primitive.ObjectID), nil
}

func (r *mongoPORepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.PurchaseOrder, error) {
    col := r.db.Collection("purchase_orders")
    var po model.PurchaseOrder
    err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&po)
    if err != nil {
        return nil, err
    }
    return &po, nil
}

func (r *mongoPORepository) List(ctx context.Context, filter POListFilter) ([]model.PurchaseOrder, int64, error) {
    col := r.db.Collection("purchase_orders")
    
    query := bson.M{}
    if filter.CompanyID != nil {
        query["company_id"] = *filter.CompanyID
    }
    if filter.Status != nil {
        query["status"] = *filter.Status
    }
    
    total, _ := col.CountDocuments(ctx, query)
    
    opts := options.Find().SetSkip(int64((filter.Page - 1) * filter.Limit)).SetLimit(int64(filter.Limit))
    cursor, err := col.Find(ctx, query, opts)
    if err != nil {
        return nil, 0, err
    }
    
    var pos []model.PurchaseOrder
    err = cursor.All(ctx, &pos)
    return pos, total, err
}

func (r *mongoPORepository) Update(ctx context.Context, id primitive.ObjectID, po model.PurchaseOrder) (*model.PurchaseOrder, error) {
    col := r.db.Collection("purchase_orders")
    _, err := col.ReplaceOne(ctx, bson.M{"_id": id}, po)
    if err != nil {
        return nil, err
    }
    return r.GetByID(ctx, id)
}
```

---

### Phase 2: Approval Workflow System (8h)

#### Task 2.1: Approval API Endpoints (2h)

**File:** [internal/modules/b2b/handler/po_approval_handler.go](../go/aystek-core/internal/modules/b2b/handler/po_approval_handler.go)

```go
package handler

import "github.com/gin-gonic/gin"

// PUT /api/v1/b2b/pos/:id/approvals/approve
func (h *POHandler) ApprovePO(c *gin.Context) {
    poID := c.Param("id")
    
    approverID := c.MustGet("user_id").(string)
    approverRole := c.MustGet("user_role").(string)
    
    body := struct {
        Notes string `json:"notes"`
    }{}
    c.BindJSON(&body)
    
    result, err := h.poSvc.ApprovePO(c, objectid(poID), approverID, approverRole, body.Notes)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    
    c.JSON(200, result)
}

// PUT /api/v1/b2b/pos/:id/approvals/reject
func (h *POHandler) RejectPO(c *gin.Context) {
    poID := c.Param("id")
    approverID := c.MustGet("user_id").(string)
    
    body := struct {
        Reason string `json:"reason"`
    }{}
    c.BindJSON(&body)
    
    result, err := h.poSvc.RejectPO(c, objectid(poID), approverID, body.Reason)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    
    c.JSON(200, result)
}

// GET /api/v1/b2b/approvals/pending?role=store_manager
func (h *POHandler) ListPendingApprovals(c *gin.Context) {
    approverRole := c.Query("role")
    approverID := c.MustGet("user_id").(string)
    
    pos, total, err := h.poSvc.ListPendingApprovalsFor(c, approverID, approverRole)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    
    c.JSON(200, gin.H{
        "items": pos,
        "total": total,
    })
}
```

#### Task 2.2: Approval Service Logic (2h)

**File:** [internal/modules/b2b/service/po_approval.go](../go/aystek-core/internal/modules/b2b/service/po_approval.go)

```go
func (s *b2bService) ApprovePO(ctx context.Context, poID primitive.ObjectID, approverID, approverRole, notes string) (*model.PurchaseOrder, error) {
    po, _ := s.poRepo.GetByID(ctx, poID)
    
    // 1. Find this approver in chain
    approvalIdx := -1
    for i, approval := range po.ApprovalChain {
        if approval.ApproverRole == approverRole {
            approvalIdx = i
            break
        }
    }
    
    if approvalIdx == -1 {
        return nil, ErrApproverNotInChain
    }
    
    // 2. Update approval status
    now := time.Now()
    po.ApprovalChain[approvalIdx].Status = "approved"
    po.ApprovalChain[approvalIdx].DecisionAt = &now
    po.ApprovalChain[approvalIdx].Notes = notes
    po.ApprovalChain[approvalIdx].ApproverID = approverID
    
    // 3. Check if all approved
    allApproved := true
    for _, approval := range po.ApprovalChain {
        if approval.Status == "pending" {
            allApproved = false
            break
        }
    }
    
    if allApproved {
        po.Status = model.POStatusApproved
        po.ApprovalStatus = model.ApprovalStatusApproved
        po.ApprovedAt = &now
        
        // Trigger notifications: PO ready to convert to order
        s.notifyPOApproved(ctx, po)
    }
    
    po.UpdatedBy = approverID
    po.UpdatedAt = now
    
    return s.poRepo.Update(ctx, poID, po)
}

func (s *b2bService) RejectPO(ctx context.Context, poID primitive.ObjectID, approverID, reason string) (*model.PurchaseOrder, error) {
    po, _ := s.poRepo.GetByID(ctx, poID)
    
    // Find and mark this approval as rejected
    for i := range po.ApprovalChain {
        po.ApprovalChain[i].Status = "rejected"
        po.ApprovalChain[i].Notes = reason
        po.ApprovalChain[i].ApproverID = approverID
        po.ApprovalChain[i].DecisionAt = ptrTime(time.Now())
    }
    
    po.Status = model.POStatusRejected
    po.ApprovalStatus = model.ApprovalStatusRejected
    po.RejectedAt = ptrTime(time.Now())
    po.UpdatedBy = approverID
    po.UpdatedAt = time.Now()
    
    // Notify buyer: PO was rejected
    s.notifyPORejected(ctx, po, reason)
    
    return s.poRepo.Update(ctx, poID, po)
}

func (s *b2bService) ListPendingApprovalsFor(ctx context.Context, approverID, approverRole string) ([]model.PurchaseOrder, int64, error) {
    // Query: approval_chain.approver_role = approverRole AND approval_chain.status = "pending"
    col := s.db.Collection("purchase_orders")
    
    query := bson.M{
        "approval_chain": bson.M{
            "$elemMatch": bson.M{
                "approver_role": approverRole,
                "status":        "pending",
            },
        },
    }
    
    total, _ := col.CountDocuments(ctx, query)
    
    cursor, _ := col.Find(ctx, query)
    var pos []model.PurchaseOrder
    cursor.All(ctx, &pos)
    
    return pos, total, nil
}
```

#### Task 2.3: Notification System for Approvers (2h)

**File:** [internal/modules/b2b/service/po_notifications.go](../go/aystek-core/internal/modules/b2b/service/po_notifications.go)

```go
func (s *b2bService) notifyApprovers(ctx context.Context, po *model.PurchaseOrder) error {
    for _, approval := range po.ApprovalChain {
        // Send notification to approver
        // Method: Email / SMS / In-app
        
        notification := &Notification{
            RecipientID: approval.ApproverID,
            Type:        "po_approval_needed",
            Title:       fmt.Sprintf("PO %s awaits your approval", po.PONumber),
            Body:        fmt.Sprintf("Purchase order from %s (TRY %.2f) requires approval", 
                         po.CompanyID, po.GrandTotal),
            ActionURL:   fmt.Sprintf("/admin/b2b/pos/%s/approve", po.ID.Hex()),
            CreatedAt:   time.Now(),
        }
        
        s.notificationSvc.Send(ctx, notification)
    }
    
    return nil
}

func (s *b2bService) notifyPOApproved(ctx context.Context, po *model.PurchaseOrder) error {
    notification := &Notification{
        RecipientID: po.CreatedBy,
        Type:        "po_approved",
        Title:       fmt.Sprintf("PO %s approved", po.PONumber),
        Body:        "Your purchase order has been approved and is ready to convert to order",
        ActionURL:   fmt.Sprintf("/b2b/pos/%s", po.ID.Hex()),
        CreatedAt:   time.Now(),
    }
    
    return s.notificationSvc.Send(ctx, notification)
}
```

#### Task 2.4: Approval Chain Configuration (2h)

**File:** [internal/modules/b2b/model/approval_config.go](../go/aystek-core/internal/modules/b2b/model/approval_config.go)

```go
type ApprovalConfig struct {
    ID                primitive.ObjectID  `bson:"_id"`
    TenantID          string              `bson:"tenant_id"`
    
    // Thresholds per role
    Thresholds        []ApprovalThreshold `bson:"thresholds"`
    
    CreatedAt         time.Time           `bson:"created_at"`
    UpdatedAt         time.Time           `bson:"updated_at"`
}

type ApprovalThreshold struct {
    Role              string              // "store_manager", "credit_department", "executive"
    MinAmount         float64             // Minimum PO amount for this role to approve
    MaxAmount         *float64            // Maximum (null = unlimited)
    UserIDs           []string            // Which users have this role (for fallback)
    IsActive          bool                // Role enabled?
}
```

**Admin Endpoint:** [internal/modules/b2b/handler/approval_config_handler.go](../go/aystek-core/internal/modules/b2b/handler/approval_config_handler.go)

```go
// PUT /api/v1/b2b/approval-config
func (h *ConfigHandler) UpdateApprovalConfig(c *gin.Context) {
    input := struct {
        Thresholds []model.ApprovalThreshold `json:"thresholds"`
    }{}
    c.BindJSON(&input)
    
    config, err := h.configSvc.UpdateApprovalConfig(c, input.Thresholds)
    if err != nil {
        c.JSON(400, err.Error())
        return
    }
    
    c.JSON(200, config)
}
```

---

### Phase 3: Payment Terms Management (4h)

#### Task 3.1: Payment Terms Model (1h)

**File:** [internal/modules/b2b/model/payment_terms.go](../go/aystek-core/internal/modules/b2b/model/payment_terms.go)

```go
type PaymentTermsType string

const (
    TermsTypeImmediate      PaymentTermsType = "immediate"         // Pay now
    TermsTypeNet30          PaymentTermsType = "net_30"            // Pay within 30 days
    TermsTypeNet60          PaymentTermsType = "net_60"            // Pay within 60 days
    TermsTypeNet90          PaymentTermsType = "net_90"            // Pay within 90 days
    TermsType2Net10         PaymentTermsType = "2_10_net_30"       // 2% discount if paid in 10 days, otherwise 30 days
    TermsTypeDeposit        PaymentTermsType = "deposit"           // Upfront deposit
)

type PaymentTerms struct {
    ID                primitive.ObjectID  `bson:"_id"`
    TenantID          string              `bson:"tenant_id"`
    Name              string              `bson:"name"`             // e.g., "Net 30"
    Type              PaymentTermsType    `bson:"type"`
    
    // Durations
    PaymentDueDays    int                 `bson:"payment_due_days"`  // Days to pay
    EarlyDiscountDays *int                `bson:"early_discount_days,omitempty"`  // Days for discount
    EarlyDiscountPct  *float64            `bson:"early_discount_pct,omitempty"`   // e.g., 0.02 for 2%
    
    // Deposit terms
    DepositPercentage *float64            `bson:"deposit_percentage,omitempty"`  // e.g., 0.50 for 50%
    DepositDueDays    *int                `bson:"deposit_due_days,omitempty"`    // Days to pay deposit
    
    Description       string              `bson:"description,omitempty"`
    IsActive          bool                `bson:"is_active"`
    
    CreatedAt         time.Time           `bson:"created_at"`
    UpdatedAt         time.Time           `bson:"updated_at"`
}
```

#### Task 3.2: Payment Terms Service (1.5h)

**File:** [internal/modules/b2b/service/payment_terms.go](../go/aystek-core/internal/modules/b2b/service/payment_terms.go)

```go
type PaymentTermsService interface {
    GetTerms(ctx context.Context, id primitive.ObjectID) (*model.PaymentTerms, error)
    ListTerms(ctx context.Context) ([]model.PaymentTerms, error)
    CreateTerms(ctx context.Context, input CreatePaymentTermsInput) (*model.PaymentTerms, error)
    UpdateTerms(ctx context.Context, id primitive.ObjectID, input UpdatePaymentTermsInput) (*model.PaymentTerms, error)
    
    CalculateDueDate(ctx context.Context, termsID primitive.ObjectID, invoiceDate time.Time) (time.Time, error)
    CalculateEarlyDiscountDate(ctx context.Context, termsID primitive.ObjectID, invoiceDate time.Time) (*time.Time, error)
    CalculateDiscountAmount(ctx context.Context, termsID primitive.ObjectID, invoiceAmount float64) (*float64, error)
}

func (s *b2bService) CalculateDueDate(ctx context.Context, termsID primitive.ObjectID, invoiceDate time.Time) (time.Time, error) {
    terms, _ := s.termsRepo.GetByID(ctx, termsID)
    return invoiceDate.AddDate(0, 0, terms.PaymentDueDays), nil
}

func (s *b2bService) CalculateEarlyDiscountDate(ctx context.Context, termsID primitive.ObjectID, invoiceDate time.Time) (*time.Time, error) {
    terms, _ := s.termsRepo.GetByID(ctx, termsID)
    
    if terms.EarlyDiscountDays == nil {
        return nil, nil
    }
    
    discountDate := invoiceDate.AddDate(0, 0, *terms.EarlyDiscountDays)
    return &discountDate, nil
}

func (s *b2bService) CalculateDiscountAmount(ctx context.Context, termsID primitive.ObjectID, invoiceAmount float64) (*float64, error) {
    terms, _ := s.termsRepo.GetByID(ctx, termsID)
    
    if terms.EarlyDiscountPct == nil {
        return nil, nil
    }
    
    discount := invoiceAmount * *terms.EarlyDiscountPct
    return &discount, nil
}
```

#### Task 3.3: Invoice Integration (1.5h)

**File:** [internal/modules/billing/invoices/model/invoice.go](../go/aystek-core/internal/modules/billing/invoices/model/invoice.go) — UPDATE

```go
type Invoice struct {
    // ... existing fields ...
    
    // ✅ ADD: Payment terms
    PaymentTermsID      *primitive.ObjectID `bson:"payment_terms_id,omitempty"`
    PaymentTermsDesc    string              `bson:"payment_terms_desc,omitempty"`    // e.g., "Net 30"
    PaymentDueDate      *time.Time          `bson:"payment_due_date,omitempty"`
    EarlyDiscountDate   *time.Time          `bson:"early_discount_date,omitempty"`
    EarlyDiscountAmount *float64            `bson:"early_discount_amount,omitempty"`
    
    // Payment status (for B2B)
    PaymentStatus       string              `bson:"payment_status,omitempty"`        // "unpaid", "partially_paid", "paid", "overdue"
    PaidAmount          float64             `bson:"paid_amount,omitempty"`
    RemainingAmount     float64             `bson:"remaining_amount,omitempty"`
}
```

**Hook:** [internal/modules/orders/service/hooks_invoice.go](../go/aystek-core/internal/modules/orders/service/hooks_invoice.go) — UPDATE

```go
func (h *InvoiceHook) createInvoice(ctx context.Context, order *model.Order) error {
    // ... existing code ...
    
    // ✅ ADD: Payment terms for B2B
    var paymentTermsID *primitive.ObjectID
    var dueDate *time.Time
    
    if order.IsB2B {  // Check if order is from B2B company
        if order.PaymentTermsID != nil {
            paymentTermsID = order.PaymentTermsID
            
            // Calculate due date
            calcDue, _ := h.paymentTermsSvc.CalculateDueDate(ctx, *paymentTermsID, time.Now())
            dueDate = &calcDue
        }
    }
    
    invoiceInput := &invoices.CreateInvoiceInput{
        // ... existing fields ...
        PaymentTermsID:  paymentTermsID,
        PaymentDueDate:  dueDate,
    }
    
    return h.invoiceService.Create(ctx, invoiceInput)
}
```

---

### Phase 4: Testing & Deployment (4h)

#### Task 4.1: Unit Tests (2h)

```go
✅ TestCreatePO — PO creation, line items, tax calculation
✅ TestSubmitPO — Approval chain building
✅ TestApprovePO — Single approval, all approved detection
✅ TestRejectPO — Rejection flow
✅ TestConvertPOToOrder — PO → Order conversion, credit reservation
✅ TestPaymentTermsDueDate — Net 30, Net 60, 2/10 Net 30
✅ TestPaymentTermsDiscount — Early discount calculation
```

#### Task 4.2: Integration Tests (1h)

```
Scenario 1: Full PO → Order → Invoice flow
├─ Create PO
├─ Submit PO
├─ Manager approval
├─ Credit dept approval
├─ Convert to Order
├─ Verify credit reserved
├─ Verify invoice with Net 30 terms
└─ Result: ✅

Scenario 2: PO Rejection
├─ Create PO
├─ Submit PO
├─ Manager rejects
├─ Verify approval chain stops
├─ Verify notifications sent
└─ Result: ✅

Scenario 3: Credit Limit Exceeded
├─ Company has 100,000 TRY credit
├─ Create PO for 150,000 TRY
├─ Submit PO
├─ Verify approval chain includes credit dept
├─ Credit dept approves increase → Approval
└─ Result: ✅
```

#### Task 4.3: Deployment (1h)

- [ ] Database migration (new collections: purchase_orders, payment_terms)
- [ ] Indexes created (po_number unique, approval_chain queries)
- [ ] Feature flag for B2B workflow (default: enabled)
- [ ] Admin training: How to set approval thresholds
- [ ] Documentation: PO workflow user guide

---

## Effort Summary

| Task | Hours |
|------|-------|
| 1.1 PO Model | 2 |
| 1.2 CRUD Service | 2 |
| 1.3 PO → Order Conversion | 2 |
| 1.4 Repository | 2 |
| 2.1 Approval API | 2 |
| 2.2 Approval Logic | 2 |
| 2.3 Notifications | 2 |
| 2.4 Config | 2 |
| 3.1 Payment Terms Model | 1 |
| 3.2 Terms Service | 1.5 |
| 3.3 Invoice Integration | 1.5 |
| **Testing** | **2** |
| **Integration** | **1** |
| **Deployment** | **1** |
| **Total** | **24** |

---

## Summary

**P1-G2 (B2B Workflows)** is **24 hours** with:
- ✅ Complete PO model + CRUD
- ✅ Multi-level approval chain
- ✅ Payment terms management
- ✅ PO → Order conversion with credit reservation
- ✅ Comprehensive testing
- ✅ Production-ready deployment

**Integration Points:**
- P0-G1 (Order state machine) → PO converted to Order
- P1-G1 (Tax rules) → PO inherits reverse charge logic
- Order status hooks → Partial fulfillment tracking on PO
