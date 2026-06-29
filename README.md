# Task Manager

A simple REST API for task management written in Go.

## Features

* Create a task
* Get a task by ID
* Get all tasks
* Update a task
* Delete a task

## Technologies

* Go
* net/http
* JSON
* Docker
* Unit Tests

## Project Structure

```text
cmd/
internal/
    api/
    models/
    service/
    storage/
```

## Running the Application

```bash
go run ./cmd/main.go
```

The server will start on:

```text
http://localhost:8080
```

## Running Tests

```bash
go test ./...
```

## Docker

Build image:

```bash
docker build -t taskmanager .
```

Run container:

```bash
docker run -p 8080:8080 taskmanager
```

## API Examples

Create a task:

```bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d "{\"title\":\"Learn Go\",\"description\":\"Practice every day\"}"
```

Get all tasks:

```bash
curl http://localhost:8080/tasks
```

Get a task by ID:

```bash
curl http://localhost:8080/tasks?id=1
```

Update a task:

```bash
curl -X PUT http://localhost:8080/tasks?id=1 \
-H "Content-Type: application/json" \
-d "{\"title\":\"Updated task\",\"description\":\"Updated description\"}"
```

Delete a task:

```bash
curl -X DELETE http://localhost:8080/tasks?id=1
```

