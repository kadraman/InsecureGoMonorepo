# API Gateway

A deliberately insecure API gateway for routing requests to microservices, designed for security testing demonstrations.

## Security Vulnerabilities (Intentional)

This service contains the following deliberate security vulnerabilities for testing purposes:

1. **Weak Authentication**: Accepts requests with hardcoded API key or no authentication
2. **Information Disclosure**: Debug endpoint exposes sensitive configuration including passwords and secrets
3. **Open Redirect**: Unvalidated redirect endpoint
4. **Header Injection**: Forwards all headers without validation
5. **Environment Variable Exposure**: Debug endpoint exposes all environment variables

## Endpoints

### Gateway Routes
- `POST /api/users` - Proxy to users-service
- `GET /api/users/:username` - Proxy to users-service
- `POST /api/login` - Proxy to users-service
- `GET /api/users/search` - Proxy to users-service

- `POST /api/products` - Proxy to products-service
- `GET /api/products/:id` - Proxy to products-service
- `PUT /api/products/:id` - Proxy to products-service
- `DELETE /api/products/:id` - Proxy to products-service
- `GET /api/products` - Proxy to products-service

- `POST /api/orders` - Proxy to orders-service
- `GET /api/orders/:id` - Proxy to orders-service
- `GET /api/orders` - Proxy to orders-service
- `PUT /api/orders/:id/status` - Proxy to orders-service

### Vulnerable Endpoints
- `GET /api/debug` - Exposes sensitive configuration (information disclosure)
- `GET /api/redirect?url=<url>` - Open redirect vulnerability

## Configuration

Environment variables:
- `PORT` - Gateway port (default: 8080)
- `USERS_SERVICE_URL` - Users service URL (default: http://localhost:8081)
- `PRODUCTS_SERVICE_URL` - Products service URL (default: http://localhost:8082)
- `ORDERS_SERVICE_URL` - Orders service URL (default: http://localhost:8083)

## Running

```bash
go run main.go
```

Or with Docker:

```bash
docker build -t api-gateway .
docker run -p 8080:8080 api-gateway
```

## Testing

```bash
go test -v
```
