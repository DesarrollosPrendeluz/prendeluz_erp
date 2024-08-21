package controllers

import (
	"log"
	"net/http"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/services/order"

	"github.com/gin-gonic/gin"
)

func AddOrder(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	serviceOrder := services.NewOrderService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Fatal(err)
		return

	}

	serviceOrder.UploadOrderExcel(file, header.Filename)

	c.JSON(http.StatusCreated, gin.H{"message": "Upload succesfully"})

}

func GetOrders(c *gin.Context) {
	results := make(map[string][]dtos.ItemInfo)

	orderService := services.NewOrderService()
	orders, err := orderService.GetOrders()
	for _, order := range orders {
		results[order.OrderCode] = order.ItemsOrdered
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
	return

}
