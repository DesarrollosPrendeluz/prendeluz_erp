package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/palletrepo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPallet(c *gin.Context) {
	var err error
	var data []models.Pallet
	var datum *models.Pallet
	var recount int64

	repo := palletrepo.NewPalletRepository(db.DB)
	pallet, _ := strconv.Atoi(c.DefaultQuery("pallet", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if pallet != 0 {
		datum, err = repo.FindByID(uint64(pallet))
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

func GetPalletByOrderID(c *gin.Context) {
	var err error
	var data []models.Pallet
	var recount int64

	repo := palletrepo.NewPalletRepository(db.DB)
	orderId, _ := strconv.Atoi(c.DefaultQuery("order_id", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	data, err = repo.GetBoxesAndLinesRaletedDataByOrderId(orderId, pageSize, page)

	recount = 1

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
	for _, dataItem := range requestBody.Data {

		repo := palletrepo.NewPalletRepository(db.DB)
		model := models.Pallet{
			OrderID: dataItem.OrderID,
			Number:  int(dataItem.Number),
			Label:   dataItem.Label,
		}
		repo.Create(&model)
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
	repo := palletrepo.NewPalletRepository(db.DB)
	for _, requestObject := range requestBody.Data {
		model, err := repo.FindByID(requestObject.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if requestObject.OrderID != nil {
			model.OrderID = *requestObject.OrderID
		}
		if requestObject.Label != nil {
			model.Label = *requestObject.Label
		}
		if requestObject.Number != nil {
			model.Number = *requestObject.Number
		}

		error := repo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Store locations are updated", "Errors": errorList}})

}
