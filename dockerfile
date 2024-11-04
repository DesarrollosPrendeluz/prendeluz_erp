# Usar la imagen base de Golang
FROM golang:1.22.6-bookworm

# Instalar Git y otros paquetes necesarios
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Establecer el directorio de trabajo en /app dentro del contenedor
WORKDIR /app

# Copiar archivos de dependencias desde la raíz
COPY . .
#RUN chmod +x /app/scripts.sh

# Compilar la aplicación apuntando a la carpeta cmd
#RUN go build -o main 

# Exponer el puerto en el que la aplicación escuchará
EXPOSE 8080

# Definir el comando de entrada
# Ordenar y limpiar las dependencias
RUN go mod tidy

# Ejecutar la aplicación con el entorno de prueba
# RUN go run cmd/prendeluz_erp/main.go --env test
#RUN go run cmd/prendeluz_erp/main.go --env test
#CMD ["scripts.sh"]
CMD ["go", "run", "cmd/prendeluz_erp/main.go", "--env", "test"]
