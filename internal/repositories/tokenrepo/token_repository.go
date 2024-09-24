package tokenrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type AccesTokensRepo interface {
	repositories.Repository[models.AccesTokens]
	CheckCredentials(token string) (bool, error)
	ReturnDataByToken(token string) (models.AccesTokens, error)
	//GenerateUserToken(idStore uint64, pageSize int, offset int) (string, error)
}
