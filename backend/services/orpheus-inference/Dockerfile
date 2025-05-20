FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orpheus-inference .

# Use a smaller image for the final build
FROM alpine:latest  

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/orpheus-inference .

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./orpheus-inference"] 