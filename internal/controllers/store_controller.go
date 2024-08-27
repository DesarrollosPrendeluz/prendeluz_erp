package controllers

import (
	"net/http"
	"prendeluz/erp/internal/services/store"

	"github.com/gin-gonic/gin"
)

func UpdateStore(c *gin.Context) {
	serviceStore := services.NewStoreService()
	serviceStore.UpdateStoreStock(1)

	c.JSON(http.StatusCreated, gin.H{"message": "Updated stock"})
}
