package controllers

import (
	"log"
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/orderstatusrepo"
	"prendeluz/erp/internal/repositories/ordertyperepo"
	services "prendeluz/erp/internal/services/order"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	orders, err := orderService.GetOrders(page, pageSize, startDate, endDate)

	for _, order := range orders {
		results[order.OrderCode] = order.ItemsOrdered
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": results})
	return

}

func GetOrderTypes(c *gin.Context) {
	repo := orderstatusrepo.NewOrderStatusRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	results, err := repo.FindAll(pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return

	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": results})
	return

}

func GetOrderStatus(c *gin.Context) {
	repo := ordertyperepo.NewOrderTypeRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	results, err := repo.FindAll(pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return

	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": results})
	return

}

func CreateOrder(c *gin.Context) {
	var requestBody dtos.OrderWithLinesRequest

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {
		order := dataItem.Order
		lines := dataItem.Lines
		repo := orderrepo.NewOrderRepository(db.DB)
		orderObject := models.Order{
			OrderStatusID: order.Status,
			OrderTypeID:   order.Type,
			Code:          "request.generated",
			Filename:      "request",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if repo.Create(&orderObject) == nil {
			createOrderLines(orderObject, lines)

		}

	}
}

func createOrderLines(order models.Order, lines []dtos.Line) error {
	repo := orderitemrepo.NewOrderItemRepository(db.DB) // Asumiendo que tienes un repositorio para las líneas

	for _, line := range lines {
		orderLine := models.OrderItem{
			OrderID:       order.ID,
			ItemID:        line.ItemID,
			Amount:        line.Quantity,
			RecivedAmount: line.RecivedQuantity,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Guardar cada línea en la base de datos
		if err := repo.Create(&orderLine); err != nil {

			return err
		}
	}

	return nil

}
