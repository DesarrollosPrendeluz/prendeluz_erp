package controllers

import (
	"net/http"
	stadisticService "prendeluz/erp/internal/services/stadistics"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrderHisotric(c *gin.Context) {

	father_id, _ := strconv.ParseUint(c.DefaultQuery("father_id", "1"), 10, 64)
	data := stadisticService.NewStadisitcService().GetChangeStadistics(father_id)

	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": data}})

}
