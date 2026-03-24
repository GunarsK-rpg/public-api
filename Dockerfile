# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy source and tidy dependencies
COPY . .
RUN go mod tidy && go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o service ./cmd/api

# Production stage
FROM alpine:3.23.3

# Security update - CACHE_BUST is set by CI to force fresh apk upgrade
ARG CACHE_BUST
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/service .

# Set ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port
EXPOSE 8090

# Run the binary
CMD ["./service"]
