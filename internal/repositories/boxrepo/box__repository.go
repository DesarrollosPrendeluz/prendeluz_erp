package boxrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type BoxRepo interface {
	repositories.Repository[models.Box]
}
