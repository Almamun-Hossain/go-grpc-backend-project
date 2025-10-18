# Go Backend Microservices Project - Local Development Guide

A comprehensive guide for running the Go Backend Microservices project locally and with Docker.

## 📋 Prerequisites

### For Local Development:
- **Go 1.21 or higher**
- **Protocol Buffers compiler (protoc)**
- **PostgreSQL 15** (if running without Docker)
- **Redis** (optional, for caching)
- **Make** (for build automation)

### For Docker:
- **Docker 20.10+**
- **Docker Compose 2.0+**

## 🚀 Method 1: Running with Docker (Recommended)

### Quick Start
```bash
# 1. Clone and navigate to the project
cd go-backend-project

# 2. Run everything with Docker Compose
docker-compose up --build
```

### Detailed Docker Setup

1. **Start all services:**
   ```bash
   docker-compose up --build -d
   ```

2. **View logs:**
   ```bash
   # All services
   docker-compose logs -f
   
   # Specific service
   docker-compose logs -f membership-service
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

4. **Stop and remove volumes:**
   ```bash
   docker-compose down -v
   ```

### Docker Services Overview
- **API Gateway**: `http://localhost:8080`
- **Membership Service**: `http://localhost:8081` (HTTP), `localhost:9081` (gRPC)
- **Website Service**: `http://localhost:8082` (HTTP), `localhost:9082` (gRPC)
- **Banking Service**: `http://localhost:8083` (HTTP), `localhost:9083` (gRPC)
- **PostgreSQL Databases**: 
  - Membership DB: `localhost:5432`
  - Website DB: `localhost:5433`
  - Banking DB: `localhost:5434`

## 🛠️ Method 2: Running Locally (Development)

### Step 1: Install Dependencies

```bash
# Install Go tools for protocol buffer generation
cd membership-service
make init

cd ../website-service
make init

cd ../banking-service
make init
```

### Step 2: Generate Protocol Buffers

```bash
# Generate protobuf code for all services
cd membership-service
make api

cd ../website-service
make api

cd ../banking-service
make api
```

### Step 3: Set Up Databases

#### Option A: Using Docker for Databases Only
```bash
# Start only the databases
docker-compose up -d membership-db website-db banking-db
```

#### Option B: Local PostgreSQL Installation
```bash
# Install PostgreSQL 15 locally
# Create databases
createdb membership
createdb website
createdb banking
```

### Step 4: Run Services Individually

Open **4 separate terminal windows**:

#### Terminal 1 - Membership Service:
```bash
cd membership-service
go run ./cmd/membership-service
```

#### Terminal 2 - Website Service:
```bash
cd website-service
go run ./cmd/website-service
```

#### Terminal 3 - Banking Service:
```bash
cd banking-service
go run ./cmd/banking-service
```

#### Terminal 4 - API Gateway:
```bash
cd api-gateway
go run ./cmd/main.go
```

## 🧪 Testing the Setup

### Health Check
```bash
curl http://localhost:8080/health
```

### Run Test Script
```bash
# Make the test script executable
chmod +x test-api.sh

# Run the test script
./test-api.sh
```

### Manual API Testing

#### 1. User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser"
  }'
```

#### 2. User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 3. Website Settings
```bash
curl http://localhost:8080/api/v1/site/settings
```

#### 4. Banking Exchange Rate
```bash
curl "http://localhost:8080/api/v1/banking/exchange-rate?from_currency=USD&to_currency=EUR" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## ⚙️ Configuration

### Environment Variables

#### Membership Service
```bash
export DATABASE_URL="postgres://membership_user:membership_pass@localhost:5432/membership?sslmode=disable"
export HTTP_PORT="8081"
export GRPC_PORT="9081"
```

#### Website Service
```bash
export DATABASE_URL="postgres://website_user:website_pass@localhost:5433/website?sslmode=disable"
export HTTP_PORT="8082"
export GRPC_PORT="9082"
```

#### Banking Service
```bash
export DATABASE_URL="postgres://banking_user:banking_pass@localhost:5434/banking?sslmode=disable"
export UPSTREAM_BANK_URL="https://api.bank.example.com"
export UPSTREAM_BANK_API_KEY="your_api_key_here"
export HTTP_PORT="8083"
export GRPC_PORT="9083"
```

#### API Gateway
```bash
export PORT=":8080"
export MEMBERSHIP_SERVICE_URL="http://localhost:8081"
export WEBSITE_SERVICE_URL="http://localhost:8082"
export BANKING_SERVICE_URL="http://localhost:8083"
```

## 🔧 Development Commands

### Build Services
```bash
# Build all services
cd membership-service && make build
cd ../website-service && make build
cd ../banking-service && make build
cd ../api-gateway && go build -o bin/api-gateway ./cmd/main.go
```

### Run Tests
```bash
# Test all services
cd membership-service && go test ./...
cd ../website-service && go test ./...
cd ../banking-service && go test ./...
cd ../api-gateway && go test ./...
```

### Regenerate Code
```bash
# Regenerate protobuf code
cd membership-service && make all
cd ../website-service && make all
cd ../banking-service && make all
```

## 🐛 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :8080
   
   # Kill the process
   kill -9 <PID>
   ```

2. **Database Connection Issues**
   ```bash
   # Check if PostgreSQL is running
   docker-compose ps
   
   # Check database logs
   docker-compose logs membership-db
   ```

3. **Protocol Buffer Generation Issues**
   ```bash
   # Install missing tools
   make init
   
   # Regenerate all code
   make all
   ```

4. **Service Not Starting**
   ```bash
   # Check service logs
   docker-compose logs -f <service-name>
   
   # Check if all dependencies are running
   docker-compose ps
   ```

### Reset Everything
```bash
# Stop and remove all containers and volumes
docker-compose down -v

# Remove all images
docker-compose down --rmi all

# Start fresh
docker-compose up --build
```

## 📊 Service Status

| Service | HTTP Port | gRPC Port | Database Port | Status |
|---------|-----------|-----------|---------------|--------|
| API Gateway | 8080 | - | - | ✅ Ready |
| Membership | 8081 | 9081 | 5432 | ✅ Ready |
| Website | 8082 | 9082 | 5433 | ✅ Ready |
| Banking | 8083 | 9083 | 5434 | ✅ Ready |

## 🎯 Development Workflow

### 1. Making Changes
```bash
# For local development
# 1. Make your changes
# 2. Regenerate protobuf if needed
make api

# 3. Restart the specific service
# Ctrl+C to stop, then restart
go run ./cmd/service-name
```

### 2. Testing Changes
```bash
# Run the test script
./test-api.sh

# Or test specific endpoints
curl http://localhost:8080/health
```

### 3. Docker Development
```bash
# Rebuild and restart specific service
docker-compose up --build membership-service

# View logs while developing
docker-compose logs -f membership-service
```

## 📝 Notes

- **Database**: Currently using mock data - database implementation is TODO
- **Authentication**: JWT validation is placeholder - needs implementation
- **File Uploads**: Logo uploads return placeholder URLs
- **External APIs**: Banking service uses mock upstream API responses

## 🆘 Getting Help

If you encounter issues:

1. Check the logs: `docker-compose logs -f`
2. Verify all services are running: `docker-compose ps`
3. Check port availability: `lsof -i :PORT`
4. Reset everything: `docker-compose down -v && docker-compose up --build`

For more detailed information, see the main README.md file.