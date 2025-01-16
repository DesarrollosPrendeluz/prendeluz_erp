package utils

import (
	"bytes"
	"encoding/base64"

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
}

func FrameGenerator(sheet string, fields map[string]string, name string) (string, string) {
	fileName := name + ".xlsx"
	f := excelize.NewFile()
	if sheet != "" {
		f.NewSheet(sheet)
	} else {
		sheet = f.GetSheetName(f.GetActiveSheetIndex())
	}

	// Crear encabezados en la primera fila
	for key, field := range fields {
		if err := f.SetCellValue(sheet, key, field); err != nil {
			return "", ""
		}
	}
	f.DeleteSheet("Sheet1")
	// Escribir el archivo Excel en un buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", ""
	}

	// Codificar el contenido del buffer en Base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64String, fileName
}
