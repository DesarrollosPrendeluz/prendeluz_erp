package controllers

import (
	"net/http"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/orderrepo"
	"strconv"

	"github.com/xuri/excelize/v2"

	"github.com/gin-gonic/gin"
)

func GetSupplierOrders(c *gin.Context) {
	var status *int
	repo := orderrepo.NewOrderRepository(db.DB)
	statusStr := c.Query("status")
	if statusStr != "" {
		// Convierte el string a un int
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			// Si hay un error en la conversión, responde con un error de Bad Request
			c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'status' debe ser un número entero."})
			return
		}
		// Asigna el valor a la variable status como puntero
		status = &statusInt
	}
	data, err := repo.GetSupplierOrders(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": data})
}
func DownloadSupplierOrderExcel(c *gin.Context) {
	// Ejemplo de obtención de datos desde el repositorio
	repo := orderrepo.NewOrderRepository(db.DB)
	var status *int

	// Manejo del parámetro de query `status`
	statusStr := c.Query("status")
	if statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'status' debe ser un número entero."})
			return
		}
		status = &statusInt
	}

	// Obtener datos del repositorio
	data, _ := repo.GetSupplierOrders(status)

	// Crear un nuevo archivo Excel
	f := excelize.NewFile()

	// Crear encabezados en la primera fila
	sheetName := "Sheet1"
	f.SetCellValue(sheetName, "A1", "ID Orden")
	f.SetCellValue(sheetName, "B1", "Nombre")
	f.SetCellValue(sheetName, "C1", "EAN")
	f.SetCellValue(sheetName, "D1", "Proveedor")
	f.SetCellValue(sheetName, "E1", "Código Proveedor")
	f.SetCellValue(sheetName, "F1", "Precio Proveedor")
	f.SetCellValue(sheetName, "G1", "Cantidad")

	// Escribir los datos en las filas siguientes
	for i, datum := range data {
		row := i + 2 // La primera fila es para los encabezados

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), datum.OrderCode)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), datum.Name)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), datum.EAN)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), datum.SupplierName)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), datum.SupplierCode)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), strconv.FormatFloat(datum.SupplierPrice, 'f', 2, 64))
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), datum.StockToBuy)
	}

	// Configura los encabezados HTTP para la descarga del archivo Excel
	c.Header("Content-Disposition", "attachment; filename=supplier_orders.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	// Escribir el archivo Excel en el cuerpo de la respuesta
	if err := f.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Error al generar el archivo Excel")
	}
}
