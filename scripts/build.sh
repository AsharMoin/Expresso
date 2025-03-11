#!/bin/bash

VERSION="v1.0.0"  # Update this for each release

# Navigate to the project root (one level up from scripts directory)
cd "$(dirname "$0")/.."

# Create build directory
mkdir -p builds

# Build for different platforms
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o builds/expresso .
tar -czf builds/expresso_${VERSION}_linux_amd64.tar.gz -C builds expresso

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o builds/expresso .
tar -czf builds/expresso_${VERSION}_darwin_amd64.tar.gz -C builds expresso

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o builds/expresso .
tar -czf builds/expresso_${VERSION}_darwin_arm64.tar.gz -C builds expresso

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o builds/expresso.exe .
zip -j builds/expresso_${VERSION}_windows_amd64.zip builds/expresso.exe

echo "Build complete! Archives ready for release in the builds directory:"
ls -la builds/expresso_${VERSION}_*.tar.gz builds/expresso_${VERSION}_*.zip