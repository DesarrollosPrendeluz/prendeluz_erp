package dtos

type ItemInfo struct {
	Sku    string
	Amount int64
}
type ItemsPerOrder struct {
	OrderCode    string
	ItemsOrdered []ItemInfo
}
type ItemAssigantion struct {
	LineID uint64 `json:"line_id"`
	UserID uint64 `json:"user_id"`
}

type ItemsAssigantion struct {
	Assignations []ItemAssigantion
}

type ItemAssigantionEdit struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
}

type ItemsAssigantionEdit struct {
	Assignations []ItemAssigantionEdit
}
