package controllers

import (
	"net/http"
	"strconv"

	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/storelocationrepo"

	"github.com/gin-gonic/gin"
)

func GetStoreLocation(c *gin.Context) {
	var err error
	var data []models.StoreLocation
	var datum *models.StoreLocation
	var recount int64

	repo := storelocationrepo.NewStoreLocationRepository(db.DB)
	storeLocation, _ := strconv.Atoi(c.DefaultQuery("store_location", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if storeLocation != 0 {
		datum, err = repo.FindByID(uint64(storeLocation))
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

func PostStoreLocation(c *gin.Context) {
	var requestBody dtos.StoreLocationCreateReq

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		repo := storelocationrepo.NewStoreLocationRepository(db.DB)
		model := models.StoreLocation{
			StoreID: dataItem.StoreID,
			Code:    dataItem.Code,
			Name:    dataItem.Name,
		}
		repo.Create(&model)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created"}})

}

func PatchStoreLocation(c *gin.Context) {
	var requestBody dtos.StoreLocationUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := storelocationrepo.NewStoreLocationRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if requestObject.StoreID != nil {
			model.StoreID = *requestObject.StoreID
		}
		if requestObject.Code != nil {
			model.Code = *requestObject.Code
		}
		if requestObject.Name != nil {
			model.Name = *requestObject.Name
		}
		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Store locations are updated", "Errors": errorList}})

}
