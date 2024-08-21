package controllers

import (
	"log"
	"net/http"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/services/order"

	"github.com/gin-gonic/gin"
)

// func ItemsExist(rawOrders []utils.ExcelOrder, filename string) ([]models.Order, map[string][]models.OrderItem, []models.ErrorOrder) {
// 	itemRepo := repositories.NewItemRepository(db.DB)
//
// 	var errorOrdersList []models.ErrorOrder
// 	var ordersList []models.Order
// 	orderItemsOk := make(map[string][]models.OrderItem)
//
// 	for _, orderCode := range rawOrders {
// 		for _, orderInfo := range orderCode.Info {
// 			item, err := itemRepo.FindByMainSku(orderInfo.MainSku)
//
// 			if err != nil {
// 				errorOrder := models.ErrorOrder{
// 					Main_Sku: orderInfo.MainSku,
// 					Error:    "Item with sku not found",
// 					Order:    orderCode.OrderCode,
// 				}
//
// 				errorOrdersList = append(errorOrdersList, errorOrder)
// 			} else {
// 				orderItem := models.OrderItem{
// 					ItemID: item.ID,
// 					Amount: orderInfo.Amount,
// 				}
// 				orderItemsOk[orderCode.OrderCode] = append(orderItemsOk[orderCode.OrderCode], orderItem)
// 			}
// 		}
//
// 		order := models.Order{
// 			Orden_compra: orderCode.OrderCode,
// 			Filename:     filename,
// 		}
// 		ordersList = append(ordersList, order)
//
// 	}
//
// 	return ordersList, orderItemsOk, errorOrdersList
// }

func AddOrder(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	serviceOrder := services.NewOrderService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Fatal(err)
		return

	}

	serviceOrder.UploadOrderExcel(file, header.Filename)

	c.JSON(http.StatusCreated, gin.H{"message": "Upload succesfully"})

}

func GetOrders(c *gin.Context) {
	results := make(map[string][]dtos.ItemInfo)

	orderService := services.NewOrderService()
	orders, err := orderService.GetOrders()
	for _, order := range orders {
		results[order.OrderCode] = order.ItemsOrdered
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
	return

}
