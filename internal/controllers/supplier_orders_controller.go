package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/orderrepo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSupplierOrders(c *gin.Context) {
	var status *int
	repo := orderrepo.NewOrderRepository(db.DB)
	statusStr := c.Query("status")
	if statusStr != "" {
		// Convierte el string a un int
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			// Si hay un error en la conversión, responde con un error de Bad Request
			c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'status' debe ser un número entero."})
			return
		}
		// Asigna el valor a la variable status como puntero
		status = &statusInt
	}
	data, err := repo.GetSupplierOrders(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": data})

}
