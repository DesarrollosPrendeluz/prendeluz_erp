package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	service "prendeluz/erp/internal/services/father_order_service"
	serviceOrderLine "prendeluz/erp/internal/services/order_lines"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFatherOrdersData(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderStatus, _ := strconv.Atoi(c.DefaultQuery("status_id", "0"))
	orderType, _ := strconv.Atoi(c.DefaultQuery("type_id", "0"))
	fatherCode := c.Query("father_order_code")

	results, recount, err := service.NewFatherOrderService().FindAllWithAssocData(fatherCode, orderType, orderStatus, pageSize, page)

	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results, "recount": recount}})

}

func ClosePickingOrder(c *gin.Context) {

	fatherOrder, _ := strconv.Atoi(c.DefaultQuery("father_order", "0"))

	err := service.NewFatherOrderService().ClosePickingOrders(uint64(fatherOrder))
	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": "ok"}})

}
func GetOrderLinesByFatherId(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	store_id, _ := strconv.Atoi(c.DefaultQuery("store_id", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	fatherCode := c.Query("father_order_code")
	eanOrder := c.Query("ean_order")
	locOrder := c.Query("loc_order")
	ean := c.Query("ean")
	supplierSku := c.Query("ref_prov")

	results, recount, err := service.NewFatherOrderService().FindLinesByFatherOrderCode(pageSize, page, fatherCode, ean, supplierSku, store_id, eanOrder, locOrder)

	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": results, "recount": recount}})
}

func DownloadPickingExcelByFatherId(c *gin.Context) {

	fatherCode := c.Query("father_order_code")

	results, _, err := service.NewFatherOrderService().FindLinesByFatherOrderCode(-1, 1, fatherCode, "", "", 1, "ASC", "ASC")

	if err != nil {
		// Manejo del error si las credenciales no son correctas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	data := serviceOrderLine.NewOrderLineServiceImpl().ReturnDownloadPickingExcel(results)

	fechaActual := time.Now().Format("2006-01-02 15:04:05")
	code := "picking_" + fatherCode + "_" + fechaActual + ".xlsx"
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{
		"file":     data,
		"filename": code,
	}})
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

func DownLoadExcelForAmazon(c *gin.Context) {
	fatherOrderID, _ := strconv.Atoi(c.DefaultQuery("fatherOrderId", "1"))
	data := service.NewFatherOrderService().DownloadOrdersExcelToAmazon(uint64(fatherOrderID))
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{
		"file":     data,
		"filename": "amazon_orders.xlsx",
	}})

}
