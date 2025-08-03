# 🎭 Avatar Feature Documentation

## Overview
The avatar feature allows users to select and customize their profile with one of 12 predefined avatars. Users can change their avatar at any time through the profile update functionality. The system uses emoji representations for fast loading and consistent display across all devices.

## Features
- **12 Unique Avatars**: Cat 🐱, Dog 🐶, Bear 🐻, Fox 🦊, Lion 🦁, Panda 🐼, Robot 🤖, Alien 👽, Ninja 🥷, Pirate 🏴‍☠️, Knight ⚔️, and Wizard 🧙‍♂️
- **Real-time Updates**: Avatar changes are immediately reflected in the user profile and chat interface
- **Emoji-based System**: Fast loading using emoji characters instead of image files
- **Default Avatar**: New users get a default avatar (👤) that can be changed anytime
- **Validation**: Only valid avatar IDs from the predefined list are accepted
- **Chat Integration**: Avatars display beside usernames in the main chat interface

## Available Avatars

| Avatar ID | Emoji | Name | Description |
|-----------|-------|------|-------------|
| `cat` | 🐱 | Cat | Cute cat avatar |
| `dog` | 🐶 | Dog | Friendly dog avatar |
| `bear` | 🐻 | Bear | Cuddly bear avatar |
| `fox` | 🦊 | Fox | Clever fox avatar |
| `lion` | 🦁 | Lion | Majestic lion avatar |
| `panda` | 🐼 | Panda | Adorable panda avatar |
| `robot` | 🤖 | Robot | Futuristic robot avatar |
| `alien` | 👽 | Alien | Mysterious alien avatar |
| `ninja` | 🥷 | Ninja | Stealthy ninja avatar |
| `pirate` | 🏴‍☠️ | Pirate | Adventurous pirate avatar |
| `knight` | ⚔️ | Knight | Noble knight avatar |
| `wizard` | 🧙‍♂️ | Wizard | Magical wizard avatar |

## API Endpoints

### Get Available Avatars
```
GET /api/avatars
```
Returns a list of all available avatars with their IDs, names, emojis, and descriptions.

**Response:**
```json
{
  "success": true,
  "message": "Available avatars retrieved",
  "avatars": [
    {
      "id": "cat",
      "name": "Cat",
      "emoji": "🐱",
      "description": "Cute cat avatar"
    },
    {
      "id": "dog", 
      "name": "Dog",
      "emoji": "🐶",
      "description": "Friendly dog avatar"
    }
    // ... all 12 avatars
  ]
}
```

### Update User Profile (including avatar)
```
PUT /api/profile
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "display_name": "My Display Name",
  "avatar": "cat"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "user": {
    "id": "...",
    "username": "john_doe",
    "email": "john@example.com",
    "display_name": "My Display Name",
    "avatar": "cat",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

## Database Schema Changes

### User Model
The `User` model has been updated to include an `avatar` field:

```go
type User struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username     string             `json:"username" bson:"username"`
    Email        string             `json:"email" bson:"email"`
    PasswordHash string             `json:"-" bson:"password_hash"`
    DisplayName  string             `json:"display_name" bson:"display_name"`
    Avatar       string             `json:"avatar" bson:"avatar"` // Avatar ID (string)
    CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
```

### Avatar Model
```go
type Avatar struct {
    ID          string `json:"id" bson:"id"`
    Name        string `json:"name" bson:"name"`
    Emoji       string `json:"emoji" bson:"emoji"`
    Description string `json:"description" bson:"description"`
}
```

## User Interface

### Profile & Avatar Manager
The avatar selection is available through the dedicated Profile & Avatar Manager interface:
- Access via `profile-demo.html` or navigation links
- Grid layout displaying all 12 avatars with emoji representations
- Visual selection with hover effects and active states
- Real-time preview of selected avatar
- Integrated with profile update functionality

### Main Chat Interface
- User's selected avatar displays beside their username in the chat header
- Consistent avatar representation across the application
- Emoji-based display for fast loading and universal compatibility

## Usage Examples

### Frontend Integration
```javascript
// Avatar emoji mapping for frontend display
const avatarEmojis = {
    'cat': '🐱',
    'dog': '🐶',
    'bear': '🐻',
    'fox': '🦊',
    'lion': '🦁',
    'panda': '🐼',
    'robot': '🤖',
    'alien': '👽',
    'ninja': '🥷',
    'pirate': '🏴‍☠️',
    'knight': '⚔️',
    'wizard': '🧙‍♂️'
};

// Get available avatars
const response = await fetch('/api/avatars');
const data = await response.json();
const avatars = data.avatars;

// Update user avatar
const updateResponse = await fetch('/api/profile', {
    method: 'PUT',
    headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
        display_name: "John Doe",
        avatar: "bear"  // Bear avatar
    })
});

// Display avatar in UI
function displayUserAvatar(user) {
    const avatarEmoji = avatarEmojis[user.avatar] || '👤';
    document.getElementById('userAvatar').textContent = avatarEmoji;
}
```

### Testing & Usage
You can test the avatar feature through multiple interfaces:
1. **Profile & Avatar Manager** (`profile-demo.html`) - Dedicated avatar selection interface
2. **Main Chat Interface** (`index.html`) - See avatars in action during chat
3. **Authentication Demo** (`index-auth.html`) - Access profile management after login

## Recent Improvements (Version 2.2.0)
- ✅ **Emoji-based System**: Switched from image files to emoji for better performance
- ✅ **Chat Integration**: Avatars now display in the main chat interface
- ✅ **Improved UI**: Better avatar selection interface with visual feedback
- ✅ **String-based IDs**: Changed from numeric IDs to descriptive string IDs
- ✅ **Enhanced Validation**: Better error handling and validation
- ✅ **Real-time Updates**: Immediate avatar updates across all interfaces
3. Select different avatars from the grid
4. Update your profile to save the changes

## Avatar Files
Avatar images are stored in `/frontend/public/static/avatars/` and are served as static files by the backend server. Each avatar is an SVG file that can be easily customized or replaced.

## Implementation Details

### Backend Changes
1. **Models**: Added `Avatar` struct and helper functions in `models/models.go`
2. **Database**: Updated `UpdateUserProfile` function to handle avatar updates
3. **Auth Handlers**: Added avatar validation and `GetAvatars` endpoint
4. **Routes**: Added `/api/avatars` endpoint and static file serving

### Validation
- Avatar IDs must be between 1 and 12
- Invalid avatar IDs are rejected with a 400 Bad Request response
- The `IsValidAvatar()` function ensures data integrity

### Default Behavior
- New users automatically receive avatar ID 1 (Cat) during registration
- If no avatar is provided during profile update, the current avatar is preserved
- The system gracefully handles missing or invalid avatar data

This feature provides a fun and engaging way for users to personalize their chat experience while maintaining data consistency and security.
