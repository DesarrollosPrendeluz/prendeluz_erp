package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StoreManager    = 4
	StoreWorker     = 5
	StoreSupervisor = 6
)

// Que el token de la request esté asinado al usuario que tenga un rol válido para la acción
// func SuperAdminStoreUsers(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	var roles []int
// 	roles = append(roles, StoreManager)
// 	checkRole(c, token, roles)

// }
func AdminStoreUsers(c *gin.Context) {
	//token := c.GetHeader("Authorization")
	// var roles []int
	// roles = append(roles, StoreManager, StoreSupervisor)
	//checkRole(c, token, roles)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	c.Abort()
	return

}
func AllStoreUsers(c *gin.Context) {
	//token := c.GetHeader("Authorization")
	// var roles []int
	// roles = append(roles, StoreManager, StoreSupervisor, StoreWorker)
	//checkRole(c, token, roles)
	c.Next()

}

// func checkRole(c *gin.Context, userToken string, roles []int) {
// 	token := c.GetHeader("Authorization")
// 	tokenRepo := tokenrepo.NewTokenRepository(db.DB)
// 	valid, err := tokenRepo.CheckCredentials(token)
// 	if err != nil || !valid {
// 		// Si hay un error o el token no es válido, detener la cadena de middleware y responder con 401
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		c.Abort() // Detener la cadena de middlewares y controladores
// 		return
// 	}

// 	// Si el token es válido, permitir que la solicitud continúe
// 	c.Next()

// }
