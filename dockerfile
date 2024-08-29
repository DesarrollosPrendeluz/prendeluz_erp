# Dockerfile

# Usar la imagen base de Golang
FROM golang:1.22.6-bookworm

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el código fuente a /app
COPY ./src /app

# Copiar el script de entrada al contenedor
COPY entrypoint.sh /app/entrypoint.sh

# Dar permisos de ejecución al script
RUN chmod +x /app/entrypoint.sh

# Comando por defecto para ejecutar el script
CMD ["/app/entrypoint.sh"]