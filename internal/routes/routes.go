package routes

import (
	"prendeluz/erp/internal/controllers"
	"prendeluz/erp/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)

	orderRoutes := router.Group("/order").Use(middlewares.Auth)
	{
		orderRoutes.POST("/add", controllers.AddOrder)
		orderRoutes.GET("", controllers.GetOrders)
	}

	storeRoutes := router.Group("/store").Use(middlewares.Auth)
	{
		storeRoutes.PATCH("/:order_code", controllers.UpdateStore)
		storeRoutes.GET("/:store_name", controllers.GetStoreStock)
	}

}
