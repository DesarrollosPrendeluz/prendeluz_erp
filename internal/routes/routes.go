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

	allUsersOrderRoutes := router.Group("/order").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersOrderRoutes.GET("", controllers.GetOrders)
		allUsersOrderRoutes.GET("/status", controllers.GetOrderStatus)
		allUsersOrderRoutes.GET("/type", controllers.GetOrderTypes)
		allUsersOrderRoutes.GET("/supplierOrders", controllers.GetSupplierOrders)
		allUsersOrderRoutes.GET("/supplierOrders/download", controllers.DownloadSupplierOrderExcel)

	}

	adminUsersOrderRoutes := router.Group("/order").Use(middlewares.Auth, middlewares.AdminStoreUsers)
	{
		adminUsersOrderRoutes.POST("/add", controllers.AddOrder)
		adminUsersOrderRoutes.POST("/addByRequest", controllers.CreateOrder)
		adminUsersOrderRoutes.PATCH("", controllers.EditOrders)
	}

	//TODO: por implementar
	// allUsersOrderLineRoutes := router.Group("/order/orderLines").Use(middlewares.Auth, middlewares.AllStoreUsers)
	// {
	// 	allUsersOrderLineRoutes.GET("", controllers.GetOrders)
	// 	allUsersOrderLineRoutes.POST("", controllers.GetOrderStatus)
	// 	allUsersOrderLineRoutes.PATCH("", controllers.EditOrdersLines)
	// }

	allUsersOrderLineAssignRoutes := router.Group("/order/orderLines/asignation").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersOrderLineAssignRoutes.POST("", controllers.CreateOrderLinesAssignation)
		allUsersOrderLineAssignRoutes.PATCH("", controllers.EditOrderLinesAssignation)
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

}
