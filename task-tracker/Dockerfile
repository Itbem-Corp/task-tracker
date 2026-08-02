# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy dependency files first to leverage Docker layer caching.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the application source code.
COPY . .

# Run tests before building to catch regressions early.
RUN go test ./...

# Build a statically-linked binary with no CGO dependencies.
# -ldflags="-s -w" strips debug info and symbol tables to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /task-tracker ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.20

# Install only CA certificates — nothing else.
RUN apk --no-cache add ca-certificates \
    && addgroup -S appgroup && adduser -S -G appgroup -u 10001 appuser

WORKDIR /app

# Copy the compiled binary from the build stage.
COPY --from=builder --chown=appuser:appgroup /task-tracker ./task-tracker

# Expose the port the application listens on.
EXPOSE 8080

ENV PORT=8080

# Drop to non-root user.
USER appuser

ENTRYPOINT ["./task-tracker"]
