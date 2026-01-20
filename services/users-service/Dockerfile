FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
WORKDIR /app/services/users-service
RUN CGO_ENABLED=0 GOOS=linux go build -o users-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/users-service/users-service .

EXPOSE 8081

CMD ["./users-service"]
