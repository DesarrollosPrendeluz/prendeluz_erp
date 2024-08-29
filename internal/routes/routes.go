package routes

import (
	"prendeluz/erp/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	orderRoutes := router.Group("/order")
	{
		orderRoutes.POST("/add", controllers.AddOrder)
		orderRoutes.GET("", controllers.GetOrders)
	}

	storeRoutes := router.Group("/store")
	{
		storeRoutes.PATCH("/:order_code", controllers.UpdateStore)
		storeRoutes.GET("/:store_name", controllers.GetStoreStock)
	}

}
