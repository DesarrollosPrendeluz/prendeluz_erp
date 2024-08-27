package routes

import (
	"prendeluz/erp/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	skuRoutes := router.Group("/skus")
	{
		skuRoutes.GET("/", controllers.GetSkus)
		skuRoutes.GET("/order", controllers.GetOrder)
	}

	orderRoutes := router.Group("/order")
	{
		orderRoutes.POST("/add", controllers.AddOrder)
		orderRoutes.GET("/", controllers.GetOrders)
	}

	storeRoutes := router.Group("/store")
	{
		storeRoutes.POST("/", controllers.UpdateStore)
	}

}
