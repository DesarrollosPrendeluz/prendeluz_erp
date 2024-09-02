package routes

import (
	"prendeluz/erp/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)

	testRoutes := router.Group("/test").Use()
	{
		testRoutes.POST("/", controllers.Login)

	}

	orderRoutes := router.Group("/order").Use()
	{
		orderRoutes.POST("/add", controllers.AddOrder)
		orderRoutes.GET("", controllers.GetOrders)
	}

	storeRoutes := router.Group("/store").Use()
	{
		storeRoutes.PATCH("/:order_code", controllers.UpdateStore)
		storeRoutes.GET("/:store_name", controllers.GetStoreStock)
	}

}
