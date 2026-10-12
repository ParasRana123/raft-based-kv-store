# Raft-Based Replicated Key-Value Store

A high-performance, concurrent, and persistent key-value store built in Go, featuring a Raft consensus engine foundation for distributed replication and fault tolerance.

---

## Architecture Overview

```
raft_kv_store/
├── api/             # HTTP REST API server and handlers
│   ├── server.go
│   └── server_test.go
├── cmd/             # Application entrypoints
│   └── node/
│       └── main.go  # Node startup binary
├── kv/              # Core storage engine
│   ├── store.go     # Thread-safe KVStore with JSON persistence
│   └── store_test.go
├── raft/            # Raft distributed consensus engine
│   ├── types.go     # Node roles, log entries, and state machine messages
│   ├── rpc.go       # RequestVote and AppendEntries RPC definitions
│   ├── raft.go      # RaftNode state machine & consensus handlers
│   └── raft_test.go # Raft consensus unit test suite
├── go.mod
└── README.md
```

---

## Features

### 1. Raft Consensus Engine (`raft/`)
- **Roles & State Transitions**: Implements `Follower`, `Candidate`, and `Leader` states.
- **Leader Election Safety (`RequestVote`)**:
  - Rejects outdated candidate terms.
  - Enforces the single-vote-per-term invariant.
  - Compares candidate log freshness before granting votes.
- **Log Replication & Heartbeats (`AppendEntries`)**:
  - Handles leader heartbeats and term step-downs.
  - Verifies log prefix consistency at `PrevLogIndex` and `PrevLogTerm`.
  - Appends new entries and resolves conflicts.
  - Updates commit index and streams committed commands via `applyCh`.

### 2. Thread-Safe Key-Value Engine (`kv/`)
- **Concurrent Access**: Uses `sync.RWMutex` for non-blocking concurrent reads and isolated exclusive writes.
- **Disk Persistence & Recovery**: Flushes state to `data.json` on mutations (`Set`/`Delete`) and restores data on node restart.
- **Clean Interface Contract**: Decoupled `KVStore` interface supporting in-memory and disk-backed implementations.

### 3. RESTful HTTP API (`api/`)
- `GET /health` — Health check endpoint.
- `PUT /kv/{key}` — Save or update key-value pair (accepts JSON or raw text).
- `GET /kv/{key}` — Retrieve value for a given key (`200 OK` or `404 Not Found`).
- `DELETE /kv/{key}` — Remove key from the store.

---

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

---

## Running Tests

Run the full automated test suite across all packages (`kv`, `api`, `raft`):

```bash
go test -v ./...
```

To run tests for a specific package:

```bash
# Raft consensus tests
go test -v ./raft

# Storage layer tests
go test -v ./kv

# HTTP API tests
go test -v ./api
```

---

## HTTP API Usage Examples

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
