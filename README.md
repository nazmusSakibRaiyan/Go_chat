# Go_chat
A lightweight real-time chat server built with Go. Supports multiple clients, basic messaging, and concurrent handling using goroutines and channels.

## Features

- Real-time messaging using WebSockets
- Multiple chat rooms
- User management
- Clean, responsive web interface
- Concurrent message handling with Go routines
- RESTful API endpoints
- Optional database integration (PostgreSQL)

## Architecture

```
backend/
├── main.go           # Application entry point
├── chat/             # WebSocket and chat logic
├── config/           # Configuration management
├── db/               # Database operations
├── models/           # Data structures
└── utils/            # Utility functions

frontend/
└── public/
    └── index.html    # Simple HTML chat interface
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Optional: PostgreSQL database

### Running the Application

1. **Start the Backend Server:**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

2. **Open the Frontend:**
   Open `frontend/public/index.html` in your web browser or use VS Code's Live Server extension.

3. **Start Chatting:**
   - Enter a username
   - Click "Connect"
   - Choose a chat room
   - Start messaging!

### Testing with Multiple Users

Open the `index.html` file in multiple browser tabs or windows with different usernames to simulate multiple users chatting.

## API Endpoints

- `GET /health` - Health check
- `GET /api/ws` - WebSocket connection endpoint
- `GET /api/rooms` - List all chat rooms
- `POST /api/rooms` - Create a new chat room
- `GET /api/rooms/:id/messages` - Get messages for a specific room

## WebSocket Message Format

### Client to Server:
```json
{
  "type": "chat_message",
  "content": "Hello, world!",
  "room_id": 1
}
```

### Server to Client:
```json
{
  "type": "chat_message",
  "username": "john_doe",
  "content": "Hello, world!",
  "room_id": 1
}
```

## Configuration

Environment variables (`.env` file):
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string (optional)
- `JWT_SECRET` - Secret key for JWT tokens

## Database Setup (Optional)

If you want to enable message persistence:

1. Install PostgreSQL
2. Create a database named `chatdb`
3. Update the `DATABASE_URL` in `.env` file
4. The application will automatically create the required tables

## Development

### Adding New Features

1. **New Message Types:** Add handlers in `chat/chat.go`
2. **Database Models:** Update `models/models.go`
3. **API Endpoints:** Add routes in `main.go`
4. **Frontend Features:** Modify `frontend/public/index.html`

### Code Structure

- **Hub Pattern:** Manages WebSocket connections and message broadcasting
- **Goroutines:** Each WebSocket connection runs in separate goroutines
- **Channels:** Used for communication between goroutines
- **Gin Framework:** Provides HTTP routing and middleware

## Technology Stack

**Backend:**
- Go 1.21
- Gin (HTTP framework)
- Gorilla WebSocket
- PostgreSQL (optional)

**Frontend:**
- HTML5
- CSS3
- Vanilla JavaScript
- WebSocket API

## License

This project is open source and available under the MIT License.
