package dtos

type Pallet struct {
	OrderID uint64 `json:"orderId"`
	Number  uint64 `json:"number"`
	Label   string `json:"label"`
}
type PalletCreateReq struct {
	Data []Pallet `json:"data"`
}

type PalletToUpdate struct {
	Id      uint64  `json:"id"`
	OrderID *uint64 `json:"orderId"`
	Number  *int    `json:"number"`
	Label   *string `json:"label"`
	IsClose *bool   `json:"is_close"`
}
type PalletUpdateReq struct {
	Data []PalletToUpdate `json:"data"`
}
