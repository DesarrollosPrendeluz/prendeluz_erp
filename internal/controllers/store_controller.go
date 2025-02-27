package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/storerepo"
	services "prendeluz/erp/internal/services/store"
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

	storeName := c.Param("store_name")
	search := c.Query("search")
	storeId, _ := strconv.Atoi(c.DefaultQuery("store_id", "1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	calcPage := page * pageSize

	service := services.NewStoreService()
	repo := storerepo.NewStoreRepository(db.DB)
	stock := service.GetStoreStock(storeName, pageSize, calcPage, search)
	recount, _ := repo.CountConditional(storeId)
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": stock, "recount": recount}})
}

func GetStores(c *gin.Context) {
	repo := storerepo.NewStoreRepository(db.DB)
	store, err := repo.FindAll(100, 0)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": store}})

}
