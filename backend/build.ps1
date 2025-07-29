# Build script for Go Chat Backend (PowerShell)

Write-Host "Building Go Chat Backend..." -ForegroundColor Green

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue "chat-server", "chat-server.exe", "*.exe"

# Build for current platform
Write-Host "Building for current platform..." -ForegroundColor Blue
go build -o chat-server.exe main.go

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build complete!" -ForegroundColor Green
    Write-Host "Ready to run: .\chat-server.exe" -ForegroundColor Cyan
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Build for different platforms if requested
if ($args[0] -eq "all") {
    Write-Host "Building for multiple platforms..." -ForegroundColor Magenta
    
    # Windows 64-bit
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o chat-server-windows-amd64.exe main.go
    
    # Linux 64-bit
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go build -o chat-server-linux-amd64 main.go
    
    # macOS 64-bit
    $env:GOOS = "darwin"
    $env:GOARCH = "amd64"
    go build -o chat-server-darwin-amd64 main.go
    
    # Reset environment variables
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    
    Write-Host "Multi-platform builds complete!" -ForegroundColor Green
}
