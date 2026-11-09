package dto

type PublishDuesRequest struct {
	PropertyID    string   `json:"property_id"`
	Period        string   `json:"period"`
	PlanID        *string  `json:"plan_id,omitempty"`
	DefaultAmount *float64 `json:"default_amount,omitempty"`
	Currency      string   `json:"currency,omitempty"`
}

type AllocatePaymentRequest struct {
	PaymentID  string   `json:"payment_id,omitempty"`
	Amount     *float64 `json:"amount,omitempty"`
	PayerRef   string   `json:"payer_ref,omitempty"`
	PropertyID string   `json:"property_id,omitempty"`
	Date       *string  `json:"date,omitempty"`
	Currency   string   `json:"currency,omitempty"`
}
