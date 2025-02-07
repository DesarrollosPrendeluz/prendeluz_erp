package controllers

import (
	"net/http"
	"prendeluz/erp/internal/dtos"
	PalletService "prendeluz/erp/internal/services/pallet"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPallet(c *gin.Context) {

	pallet, _ := strconv.Atoi(c.DefaultQuery("pallet", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	data, recount, err := PalletService.NewPalletService().Get(pallet, page, pageSize)

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})

}

func GetPalletByOrderID(c *gin.Context) {

	orderId, _ := strconv.Atoi(c.DefaultQuery("order_id", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	data, recount, err := PalletService.NewPalletService().GetPalletByOrder(orderId, page, pageSize)

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})

}

func PostPallet(c *gin.Context) {
	var requestBody dtos.PalletCreateReq

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	errList := PalletService.NewPalletService().Create(requestBody)

	// Acceder a los valores del cuerpo
	if len(errList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errList}})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created"}})

}

func PatchPallet(c *gin.Context) {
	var requestBody dtos.PalletUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	errList := PalletService.NewPalletService().Update(requestBody)

	// Acceder a los valores del cuerpo
	if len(errList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errList}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Store locations are updated", "Errors": errorList}})

}
