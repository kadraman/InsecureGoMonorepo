FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
WORKDIR /app/services/orders-service
RUN CGO_ENABLED=0 GOOS=linux go build -o orders-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/orders-service/orders-service .

EXPOSE 8083

CMD ["./orders-service"]
