package services

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/repositories/erpupdateorderlinehistoryrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
)

type StadisitcsImpl struct {
	orderrepo                     orderrepo.OrderRepoImpl
	erpupdateorderlinehistoryrepo erpupdateorderlinehistoryrepo.ErpUpdateOrderLineHistoryImpl
	orderitemrepo                 orderitemrepo.OrderItemRepoImpl
}
type OriginalOrderLine struct {
	OrderLineID uint64
	Quantity    int64
}

func NewStadisitcService() *StadisitcsImpl {
	erpupdateorderlinehistoryrepo := *erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	orderrepo := *orderrepo.NewOrderRepository(db.DB)

	return &StadisitcsImpl{
		erpupdateorderlinehistoryrepo: erpupdateorderlinehistoryrepo,
		orderitemrepo:                 orderitemrepo,
		orderrepo:                     orderrepo,
	}
}

func (s *StadisitcsImpl) GetChangeStadistics(fatherId uint64) []dtos.OrderLineStat {
	orderIdList := []uint64{}

	orders, _ := s.orderrepo.FindByFatherId(fatherId)
	for _, order := range orders {
		orderIdList = append(orderIdList, order.ID)
	}
	fmt.Println(orderIdList)

	data, _ := getFirstStateOrderLines(orderIdList, fatherId)
	return data

}

func getFirstStateOrderLines(orderId []uint64, fatherId uint64) ([]dtos.OrderLineStat, error) {
	var linedata []dtos.OrderLineStat
	var orderIdList []uint64
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

	}

	return linedata, nil

}
