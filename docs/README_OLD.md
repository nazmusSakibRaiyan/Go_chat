# Go Chat Application with MongoDB

A real-time chat application built with Go, WebSockets, and MongoDB. Features include real-time messaging, multiple chat rooms, persistent message history, and a clean web interface.

## 🚀 Features

- **Real-time messaging** using WebSockets
- **MongoDB integration** for persistent data storage
- **Multiple chat rooms** with dynamic loading
- **Message history** - see previous messages when joining rooms
- **User management** with unique usernames
- **Responsive web interface** - works on desktop and mobile
- **Concurrent handling** using Go routines and channels
- **RESTful API** for room and message management
- **MongoDB Atlas support** - works with cloud databases

## 🏗️ Project Structure

```
Go_chat/
├── README.md                 # Project documentation
├── MONGODB_SETUP.md         # MongoDB setup guide
├── .gitignore              # Git ignore rules
├── setup.ps1               # Development setup script
│
├── backend/                # Go backend server
│   ├── main.go            # Application entry point
│   ├── go.mod             # Go module dependencies
│   ├── go.sum             # Dependency checksums
│   ├── .env.example       # Environment variables template
│   ├── build.ps1          # Build script (PowerShell)
│   ├── build.sh           # Build script (Bash)
│   │
│   ├── auth/              # Authentication utilities
│   │   └── auth.go        # Password hashing and validation
│   │
│   ├── chat/              # WebSocket and chat logic
│   │   └── chat.go        # Hub pattern and message handling
│   │
│   ├── config/            # Configuration management
│   │   └── config.go      # Environment variables and settings
│   │
│   ├── db/                # Database operations
│   │   └── database.go    # MongoDB connection and CRUD
│   │
│   ├── models/            # Data structures
│   │   └── models.go      # MongoDB models with BSON tags
│   │
│   └── utils/             # Helper functions
│       └── utils.go       # String sanitization and utilities
│
└── frontend/              # Web interface
    ├── README.md          # Frontend documentation
    ├── package.json       # Development dependencies
    └── public/
        └── index.html     # Complete chat application
```

## 🛠️ Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **MongoDB** - Choose one option:
  - **MongoDB Atlas** (Cloud) - [Sign up free](https://www.mongodb.com/atlas)
  - **Local MongoDB** - [Download MongoDB](https://www.mongodb.com/try/download/community)
  - **Docker** - `docker run -d -p 27017:27017 mongo:latest`

## ⚡ Quick Start

### 1. Setup (Automated)

**Windows:**
```powershell
.\setup.ps1
```

**Manual Setup:**
```bash
# Clone repository
git clone <your-repo-url>
cd Go_chat

# Setup backend
cd backend
go mod tidy
cp .env.example .env
# Edit .env with your MongoDB connection details

# Optional: Setup frontend
cd ../frontend
npm install
```

### 2. Configure Database

Edit `backend/.env` with your MongoDB connection:

```env
PORT=8080
MONGODB_URI=mongodb+srv://your_username:your_password@your_cluster.mongodb.net/
MONGODB_DATABASE=go_chat_db
JWT_SECRET=your-secret-key-change-this-in-production
```

⚠️ **Important:** Never commit your actual `.env` file to git. The `.env.example` file is a template showing what variables are needed.

**For MongoDB Atlas:**
- Use your Atlas connection string
- Make sure to replace `username`, `password`, and `cluster` with your actual values

**For Local MongoDB:**
```env
MONGODB_URI=mongodb://localhost:27017
```

### 3. Run the Application

```bash
# Start the backend server
cd backend
go run main.go

# Or build and run
.\build.ps1
.\chat-server.exe

# Open frontend in browser
# Navigate to: frontend/public/index.html
```

## 🎯 Usage

### Starting a Chat Session

1. **Open the frontend** in your web browser
2. **Enter a username** (3-20 characters)
3. **Click "Connect"** to join the chat
4. **Choose a room** from the available options
5. **Start chatting!** Messages are saved to MongoDB

### Multiple Users

- Open the application in **multiple browser tabs** with different usernames
- Messages are delivered in **real-time** to all connected users
- **Message history** is loaded when joining rooms

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/api/ws` | WebSocket connection |
| `GET` | `/api/rooms` | List all chat rooms |
| `POST` | `/api/rooms` | Create a new room |
| `GET` | `/api/rooms/:id/messages` | Get room message history |

### WebSocket Message Format

**Client → Server:**
```json
{
  "type": "chat_message",
  "content": "Hello, world!",
  "room_id": "ObjectId"
}
```

**Server → Client:**
```json
{
  "type": "chat_message",
  "username": "john_doe",
  "content": "Hello, world!",
  "room_id": "ObjectId"
}
```

## 🗄️ Database Schema

### Collections

**Rooms Collection:**
```javascript
{
  "_id": ObjectId,
  "name": "general",
  "description": "General chat room",
  "created_at": ISODate,
  "updated_at": ISODate
}
```

**Messages Collection:**
```javascript
{
  "_id": ObjectId,
  "room_id": ObjectId,
  "user_id": ObjectId, // optional
  "username": "john_doe",
  "content": "Hello world!",
  "type": "text",
  "created_at": ISODate
}
```

**Users Collection:**
```javascript
{
  "_id": ObjectId,
  "username": "john_doe",
  "email": "john@example.com",
  "password_hash": "bcrypt_hash",
  "created_at": ISODate,
  "updated_at": ISODate
}
```

## 🔧 Development

### Project Structure

- **Hub Pattern**: Manages WebSocket connections and message broadcasting
- **Goroutines**: Each WebSocket connection runs in separate goroutines
- **Channels**: Communication between goroutines for thread safety
- **MongoDB Driver**: Official Go driver for MongoDB operations
- **Gin Framework**: HTTP routing and middleware

### Adding Features

1. **New Message Types**: Add handlers in `chat/chat.go`
2. **Database Models**: Update `models/models.go`
3. **API Endpoints**: Add routes in `main.go`
4. **Frontend Features**: Modify `frontend/public/index.html`

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | `go_chat_db` |
| `JWT_SECRET` | Secret for JWT tokens | `your-secret-key` |

## 🚀 Deployment

### Production Checklist

- [ ] Change `JWT_SECRET` to a secure random value
- [ ] Set `GIN_MODE=release` for production
- [ ] Use MongoDB Atlas or secure MongoDB instance
- [ ] Enable HTTPS/TLS
- [ ] Set up proper CORS policies
- [ ] Configure firewall rules
- [ ] Set up monitoring and logging

### Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o chat-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/chat-server .
COPY --from=builder /app/.env .
CMD ["./chat-server"]
```

## 🛠️ Troubleshooting

### Common Issues

**Connection Failed:**
- Check MongoDB connection string in `.env`
- Verify MongoDB Atlas IP whitelist settings
- Ensure MongoDB service is running (for local installations)

**WebSocket Connection Failed:**
- Check if port 8080 is available
- Verify CORS settings allow your frontend domain
- Check browser developer console for errors

**Messages Not Persisting:**
- Verify MongoDB connection is successful
- Check server logs for database errors
- Ensure proper database permissions

### Logging

The application logs important events:
- MongoDB connection status
- WebSocket connections/disconnections
- Message broadcasting
- API requests

## 🔐 Security Features

- **Password hashing** using bcrypt
- **Input sanitization** for usernames and messages
- **CORS protection** for API endpoints
- **Connection validation** for WebSocket upgrades

## 📈 Performance

- **Automatic indexing** for optimized queries
- **Connection pooling** via MongoDB driver
- **Efficient message broadcasting** using Go channels
- **Concurrent request handling** with Gin framework

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

## 🎉 Success!

Your Go Chat Application with MongoDB is now ready! 

- **Backend**: Connected to MongoDB Atlas ✅
- **Frontend**: Dynamic room loading ✅
- **Real-time**: WebSocket messaging ✅
- **Persistence**: Message history in MongoDB ✅

Start chatting and enjoy your new real-time chat application! 🚀
