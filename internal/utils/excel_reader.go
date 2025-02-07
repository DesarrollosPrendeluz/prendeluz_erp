package utils

import (
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type OrderInfo struct {
	ASIN      string
	MainSku   string
	Amount    int64
	ParentSku string
	Store     uint64
	Client    uint64
}

type ExcelOrder struct {
	OrderCode string
	Info      []OrderInfo
}

type ExcelModifyOrder struct {
	Sku      string
	Quantity int64
	Type     uint64
}

type ExcelUpdateStocks struct {
	Sku      string
	Loc      string
	Quantity int64
}

func ExceltoJSON(file io.Reader) ([]ExcelOrder, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	orders := make(map[string][]OrderInfo)
	rows, err := f.GetRows(NewOrderSheetName)

	if err != nil {
		return nil, err
	}

	for _, row := range rows[3:] {
		if len(row) > 4 {
			amount, _ := strconv.ParseInt(row[3], 10, 64)
			store, _ := strconv.ParseUint(row[6], 10, 64)
			client, _ := strconv.ParseUint(row[7], 10, 64)

			item := OrderInfo{
				ASIN:      row[1],
				MainSku:   row[2],
				Amount:    amount,
				ParentSku: row[4],
				Store:     store,
				Client:    client,
			}
			// if item.MainSku == "#N/D" || item.MainSku == "#N/A" {
			// 	continue
			// }
			code := row[0]
			orders[code] = append(orders[code], item)
		}
	}
	var result []ExcelOrder
	for code, info := range orders {
		result = append(result, ExcelOrder{
			OrderCode: code,
			Info:      info,
		})
	}

	return result, nil

}

func ExcelToJSONOrder(file io.Reader) ([]ExcelModifyOrder, error) {
	var result []ExcelModifyOrder

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	rows, err := f.GetRows(ModifyOrderSheetName)

	if err != nil {
		return nil, err
	}

	for _, row := range rows[1:] {

		amount, _ := strconv.ParseInt(row[1], 10, 64)
		updateType, _ := strconv.ParseUint(row[2], 10, 64)

		item := ExcelModifyOrder{
			Sku:      strings.Trim(row[0], " "),
			Quantity: amount,
			Type:     updateType}

		result = append(result, item)

	}

	return result, nil

}

func ExcelToJsonUpdateStocks(file io.Reader) ([]ExcelUpdateStocks, error) {
	var result []ExcelUpdateStocks

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	rows, err := f.GetRows(UploadStockSheetName)

	if err != nil {
		return nil, err
	}

	for _, row := range rows[1:] {

		amount, _ := strconv.ParseInt(row[2], 10, 64)

		item := ExcelUpdateStocks{
			Sku:      strings.Trim(row[0], " "),
			Loc:      strings.Trim(row[1], " "),
			Quantity: amount}

		result = append(result, item)

	}

	return result, nil

}
