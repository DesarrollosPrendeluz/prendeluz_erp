package dtos

type OrderLineStat struct {
	Line            uint64 `json:"lineId"`
	OrderID         uint64 `json:"orderId"`
	FatherId        uint64 `json:"fatherId"`
	Quantity        int    `json:"quantity"`
	RecivedQuantity int    `json:"recivedQuantity"`
}
