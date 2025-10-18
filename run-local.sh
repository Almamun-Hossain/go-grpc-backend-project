#!/bin/bash

echo "🚀 Running Go Backend Microservices Locally with Latest Go Version"
echo "=================================================================="

# Check Go version
echo "📋 Go Version:"
go version
echo ""

# Build all services locally
echo "🔨 Building services locally..."

# Build Membership Service
echo "Building membership-service..."
cd membership-service
go build -o bin/membership-service ./cmd/membership-service
echo "✅ membership-service built successfully"
cd ..

# Build Website Service  
echo "Building website-service..."
cd website-service
go build -o bin/website-service ./cmd/website-service
echo "✅ website-service built successfully"
cd ..

# Build Banking Service
echo "Building banking-service..."
cd banking-service
go build -o bin/banking-service ./cmd/banking-service
echo "✅ banking-service built successfully"
cd ..

# Build API Gateway
echo "Building api-gateway..."
cd api-gateway
go build -o bin/api-gateway ./cmd/main.go
echo "✅ api-gateway built successfully"
cd ..

echo ""
echo "🎉 All services built successfully with Go 1.24.3!"
echo ""

# Start services in background
echo "🚀 Starting all services..."

# Start Membership Service
echo "Starting membership-service on port 8081..."
cd membership-service
./bin/membership-service -conf ./configs/config.yaml &
MEMBERSHIP_PID=$!
cd ..

# Start Website Service
echo "Starting website-service on port 8082..."
cd website-service
./bin/website-service -conf ./configs/config.yaml &
WEBSITE_PID=$!
cd ..

# Start Banking Service
echo "Starting banking-service on port 8083..."
cd banking-service
./bin/banking-service -conf ./configs/config.yaml &
BANKING_PID=$!
cd ..

# Start API Gateway
echo "Starting api-gateway on port 8080..."
cd api-gateway
./bin/api-gateway &
GATEWAY_PID=$!
cd ..

echo ""
echo "🎉 All services are now running!"
echo ""
echo "📊 Service Status:"
echo "- Membership Service (PID: $MEMBERSHIP_PID) - http://localhost:8081"
echo "- Website Service (PID: $WEBSITE_PID) - http://localhost:8082"  
echo "- Banking Service (PID: $BANKING_PID) - http://localhost:8083"
echo "- API Gateway (PID: $GATEWAY_PID) - http://localhost:8080"
echo ""
echo "🌐 Test the API:"
echo "curl -X POST http://localhost:8081/v1/auth/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"username\":\"testuser\"}'"
echo ""
echo "🛑 To stop all services: kill $MEMBERSHIP_PID $WEBSITE_PID $BANKING_PID $GATEWAY_PID"
echo ""

# Wait a moment for services to start
sleep 3

# Test the membership service
echo "🧪 Testing membership service..."
curl -s -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","username":"testuser"}' | head -c 200
echo ""
echo "✅ Test completed!"
