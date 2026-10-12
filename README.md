# Raft-Based Replicated Key-Value Store

A high-performance, concurrent, and persistent key-value store built in Go, designed as the foundation for a Raft-replicated distributed system.

## Features

- **In-Memory & Persistent Storage**: Fast concurrent reads with `sync.RWMutex` combined with persistent JSON file storage (`data.json`) across server restarts.
- **Thread-Safe**: Reader-writer mutex locking ensuring concurrent access without race conditions.
- **RESTful HTTP API**:
  - `GET /health` — Health check endpoint.
  - `PUT /kv/{key}` — Save or update key-value pair (supports JSON or raw text).
  - `GET /kv/{key}` — Retrieve value for a given key.
  - `DELETE /kv/{key}` — Delete key from the store.
- **Automated Test Suite**: Full unit and integration tests covering CRUD operations, concurrency, persistence, and restart recovery.

## Getting Started

### Prerequisites

- Go 1.22+ installed.

### Running the Server

Start the node server:

```bash
go run ./cmd/node
```

By default, the server starts on port `:8080` and persists state to `data.json`.
You can customize the persistence file with the `DATA_FILE` environment variable:

```bash
DATA_FILE="custom_data.json" go run ./cmd/node
```

### Running Tests

Run the full automated test suite:

```bash
go test -v ./...
```

## API Usage Examples

### 1. Health Check
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 2. Set / Update a Key
```bash
curl -X PUT http://localhost:8080/kv/user \
  -H "Content-Type: application/json" \
  -d '{"value":"Paras"}'
# {"key":"user","value":"Paras"}
```

### 3. Retrieve a Key
```bash
curl http://localhost:8080/kv/user
# {"key":"user","value":"Paras"}
```

### 4. Delete a Key
```bash
curl -X DELETE http://localhost:8080/kv/user
# {"message":"key deleted"}
```
