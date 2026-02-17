# Go Books API

A simple REST API for managing books, built with Go (Gin), Gorm, and MySQL.
Refactored to follow the **Controller-Service-Repository** pattern.

## Prerequisites

- Go 1.24+
- Docker & Docker Compose

## Running the Application

### Option 1: Docker (Recommended)

This spins up the API and a MySQL database in containers.

```bash
make docker-up
# Or manually:
docker-compose up -d --build
```

The API will be available at `http://localhost:8080`.

To stop:
```bash
make docker-down
```

### Option 2: Local (Requires local DB)

If you have a local MySQL instance running:

1. Update `platform/database.go` or set environment variables for DB connection.
2. Run:
   ```bash
   make run
   # Or manually:
   go run cmd/api/main.go
   ```

## Development

### Running Tests

```bash
make test
```

### Checking Coverage

```bash
make coverage
```

## API Endpoints

### Create a Book

**POST** `/books`

```json
{
    "title": "The Refactoring",
    "author": "Martin Fowler"
}
```

Example Curl:
```bash
curl -v -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Refactoring", "author": "Martin Fowler"}'
```

### Get All Books

**GET** `/books`

Example Curl:
```bash
curl -v http://localhost:8080/books
```

### Get One Book

**GET** `/books/:id`

Example Curl:
```bash
curl -v http://localhost:8080/books/1
```

## creating Implementation Plan

The architecture follows a standard layered approach:

- **Controller** (`internal/controller`): Handles HTTP requests.
- **Service** (`internal/service`): Contains business logic.
- **Repository** (`internal/repository`): Handles data persistence.
- **Domain** (`internal/domain`): Contains entities and repository interfaces.
- **Platform** (`platform/`): Generic infrastructure (DB, Telemetry).
