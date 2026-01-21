# Users Service

[![Users Service](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/users-service.yml/badge.svg?branch=main)](https://github.com/kadraman/InsecureGoMonorepo/actions/workflows/users-service.yml)

A deliberately insecure user management microservice for security testing demonstrations.

## Security Vulnerabilities (Intentional)

This service contains the following deliberate security vulnerabilities for testing purposes:

1. **SQL Injection**: Direct string concatenation in SQL queries
2. **Weak Password Hashing**: Uses MD5 instead of bcrypt
3. **Hardcoded Secrets**: JWT secret and API keys hardcoded in config
4. **Command Injection**: Through log file export functionality
5. **Insecure Token Generation**: Uses MD5 with hardcoded secret

## Endpoints

- `POST /users` - Create a new user
- `GET /users/:username` - Get user by username
- `POST /login` - Authenticate user
- `GET /search?q=<term>` - Search users (vulnerable to SQL injection)
- `GET /export?filename=<name>` - Export users (vulnerable to command injection)

## Running

```bash
go run main.go
```

Or with Docker:

```bash
docker build -t users-service .
docker run -p 8081:8081 users-service
```

## Testing

```bash
go test -v
```
