package dtos

import "prendeluz/erp/internal/models"

type FatherOrderWithRecount struct {
	ID                          uint    `json:"id"`
	Code                        string  `json:"code"`
	OrderStatusID               uint    `json:"status_id"`
	OrderTypeID                 uint    `json:"type_id"`
	Status                      string  `json:"status"`
	Type                        string  `json:"type"`
	TotalStock                  float64 `json:"total_stock"`
	PendingStock                float64 `json:"pending_stock"`
	TotalPickingStock           float64 `json:"total_picking_stock"`
	TotalRecivedPickingQuantity float64 `json:"total_recived_picking_quantity"`
}
type FatherOrder struct {
	ID              uint64                `json:"id"`
	Code            string                `json:"code"`
	OrderStatusID   uint                  `json:"status_id"`
	OrderTypeID     uint                  `json:"type_id"`
	Status          string                `json:"status"`
	Type            string                `json:"type"`
	Quantity        uint64                `json:"quantity"`
	RecivedQuantity uint64                `json:"recived_quantity"`
	GenericSupplier *models.SupplierOrder `json:"supplier_order"`
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
	FatherMainSku   string   `json:"father_main_sku"`
	Ean             string   `json:"ean"`
	Name            string   `json:"name"`
	SupplierName    string   `json:"supplier"`
	Pallet          *string  `json:"pallet"`
	Box             *string  `json:"box"`
	SupplierRef     string   `json:"supplier_reference"`
	Location        []string `json:"locations"`
	AssignedUser    AssignedUserToOrderItem
}

type FatherOrderId struct {
	FatherOrderId int64 `json:"fatherOrderId"`
}
