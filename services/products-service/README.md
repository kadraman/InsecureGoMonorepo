# Products Service

[![Products Service](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/products-service.yml/badge.svg?branch=main)](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/products-service.yml)

A deliberately insecure product management microservice for security testing demonstrations.

## Security Vulnerabilities (Intentional)

This service contains the following deliberate security vulnerabilities for testing purposes:

1. **SQL Injection**: Direct string concatenation in SQL queries for all CRUD operations
2. **Path Traversal**: Unrestricted file path access in image retrieval
3. **Command Injection**: Direct command execution through execute endpoint
4. **Unrestricted File Upload**: No validation on file type or size
5. **Information Disclosure**: Error messages reveal system information

## Endpoints

- `POST /products` - Create a new product (SQL injection vulnerable)
- `GET /products/:id` - Get product by ID (SQL injection vulnerable)
- `PUT /products/:id` - Update product (SQL injection vulnerable)
- `DELETE /products/:id` - Delete product (SQL injection vulnerable)
- `GET /products` - List products with filters (SQL injection vulnerable)
- `GET /images/:filename` - Get image (path traversal vulnerable)
- `POST /images` - Upload image (unrestricted upload)
- `GET /execute?cmd=<command>` - Execute command (command injection)

## Running

```bash
go run main.go
```

Or with Docker:

```bash
docker build -t products-service .
docker run -p 8082:8082 products-service
```

## Testing

```bash
go test -v
```
