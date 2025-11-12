# ---- Build ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# Dependencias del sistema que pueden necesitar los módulos
RUN apk add --no-cache git ca-certificates

# Copiá los manifests primero (si tenés go.sum también)
COPY go.mod ./
# COPY go.sum ./

# Copiá el código antes de tidy para que detecte imports reales
COPY . .

# Resuelve módulos y genera go.sum
RUN go mod tidy

# Compila binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api .

# ---- Runtime ----
FROM alpine:3.20
WORKDIR /app
COPY --from=build /api /api
EXPOSE 8080
CMD ["/api"]
