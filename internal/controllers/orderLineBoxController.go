package controllers

import (
	"net/http"
	"strconv"

	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/orderlineboxrepo"
	pallet_boxes_order_lines "prendeluz/erp/internal/services/pallet_boxes_order_lines"

	"github.com/gin-gonic/gin"
)

func GetOrderLineBox(c *gin.Context) {
	var err error
	var data []models.OrderLineBox
	var datum *models.OrderLineBox
	var recount int64

	repo := orderlineboxrepo.NewOrderLineBoxRepository(db.DB)
	orderlineboxrepo, _ := strconv.Atoi(c.DefaultQuery("orderlineboxrepo", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if orderlineboxrepo != 0 {
		datum, err = repo.FindByID(uint64(orderlineboxrepo))
		if datum != nil { // Verificar si datum no es nil
			data = append(data, *datum)
		}
		recount = 1
	} else {
		data, err = repo.FindAll(pageSize, page)
		recount, _ = repo.CountAll()

	}

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})

}

func PostOrderLineBox(c *gin.Context) {
	var requestBody dtos.OrderLineBoxCreateReq

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		repo := orderlineboxrepo.NewOrderLineBoxRepository(db.DB)
		model := models.OrderLineBox{
			OrderLineID: int(dataItem.OrderLineID),
			BoxID:       int(dataItem.BoxID),
			Quantity:    dataItem.Quantity,
		}
		repo.Create(&model)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created"}})

}

func PostOrderLineBoxWithProcess(c *gin.Context) {
	var requestBody dtos.OrderLineBoxProcessedCreateReq
	var response []string
	var err []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		service := pallet_boxes_order_lines.NewStockDeficitService()
		response, err = service.CheckAndCreateBoxOrderLines(dataItem.OrderLineID, dataItem.Pallet, dataItem.Box, dataItem.Quantity)

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Responses": response, "Errors": err}})

}

func PatchOrderLineBox(c *gin.Context) {
	var requestBody dtos.OrderLineBoxUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := orderlineboxrepo.NewOrderLineBoxRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if requestObject.BoxID != nil {
			model.BoxID = int(*requestObject.BoxID)
		}
		if requestObject.OrderLineID != nil {
			model.OrderLineID = int(*requestObject.OrderLineID)
		}

		if requestObject.Quantity != nil {
			model.Quantity = *requestObject.Quantity
		}

		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Store locations are updated", "Errors": errorList}})

}
