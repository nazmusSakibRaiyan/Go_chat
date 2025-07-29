# Development Setup Script for Go Chat

Write-Host "Setting up Go Chat Development Environment..." -ForegroundColor Green

# Check if Go is installed
Write-Host "Checking Go installation..." -ForegroundColor Blue
$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Go is not installed! Please install Go 1.21+ from https://golang.org/dl/" -ForegroundColor Red
    exit 1
}
Write-Host "Found: $goVersion" -ForegroundColor Green

# Setup backend
Write-Host "Setting up backend dependencies..." -ForegroundColor Blue
Set-Location backend
go mod tidy

if ($LASTEXITCODE -eq 0) {
    Write-Host "Backend dependencies installed!" -ForegroundColor Green
} else {
    Write-Host "Failed to install backend dependencies!" -ForegroundColor Red
    exit 1
}

# Check for .env file
if (-not (Test-Path ".env")) {
    Write-Host "No .env file found. Please create one with your MongoDB connection details" -ForegroundColor Yellow
    Write-Host "Example content:" -ForegroundColor Cyan
    Write-Host "MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/database" -ForegroundColor Gray
    Write-Host "PORT=8080" -ForegroundColor Gray
}

# Setup frontend (optional)
Set-Location ../frontend
if (Test-Path "package.json") {
    Write-Host "Setting up frontend dependencies..." -ForegroundColor Blue
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        npm install
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Frontend dependencies installed!" -ForegroundColor Green
        } else {
            Write-Host "Frontend setup failed, but you can still use the HTML file directly" -ForegroundColor Yellow
        }
    } else {
        Write-Host "npm not found. You can still use the HTML file directly" -ForegroundColor Yellow
    }
}

Set-Location ..

Write-Host ""
Write-Host "Setup complete!" -ForegroundColor Green
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Edit backend/.env with your MongoDB Atlas connection string" -ForegroundColor White
Write-Host "  2. Run: cd backend && go run main.go" -ForegroundColor White
Write-Host "  3. Open frontend/public/index.html in your browser" -ForegroundColor White
Write-Host "  4. Start chatting!" -ForegroundColor White
