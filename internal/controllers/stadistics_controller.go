package controllers

import (
	"net/http"
	stadisticService "prendeluz/erp/internal/services/stadistics"

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
