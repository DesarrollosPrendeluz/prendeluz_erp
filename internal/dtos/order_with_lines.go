package dtos

type Order struct {
	Status uint64 `json:"status"`
	Type   uint64 `json:"type"`
}

type Line struct {
	ItemID          uint64 `json:"item_id"`
	Quantity        int64  `json:"quantity"`
	RecivedQuantity int64  `json:"recived_quantity"`
}

type DataItem struct {
	Order Order  `json:"order"`
	Lines []Line `json:"lines"`
}

type OrderWithLinesRequest struct {
	Data []DataItem `json:"data"`
}
