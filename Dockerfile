# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary for linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o branch-aware-ci .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/branch-aware-ci /app/branch-aware-ci

# Ensure the binary is executable
RUN chmod +x /app/branch-aware-ci

ENTRYPOINT ["/app/branch-aware-ci"]

