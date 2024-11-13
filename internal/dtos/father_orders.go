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
type FatherOrder struct {
	ID            uint64 `json:"id"`
	Code          string `json:"code"`
	OrderStatusID uint   `json:"status_id"`
	OrderTypeID   uint   `json:"type_id"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Childs        []ChildOrder
}
type ChildOrder struct {
	ID            uint64 `json:"id"`
	Code          string `json:"code"`
	OrderStatusID uint   `json:"status_id"`
	Status        string `json:"status"`
}

type FatherOrderOrdersAndLines struct {
	FatherOrder FatherOrder
	Lines       []LinesInfo
}

type LinesInfo struct {
	LineID          uint     `json:"id"`
	Quantity        int      `json:"quantity"`
	RecivedQuantity int      `json:"recived_quantity"`
	MainSku         string   `json:"main_sku"`
	Ean             string   `json:"ean"`
	Name            string   `json:"name"`
	SupplierName    string   `json:"supplier"`
	Location        []string `json:"locations"`
	AssignedUser    AssignedUserToOrderItem
}
