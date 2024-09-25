package dtos

type ItemInfo struct {
	Id            uint64
	Sku           string
	Amount        int64
	RecivedAmount int64
}
type ItemsPerOrder struct {
	Id           uint64
	OrderCode    string
	Status       int64
	Type         int64
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
