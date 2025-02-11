# Phonon Service

Phonon is a scalable audio processing service that handles audio file uploads, format conversions, and storage management. It provides a robust API for audio file operations with asynchronous processing capabilities using Kafka for better scalability.

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database Design](#database-design)
- [Development](#development)
- [Improvement Plan](#improvement-plan)

## Overview

Phonon is designed to handle audio file processing at scale. It supports various audio formats and provides a clean REST API for file operations. The service uses a microservices architecture with the following main components:

- REST API for file uploads and retrievals
- Kafka-based message queue for asynchronous processing
- Audio format converter using FFmpeg
- Flexible storage backend (Local/S3)
- Cleanup service (Janitor) for temporary files

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│             │         │             │         │             │
│  REST API   ├────────►│   Kafka    ├────────►│  Converter  │
│             │         │             │         │             │
└─────────────┘         └─────────────┘         └─────────────┘
       ▲                                               │
       │                                               ▼
       │                                        ┌─────────────┐
       │                                        │             │
       └────────────────────────────────────────┤   Storage   │
                                                │             │
                                                └─────────────┘
```

### Component Details

1. **REST API**
   - Handles file uploads and downloads
   - User authentication and authorization
   - Request validation and error handling

2. **Kafka Queue**
   - Manages asynchronous processing tasks
   - Handles cleanup job distribution
   - Ensures system scalability

3. **Converter**
   - FFmpeg-based audio format conversion
   - Supports multiple audio formats
   - Optimized for performance

4. **Storage**
   - Pluggable storage backend
   - Supports local filesystem and S3
   - Efficient file management

## Features

- Audio file upload and download
- Multiple audio format support
- Asynchronous processing
- Scalable architecture
- Automatic temporary file cleanup
- Configurable storage backend

## Getting Started

### Prerequisites

- Go 1.19 or later
- Docker and Docker Compose
- FFmpeg
- MySQL/SQLite
- Kafka

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/phonon.git
   cd phonon
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the service:
   ```bash
   make build
   ```

### Configuration

The service uses a `config.yaml` file for configuration. Here's an example configuration:

```yaml
server:
  port: 8080
  host: localhost

database:
  driver: mysql
  dsn: user:password@tcp(localhost:3306)/phonon

kafka:
  brokers:
    - localhost:9092
  topics:
    cleanup: cleanup-topic

storage:
  type: local
  path: /path/to/storage
  # For S3:
  # type: s3
  # bucket: your-bucket
  # region: us-west-2
```

## API Documentation

### Upload Audio
```http
POST /api/v1/users/{user_id}/phrases/{phrase_id}/audio
Content-Type: multipart/form-data

Form Data:
- audio_file: The audio file to upload
```

### Get Audio
```http
GET /api/v1/users/{user_id}/phrases/{phrase_id}/audio/{format}
```

## Database Design

### Users Table
```sql
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Audio Files Table
```sql
CREATE TABLE audio_files (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    phrase_id INT NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    format VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## Development

### Makefile Commands

- `make build`: Build the service
- `make test`: Run tests
- `make run`: Run the service
- `make docker`: Build Docker image
- `make docker-compose`: Run with Docker Compose
- `make clean`: Clean build artifacts

### Project Structure

```
├── cmd/                 # Command line applications
│   ├── janitor/        # Cleanup service
│   └── phonon/         # Main service
├── pkg/                # Package code
│   ├── api/            # REST API handlers
│   ├── converter/      # Audio conversion
│   ├── queue/          # Message queue
│   ├── service/        # Business logic
│   └── storage/        # Storage backends
└── scripts/            # Utility scripts
```

## Improvement Plan

### Short Term
1. Add user authentication and authorization
2. Implement rate limiting
3. Add API documentation using Swagger
4. Improve error handling and logging
5. Add metrics and monitoring

### Medium Term
1. Add support for more audio formats
2. Implement audio processing pipeline
3. Add caching layer
4. Improve test coverage
5. Add API versioning

### Long Term
1. Implement streaming support
2. Add audio analysis features
3. Support for multiple storage backends
4. Implement audio processing plugins
5. Add support for batch processing

## Design Decision
1. Using Asynchronous processing for better scalability as it doesn't explicitly stated that the API should reflect the audio retrieval