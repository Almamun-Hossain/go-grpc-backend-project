#!/bin/bash

# Go Backend Microservices API Gateway Test Script
# This script tests all the API endpoints through the API Gateway

BASE_URL="http://localhost:8080"
echo "🚀 Testing Go Backend Microservices API Gateway"
echo "================================================"

# Test Health Check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/health" | jq . || echo "Health check failed"
echo ""

# Test Authentication
echo "2. Testing User Registration..."
curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","username":"testuser"}' | jq . || echo "Registration failed"
echo ""

echo "3. Testing User Login..."
curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq . || echo "Login failed"
echo ""

# Test Website Service
echo "4. Testing Website Service - Get Site Settings..."
curl -s "$BASE_URL/api/v1/site/settings" | jq . || echo "Site settings failed"
echo ""

echo "5. Testing Website Service - Get Logo..."
curl -s "$BASE_URL/api/v1/site/logo" | jq . || echo "Logo failed"
echo ""

echo "6. Testing Website Service - Get Dynamic Settings..."
curl -s "$BASE_URL/api/v1/site/dynamic-settings" | jq . || echo "Dynamic settings failed"
echo ""

# Test Banking Service
echo "7. Testing Banking Service - List Cards..."
curl -s "$BASE_URL/api/v1/banking/cards" \
  -H "Authorization: Bearer test_token" | jq . || echo "List cards failed"
echo ""

echo "8. Testing Banking Service - Get Exchange Rate..."
curl -s "$BASE_URL/api/v1/banking/exchange-rate?from_currency=USD&to_currency=EUR" \
  -H "Authorization: Bearer test_token" | jq . || echo "Exchange rate failed"
echo ""

# Test User Management (Protected Routes)
echo "9. Testing User Profile..."
curl -s "$BASE_URL/api/v1/users/profile" \
  -H "Authorization: Bearer test_token" | jq . || echo "User profile failed"
echo ""

echo "10. Testing User Assets..."
curl -s "$BASE_URL/api/v1/users/assets" \
  -H "Authorization: Bearer test_token" | jq . || echo "User assets failed"
echo ""

echo "✅ API Gateway testing completed!"
echo "================================================"
