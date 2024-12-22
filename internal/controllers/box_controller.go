package controllers

import (
	"net/http"
	"strconv"

	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"

	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/boxrepo"

	"github.com/gin-gonic/gin"
)

func GetBox(c *gin.Context) {
	var err error
	var data []models.Box
	var datum *models.Box
	var recount int64

	repo := boxrepo.NewBoxRepository(db.DB)
	box, _ := strconv.Atoi(c.DefaultQuery("box", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if box != 0 {
		datum, err = repo.FindByID(uint64(box))
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

func PostBox(c *gin.Context) {
	var requestBody dtos.BoxCreateReq

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		repo := boxrepo.NewBoxRepository(db.DB)
		model := models.Box{
			PalletID: dataItem.PalletID,
			Number:   int(dataItem.Number),
			Label:    dataItem.Label,
			Quantity: dataItem.Quantity,
		}
		repo.Create(&model)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Registers are created"}})

}

func PatchBox(c *gin.Context) {
	var requestBody dtos.BoxUpdateReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	repo := boxrepo.NewBoxRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if requestObject.PalletID != nil {
			model.PalletID = *requestObject.PalletID
		}
		if requestObject.Label != nil {
			model.Label = *requestObject.Label
		}
		if requestObject.Number != nil {
			model.Number = int(*requestObject.Number)
		}

		if requestObject.Quantity != nil {
			model.Number = *requestObject.Quantity
		}

		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Store locations are updated", "Errors": errorList}})

}
