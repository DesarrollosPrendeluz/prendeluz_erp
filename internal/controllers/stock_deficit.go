package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStockDeficit(c *gin.Context) {
	store, _ := strconv.Atoi(c.DefaultQuery("store_id", "1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	repo := stockdeficitrepo.NewStockDeficitRepository(db.DB)
	recount, _ := repo.CountConditional(store)
	stockDeficits, _ := repo.GetallByStore(store, pageSize, page)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": stockDeficits, "recount": recount}})

}
