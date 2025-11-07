# Use golang:1.21-alpine as base image
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o search-api .

# Expose port 8082
EXPOSE 8082

# Run the binary
CMD ["./search-api"]
