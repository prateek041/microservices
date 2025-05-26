# Multi Stage Docker build process
# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o product-catalog ./cmd/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/product-catalog ./
# If there are configuration files
# COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./product-catalog"]
