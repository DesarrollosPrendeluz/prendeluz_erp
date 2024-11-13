package dtos

type FatherOrderWithRecount struct {
	ID            uint    `json:"id"`
	Code          string  `json:"code"`
	OrderStatusID uint    `json:"status_id"`
	OrderTypeID   uint    `json:"type_id"`
	Status        string  `json:"status"`
	Type          string  `json:"type"`
	TotalStock    float64 `json:"total_stock"`
	PendingStock  float64 `json:"pending_stock"`
}
