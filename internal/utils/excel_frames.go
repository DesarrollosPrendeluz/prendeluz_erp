package utils

import (
	"github.com/xuri/excelize/v2"
)

var NewOrderSheetName = "OC SQL"
var NewOrder = map[string]string{
	"A1": "Orden de compra",
	"B1": "Asin",
	"C1": "Sku",
	"D1": "Amount",
	"E1": "Parent sku",
	"G1": "Store",
	"H1": "Client",
}
var UploadStockSheetName = "Stock Modify"
var UploadStock = map[string]string{
	"A1": "Sku",
	"B1": "Loc_code",
	"C1": "Quantity",
}
var ModifyOrderSheetName = "Modify Order"
var ModifyOrder = map[string]string{
	"A1": "Sku",
	"B1": "Quantity",
	"C1": "Update_Reason_Id",
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
