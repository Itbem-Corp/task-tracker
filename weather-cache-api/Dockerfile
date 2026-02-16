# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy module file first for layer caching
COPY go.mod ./

# Copy all source code
COPY . .

# Build static binary with stripped debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /weather-cache-api .

# Stage 2: Runtime
FROM alpine:3.19

# Install CA certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /weather-cache-api /weather-cache-api

# Expose API port
EXPOSE 8080

# Set default port
ENV PORT=8080

# Run as non-root
USER appuser

ENTRYPOINT ["/weather-cache-api"]
