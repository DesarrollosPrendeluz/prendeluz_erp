# Prendeluz ERP
New prendeluz erp backend  
*Es necesario tener go instalado con al menos la versión 1.22.5*



## Installation

Install and run in develop prendeluz erp with go

```bash
  go mod tidy
  go run cmd/prendeluz_erp/main.go  --env test
```

## Run tests
Run tests, ejecutar el comando siguiente en el directorio donde esten los test

```bash
  go test 
```
## Run docker image
Revisar la ruta hacia la que hace las peticiones el en o en este caso el archivo dockerfile

```bash
  docker build -t backend-erp-prod .
  docker run -p 8080:8080 -e HOST=0.0.0.0 backend-erp-prod
```
## Upload image to Dockerhub

```bash
  docker tag backend-erp-prod desarrollospr/backend-erp-prod:1.1.0 ... (tags varios)
  docker push desarrollospr/backend-erp-prod:1.1.0  (tags varios)
```

## Upload image to Plesk via Ssh

```bash
  docker login
  docker pull ...
  
```
*En plesk vincular al puerto 8888*
Hay que vincular mediante el docker proxy con el dominio que se desee usar.