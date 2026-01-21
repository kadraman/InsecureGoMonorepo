FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
WORKDIR /app/services/products-service
RUN CGO_ENABLED=0 GOOS=linux go build -o products-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/products-service/products-service .

EXPOSE 8082

CMD ["./products-service"]
