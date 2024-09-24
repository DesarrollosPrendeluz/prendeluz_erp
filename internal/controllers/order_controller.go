package controllers

import (
	"fmt"
	"log"
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/orderstatusrepo"
	"prendeluz/erp/internal/repositories/ordertyperepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
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
	orderRepo := orderrepo.NewOrderRepository(db.DB)
	orderService := services.NewOrderService()
	recount, _ := orderRepo.CountAll()
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

	c.IndentedJSON(http.StatusOK, gin.H{"data": results, "recount": recount})

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
		fmt.Println("pasa")
		if err := db.DB.Raw("CALL UpdateStockDeficitByStore();").Error; err != nil {
			log.Printf("Error ejecutando UpdateStockDeficitByStore: %v", err)
		}
		fmt.Println("pasa2")

		// Llamada al segundo procedimiento almacenado
		if err := db.DB.Raw("CALL UpdatePendingStocks();").Error; err != nil {
			log.Printf("Error ejecutando UpdatePendingStocks: %v", err)
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
			StoreID:       line.StoreID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Guardar cada línea en la base de datos
		if err := repo.Create(&orderLine); err != nil {

			return err
		}
		if order.OrderTypeID == uint64(2) && line.ClientID != nil {
			outRelRepo := outorderrelationrepo.NewOutOrderRelationRepository(db.DB)
			outRel := models.OutOrderRelation{
				ClientID:    *line.ClientID,
				OrderLineID: orderLine.ID,
			}
			outRelRepo.Create(&outRel)

		}
	}

	return nil

}

func EditOrders(c *gin.Context) {
	var requestBody dtos.OrdersToUpdatePartially
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	order := orderrepo.NewOrderRepository(db.DB)
	for _, dataItem := range requestBody.Data {
		model, err := order.FindByID(dataItem.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if dataItem.Status != nil {
			model.OrderStatusID = *dataItem.Status
		}
		if dataItem.Type != nil {
			model.OrderTypeID = *dataItem.Type
		}
		error := order.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Ok": "Orders are updated", "Errors": errorList})

}

func EditOrdersLines(c *gin.Context) {
	var requestBody dtos.OrdersLinesToUpdatePartially
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	orderLines := orderitemrepo.NewOrderItemRepository(db.DB)
	for _, dataItem := range requestBody.Data {
		model, err := orderLines.FindByID(dataItem.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if dataItem.ItemID != nil {
			model.ItemID = *dataItem.ItemID
		}
		if dataItem.RecivedQuantity != nil {
			model.RecivedAmount = *dataItem.RecivedQuantity
		}
		if dataItem.Quantity != nil {
			model.Amount = *dataItem.Quantity
		}
		if dataItem.StoreID != nil {
			model.Amount = *dataItem.StoreID
		}
		error := orderLines.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Ok": "Orders lines are updated", "Errors": errorList})

}
