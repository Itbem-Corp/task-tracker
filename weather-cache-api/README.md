# Weather Cache API

A lightweight, high-performance in-memory cache service built in Go with REST API endpoints. Features thread-safe operations, TTL (Time-To-Live) support, and automatic cleanup of expired entries.

## Features

- **In-Memory Cache Store**: Fast key-value storage with optional TTL support
- **Thread-Safe Operations**: Concurrent reads and writes using sync.RWMutex
- **Automatic Cleanup**: Background goroutine cleans expired entries every 10 seconds
- **REST API**: Full CRUD operations via HTTP endpoints
- **No External Dependencies**: Uses Go standard library only
- **Lightweight Docker Image**: Multi-stage build produces ~15MB image
- **Comprehensive Tests**: Full test coverage for cache and HTTP handlers

## Project Structure

```
weather-cache-api/
├── main.go                          # Entry point with HTTP server setup
├── go.mod                           # Go module definition
├── Dockerfile                       # Multi-stage Docker build
├── README.md                        # This file
├── internal/
│   ├── cache/
│   │   ├── cache.go                # In-memory cache implementation
│   │   └── cache_test.go           # Cache unit tests
│   └── handlers/
│       ├── handlers.go             # HTTP endpoint handlers
│       └── handlers_test.go        # Handler unit tests
```

## API Endpoints

### 1. Health Check

**Endpoint:** `GET /health`

**Description:** Returns the service health status.

**Response:**
```json
{
  "status": "healthy"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### 2. Store Key-Value Pair

**Endpoint:** `POST /cache`

**Description:** Stores a key-value pair with optional TTL.

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "key": "string (required)",
  "value": "string (required)",
  "ttl_seconds": "integer (optional)"
}
```

**Parameters:**
- `key`: The cache key (required, non-empty)
- `value`: The cache value (required, non-empty)
- `ttl_seconds`: Time-to-live in seconds (optional, default: no expiration)

**Response (Success):**
```json
{
  "status": "ok",
  "key": "mykey"
}
```

**Response (Error):**
```json
{
  "error": "key and value are required"
}
```

**Status Codes:**
- `201 Created`: Successfully stored
- `400 Bad Request`: Invalid request (missing fields, invalid JSON)

**Examples:**

Store without TTL:
```bash
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "value": "John Doe"}'
```

Store with 1 hour TTL:
```bash
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key": "session:abc", "value": "token123", "ttl_seconds": 3600}'
```

---

### 3. Retrieve Value by Key

**Endpoint:** `GET /cache/{key}`

**Description:** Retrieves the value for a given key if it exists and hasn't expired.

**Response (Success):**
```json
{
  "key": "mykey",
  "value": "myvalue"
}
```

**Response (Error - Not Found):**
```json
{
  "error": "key not found"
}
```

**Status Codes:**
- `200 OK`: Key found
- `400 Bad Request`: Key parameter missing
- `404 Not Found`: Key doesn't exist or has expired

**Examples:**

```bash
curl http://localhost:8080/cache/user:123
```

---

### 4. Delete Key

**Endpoint:** `DELETE /cache/{key}`

**Description:** Removes a key from the cache.

**Response (Success):**
```json
{
  "status": "deleted",
  "key": "mykey"
}
```

**Response (Error):**
```json
{
  "error": "key not found"
}
```

**Status Codes:**
- `200 OK`: Successfully deleted
- `400 Bad Request`: Key parameter missing
- `404 Not Found`: Key doesn't exist

**Examples:**

```bash
curl -X DELETE http://localhost:8080/cache/user:123
```

---

### 5. List All Keys

**Endpoint:** `GET /cache`

**Description:** Returns all non-expired keys with their expiry information.

**Response (Success):**
```json
{
  "keys": [
    {
      "key": "permanent_key",
      "expires_at": null
    },
    {
      "key": "session:abc",
      "expires_at": "2025-02-16T12:30:45Z"
    }
  ]
}
```

**Status Codes:**
- `200 OK`: Always successful

**Response Details:**
- `expires_at` is `null` for keys without TTL
- `expires_at` is an RFC3339 timestamp for keys with TTL
- Keys are sorted alphabetically
- Only non-expired keys are included

**Examples:**

```bash
curl http://localhost:8080/cache
```

---

## Installation

### Prerequisites

- Go 1.21 or later
- Docker (for containerized deployment)

### Building Locally

```bash
# Clone or download the project
cd weather-cache-api

# Download dependencies
go mod tidy

# Build the binary
go build -o weather-cache-api

# Run the service
./weather-cache-api
```

The service will start on port 8080 by default.

### Using Docker

**Build the image:**
```bash
docker build -t weather-cache-api:latest .
```

**Run a container:**
```bash
docker run -d \
  -p 8080:8080 \
  --name weather-cache \
  weather-cache-api:latest
```

**With custom port:**
```bash
docker run -d \
  -p 9000:9000 \
  -e PORT=9000 \
  --name weather-cache \
  weather-cache-api:latest
```

### Environment Variables

- `PORT`: The port the service listens on (default: `8080`)

```bash
# Set custom port
export PORT=9000
./weather-cache-api
```

---

## Running Tests

### Run all tests:
```bash
go test -v ./...
```

### Run tests for cache package only:
```bash
go test -v ./internal/cache
```

### Run tests for handlers package only:
```bash
go test -v ./internal/handlers
```

### Run with coverage:
```bash
go test -cover ./...
```

### Run with detailed coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Cache Behavior

### TTL (Time-To-Live)

- When you set a key with `ttl_seconds`, the key automatically expires after that duration
- Expired keys are treated as non-existent (GET returns 404)
- A background goroutine cleans up expired entries every 10 seconds
- Keys without TTL never expire (unless explicitly deleted)

### Thread Safety

- All cache operations are thread-safe
- Multiple goroutines can safely read and write concurrently
- Uses `sync.RWMutex` for efficient concurrent access

### Memory Management

- In-memory only (data is lost on restart)
- No persistence layer
- Suitable for caching temporary data like sessions, rate limits, etc.

---

## Usage Example: Session Cache

```bash
# Store a session
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "session:user123",
    "value": "{\"user_id\": 123, \"username\": \"john\"}",
    "ttl_seconds": 3600
  }'

# Response: 201 Created
# {"status":"ok","key":"session:user123"}

# Retrieve the session
curl http://localhost:8080/cache/session:user123

# Response: 200 OK
# {"key":"session:user123","value":"{\"user_id\": 123, \"username\": \"john\"}"}

# List all sessions
curl http://localhost:8080/cache

# Response: 200 OK
# {"keys":[{"key":"session:user123","expires_at":"2025-02-16T13:30:45Z"}]}

# Delete the session
curl -X DELETE http://localhost:8080/cache/session:user123

# Response: 200 OK
# {"status":"deleted","key":"session:user123"}
```

---

## Architecture

### In-Memory Cache Store (`internal/cache/`)

- **Thread-Safe**: Uses `sync.RWMutex` for concurrent access
- **TTL Support**: Tracks expiration time for each entry
- **Automatic Cleanup**: Background goroutine removes expired entries

**Key Methods:**
- `Set(key, value string, ttl *time.Duration)`: Store a key-value pair
- `Get(key string) (value string, ok bool)`: Retrieve a value
- `Delete(key string) bool`: Remove a key
- `Keys() []KeyInfo`: List all non-expired keys
- `CleanExpired()`: Remove expired entries

### HTTP Handlers (`internal/handlers/`)

- **RESTful Design**: Standard HTTP methods (GET, POST, DELETE)
- **JSON API**: All requests/responses use JSON
- **Error Handling**: Consistent error responses with status codes

**Handler Methods:**
- `PostCache`: Store a key-value pair
- `GetCacheByKey`: Retrieve by key
- `DeleteCacheByKey`: Delete by key
- `ListCache`: List all keys
- `HealthCheck`: Service health status

### Main Server (`main.go`)

- **HTTP Server**: Uses Go's net/http package
- **Request Router**: Custom routing logic in `cacheRouter`
- **Background Cleanup**: Goroutine runs cleanup every 10 seconds
- **Graceful Startup**: Logs when service starts

---

## Performance Characteristics

- **Read Performance**: O(1) average case (hash map lookup)
- **Write Performance**: O(1) average case (hash map insertion)
- **Delete Performance**: O(1) average case (hash map deletion)
- **List Performance**: O(n) where n = number of keys (single scan)
- **Cleanup Performance**: O(n) every 10 seconds (full scan for expired items)

---

## Limitations

- **In-Memory Only**: Data is lost when the service restarts
- **Single Instance**: No replication or clustering
- **No Persistence**: No disk storage or database backend
- **No Authentication**: No access control or API keys
- **No Rate Limiting**: No built-in rate limit protection

---

## Future Enhancements

- [ ] Persistent storage backend (Redis integration)
- [ ] API authentication and authorization
- [ ] Rate limiting per client
- [ ] Metrics and monitoring (Prometheus)
- [ ] Distributed cache (cluster support)
- [ ] Data export/import functionality
- [ ] Admin dashboard

---

## License

MIT License

---

## Support

For issues, feature requests, or contributions, please refer to the project repository.
