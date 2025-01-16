package utils

import (
	"bytes"
	"encoding/base64"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type returnDoc func(*excelize.File, string) *excelize.File

type StockUpdateError struct {
	FatherSku string
	Loc       string
	Error     string
}
type UpdateOrderError struct {
	FatherSku string
	Error     string
}

var UpdateStockErr = map[string]string{
	"A1": "Sku",
	"B1": "Loc",
	"C1": "Err",
}
var UpdateOrderErr = map[string]string{
	"A1": "Sku",
	"B1": "Err",
}

func ReturnUpdateErrorsExcel(data []StockUpdateError) string {
	//s.stockRepo.FindByStore(store_id);
	f := excelize.NewFile()
	callback := func(f *excelize.File, sheetName string) *excelize.File {
		for totalIndex, datum := range data {
			totalRow := totalIndex + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(totalRow), datum.FatherSku)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(totalRow), datum.Loc)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(totalRow), datum.Error)
		}
		return f
	}
	genericSheetCreator(f, "Update Errors", UpdateStockErr, callback)
	return base64ExcelEncoder(f)
}

func ReturnUpdateOrdersErrorsExcel(data []UpdateOrderError) string {

	f := excelize.NewFile()
	callback := func(f *excelize.File, sheetName string) *excelize.File {
		for totalIndex, datum := range data {
			totalRow := totalIndex + 2
			f.SetCellValue(sheetName, "A"+strconv.Itoa(totalRow), datum.FatherSku)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(totalRow), datum.Error)
		}
		return f
	}
	genericSheetCreator(f, "Update Errors", UpdateOrderErr, callback)
	return base64ExcelEncoder(f)

}

func genericSheetCreator(file *excelize.File, sheet string, fields map[string]string, completeData returnDoc) *excelize.File {

	file.NewSheet(sheet)
	// Crear encabezados en la primera fila
	for key, field := range fields {
		if err := file.SetCellValue(sheet, key, field); err != nil {
		}
	}

	completeData(file, sheet)
	return file

}

func base64ExcelEncoder(f *excelize.File) string {
	f.DeleteSheet("Sheet1")
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {

		return ""
	}

	// Codificar el contenido del buffer en Base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String
}
