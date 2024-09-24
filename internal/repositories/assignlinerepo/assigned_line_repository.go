package assignlinerepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type AssignedLineRepoRepo interface {
	repositories.Repository[models.AssignedLine]
}
