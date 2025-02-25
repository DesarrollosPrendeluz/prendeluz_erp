package orderlineboxrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderLineBoxRepo interface {
	repositories.Repository[models.OrderLineBox]
	GetByLineId(lineId int) ([]models.OrderLineBox, error)
	GetByBox(boxId int) ([]models.OrderLineBox, error)
	GetOrCreateByOrderLineAndBoxId(orderLineId int, boxId int, quantity int) (models.OrderLineBox, bool, error)
}
