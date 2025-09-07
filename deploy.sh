#!/bin/bash

# Deploy script for Aloha Email Processor
set -e

echo "🚀 Deploying Aloha Email Processor..."

# Build Lambda function
echo "📦 Building Lambda function..."
./build-lambda.sh

# Synthesize CDK stack
echo "🔧 Synthesizing CDK stack..."
go run main.go

echo "✅ Build completed successfully!"
echo ""
echo "To deploy to AWS, run:"
echo "  cdk deploy"
echo ""
echo "To destroy the stack, run:"
echo "  cdk destroy"