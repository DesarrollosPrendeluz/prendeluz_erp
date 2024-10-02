package pdf

import (
	"fmt"

	"github.com/go-pdf/fpdf"
)

func generateItemLabel() {
	// Dimensiones de un cuarto de una página A4
	width := 105.0 // Ancho de un cuarto de A4 en mm
	height := 48.5 // Alto de un cuarto de A4 en mm

	// Crear una nueva instancia de FPDF con el tamaño ajustado
	pdf := fpdf.NewCustom(&fpdf.InitType{
		UnitStr: "mm",                                 // Unidades en milímetros
		Size:    fpdf.SizeType{Wd: width, Ht: height}, // Tamaño de un cuarto de A4
	})

	// Agregar una página
	pdf.AddPage()

	// Configurar la fuente
	pdf.SetFont("Arial", "B", 16)

	// Escribir texto en el PDF
	pdf.CellFormat(10, 7, "Texto 1", "", 1, "L", false, 0, "")
	pdf.CellFormat(10, 7, "Texto 2", "", 1, "L", false, 0, "")
	pdf.CellFormat(10, 6, "Texto 3", "", 1, "L", false, 0, "") // Ahora se pasa a la siguiente línea

	// Guardar el archivo PDF en el directorio de destino
	outputPath := `C:\Users\Bruno\Desktop\prueba_gofpdf_cuarto_A24.pdf`
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		fmt.Println("Error al generar el PDF:", err)
		return
	}

	fmt.Println("PDF generado correctamente en:", outputPath)
}
