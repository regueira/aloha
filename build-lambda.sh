#!/bin/bash

# Build script for Lambda function
set -e

echo "Building email processor Lambda function..."

cd lambda/email-processor

# Download dependencies
go mod tidy

# Build for Linux x86_64
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

echo "Lambda function built successfully!"
echo "Binary location: $(pwd)/bootstrap"