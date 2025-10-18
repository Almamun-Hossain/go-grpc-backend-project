# Go gRPC Backend Project

A comprehensive microservices backend project built with Go, featuring gRPC communication, JWT authentication, and modern architecture patterns.

## 🚀 Features

- **Microservices Architecture**: Modular services with clear separation of concerns
- **gRPC Communication**: High-performance inter-service communication
- **JWT Authentication**: Secure token-based authentication with bcrypt password hashing
- **API Gateway**: Centralized request routing and management
- **Database Integration**: PostgreSQL support for data persistence
- **Docker Support**: Containerized deployment with Docker Compose
- **Latest Go Version**: Built with Go 1.24.3

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │  Membership     │    │   Website       │
│   (Port 8080)   │◄──►│  Service        │◄──►│   Service       │
│                 │    │  (Port 8081)    │    │   (Port 8082)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Banking       │    │  Membership     │    │   Website       │
│   Service       │    │  Database       │    │   Database      │
│   (Port 8083)   │    │  (Port 5432)    │    │   (Port 5433)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐
│   Banking       │
│   Database      │
│   (Port 5434)   │
└─────────────────┘
```

## 📦 Services

### 1. API Gateway (Port 8080)
- Centralized request routing
- Authentication middleware
- Rate limiting and monitoring
- Health check endpoints

### 2. Membership Service (Port 8081)
- User registration and authentication
- JWT token management
- Password hashing with bcrypt
- User profile management
- Referral system

### 3. Website Service (Port 8082)
- Site settings management
- Logo and branding
- Dynamic configuration
- Content management

### 4. Banking Service (Port 8083)
- Credit card management
- Remittance processing
- Exchange rate services
- Transaction history

## 🔧 Quick Start

### Prerequisites
- Go 1.24.3 or later
- Docker and Docker Compose (optional)
- PostgreSQL (for production)

### Local Development

1. **Clone the repository**
   ```bash
   git clone git@github.com:Almamun-Hossain/go-grpc-backend-project.git
   cd go-grpc-backend-project
   ```

2. **Run all services locally**
   ```bash
   ./run-local.sh
   ```

3. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # User registration
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123","username":"testuser"}'
   ```

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker compose up --build -d
   ```

2. **View service status**
   ```bash
   docker compose ps
   ```

## 🔐 Authentication

The project implements secure JWT-based authentication:

- **Password Hashing**: bcrypt with default cost
- **Access Tokens**: Short-lived (15 minutes)
- **Refresh Tokens**: Long-lived (7 days)
- **Token Validation**: Secure parsing and validation

### Authentication Endpoints

```bash
# Register new user
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password123",
  "username": "username"
}

# User login
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Refresh token
POST /api/v1/auth/refresh
{
  "refreshToken": "your_refresh_token"
}
```

## 📊 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh

### User Management
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `GET /api/v1/users/assets` - Get user assets
- `GET /api/v1/users/referrals` - Get user referrals

### Website Management
- `GET /api/v1/site/settings` - Get site settings
- `PUT /api/v1/site/settings` - Update site settings
- `GET /api/v1/site/logo` - Get site logo

### Banking Services
- `GET /api/v1/banking/cards` - List credit cards
- `POST /api/v1/banking/remittance` - Initiate remittance
- `GET /api/v1/banking/exchange-rate` - Get exchange rates

## 🛠️ Development

### Project Structure
```
go-grpc-backend-project/
├── api-gateway/          # API Gateway service
├── membership-service/   # Membership service
├── website-service/      # Website service
├── banking-service/      # Banking service
├── docker-compose.yml    # Docker Compose configuration
├── run-local.sh         # Local development script
└── README.md            # This file
```

### Building Services
```bash
# Build individual service
cd membership-service
go build -o bin/membership-service ./cmd/membership-service

# Build all services
./run-local.sh
```

### Testing
```bash
# Test authentication
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","username":"testuser"}'

# Test via API Gateway
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","username":"testuser"}'
```

## 🐳 Docker

### Docker Compose Services
- **API Gateway**: `api-gateway:latest`
- **Membership Service**: `membership-service:latest`
- **Website Service**: `website-service:latest`
- **Banking Service**: `banking-service:latest`
- **PostgreSQL Databases**: `postgres:15-alpine`

### Docker Commands
```bash
# Build and start all services
docker compose up --build -d

# View logs
docker compose logs -f [service-name]

# Stop all services
docker compose down

# Rebuild specific service
docker compose up --build -d [service-name]
```

## 🔧 Configuration

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing secret
- `PORT`: Service port number

### Configuration Files
Each service has its own `configs/config.yaml` file for local configuration.

## 📈 Performance

- **Startup Time**: ~3 seconds for all services
- **Memory Usage**: Minimal (Go's efficient runtime)
- **Response Time**: < 100ms for most API calls
- **Concurrency**: Handles multiple requests efficiently

## 🚀 Deployment

### Production Checklist
- [ ] Set up PostgreSQL databases
- [ ] Configure environment variables
- [ ] Set up SSL certificates
- [ ] Configure load balancer
- [ ] Set up monitoring and logging
- [ ] Configure backup strategy

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Almamun Hossain**
- Website: [almamun.me](https://almamun.me)
- GitHub: [@Almamun-Hossain](https://github.com/Almamun-Hossain)

## 🙏 Acknowledgments

- Go team for the excellent language and tooling
- gRPC team for the high-performance RPC framework
- All open-source contributors who made this project possible

---

**Built with ❤️ using Go 1.24.3**