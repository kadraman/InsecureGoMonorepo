# InsecureGoMonorepo

⚠️ **WARNING: This repository contains intentionally vulnerable code for security testing demonstrations. DO NOT use in production!**

An insecure Go microservices monorepo designed for application security testing scenarios with GitHub Advanced Security, Fortify/OpenText, and other security scanning tools.

## Overview

This monorepo contains multiple Go microservices with deliberate security vulnerabilities for educational and testing purposes. Each service uses the Gin web framework for routing and shares common packages for logging, configuration, and database operations.

## Architecture

### Services

- **users-service** (Port 8081) - User management and authentication
- **products-service** (Port 8082) - Product catalog management
- **orders-service** (Port 8083) - Order processing
- **api-gateway** (Port 8080) - API gateway for routing requests

### Shared Packages

- **pkg/logging** - Logging utilities
- **pkg/config** - Configuration management
- **pkg/database** - Database operations

## Security Vulnerabilities (Intentional)

This repository demonstrates various security vulnerabilities including:

### SQL Injection
- Direct string concatenation in SQL queries across all services
- Vulnerable search, filter, and sorting functionality

### Command Injection
- Shell command execution with user input
- Logging functionality that executes system commands

### Hardcoded Secrets
- API keys, JWT secrets, and database passwords in source code
- Credentials exposed through debug endpoints

### Authentication & Authorization
- Weak password hashing (MD5)
- Insecure token generation
- Missing authentication checks

### Path Traversal
- Unrestricted file path access
- No validation on file operations

### XXE (XML External Entity)
- Unsafe XML parsing in import functionality

### Information Disclosure
- Debug endpoints exposing sensitive configuration
- Detailed error messages revealing system information

### Other Vulnerabilities
- Open redirect
- Unrestricted file upload
- ORDER BY SQL injection
- Header injection

## Quick Start

### Prerequisites

- Go 1.24 or later
- Docker (optional)
- Make (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/kadraman/InsecureGoMonorepo.git
cd InsecureGoMonorepo

# Install dependencies
go mod download
```

### Running Services

#### Using Make

```bash
# Build all services
make build

# Run individual services
make run-users      # Users service on :8081
make run-products   # Products service on :8082
make run-orders     # Orders service on :8083
make run-gateway    # API Gateway on :8080
```

#### Using Go directly

```bash
# Users service
cd services/users-service && go run main.go

# Products service
cd services/products-service && go run main.go

# Orders service
cd services/orders-service && go run main.go

# API Gateway
cd services/api-gateway && go run main.go
```

#### Using Docker

```bash
# Build Docker images
make docker-build

# Or build individually
docker build -t users-service -f services/users-service/Dockerfile .
docker build -t products-service -f services/products-service/Dockerfile .
docker build -t orders-service -f services/orders-service/Dockerfile .
docker build -t api-gateway -f services/api-gateway/Dockerfile .
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Test specific package
go test ./pkg/logging/...
go test ./services/users-service/...
```

## API Examples

### Users Service

```bash
# Create a user (SQL Injection vulnerable)
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"password123"}'

# Search users (SQL Injection vulnerable)
curl "http://localhost:8081/search?q=john' OR '1'='1"

# Login
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"password123"}'
```

### Products Service

```bash
# Create a product
curl -X POST http://localhost:8082/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":999.99,"category":"Electronics"}'

# Execute command (Command Injection vulnerable)
curl "http://localhost:8082/execute?cmd=ls -la"
```

### Orders Service

```bash
# Create an order
curl -X POST http://localhost:8083/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product_id":1,"quantity":2,"total_price":1999.98}'

# List orders with ORDER BY injection
curl "http://localhost:8083/orders?sort_by=id; DROP TABLE orders--"
```

### API Gateway

```bash
# Access debug endpoint (Information Disclosure)
curl http://localhost:8080/api/debug

# Open redirect vulnerability
curl "http://localhost:8080/api/redirect?url=http://malicious-site.com"
```

## Development

### Code Formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

### Cleaning Build Artifacts

```bash
make clean
```

## Project Structure

```
InsecureGoMonorepo/
├── services/
│   ├── users-service/
│   │   ├── main.go
│   │   ├── main_test.go
│   │   ├── Dockerfile
│   │   └── README.md
│   ├── products-service/
│   │   ├── main.go
│   │   ├── main_test.go
│   │   ├── Dockerfile
│   │   └── README.md
│   ├── orders-service/
│   │   ├── main.go
│   │   ├── main_test.go
│   │   ├── Dockerfile
│   │   └── README.md
│   └── api-gateway/
│       ├── main.go
│       ├── main_test.go
│       ├── Dockerfile
│       └── README.md
├── pkg/
│   ├── logging/
│   │   ├── logger.go
│   │   └── logger_test.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   └── database/
│       ├── database.go
│       └── database_test.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Security Testing

This repository is designed to be scanned with:

- **GitHub Advanced Security** - Code scanning and secret scanning
- **Fortify/OpenText** - Static application security testing (SAST)
- **Snyk** - Dependency scanning
- **SonarQube** - Code quality and security analysis
- **OWASP ZAP** - Dynamic application security testing (DAST)

## Contributing

This is a demonstration repository. If you'd like to add more vulnerability examples or improve existing ones, please feel free to submit a pull request.

## License

See LICENSE file for details.

## Disclaimer

⚠️ **IMPORTANT**: This code is intentionally insecure and should NEVER be used in production environments. It is designed solely for security testing, training, and demonstration purposes. The authors are not responsible for any misuse of this code.
