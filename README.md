# Go Microservices Template

Production-ready microservices template with gRPC, REST, Event-Driven messaging.

## Architecture

- **Clean Architecture** (Hexagonal)
- **gRPC + REST** via gRPC-gateway
- **Event-Driven** with RabbitMQ
- **PostgreSQL** with sqlc
- **Prometheus + Grafana** monitoring

## Tech Stack

- Go 1.21+
- gRPC / Protocol Buffers
- RabbitMQ (AMQP)
- PostgreSQL + sqlc
- Docker & Kubernetes
- Prometheus & Grafana

## Project Structure

```
├── cmd/ # Application entrypoints
│ └── api/ # Main API service
├── internal/ # Private application code
│ ├── domain/ # Business entities
│ ├── usecase/ # Business logic
│ └── infrastructure/ # External implementations
├── api/proto/ # Protocol buffer definitions
├── pkg/ # Public reusable packages
├── deployments/ # Docker and K8s configs
└── scripts/ # Build and utility scripts
```


## Getting Started

Coming soon...

