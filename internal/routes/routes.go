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
		//https://erp.zarivy.com
		//AllowOrigins: []string{"http://127.0.0.1:3000"},
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:3001", "http://127.0.0.1:3001", "https://testerp.zarivy.com", "https://erp.zarivy.com"},
		//AllowOrigins: []string{"https://erp.zarivy.com"},
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

	//Order routes
	allUsersOrderRoutes := router.Group("/order").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersOrderRoutes.GET("", controllers.GetOrders)
		allUsersOrderRoutes.GET("/status", controllers.GetOrderStatus)
		allUsersOrderRoutes.GET("/type", controllers.GetOrderTypes)
		allUsersOrderRoutes.GET("/supplierOrders", controllers.GetSupplierOrders)
		allUsersOrderRoutes.PATCH("/closeOrders", controllers.CloseOrderLines)
		allUsersOrderRoutes.PATCH("/openOrders", controllers.OpenOrderLines)
		allUsersOrderRoutes.POST("/editOrders", controllers.UpdateOrderByExcel)
		allUsersOrderRoutes.POST("/editSupplierOrders", controllers.UpdateSupplierOrderByExcel)
		allUsersOrderRoutes.GET("/editOrders/frame", controllers.DownloadUpdateOrderByExcelFrame)
		allUsersOrderRoutes.GET("/editSupplierOrders/frame", controllers.DownloadUpdateSupplierOrderByExcelFrame)

		allUsersOrderRoutes.GET("/supplierOrders/download", controllers.DownloadSupplierOrderExcel)

	}
	adminUsersOrderRoutes := router.Group("/order").Use(middlewares.Auth, middlewares.AdminStoreUsers)
	{
		adminUsersOrderRoutes.POST("/add", controllers.AddOrder)
		adminUsersOrderRoutes.GET("/add/frame", controllers.DownloadAddOrderFrame)
		adminUsersOrderRoutes.POST("/addByRequest", controllers.CreateOrder)
		adminUsersOrderRoutes.PATCH("", controllers.EditOrders)
	}

	//Order lines
	// adminUsersOrderLineRoutes := router.Group("/order/orderLines").Use(middlewares.Auth, middlewares.AdminStoreUsers)
	// {

	// }

	allUsersOrderLineRoutes := router.Group("/order/orderLines").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersOrderLineRoutes.GET("", controllers.GetOrders)
		allUsersOrderLineRoutes.GET("/labels", controllers.OrderLineLabels)
		//allUsersOrderLineRoutes.POST("", controllers.GetOrderStatus)
		allUsersOrderLineRoutes.PATCH("/add", controllers.AddQuantityToOrdersLines)
		allUsersOrderLineRoutes.PATCH("/remove", controllers.RemoveQuantityToOrdersLines)
		allUsersOrderLineRoutes.PATCH("", controllers.EditOrdersLines)
	}

	allUsersOrderLineAssignRoutes := router.Group("/order/orderLines/asignation").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersOrderLineAssignRoutes.POST("", controllers.CreateOrderLinesAssignation)
		allUsersOrderLineAssignRoutes.PATCH("", controllers.EditOrderLinesAssignation)
	}

	//Father routes
	allUsersFatherOrderRoutes := router.Group("/fatherOrder").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		allUsersFatherOrderRoutes.GET("", controllers.GetFatherOrdersData)
		allUsersFatherOrderRoutes.PATCH("", controllers.UpdateFatherOrders)
		allUsersFatherOrderRoutes.GET("/orderLines", controllers.GetOrderLinesByFatherId)
		allUsersFatherOrderRoutes.GET("/orderLines/downloadPicking", controllers.DownloadPickingExcelByFatherId)
		allUsersFatherOrderRoutes.GET("/amazonExcel", controllers.DownLoadExcelForAmazon)
		allUsersFatherOrderRoutes.POST("/closePicking", controllers.ClosePickingOrder)
		allUsersFatherOrderRoutes.PATCH("/close", controllers.CloseOrderLines)
		allUsersFatherOrderRoutes.PATCH("/open", controllers.OpenOrderLines)

	}
	//Store
	storeRoutes := router.Group("/store").Use(middlewares.Auth)
	{
		storeRoutes.PATCH("/:order_code", controllers.UpdateStore)
		storeRoutes.GET("/:store_name", controllers.GetStoreStock)
		storeRoutes.GET("", controllers.GetStores)
		storeRoutes.POST("/excel", controllers.UpdateStockByExcel)
		storeRoutes.GET("/excel/frame", controllers.DownloadUpdateStockByExcelFrame)

	}
	//stock deficit
	stockDeficit := router.Group("/stock_deficit").Use(middlewares.Auth)
	{
		stockDeficit.GET("", controllers.GetStockDeficit)
		stockDeficit.GET("calc", controllers.CalcStockDeficitByOrder)
		stockDeficit.GET("/download", controllers.DownloadStockDeficitExcel)
	}
	//stock
	stock := router.Group("/stock").Use(middlewares.Auth)
	{
		stock.GET("getExcel", controllers.GetStockExcelData)
	}

	//stock deficit
	storeLocations := router.Group("/store_location").Use(middlewares.Auth)
	{
		storeLocations.GET("", controllers.GetStoreLocation)
		storeLocations.POST("", controllers.PostStoreLocation)
		storeLocations.PATCH("", controllers.PatchStoreLocation)

	}

	pallets := router.Group("/pallet").Use(middlewares.Auth)
	{
		pallets.GET("", controllers.GetPallet)
		pallets.GET("/crossDataByOrderId", controllers.GetPalletByOrderID)
		pallets.POST("", controllers.PostPallet)
		pallets.PATCH("", controllers.PatchPallet)

	}
	boxes := router.Group("/box").Use(middlewares.Auth)
	{
		boxes.GET("", controllers.GetBox)
		boxes.POST("", controllers.PostBox)
		boxes.PATCH("", controllers.PatchBox)

	}
	boxAdmin := router.Group("/box").Use(middlewares.Auth, middlewares.AdminStoreUsers)
	{

		boxAdmin.DELETE("", controllers.DeleteBox)

	}
	order_lines_boxes := router.Group("/order_line_boxes").Use(middlewares.Auth)
	{
		order_lines_boxes.GET("", controllers.GetOrderLineBox)
		order_lines_boxes.POST("", controllers.PostOrderLineBox)
		order_lines_boxes.POST("/withProcess", controllers.PostOrderLineBoxWithProcess)
		order_lines_boxes.PATCH("", controllers.PatchOrderLineBox)

	}

	itemStockLocations := router.Group("/item_stock_location").Use(middlewares.Auth)
	{
		itemStockLocations.GET("", controllers.GetItemStockLocation)
		itemStockLocations.DELETE("", controllers.DropItemStockLocation)
		itemStockLocations.POST("", controllers.PostItemStockLocation)
		itemStockLocations.PATCH("", controllers.PatchItemStockLocation)
		itemStockLocations.PATCH("/stockChanges", controllers.StockChanges)
		itemStockLocations.PATCH("/stockMovement", controllers.StockMovements)

	}

	//supplier
	suppliers := router.Group("/supplier").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		suppliers.GET("", controllers.GetSuppliers)
	}

	stadistics := router.Group("/stadistics").Use(middlewares.Auth, middlewares.AllStoreUsers)
	{
		stadistics.GET("olHisotricByFatherOrder", controllers.GetOrderHisotric)
		stadistics.GET("lines", controllers.GetLines)
		stadistics.GET("items", controllers.GetItems)
	}

}
