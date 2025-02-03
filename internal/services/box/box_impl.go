package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/boxrepo"
)

type BoxImpl struct {
	boxRepo boxrepo.BoxImpl
}

func NewBoxService() *BoxImpl {
	boxRepo := *boxrepo.NewBoxRepository(db.DB)

	return &BoxImpl{
		boxRepo: boxRepo}
}

func (s *BoxImpl) GetBox(box int, page int, pageSize int) ([]models.Box, int64, error) {
	var err error
	var data []models.Box
	var datum *models.Box
	var recount int64

	if box != 0 {
		datum, err = s.boxRepo.FindByID(uint64(box))
		if datum != nil { // Verificar si datum no es nil
			data = append(data, *datum)
		}
		recount = 1
	} else {
		data, err = s.boxRepo.FindAll(pageSize, page)
		recount, _ = s.boxRepo.CountAll()

	}

	if err != nil {
		return nil, 0, err

	}
	return data, recount, nil

}

func (s *BoxImpl) CreateBox(data dtos.BoxCreateReq) []error {
	var errorList []error
	for _, dataItem := range data.Data {
		model := models.Box{
			PalletID: dataItem.PalletID,
			Number:   int(dataItem.Number),
			Label:    dataItem.Label,
			Quantity: dataItem.Quantity,
		}
		error := s.boxRepo.Create(&model)
		if error != nil {
			errorList = append(errorList, error)
		}
	}
	return errorList
}

func (s *BoxImpl) UpdateBox(data dtos.BoxUpdateReq) []error {
	var errorList []error
	for _, dataItem := range data.Data {
		model, err := s.boxRepo.FindByID(dataItem.Id)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		if dataItem.PalletID != nil {
			model.PalletID = *dataItem.PalletID
		}
		if dataItem.Label != nil {
			model.Label = *dataItem.Label
		}
		if dataItem.Number != nil {
			model.Number = int(*dataItem.Number)
		}

		if dataItem.Quantity != nil {
			model.Quantity = *dataItem.Quantity
		}

		if dataItem.IsClose != nil {
			if *dataItem.IsClose {
				model.IsClose = 1
			} else {
				model.IsClose = 0
			}

		}

		error := s.boxRepo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}
	}
	return errorList
}
