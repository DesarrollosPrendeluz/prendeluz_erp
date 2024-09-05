package orderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

var Order_Status = map[string]int{
	"inicada":    1,
	"finalizada": 2,
	"en_proceso": 3,
	"en_espera":  4,
}

type OrderRepoImpl struct {
	*repositories.GORMRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) *OrderRepoImpl {
	return &OrderRepoImpl{repositories.NewGORMRepository(db, models.Order{})}

}

// Busca una orden por el codigo alfanumerico asociado
func (repo *OrderRepoImpl) FindByOrderCode(orderCode string) (models.Order, error) {
	var order models.Order

	results := repo.DB.Where("code LIKE ?", "%"+orderCode+"%").First(&order)

	return order, results.Error
}

// FindOrderByDate recupera un pedido basado en un rango de fechas.
// Acepta startDate y endDate como cadenas de texto con el formato: "YYYY-MM-DD".
// Si se proporcionan tanto startDate como endDate, devuelve los pedidos entre esas fechas.
// Si solo se proporciona startDate, devuelve los pedidos para esa fecha específica.
// Si solo se proporciona endDate, devuelve los pedidos para esa fecha específica.
// Devuelve una estructura Order y un error si algo sale mal.

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

// Actualiza el estado de una orden
// Recibe el id del nuevo estado y el id de la orden
func (repo *OrderRepoImpl) UpdateStatus(newStatus int, orderID uint64) error {
	results := repo.DB.Model(models.Order{}).Where("id = ?", orderID).Update("status_id", newStatus)

	return results.Error

}
