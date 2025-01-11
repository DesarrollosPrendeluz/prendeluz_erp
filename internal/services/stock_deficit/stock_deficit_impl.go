package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"strconv"

	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"

	"github.com/xuri/excelize/v2"
)

type ParentItemResult struct {
	ParentItemID int    `gorm:"column:parent_item_id"`
	MainSKU      string `gorm:"column:main_sku"`
}
type StockDeficitResult struct {
	Deficit float64 `gorm:"column:deficit"`
}

type StockDeficitServiceImpl struct {
	stockDeficitRepo stockdeficitrepo.StockDeficitImpl
	itemsRepo        itemsrepo.ItemRepoImpl
	orderErrorRepo   repositories.GORMRepository[models.ErrorOrder]
}

func NewStockDeficitService() *StockDeficitServiceImpl {
	stockDeficitRepo := *stockdeficitrepo.NewStockDeficitRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &StockDeficitServiceImpl{
		stockDeficitRepo: stockDeficitRepo,
		orderErrorRepo:   errorOrderRepo,
		itemsRepo:        itemsRepo}
}

func (s *StockDeficitServiceImpl) SearchBySkuAndEan(filter string, store int, page int, pageSize int) ([]models.StockDeficit, []error) {

	var fatherSkus []string
	var errArray []error
	//subQuery := s.stockDeficitRepo.
	items, err1 := s.itemsRepo.FindByFathersMainSkuOrEan(filter)
	for _, item := range items {
		fatherSkus = append(fatherSkus, item.MainSKU)
	}

	stockDef, err2 := s.stockDeficitRepo.GetByRegsitersByFatherSkuIn(fatherSkus, store, page, pageSize)

	if err1 != nil || err2 != nil {
		errArray = append(errArray, err1)
		errArray = append(errArray, err2)
		return nil, errArray
	}
	return stockDef, errArray

}

func (s *StockDeficitServiceImpl) AddPendingStockByItem(child_item_id uint64, store_id int64, quantity int) error {
	//TODO: Refactorizar este método hay que separar la lógica de la consulta de la lógica de la actualización

	result := returnParentItemById(child_item_id)

	stockDef, err := s.stockDeficitRepo.FindOrCreateByFatherAndStore(result.MainSKU, store_id)

	if err != nil {
		fmt.Errorf("error al buscar registro existente: %w", err)
		return err
	}

	stockDef.PendingAmount += int64(quantity)

	err2 := s.stockDeficitRepo.Update(&stockDef)
	if err2 != nil {
		fmt.Errorf("error al actualizar: %w", err)
		return err
	}
	return nil

}

func (s *StockDeficitServiceImpl) CalcStockDeficitByFatherOrder(father_order_id uint64) {
	fmt.Println("hola pasa por aqui")
	data, err := s.stockDeficitRepo.StockDeficitByFatherOrder(father_order_id)
	fmt.Println(len(data))
	if data != nil && err == nil {
		for _, datum := range data {
			result := returnParentItemById(datum.ItemID)
			existing, err3 := s.stockDeficitRepo.FindOrCreateByFatherAndStore(result.MainSKU, 2)
			if err3 != nil {
				fmt.Errorf("error al buscar registro existente: %w", err3)
				//return
			}
			// El registro ya existe, realizar una actualización
			existing.Amount = existing.Amount + int64(datum.Deficit)
			s.stockDeficitRepo.Update(&existing)

		}
	} else {
		fmt.Printf("Error al ejecutar la consulta: %v", err)
	}
}

func (s *StockDeficitServiceImpl) CalcStockDeficitByItem(child_item_id uint64, store_id int64) {
	//TODO: Refactorizar este método hay que separar la lógica de la consulta de la lógica de la actualización
	var existing models.StockDeficit
	var deficit StockDeficitResult
	var pending StockDeficitResult

	result := returnParentItemById(child_item_id)

	err1 := s.stockDeficitRepo.DB.
		Table("order_lines AS ol").
		Select("GREATEST(0, -(IFNULL(AVG(ss.quantity), 0) - (SUM(ol.quantity) - SUM(ol.recived_quantity)))) AS deficit").
		Joins("LEFT JOIN item_parents ip ON ip.child_item_id = ol.item_id").
		Joins("LEFT JOIN orders ord ON ord.id = ol.order_id").
		Joins("INNER JOIN father_orders fo ON fo.id = ord.father_order_id AND fo.order_type_id = 2 AND fo.order_status_id != 3").
		Joins("LEFT JOIN store_stocks ss ON ss.parent_main_sku = ?", result.MainSKU).
		Where("ip.parent_item_id = ?", result.ParentItemID).
		Where("ol.store_id = ?", store_id).
		Group("ip.parent_item_id").
		Take(&deficit).Error
	if err1 != nil {
		fmt.Printf("Error al ejecutar la consulta: %v", err1)
	}
	err2 := s.stockDeficitRepo.DB.
		Table("order_lines AS ol").
		Select(" SUM(ol.quantity) - SUM(ol.recived_quantity) AS deficit").
		Joins("LEFT JOIN item_parents ip ON ip.child_item_id = ol.item_id").
		Joins("LEFT JOIN orders ord ON ord.id = ol.order_id").
		Joins("INNER JOIN father_orders fo ON fo.id = ord.father_order_id AND fo.order_type_id = 1 AND fo.order_status_id != 3").
		Where("ip.parent_item_id = ?", result.ParentItemID).
		Where("ol.store_id = ?", store_id).
		Group("ip.parent_item_id").
		Take(&pending).Error

	if err2 != nil {
		fmt.Printf("Error al ejecutar la consulta: %v", err2)
	}

	existing, err3 := s.stockDeficitRepo.FindOrCreateByFatherAndStore(result.MainSKU, store_id)
	if err3 != nil {
		fmt.Errorf("error al buscar registro existente: %w", err3)
		return
	}

	// El registro ya existe, realizar una actualización
	existing.Amount = int64(deficit.Deficit)
	existing.PendingAmount = int64(pending.Deficit)
	s.stockDeficitRepo.Update(&existing)

}

func returnParentItemById(id uint64) (parent ParentItemResult) {
	//TODO: Refactorizar este método hay que separar la lógica de la consulta de la lógica de la actualización
	var result ParentItemResult

	item, _ := itemsrepo.NewItemRepository(db.DB).FindByID(id)

	if item.ItemType == "father" {
		result.MainSKU = item.MainSKU
		result.ParentItemID = int(item.ID)

	} else {
		parent, err := itemsparentsrepo.NewItemParentRepository(db.DB).FindByChild(id)
		if err != nil {
			fmt.Printf("Error al ejecutar la consulta: %v", err)
		}
		result.MainSKU = parent.Parent.MainSKU
		result.ParentItemID = int(parent.Parent.ID)

	}

	return result

}

func (s *StockDeficitServiceImpl) ReturnDownloadStockDeficitExcel(store_id int) string {
	//s.stockRepo.FindByStore(store_id);
	stockDeficits, _ := s.stockDeficitRepo.GetallByStore(2, -1, -1)
	f := excelize.NewFile()

	// Inicia en la fila 2 para Locations

	// Crear encabezados en la primera fila
	sheetNameTotals := "Stock Deficit"

	f.NewSheet(sheetNameTotals)
	//f.SetCellValue(sheetNameTotals, "A1", "Name")
	f.SetCellValue(sheetNameTotals, "A1", "MainSku")
	f.SetCellValue(sheetNameTotals, "B1", "Ean")
	f.SetCellValue(sheetNameTotals, "C1", "Total")
	f.SetCellValue(sheetNameTotals, "D1", "ToRecive")
	f.SetCellValue(sheetNameTotals, "E1", "ToOrder")

	for totalIndex, deficit := range stockDeficits {
		totalRow := totalIndex + 2
		toOrder := deficit.Amount - deficit.PendingAmount
		if toOrder < 0 {
			toOrder = 0
		}
		//f.SetCellValue(sheetNameTotals, "A"+strconv.Itoa(totalRow), deficit.Item.Name)
		f.SetCellValue(sheetNameTotals, "A"+strconv.Itoa(totalRow), deficit.Item.MainSKU)
		f.SetCellValue(sheetNameTotals, "B"+strconv.Itoa(totalRow), deficit.Item.EAN)
		f.SetCellValue(sheetNameTotals, "C"+strconv.Itoa(totalRow), deficit.Amount)
		f.SetCellValue(sheetNameTotals, "D"+strconv.Itoa(totalRow), deficit.PendingAmount)
		f.SetCellValue(sheetNameTotals, "E"+strconv.Itoa(totalRow), toOrder)

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
