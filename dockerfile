
FROM golang:1.22.6-bookworm

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "cmd/prendeluz_erp/main.go", "--env", "test"]
