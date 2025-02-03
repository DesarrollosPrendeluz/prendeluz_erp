package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/palletrepo"
)

type PalletImpl struct {
	palletrepo palletrepo.PalletImpl
}

func NewPalletService() *PalletImpl {
	palletrepo := *palletrepo.NewPalletRepository(db.DB)

	return &PalletImpl{
		palletrepo: palletrepo}
}

func (s *PalletImpl) Get(pallet int, page int, pageSize int) ([]models.Pallet, int64, error) {
	var err error
	var data []models.Pallet
	var datum *models.Pallet
	var recount int64

	if pallet != 0 {
		datum, err = s.palletrepo.FindByID(uint64(pallet))
		if datum != nil { // Verificar si datum no es nil
			data = append(data, *datum)
		}
		recount = 1
	} else {
		data, err = s.palletrepo.FindAll(pageSize, page)
		recount, _ = s.palletrepo.CountAll()

	}
	return data, recount, err

}

func (s *PalletImpl) GetPalletByOrder(orderId int, page int, pageSize int) ([]models.Pallet, int64, error) {
	var recount int64
	data, err := s.palletrepo.GetBoxesAndLinesRaletedDataByOrderId(orderId, pageSize, page)

	recount = 1
	return data, recount, err

}

func (s *PalletImpl) Create(data dtos.PalletCreateReq) []error {
	var errorList []error
	for _, dataItem := range data.Data {

		model := models.Pallet{
			OrderID: dataItem.OrderID,
			Number:  int(dataItem.Number),
			Label:   dataItem.Label,
		}

		error := s.palletrepo.Create(&model)
		if error != nil {
			errorList = append(errorList, error)
		}
	}
	return errorList
}

func (s *PalletImpl) Update(data dtos.PalletUpdateReq) []error {
	var errorList []error
	for _, requestObject := range data.Data {
		model, err := s.palletrepo.FindByID(requestObject.Id)
		if err != nil {
			errorList = append(errorList, err)
			return errorList
		}
		if requestObject.OrderID != nil {
			model.OrderID = *requestObject.OrderID
		}
		if requestObject.Label != nil {
			model.Label = *requestObject.Label
		}
		if requestObject.Number != nil {
			model.Number = *requestObject.Number
		}
		if requestObject.IsClose != nil {
			if *requestObject.IsClose {
				model.IsClose = 1
			} else {
				model.IsClose = 0
			}
		}

		error := s.palletrepo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	return errorList
}
