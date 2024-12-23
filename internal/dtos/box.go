package dtos

type Box struct {
	PalletID uint64 `json:"palletId"`
	Number   uint64 `json:"number"`
	Label    string `json:"label"`
	Quantity int    `json:"quantity"`
}
type BoxCreateReq struct {
	Data []Box `json:"data"`
}

type BoxToUpdate struct {
	Id       uint64  `json:"id"`
	PalletID *uint64 `json:"palletId"`
	Number   *uint64 `json:"number"`
	Label    *string `json:"label"`
	Quantity *int    `json:"quantity"`
}
type BoxUpdateReq struct {
	Data []BoxToUpdate `json:"data"`
}
