package userrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type UserRepo interface {
	repositories.Repository[models.User]
	CheckCredentials(parent_sku string) (models.User, error)
	GenerateUserToken(idStore uint64, pageSize int, offset int) (string, error)
}
