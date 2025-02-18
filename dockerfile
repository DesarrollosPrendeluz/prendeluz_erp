
FROM golang:1.22.6-bookworm

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
# Define los argumentos de construcción
ARG DB_PSWD
ARG DB_NAME
# Opcional: Pasar los valores de ARG a ENV para que estén disponibles en tiempo de ejecución
ENV SERVER_PORT=8880
ENV DB_HOST="70.35.197.24"
ENV DB_PORT=32991
ENV DB_USERNAME="root"
ENV DB_PASSWORD=${DB_PSWD}
ENV DB_NAME=${DB_NAME}
#RUN echo "Valor de MY_ENV_VAR: $DB_PSWD"

RUN go mod tidy

EXPOSE 8880

CMD ["go", "run", "cmd/prendeluz_erp/main.go", "--env", "vars"]
