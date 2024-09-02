package middlewares

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/tokenrepo"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	tokenRepo := tokenrepo.NewTokenRepository(db.DB)
	valid, err := tokenRepo.CheckCredentials(token)
	if err != nil || !valid {
		// Si hay un error o el token no es válido, detener la cadena de middleware y responder con 401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort() // Detener la cadena de middlewares y controladores
		return
	}

	// Si el token es válido, permitir que la solicitud continúe
	c.Next()

}
