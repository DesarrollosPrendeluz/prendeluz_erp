package controllers

import (
	"net/http"
	"strconv"

	"prendeluz/erp/internal/dtos"

	service "prendeluz/erp/internal/services/box"

	"github.com/gin-gonic/gin"
)

func GetBox(c *gin.Context) {

	box, _ := strconv.Atoi(c.DefaultQuery("box", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	data, recount, err := service.NewBoxService().GetBox(box, page, pageSize)

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
	errList := service.NewBoxService().CreateBox(requestBody)

	// Acceder a los valores del cuerpo
	if len(errList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errList}})
		return
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

	errList := service.NewBoxService().UpdateBox(requestBody)
	if len(errList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errList}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Box updated", "Errors": errorList}})

}

func DeleteBox(c *gin.Context) {
	var requestBody dtos.BoxDeleteReq
	var errorList []error

	// Intentar bindear los datos del cuerpo de la request al struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errList := service.NewBoxService().DeleteBox(requestBody)
	if len(errList) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Results": gin.H{"error": errList}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{"Ok": "Boxes are deleted", "Errors": errorList}})

}
