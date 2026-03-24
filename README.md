# Public API

![CI](https://github.com/GunarsK-rpg/public-api/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/GunarsK-rpg/public-api/graph/badge.svg)](https://codecov.io/gh/GunarsK-rpg/public-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/GunarsK-rpg/public-api)](https://goreportcard.com/report/github.com/GunarsK-rpg/public-api)
[![CodeRabbit](https://img.shields.io/coderabbit/prs/github/GunarsK-rpg/public-api?label=CodeRabbit&color=2ea44f)](https://coderabbit.ai)

RESTful API for Cosmere RPG with PostgreSQL stored functions (JSONB passthrough).

## Features

- Classifier endpoints (30+ read-only reference data)
- Hero CRUD with 13 sub-resources (attributes, talents, equipment, etc.)
- Campaign CRUD
- JWT authentication via portfolio auth-service
- Permission-based access control (read/edit/delete levels)
- Transaction-scoped audit context
- User sync from JWT to RPG database

## Tech Stack

- **Language**: Go 1.26
- **Framework**: Gin
- **Database**: PostgreSQL 17+ (pgx, stored functions)
- **Common**: [portfolio-common][pc] (auth, logging, metrics)
- **Auth**: JWT validation via [auth-service][as]

[pc]: https://github.com/GunarsK-portfolio/portfolio-common
[as]: https://github.com/GunarsK-portfolio/auth-service

## Prerequisites

- Go 1.26+
- Node.js 22+ and npm 11+
- PostgreSQL with Cosmere RPG database schema
- auth-service running

## Project Structure

```text
public-api/
├── cmd/
│   └── api/              # Application entrypoint
├── internal/
│   ├── cache/            # Redis + in-memory caching
│   ├── config/           # Configuration
│   ├── constants/        # Resource constants
│   ├── database/         # pgx pool and health checker
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # User sync, body limit
│   ├── models/           # Request models
│   ├── repository/       # Data access layer (DB function calls)
│   └── routes/           # Route definitions
```

## Quick Start

### Using Docker Compose

```bash
docker-compose up -d
```

### Local Development

1. Copy environment file:

```bash
cp .env.example .env
```

1. Update `.env` with your configuration:

```env
PORT=8182
DB_HOST=localhost
DB_PORT=5432
DB_USER=cosmere_app
DB_PASSWORD=cosmere_app_dev_pass
DB_NAME=cosmere_rpg
JWT_SECRET=your-secret-matching-auth-service
```

1. Run the service:

```bash
go run cmd/api/main.go
```

## Available Commands

Using Task:

```bash
# Development
task dev:install-tools   # Install dev tools (golangci-lint, govulncheck, etc.)

# Build and run
task build               # Build binary
task test                # Run tests
task test:coverage       # Run tests with coverage report
task clean               # Clean build artifacts

# Code quality
task format              # Format code with gofmt
task tidy                # Tidy and verify go.mod
task lint                # Run golangci-lint
task vet                 # Run go vet

# Security
task security:scan       # Run gosec security scanner
task security:vuln       # Check for vulnerabilities with govulncheck

# Docker
task docker:build        # Build Docker image
task docker:run          # Run service in Docker container
task docker:stop         # Stop running Docker container
task docker:logs         # View Docker container logs

# CI/CD
task ci:all              # Run all CI checks
```

Using Go directly:

```bash
go run cmd/api/main.go                      # Run
go build -o bin/public-api cmd/api/main.go  # Build
go test ./...                               # Test
```

## Environment Variables

| Variable          | Description                          | Default       |
| ----------------- | ------------------------------------ | ------------- |
| `PORT`            | Server port                          | `8182`        |
| `DB_HOST`         | PostgreSQL host                      | `localhost`   |
| `DB_PORT`         | PostgreSQL port                      | `5432`        |
| `DB_USER`         | Database user                        | `cosmere_app` |
| `DB_PASSWORD`     | Database password                    | -             |
| `DB_NAME`         | Database name                        | `cosmere_rpg` |
| `DB_SSLMODE`      | PostgreSQL SSL mode                  | `disable`     |
| `JWT_SECRET`      | JWT secret (must match auth-service) | -             |
| `ALLOWED_ORIGINS` | CORS allowed origins                 | -             |
| `LOG_LEVEL`       | Log level (debug/info/warn/error)    | `info`        |
| `LOG_FORMAT`      | Log format (text/json)               | `text`        |
| `MAX_BODY_SIZE`   | Max request body size in bytes       | `65536`       |

## Testing

```bash
task test              # Run all tests
task test:coverage     # Coverage report
```

See [TESTING.md](TESTING.md) for detailed testing documentation.

## Authentication

This API validates JWT tokens issued by
[auth-service][as] using the [portfolio-common][pc]
auth middleware.
Tokens must include:

- `user_id`: User's numeric ID
- `username`: User's display name
- `scopes`: Permission map, e.g., `{"heroes": "edit", "campaigns": "read"}`

Permission levels are hierarchical: `none < read < edit < delete`

All endpoints except `/health` and `/metrics` require authentication.

## License

[MIT](LICENSE)
