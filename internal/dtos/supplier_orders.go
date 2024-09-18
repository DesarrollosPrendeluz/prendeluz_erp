package dtos

type SupplierOrders struct {
	OrderCode     uint64  `json:"order_code"`
	StockToBuy    int     `json:"stock_to_buy"`
	ItemSKU       string  `json:"item_sku"`
	ItemID        string  `json:"item_id"`
	FatherId      string  `json:"father_id"`
	Name          string  `json:"name"`
	EAN           string  `json:"ean"`
	SupplierName  string  `json:"supplier_name"`
	SupplierCode  string  `json:"supplier_sku"`
	SupplierPrice float64 `json:"supplier_price"`
}
