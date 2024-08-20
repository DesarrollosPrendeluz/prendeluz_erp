package utils

import (
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type OrderInfo struct {
	ASIN      string
	MainSku   string
	Amount    int64
	ParentSku string
}

type ExcelOrder struct {
	OrderCode string
	Info      []OrderInfo
}

func ExceltoJSON(file io.Reader) ([]ExcelOrder, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	orders := make(map[string][]OrderInfo)
	rows, err := f.GetRows("OC SQL")

	if err != nil {
		return nil, err
	}

	for _, row := range rows[3:] {
		if len(row) > 4 {
			amount, _ := strconv.ParseInt(row[3], 10, 64)
			item := OrderInfo{
				ASIN:      row[1],
				MainSku:   row[2],
				Amount:    amount,
				ParentSku: row[4],
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
