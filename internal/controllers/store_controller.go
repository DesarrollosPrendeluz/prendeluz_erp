package controllers

import (
	"net/http"
	"prendeluz/erp/internal/services/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateStore(c *gin.Context) {
	serviceStore := services.NewStoreService()
	orderCode := c.Param("order_code")
	serviceStore.UpdateStoreStock(orderCode)
	c.JSON(http.StatusOK, gin.H{"message": "Updated stock"})
}

func GetStoreStock(c *gin.Context) {
	services := services.NewStoreService()
	storeName := c.Param("store_name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize
	stock := services.GetStoreStock(storeName, pageSize, offset)
	c.IndentedJSON(http.StatusOK, gin.H{"results": stock})
}
