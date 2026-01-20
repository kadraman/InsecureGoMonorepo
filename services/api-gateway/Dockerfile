FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/api-gateway/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
