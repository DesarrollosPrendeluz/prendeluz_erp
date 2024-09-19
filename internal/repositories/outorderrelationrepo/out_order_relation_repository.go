package outorderrelationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type InOrderRelationRepo interface {
	repositories.Repository[models.OutOrderRelation]
}
