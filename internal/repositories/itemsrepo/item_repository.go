package itemsrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type ItemRepo interface {
	repositories.Repository[models.Item]
	FindByMainSku(sku string) (models.Item, error)
	FindSonId(id uint64) (uint64, error)
	FindByIdExtraData(id int) (models.Item, error)
	FindByMainSkus(skus []string) (map[string]models.Item, error)
}
