# Multi-stage Dockerfile for Go URL shortening service
FROM golang:1.25-alpine AS builder

# install necessary build tools
# git is needed to fetch dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
# CGO_ENABLED=0: Disable CGO for static binary
# -ldflags="-w -s": Strip debugging information to reduce binary size
RUN CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o /app/bin/main \
    ./cmd/api


# Start a new stage from scratch
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Taipei

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/main .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 80

# RUN the application
CMD ["sh"]