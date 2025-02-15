# Build stage
FROM golang:1.23-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o phonon ./cmd/phonon

# Final stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/phonon .
COPY config.yaml .

EXPOSE 8080

CMD ["./phonon"]