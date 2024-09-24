package controllers

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/assignlinerepo"

	"net/http"
	"prendeluz/erp/internal/db"

	"github.com/gin-gonic/gin"
)

func CreateOrderLinesAssignation(c *gin.Context) {
	var requestBody dtos.ItemsAssigantion

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": err.Error()}})
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Assignations {
		line := dataItem.LineID
		user := dataItem.UserID
		repo := assignlinerepo.NewAssignedLineImplRepository(db.DB)
		asignation := models.AssignedLine{
			UserID:      user,
			OrderLineID: line,
		}
		repo.Create(&asignation)
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Asignations are created"}})
}

func EditOrderLinesAssignation(c *gin.Context) {
	var requestBody dtos.ItemsAssigantionEdit
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los valores del cuerpo
	asignation := assignlinerepo.NewAssignedLineImplRepository(db.DB)
	for _, dataItem := range requestBody.Assignations {
		model, err := asignation.FindByID(dataItem.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		model.UserID = dataItem.UserID

		error := asignation.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Asignations are updated", "Errors": errorList}})

}
