# API Documentation

## Overview

The microservices template provides both **gRPC** and **REST** APIs through gRPC-gateway. All endpoints support JSON format.

## Base URLs

- **REST API**: `http://localhost:8080/v1`
- **gRPC**: `localhost:50051`
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

---

## Authentication

> **Note**: Authentication is not implemented in this template. Add JWT/OAuth2 middleware for production use.

---

## User Service

### Create User

Creates a new user in the system.

**gRPC Method**: `UserService.CreateUser`

**REST Endpoint**: `POST /v1/users`

**Request Body**:
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securePassword123"
}
```

**Response** (201 Created):
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-10-30T19:00:00Z",
    "updated_at": "2025-10-30T19:00:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input (missing email, name, or weak password)
- `409 Conflict`: User with this email already exists
- `500 Internal Server Error`: Server error

**cURL Example**:
```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "securePassword123"
  }'
```

---

### Get User

Retrieves a user by their ID.

**gRPC Method**: `UserService.GetUser`

**REST Endpoint**: `GET /v1/users/{user_id}`

**Path Parameters**:
- `user_id` (string, required): UUID of the user

**Response** (200 OK):
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-10-30T19:00:00Z",
    "updated_at": "2025-10-30T19:00:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Missing user_id
- `404 Not Found`: User does not exist
- `500 Internal Server Error`: Server error

**cURL Example**:
```bash
curl http://localhost:8080/v1/users/550e8400-e29b-41d4-a716-446655440000
```

---

### Update User

Updates an existing user's profile.

**gRPC Method**: `UserService.UpdateUser`

**REST Endpoint**: `PUT /v1/users/{user_id}`

**Path Parameters**:
- `user_id` (string, required): UUID of the user

**Request Body**:
```json
{
  "name": "Jane Doe"
}
```

**Response** (200 OK):
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "Jane Doe",
    "created_at": "2025-10-30T19:00:00Z",
    "updated_at": "2025-10-30T19:05:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input
- `404 Not Found`: User does not exist
- `500 Internal Server Error`: Server error

**Status**: 🚧 Not implemented yet

---

### Delete User

Deletes a user from the system.

**gRPC Method**: `UserService.DeleteUser`

**REST Endpoint**: `DELETE /v1/users/{user_id}`

**Path Parameters**:
- `user_id` (string, required): UUID of the user

**Response** (204 No Content):
```json
{}
```

**Error Responses**:
- `400 Bad Request`: Missing user_id
- `404 Not Found`: User does not exist
- `500 Internal Server Error`: Server error

**Status**: 🚧 Not implemented yet

---

### List Users

Retrieves a paginated list of users.

**gRPC Method**: `UserService.ListUsers`

**REST Endpoint**: `GET /v1/users`

**Query Parameters**:
- `limit` (int32, optional): Number of users per page (default: 10, max: 100)
- `offset` (int32, optional): Offset for pagination (default: 0)

**Response** (200 OK):
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user1@example.com",
      "name": "John Doe",
      "created_at": "2025-10-30T19:00:00Z",
      "updated_at": "2025-10-30T19:00:00Z"
    },
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "email": "user2@example.com",
      "name": "Jane Smith",
      "created_at": "2025-10-30T18:30:00Z",
      "updated_at": "2025-10-30T18:30:00Z"
    }
  ],
  "total": 42
}
```

**cURL Example**:
```bash
curl "http://localhost:8080/v1/users?limit=20&offset=0"
```

**Status**: 🚧 Not implemented yet

---

## gRPC Testing

### Using grpcurl

Install grpcurl:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

**List services**:
```bash
grpcurl -plaintext localhost:50051 list
```

**Create user**:
```bash
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "name": "Test User",
  "password": "password123"
}' localhost:50051 user.v1.UserService/CreateUser
```

**Get user**:
```bash
grpcurl -plaintext -d '{
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:50051 user.v1.UserService/GetUser
```

---

## Health & Monitoring Endpoints

### Health Check

**Endpoint**: `GET /health`

Returns `200 OK` if the service is running.

```bash
curl http://localhost:8080/health
```

Response:
```
OK
```

---

### Readiness Check

**Endpoint**: `GET /ready`

Returns `200 OK` if the service is ready to handle requests (database connected, etc.).

```bash
curl http://localhost:8080/ready
```

Response:
```
Ready
```

Or `503 Service Unavailable` if not ready:
```
Database not ready
```

---

### Prometheus Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics in text format.

```bash
curl http://localhost:8080/metrics
```

**Key Metrics**:
- `microservices_http_requests_total` - Total HTTP requests
- `microservices_http_request_duration_seconds` - HTTP request latency histogram
- `microservices_grpc_requests_total` - Total gRPC requests
- `microservices_grpc_request_duration_seconds` - gRPC request latency histogram
- `microservices_database_queries_total` - Total database queries
- `microservices_events_published_total` - Total events published

---

## Error Handling

All endpoints return errors in a consistent format:

**REST Error Response**:
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "email is required",
  "details": {
    "field": "email"
  }
}
```

**gRPC Error Codes**:
- `INVALID_ARGUMENT` (3): Bad request
- `NOT_FOUND` (5): Resource not found
- `ALREADY_EXISTS` (6): Resource already exists
- `INTERNAL` (13): Internal server error
- `UNIMPLEMENTED` (12): Feature not implemented

---

## Rate Limiting

> **Note**: Rate limiting is not implemented in this template. Add middleware for production use.

Recommended libraries:
- [`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate)
- [`github.com/ulule/limiter/v3`](https://github.com/ulule/limiter)

---

## Versioning

API versioning is handled through URL path:
- Current version: `/v1/`
- Future version: `/v2/` (when needed)

gRPC uses package versioning:
- Current: `user.v1`
- Future: `user.v2`

---

## Further Reading

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [gRPC-gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)