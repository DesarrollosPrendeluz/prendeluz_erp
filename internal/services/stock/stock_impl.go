package services

import (
	"bytes"
	"encoding/base64"
	"prendeluz/erp/internal/db"
	dtos "prendeluz/erp/internal/dtos/api"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"prendeluz/erp/internal/repositories/asinrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type StockServiceImpl struct {
	stockRepo      storestockrepo.StoreStockRepoImpl
	asinRepo       asinrepo.AsinRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewStockService() *StockServiceImpl {
	stockRepo := *storestockrepo.NewStoreStockRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})

	return &StockServiceImpl{
		stockRepo:      stockRepo,
		orderErrorRepo: errorOrderRepo}
}

func (s *StockServiceImpl) FreeReservedStock(quantity int64, parent_sku string) error {
	itemStock, _ := s.stockRepo.FindByItemAndStore(parent_sku, "1")

	itemStock.ReservedAmount = itemStock.ReservedAmount - quantity

	err := s.stockRepo.Update(&itemStock)
	return err
}

func (s *StockServiceImpl) ReturnDownloadStockExcel(store_id int) string {
	//s.stockRepo.FindByStore(store_id);
	stocks, _ := s.stockRepo.FindByStoreWithLocations(uint64(store_id))
	f := excelize.NewFile()

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

	for totalIndex, stock := range stocks.StoreStocks {
		totalRow := totalIndex + 2
		f.SetCellValue(sheetNameTotals, "A"+strconv.Itoa(totalRow), stock.SKU_Parent)
		f.SetCellValue(sheetNameTotals, "B"+strconv.Itoa(totalRow), stock.Item.EAN)
		f.SetCellValue(sheetNameTotals, "C"+strconv.Itoa(totalRow), stock.ReservedAmount)
		f.SetCellValue(sheetNameTotals, "D"+strconv.Itoa(totalRow), stock.Amount)
	}

	for locIndex, location := range stocks.ItemsLocation {

		totalRow := locIndex + 2
		f.SetCellValue(sheetNamePartials, "A"+strconv.Itoa(totalRow), location.ItemSku)
		f.SetCellValue(sheetNamePartials, "B"+strconv.Itoa(totalRow), location.Ean)
		f.SetCellValue(sheetNamePartials, "C"+strconv.Itoa(totalRow), location.Code)
		f.SetCellValue(sheetNamePartials, "D"+strconv.Itoa(totalRow), location.Stock)

	}

	f.DeleteSheet("Sheet1")
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {

		return ""
	}

	// Codificar el contenido del buffer en Base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String

}

func (s *StockServiceImpl) ReturnStockByAsins(asins []string) []dtos.StockItem {
	// Declaración de variables
	var stockItems []dtos.StockItem
	var itemIds []uint64
	var fatherSkus []string
	itemToAsin := make(map[uint64]string)
	skuToAsin := make(map[string]string)
	// Obtener los datos de los ASINs
	ainsData, err := s.asinRepo.FindByAsins(asins)
	if err != nil {
		return nil
	}
	// Obtener los IDs de los items relacionados con los ASINs encontrados
	for _, asin := range ainsData {
		itemIds = append(itemIds, asin.ItemID)
		itemToAsin[asin.ItemID] = asin.Code
	}
	mixedItems, err := s.itemsRepo.FindByIds(itemIds)
	if err != nil {
		return nil
	}
	// Obtener los SKUs padre de los items encontrados
	for _, item := range mixedItems {
		if item.ItemType == "father" {
			fatherSkus = append(fatherSkus, item.MainSKU)
			skuToAsin[item.MainSKU] = itemToAsin[item.ID]
		} else {
			ItemWithFather, err := s.itemsRepo.FindByIdWithFatherPreload(item.ID)
			if err != nil {
				return nil
			}
			fatherSkus = append(fatherSkus, ItemWithFather.FatherRel.Parent.MainSKU)
			skuToAsin[ItemWithFather.FatherRel.Parent.MainSKU] = itemToAsin[ItemWithFather.ID]
		}
	}
	// Obtener el stock de los SKUs padre para la tienda específica
	stockData, err := s.stockRepo.FindByParentSkusAndStore(fatherSkus, 1)
	if err != nil {
		return nil
	}
	for _, stock := range stockData {
		var stockItem dtos.StockItem
		reserved := stock.ReservedAmount
		if reserved < 0 {
			reserved = 0
		}
		aviableStock := stock.Amount - reserved
		if aviableStock < 0 {
			aviableStock = 0
		}

		stockItem.StoreID = stock.StoreID
		stockItem.ASIN = skuToAsin[stock.SKU_Parent]
		stockItem.Quantity = int(aviableStock)

		stockItems = append(stockItems, stockItem)
	}

	return stockItems

}
