#!/bin/bash

echo "🚀 Building Go Backend Microservices with Latest Go Version"
echo "=========================================================="

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
cd ..

# Build Website Service  
echo "Building website-service..."
cd website-service
go build -o bin/website-service ./cmd/website-service
cd ..

# Build Banking Service
echo "Building banking-service..."
cd banking-service
go build -o bin/banking-service ./cmd/banking-service
cd ..

# Build API Gateway
echo "Building api-gateway..."
cd api-gateway
go build -o bin/api-gateway ./cmd/main.go
cd ..

echo "✅ All services built successfully!"
echo ""

# Create simple Docker images using local builds
echo "🐳 Creating Docker images..."

# Create a simple Dockerfile for all services
cat > Dockerfile.local << 'EOF'
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY bin/ /app/
COPY configs/ /app/configs/

EXPOSE 8080 8081 8082 8083 9081 9082 9083

CMD ["./membership-service", "-conf", "./configs/config.yaml"]
EOF

# Build membership service image
echo "Building membership-service Docker image..."
docker build -f Dockerfile.local -t membership-service:latest ./membership-service/

# Build website service image  
echo "Building website-service Docker image..."
docker build -f Dockerfile.local -t website-service:latest ./website-service/

# Build banking service image
echo "Building banking-service Docker image..."
docker build -f Dockerfile.local -t banking-service:latest ./banking-service/

# Build API gateway image
echo "Building api-gateway Docker image..."
docker build -f Dockerfile.local -t api-gateway:latest ./api-gateway/

echo "✅ All Docker images created successfully!"
echo ""

# Start services with docker-compose
echo "🚀 Starting all services..."
docker compose up -d

echo ""
echo "🎉 All services are now running!"
echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "🌐 API Endpoints:"
echo "- API Gateway: http://localhost:8080"
echo "- Membership Service: http://localhost:8081"  
echo "- Website Service: http://localhost:8082"
echo "- Banking Service: http://localhost:8083"
echo ""
echo "🔍 To view logs: docker compose logs -f [service-name]"
echo "🛑 To stop: docker compose down"
