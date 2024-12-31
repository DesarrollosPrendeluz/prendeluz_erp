package controllers

import (
	"net/http"
	stockService "prendeluz/erp/internal/services/stock"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStockExcelData(c *gin.Context) {
	stockService := stockService.NewStockService()
	store_id, _ := strconv.Atoi(c.DefaultQuery("store_id", "1"))

	data := stockService.ReturnDownloadStockExcel(store_id)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"file": data, "name": "stock.xlsx"}})

}
