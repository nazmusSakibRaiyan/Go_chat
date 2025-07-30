# Go Chat Authentication System

## 🚀 Feature Overview

This is a complete authentication system for the Go Chat application with user registration, login, logout functionality, and secure JWT-based authentication.

## ✨ Features Implemented

### Backend Features
- **User Registration** with validation
- **User Login** with JWT token generation
- **User Logout** endpoint
- **Password Hashing** using bcrypt
- **JWT Authentication** with middleware
- **Rate Limiting** to prevent brute force attacks
- **Input Validation** for security
- **Password Strength Validation**
- **CORS Configuration** for frontend access

### Frontend Features
- **Interactive Login/Register Forms**
- **Real-time Form Validation**
- **Password Strength Indicator**
- **Loading States** with spinners
- **Error/Success Message Display**
- **Auto-login** with stored tokens
- **Responsive Design**
- **Modern UI** with animations

## 🛠 Technical Implementation

### Backend (Go)
- **Framework**: Gin Web Framework
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Rate Limiting**: Custom implementation
- **Validation**: Gin validation with custom rules

### Frontend (HTML/CSS/JavaScript)
- **Pure JavaScript** (no frameworks)
- **Modern CSS** with flexbox and animations
- **LocalStorage** for token persistence
- **Fetch API** for HTTP requests
- **Real-time validation** with visual feedback

## 📁 File Structure

```
backend/
├── auth/
│   ├── auth.go           # Password hashing & validation
│   ├── handlers.go       # Auth route handlers
│   ├── middleware.go     # JWT middleware
│   └── ratelimit.go      # Rate limiting
├── models/
│   └── models.go         # User model
├── db/
│   └── database.go       # MongoDB connection
└── main.go               # Server setup

frontend/public/
├── index-auth.html       # Full chat app with auth
├── auth-demo.html        # Authentication demo
└── index.html           # Basic chat (no auth)
```

## 🔐 Security Features

### Password Security
- Minimum 6 characters
- Must contain letters and numbers
- bcrypt hashing with cost 14

### Authentication Security
- JWT tokens with 24-hour expiration
- Bearer token authentication
- Secure token validation

### Rate Limiting
- 5 login attempts per 15 minutes per IP
- Automatic cleanup of old entries
- Prevents brute force attacks

### Input Validation
- Username: 3-20 characters, alphanumeric
- Email: Valid email format
- Password: Strength requirements
- SQL injection prevention

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- MongoDB
- Modern web browser

### Backend Setup
1. **Install Dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

2. **Set Environment Variables**
   ```bash
   export MONGO_URI="mongodb://localhost:27017"
   export MONGO_DB_NAME="go_chat"
   export JWT_SECRET="your-secret-key"
   export PORT="8080"
   ```

3. **Run the Server**
   ```bash
   go run main.go
   ```

### Frontend Usage
1. **Open the Authentication Demo**
   ```
   Open: frontend/public/auth-demo.html
   ```

2. **Or use the Full Chat App**
   ```
   Open: frontend/public/index-auth.html
   ```

## 🎯 API Endpoints

### Authentication Routes
```
POST /api/auth/register
POST /api/auth/login
POST /api/auth/logout
```

### Protected Routes
```
GET /api/me                    # Get current user info
```

### Example Usage

#### Register a New User
```javascript
POST /api/auth/register
Content-Type: application/json

{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123"
}
```

#### Login
```javascript
POST /api/auth/login
Content-Type: application/json

{
    "username": "johndoe",
    "password": "securepass123"
}
```

## 🎨 UI Features

### Login/Register Forms
- Clean, modern design
- Real-time validation feedback
- Loading states with spinners
- Error/success messages
- Password strength indicator
- Smooth animations

### Responsive Design
- Mobile-friendly
- Tablet optimized
- Desktop enhanced
- Cross-browser compatible

## 🔧 Configuration

### Environment Variables
```bash
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=go_chat
JWT_SECRET=your-super-secret-key
PORT=8080
```

### CORS Settings
Configured to allow requests from:
- http://localhost:3000
- http://localhost:5173
- http://127.0.0.1:5500

## 🐛 Error Handling

### Backend Errors
- Invalid request data
- User already exists
- Invalid credentials
- Database connection errors
- Token validation errors
- Rate limit exceeded

### Frontend Errors
- Network connectivity issues
- Invalid form data
- Authentication failures
- Token expiration
- Server unavailable

## 🧪 Testing

### Manual Testing
1. **Registration**: Create new accounts
2. **Login**: Sign in with credentials
3. **Token Persistence**: Refresh page, stay logged in
4. **Logout**: Clear session data
5. **Validation**: Test form validation
6. **Rate Limiting**: Test multiple failed logins

### Example Test Cases
- Valid registration with strong password
- Duplicate username/email handling
- Invalid credentials rejection
- Token expiration handling
- Password strength validation
- Rate limiting functionality

## 🔮 Future Enhancements

### Planned Features
- **Email Verification** during registration
- **Password Reset** functionality
- **Two-Factor Authentication** (2FA)
- **Social Login** (Google, GitHub)
- **Account Management** (profile updates)
- **Session Management** (multiple devices)
- **Audit Logging** for security events

### Technical Improvements
- **Redis** for session storage
- **Refresh Tokens** for better security
- **OAuth2** implementation
- **Password Policies** configuration
- **Account Lockout** after failed attempts
- **Email Notifications** for security events

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

If you encounter any issues:
1. Check the server logs
2. Verify MongoDB connection
3. Confirm environment variables
4. Check browser console for errors
5. Test API endpoints with Postman

---

**Ready to Chat!** 🎉 Your authentication system is now complete and secure!
