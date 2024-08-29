package orderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	*repositories.GORMRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) *OrderRepoImpl {
	return &OrderRepoImpl{repositories.NewGORMRepository(db, models.Order{})}

}

func (repo *OrderRepoImpl) FindByOrderCode(orderCode string) (models.Order, error) {
	var order models.Order

	results := repo.DB.Where("orden_compra LIKE ?", "%"+orderCode+"%").First(&order)

	return order, results.Error
}

func (repo *OrderRepoImpl) UpdateStatus(newStatus string, orderID uint64) error {
	results := repo.DB.Model(models.Order{}).Where("id = ?", orderID).Update("status", newStatus)

	return results.Error

}
