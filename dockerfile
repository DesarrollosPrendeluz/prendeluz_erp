# Usar la imagen base de Golang
FROM golang:1.22.6-bookworm

# Instalar Git y otros paquetes necesarios
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Establecer el directorio de trabajo en /app dentro del contenedor
WORKDIR /app

# Copiar archivos de dependencias desde la raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente de todas las carpetas necesarias al contenedor
COPY cmd/ /app/cmd
COPY internal/ /app/internal
COPY test/ /app/test

# Copiar el script de entrada si lo tienes
COPY entrypoint.sh /app/entrypoint.sh
#RUN chmod +x /app/entrypoint.sh

# Compilar la aplicación apuntando a la carpeta cmd
#RUN go build -o main 

# Exponer el puerto en el que la aplicación escuchará
EXPOSE 8080

# Definir el comando de entrada
CMD ["/app/entrypoint.sh"]
