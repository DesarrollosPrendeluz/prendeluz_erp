package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	log.Println("hola")
	email := c.Param("email")
	password := c.Param("password")
	log.Println(email)

	log.Println(password)

}

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token != "Bearer mysecrettoken" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort() // Detiene la cadena de middlewares y controladores
		return
	}

	c.Next()

}
