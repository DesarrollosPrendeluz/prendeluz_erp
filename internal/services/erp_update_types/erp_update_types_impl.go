package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/updateerptyperepo"
)

type ErpUpdateTypeImpl struct {
	updateerptyperepo updateerptyperepo.UpdateErpTypeImpl
}

func NewErpUpdateTypeService() *ErpUpdateTypeImpl {
	updateerptyperepo := *updateerptyperepo.NewUpdateErpTypeRepository(db.DB)

	return &ErpUpdateTypeImpl{
		updateerptyperepo: updateerptyperepo}
}

func (s *ErpUpdateTypeImpl) GetAll() []models.UpdateErpType {
	data, _ := s.updateerptyperepo.FindAll(-1, -1)
	return data

}
