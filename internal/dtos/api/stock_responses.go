package dtos

type StockItem struct {
	StoreID  uint64 `json:"storeId"`
	ASIN     string `json:"asin"`
	Quantity int    `json:"quantity"`
}
