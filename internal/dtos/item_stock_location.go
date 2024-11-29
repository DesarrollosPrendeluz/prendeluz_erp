package dtos

type ItemStockLocation struct {
	ItemMainSku     string `json:"itemMainSku"`
	StoreLocationID uint64 `json:"storeLocationId"`
	Stock           int    `json:"stock"`
}
type ItemStockLocationCreateReq struct {
	Data []ItemStockLocation `json:"data"`
}

//	Id     uint64  `json:"id"`

type ItemStockLocationToUpdate struct {
	Id              uint64  `json:"id"`
	ItemMainSku     *string `json:"itemMainSku"`
	StoreLocationID *uint64 `json:"storeLocationId"`
	Stock           *int    `json:"stock"`
}
type ItemStockLocationUpdateReq struct {
	Data []ItemStockLocationToUpdate `json:"data"`
}
