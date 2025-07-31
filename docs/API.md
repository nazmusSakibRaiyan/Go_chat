# 📡 API Documentation

Complete API reference for the Go Chat application backend services.

## 🌐 Base URL

```
http://localhost:8080/api
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 📋 Response Format

All API responses follow this standard format:

```json
{
  "success": boolean,
  "message": string,
  "token": string (optional),
  "user": object (optional),
  "data": object (optional)
}
```

## 🔑 Authentication Endpoints

### Register User

**POST** `/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "username": "string (3-20 chars, required)",
  "email": "string (valid email, required)",
  "password": "string (min 6 chars, required)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input data
- `409 Conflict` - Username or email already exists
- `500 Internal Server Error` - Server error

**Example:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### Login User

**POST** `/auth/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Invalid credentials
- `429 Too Many Requests` - Rate limit exceeded
- `503 Service Unavailable` - Database unavailable

**Example:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepass123"
  }'
```

### Logout User

**POST** `/auth/logout`

Logout user (mainly for client-side token cleanup).

**Headers:**
```
Authorization: Bearer <token> (optional)
```

**Response:**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 👤 User Endpoints

### Get Current User

**GET** `/me` 🔒

Get current authenticated user information.

**Headers:**
```
Authorization: Bearer <token> (required)
```

**Response:**
```json
{
  "success": true,
  "message": "User data retrieved",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Server error

**Example:**
```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Update User Profile

**PUT** `/profile` 🔒

Update current user's profile information.

**Headers:**
```
Authorization: Bearer <token> (required)
Content-Type: application/json
```

**Request Body:**
```json
{
  "display_name": "string (3-50 chars, optional)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T12:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid display name (too short, too long, or empty)
- `401 Unauthorized` - Invalid or missing token
- `500 Internal Server Error` - Server error

**Validation Rules:**
- Display name must be 3-50 characters long
- Display name cannot be only whitespace
- Special characters are allowed

**Example:**
```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"display_name": "John Doe"}'
```

## 💬 Chat Room Endpoints

### Get All Rooms

**GET** `/rooms`

Retrieve list of all available chat rooms.

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "General",
    "description": "General discussion room",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "name": "Tech Talk",
    "description": "Technology discussions",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  }
]
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/rooms
```

### Create Room

**POST** `/rooms`

Create a new chat room.

**Request Body:**
```json
{
  "name": "string (required)",
  "description": "string (optional)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Room created successfully",
  "room": {
    "id": "507f1f77bcf86cd799439013",
    "name": "New Room",
    "description": "A new chat room",
    "created_at": "2025-07-30T10:00:00Z",
    "updated_at": "2025-07-30T10:00:00Z"
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Room",
    "description": "A new chat room"
  }'
```

### Get Room Messages

**GET** `/rooms/{roomId}/messages`

Get message history for a specific room.

**Path Parameters:**
- `roomId` - Room ID (MongoDB ObjectID)

**Query Parameters:**
- `limit` - Number of messages to retrieve (default: 50)
- `offset` - Number of messages to skip (default: 0)

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439014",
    "room_id": "507f1f77bcf86cd799439011",
    "user_id": "507f1f77bcf86cd799439010",
    "username": "johndoe",
    "content": "Hello everyone!",
    "type": "message",
    "created_at": "2025-07-30T10:00:00Z"
  }
]
```

**Example:**
```bash
curl -X GET "http://localhost:8080/api/rooms/507f1f77bcf86cd799439011/messages?limit=20"
```

## 🔌 WebSocket Connection

### Connect to WebSocket

**GET** `/ws`

Establish WebSocket connection for real-time messaging.

**Query Parameters:**
- `username` - Username for the connection
- `room_id` - Initial room to join

**WebSocket URL:**
```
ws://localhost:8080/api/ws?username=johndoe&room_id=507f1f77bcf86cd799439011
```

**Message Format:**
```json
{
  "type": "message|user_joined|user_left",
  "room_id": "507f1f77bcf86cd799439011",
  "username": "johndoe",
  "content": "Message content",
  "timestamp": "2025-07-30T10:00:00Z"
}
```

**Example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws?username=johndoe&room_id=507f1f77bcf86cd799439011');

// Send message
ws.send(JSON.stringify({
  type: 'message',
  room_id: '507f1f77bcf86cd799439011',
  username: 'johndoe',
  content: 'Hello everyone!'
}));

// Receive messages
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

## 🏥 System Endpoints

### Health Check

**GET** `/health`

Check if the server is running and healthy.

**Response:**
```json
{
  "status": "ok"
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/health
```

## 🔒 Security

### Rate Limiting

Currently **disabled for development**. In production:
- **Login endpoint**: 10 attempts per 5 minutes per IP
- **Registration endpoint**: No limit (but could be added)

### JWT Token

- **Algorithm**: HS256
- **Expiration**: 24 hours
- **Claims**: user_id, username, issued_at, expires_at

### Password Security

- **Hashing**: bcrypt with cost 14
- **Minimum length**: 6 characters
- **Requirements**: Must contain letters and numbers

## ❌ Error Codes

### HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

### Common Error Messages

```json
{
  "success": false,
  "message": "Invalid request data: validation failed"
}
```

```json
{
  "success": false,
  "message": "Username already exists"
}
```

```json
{
  "success": false,
  "message": "Invalid username or password"
}
```

```json
{
  "success": false,
  "message": "Too many login attempts. Please try again later."
}
```

## 📝 Request Examples

### Complete Authentication Flow

```bash
# 1. Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# Response: {"success": true, "token": "...", "user": {...}}

# 2. Login user
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# Response: {"success": true, "token": "eyJ...", "user": {...}}

# 3. Access protected endpoint
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Response: {"success": true, "user": {...}}

# 4. Logout
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Response: {"success": true, "message": "Logout successful"}
```

## 🧪 Testing

Use the built-in API tester:
```
http://localhost:8000/api-tester.html
```

Or use tools like:
- **Postman** - GUI API client
- **curl** - Command line client
- **HTTPie** - User-friendly command line client
- **Insomnia** - REST client

## 📚 Additional Resources

- [Authentication Documentation](AUTHENTICATION.md)
- [Deployment Guide](DEPLOYMENT.md)
- [Main Documentation](README.md)

---

**Last Updated**: July 30, 2025  
**API Version**: 1.0  
**Backend Version**: Go 1.19+
