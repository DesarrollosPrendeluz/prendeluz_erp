package dtos

type Order struct {
	Status uint64 `json:"status"`
	Type   uint64 `json:"type"`
}

type OrderToUpdate struct {
	Id     uint64  `json:"id"`
	Status *uint64 `json:"status"`
	Type   *uint64 `json:"type"`
}
type LineToUpdate struct {
	Id              uint64  `json:"id"`
	ItemID          *uint64 `json:"item_id"`
	Quantity        *int64  `json:"quantity"`
	RecivedQuantity *int64  `json:"recived_quantity"`
	StoreID         *int64  `json:"store_id"`
	ClientID        *uint64 `json:"client_id"`
}

type Line struct {
	ItemID          uint64  `json:"item_id"`
	Quantity        int64   `json:"quantity"`
	RecivedQuantity int64   `json:"recived_quantity"`
	StoreID         int64   `json:"store_id"`
	ClientID        *uint64 `json:"client_id"`
}

type DataItem struct {
	Order Order  `json:"order"`
	Lines []Line `json:"lines"`
}

type OrderWithLinesRequest struct {
	Data []DataItem `json:"data"`
}

type OrdersToUpdatePartially struct {
	Data []OrderToUpdate `json:"data"`
}
type OrdersLinesToUpdatePartially struct {
	Data []LineToUpdate `json:"data"`
}
type Assign struct {
	ID int `json:"id"`
}
