# Orders Service

[![Orders Service](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/orders-service.yml/badge.svg?branch=main)](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/orders-service.yml)

A deliberately insecure order management microservice for security testing demonstrations.

## Security Vulnerabilities (Intentional)

This service contains the following deliberate security vulnerabilities for testing purposes:

1. **SQL Injection**: Direct string concatenation in all SQL queries
2. **ORDER BY Injection**: Vulnerable sorting functionality
3. **XXE (XML External Entity)**: Unsafe XML parsing in import endpoint
4. **Information Disclosure**: Exposes sensitive data through export endpoint
5. **No Input Validation**: Accepts arbitrary values in query parameters

## Endpoints

- `POST /orders` - Create a new order (SQL injection vulnerable)
- `GET /orders/:id` - Get order by ID (SQL injection vulnerable)
- `GET /orders` - List orders with filters (SQL injection, ORDER BY injection)
- `PUT /orders/:id/status` - Update order status (SQL injection vulnerable)
- `POST /orders/import` - Import orders from XML (XXE vulnerable)
- `GET /orders/export` - Export orders (information disclosure)
- `DELETE /orders/:id` - Delete order (SQL injection vulnerable)

## Running

```bash
go run main.go
```

Or with Docker:

```bash
docker build -t orders-service .
docker run -p 8083:8083 orders-service
```

## Testing

```bash
go test -v
```
