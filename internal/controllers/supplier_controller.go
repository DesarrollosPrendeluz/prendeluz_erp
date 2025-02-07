package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/supplierrepo"

	"github.com/gin-gonic/gin"
)

func GetSuppliers(c *gin.Context) {
	repo := supplierrepo.NewSupplierRepository(db.DB)
	store, err := repo.FindAllOrdered(1000, 0)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"Error": gin.H{"err": err}})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Results": gin.H{"data": store}})

}
