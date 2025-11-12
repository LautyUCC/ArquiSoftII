# Build stage
FROM golang:1.22-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api .

# Runtime stage
FROM alpine:3.20

COPY --from=build /api /bin/api

EXPOSE 8082

ENTRYPOINT ["/bin/api"]