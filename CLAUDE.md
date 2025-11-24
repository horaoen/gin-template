# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend application implementing Clean Architecture patterns with the following key technologies:
- **Web Framework**: Gin (HTTP router and middleware)
- **Database**: MongoDB with custom abstraction layer
- **Authentication**: JWT tokens with access/refresh token pattern
- **Logging**: Zerolog for structured logging
- **Configuration**: Viper for environment-based configuration

## Architecture

The project follows Clean Architecture principles with clear separation of concerns:

### Core Layers
- **`domain/`**: Business entities and interfaces (User, Task, JWT claims)
- **`usecase/`**: Business logic implementation (auth, profile, task management)
- **`repository/`**: Data access interfaces and MongoDB implementations
- **`api/`**: HTTP controllers, middleware, and routing
- **`bootstrap/`**: Application initialization and dependency injection

### Key Components
- **`mongo/`**: Custom MongoDB client abstraction with interface-based design
- **`internal/`**: Internal utilities (JWT token handling, test utilities)
- **`cmd/`**: Application entry point

### Request Flow
1. HTTP request hits Gin router → Controller → Usecase → Repository → MongoDB
2. JWT middleware protects private routes with token validation
3. Environment-based configuration manages database connections and secrets

## Development Commands

### Build and Run
```bash
# Build the application
go build -o main cmd/main.go

# Run directly
go run cmd/main.go

# Run with Docker Compose (includes MongoDB)
docker-compose up --build
```

### Testing
```bash
# Run all tests
go test -v ./...

# Run specific test file
go test -v ./api/controller/profile_controller_test.go

# Run tests for specific package
go test -v ./usecase/
```

### Code Quality
```bash
# Run comprehensive code quality checks
golangci-lint run

# Run specific checks (errcheck, staticcheck, unused)
golangci-lint run --enable-only=errcheck,staticcheck,unused

# Format code
go fmt ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy up go.mod
go mod tidy

# Update dependencies
go get -u ./...
```

## Key Patterns and Conventions

### Interface-Driven Design
- All major components have interfaces for testability
- MongoDB operations are abstracted through custom interfaces in `mongo/`
- Repository pattern with mock implementations in `domain/mocks/`

### JWT Implementation
- Uses `jwt.RegisteredClaims` (not deprecated `StandardClaims`)
- Separate access and refresh token utilities in `internal/tokenutil/`
- Custom claims defined in `domain/jwt_custom.go`

### Error Handling
- Follow Go conventions: error messages start with lowercase
- Structured error responses defined in `domain/error_response.go`
- Consistent error format across API endpoints

### Configuration
- Environment-based config via `.env` file (see `.env.example`)
- Centralized in `bootstrap/env.go`
- Supports different environments through environment variables

## Testing Strategy

- Unit tests exist for controllers, usecases, and repositories
- Test files follow `*_test.go` naming convention
- Mock implementations available in `domain/mocks/` directory
- Tests use testify for assertions and mocking

## Database

- MongoDB connection and client abstracted in `mongo/mongo.go`
- Custom interface-based design allows for easy testing and swapping
- Database initialization handled in `bootstrap/database.go`
- Connection cleanup handled in application lifecycle

## Important Notes

- Application uses Go 1.24+ features
- All external dependencies managed through go.mod
- JWT secrets and database credentials should be configured via environment variables
- The codebase follows Go naming conventions and clean architecture principles
- Error strings should be lowercase (per Go conventions)
