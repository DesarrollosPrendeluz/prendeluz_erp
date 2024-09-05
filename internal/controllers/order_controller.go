package controllers

import (
	"log"
	"net/http"
	"prendeluz/erp/internal/dtos"
	services "prendeluz/erp/internal/services/order"
	"strconv"

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	orders, err := orderService.GetOrders(page, pageSize, startDate, endDate)

	for _, order := range orders {
		results[order.OrderCode] = order.ItemsOrdered
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": results})
	return

}
