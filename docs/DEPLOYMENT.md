# 🚀 Deployment Guide

Complete guide for deploying the Go Chat application to production environments.

## 📋 Table of Contents

- [Prerequisites](#-prerequisites)
- [Environment Setup](#-environment-setup)
- [Local Development](#-local-development)
- [Production Deployment](#-production-deployment)
- [Docker Deployment](#-docker-deployment)
- [Cloud Deployment](#-cloud-deployment)
- [Security Considerations](#-security-considerations)
- [Monitoring & Maintenance](#-monitoring--maintenance)
- [Troubleshooting](#-troubleshooting)

## 🔧 Prerequisites

### System Requirements

**Minimum:**
- **CPU**: 1 core
- **RAM**: 512MB
- **Storage**: 1GB
- **Network**: Port 8080 (backend), Port 80/443 (frontend)

**Recommended:**
- **CPU**: 2+ cores
- **RAM**: 2GB+
- **Storage**: 10GB+
- **Network**: Load balancer, CDN

### Software Dependencies

- **Go**: 1.19 or higher
- **MongoDB**: 4.4 or higher
- **Node.js**: 16+ (for development tools)
- **Git**: For version control
- **Web Server**: Nginx, Apache, or similar (for frontend)

## 🌍 Environment Setup

### Environment Variables

Create a `.env` file or set system environment variables:

```bash
# Database Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=go_chat_prod

# Security
JWT_SECRET=your-super-secure-256-bit-secret-key-here
BCRYPT_COST=14

# Server Configuration
PORT=8080
GIN_MODE=release

# CORS Configuration (Production)
ALLOWED_ORIGINS=https://yourchatapp.com,https://www.yourchatapp.com

# Rate Limiting (Production)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=5
RATE_LIMIT_WINDOW=15m

# Logging
LOG_LEVEL=info
LOG_FILE=/var/log/gochat/app.log
```

### Security Environment Variables

```bash
# Generate secure JWT secret (256-bit)
openssl rand -hex 32

# Example secure values
JWT_SECRET=a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456
DB_PASSWORD=SecureP@ssw0rd123!
```

## 💻 Local Development

### Quick Start

1. **Clone Repository**
   ```bash
   git clone https://github.com/nazmusSakibRaiyan/Go_chat.git
   cd Go_chat
   ```

2. **Backend Setup**
   ```bash
   cd backend
   
   # Install dependencies
   go mod tidy
   
   # Set environment variables
   export MONGO_URI="mongodb://localhost:27017"
   export MONGO_DB_NAME="go_chat_dev"
   export JWT_SECRET="dev-secret-key"
   export PORT="8080"
   
   # Run in development mode
   go run main.go
   ```

3. **Frontend Setup**
   ```bash
   cd frontend/public
   
   # Start development server
   python3 -m http.server 8000
   # OR
   npx serve -p 8000
   ```

4. **Access Application**
   - Frontend: http://localhost:8000
   - Backend API: http://localhost:8080/api
   - Health Check: http://localhost:8080/health

### Development Tools

```bash
# Install Go development tools
go install github.com/cosmtrek/air@latest  # Live reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting

# Use Air for live reload
air
```

## 🏭 Production Deployment

### 1. Server Preparation

**Ubuntu/Debian:**
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y git curl wget

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install MongoDB
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt update
sudo apt install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 2. Application Deployment

```bash
# Create application user
sudo useradd -m -s /bin/bash gochat
sudo su - gochat

# Clone and build application
git clone https://github.com/nazmusSakibRaiyan/Go_chat.git
cd Go_chat/backend

# Build production binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gochat-server main.go

# Set permissions
chmod +x gochat-server
```

### 3. Systemd Service

Create `/etc/systemd/system/gochat.service`:

```ini
[Unit]
Description=Go Chat Server
After=network.target mongod.service
Wants=mongod.service

[Service]
Type=simple
User=gochat
Group=gochat
WorkingDirectory=/home/gochat/Go_chat/backend
ExecStart=/home/gochat/Go_chat/backend/gochat-server
Restart=always
RestartSec=5

# Environment variables
Environment=MONGO_URI=mongodb://localhost:27017
Environment=MONGO_DB_NAME=go_chat_prod
Environment=JWT_SECRET=your-production-secret-here
Environment=PORT=8080
Environment=GIN_MODE=release

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/gochat

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gochat

[Install]
WantedBy=multi-user.target
```

**Enable and start service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable gochat
sudo systemctl start gochat
sudo systemctl status gochat
```

### 4. Nginx Configuration

Create `/etc/nginx/sites-available/gochat`:

```nginx
upstream gochat_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name yourchatapp.com www.yourchatapp.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourchatapp.com www.yourchatapp.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/yourchatapp.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourchatapp.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Frontend static files
    location / {
        root /home/gochat/Go_chat/frontend/public;
        index index.html;
        try_files $uri $uri/ /index.html;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Backend API
    location /api/ {
        proxy_pass http://gochat_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 86400;
    }

    # WebSocket support
    location /api/ws {
        proxy_pass http://gochat_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }

    # Health check
    location /health {
        proxy_pass http://gochat_backend;
        access_log off;
    }
}
```

**Enable site:**
```bash
sudo ln -s /etc/nginx/sites-available/gochat /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 5. SSL Certificate

Using Let's Encrypt:
```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d yourchatapp.com -d www.yourchatapp.com
sudo systemctl reload nginx
```

## 🐳 Docker Deployment

### 1. Dockerfile (Backend)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
```

### 2. Docker Compose

```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:6.0
    container_name: gochat-mongo
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_ROOT_PASSWORD}
      MONGO_INITDB_DATABASE: go_chat
    volumes:
      - mongodb_data:/data/db
      - ./init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro
    networks:
      - gochat-network
    ports:
      - "27017:27017"

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gochat-backend
    restart: unless-stopped
    environment:
      MONGO_URI: mongodb://gochat_user:${MONGO_PASSWORD}@mongodb:27017/go_chat
      MONGO_DB_NAME: go_chat
      JWT_SECRET: ${JWT_SECRET}
      PORT: 8080
      GIN_MODE: release
    depends_on:
      - mongodb
    networks:
      - gochat-network
    ports:
      - "8080:8080"

  nginx:
    image: nginx:alpine
    container_name: gochat-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./frontend/public:/usr/share/nginx/html:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - backend
    networks:
      - gochat-network

volumes:
  mongodb_data:

networks:
  gochat-network:
    driver: bridge
```

### 3. Environment File

Create `.env`:
```bash
MONGO_ROOT_PASSWORD=secure_root_password
MONGO_PASSWORD=secure_user_password
JWT_SECRET=your-production-jwt-secret-256-bit
```

### 4. Deploy with Docker

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Scale backend service
docker-compose up -d --scale backend=3

# Update application
docker-compose pull
docker-compose up -d --force-recreate
```

## ☁️ Cloud Deployment

### AWS Deployment

**1. EC2 Instance:**
```bash
# Launch EC2 instance
aws ec2 run-instances \
  --image-id ami-0c02fb55956c7d316 \
  --instance-type t3.micro \
  --key-name your-key-pair \
  --security-group-ids sg-xxxxxxxxx \
  --subnet-id subnet-xxxxxxxxx
```

**2. RDS MongoDB:**
```bash
# Create DocumentDB cluster
aws docdb create-db-cluster \
  --db-cluster-identifier gochat-cluster \
  --engine docdb \
  --master-username admin \
  --master-user-password SecurePassword123
```

**3. Load Balancer:**
```bash
# Create Application Load Balancer
aws elbv2 create-load-balancer \
  --name gochat-alb \
  --subnets subnet-xxxxxxxx subnet-yyyyyyyy \
  --security-groups sg-xxxxxxxxx
```

### Google Cloud Platform

```bash
# Create GKE cluster
gcloud container clusters create gochat-cluster \
  --num-nodes=3 \
  --zone=us-central1-a

# Deploy to GKE
kubectl apply -f k8s/
```

### Heroku Deployment

```bash
# Install Heroku CLI
npm install -g heroku

# Login and create app
heroku login
heroku create your-gochat-app

# Set environment variables
heroku config:set JWT_SECRET=your-secret
heroku config:set MONGO_URI=your-mongo-uri

# Deploy
git push heroku main
```

## 🔒 Security Considerations

### 1. Production Security Checklist

- [ ] **Environment Variables**: Never commit secrets to git
- [ ] **JWT Secret**: Use 256-bit random key
- [ ] **HTTPS**: Always use SSL/TLS in production
- [ ] **CORS**: Restrict to specific domains
- [ ] **Rate Limiting**: Enable for all public endpoints
- [ ] **Input Validation**: Validate all user inputs
- [ ] **Database Security**: Use authentication and encryption
- [ ] **Firewall**: Only open necessary ports
- [ ] **Updates**: Keep dependencies updated
- [ ] **Monitoring**: Set up security monitoring

### 2. Security Headers

```go
// Add to main.go
func securityHeaders() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    })
}
```

### 3. Rate Limiting (Production)

```go
// Update rate limiter for production
rateLimiter := NewRateLimiter(5, 15*time.Minute)
```

## 📊 Monitoring & Maintenance

### 1. Health Checks

```bash
# Backend health check
curl -f http://localhost:8080/health || exit 1

# Database health check
mongosh --eval "db.adminCommand('ping')"
```

### 2. Logging

**Application Logs:**
```bash
# View systemd logs
sudo journalctl -u gochat -f

# Application log file
tail -f /var/log/gochat/app.log
```

**Nginx Logs:**
```bash
# Access logs
tail -f /var/log/nginx/access.log

# Error logs
tail -f /var/log/nginx/error.log
```

### 3. Monitoring Tools

**Basic Monitoring:**
```bash
# CPU and Memory usage
htop

# Disk usage
df -h

# Network connections
ss -tulpn | grep :8080
```

**Advanced Monitoring:**
- **Prometheus + Grafana**: Metrics and dashboards
- **ELK Stack**: Centralized logging
- **Uptime Robot**: External monitoring
- **New Relic**: Application performance monitoring

### 4. Backup Strategy

**Database Backup:**
```bash
# Create backup
mongodump --db go_chat_prod --out /backup/$(date +%Y%m%d)

# Automated backup script
#!/bin/bash
BACKUP_DIR="/backup/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR
mongodump --db go_chat_prod --out $BACKUP_DIR
tar -czf $BACKUP_DIR.tar.gz $BACKUP_DIR
rm -rf $BACKUP_DIR

# Upload to S3 (optional)
aws s3 cp $BACKUP_DIR.tar.gz s3://your-backup-bucket/
```

**Application Backup:**
```bash
# Backup application files
tar -czf gochat-app-$(date +%Y%m%d).tar.gz /home/gochat/Go_chat
```

## 🔧 Troubleshooting

### Common Issues

**1. Port Already in Use**
```bash
# Find process using port
sudo lsof -i :8080
# OR
sudo netstat -tulpn | grep :8080

# Kill process
sudo kill -9 <PID>
```

**2. MongoDB Connection Issues**
```bash
# Check MongoDB status
sudo systemctl status mongod

# Check MongoDB logs
sudo tail -f /var/log/mongodb/mongod.log

# Test connection
mongosh --eval "db.adminCommand('ping')"
```

**3. Nginx Configuration Issues**
```bash
# Test configuration
sudo nginx -t

# Reload configuration
sudo systemctl reload nginx

# Check error logs
sudo tail -f /var/log/nginx/error.log
```

**4. SSL Certificate Issues**
```bash
# Check certificate expiry
openssl x509 -in /etc/letsencrypt/live/yourdomain.com/cert.pem -text -noout | grep "Not After"

# Renew certificate
sudo certbot renew --dry-run
```

### Performance Troubleshooting

**1. High CPU Usage**
```bash
# Monitor Go application
top -p $(pgrep gochat-server)

# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile
```

**2. High Memory Usage**
```bash
# Check memory usage
free -h

# Monitor Go memory
go tool pprof http://localhost:8080/debug/pprof/heap
```

**3. Database Performance**
```bash
# MongoDB performance monitoring
mongosh --eval "db.currentOp()"
mongosh --eval "db.serverStatus()"
```

### Log Analysis

**Search for errors:**
```bash
# Backend errors
sudo journalctl -u gochat | grep -i error

# Nginx errors
sudo grep -i error /var/log/nginx/error.log

# MongoDB errors
sudo grep -i error /var/log/mongodb/mongod.log
```

## 📞 Support

For deployment issues:

1. Check this deployment guide
2. Review application logs
3. Check system resources
4. Verify network connectivity
5. Consult the main [documentation](README.md)

---

**Last Updated**: July 30, 2025  
**Deployment Guide Version**: 1.0
