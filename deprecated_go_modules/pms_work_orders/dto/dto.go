package dto

type CostItem struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Note      string  `json:"note,omitempty"`
	VendorRef string  `json:"vendor_ref,omitempty"`
}

type Attachment struct {
	MediaID   string `json:"media_id,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	MimeType  string `json:"mime_type,omitempty"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

type CreateRequest struct {
	PropertyID  string       `json:"property_id"`
	UnitID      *string      `json:"unit_id,omitempty"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Priority    string       `json:"priority"`
	SLADueAt    *string      `json:"sla_due_at,omitempty"`
	ReporterRef string       `json:"reporter_ref,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	CostItems   []CostItem   `json:"cost_items,omitempty"`
}

type UpdateRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Priority    string       `json:"priority"`
	Tags        []string     `json:"tags,omitempty"`
	SLADueAt    *string      `json:"sla_due_at,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	CostItems   []CostItem   `json:"cost_items,omitempty"`
}

type AssignRequest struct {
	AssigneeUserID string `json:"assignee_user_id"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

type CommentRequest struct {
	Comment string `json:"comment"`
}
