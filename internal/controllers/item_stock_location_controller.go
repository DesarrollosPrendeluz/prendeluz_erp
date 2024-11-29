package controllers

import (
	"net/http"
	"strconv"

	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/itemlocationrepo"

	"github.com/gin-gonic/gin"
)

func GetItemStockLocation(c *gin.Context) {
	var err error
	var data []models.ItemLocation
	var datum *models.ItemLocation
	var recount int64

	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	storeLocation, _ := strconv.Atoi(c.DefaultQuery("item_stock_location", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if storeLocation != 0 {
		datum, err = repo.FindByID(uint64(storeLocation))
		if datum != nil {
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

func PostItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationCreateReq

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
		model := models.ItemLocation{
			ItemMainSku:     dataItem.ItemMainSku,
			StoreLocationID: dataItem.StoreLocationID,
			Stock:           dataItem.Stock,
		}
		repo.Create(&model)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created"}})

}

func PatchItemStockLocation(c *gin.Context) {
	var requestBody dtos.ItemStockLocationUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if requestObject.ItemMainSku != nil {
			model.ItemMainSku = *requestObject.ItemMainSku
		}
		if requestObject.StoreLocationID != nil {
			model.StoreLocationID = *requestObject.StoreLocationID
		}
		if requestObject.Stock != nil {
			model.Stock = *requestObject.Stock
		}
		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Stock locations are updated", "Errors": errorList}})

}
