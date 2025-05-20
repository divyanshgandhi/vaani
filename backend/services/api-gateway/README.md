# API Gateway Service

The API Gateway service is the main entry point for the Vaani platform. It handles user authentication, request routing, rate limiting, and streaming of previews.

## Features

- HTTP/WebSocket gateway
- User authentication via Firebase PhoneAuth JWT
- Request routing to appropriate services
- Rate limiting (token bucket per user)
- Quota enforcement
- Language detection
- Firestore integration for job persistence

## Development

### Prerequisites

- Go 1.21 or higher
- Make

### Setup

1. Clone the repository
2. Navigate to the api-gateway directory
3. Install dependencies:

```bash
make deps
```

### Running Locally

```bash
make run
```

Or for development with hot reload (requires [Air](https://github.com/cosmtrek/air)):

```bash
make dev
```

### Testing

```bash
make test
```

### Building

```bash
make build
```

## API Endpoints

- `GET /healthz` - Health check endpoint
- `POST /v1/generate` - Generate TTS audio (protected)
- `WebSocket /v1/preview` - Stream audio previews (protected)
- `POST /v1/detect-language` - Detect language of text (public)

### Generate TTS Audio

The generate endpoint (`POST /v1/generate`) accepts text and voice parameters to create a TTS audio job.

Example request:
```json
{
  "text": "Your text to synthesize as speech",
  "voice_id": "voice-id-here",
  "speed": 1.0,
  "pitch": 0,
  "emotion": 50,
  "language": "en",
  "output_type": "mp3"
}
```

Example success response (202 Accepted):
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "queued",
  "queue_time": "2023-04-25T12:34:56Z",
  "message": "Job has been queued for processing"
}
```

Validation:
- Text must not be empty
- Text must be 5000 characters or less
- Voice ID is required

### Preview Audio Stream

The preview endpoint (`WebSocket /v1/preview`) allows for streaming real-time audio previews. This is useful for giving users immediate feedback before committing to a full TTS generation.

To use the WebSocket endpoint:

1. Connect to `/v1/preview` using the WebSocket protocol
2. Send a JSON message with the text to preview:

```json
{
  "type": "text",
  "text": "Text to preview as speech",
  "voice_id": "voice-id-here",
  "speed": 1.0,
  "pitch": 0,
  "emotion": 50
}
```

The server will respond with a series of messages:

1. Start message:
```json
{
  "type": "start"
}
```

2. Audio chunks (multiple):
```json
{
  "type": "audio",
  "data": "base64-encoded-audio-data",
  "duration": 500
}
```

3. End message:
```json
{
  "type": "end",
  "duration": 1500
}
```

Error responses:
```json
{
  "type": "error",
  "error": "Error message"
}
```

### Language Detection

The language detection endpoint (`POST /v1/detect-language`) allows you to determine the language of a given text. This is used internally for routing TTS requests to the appropriate backend (Bulbul for Indic languages, Orpheus for English).

Example request:
```json
{
  "text": "Your text to detect language here"
}
```

Example response:
```json
{
  "language": "en",
  "language_name": "English",
  "is_indic": false,
  "text_length": 30,
  "supported_languages": ["en", "hi", "ta", "te", "kn", "ml", "bn", "gu", "mr"]
}
```

Supported languages:
- English (en)
- Hindi (hi)
- Tamil (ta)
- Telugu (te)
- Kannada (kn)
- Malayalam (ml)
- Bengali (bn)
- Gujarati (gu)
- Marathi (mr)

## Rate Limiting

The API Gateway uses a token bucket rate limiter to control request rates on a per-user basis. When authenticated, the rate limit is applied using the user's ID. For unauthenticated requests, the client IP address is used.

By default, clients are limited to:
- 5 requests per second
- Bucket resets fully after 60 seconds of inactivity

When a client exceeds their rate limit, they receive a 429 Too Many Requests response with a JSON body containing:
```json
{
  "error": "Rate limit exceeded",
  "requests_per_sec": 5,
  "reset_after_sec": 60,
  "retry_after_sec": 1
}
```

The API also sets the `Retry-After` header to indicate when the client can retry.

## Environment Variables

- `PORT` - Port to listen on (default: 8080)
- `APP_ENV` - Application environment (development, production)
- `LOG_LEVEL` - Logging level (debug, info, error)

### Rate Limiting

- `RATE_LIMIT_ENABLED` - Enable rate limiting (default: true)
- `RATE_LIMIT_REQUESTS_PER_SECOND` - Maximum requests per second per user (default: 5)
- `RATE_LIMIT_RESET_AFTER_SECONDS` - Duration in seconds after which to reset rate limits (default: 60)

### Firebase Authentication & Firestore

- `FIREBASE_PROJECT_ID` - Firebase project ID
- `FIREBASE_CREDENTIALS_FILE` - Path to Firebase service account JSON
- `FIREBASE_USE_EMULATOR` - Set to "true" to use Firebase emulators
- `FIREBASE_AUTH_EMULATOR_HOST` - Firebase Auth emulator host:port
- `FIREBASE_FIRESTORE_EMULATOR_HOST` - Firebase Firestore emulator host:port

## Authentication

This service uses Firebase Phone Authentication. Protected routes require a valid Firebase ID token in the Authorization header:

```
Authorization: Bearer <firebase-id-token>
```

### Setting Up Firebase

1. Create a Firebase project at https://console.firebase.google.com/
2. Enable Phone Authentication in the Authentication > Sign-in method section
3. Generate a service account key in Project settings > Service accounts
4. Save the JSON file to `credentials/firebase-service-account.json` (gitignored)
5. Set the `FIREBASE_PROJECT_ID` environment variable to your Firebase project ID

### Testing with Firebase Emulator

For local development, you can use the Firebase Auth Emulator:

1. Install the Firebase CLI: `npm install -g firebase-tools`
2. Initialize Firebase in a separate directory: `firebase init`
3. Start the emulator: `firebase emulators:start --only auth`
4. Set environment variables:
   ```
   FIREBASE_USE_EMULATOR=true
   FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
   ```

## Docker

Build and run with Docker:

```bash
docker build -t vaani-api-gateway .
docker run -p 8080:8080 vaani-api-gateway
```

Or using Docker Compose:

```bash
docker-compose up
``` 