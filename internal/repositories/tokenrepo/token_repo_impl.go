package tokenrepo

import (
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type TokenImpl struct {
	*repositories.GORMRepository[models.AccesTokens]
}

func NewTokenRepository(db *gorm.DB) *TokenImpl {
	return &TokenImpl{repositories.NewGORMRepository(db, models.AccesTokens{})}
}

// Verifica si el token dado están en base de datos y es válido
func (repo *TokenImpl) CheckCredentials(token string) (bool, error) {
	var tokenModel models.AccesTokens

	// Buscar el usuario por correo electrónico
	results := repo.DB.Where("token = ?", token).Where("valid = ?", 1).First(&tokenModel)
	if results.Error != nil {
		// Si ocurre un error al buscar el usuario, retornamos el error
		return false, fmt.Errorf("no se ha encontrado el token o ha caducado")
	}

	return true, nil
}

// Verifica si el token dado están en base de datos y es válido
func (repo *TokenImpl) ReturnDataByToken(token string) (models.AccesTokens, error) {
	var tokenModel models.AccesTokens

	// Buscar el usuario por correo electrónico
	results := repo.DB.Where("token = ?", token).First(&tokenModel)
	if results.Error != nil {
		// Si ocurre un error al buscar el usuario, retornamos el error
		fmt.Println("no se ha encontrado el token o ha caducado")
		return tokenModel, results.Error
	}

	return tokenModel, nil
}
