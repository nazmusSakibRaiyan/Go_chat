# Go Chat Application

A complete real-time chat application built with Go (backend) and vanilla JavaScript (frontend), featuring secure user authentication, WebSocket communication, and MongoDB for data persistence.

## 🚀 Features

- **🔐 User Authentication**: Secure registration and login with JWT tokens
- **💬 Real-time Chat**: WebSocket-based instant messaging
- **🛡️ Password Security**: Bcrypt hashing with strength validation
- **⚡ Rate Limiting**: Protection against brute force attacks
- **📱 Responsive UI**: Modern, mobile-friendly interface
- **🗄️ MongoDB Integration**: Efficient data storage and retrieval
- **🔍 API Testing**: Built-in API testing interface
- **📚 Comprehensive Documentation**: Complete guides for development and deployment

## 🛠️ Tech Stack

**Backend:**
- Go 1.19+ with Gin Web Framework
- MongoDB with official Go driver
- JWT authentication with 24-hour expiration
- bcrypt password hashing (cost 14)
- CORS middleware for cross-origin support
- Rate limiting for security

**Frontend:**
- Vanilla JavaScript (ES6+) with Fetch API
- WebSocket API for real-time communication
- Modern CSS with animations and transitions
- Responsive design for all devices
- localStorage for session persistence
- Multiple UI interfaces (chat, auth demo, API tester)

## 📚 Documentation

Comprehensive documentation is available in the [`docs/`](docs/) folder:

- **[📖 Main Documentation](docs/README.md)** - Complete project overview and setup guide
- **[🔗 API Documentation](docs/API.md)** - Detailed API endpoints and usage
- **[🔐 Authentication Guide](docs/AUTHENTICATION.md)** - Security implementation details
- **[🚀 Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment instructions
- **[🛠️ Development Guide](docs/DEVELOPMENT.md)** - Development setup and contributing

## 📁 Project Structure

```
Go_chat/
├── README.md                 # This file
├── setup.ps1                # Quick setup script
├── docker-compose.yml       # Docker deployment
├── Dockerfile              # Container build
│
├── backend/                # Go backend application
│   ├── main.go            # Application entry point
│   ├── go.mod             # Go module definition
│   ├── go.sum             # Dependency checksums
│   │
│   ├── auth/              # Authentication system
│   │   ├── auth.go        # JWT utilities
│   │   ├── handlers.go    # Auth HTTP handlers
│   │   ├── middleware.go  # Auth middleware
│   │   └── ratelimit.go   # Rate limiting
│   │
│   ├── chat/              # Chat functionality
│   │   └── chat.go        # WebSocket chat handlers
│   │
│   ├── config/            # Configuration management
│   │   └── config.go      # Application configuration
│   │
│   ├── db/                # Database layer
│   │   └── database.go    # MongoDB connection
│   │
│   ├── models/            # Data models
│   │   └── models.go      # User and message structures
│   │
│   └── utils/             # Utility functions
│       └── utils.go       # Helper functions
│
├── frontend/              # Frontend application
│   ├── package.json       # Frontend dependencies
│   ├── README.md          # Frontend documentation
│   └── public/            # Static web files
│       ├── index.html     # Main chat interface
│       ├── index-auth.html # Full authentication demo
│       ├── api-tester.html # API testing interface
│       └── launcher.html   # Frontend navigation
│
└── docs/                  # Documentation
    ├── README.md          # Main project documentation
    ├── API.md             # API reference
    ├── AUTHENTICATION.md  # Authentication guide
    ├── DEPLOYMENT.md      # Deployment instructions
    └── DEVELOPMENT.md     # Development guide
```

## ⚡ Quick Start

### Prerequisites

- **Go 1.19+** - [Download Go](https://golang.org/dl/)
- **MongoDB** - [MongoDB Atlas](https://www.mongodb.com/atlas) or local installation
- **Git** - [Download Git](https://git-scm.com/downloads)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/nazmusSakibRaiyan/Go_chat.git
   cd Go_chat
   ```

2. **Setup backend**
   ```bash
   cd backend
   go mod tidy
   
   # Set environment variables (or create .env file)
   export MONGO_URI="mongodb://localhost:27017"
   export MONGO_DB_NAME="go_chat"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   ```

3. **Run the application**
   ```bash
   # Start backend (from backend/ directory)
   go run main.go
   
   # Start frontend server (from frontend/public/ directory)
   python3 -m http.server 8000
   # OR
   npx serve -p 8000
   ```

4. **Access the application**
   - Frontend: http://localhost:8000
   - Backend API: http://localhost:8080/api
   - Health Check: http://localhost:8080/health

### Using VS Code Task

If you're using VS Code, you can use the built-in task:

```bash
# Start backend server
Ctrl+Shift+P → "Tasks: Run Task" → "Start Backend Server"
```

## 🎮 Usage

### User Registration & Login

1. **Navigate to the application** at http://localhost:8000
2. **Register a new account** with username, email, and password
3. **Login** with your credentials
4. **Start chatting** in real-time!

### Multiple Interfaces

- **Main Chat App** (`index.html`) - Full featured chat application
- **Auth Demo** (`index-auth.html`) - Simple authentication demonstration
- **API Tester** (`api-tester.html`) - Test API endpoints directly
- **Launcher** (`launcher.html`) - Navigate between interfaces

### API Testing

The built-in API tester allows you to:
- Test user registration and login
- View user profiles
- Test protected endpoints
- Monitor API responses in real-time

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt with cost factor 14
- **Rate Limiting**: Prevents brute force attacks
- **Input Validation**: Server-side validation for all inputs
- **CORS Configuration**: Secure cross-origin resource sharing
- **Secure Headers**: Additional security headers in responses

## 🐳 Docker Deployment

```bash
# Using Docker Compose
docker-compose up -d

# Manual Docker build
docker build -t go-chat .
docker run -p 8080:8080 go-chat
```

## 🧪 Testing

```bash
# Run backend tests
cd backend
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test -v ./auth
```

## 🤝 Contributing

We welcome contributions! Please see our [Development Guide](docs/DEVELOPMENT.md) for details on:

- Setting up the development environment
- Code style and conventions
- Testing requirements
- Pull request process

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For help and support:

1. **Check the documentation** in the [`docs/`](docs/) folder
2. **Review existing issues** on GitHub
3. **Create a new issue** if you can't find a solution
4. **Join our community discussions**

## 🗂️ Key Documentation Links

- **[Complete Setup Guide](docs/README.md)** - Detailed installation and configuration
- **[API Reference](docs/API.md)** - All endpoints with examples
- **[Authentication Guide](docs/AUTHENTICATION.md)** - Security implementation details
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment instructions
- **[Development Guide](docs/DEVELOPMENT.md)** - Contributing and development workflow

---

**Built with ❤️ using Go and JavaScript**

**Last Updated**: July 30, 2025  
**Version**: 2.0.0
