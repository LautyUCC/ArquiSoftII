# ============================================
# STAGE 1: BUILD
# ============================================
# Este stage compila la aplicación Go
# Usa una imagen de Go con Alpine para reducir el tamaño
FROM golang:1.21-alpine AS build

# Establecer el directorio de trabajo dentro del contenedor
# Todas las operaciones siguientes se ejecutarán en este directorio
WORKDIR /app

# Copiar go.mod y go.sum primero
# Esto permite aprovechar el cache de Docker si las dependencias no cambian
# Solo si estos archivos cambian, se volverá a ejecutar go mod download
COPY go.mod go.sum ./

# Descargar las dependencias del proyecto
# Esto se ejecuta antes de copiar el código para aprovechar el cache de Docker
# Las dependencias se descargan en el cache de Go del contenedor
RUN go mod download

# Copiar todo el código fuente de la aplicación
# Esto incluye todos los archivos .go y otros recursos necesarios
COPY . .

# Compilar la aplicación Go
# -o main: especifica el nombre del binario de salida
# .: compila el paquete en el directorio actual (main.go)
# El binario resultante será un ejecutable estático de Linux
RUN go build -o main .

# ============================================
# STAGE 2: RUNTIME
# ============================================
# Este stage crea la imagen final de producción
# Usa Alpine Linux, una distribución mínima para reducir el tamaño de la imagen
FROM alpine:latest

# Instalar ca-certificates para permitir conexiones HTTPS/TLS
# Necesario para comunicarse con servicios externos (APIs, MongoDB, RabbitMQ)
# --no-cache: no almacena la lista de paquetes en el cache para reducir tamaño
RUN apk add --no-cache ca-certificates

# Establecer el directorio de trabajo en el contenedor
# El binario se ejecutará desde este directorio
WORKDIR /root/

# Copiar el binario compilado desde el stage de build
# El binario "main" del stage anterior se copia al directorio actual
COPY --from=build /app/main .

# Exponer el puerto 8081
# Este es el puerto donde la aplicación escuchará las peticiones HTTP
# La documentación del Dockerfile indica qué puerto usa la aplicación
EXPOSE 8081

# Comando por defecto para ejecutar la aplicación
# Cuando se inicie el contenedor, se ejecutará el binario ./main
# Esto iniciará el servidor HTTP en el puerto 8081
CMD ["./main"]
