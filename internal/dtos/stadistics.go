package dtos

type HistoricStats struct {
	Results []OrderLinesStats `json:"results"`
}

type OrderLinesStats struct {
	TotaOrder int             `json:"totalOrder"`
	Code      string          `json:"code"`
	Lines     []OrderLineStat `json:"lines"`
}

type OrderLineStat struct {
	Line            uint64 `json:"lineId"`
	OrderID         uint64 `json:"orderId"`
	FatherId        uint64 `json:"fatherId"`
	Quantity        int    `json:"quantity"`
	RecivedQuantity int    `json:"recivedQuantity"`
}

type RecivedHistory struct {
	PickingProcess        ProcessTotalsAndPartials
	StaggingProcess       ProcessTotalsAndPartials
	UserPickingProcessed  []UserProcessed
	UserStaggingProcessed []UserProcessed
}
type UserProcessed struct {
	UserId        uint64
	UserName      string
	UserProcessed int
}
type ProcessTotalsAndPartials struct {
	TotalToProcess int
	NotProcessed   int
}
