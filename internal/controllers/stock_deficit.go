package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStockDeficit(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orderRepo := stockdeficitrepo.NewStockDeficitRepository(db.DB)

	stockDeficits, err := orderRepo.FindAll(pageSize, page)
	//stockDeficits, err := orderRepo.FindByID(1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": stockDeficits}})

}
