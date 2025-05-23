# Vaani Backend

A microservices-based backend for the Vaani voice AI studio, providing text-to-speech generation for English and 11 Indian languages.

## Architecture

The backend consists of 4 main services:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   API Gateway   │────│  Bulbul Adapter  │────│   Sarvam API    │
│   (Port 8080)   │    │   (Port 8082)    │    │   (External)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │
         ├─────────────────┐
         │                 │
┌─────────────────┐    ┌──────────────────┐
│ Orpheus Service │    │  Media Service   │
│   (Port 8081)   │    │   (Port 8083)    │
└─────────────────┘    └──────────────────┘
```

### Services

1. **API Gateway** (`services/api-gateway/`)
   - Main entry point for all client requests
   - Handles authentication (Firebase JWT)
   - Routes requests to appropriate TTS services
   - Manages job queue and status tracking
   - Language detection and routing logic

2. **Bulbul Adapter** (`services/bulbul-adapter/`)
   - Adapter for Sarvam AI's Bulbul v2 API
   - Handles Indic language TTS generation
   - Emotion mapping and voice selection
   - Response caching for cost optimization

3. **Orpheus Inference** (`services/orpheus-inference/`)
   - Self-hosted Orpheus TTS for English/Hinglish
   - Zero-shot voice cloning capabilities
   - Streaming audio generation
   - GPU-accelerated inference

4. **Media Service** (`services/media-service/`)
   - Audio file processing and encoding
   - Storage management (Wasabi/S3)
   - Pre-signed URL generation
   - Format conversion (WAV/MP3)

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Firebase project (for authentication)
- Sarvam AI API key (for Indic TTS)

### 1. Build All Services

```bash
# Make build script executable and run
chmod +x test-build.sh
./test-build.sh
```

### 2. Environment Setup

Copy environment files and configure:

```bash
# API Gateway
cp services/api-gateway/env.example services/api-gateway/.env

# Edit the .env file with your configuration:
# - FIREBASE_PROJECT_ID
# - SARVAM_API_KEY
# - Service URLs (if different from defaults)
```

### 3. Run with Docker Compose

```bash
# Start all services
make dev-up

# Stop all services
make dev-down
```

### 4. Run Individual Services

```bash
# Terminal 1: API Gateway
cd services/api-gateway
./api-gateway

# Terminal 2: Bulbul Adapter
cd services/bulbul-adapter
./bulbul-adapter

# Terminal 3: Orpheus Inference
cd services/orpheus-inference
./orpheus-inference

# Terminal 4: Media Service
cd services/media-service
./media-service
```

## API Endpoints

### Public Endpoints

- `GET /healthz` - Health check
- `POST /v1/detect-language` - Language detection

### Protected Endpoints (Require Authentication)

- `POST /v1/generate` - Generate TTS audio
- `GET /v1/jobs/{jobID}` - Get job status
- `GET /v1/download/{jobID}` - Download generated audio
- `GET /v1/preview` - WebSocket for real-time preview
- `POST /v1/clone` - Voice cloning (placeholder)
- `GET /v1/voices` - List available voices (placeholder)

### Example Usage

```bash
# 1. Generate TTS
curl -X POST http://localhost:8080/v1/generate \
  -H "Authorization: Bearer YOUR_FIREBASE_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "नमस्ते दुनिया",
    "voice_id": "hi_female_1",
    "emotion": 50,
    "language": "hi"
  }'

# Response:
# {
#   "job_id": "uuid-here",
#   "status": "queued",
#   "queue_time": "2024-01-01T12:00:00Z"
# }

# 2. Check job status
curl -X GET http://localhost:8080/v1/jobs/uuid-here \
  -H "Authorization: Bearer YOUR_FIREBASE_JWT"

# 3. Download when complete
curl -X GET http://localhost:8080/v1/download/uuid-here \
  -H "Authorization: Bearer YOUR_FIREBASE_JWT"
```

## Configuration

### Environment Variables

#### API Gateway
```bash
# Application
APP_ENV=development
PORT=8080
LOG_LEVEL=info

# Firebase
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_CREDENTIALS_FILE=./credentials/firebase-service-account.json

# Service URLs
BULBUL_ADAPTER_URL=http://localhost:8082
ORPHEUS_SERVICE_URL=http://localhost:8081
MEDIA_SERVICE_URL=http://localhost:8083

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_SECOND=5
```

#### Bulbul Adapter
```bash
PORT=8082
SARVAM_API_KEY=your-sarvam-api-key
```

#### Orpheus Inference
```bash
PORT=8081
MODEL_PATH=/path/to/orpheus/model
GPU_MEMORY_FRACTION=0.8
```

#### Media Service
```bash
PORT=8083
WASABI_ACCESS_KEY=your-access-key
WASABI_SECRET_KEY=your-secret-key
WASABI_BUCKET=vaani-media-storage
```

## Development

### Adding New Features

1. **New API Endpoint**: Add to `services/api-gateway/router/router.go`
2. **New Handler**: Create in `services/api-gateway/handlers/`
3. **New Service**: Follow the existing service structure
4. **Database Changes**: Update models in `services/api-gateway/models/`

### Testing

```bash
# Run tests for all services
make test

# Run tests for specific service
cd services/api-gateway
go test ./...
```

### Code Structure

```
backend/
├── services/
│   ├── api-gateway/
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # Auth, rate limiting, etc.
│   │   ├── models/           # Data models
│   │   ├── database/         # Firestore client
│   │   ├── services/         # Business logic
│   │   ├── utils/            # Utilities (lang detection)
│   │   └── config/           # Configuration
│   ├── bulbul-adapter/       # Sarvam API adapter
│   ├── orpheus-inference/    # Orpheus TTS service
│   └── media-service/        # Media processing
├── infra/
│   ├── docker/              # Docker compose files
│   └── terraform/           # Infrastructure as code
└── libs/                    # Shared libraries
```

## Deployment

### Production Deployment

1. **Build Docker Images**:
   ```bash
   make docker
   ```

2. **Deploy to Cloud**:
   - Use provided Terraform configurations
   - Set up proper secrets management
   - Configure monitoring and logging

3. **Environment-Specific Configs**:
   - Development: Local services, Firebase emulator
   - Staging: Cloud services, test data
   - Production: Full cloud setup, monitoring

### Monitoring

- **Logs**: Structured JSON logging with Zap
- **Metrics**: Prometheus-compatible metrics
- **Health Checks**: `/healthz` endpoints on all services
- **Tracing**: OpenTelemetry integration (planned)

## Troubleshooting

### Common Issues

1. **Firebase Authentication Errors**:
   - Check `FIREBASE_PROJECT_ID` is correct
   - Verify service account credentials
   - Ensure JWT tokens are valid

2. **Sarvam API Errors**:
   - Verify `SARVAM_API_KEY` is set
   - Check API quota and rate limits
   - Review supported languages and voices

3. **Service Communication Errors**:
   - Verify all services are running
   - Check service URLs in environment
   - Review network connectivity

4. **Build Errors**:
   - Ensure Go 1.21+ is installed
   - Run `go mod download` in each service
   - Check for missing dependencies

### Debug Mode

```bash
# Enable debug logging
export LOG_LEVEL=debug

# Run with verbose output
./api-gateway --verbose
```

## Contributing

1. Follow Go best practices and conventions
2. Add tests for new functionality
3. Update documentation for API changes
4. Use structured logging with appropriate levels
5. Handle errors gracefully with proper HTTP status codes

## License

This project is part of the Vaani voice AI studio. See the main repository for license information. 