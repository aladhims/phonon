# Build stage
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o phonon ./cmd/phonon

# Final stage
FROM debian:bookworm-slim

# Install ffmpeg and required libraries
RUN apt-get update && apt-get install -y \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/phonon .
# Copy config file
COPY config.yaml .

EXPOSE 8080

CMD ["./phonon"]