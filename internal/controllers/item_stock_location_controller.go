package controllers

import (
	"net/http"
	"strconv"

	"errors"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	ItemStockLocationService "prendeluz/erp/internal/services/item_stock_locations"

	"github.com/gin-gonic/gin"
)

func GetItemStockLocation(c *gin.Context) {
	var err error
	var data []models.ItemLocation
	var recount int64

	main_sku := c.Query("main_sku")
	store_id, _ := strconv.Atoi(c.DefaultQuery("store_id", "0"))
	storeLocation, _ := strconv.Atoi(c.DefaultQuery("item_stock_location", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	data, recount, err = ItemStockLocationService.NewItemStockLocationService().GetItemStockLocation(main_sku, store_id, storeLocation, page, pageSize)

	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})

}

func DropItemStockLocation(c *gin.Context) {
	location_id, _ := strconv.Atoi(c.DefaultQuery("location_id", "0"))
	var error error
	error = errors.New("se ha enviado id = 0")
	if location_id != 0 {
		error = ItemStockLocationService.NewItemStockLocationService().DropItemLocation(uint64(location_id))
	}
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": error.Error()}})

	} else {
		c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Se ha eliminado el registro de id ": location_id}})

	}

}

func PostItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationCreateReq
	var idArray []uint64

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	idArray = ItemStockLocationService.NewItemStockLocationService().PostItemStockLocation(requestBody)
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created", "CreatedIds": idArray}})

}

func PatchItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errorList = ItemStockLocationService.NewItemStockLocationService().PatchItemStockLocation(requestBody)

	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}

// n las request hay que enviarle el stock de articulo nuevo el total es decir si tenemos 100 y le restas 4 le mandamos 96
func StockChanges(c *gin.Context) {

	var requestBody dtos.ItemStockLocationStockChangeRequest
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	errorList = ItemStockLocationService.NewItemStockLocationService().StockChanges(requestBody)

	if len(errorList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errorList}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}

func StockMovements(c *gin.Context) {

	var requestBody dtos.ItemStockLocationStockMovementRequest
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errorList = ItemStockLocationService.NewItemStockLocationService().StockMovements(requestBody)

	if len(errorList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errorList}})
		return

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}

func DeleteZeroStock(c *gin.Context) {
	err := ItemStockLocationService.NewItemStockLocationService().DeleteZeroStock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Results": gin.H{"message": "Registros con stock 0 eliminados"}})
}
