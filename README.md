# Task Manager

REST API for task management written in Go.

The project supports basic CRUD operations for tasks and uses PostgreSQL for persistent data storage.

## Features

* Create a task
* Get a task by ID
* Get all tasks
* Update a task
* Delete a task
* PostgreSQL data storage
* Database migrations

## Technologies

* Go
* PostgreSQL
* pgx
* net/http
* JSON
* Docker
* golang-migrate
* Unit Tests

## Architecture

The application follows a layered architecture:

```text
HTTP
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
PostgreSQL
```

* **Handler** — handles HTTP requests and responses
* **Service** — contains business logic
* **Repository** — provides data access
* **PostgreSQL** — stores application data

The service depends on a repository interface, which allows the storage implementation to be changed without modifying the business logic.

## Project Structure

```text
cmd/
internal/
    handlers/
    services/
    repository/
    store/
    models/
migrations/
```

## Database Migrations

Database schema changes are managed with `golang-migrate`.

Example migration:

```text
migrations/
├── 000001_create_tasks.up.sql
└── 000001_create_tasks.down.sql
```

To apply migrations:

```bash
migrate -path migrations -database "postgres://USERNAME:PASSWORD@localhost:5432/task_manager?sslmode=disable" up
```

## Running the Application

Configure the database connection and application settings, then run:

```bash
go run ./cmd/main.go
```

## Running Tests

```bash
go test ./...
```

## API Examples

### Create a task

```bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d "{\"title\":\"Learn Go\",\"description\":\"Practice every day\"}"
```

### Get all tasks

```bash
curl http://localhost:8080/tasks
```

### Get a task by ID

```bash
curl http://localhost:8080/tasks?id=1
```

### Update a task

```bash
curl -X PUT http://localhost:8080/tasks?id=1 \
-H "Content-Type: application/json" \
-d "{\"title\":\"Updated task\",\"description\":\"Updated description\"}"
```

### Delete a task

```bash
curl -X DELETE http://localhost:8080/tasks?id=1
```
