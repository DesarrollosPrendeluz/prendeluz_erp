package utils

import (
	"github.com/xuri/excelize/v2"
)

var NewOrderSheetName = "OC SQL"
var NewOrder = map[string]string{
	"A3": "Orden de compra",
	"B3": "Asin",
	"C3": "Sku",
	"D3": "Amount",
	"E3": "Parent sku",
	"G3": "Store",
	"H3": "Client",
}
var UploadStockSheetName = "Stock Modify"
var UploadStock = map[string]string{
	"A1": "Sku",
	"B1": "Loc_code",
	"C1": "Quantity",
}
var ModifyOrderSheetName = "Modify Order"
var ModifySuppOrder = map[string]string{
	"A1": "Sku",
	"B1": "Quantity",
	"C1": "Update_Reason_Id",
}

var ModifyOrder = map[string]string{
	"A1": "Sku",
	"B1": "Quantity",
	"C1": "Update_Reason_Id",
	"D1": "Order Code",
}

func FrameGenerator(sheet string, fields map[string]string, name string) (string, string) {
	fileName := name + ".xlsx"
	f := excelize.NewFile()
	callback := func(f *excelize.File, sheetName string) *excelize.File {
		return f
	}
	genericSheetCreator(f, sheet, fields, callback)

	return base64ExcelEncoder(f), fileName
}
