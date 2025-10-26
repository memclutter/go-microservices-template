#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Generating protobuf code..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    exit 1
fi

# Generate Go code and gRPC gateway
protoc \
  --proto_path=api/proto \
  --proto_path=third_party \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=generate_unbound_methods=true \
  api/proto/**/*.proto

echo -e "${GREEN}✓ Protobuf code generated successfully${NC}"
