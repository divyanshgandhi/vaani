#!/bin/bash
set -e

# This script generates code from protobuf definitions for both Go and Dart

# Ensure protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed. Please install Protocol Buffers compiler."
    exit 1
fi

echo "Generating protobuf code..."

# Generate Go code
BACKEND_OUT_DIR="../../backend/libs/protobufs"
mkdir -p $BACKEND_OUT_DIR

# Generate Dart code
FRONTEND_OUT_DIR="../../frontend/packages/vaani_models/lib/src/generated"
mkdir -p $FRONTEND_OUT_DIR

# Find all proto files
PROTO_FILES=$(find . -name "*.proto")

for proto_file in $PROTO_FILES; do
    echo "Processing $proto_file..."
    
    # Generate Go code
    protoc \
        --go_out=$BACKEND_OUT_DIR \
        --go_opt=paths=source_relative \
        --go-grpc_out=$BACKEND_OUT_DIR \
        --go-grpc_opt=paths=source_relative \
        $proto_file

    # Generate Dart code
    protoc \
        --dart_out=grpc:$FRONTEND_OUT_DIR \
        $proto_file
done

echo "Generated protobuf code successfully!" 