package services

import (
	"prendeluz/erp/internal/dtos"
)

type OrderLineService interface {
	OrderLineLabel(id int) (dtos.OrderLineLable, error)
}
