package print

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func printerList() ([]string, error) {
	// Ejecutar el comando wmic para listar impresoras
	cmd := exec.Command("wmic", "printer", "get", "name")

	// Capturar la salida del comando
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar el comando: %v", err)
	}

	// Convertir la salida en string y eliminar espacios y saltos de línea
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// Crear un slice para almacenar los nombres de las impresoras
	var impresoras []string

	// Recorrer las líneas y agregar los nombres de las impresoras al slice, ignorando la primera línea (encabezado)
	for _, line := range lines[1:] { // Se ignora la primera línea porque es el encabezado "Name"
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 {
			impresoras = append(impresoras, trimmedLine)
		}
	}

	return impresoras, nil
}

func powerShellTextPrint(filePath string, printerName string) error {
	// Cambiar el comando para usar notepad como el programa para manejar el archivo de texto
	cmd := exec.Command("PowerShell", "-Command", fmt.Sprintf("Start-Process -FilePath 'SumatraPdf' -ArgumentList '/p', '%s' -NoNewWindow", filePath))

	// Capturar la salida de error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar PowerShell: %v, detalles: %s", err, stderr.String())
	}
	return nil
}

// requiere tener instaldas las impresoras
func sumatraPdfPrint(pdfPath string, printerName string) error {
	// Ruta completa al ejecutable de SumatraPDF
	sumatraPDF := `SumatraPDF`

	// Comando para imprimir el PDF usando SumatraPDF
	cmd := exec.Command(sumatraPDF, "-print-to", printerName, pdfPath)

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error al imprimir el PDF: %v", err)
	}
	return nil
}
