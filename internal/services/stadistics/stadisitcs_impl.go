package services

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/repositories/erpupdateorderlinehistoryrepo"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/userrepo"
)

type StadisitcsImpl struct {
	userrepo                      userrepo.UserImpl
	fatherorderrepo               fatherorderrepo.FatherOrderImpl
	orderrepo                     orderrepo.OrderRepoImpl
	erpupdateorderlinehistoryrepo erpupdateorderlinehistoryrepo.ErpUpdateOrderLineHistoryImpl
	orderitemrepo                 orderitemrepo.OrderItemRepoImpl
	itemsrepo                     itemsrepo.ItemRepoImpl
}
type OriginalOrderLine struct {
	OrderLineID uint64
	Quantity    int64
}

func NewStadisitcService() *StadisitcsImpl {
	userrepo := *userrepo.NewUsersRepository(db.DB)
	fatherorderrepo := *fatherorderrepo.NewFatherOrderRepository(db.DB)
	erpupdateorderlinehistoryrepo := *erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	orderrepo := *orderrepo.NewOrderRepository(db.DB)
	itemsrepo := *itemsrepo.NewItemRepository(db.DB)

	return &StadisitcsImpl{
		fatherorderrepo:               fatherorderrepo,
		erpupdateorderlinehistoryrepo: erpupdateorderlinehistoryrepo,
		orderitemrepo:                 orderitemrepo,
		orderrepo:                     orderrepo,
		userrepo:                      userrepo,
		itemsrepo:                     itemsrepo,
	}
}

func (s *StadisitcsImpl) GetChangeStadistics(fatherCode string) dtos.HistoricStats {
	orderIdList := []uint64{}
	var fatherId uint64
	var returnData dtos.HistoricStats
	var data *dtos.OrderLinesStats
	if fatherCode != "" {
		fatherData, _ := s.fatherorderrepo.FindByCode(fatherCode)
		fatherId = fatherData.ID
		orders, _ := s.orderrepo.FindByFatherId(fatherId)
		for _, order := range orders {
			orderIdList = append(orderIdList, order.ID)
		}

		firstData, _ := getFirstStateOrderLines(orderIdList, fatherId, &returnData)
		data = &firstData
		codes, _ := s.erpupdateorderlinehistoryrepo.FindUpdateCodesByOrders(orderIdList)
		for _, v := range codes {
			historicData, _ := getHistoricLines(data, v.Code, &returnData)
			data = &historicData

		}

	}

	return returnData

}

func (s *StadisitcsImpl) GetRecivedStadistics(fatherCode string) dtos.RecivedHistory {
	orderIdList := []uint64{}
	var fatherId uint64
	var returnData dtos.RecivedHistory
	if fatherCode != "" {
		fatherData, _ := s.fatherorderrepo.FindByCode(fatherCode)
		fatherId = fatherData.ID
		orders, _ := s.orderrepo.FindByFatherId(fatherId)
		for _, order := range orders {
			orderIdList = append(orderIdList, order.ID)
		}
		fatherDataComplete, _, _ := s.fatherorderrepo.FindAllWithAssocData(1, 0, fatherCode, 0, 0)
		returnData.PickingProcess = dtos.ProcessTotalsAndPartials{
			TotalToProcess: int(fatherDataComplete[0].TotalPickingStock),
			NotProcessed:   int(fatherDataComplete[0].TotalPickingStock) - int(fatherDataComplete[0].TotalRecivedPickingQuantity),
		}
		returnData.StaggingProcess = dtos.ProcessTotalsAndPartials{
			TotalToProcess: int(fatherDataComplete[0].TotalStock),
			NotProcessed:   int(fatherDataComplete[0].TotalStock) - int(fatherDataComplete[0].PendingStock),
		}

		pickingData, _ := s.erpupdateorderlinehistoryrepo.FindDonePrecentByCode("1", orderIdList)
		returnData.UserPickingProcessed = returnRecivedWithUser(pickingData)
		satggingData, _ := s.erpupdateorderlinehistoryrepo.FindDonePrecentByCode("4", orderIdList)
		returnData.UserStaggingProcessed = returnRecivedWithUser(satggingData)

	}

	return returnData

}

func (s *StadisitcsImpl) GetOrderLineStadistics(fatherCode string) []dtos.LinePickingStats {
	var results []dtos.LinePickingStats

	fatherOrder, _ := s.fatherorderrepo.FindByCode(fatherCode)

	var ordersIDs []uint64
	fmt.Println(fatherOrder)
	for _, order := range *fatherOrder.ChildOrders {
		ordersIDs = append(ordersIDs, order.ID)
	}

	rawData := s.erpupdateorderlinehistoryrepo.FindAllByOrderId(ordersIDs)

	for _, itemData := range rawData {
		user, _ := s.userrepo.FindByID(uint64(itemData.UserID))
		item, _ := s.itemsrepo.FindByID(uint64(itemData.ItemID))
		dtoItem := dtos.LinePickingStats{Worker: user.Name, Item: *item.Name, Ean: item.EAN, Quantity: itemData.Quantity, CurrentTime: itemData.UpdatedAt}
		results = append(results, dtoItem)
	}

	return results

}
func returnRecivedWithUser(data []erpupdateorderlinehistoryrepo.Result) []dtos.UserProcessed {
	var returnData []dtos.UserProcessed
	for _, datum := range data {
		user, _ := userrepo.NewUsersRepository(db.DB).FindByID(datum.UserID)
		returns := dtos.UserProcessed{
			UserId:        datum.UserID,
			UserName:      user.Name,
			UserProcessed: datum.ModificationDif,
		}
		returnData = append(returnData, returns)

	}
	return returnData

}

func getFirstStateOrderLines(orderId []uint64, fatherId uint64, list *dtos.HistoricStats) (dtos.OrderLinesStats, error) {
	var linedata []dtos.OrderLineStat
	var orderIdList []uint64
	var total int = 0
	hisotricLineRepo := erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	orderLineMap := make(map[uint64]OriginalOrderLine)
	linesBase, _ := hisotricLineRepo.FindByOrders(orderId)

	for _, line := range linesBase {
		orderLineMap[line.OrderLineID] = OriginalOrderLine{
			OrderLineID: line.OrderLineID,
			Quantity:    line.Quantity,
		}
		orderIdList = append(orderIdList, line.OrderLineID)

	}

	orderItems, _ := orderItemRepo.FindIdWhereIn(orderIdList)
	orderData, _ := orderItemRepo.FindByOrderExludingIds(orderIdList, orderId)
	for _, orderItem := range orderItems {
		retrive := dtos.OrderLineStat{
			Line:            orderItem.ID,
			OrderID:         orderItem.OrderID,
			FatherId:        fatherId,
			Quantity:        int(orderLineMap[orderItem.ID].Quantity),
			RecivedQuantity: int(0),
		}
		linedata = append(linedata, retrive)
		total += int(orderLineMap[orderItem.ID].Quantity)

	}
	for _, orderItem := range orderData {
		retrive := dtos.OrderLineStat{
			Line:            orderItem.ID,
			OrderID:         orderItem.OrderID,
			FatherId:        fatherId,
			Quantity:        int(orderItem.Amount),
			RecivedQuantity: int(orderItem.RecivedAmount),
		}
		linedata = append(linedata, retrive)
		total += int(orderItem.Amount)

	}
	returnData := dtos.OrderLinesStats{
		TotaOrder: total,
		Code:      "Base",
		Lines:     linedata,
	}
	list.Results = append(list.Results, returnData)

	return returnData, nil

}

func getHistoricLines(data *dtos.OrderLinesStats, code string, list *dtos.HistoricStats) (dtos.OrderLinesStats, error) {
	var newData dtos.OrderLinesStats
	var total int = 0
	hisotricLineRepo := erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)
	modLines, _ := hisotricLineRepo.FindHistoryLinesByCode(code, []int{4, 2})
	for _, datum := range data.Lines {
		var orderLineStat dtos.OrderLineStat
		if value, exists := modLines[datum.Line]; exists {
			orderLineStat = dtos.OrderLineStat{
				FatherId: datum.FatherId,
				Line:     datum.Line,
				OrderID:  datum.OrderID,
				Quantity: value.NewQuantity}
			total += value.NewQuantity

		} else {
			orderLineStat = dtos.OrderLineStat{
				FatherId: datum.FatherId,
				Line:     datum.Line,
				OrderID:  datum.OrderID,
				Quantity: datum.Quantity}
			total += datum.Quantity
		}
		newData.Lines = append(newData.Lines, orderLineStat)

	}
	newData.TotaOrder = total
	newData.Code = code
	list.Results = append(list.Results, newData)
	return newData, nil
}
