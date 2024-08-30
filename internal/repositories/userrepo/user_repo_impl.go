package userrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type UserImpl struct {
	*repositories.GORMRepository[models.User]
}

func NewUsersitory(db *gorm.DB) *UserImpl {
	return &UserImpl{repositories.NewGORMRepository(db, models.User{})}
}

func (repo *UserImpl) CheckCredentials(sku_parent string) (models.User, error) {
	var storeStocks models.User

	results := repo.DB.Where("parent_sku LIKE ?", "%"+sku_parent+"%").First(&storeStocks)

	return storeStocks, results.Error
}

func (repo *UserImpl) GenerateUserToken(idStore uint64, pageSize int, offset int) (string, error) {
	var token string = "pepe"

	results := repo.DB.Where("parent_sku LIKE ?", "%"+"sku_parent"+"%") //.First(&"storeStocks")

	return token, results.Error

}
