package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	services "prendeluz/erp/internal/services/stock_deficit"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStockDeficit(c *gin.Context) {
	store, _ := strconv.Atoi(c.DefaultQuery("store_id", "1"))
	supplier, _ := strconv.Atoi(c.DefaultQuery("supplier", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	filter := c.Query("filter")

	repo := stockdeficitrepo.NewStockDeficitRepository(db.DB)
	recount, _ := repo.CountConditional(store)
	var stockDeficits []models.StockDeficit
	if filter != "" {
		service := services.NewStockDeficitService()
		stockDeficits, _ = service.SearchBySkuAndEan(filter, store, page, pageSize)
		recount = 1
	} else {
		if supplier == 0 {
			stockDeficits, _ = repo.GetallByStore(store, pageSize, page)
		} else {
			stockDeficits, _ = repo.GetallByStoreAndSupplier(store, supplier, pageSize, page)

		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": stockDeficits, "recount": recount}})

}

func DownloadStockDeficitExcel(c *gin.Context) {
	data := services.NewStockDeficitService().ReturnDownloadStockDeficitExcel(2)
	fechaActual := time.Now().Format("2006-01-02 15:04:05")
	code := "stock_deficit_" + fechaActual + ".xlsx"
	c.JSON(http.StatusAccepted, gin.H{"Results": gin.H{
		"file":     data,
		"filename": code,
	}})

}

func CalcStockDeficitByOrder(c *gin.Context) {
	services.NewStockDeficitService().CalcStockDeficitByFatherOrder(7)
}
