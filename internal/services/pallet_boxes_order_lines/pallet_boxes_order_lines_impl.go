package services

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/boxrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderlineboxrepo"
	"prendeluz/erp/internal/repositories/palletrepo"
)

type PalletBoxesOrderLinesServiceImpl struct {
	boxRepo          boxrepo.BoxImpl
	palletRepo       palletrepo.PalletImpl
	orderLineBoxRepo orderlineboxrepo.OrderLineBoxImpl
	orderItemRepo    orderitemrepo.OrderItemRepoImpl
}

func NewStockDeficitService() *PalletBoxesOrderLinesServiceImpl {

	boxRepo := *boxrepo.NewBoxRepository(db.DB)
	palletRepo := *palletrepo.NewPalletRepository(db.DB)
	orderLineBoxRepo := *orderlineboxrepo.NewOrderLineBoxRepository(db.DB)
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)

	return &PalletBoxesOrderLinesServiceImpl{
		boxRepo:          boxRepo,
		palletRepo:       palletRepo,
		orderLineBoxRepo: orderLineBoxRepo,
		orderItemRepo:    orderItemRepo,
	}
}

func (s *PalletBoxesOrderLinesServiceImpl) CheckAndCreateBoxOrderLines(orderLineId int, palletNumber int, boxNumber int, quantity int) ([]string, []error) {

	var checks []string
	var errArray []error
	var boxErr error
	orderLine, errOL := s.orderItemRepo.FindByID(uint64(orderLineId))

	if orderLine != nil && errOL == nil {
		pallet, create, errPallet := s.palletRepo.GetOrCreatePalletByOrderIdAndNumber(int(orderLine.OrderID), palletNumber)
		checksAndErrors(create, errPallet, int(pallet.ID), "el palet", &checks, &errArray)

		box, create2, errBox := s.boxRepo.GetOrCreateBoxByPalletIdAndNumber(int(pallet.ID), boxNumber, quantity)
		checksAndErrors(create2, errBox, int(box.ID), "la caja", &checks, &errArray)
		if !create2 {
			box.Quantity = (box.Quantity + quantity)
			if box.Quantity >= 0 {
				s.boxRepo.Update(&box)
			} else {
				boxErr = fmt.Errorf("la cantidad de cajas no puede ser negativa")
				checksAndErrors(false, boxErr, int(box.ID), "la caja", &checks, &errArray)
			}

		}

		olBox, create3, errOlBox := s.orderLineBoxRepo.GetOrCreateByOrderLineAndBoxId(int(orderLine.ID), int(box.ID), quantity)
		checksAndErrors(create3, errOlBox, olBox.ID, "la relación OL box", &checks, &errArray)
		if !create3 && boxErr == nil {
			olBox.Quantity = (olBox.Quantity + quantity)
			s.orderLineBoxRepo.Update(&olBox)
		}

	} else {
		checks = append(checks, "Error en la consulta del order line")
		errArray = append(errArray, errOL)
	}

	return checks, errArray

}

func checksAndErrors(create bool, err error, id int, name string, checks *[]string, errArray *[]error) {
	if create {
		*checks = append(*checks, "Se ha creado "+name+" "+fmt.Sprint(id))
	}
	if err != nil {
		*errArray = append(*errArray, err)
	}

}
