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
} //Esto no tiene sentido

type ItemStockLocationStockChange struct {
	Id    uint64 `json:"id"`
	Stock int    `json:"stock"`
}
type ItemStockLocationStockChangeRequest struct {
	Data []ItemStockLocationStockChange `json:"data"`
}

type ItemStockLocationStockMovement struct {
	MainSku          string `json:"productSku"`
	StoreLocationID1 uint64 `json:"beforeStoreLocationId"`
	StoreLocationID2 uint64 `json:"aftherStoreLocationId"`
	Stock            int    `json:"stock"`
}
type ItemStockLocationStockMovementRequest struct {
	Data []ItemStockLocationStockMovement `json:"data"`
}
