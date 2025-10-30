# Contributing to go-microservices-template

Thank you for considering contributing to this project! 🎉

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

---

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://go.dev/conduct). Please be respectful and professional in all interactions.

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- PostgreSQL 15+
- RabbitMQ 3.12+
- Make
- buf CLI v1.59.0+ (for protobuf)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/go-microservices-template.git
cd go-microservices-template
```

3. Add upstream remote:
```bash
git remote add upstream https://github.com/memclutter/go-microservices-template.git
```

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
make install-tools
```

### Start Development Environment

```bash
# Start infrastructure (PostgreSQL, RabbitMQ, etc.)
make docker-up

# Run database migrations
make migrate-up

# Generate code (protobuf, sqlc, wire)
make generate
```

---

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Adding tests

### 2. Make Your Changes

Follow the [Coding Standards](#coding-standards) below.

### 3. Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run linter
make lint
```

### 4. Commit Your Changes

Follow [Commit Message Guidelines](#commit-message-guidelines).

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

---

## Coding Standards

### Go Code Style

Follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

**Key points**:

1. **Formatting**: Use `gofmt` and `goimports`
   ```bash
   gofmt -s -w .
   goimports -w .
   ```

2. **Naming**:
    - Packages: lowercase, single word
    - Interfaces: `-er` suffix (e.g., `Reader`, `Publisher`)
    - Unexported: camelCase
    - Exported: PascalCase

3. **Error Handling**: Always check errors
   ```go
   // Good
   if err := doSomething(); err != nil {
       return fmt.Errorf("failed to do something: %w", err)
   }

   // Bad
   doSomething() // ignoring error
   ```

4. **Context**: Always pass `context.Context` as first parameter
   ```go
   func DoSomething(ctx context.Context, arg string) error
   ```

5. **Comments**:
    - Exported functions must have comments
    - Comments should be full sentences

### Project Structure

Follow Clean Architecture principles:

```
internal/
├── domain/          # Business entities (no dependencies)
├── usecase/         # Application logic
└── infrastructure/  # External adapters (DB, messaging, etc.)
```

**Rules**:
- Domain layer: No external dependencies
- Use case layer: Can depend on domain only
- Infrastructure layer: Implements interfaces from domain

### SQL Queries

Use `sqlc` for type-safe SQL:

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
```

Then run:
```bash
make sqlc-generate
```

### Protobuf

Follow [Buf Style Guide](https://buf.build/docs/best-practices/style-guide):

- Use `snake_case` for field names
- Use PascalCase for message names
- Include comments for all messages and fields

---

## Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting changes
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding tests
- `chore`: Changes to build process or auxiliary tools

### Examples

```
feat(user): add email verification

- Implement email verification use case
- Add verification token generation
- Send verification email via SMTP

Closes #123
```

```
fix(database): handle connection pool exhaustion

When connection pool was exhausted, the application would hang
indefinitely. Now it returns an error after timeout.

Fixes #456
```

---

## Pull Request Process

### Before Submitting

1. ✅ All tests pass (`make test`)
2. ✅ Linter passes (`make lint`)
3. ✅ Code is formatted (`gofmt -s`)
4. ✅ No unnecessary files (check `.gitignore`)
5. ✅ Documentation updated if needed
6. ✅ Commit messages follow guidelines

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
Describe how you tested your changes

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review performed
- [ ] Comments added where needed
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] No new warnings
```

### Review Process

1. Maintainer reviews code
2. Address feedback
3. Once approved, squash and merge

---

## Testing

### Unit Tests

```go
func TestUserCreate(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserInput{
                Email:    "test@example.com",
                Name:     "Test",
                Password: "password123",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

Use `testcontainers`:

```go
func TestUserRepository(t *testing.T) {
    // Start PostgreSQL container
    ctx := context.Background()
    postgres, err := testcontainers.GenericContainer(ctx, /* ... */)
    require.NoError(t, err)
    defer postgres.Terminate(ctx)

    // Run tests against real database
}
```

### Coverage Goals

- Unit tests: >80%
- Integration tests: Critical paths
- E2E tests: Main user flows

Run coverage:
```bash
make test-coverage
open coverage.html
```

---

## Documentation

### Code Comments

```go
// CreateUser creates a new user in the system.
// It validates the input, generates a unique ID, hashes the password,
// and stores the user in the database.
//
// Returns ErrUserAlreadyExists if email is taken.
func CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // ...
}
```

### README Updates

Update README.md if you:
- Add new features
- Change configuration
- Update dependencies
- Modify API endpoints

### API Documentation

Update `docs/api.md` for any API changes.

---

## Questions?

- Open an [issue](https://github.com/memclutter/go-microservices-template/issues)
- Ask in [discussions](https://github.com/memclutter/go-microservices-template/discussions)

---

**Thank you for contributing! 🚀**