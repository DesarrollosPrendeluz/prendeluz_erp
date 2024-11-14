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
	ID              uint64 `json:"id"`
	Code            string `json:"code"`
	OrderStatusID   uint   `json:"status_id"`
	OrderTypeID     uint   `json:"type_id"`
	Status          string `json:"status"`
	Type            string `json:"type"`
	Quantity        uint64 `json:"quantity"`
	RecivedQuantity uint64 `json:"recived_quantity"`
	Childs          []ChildOrder
}
type ChildOrder struct {
	ID              uint64 `json:"id"`
	Code            string `json:"code"`
	OrderStatusID   uint   `json:"status_id"`
	Status          string `json:"status"`
	Quantity        uint64 `json:"quantity"`
	RecivedQuantity uint64 `json:"recived_quantity"`
}

type FatherOrderOrdersAndLines struct {
	FatherOrder FatherOrder
	Lines       []LinesInfo
}

type LinesInfo struct {
	LineID          uint     `json:"id"`
	OrderCode       uint64   `json:"order_id"`
	Quantity        int      `json:"quantity"`
	RecivedQuantity int      `json:"recived_quantity"`
	MainSku         string   `json:"main_sku"`
	Ean             string   `json:"ean"`
	Name            string   `json:"name"`
	SupplierName    string   `json:"supplier"`
	Location        []string `json:"locations"`
	AssignedUser    AssignedUserToOrderItem
}
