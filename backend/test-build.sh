#!/bin/bash

# Test build script for Vaani backend services
set -e

echo "🔨 Building Vaani backend services..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Build each service
services=("api-gateway" "bulbul-adapter" "orpheus-inference" "media-service")

for service in "${services[@]}"; do
    echo ""
    echo "Building $service..."
    
    cd "services/$service"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found in $service"
        cd ../..
        continue
    fi
    
    # Download dependencies
    echo "  Downloading dependencies..."
    go mod download
    
    # Run tests if they exist
    if ls *_test.go 1> /dev/null 2>&1; then
        echo "  Running tests..."
        if go test ./...; then
            print_status "Tests passed for $service"
        else
            print_warning "Some tests failed for $service"
        fi
    else
        print_warning "No tests found for $service"
    fi
    
    # Build the service
    echo "  Building binary..."
    if go build -o "$service" .; then
        print_status "Successfully built $service"
    else
        print_error "Failed to build $service"
        cd ../..
        exit 1
    fi
    
    cd ../..
done

echo ""
print_status "All services built successfully!"

echo ""
echo "🚀 To run the services locally:"
echo "  1. Set up environment variables (copy env.example files)"
echo "  2. Run: make dev-up"
echo "  3. Or run individual services:"
for service in "${services[@]}"; do
    echo "     cd services/$service && ./$service"
done

echo ""
echo "📋 Next steps:"
echo "  - Configure Firebase credentials"
echo "  - Set up Sarvam API key for Bulbul adapter"
echo "  - Configure Wasabi/S3 storage for media service"
echo "  - Set up Orpheus TTS model (see orpheus-inference/README.md)" 