package services

import (
	"io"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/utils"

	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	orderRepo      repositories.OrderRepo
	orderItemsRepo repositories.OrderItemRepo
	orderErrorRepo repositories.ErrorOrderRepo
}

func NewOrderService() *OrderServiceImpl {

	orderRepo := repositories.NewOrderRepository(db.DB)
	errorOrderRepo := repositories.NewErrorOrderRepository(db.DB)
	orderItemRepo := repositories.NewOrderItemRepository(db.DB)
	return &OrderServiceImpl{orderRepo: *orderRepo, orderItemsRepo: *orderItemRepo, orderErrorRepo: *errorOrderRepo}
}

func itemsExist(rawOrders []utils.ExcelOrder, filename string) ([]models.Order, map[string][]models.OrderItem, []models.ErrorOrder) {
	itemRepo := repositories.NewItemRepository(db.DB)

	var errorOrdersList []models.ErrorOrder
	var ordersList []models.Order
	orderItemsOk := make(map[string][]models.OrderItem)

	for _, orderCode := range rawOrders {
		for _, orderInfo := range orderCode.Info {
			item, err := itemRepo.FindByMainSku(orderInfo.MainSku)

			if err != nil {
				errorOrder := models.ErrorOrder{
					Main_Sku: orderInfo.MainSku,
					Error:    "Item with sku not found",
					Order:    orderCode.OrderCode,
				}

				errorOrdersList = append(errorOrdersList, errorOrder)
			} else {
				orderItem := models.OrderItem{
					ItemID: item.ID,
					Amount: orderInfo.Amount,
				}
				orderItemsOk[orderCode.OrderCode] = append(orderItemsOk[orderCode.OrderCode], orderItem)
			}
		}

		order := models.Order{
			Orden_compra: orderCode.OrderCode,
			Filename:     filename,
		}
		ordersList = append(ordersList, order)

	}
	return ordersList, orderItemsOk, errorOrdersList
}
func (s *OrderServiceImpl) UploadOrderExcel(file io.Reader, filename string) error {

	excelOrderList, err := utils.ExceltoJSON(file)
	if err != nil {
		return err
	}
	succesOrders, orderItems, errorOrders := itemsExist(excelOrderList, filename)

	return db.DB.Transaction(func(tx *gorm.DB) error {

		orderRepo := repositories.NewOrderRepository(tx)
		errorOrderRepo := repositories.NewErrorOrderRepository(tx)
		orderItemRepo := repositories.NewOrderItemRepository(tx)

		_, err = orderRepo.CreateAll(&succesOrders)
		if err != nil {
			return err
		}

		if len(errorOrders) > 0 {
			_, err = errorOrderRepo.CreateAll(&errorOrders)
			if err != nil {
				return err
			}
		}
		var orderItemToInsert []models.OrderItem
		for _, order := range succesOrders {
			for _, tmp := range orderItems[order.Orden_compra] {
				tmp.OrderID = order.ID
				orderItemToInsert = append(orderItemToInsert, tmp)
			}
		}
		_, err = orderItemRepo.CreateAll(&orderItemToInsert)
		if err != nil {
			return err
		}
		return nil

	})
}

func (s *OrderServiceImpl) GetOrders() ([]dtos.ItemsPerOrder, error) {

	var results []dtos.ItemsPerOrder

	orderRepo := repositories.NewOrderRepository(db.DB)
	orderItemRepo := repositories.NewOrderItemRepository(db.DB)
	itemRepo := repositories.NewItemRepository(db.DB)

	orders, err := orderRepo.FindAll()
	if err != nil {
		return results, err
	}

	for _, order := range orders {
		var itemOrder dtos.ItemsPerOrder
		itemOrder.OrderCode = order.Orden_compra
		orderItemList, _ := orderItemRepo.FindByOrder(order.ID)

		for _, orderItem := range orderItemList {
			var itemInfo dtos.ItemInfo
			item, _ := itemRepo.FindByID(orderItem.ItemID)
			itemInfo.Amount = orderItem.Amount
			itemInfo.Sku = item.MainSKU
			itemOrder.ItemsOrdered = append(itemOrder.ItemsOrdered, itemInfo)
		}
		results = append(results, itemOrder)

	}

	return results, nil

}
