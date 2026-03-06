# Hello World HTTP Server

A simple Go HTTP server with two endpoints.

## Endpoints

- `GET /` - Returns "Hello, World!"
- `GET /health` - Returns JSON status information

## Requirements

- Go 1.22 or higher

## Running the Server

```bash
go run main.go
```

The server will start on port 8080. You should see:
```
Starting server on :8080
```

## Testing

Test the root endpoint:
```bash
curl http://localhost:8080/
```

Test the health check:
```bash
curl http://localhost:8080/health
```

Expected health response:
```json
{
  "status": "ok",
  "timestamp": "2026-03-06T12:34:56Z",
  "service": "hello-world-server"
}
```
