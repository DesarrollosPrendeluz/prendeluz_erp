package itemsrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ItemRepoImpl struct {
	*repositories.GORMRepository[models.Item]
}

func NewItemRepository(db *gorm.DB) *ItemRepoImpl {
	return &ItemRepoImpl{repositories.NewGORMRepository(db, models.Item{})}
}
func (repo *ItemRepoImpl) FindByMainSku(sku string) (models.Item, error) {
	var item models.Item

	result := repo.DB.Where("main_sku LIKE ?", "%"+sku+"%").First(&item)

	return item, result.Error

}
