# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /app/bin/branch-detector \
    ./services/branch-detector

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/branch-detector .

# Expose ports
EXPOSE 50051 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

ENTRYPOINT ["./branch-detector"]

