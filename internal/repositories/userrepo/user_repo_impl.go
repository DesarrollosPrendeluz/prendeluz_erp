package userrepo

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/tokenrepo"

	"prendeluz/erp/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserImpl struct {
	*repositories.GORMRepository[models.User]
}

func NewUsersRepository(db *gorm.DB) *UserImpl {
	return &UserImpl{repositories.NewGORMRepository(db, models.User{})}
}

func (repo *UserImpl) CheckCredentials(email string, password string) (string, error) {
	var user models.User

	// Buscar el usuario por correo electrónico
	results := repo.DB.Where("email = ?", email).First(&user)
	if results.Error != nil {
		// Si ocurre un error al buscar el usuario, retornamos el error
		return "", fmt.Errorf("usuario no encontrado o error en la consulta")
	}

	// Comparar la contraseña ingresada con la

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		// Si las contraseñas no coinciden, retornamos un error
		return "", fmt.Errorf("credenciales incorrectas")
	}

	// Si las credenciales son correctas, generar y retornar el token de usuario
	token := repo.GenerateUserToken(uint64(user.ID))
	return token, nil
}

func (repo *UserImpl) GenerateUserToken(userId uint64) string {
	token := utils.GenerateRandomString(240)
	tokenRepo := *tokenrepo.NewTokenRepository(db.DB)
	newToken := &models.AccesTokens{
		UserId: userId,
		Token:  token,
		Valid:  true,
	}
	tokenRepo.Create(newToken)
	//results := repo.DB.Where("parent_sku LIKE ?", "%"+"sku_parent"+"%") //.First(&"storeStocks")

	return token

}
