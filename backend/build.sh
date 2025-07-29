#!/bin/bash

# Build script for Go Chat Backend

echo "🚀 Building Go Chat Backend..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f chat-server chat-server.exe

# Build for current platform
echo "🔨 Building for current platform..."
go build -o chat-server main.go

# Build for different platforms (optional)
if [ "$1" == "all" ]; then
    echo "🌍 Building for multiple platforms..."
    
    # Windows
    GOOS=windows GOARCH=amd64 go build -o chat-server-windows-amd64.exe main.go
    
    # Linux
    GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 main.go
    
    # macOS
    GOOS=darwin GOARCH=amd64 go build -o chat-server-darwin-amd64 main.go
    
    echo "✅ Multi-platform builds complete!"
else
    echo "✅ Build complete!"
fi

echo "🎉 Ready to run: ./chat-server"
