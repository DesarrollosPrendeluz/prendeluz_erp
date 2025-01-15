package controllers

import (
	"log"
	"net/http"
	stockService "prendeluz/erp/internal/services/stock"
	storeService "prendeluz/erp/internal/services/store"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStockExcelData(c *gin.Context) {
	store_id, _ := strconv.Atoi(c.DefaultQuery("store_id", "1"))

	data := stockService.NewStockService().ReturnDownloadStockExcel(store_id)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"file": data, "name": "stock.xlsx"}})

}

func UpdateStockByExcel(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Fatal(err)
		return

	}
	fileData, fileName, _ := storeService.NewStoreService().UploadStocks(file, header.Filename)

	c.JSON(http.StatusCreated, gin.H{"Results": gin.H{"File": fileData, "FileName": fileName}})

}
