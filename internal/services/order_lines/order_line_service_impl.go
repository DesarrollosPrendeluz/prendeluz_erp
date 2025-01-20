package services

import (
	"bytes"
	"encoding/base64"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/tokenrepo"
	stockDeficit "prendeluz/erp/internal/services/stock_deficit"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type OrderLineServiceImpl struct {
	orderRepo      orderrepo.OrderRepoImpl
	orderItemsRepo orderitemrepo.OrderItemRepoImpl
	//itemRepo       itemsrepo.ItemRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewOrderLineServiceImpl() *OrderLineServiceImpl {
	orderRepo := *orderrepo.NewOrderRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &OrderLineServiceImpl{
		orderRepo:      orderRepo,
		orderItemsRepo: orderItemRepo,
		orderErrorRepo: errorOrderRepo,
		itemsRepo:      itemsRepo,
	}
}

func (s *OrderLineServiceImpl) OrderLineLabel(id int) (dtos.OrderLineLable, error) {
	var orderItem models.OrderItem
	//var item models.Item
	var label dtos.OrderLineLable
	//itemsRepo := itemsrepo.NewItemRepository(db.DB)
	s.orderItemsRepo.DB.
		Preload("Item").
		Preload("Item.AsinRel.Brand").
		Preload("Item.FatherRel").
		Where("id=?", id).
		First(&orderItem)

	// s.itemsRepo.DB.
	// 	Preload("SupplierItems").
	// 	Preload("SupplierItems.Brand").
	// 	Where("id=?", orderItem.Item.FatherRel.ParentItemID).
	// 	First(&item)

	if orderItem.Item.AsinRel != nil {
		label.Ean = orderItem.Item.AsinRel.Ean
		label.Asin = &orderItem.Item.AsinRel.Code
		label.Brand = orderItem.Item.AsinRel.Brand.Name
		label.BrandAddress = orderItem.Item.AsinRel.Brand.Address
		label.BrandEmail = orderItem.Item.AsinRel.Brand.Email
		label.Company = orderItem.Item.AsinRel.Brand.Company
	}

	// for _, supplierItem := range *item.SupplierItems {

	// 	if supplierItem.Order == 1 {
	// 		label.Brand = supplierItem.Brand.Name
	// 		label.BrandAddress = supplierItem.Brand.Address
	// 		label.BrandEmail = supplierItem.Brand.Email
	// 		label.Company = supplierItem.Brand.Company
	// 		break
	// 	}
	// }

	return label, nil
}

func (s *OrderLineServiceImpl) ReturnDownloadPickingExcel(data dtos.FatherOrderOrdersAndLines) string {
	//s.stockRepo.FindByStore(store_id);
	//stockDeficits, _ := s.stockDeficitRepo.GetallByStore(2, -1, -1)
	f := excelize.NewFile()

	// Inicia en la fila 2 para Locations

	// Crear encabezados en la primera fila
	sheetNameTotals := "Stock Deficit"

	f.NewSheet(sheetNameTotals)
	//f.SetCellValue(sheetNameTotals, "A1", "Name")
	f.SetCellValue(sheetNameTotals, "A1", "MainSku")
	f.SetCellValue(sheetNameTotals, "B1", "Ean")
	f.SetCellValue(sheetNameTotals, "C1", "Nombre")
	f.SetCellValue(sheetNameTotals, "D1", "RefProveedor")
	f.SetCellValue(sheetNameTotals, "E1", "Proveedor")
	f.SetCellValue(sheetNameTotals, "F1", "Total")
	f.SetCellValue(sheetNameTotals, "G1", "Procesado")
	f.SetCellValue(sheetNameTotals, "H1", "Responsable")
	f.SetCellValue(sheetNameTotals, "I1", "Ubicaciones")

	for totalIndex, datum := range data.Lines {
		totalRow := totalIndex + 2
		f.SetCellValue(sheetNameTotals, "A"+strconv.Itoa(totalRow), datum.MainSku)
		f.SetCellValue(sheetNameTotals, "B"+strconv.Itoa(totalRow), datum.Ean)
		f.SetCellValue(sheetNameTotals, "C"+strconv.Itoa(totalRow), datum.Name)
		f.SetCellValue(sheetNameTotals, "D"+strconv.Itoa(totalRow), datum.SupplierRef)
		f.SetCellValue(sheetNameTotals, "E"+strconv.Itoa(totalRow), datum.SupplierName)
		f.SetCellValue(sheetNameTotals, "F"+strconv.Itoa(totalRow), datum.Quantity)
		f.SetCellValue(sheetNameTotals, "G"+strconv.Itoa(totalRow), datum.RecivedQuantity)
		f.SetCellValue(sheetNameTotals, "H"+strconv.Itoa(totalRow), datum.AssignedUser.UserName)
		f.SetCellValue(sheetNameTotals, "I"+strconv.Itoa(totalRow), datum.Location)

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
func (s *OrderLineServiceImpl) UpdateOrderLineHandler(

	c *gin.Context,
	requestBody dtos.OrdersLinesToUpdatePartially,
	token string,
	failedIds *[]int,
	errorList *[]error,
	callback func(*gin.Context, dtos.LineToUpdate, *models.OrderItem, error, *[]error),
	admin bool) {

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		*errorList = append(*errorList, err)
		return
	}

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {
		var assign dtos.Assign
		repo := tokenrepo.NewTokenRepository(db.DB)
		user, _ := repo.ReturnDataByToken(token)
		query := `SELECT id FROM assigned_lines WHERE  order_line_id = ? and user_id = ? LIMIT 1`

		err := db.DB.Raw(query, dataItem.Id, user.UserId).Scan(&assign).Error

		if (err != nil || assign.ID == 0) && !admin {
			*failedIds = append(*failedIds, int(dataItem.Id))

		} else {

			updateOrderLine(c, dataItem, errorList, callback)
		}

	}

}

func updateOrderLine(
	c *gin.Context,
	dataItem dtos.LineToUpdate,
	errorList *[]error,
	callback func(*gin.Context, dtos.LineToUpdate, *models.OrderItem, error, *[]error)) {
	orderLines := orderitemrepo.NewOrderItemRepository(db.DB)
	stockService := stockDeficit.NewStockDeficitService()
	model, err := orderLines.FindByID(dataItem.Id)

	callback(c, dataItem, model, err, errorList)
	error := orderLines.Update(model)
	if model.StoreID == 1 {
		stockService.CalcStockDeficitByItem(model.ItemID, model.StoreID)

	}
	if error != nil {
		*errorList = append(*errorList, error)
	}

}
