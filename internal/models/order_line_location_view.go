package models

type OrderLineLocationView struct {
	FatherOrderID     uint   `json:"father_order_id" gorm:"column:father_order_id"`
	OrderID           uint   `json:"order_id" gorm:"column:order_id"`
	OrderLineID       uint   `json:"order_line_id" gorm:"column:order_line_id"`
	OrderLineItemID   uint   `json:"order_line_item_id" gorm:"column:order_line_item_id"`
	StoreID           uint   `json:"store_id" gorm:"column:store_id"`
	ItemSKU           string `json:"item_sku" gorm:"column:item_sku"`
	OrderLineItemEAN  string `json:"order_line_item_ean" gorm:"column:order_line_item_ean"`
	OrderLineItemType string `json:"order_line_item_type" gorm:"column:order_line_item_type"`
	FatherID          uint   `json:"father_id" gorm:"column:father_id"`
	FatherSKU         string `json:"father_sku" gorm:"column:father_sku"`
	StoreLocationCode string `json:"store_location_code" gorm:"column:store_location_code"`
	StoreLocationID   *uint  `json:"store_location_id" gorm:"column:store_location_id"`

	AssignedRel AssignedLine `gorm:"foreignKey:OrderLineID;references:OrderLineID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FatherItem  Item         `gorm:"foreignKey:ID;references:FatherID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la vista en la base de datos
func (OrderLineLocationView) TableName() string {
	return "order_lines_locations" // Cambia esto al nombre real de la vista
}
