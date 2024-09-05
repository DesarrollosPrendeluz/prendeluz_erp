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

	results := repo.DB.Where("code LIKE ?", "%"+orderCode+"%").First(&order)

	return order, results.Error
}

// FindOrderByDate retrieves an order based on a date range.
// It accepts startDate and endDate as strings in the format "YYYY-MM-DD".
// If both startDate and endDate are provided, it returns orders between these dates.
// If only startDate is provided, it returns orders for that specific date.
// If only endDate is provided, it returns orders for that specific date.
// Returns an Order struct and an error if something goes wrong.
func (repo *OrderRepoImpl) FindOrderByDate(startDate string, endDate string) ([]models.Order, error) {
	var order []models.Order
	var results *gorm.DB

	if startDate != "" && endDate != "" {

		results = repo.DB.Where("date(created_at) BETWEEN ? AND ?", startDate, endDate).First(&order)
	} else if startDate != "" {
		results = repo.DB.Where("date(created_at) = ?", startDate).First(&order)
	} else {
		results = repo.DB.Where("date(created_at) = ?", endDate).First(&order)
	}

	return order, results.Error
}

func (repo *OrderRepoImpl) UpdateStatus(newStatus string, orderID uint64) error {
	results := repo.DB.Model(models.Order{}).Where("id = ?", orderID).Update("status", newStatus)

	return results.Error

}
