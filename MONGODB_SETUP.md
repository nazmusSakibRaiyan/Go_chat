# MongoDB Setup for Go Chat Application

## Option 1: Install MongoDB locally

### Windows
1. Download MongoDB Community Server from https://www.mongodb.com/try/download/community
2. Install MongoDB with default settings
3. Start MongoDB service:
   ```
   net start MongoDB
   ```

### macOS (using Homebrew)
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb/brew/mongodb-community
```

### Linux (Ubuntu/Debian)
```bash
# Import the public key
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -

# Create list file
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list

# Update and install
sudo apt-get update
sudo apt-get install -y mongodb-org

# Start service
sudo systemctl start mongod
sudo systemctl enable mongod
```

## Option 2: Use MongoDB Atlas (Cloud)

1. Go to https://www.mongodb.com/atlas
2. Create a free account
3. Create a new cluster
4. Get your connection string
5. Update `.env` file with your connection string:
   ```
   MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/
   ```

## Option 3: Use Docker

```bash
# Run MongoDB in a Docker container
docker run -d --name mongodb -p 27017:27017 mongo:latest

# Or with persistent data
docker run -d --name mongodb -p 27017:27017 -v mongodb_data:/data/db mongo:latest
```

## Verify Connection

After setting up MongoDB, you can verify the connection by running:

```bash
# Test connection with MongoDB shell
mongosh

# Or test with our Go application
go run main.go
```

## Default Configuration

The application will create:
- Database: `go_chat_db`
- Collections: `rooms`, `users`, `messages`
- Default rooms: general, random, tech

## Environment Variables

Update your `.env` file:
```
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=go_chat_db
JWT_SECRET=your-secret-key-change-this-in-production
```

## Troubleshooting

### Connection Issues
- Make sure MongoDB is running on port 27017
- Check firewall settings
- Verify MONGODB_URI in .env file

### Permission Issues
- Ensure the user has read/write permissions to the database
- For MongoDB Atlas, check IP whitelist settings

### Performance
- The application creates indexes automatically for better performance
- Monitor query performance in production
