package services

type PalletBoxesOrderLinesService interface {
	CheckAndCreateBoxOrderLines(orderLineId int, palletNumber int, BoxNumber int, Quantity int) ([]string, []error)
}
