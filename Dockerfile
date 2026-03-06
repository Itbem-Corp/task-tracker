# Multi-stage build for Go hello-world server
# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary
# CGO_ENABLED=0 for static binary, compatible with distroless
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server .

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server /app/server

# Use non-root user (distroless nonroot is uid 65532)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/app/server"]
