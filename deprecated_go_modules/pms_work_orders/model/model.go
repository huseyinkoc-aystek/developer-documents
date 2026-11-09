package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkOrderCategory string

type WorkOrderPriority string

type WorkOrderStatus string

type WorkOrderEventType string

const (
	CategoryPlumbing WorkOrderCategory = "plumbing"
	CategoryElectric WorkOrderCategory = "electric"
	CategorySecurity WorkOrderCategory = "security"
	CategoryCleaning WorkOrderCategory = "cleaning"
	CategoryOther    WorkOrderCategory = "other"
)

const (
	PriorityLow    WorkOrderPriority = "low"
	PriorityNormal WorkOrderPriority = "normal"
	PriorityHigh   WorkOrderPriority = "high"
	PriorityUrgent WorkOrderPriority = "urgent"
)

const (
	StatusOpen       WorkOrderStatus = "open"
	StatusInProgress WorkOrderStatus = "in_progress"
	StatusBlocked    WorkOrderStatus = "blocked"
	StatusDone       WorkOrderStatus = "done"
	StatusCancelled  WorkOrderStatus = "cancelled"
)

const (
	EventCreated       WorkOrderEventType = "created"
	EventAssigned      WorkOrderEventType = "assigned"
	EventStatusChanged WorkOrderEventType = "status_changed"
	EventCommented     WorkOrderEventType = "commented"
	EventClosed        WorkOrderEventType = "closed"
	EventReopened      WorkOrderEventType = "reopened"
)

func (v WorkOrderCategory) IsValid() bool {
	switch v {
	case CategoryPlumbing, CategoryElectric, CategorySecurity, CategoryCleaning, CategoryOther:
		return true
	default:
		return false
	}
}

func (v WorkOrderPriority) IsValid() bool {
	switch v {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

func (v WorkOrderStatus) IsValid() bool {
	switch v {
	case StatusOpen, StatusInProgress, StatusBlocked, StatusDone, StatusCancelled:
		return true
	default:
		return false
	}
}

type Attachment struct {
	MediaID   string `bson:"media_id,omitempty" json:"media_id,omitempty"`
	FileName  string `bson:"file_name,omitempty" json:"file_name,omitempty"`
	MimeType  string `bson:"mime_type,omitempty" json:"mime_type,omitempty"`
	SizeBytes int64  `bson:"size_bytes,omitempty" json:"size_bytes,omitempty"`
}

type CostItem struct {
	Amount    float64 `bson:"amount" json:"amount"`
	Currency  string  `bson:"currency" json:"currency"`
	Note      string  `bson:"note,omitempty" json:"note,omitempty"`
	VendorRef string  `bson:"vendor_ref,omitempty" json:"vendor_ref,omitempty"`
}

type WorkOrder struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID       string              `bson:"tenant_id" json:"tenant_id"`
	SiteID         string              `bson:"site_id" json:"site_id"`
	PropertyID     primitive.ObjectID  `bson:"property_id" json:"property_id"`
	UnitID         *primitive.ObjectID `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	Title          string              `bson:"title" json:"title"`
	Description    string              `bson:"description,omitempty" json:"description,omitempty"`
	Category       WorkOrderCategory   `bson:"category" json:"category"`
	Priority       WorkOrderPriority   `bson:"priority" json:"priority"`
	Status         WorkOrderStatus     `bson:"status" json:"status"`
	SLADueAt       *time.Time          `bson:"sla_due_at,omitempty" json:"sla_due_at,omitempty"`
	AssigneeUserID *primitive.ObjectID `bson:"assignee_user_id,omitempty" json:"assignee_user_id,omitempty"`
	ReporterRef    string              `bson:"reporter_ref,omitempty" json:"reporter_ref,omitempty"`
	Tags           []string            `bson:"tags,omitempty" json:"tags,omitempty"`
	Attachments    []Attachment        `bson:"attachments,omitempty" json:"attachments,omitempty"`
	CostItems      []CostItem          `bson:"cost_items,omitempty" json:"cost_items,omitempty"`
	ClosedAt       *time.Time          `bson:"closed_at,omitempty" json:"closed_at,omitempty"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at" json:"updated_at"`
	CreatedBy      *primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy      *primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	DeletedAt      *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	DeletedBy      *primitive.ObjectID `bson:"deleted_by,omitempty" json:"deleted_by,omitempty"`
}

type WorkOrderEvent struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	TenantID    string                 `bson:"tenant_id" json:"tenant_id"`
	SiteID      string                 `bson:"site_id" json:"site_id"`
	WorkOrderID primitive.ObjectID     `bson:"work_order_id" json:"work_order_id"`
	Type        WorkOrderEventType     `bson:"type" json:"type"`
	Payload     map[string]interface{} `bson:"payload,omitempty" json:"payload,omitempty"`
	CreatedBy   *primitive.ObjectID    `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
}

type OutboxEvent struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty"`
	TenantID   string                 `bson:"tenant_id"`
	SiteID     string                 `bson:"site_id"`
	EventName  string                 `bson:"event_name"`
	EntityID   primitive.ObjectID     `bson:"entity_id"`
	ActorID    *primitive.ObjectID    `bson:"actor_id,omitempty"`
	OccurredAt time.Time              `bson:"occurred_at"`
	Version    int                    `bson:"version"`
	Diff       map[string]interface{} `bson:"diff,omitempty"`
	Payload    map[string]interface{} `bson:"payload,omitempty"`
	Processed  bool                   `bson:"processed"`
	CreatedAt  time.Time              `bson:"created_at"`
}

type WorkOrderDetail struct {
	WorkOrder WorkOrder        `json:"work_order"`
	Events    []WorkOrderEvent `json:"events"`
}
