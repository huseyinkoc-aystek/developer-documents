package dto

type ListQuery struct {
	PropertyID string `form:"property_id"`
	Q          string `form:"q"`
	Page       int64  `form:"page"`
	Limit      int64  `form:"limit"`
}
