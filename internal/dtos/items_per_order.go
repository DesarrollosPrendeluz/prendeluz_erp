package dtos

type ItemInfo struct {
	Sku    string
	Amount int64
}
type ItemsPerOrder struct {
	OrderCode    string
	ItemsOrdered []ItemInfo
}
