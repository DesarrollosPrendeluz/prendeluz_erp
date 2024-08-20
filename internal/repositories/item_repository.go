package repositories

import (
	"gorm.io/gorm"
	"prendeluz/erp/internal/models"
)

type ItemRepo struct {
	GORMRepository[models.Item]
}

func NewItemRepository(db *gorm.DB) *ItemRepo {
	return &ItemRepo{*NewGORMRepository(db, models.Item{})}
}

func (repo *ItemRepo) FindByMainSku(sku string) (models.Item, error) {
	var item models.Item

	result := repo.db.Where("main_sku LIKE ?", "%"+sku+"%").First(&item)

	return item, result.Error

}
