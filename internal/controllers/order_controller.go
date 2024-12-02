package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/orderstatusrepo"
	"prendeluz/erp/internal/repositories/ordertyperepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
	"prendeluz/erp/internal/repositories/tokenrepo"
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
	if err := db.DB.Exec("CALL UpdateStockDeficitByStore();").Error; err != nil {
		log.Printf("Error ejecutando UpdateStockDeficitByStore: %v", err)
	}

	// Llamada al segundo procedimiento almacenado
	if err := db.DB.Exec("CALL UpdatePendingStocks();").Error; err != nil {
		log.Printf("Error ejecutando UpdatePendingStocks: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{"Results": gin.H{"Ok": "Upload succesfully"}})

}

func GetOrders(c *gin.Context) {

	orderService := services.NewOrderService()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	statusType, _ := strconv.Atoi(c.DefaultQuery("status_id", "0"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	code := c.Query("order_code")

	orders, recount, err := orderService.GetOrders(page, pageSize, startDate, endDate, statusType, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Results": gin.H{"error": err}})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": orders, "recount": recount}})

}

func GetOrderTypes(c *gin.Context) {
	repo := ordertyperepo.NewOrderTypeRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	results, err := repo.FindAll(pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Results": gin.H{"error": err}})
		return

	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results}})

}

func GetOrderStatus(c *gin.Context) {
	repo := orderstatusrepo.NewOrderStatusRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	results, err := repo.FindAll(pageSize, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Results": gin.H{"error": err}})
		return

	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results}})

}

func CreateOrder(c *gin.Context) {
	var requestBody dtos.OrderWithLinesRequest

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}
	fechaActual := time.Now().Format("2006-01-02 15:04:05")

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {
		order := dataItem.Order
		lines := dataItem.Lines
		fatherRepo := fatherorderrepo.NewFatherOrderRepository(db.DB)
		repo := orderrepo.NewOrderRepository(db.DB)
		fatherObject := models.FatherOrder{
			OrderStatusID: order.Status,
			OrderTypeID:   order.Type,
			Code:          "request.generated." + fechaActual,
			Filename:      "request",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if fatherRepo.Create(&fatherObject) == nil {
			orderObject := models.Order{
				OrderStatusID: order.Status,
				FatherOrderID: fatherObject.ID,
				Code:          "request.generated." + fechaActual,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if repo.Create(&orderObject) == nil {
				createOrderLines(fatherObject, orderObject, lines)

			}

		}

		if err := db.DB.Exec("CALL UpdateStockDeficitByStore();").Error; err != nil {
			log.Printf("Error ejecutando UpdateStockDeficitByStore: %v", err)
		} else {
			fmt.Println("en teoría se ha ejecutado: CALL UpdateStockDeficitByStore();")

		}

		// Llamada al segundo procedimiento almacenado
		if err := db.DB.Exec("CALL UpdatePendingStocks();").Error; err != nil {
			log.Printf("Error ejecutando UpdatePendingStocks: %v", err)
		} else {
			fmt.Println("en teoría se ha ejecutado: CALL UpdatePendingStocks()")

		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders are created"}})
}

func createOrderLines(fatherOrder models.FatherOrder, order models.Order, lines []dtos.Line) error {
	repo := orderitemrepo.NewOrderItemRepository(db.DB) // Asumiendo que tienes un repositorio para las líneas
	itemRepo := itemsrepo.NewItemRepository(db.DB)

	for _, line := range lines {
		son, _ := itemRepo.FindSonId(line.ItemID)
		orderLine := models.OrderItem{
			OrderID:       order.ID,
			ItemID:        son,
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
		if fatherOrder.OrderTypeID == uint64(2) && line.ClientID != nil {
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
		// if dataItem.Type != nil {
		// 	model.OrderTypeID = *dataItem.Type
		// }
		error := order.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders are updated", "Errors": errorList}})

}

func EditOrdersLines(c *gin.Context) {
	var requestBody dtos.OrdersLinesToUpdatePartially
	var errorList []error
	var failedIds []int

	token := c.GetHeader("Authorization")

	updateCallback := func(c *gin.Context, dataItem dtos.LineToUpdate, model *models.OrderItem, err error, errorList *[]error) {
		if err != nil {
			*errorList = append(*errorList, err)
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

	}
	updateOrderLineHandler(c, requestBody, token, &failedIds, &errorList, updateCallback, true)

	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders lines are updated", "Errors": errorList, "Not_permited_lines_ids": failedIds}})

}

func AddQuantityToOrdersLines(c *gin.Context) {
	var requestBody dtos.OrdersLinesToUpdatePartially
	var errorList []error
	var failedIds []int
	var list string

	token := c.GetHeader("Authorization")

	updateCallback := func(c *gin.Context, dataItem dtos.LineToUpdate, model *models.OrderItem, err error, errorList *[]error) {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if dataItem.RecivedQuantity != nil {
			fmt.Println("entra en el callback")
			newQuantity := *dataItem.RecivedQuantity + model.RecivedAmount

			if model.Amount >= newQuantity {
				model.RecivedAmount = newQuantity
			} else {
				customError := errors.New("se ha intentado actualizar la cantidad por encima del límite máximo")
				*errorList = append(*errorList, customError)
				return
			}

		}

	}
	updateOrderLineHandler(c, requestBody, token, &failedIds, &errorList, updateCallback, false)

	if len(errorList) != 0 {
		list = "Se ha intenado aumentar la cantidad mas allá del máximo"
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders lines are updated", "Errors": list, "Not_permited_lines_ids": failedIds}})
}

func RemoveQuantityToOrdersLines(c *gin.Context) {
	var requestBody dtos.OrdersLinesToUpdatePartially
	var errorList []error
	var failedIds []int
	var list string

	token := c.GetHeader("Authorization")

	updateCallback := func(c *gin.Context, dataItem dtos.LineToUpdate, model *models.OrderItem, err error, errorList *[]error) {
		if err != nil {
			*errorList = append(*errorList, err)
			return
		}
		if dataItem.RecivedQuantity != nil {

			newQuantity := model.RecivedAmount - *dataItem.RecivedQuantity

			if newQuantity >= 0 {
				model.RecivedAmount = newQuantity
			} else {

				customError := errors.New("se ha intendo actualizar la cantidad por debajo del límite máximo")
				*errorList = append(*errorList, customError)
				return
			}

		}

	}
	updateOrderLineHandler(c, requestBody, token, &failedIds, &errorList, updateCallback, false)

	if len(errorList) != 0 {
		list = "Se ha intenado reducir la cantidad por debajo de 0"
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders lines are updated", "Errors": list, "Not_permited_lines_ids": failedIds}})
}

func updateOrderLineHandler(

	c *gin.Context,
	requestBody dtos.OrdersLinesToUpdatePartially,
	token string,
	failedIds *[]int,
	errorList *[]error,
	callback func(*gin.Context, dtos.LineToUpdate, *models.OrderItem, error, *[]error),
	admin bool) {

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		*errorList = append(*errorList, err)
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {
		var assign dtos.Assign
		repo := tokenrepo.NewTokenRepository(db.DB)
		user, _ := repo.ReturnDataByToken(token)
		query := `SELECT id FROM assigned_lines WHERE  order_line_id = ? and user_id = ? LIMIT 1`

		err := db.DB.Raw(query, dataItem.Id, user.UserId).Scan(&assign).Error

		if (err != nil || assign.ID == 0) && !admin {
			*failedIds = append(*failedIds, int(dataItem.Id))

		} else {

			updateOrderLine(c, dataItem, errorList, callback)
		}

	}

}

func updateOrderLine(
	c *gin.Context,
	dataItem dtos.LineToUpdate,
	errorList *[]error,
	callback func(*gin.Context, dtos.LineToUpdate, *models.OrderItem, error, *[]error)) {
	orderLines := orderitemrepo.NewOrderItemRepository(db.DB)
	//repoStockDef := stockdeficitrepo.NewStockDeficitRepository(db.DB)
	model, err := orderLines.FindByID(dataItem.Id)

	callback(c, dataItem, model, err, errorList)
	error := orderLines.Update(model)
	// if(model.StoreID == 1 && ){
	// 	repoStockDef.CalcStockDeficitByItem(model.ItemID, model.StoreID)

	// }
	if error != nil {
		*errorList = append(*errorList, error)
	}

}
