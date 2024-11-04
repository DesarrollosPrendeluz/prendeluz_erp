#!/bin/sh

# Instalar dependencias si es necesario
#go mod download

# Ordenar y limpiar las dependencias
go mod tidy

# Ejecutar la aplicación con el entorno de prueba
go run cmd/prendeluz_erp/main.go --env test