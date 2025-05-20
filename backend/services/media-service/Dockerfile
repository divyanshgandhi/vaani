FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o media-service .

# Use a smaller image for the final build
FROM alpine:latest  

WORKDIR /app

# Create a directory for media storage
RUN mkdir -p /app/data

# Copy the binary from builder
COPY --from=builder /app/media-service .

# Expose port
EXPOSE 8083

# Set volume for media data
VOLUME ["/app/data"]

# Run the binary
CMD ["./media-service"] 