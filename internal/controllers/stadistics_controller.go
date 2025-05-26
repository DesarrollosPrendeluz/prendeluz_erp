package controllers

import (
	"net/http"
	itemsService "prendeluz/erp/internal/services/items"
	stadisticService "prendeluz/erp/internal/services/stadistics"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrderHisotric(c *gin.Context) {

	father_code := c.Query("father_code")
	data := stadisticService.NewStadisitcService().GetChangeStadistics(father_code)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data}})

}

func GetLines(c *gin.Context) {
	father_code := c.Query("father_code")
	data := stadisticService.NewStadisitcService().GetOrderLineStadistics(father_code)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data}})
}

func GetItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	flag := c.Query("flag")
	offset := (page - 1) * pageSize
	data, recount := itemsService.NewItemsServiceImpl().GetItemsForDashboard(flag, offset, pageSize)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data, "recount": recount}})
}
