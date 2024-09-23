package ordererrorrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderErrRepoImpl struct {
	*repositories.GORMRepository[models.ErrorOrder]
}

func NewOrderErrRepository(db *gorm.DB) *OrderErrRepoImpl {
	return &OrderErrRepoImpl{repositories.NewGORMRepository(db, models.ErrorOrder{})}
}

// Se retornan las lineas de un pedido por el id del pedido
