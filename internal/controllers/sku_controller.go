package controllers

import (
	"log"
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"github.com/gin-gonic/gin"
)

func GetSkus(c *gin.Context) {
	var skus []models.Sku
	skuRepostory := repositories.NewSkuRepository(db.DB)

	skus, err := skuRepostory.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": skus})
	return

}

// EndPoint de prueba para leer el excel
func GetOrder(c *gin.Context) {
	var data []byte

	c.Data(http.StatusOK, "application/json", data)
}
