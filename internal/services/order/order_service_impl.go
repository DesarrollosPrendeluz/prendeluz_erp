package services

import (
	"fmt"
	"io"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/ordererrorrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
	"prendeluz/erp/internal/utils"
)

type OrderServiceImpl struct {
	orderRepo      orderrepo.OrderRepoImpl
	orderItemsRepo orderitemrepo.OrderItemRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewOrderService() *OrderServiceImpl {
	orderRepo := *orderrepo.NewOrderRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &OrderServiceImpl{orderRepo: orderRepo, orderItemsRepo: orderItemRepo,
		orderErrorRepo: errorOrderRepo, itemsRepo: itemsRepo}
}

// retorna datos para crear ordenes las línea de las ordenes y los errores correspondientes
func generateOrdersAndOrderLines(rawOrders []utils.ExcelOrder, filename string) ([]models.Order, map[string][]models.OrderItem, []models.ErrorOrder) {
	itemsRepo := itemsrepo.NewItemRepository(db.DB)

	var errorOrdersList []models.ErrorOrder
	var ordersList []models.Order
	orderItemsOk := make(map[string][]models.OrderItem)

	for _, orderCode := range rawOrders {
		for _, orderInfo := range orderCode.Info {
			item, err := itemsRepo.FindByMainSku(orderInfo.MainSku)

			if err != nil {
				errorOrder := models.ErrorOrder{
					Main_Sku: orderInfo.MainSku,
					Error:    "Item with sku " + orderInfo.MainSku + " not found",
					Order:    orderCode.OrderCode,
				}

				errorOrdersList = append(errorOrdersList, errorOrder)
			} else {
				orderItem := models.OrderItem{
					ItemID:        item.ID,
					Amount:        orderInfo.Amount,
					RecivedAmount: 0,
					StoreID:       int64(orderInfo.Store),
					ClientID:      orderInfo.Client,
				}
				orderItemsOk[orderCode.OrderCode] = append(orderItemsOk[orderCode.OrderCode], orderItem)
			}
		}

		order := models.Order{
			OrderStatusID: uint64(orderrepo.Order_Status["iniciada"]),
			OrderTypeID:   uint64(orderrepo.Order_Types["venta"]),
			Code:          orderCode.OrderCode,
			Filename:      filename,
		}
		ordersList = append(ordersList, order)

	}
	return ordersList, orderItemsOk, errorOrdersList
}

// Carga el excel y crea las nuevas ordenes en este caso solo de ventas por el momento
func (s *OrderServiceImpl) UploadOrderExcel(file io.Reader, filename string) error {
	fmt.Println("entra")
	excelOrderList, err := utils.ExceltoJSON(file)

	if err != nil {
		return err
	}
	succesOrders, orderItems, errorOrders := generateOrdersAndOrderLines(excelOrderList, filename)
	orderRepo := orderrepo.NewOrderRepository(db.DB)
	orderItem := orderitemrepo.NewOrderItemRepository(db.DB)
	out := outorderrelationrepo.NewOutOrderRelationRepository(db.DB)
	orderErr := ordererrorrepo.NewOrderErrRepository(db.DB)

	_, errs := orderRepo.CreateAll(&succesOrders)
	if errs != nil {
		return err
	}

	if len(errorOrders) > 0 {
		_, err := orderErr.CreateAll(&errorOrders)
		if err != nil {
			return err
		}
	}
	var orderItemToInsert []models.OrderItem
	for _, order := range succesOrders {
		for _, tmp := range orderItems[order.Code] {
			tmp.OrderID = order.ID
			orderItemToInsert = append(orderItemToInsert, tmp)
		}
	}
	for _, orderLine := range orderItemToInsert {
		err = orderItem.Create(&orderLine)
		out.Create(&models.OutOrderRelation{
			ClientID:    orderLine.ClientID,
			OrderLineID: orderLine.ID,
		})

	}

	if err != nil {
		return err
	}
	return nil

}

// Obtiene las ordenes paginadas en base a los parámetros page y pagesize
// A su vez si recibe los parámetros start dat y end date se filtran dichas ordenes por fecha de creación
func (s *OrderServiceImpl) GetOrders(page int, pageSize int, startDate string, endDate string) ([]dtos.ItemsPerOrder, error) {

	var results []dtos.ItemsPerOrder
	var orders []models.Order
	var err error
	offset := (page - 1) * pageSize
	if startDate != "" || endDate != "" {
		orders, err = s.orderRepo.FindOrderByDate(startDate, endDate)
	} else {
		orders, err = s.orderRepo.FindAll(pageSize, offset)
	}

	if err != nil {
		return results, err
	}

	for _, order := range orders {
		var itemOrder dtos.ItemsPerOrder
		itemOrder.OrderCode = order.Code
		orderItemList, _ := s.orderItemsRepo.FindByOrder(order.ID)

		for _, orderItem := range orderItemList {
			var itemInfo dtos.ItemInfo
			item, _ := s.itemsRepo.FindByID(orderItem.ItemID)
			itemInfo.Amount = orderItem.Amount
			itemInfo.Sku = item.MainSKU
			itemOrder.ItemsOrdered = append(itemOrder.ItemsOrdered, itemInfo)
		}
		results = append(results, itemOrder)

	}

	return results, nil

}

// Actualiza el estado de una orden a completada
func (s *OrderServiceImpl) OrderComplete(orderCode string) error {
	order, err := s.orderRepo.FindByOrderCode(orderCode)
	if err != nil {
		return err
	}

	s.orderRepo.UpdateStatus(orderrepo.Order_Status["finalizada"], order.ID)

	return nil
}
