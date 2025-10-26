# Go Microservices Template

[![CI](https://github.com/memclutter/go-microservices-template/actions/workflows/ci.yml/badge.svg)](https://github.com/memclutter/go-microservices-template/actions)
[![codecov](https://codecov.io/gh/memclutter/go-microservices-template/branch/main/graph/badge.svg)](https://codecov.io/gh/memclutter/go-microservices-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/memclutter/go-microservices-template)](https://goreportcard.com/report/github.com/memclutter/go-microservices-template)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Production-ready microservices template demonstrating modern Go backend practices with gRPC, REST, event-driven messaging, and full observability.

## ✨ Features

- **Clean Architecture** - Hexagonal architecture with clear separation of concerns
- **gRPC + REST** - Dual API support via gRPC-gateway
- **Event-Driven** - Asynchronous messaging with RabbitMQ
- **Type-Safe SQL** - Code generation with sqlc
- **Monitoring** - Prometheus metrics + Grafana dashboards
- **Docker** - Multi-stage builds for minimal images (~20MB)
- **Kubernetes** - Production-ready manifests with HPA
- **CI/CD** - Automated testing and deployment

## 📋 Table of Contents

- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Monitoring](#monitoring)
- [Contributing](#contributing)

## 🏗 Architecture

This project follows **Clean Architecture** (Hexagonal) principles:

```
┌─────────────────────────────────────────┐
│          API Layer (Adapters)           │
│     gRPC Server │ REST Gateway          │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│         Application Layer               │
│         (Use Cases / Business Logic)    │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│          Domain Layer                   │
│    (Entities, Value Objects, Rules)     │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      Infrastructure Layer (Adapters)    │
│  PostgreSQL │ RabbitMQ │ External APIs  │
└─────────────────────────────────────────┘
```

### Key Design Patterns

- **Repository Pattern** - Abstract data access
- **Dependency Injection** - Google Wire for compile-time DI
- **CQRS Lite** - Separate read/write use cases
- **Event Sourcing** - Domain events for async communication

## 🛠 Tech Stack

| Category | Technology |
|----------|-----------|
| **Language** | Go 1.21+ |
| **API** | gRPC, gRPC-gateway (REST) |
| **Database** | PostgreSQL 15+ |
| **ORM/Query** | sqlc (type-safe SQL) |
| **Messaging** | RabbitMQ (AMQP) |
| **Monitoring** | Prometheus, Grafana |
| **Logging** | slog (structured) |
| **Configuration** | Viper |
| **Testing** | testify, gomock |
| **Container** | Docker, Docker Compose |
| **Orchestration** | Kubernetes |
| **CI/CD** | GitHub Actions |

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- RabbitMQ 3.12+
- Make

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/memclutter/go-microservices-template.git
cd go-microservices-template
```

2. **Install development tools**
```bash
make install-tools
```

3. **Start infrastructure**
```bash
make docker-up
```

4. **Run database migrations**
```bash
make migrate-up
```

5. **Generate code** (protobuf + sqlc + wire)
```bash
make generate
```

6. **Run the application**
```bash
make run
```

The API will be available at:
- REST: `http://localhost:8080`
- gRPC: `localhost:50051`
- Metrics: `http://localhost:8080/metrics`

### Quick Test

```bash
# Create a user via REST
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","password":"password123"}'

# Get user
curl http://localhost:8080/v1/users/{user_id}
```

## 📁 Project Structure

```
.
├── api/
│   └── proto/              # Protocol buffer definitions
├── cmd/
│   └── api/               # Application entrypoint
├── internal/
│   ├── domain/            # Business entities & logic
│   ├── usecase/           # Application business flows
│   └── infrastructure/    # External implementations
│       ├── grpc/          # gRPC server
│       ├── http/          # REST gateway
│       ├── repository/    # Database access
│       └── messaging/     # Event pub/sub
├── pkg/                   # Public reusable packages
│   ├── config/           # Configuration
│   ├── logger/           # Logging
│   └── metrics/          # Prometheus metrics
├── db/
│   ├── migrations/       # Database migrations
│   └── queries/          # SQL queries for sqlc
├── deployments/
│   ├── docker/           # Dockerfiles & compose
│   └── k8s/              # Kubernetes manifests
├── scripts/              # Build & utility scripts
└── docs/                 # Additional documentation
```

## 📚 API Documentation

### REST Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/users` | Create user |
| GET | `/v1/users/{id}` | Get user by ID |
| PUT | `/v1/users/{id}` | Update user |
| DELETE | `/v1/users/{id}` | Delete user |
| GET | `/v1/users` | List users |

### gRPC Methods

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse)
  rpc GetUser(GetUserRequest) returns (GetUserResponse)
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse)
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse)
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse)
}
```

See [API Documentation](docs/api.md) for detailed request/response schemas.

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
```

### Test Coverage

- Unit tests for all layers
- Integration tests with testcontainers
- Mocks generated with gomock
- Table-driven tests

## 🚢 Deployment

### Docker

```bash
# Build image
make docker-build

# Run with docker-compose
make docker-up

# View logs
make docker-logs
```

### Kubernetes

```bash
# Deploy to cluster
make k8s-apply

# Check status
make k8s-status

# Scale deployment
kubectl scale deployment api -n microservices --replicas=5
```

### Production Checklist

- [ ] Use external secret management (Vault, AWS Secrets Manager)
- [ ] Configure TLS for gRPC and HTTPS
- [ ] Set up log aggregation (ELK, Loki)
- [ ] Configure distributed tracing (Jaeger, Zipkin)
- [ ] Enable rate limiting
- [ ] Set up backups for PostgreSQL
- [ ] Configure alerts in Prometheus

## 📊 Monitoring

### Prometheus Metrics

Access metrics at `http://localhost:8080/metrics`

Key metrics:
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `database_queries_total` - Database operations
- `events_published_total` - Published events

### Grafana Dashboards

Access Grafana at `http://localhost:3000` (admin/admin)

Pre-configured dashboards:
- Service Overview
- HTTP Performance
- Database Performance
- RabbitMQ Metrics

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file.

## 🙏 Acknowledgments

- Clean Architecture by Robert Martin
- Domain-Driven Design by Eric Evans
- Go community and maintainers

---

**Happy coding! 🚀**
