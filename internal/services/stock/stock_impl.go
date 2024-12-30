package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"strconv"

	"prendeluz/erp/internal/repositories/storestockrepo"

	"github.com/xuri/excelize/v2"
)

type StockServiceImpl struct {
	stockRepo      storestockrepo.StoreStockRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
}

func NewStockService() *StockServiceImpl {
	stockRepo := *storestockrepo.NewStoreStockRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})

	return &StockServiceImpl{
		stockRepo:      stockRepo,
		orderErrorRepo: errorOrderRepo}
}

func (s *StockServiceImpl) ReturnDownloadStockExcel(store_id int) string {
	//s.stockRepo.FindByStore(store_id);
	stocks, _ := s.stockRepo.FindByStoreWithLocations(uint64(store_id))
	f := excelize.NewFile()
	f.DeleteSheet("Sheet1")
	// Inicia en la fila 2 para Locations

	// Crear encabezados en la primera fila
	sheetNameTotals := "Totals"

	f.NewSheet(sheetNameTotals)
	f.SetCellValue(sheetNameTotals, "A1", "Sku Padre")
	f.SetCellValue(sheetNameTotals, "B1", "Ean Proveedor")
	f.SetCellValue(sheetNameTotals, "C1", "Stock Total Reservado")
	f.SetCellValue(sheetNameTotals, "D1", "Stock Total")

	sheetNamePartials := "Locations"
	f.NewSheet(sheetNamePartials)
	f.SetCellValue(sheetNamePartials, "A1", "Sku Padre")
	f.SetCellValue(sheetNamePartials, "B1", "Ean Proveedor")
	f.SetCellValue(sheetNamePartials, "C1", "Codigo Localización")
	f.SetCellValue(sheetNamePartials, "D1", "Stock Localización")

	for totalIndex, stock := range stocks {
		totalRow := totalIndex + 2
		f.SetCellValue(sheetNameTotals, "A"+strconv.Itoa(totalRow), stock.SKU_Parent)
		f.SetCellValue(sheetNameTotals, "B"+strconv.Itoa(totalRow), stock.Item.EAN)
		f.SetCellValue(sheetNameTotals, "C"+strconv.Itoa(totalRow), stock.ReservedAmount)
		f.SetCellValue(sheetNameTotals, "D"+strconv.Itoa(totalRow), stock.Amount)

		for locIndex, location := range *stock.Locations {
			locRow := totalRow + locIndex
			locationCode := ""
			if location.StoreLocations != nil {
				locationCode = location.StoreLocations.Code
			}
			fmt.Printf("Location: %+v\n", location)
			f.SetCellValue(sheetNamePartials, "A"+strconv.Itoa(locRow), stock.SKU_Parent)
			f.SetCellValue(sheetNamePartials, "B"+strconv.Itoa(locRow), stock.Item.EAN)

			f.SetCellValue(sheetNamePartials, "C"+strconv.Itoa(locRow), locationCode)
			f.SetCellValue(sheetNamePartials, "D"+strconv.Itoa(locRow), location.Stock)
		}

	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {

		return ""
	}

	// Codificar el contenido del buffer en Base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String

}
