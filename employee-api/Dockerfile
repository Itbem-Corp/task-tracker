# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod first to leverage Docker layer caching.
# If go.mod hasn't changed, Docker reuses the cached layer and skips
# re-downloading dependencies on subsequent builds.
COPY go.mod ./

# Copy the application source code.
COPY . .

# Build a statically-linked binary with no CGO dependencies.
# -ldflags="-s -w" strips debug info and symbol tables to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ---- Runtime stage ----
FROM alpine:3.20

# Install CA certificates so the application can make outbound HTTPS calls.
RUN apk --no-cache add ca-certificates

# Create a non-root user and group for running the application.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the compiled binary from the build stage.
COPY --from=builder /app/server .

# Expose the port the application listens on.
EXPOSE 8080

# Set the default port as an environment variable.
ENV PORT=8080

# Switch to the non-root user.
USER appuser

# Use the binary as the container entrypoint.
ENTRYPOINT ["./server"]
