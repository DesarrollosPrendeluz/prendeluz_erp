package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFatherOrdersData(c *gin.Context) {
	repo := fatherorderrepo.NewFatherOrderRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderStatus, _ := strconv.Atoi(c.DefaultQuery("status_id", "0"))
	orderType, _ := strconv.Atoi(c.DefaultQuery("type_id", "0"))
	fatherCode := c.Query("father_order_code")

	results, recount, err := repo.FindAllWithAssocData(pageSize, page, fatherCode, orderType, orderStatus)

	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results, "recount": recount}})

}
func GetOrderLinesByFatherId(c *gin.Context) {

	repo := fatherorderrepo.NewFatherOrderRepository(db.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	fatherCode := c.Query("father_order_code")

	results, recount, err := repo.FindLinesByFatherOrderCode(pageSize, page, fatherCode)

	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results, "recount": recount}})
}
func UpdateFatherOrders(c *gin.Context) {
	var requestBody dtos.OrdersToUpdatePartially
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	fatherOrder := fatherorderrepo.NewFatherOrderRepository(db.DB)
	for _, dataItem := range requestBody.Data {
		model, err := fatherOrder.FindByID(dataItem.Id)
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
		error := fatherOrder.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Orders are updated", "Errors": errorList}})

}
