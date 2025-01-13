package orderlinelocationviewrepo

import (
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"strings"

	"gorm.io/gorm"
)

type OrderLineLocationViewImpl struct {
	*repositories.GORMRepository[models.OrderLineLocationView]
}

func NewOrderLineLocationViewRepository(db *gorm.DB) *OrderLineLocationViewImpl {
	return &OrderLineLocationViewImpl{repositories.NewGORMRepository(db, models.OrderLineLocationView{})}
}

func (repo *OrderLineLocationViewImpl) FindByFatherAndStoreWithOrder(father_id uint64, idStore int, orderByLocation string, orderByEan string) ([]uint64, string, error) {
	var lineIds []uint64
	var data []models.OrderLineLocationView

	query := repo.DB.
		Where("father_order_id= ?", father_id).Where("store_id = ?", idStore).Find(&data)

	if orderByLocation != "" {
		query = query.Order("store_location_code " + orderByLocation) // Se espera que `orderByLocation` sea algo como "location ASC" o "location DESC"
	}

	if orderByEan != "" {
		query = query.Order("order_line_item_ean " + orderByEan) // Similarmente, `orderByEan` puede ser "ean ASC" o "ean DESC"
	}

	// Ejecutar la consulta
	result := query.Find(&data)
	if result.Error != nil {
		return nil, "", result.Error
	}

	// Extraer los IDs en el orden en que se devuelven
	for _, datum := range data {
		lineIds = append(lineIds, uint64(datum.OrderLineID)) // Reemplaza `stock.ID` con el campo correspondiente al ID único
	}
	var idStrings []string
	for _, id := range lineIds {
		idStrings = append(idStrings, fmt.Sprintf("%d", id))
	}
	orderQuery := fmt.Sprintf("FIELD(%s, %s)", "id", strings.Join(idStrings, ", "))
	return lineIds, orderQuery, nil
}
