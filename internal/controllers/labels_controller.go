package controllers

import (
	"net/http"
	services "prendeluz/erp/internal/services/order_lines"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OrderLineLabels(c *gin.Context) {
	line, _ := strconv.Atoi(c.DefaultQuery("line_id", "0"))

	if line == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"Results": gin.H{"error": "No se ha proporcionado id de linea "}})
		return
	}
	orderLine, _ := services.NewOrderLineServiceImpl().OrderLineLabel(line)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": orderLine}})

}
