# 🛠️ Development Guide

Complete guide for developing and contributing to the Go Chat application.

## 📋 Table of Contents

- [Getting Started](#-getting-started)
- [Development Environment](#-development-environment)
- [Project Structure](#-project-structure)
- [Backend Development](#-backend-development)
- [Frontend Development](#-frontend-development)
- [Database Management](#-database-management)
- [Testing](#-testing)
- [Code Quality](#-code-quality)
- [Debugging](#-debugging)
- [Contributing](#-contributing)
- [Best Practices](#-best-practices)

## 🚀 Getting Started

### Prerequisites

- **Go**: 1.19 or higher
- **MongoDB**: 4.4 or higher
- **Git**: Latest version
- **Code Editor**: VS Code (recommended) or your preferred editor
- **Terminal**: PowerShell, Bash, or similar

### Quick Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/nazmusSakibRaiyan/Go_chat.git
   cd Go_chat
   ```

2. **Install Go dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   # Create .env file in backend directory
   echo "MONGO_URI=mongodb://localhost:27017" > .env
   echo "MONGO_DB_NAME=go_chat_dev" >> .env
   echo "JWT_SECRET=dev-secret-key-change-in-production" >> .env
   echo "PORT=8080" >> .env
   echo "GIN_MODE=debug" >> .env
   ```

4. **Start MongoDB**
   ```bash
   # Windows
   net start MongoDB
   
   # Linux/macOS
   sudo systemctl start mongod
   # OR
   brew services start mongodb-community
   ```

5. **Run the application**
   ```bash
   # Backend (from backend/ directory)
   go run main.go
   
   # Frontend (from frontend/public/ directory)
   python3 -m http.server 8000
   # OR
   npx serve -p 8000
   ```

## 💻 Development Environment

### VS Code Setup

**Recommended Extensions:**
```json
{
  "recommendations": [
    "golang.go",
    "mongodb.mongodb-vscode",
    "bradlc.vscode-tailwindcss",
    "ms-vscode.vscode-json",
    "esbenp.prettier-vscode",
    "ms-vscode.live-server"
  ]
}
```

**VS Code Settings** (`.vscode/settings.json`):
```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "files.autoSave": "afterDelay",
  "files.autoSaveDelay": 1000
}
```

**Launch Configuration** (`.vscode/launch.json`):
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Go Chat Backend",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/main.go",
      "cwd": "${workspaceFolder}/backend",
      "env": {
        "MONGO_URI": "mongodb://localhost:27017",
        "MONGO_DB_NAME": "go_chat_dev",
        "JWT_SECRET": "dev-secret-key",
        "PORT": "8080",
        "GIN_MODE": "debug"
      }
    }
  ]
}
```

### Go Development Tools

```bash
# Install essential Go tools
go install -a github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

**Air Configuration** (`.air.toml`):
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

## 🏗️ Project Structure

```
Go_chat/
├── backend/                    # Go backend application
│   ├── main.go                # Application entry point
│   ├── go.mod                 # Go module definition
│   ├── go.sum                 # Go module checksums
│   ├── .env                   # Environment variables (development)
│   ├── .air.toml              # Air live reload configuration
│   ├── auth/                  # Authentication module
│   │   ├── auth.go           # JWT utilities
│   │   ├── handlers.go       # Auth HTTP handlers
│   │   ├── middleware.go     # Auth middleware
│   │   └── ratelimit.go      # Rate limiting
│   ├── chat/                  # Chat functionality
│   │   └── chat.go           # WebSocket chat handlers
│   ├── config/                # Configuration management
│   │   └── config.go         # App configuration
│   ├── db/                    # Database layer
│   │   └── database.go       # MongoDB connection
│   ├── models/                # Data models
│   │   └── models.go         # User and message models
│   └── utils/                 # Utility functions
│       └── utils.go          # Helper functions
├── frontend/                   # Frontend application
│   ├── public/                # Static web files
│   │   ├── index.html        # Main chat interface
│   │   ├── index-auth.html   # Authentication demo
│   │   ├── api-tester.html   # API testing interface
│   │   └── launcher.html     # Frontend navigation
│   ├── package.json          # Frontend dependencies
│   └── README.md             # Frontend documentation
├── docs/                      # Documentation
│   ├── README.md             # Main project documentation
│   ├── API.md                # API documentation
│   ├── AUTHENTICATION.md     # Auth guide
│   ├── DEPLOYMENT.md         # Deployment guide
│   └── DEVELOPMENT.md        # This file
├── .gitignore                # Git ignore rules
├── docker-compose.yml        # Docker composition
├── Dockerfile               # Docker build configuration
└── README.md               # Project overview
```

### Module Organization

**Backend Architecture:**
```
┌─────────────────┐
│    main.go      │ ← Entry point, router setup
└─────────────────┘
         │
    ┌────▼────┐
    │ Config  │ ← Environment variables, settings
    └─────────┘
         │
    ┌────▼────┐
    │Database │ ← MongoDB connection
    └─────────┘
         │
    ┌────▼────┐
    │ Models  │ ← Data structures
    └─────────┘
         │
    ┌────▼────┐
    │  Auth   │ ← Authentication logic
    └─────────┘
         │
    ┌────▼────┐
    │  Chat   │ ← WebSocket chat
    └─────────┘
```

## 🔧 Backend Development

### Adding New API Endpoints

1. **Define the route in main.go**
   ```go
   // Add to setupRoutes function
   api.GET("/new-endpoint", handlers.NewEndpointHandler)
   ```

2. **Create handler function**
   ```go
   // In appropriate module (e.g., auth/handlers.go)
   func NewEndpointHandler(c *gin.Context) {
       // Implementation here
       c.JSON(http.StatusOK, gin.H{
           "message": "Success",
           "data":    responseData,
       })
   }
   ```

3. **Add middleware if needed**
   ```go
   // Protected endpoint
   api.GET("/protected", auth.RequireAuth(), handlers.ProtectedHandler)
   ```

### Database Operations

**Creating a new model:**
```go
// In models/models.go
type NewModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name      string            `bson:"name" json:"name"`
    CreatedAt time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}
```

**Database operations:**
```go
// Create
collection := database.GetCollection("new_models")
result, err := collection.InsertOne(context.Background(), newModel)

// Read
var model NewModel
err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&model)

// Update
update := bson.M{"$set": bson.M{"name": "new name", "updated_at": time.Now()}}
result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)

// Delete
result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
```

### Error Handling

**Standard error response:**
```go
func handleError(c *gin.Context, statusCode int, message string, err error) {
    if err != nil {
        log.Printf("Error: %v", err)
    }
    c.JSON(statusCode, gin.H{
        "error":   message,
        "success": false,
    })
}

// Usage
if err != nil {
    handleError(c, http.StatusInternalServerError, "Database error", err)
    return
}
```

### Middleware Development

**Creating custom middleware:**
```go
func CustomMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // Before request processing
        start := time.Now()
        
        // Process request
        c.Next()
        
        // After request processing
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
    })
}

// Apply middleware
router.Use(CustomMiddleware())
```

### Configuration Management

**Adding new config options:**
```go
// In config/config.go
type Config struct {
    MongoURI     string
    DatabaseName string
    JWTSecret    string
    Port         string
    NewOption    string `mapstructure:"NEW_OPTION"`
}

// Load from environment
func LoadConfig() (*Config, error) {
    config := &Config{
        MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
        DatabaseName: getEnv("MONGO_DB_NAME", "go_chat"),
        JWTSecret:    getEnv("JWT_SECRET", "secret"),
        Port:         getEnv("PORT", "8080"),
        NewOption:    getEnv("NEW_OPTION", "default_value"),
    }
    return config, nil
}
```

## 🎨 Frontend Development

### Adding New Pages

1. **Create HTML file in `frontend/public/`**
   ```html
   <!DOCTYPE html>
   <html lang="en">
   <head>
       <meta charset="UTF-8">
       <meta name="viewport" content="width=device-width, initial-scale=1.0">
       <title>New Page</title>
       <link rel="stylesheet" href="styles.css">
   </head>
   <body>
       <div id="app">
           <!-- Page content -->
       </div>
       <script src="js/new-page.js"></script>
   </body>
   </html>
   ```

2. **Create JavaScript file**
   ```javascript
   // js/new-page.js
   class NewPage {
       constructor() {
           this.init();
       }
   
       init() {
           this.setupEventListeners();
           this.loadData();
       }
   
       setupEventListeners() {
           // Event handlers
       }
   
       async loadData() {
           try {
               const response = await fetch('/api/data');
               const data = await response.json();
               this.renderData(data);
           } catch (error) {
               console.error('Error loading data:', error);
           }
       }
   
       renderData(data) {
           // Render data to DOM
       }
   }
   
   // Initialize when DOM is loaded
   document.addEventListener('DOMContentLoaded', () => {
       new NewPage();
   });
   ```

### API Integration

**Standard API call pattern:**
```javascript
class APIClient {
    constructor() {
        this.baseURL = '/api';
        this.token = localStorage.getItem('token');
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        };

        if (this.token) {
            config.headers.Authorization = `Bearer ${this.token}`;
        }

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // Convenience methods
    get(endpoint) {
        return this.request(endpoint);
    }

    post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    put(endpoint, data) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    delete(endpoint) {
        return this.request(endpoint, {
            method: 'DELETE',
        });
    }
}

// Usage
const api = new APIClient();
const users = await api.get('/users');
```

### WebSocket Integration

**WebSocket client pattern:**
```javascript
class WebSocketClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 1000;
    }

    connect() {
        try {
            this.ws = new WebSocket(this.url);
            this.setupEventListeners();
        } catch (error) {
            console.error('WebSocket connection error:', error);
            this.reconnect();
        }
    }

    setupEventListeners() {
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.reconnectAttempts = 0;
        };

        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleMessage(data);
        };

        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
            this.reconnect();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    handleMessage(data) {
        // Override in subclass or set callback
        console.log('Received:', data);
    }

    send(data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        }
    }

    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            setTimeout(() => {
                console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
                this.connect();
            }, this.reconnectInterval * this.reconnectAttempts);
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}
```

## 🗄️ Database Management

### MongoDB Operations

**Development database setup:**
```bash
# Connect to MongoDB
mongosh

# Create development database
use go_chat_dev

# Create user collection with validation
db.createCollection("users", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["username", "email", "password"],
         properties: {
            username: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            email: {
               bsonType: "string",
               pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
               description: "must be a valid email address"
            },
            password: {
               bsonType: "string",
               minLength: 6,
               description: "must be a string and at least 6 characters"
            }
         }
      }
   }
})

# Create indexes
db.users.createIndex({"email": 1}, {unique: true})
db.users.createIndex({"username": 1}, {unique: true})
```

**Database migration pattern:**
```go
// migrations/001_create_indexes.go
func CreateIndexes(db *mongo.Database) error {
    // Users collection indexes
    userIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }
    
    _, err := db.Collection("users").Indexes().CreateMany(context.Background(), userIndexes)
    return err
}
```

### Data Seeding

**Development data seeder:**
```go
// utils/seeder.go
func SeedDevelopmentData(db *mongo.Database) error {
    users := []interface{}{
        models.User{
            Username: "testuser1",
            Email:    "test1@example.com",
            Password: "hashedpassword1",
            CreatedAt: time.Now(),
        },
        models.User{
            Username: "testuser2",
            Email:    "test2@example.com",
            Password: "hashedpassword2",
            CreatedAt: time.Now(),
        },
    }
    
    _, err := db.Collection("users").InsertMany(context.Background(), users)
    return err
}
```

## 🧪 Testing

### Unit Testing

**Test file structure:**
```
backend/
├── auth/
│   ├── handlers.go
│   ├── handlers_test.go
│   ├── middleware.go
│   └── middleware_test.go
```

**Example test:**
```go
// auth/handlers_test.go
package auth

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/register", RegisterHandler)
    
    // Test cases
    tests := []struct {
        name           string
        body           map[string]string
        expectedStatus int
        expectedError  string
    }{
        {
            name: "Valid registration",
            body: map[string]string{
                "username": "testuser",
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedStatus: http.StatusCreated,
        },
        {
            name: "Missing username",
            body: map[string]string{
                "email":    "test@example.com",
                "password": "password123",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Username is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Prepare request
            jsonBody, _ := json.Marshal(tt.body)
            req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            
            // Execute request
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            // Assert
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            if tt.expectedError != "" {
                var response map[string]interface{}
                json.Unmarshal(w.Body.Bytes(), &response)
                assert.Contains(t, response["error"], tt.expectedError)
            }
        })
    }
}
```

**Running tests:**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./auth -run TestRegisterHandler

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Testing

**API integration test:**
```go
// tests/integration_test.go
func TestUserRegistrationFlow(t *testing.T) {
    // Setup test database
    testDB := setupTestDB()
    defer cleanupTestDB(testDB)
    
    // Setup test server
    router := setupTestRouter(testDB)
    server := httptest.NewServer(router)
    defer server.Close()
    
    // Test registration
    registerData := map[string]string{
        "username": "integrationtest",
        "email":    "integration@test.com",
        "password": "testpassword123",
    }
    
    registerResp := makeRequest(t, server.URL+"/api/auth/register", "POST", registerData)
    assert.Equal(t, http.StatusCreated, registerResp.StatusCode)
    
    // Test login
    loginData := map[string]string{
        "email":    "integration@test.com",
        "password": "testpassword123",
    }
    
    loginResp := makeRequest(t, server.URL+"/api/auth/login", "POST", loginData)
    assert.Equal(t, http.StatusOK, loginResp.StatusCode)
    
    // Extract token
    var loginResponse map[string]interface{}
    json.NewDecoder(loginResp.Body).Decode(&loginResponse)
    token := loginResponse["token"].(string)
    
    // Test protected endpoint
    protectedResp := makeAuthenticatedRequest(t, server.URL+"/api/auth/profile", "GET", nil, token)
    assert.Equal(t, http.StatusOK, protectedResp.StatusCode)
}
```

### Frontend Testing

**JavaScript unit tests:**
```javascript
// tests/auth.test.js
describe('Authentication', () => {
    let apiClient;
    
    beforeEach(() => {
        apiClient = new APIClient();
        localStorage.clear();
    });
    
    describe('login', () => {
        it('should store token on successful login', async () => {
            // Mock fetch
            global.fetch = jest.fn().mockResolvedValue({
                ok: true,
                json: () => Promise.resolve({
                    token: 'test-token',
                    user: { id: '1', username: 'testuser' }
                })
            });
            
            await apiClient.login('test@example.com', 'password');
            
            expect(localStorage.getItem('token')).toBe('test-token');
        });
        
        it('should throw error on invalid credentials', async () => {
            global.fetch = jest.fn().mockResolvedValue({
                ok: false,
                json: () => Promise.resolve({
                    error: 'Invalid credentials'
                })
            });
            
            await expect(apiClient.login('invalid@example.com', 'wrong'))
                .rejects.toThrow('Invalid credentials');
        });
    });
});
```

## ✅ Code Quality

### Linting

**golangci-lint configuration** (`.golangci.yml`):
```yaml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true

linters:
  disable-all: true
  enable:
    - bodyclose
    - deadcode
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - structcheck
    - typecheck
    - unused
    - varcheck
    - gofmt
    - goimports
    - golint
    - gosec

linters-settings:
  gosec:
    excludes:
      - G204  # Subprocess launched with variable
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/nazmusSakibRaiyan/Go_chat
```

**Run linting:**
```bash
# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix

# Run specific linter
golangci-lint run --enable=gosec
```

### Code Formatting

**Go formatting:**
```bash
# Format all Go files
gofmt -w .

# Format with imports
goimports -w .

# Check formatting without making changes
gofmt -d .
```

**JavaScript formatting (Prettier):**
```bash
# Install prettier
npm install -g prettier

# Format JavaScript files
prettier --write "frontend/public/js/**/*.js"

# Check formatting
prettier --check "frontend/public/js/**/*.js"
```

### Pre-commit Hooks

**Git hooks setup:**
```bash
#!/bin/sh
# .git/hooks/pre-commit

echo "Running pre-commit checks..."

# Go formatting
echo "Checking Go formatting..."
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
    echo "Some files are not properly formatted:"
    echo "$unformatted"
    exit 1
fi

# Go linting
echo "Running Go linter..."
golangci-lint run
if [ $? -ne 0 ]; then
    echo "Linting failed"
    exit 1
fi

# Go tests
echo "Running Go tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

echo "All checks passed!"
```

## 🐛 Debugging

### Backend Debugging

**Using VS Code Debugger:**
1. Set breakpoints in your Go code
2. Press F5 or use "Run and Debug"
3. Select "Launch Go Chat Backend"

**Using Delve (command line):**
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main.go
dlv debug main.go

# Debug with arguments
dlv debug main.go -- --port 8080

# Attach to running process
dlv attach <PID>
```

**Common debugging commands:**
```bash
(dlv) break main.main       # Set breakpoint
(dlv) continue              # Continue execution
(dlv) next                  # Step over
(dlv) step                  # Step into
(dlv) print variableName    # Print variable
(dlv) locals                # Show local variables
(dlv) goroutines            # List goroutines
```

### Frontend Debugging

**Browser Developer Tools:**
```javascript
// Console debugging
console.log('Debug info:', data);
console.error('Error occurred:', error);
console.table(arrayData);

// Breakpoint in code
debugger;

// Performance timing
console.time('API Call');
await api.getData();
console.timeEnd('API Call');
```

**Network debugging:**
1. Open Developer Tools (F12)
2. Go to Network tab
3. Monitor API calls
4. Check request/response headers
5. Verify payload data

### Database Debugging

**MongoDB debugging:**
```javascript
// Enable profiling
db.setProfilingLevel(2)

// View slow operations
db.system.profile.find().limit(5).sort({ts: -1}).pretty()

// Check current operations
db.currentOp()

// Explain query execution
db.users.find({email: "test@example.com"}).explain("executionStats")
```

### Logging

**Structured logging in Go:**
```go
// Use structured logging
import "github.com/sirupsen/logrus"

log := logrus.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "login",
})
log.Info("User login attempt")

// Error logging with context
log.WithError(err).Error("Database connection failed")
```

**Log levels:**
```go
log.Trace("Very detailed information")
log.Debug("Debug information")
log.Info("General information")
log.Warn("Warning information")
log.Error("Error information")
log.Fatal("Fatal error - exits program")
log.Panic("Panic - calls panic()")
```

## 🤝 Contributing

### Git Workflow

**Branch naming convention:**
```
feature/description        # New features
bugfix/description        # Bug fixes
hotfix/description        # Critical fixes
refactor/description      # Code refactoring
docs/description          # Documentation updates
```

**Commit message format:**
```
type(scope): description

feat(auth): add password strength validation
fix(chat): resolve WebSocket connection issue
docs(api): update authentication endpoints
refactor(db): optimize user query performance
test(auth): add integration tests for login
```

**Pull Request Process:**
1. Create feature branch from `main`
2. Make changes with descriptive commits
3. Add/update tests
4. Update documentation
5. Create pull request
6. Address review feedback
7. Merge after approval

### Code Review Checklist

**Backend Review:**
- [ ] Code follows Go conventions
- [ ] Error handling is proper
- [ ] Security considerations addressed
- [ ] Tests cover new functionality
- [ ] Database operations are optimized
- [ ] Logging is appropriate
- [ ] Documentation is updated

**Frontend Review:**
- [ ] Code follows JavaScript best practices
- [ ] UI is responsive and accessible
- [ ] Error handling is user-friendly
- [ ] Performance considerations addressed
- [ ] Browser compatibility maintained

### Setting up Development Branch

```bash
# Create and switch to feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "feat(feature): implement new feature"

# Push branch
git push origin feature/new-feature

# Create pull request on GitHub
```

## 📚 Best Practices

### Go Best Practices

**1. Error Handling**
```go
// Good
if err != nil {
    return fmt.Errorf("failed to connect to database: %w", err)
}

// Bad
if err != nil {
    log.Fatal(err)  // Don't use Fatal in libraries
}
```

**2. Context Usage**
```go
// Pass context through call chain
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.FindUser(ctx, id)
}

// Set timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

**3. Interface Design**
```go
// Keep interfaces small
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
}

// Accept interfaces, return structs
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### JavaScript Best Practices

**1. Async/Await**
```javascript
// Good
async function fetchUserData(id) {
    try {
        const response = await fetch(`/api/users/${id}`);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Failed to fetch user:', error);
        throw error;
    }
}

// Bad
function fetchUserData(id) {
    return fetch(`/api/users/${id}`)
        .then(response => response.json())
        .catch(error => console.log(error));
}
```

**2. Error Handling**
```javascript
// Good - specific error handling
try {
    const data = await api.getData();
    updateUI(data);
} catch (error) {
    if (error.name === 'NetworkError') {
        showNetworkErrorMessage();
    } else if (error.status === 401) {
        redirectToLogin();
    } else {
        showGenericErrorMessage();
    }
}
```

**3. DOM Manipulation**
```javascript
// Good - efficient DOM updates
const fragment = document.createDocumentFragment();
items.forEach(item => {
    const element = createItemElement(item);
    fragment.appendChild(element);
});
container.appendChild(fragment);

// Bad - multiple DOM updates
items.forEach(item => {
    container.appendChild(createItemElement(item));
});
```

### Security Best Practices

**1. Input Validation**
```go
// Validate all inputs
func validateUser(user *User) error {
    if len(user.Username) < 3 || len(user.Username) > 50 {
        return errors.New("username must be 3-50 characters")
    }
    if !isValidEmail(user.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

**2. Authentication**
```go
// Use secure password hashing
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// Validate JWT tokens properly
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, err
}
```

**3. CORS Configuration**
```go
// Configure CORS properly for production
config := cors.Config{
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

### Performance Best Practices

**1. Database Optimization**
```go
// Use indexes
db.users.createIndex({"email": 1}, {unique: true})

// Limit query results
cursor, err := collection.Find(ctx, filter, options.Find().SetLimit(100))

// Use projection to limit fields
projection := bson.M{"password": 0, "created_at": 0}
options := options.FindOne().SetProjection(projection)
```

**2. Caching**
```go
// Implement simple in-memory cache
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}
```

**3. Frontend Performance**
```javascript
// Debounce user input
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Use it for search inputs
const debouncedSearch = debounce(performSearch, 300);
searchInput.addEventListener('input', debouncedSearch);
```

---

**Happy Coding! 🚀**

For questions or issues, please refer to the [main documentation](README.md) or open an issue on GitHub.

---

**Last Updated**: July 30, 2025  
**Development Guide Version**: 1.0
