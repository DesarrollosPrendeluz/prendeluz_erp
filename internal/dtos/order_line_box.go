package dtos

type OrderLineBox struct {
	BoxID       uint64 `json:"boxId"`
	OrderLineID uint64 `json:"OrderLineId"`
	Quantity    int    `json:"quantity"`
}
type OrderLineBoxCreateReq struct {
	Data []OrderLineBox `json:"data"`
}

type OrderLineBoxToUpdate struct {
	Id          uint64  `json:"id"`
	BoxID       *uint64 `json:"boxId"`
	OrderLineID *uint64 `json:"OrderLineId"`
	Quantity    *int    `json:"quantity"`
}
type OrderLineBoxUpdateReq struct {
	Data []OrderLineBoxToUpdate `json:"data"`
}
