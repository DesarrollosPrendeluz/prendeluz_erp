package dtos

type ItemInfo struct {
	Id            uint64
	Sku           string
	Amount        int64
	RecivedAmount int64
	Name          string
	Ean           string
	Supplier      string
	AssignedUser  AssignedUserToOrderItem
	Locations     []string
	//FIXME:quitar una vez se haya externalizado
	Box    string
	Pallet string
}
type AssignedUserToOrderItem struct {
	AssignationId uint64 `json:"assignation_id"`
	UserId        uint64 `json:"user_id"`
	UserName      string `json:"user_name"`
}
type ItemsPerOrder struct {
	Id           uint64
	OrderCode    string
	StatusID     int64
	TypeID       int64
	Status       string
	Type         string
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
