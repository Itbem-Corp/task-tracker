# Weather Cache API - Quick Start Guide

## 30-Second Setup

### Run Locally
```bash
cd weather-cache-api
go build -o weather-cache-api
./weather-cache-api
```
✅ Service runs on `http://localhost:8080`

### Run in Docker
```bash
docker build -t weather-cache-api .
docker run -d -p 8080:8080 weather-cache-api
```

### Test with curl

**Store a value (with 1 hour TTL):**
```bash
curl -X POST http://localhost:8080/cache \
  -H "Content-Type: application/json" \
  -d '{"key": "mykey", "value": "myvalue", "ttl_seconds": 3600}'
```

**Get a value:**
```bash
curl http://localhost:8080/cache/mykey
```

**List all keys:**
```bash
curl http://localhost:8080/cache
```

**Delete a key:**
```bash
curl -X DELETE http://localhost:8080/cache/mykey
```

**Health check:**
```bash
curl http://localhost:8080/health
```

---

## API Reference

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | /health | Health check |
| POST | /cache | Store key-value |
| GET | /cache | List all keys |
| GET | /cache/{key} | Get value by key |
| DELETE | /cache/{key} | Delete key |

---

## Run Tests

```bash
go test -v ./...
```

---

## Configuration

### Port
Set `PORT` environment variable (default: 8080)
```bash
PORT=9000 ./weather-cache-api
```

---

## Key Features

- ✅ In-memory cache with optional TTL
- ✅ Thread-safe operations
- ✅ Automatic cleanup every 10 seconds
- ✅ Zero external dependencies
- ✅ Production-ready Docker image (~15MB)

---

For detailed documentation, see [README.md](./README.md)
