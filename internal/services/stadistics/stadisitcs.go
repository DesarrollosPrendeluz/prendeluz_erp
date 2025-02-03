package services

import (
	"prendeluz/erp/internal/dtos"
)

type StadisitcsService interface {
	GetChangeStadistics(fatherId uint64) dtos.HistoricStats
	GetRecivedStadistics(fatherId uint64) dtos.RecivedHistory
}
