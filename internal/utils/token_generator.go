package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) string {
	// Calcular el número de bytes necesarios para obtener la longitud deseada
	// en caracteres base64, considerando que 4 caracteres base64 son 3 bytes.
	byteLength := (length * 3) / 4

	// Crear un slice de bytes de longitud calculada
	randomBytes := make([]byte, byteLength)

	// Rellenar el slice con bytes aleatorios
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Manejar el error adecuadamente en producción
	}

	// Codificar los bytes en base64
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// Asegurar que la longitud sea exactamente la deseada
	return randomString[:length]
}
