package main

import (
	"flag"
	"log"
	"prendeluz/erp/internal/config"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	env := flag.String("env", "dev", "Variable to determine the enviroment (prod,dev,test)")
	flag.Parse()

	cfg, err := config.LoadConfig(*env)
	if err != nil {
		log.Fatal(err)

	}
	db.InitDb(&cfg.Database)
	// database := db.DB
	//database.AutoMigrate(&models.Order{})

	router := gin.Default()

	routes.RegisterRoutes(router)

	router.Run(":8080")

}
