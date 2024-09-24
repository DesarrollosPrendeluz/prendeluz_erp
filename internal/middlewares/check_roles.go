package middlewares

import (
	"net/http"
	"prendeluz/erp/internal/db"

	"prendeluz/erp/internal/repositories/tokenrepo"
	"prendeluz/erp/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	StoreManager    = 4
	StoreWorker     = 5
	StoreSupervisor = 6
)

type Role struct {
	RoleID int `json:"role_id"`
}
type Assign struct {
	ID int `json:"id"`
}

// Que el token de la request esté asinado al usuario que tenga un rol válido para la acción
// func SuperAdminStoreUsers(c *gin.Context) {
// 	token := c.GetHeader("Authorization")
// 	var roles []int
// 	roles = append(roles, StoreManager)
// 	checkRole(c, token, roles)

// }
func AdminStoreUsers(c *gin.Context) {
	token := c.GetHeader("Authorization")
	var roles []int
	roles = append(roles, StoreManager, StoreSupervisor)
	checkRole(c, token, roles)

}
func AllStoreUsers(c *gin.Context) {
	token := c.GetHeader("Authorization")
	var roles []int
	roles = append(roles, StoreManager, StoreSupervisor, StoreWorker)
	checkRole(c, token, roles)

}

func checkRole(c *gin.Context, userToken string, roles []int) {
	flag := false
	result, err := ObtainRole(userToken)

	if err != nil || !utils.ContainsInt(roles, result.RoleID) {
		flag = true
	}

	if flag {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "El rol de su suario carece de los permisos necesarios"})
		c.Abort()
		return
	}
	c.Next()

}

func ObtainRole(userToken string) (Role, error) {
	var result Role

	repo := tokenrepo.NewTokenRepository(db.DB)
	model, err := repo.ReturnDataByToken(userToken)
	if err != nil {
		return result, err
	}
	query := `SELECT role_id FROM model_has_roles WHERE  model_type= 'App\\Models\\User' and model_id = ? LIMIT 1`
	err = db.DB.Raw(query, model.UserId).Scan(&result).Error
	return result, err

}
