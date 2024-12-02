package dtos

type StoreLocation struct {
	StoreID uint64 `json:"storeId"`
	Code    string `json:"code"`
	Name    string `json:"name"`
}
type StoreLocationCreateReq struct {
	Data []StoreLocation `json:"data"`
}

type StoreLocationToUpdate struct {
	Id      uint64  `json:"id"`
	StoreID *uint64 `json:"storeId"`
	Code    *string `json:"code"`
	Name    *string `json:"name"`
}
type StoreLocationUpdateReq struct {
	Data []StoreLocationToUpdate `json:"data"`
}
