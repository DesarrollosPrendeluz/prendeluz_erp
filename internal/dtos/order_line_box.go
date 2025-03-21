package dtos

type OrderLineBox struct {
	BoxID       uint64 `json:"boxId"`
	OrderLineID uint64 `json:"orderLineId"`
	Quantity    int    `json:"quantity"`
}
type OrderLineBoxCreateReq struct {
	Data []OrderLineBox `json:"data"`
}

type OrderLineBoxToUpdate struct {
	Id          uint64  `json:"id"`
	BoxID       *uint64 `json:"boxId"`
	OrderLineID *uint64 `json:"orderLineId"`
	Quantity    *int    `json:"quantity"`
}
type OrderLineBoxUpdateReq struct {
	Data []OrderLineBoxToUpdate `json:"data"`
}

type OrderLineBoxProcessed struct {
	Box         int `json:"boxNumber"`
	Pallet      int `json:"palletNumber"`
	OrderLineID int `json:"orderLineId"`
	Quantity    int `json:"quantity"`
}
type OrderLineBoxProcessedCreateReq struct {
	Data []OrderLineBoxProcessed `json:"data"`
}
