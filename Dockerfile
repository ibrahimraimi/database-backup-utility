# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
  ca-certificates \
  mysql-client \
  postgresql-client \
  mongodb-tools \
  sqlite

# Create non-root user
RUN adduser -D -s /bin/sh dbbackup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/build/dbu /usr/local/bin/dbu

# Create backup directory
RUN mkdir -p /backups && chown dbbackup:dbbackup /backups

# Switch to non-root user
USER dbbackup

# Set default environment variables
ENV BACKUP_DIR=/backups
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json

# Expose volume for backups
VOLUME ["/backups"]

# Set entrypoint
ENTRYPOINT ["dbu"]

# Default command
CMD ["--help"]
