package routes

import (
	"prendeluz/erp/internal/controllers"
	"prendeluz/erp/internal/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		// Permitir solicitudes desde http://localhost:3000
		AllowOrigins: []string{"http://localhost:3000"},
		// Permitir métodos HTTP
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// Permitir encabezados específicos
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		// Permitir enviar cookies en las solicitudes
		AllowCredentials: true,
		// Definir el tiempo de caché para la respuesta preflight
		MaxAge: 12 * time.Hour,
	}))
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

	stockDeficit := router.Group("/stock_deficit").Use(middlewares.Auth)
	{
		stockDeficit.GET("", controllers.GetStockDeficit)
	}

	// orderRoutes := router.Group("/order").Use(middlewares.Auth)
	// {
	// 	orderRoutes.PATCH("", controllers.UpdateStore)
	// 	orderRoutes.GET("", controllers.GetStoreStock)
	// }

}
