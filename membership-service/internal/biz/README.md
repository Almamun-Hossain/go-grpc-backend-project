# Authentication Implementation

This package now includes secure password hashing and JWT token management for the membership service.

## Features Implemented

### Password Security
- **bcrypt hashing**: Passwords are securely hashed using bcrypt with default cost
- **Password verification**: Secure password comparison during login

### JWT Token Management
- **Access tokens**: Short-lived tokens (15 minutes) for API access
- **Refresh tokens**: Long-lived tokens (7 days) for token renewal
- **Token validation**: Secure token parsing and validation
- **Token refresh**: Ability to generate new tokens using refresh tokens

## Configuration

The authentication system uses `AuthConfig` for JWT settings:

```go
type AuthConfig struct {
    AccessTokenSecret  string        // Secret key for access tokens
    RefreshTokenSecret string        // Secret key for refresh tokens
    AccessTokenExpiry  time.Duration // Access token lifetime (default: 15 minutes)
    RefreshTokenExpiry time.Duration // Refresh token lifetime (default: 7 days)
}
```

## Usage Examples

### User Registration
```go
user, accessToken, refreshToken, err := membershipUsecase.Register(
    ctx, 
    "user@example.com", 
    "password123", 
    "username", 
    "REF123456" // optional referral code
)
```

### User Login
```go
user, accessToken, refreshToken, err := membershipUsecase.Login(
    ctx, 
    "user@example.com", 
    "password123"
)
```

### Token Validation
```go
claims, err := membershipUsecase.ValidateToken(ctx, accessToken)
if err != nil {
    // Handle invalid token
}
// Use claims.UserID, claims.Email, claims.Username
```

### Token Refresh
```go
newAccessToken, newRefreshToken, err := membershipUsecase.RefreshToken(ctx, refreshToken)
```

## Security Notes

1. **Change default secrets**: Update the default JWT secrets in production
2. **Environment variables**: Consider loading secrets from environment variables
3. **Token storage**: Store refresh tokens securely (database, Redis, etc.)
4. **HTTPS only**: Always use HTTPS in production for token transmission

## Dependencies Added

- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT token handling