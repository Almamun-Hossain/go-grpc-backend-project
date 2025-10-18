package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// User represents a user in the system
type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Phone        string
	AvatarURL    string
	ReferralCode string
	ReferredBy   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Asset represents a user's asset
type Asset struct {
	ID        string
	UserID    string
	AssetType string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Referral represents a referral relationship
type Referral struct {
	ID               string
	ReferrerID       string
	ReferredUserID   string
	ReferredUsername string
	Earnings         float64
	Status           string
	CreatedAt        time.Time
}

// Withdrawal represents a withdrawal request
type Withdrawal struct {
	ID             string
	UserID         string
	Amount         float64
	PaymentMethod  string
	PaymentDetails string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MembershipRepo defines the interface for membership data operations
type MembershipRepo interface {
	// User operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)

	// Asset operations
	GetUserAssets(ctx context.Context, userID string) ([]*Asset, error)
	GetAssetByID(ctx context.Context, id string) (*Asset, error)
	UpdateAsset(ctx context.Context, asset *Asset) (*Asset, error)
	CreateAsset(ctx context.Context, asset *Asset) (*Asset, error)

	// Referral operations
	GetReferrals(ctx context.Context, userID string, page, pageSize int) ([]*Referral, int, error)
	CreateReferral(ctx context.Context, referral *Referral) (*Referral, error)

	// Withdrawal operations
	CreateWithdrawal(ctx context.Context, withdrawal *Withdrawal) (*Withdrawal, error)
	GetWithdrawals(ctx context.Context, userID string, page, pageSize int) ([]*Withdrawal, int, error)
	UpdateWithdrawal(ctx context.Context, withdrawal *Withdrawal) (*Withdrawal, error)
}

// MembershipUsecase handles membership business logic
type MembershipUsecase struct {
	repo   MembershipRepo
	log    *log.Helper
	config *AuthConfig
}

// NewMembershipUsecase creates a new MembershipUsecase
func NewMembershipUsecase(repo MembershipRepo, logger log.Logger) *MembershipUsecase {
	return &MembershipUsecase{
		repo:   repo,
		log:    log.NewHelper(logger),
		config: DefaultAuthConfig(),
	}
}

// Register creates a new user account
func (uc *MembershipUsecase) Register(ctx context.Context, email, password, username, referralCode string) (*User, string, string, error) {
	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		uc.log.Errorf("Failed to hash password: %v", err)
		return nil, "", "", err
	}

	user := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		ReferralCode: generateReferralCode(),
		ReferredBy:   referralCode,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	// Generate JWT tokens
	accessToken, err := GenerateAccessToken(createdUser.ID, createdUser.Email, createdUser.Username, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate access token: %v", err)
		return nil, "", "", err
	}

	refreshToken, err := GenerateRefreshToken(createdUser.ID, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate refresh token: %v", err)
		return nil, "", "", err
	}

	return createdUser, accessToken, refreshToken, nil
}

// Login authenticates a user
func (uc *MembershipUsecase) Login(ctx context.Context, email, password string) (*User, string, string, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.Errorf("Failed to get user by email: %v", err)
		return nil, "", "", err
	}

	// Verify password hash
	err = VerifyPassword(user.PasswordHash, password)
	if err != nil {
		uc.log.Errorf("Invalid password for user %s: %v", email, err)
		return nil, "", "", err
	}

	// Generate JWT tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate access token: %v", err)
		return nil, "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user.ID, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate refresh token: %v", err)
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// GetProfile retrieves user profile
func (uc *MembershipUsecase) GetProfile(ctx context.Context, userID string) (*User, error) {
	return uc.repo.GetUserByID(ctx, userID)
}

// UpdateProfile updates user profile
func (uc *MembershipUsecase) UpdateProfile(ctx context.Context, userID, username, phone, avatarURL string) (*User, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Phone = phone
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	return uc.repo.UpdateUser(ctx, user)
}

// GetAssets retrieves user assets
func (uc *MembershipUsecase) GetAssets(ctx context.Context, userID string) ([]*Asset, float64, error) {
	assets, err := uc.repo.GetUserAssets(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	var totalBalance float64
	for _, asset := range assets {
		totalBalance += asset.Balance
	}

	return assets, totalBalance, nil
}

// UpdateAsset updates an asset balance
func (uc *MembershipUsecase) UpdateAsset(ctx context.Context, assetID string, amount float64, operation string) (*Asset, error) {
	asset, err := uc.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	switch operation {
	case "add":
		asset.Balance += amount
	case "subtract":
		asset.Balance -= amount
	}

	asset.UpdatedAt = time.Now()

	return uc.repo.UpdateAsset(ctx, asset)
}

// GetReferrals retrieves user referrals
func (uc *MembershipUsecase) GetReferrals(ctx context.Context, userID string, page, pageSize int) ([]*Referral, int, float64, error) {
	referrals, total, err := uc.repo.GetReferrals(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}

	var totalEarnings float64
	for _, ref := range referrals {
		totalEarnings += ref.Earnings
	}

	return referrals, total, totalEarnings, nil
}

// RequestWithdrawal creates a withdrawal request
func (uc *MembershipUsecase) RequestWithdrawal(ctx context.Context, userID string, amount float64, paymentMethod, paymentDetails string) (*Withdrawal, error) {
	withdrawal := &Withdrawal{
		UserID:         userID,
		Amount:         amount,
		PaymentMethod:  paymentMethod,
		PaymentDetails: paymentDetails,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return uc.repo.CreateWithdrawal(ctx, withdrawal)
}

// GetWithdrawals retrieves user withdrawals
func (uc *MembershipUsecase) GetWithdrawals(ctx context.Context, userID string, page, pageSize int) ([]*Withdrawal, int, error) {
	return uc.repo.GetWithdrawals(ctx, userID, page, pageSize)
}

// RefreshToken generates new access and refresh tokens using a valid refresh token
func (uc *MembershipUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := ValidateRefreshToken(refreshToken, uc.config)
	if err != nil {
		uc.log.Errorf("Invalid refresh token: %v", err)
		return "", "", err
	}

	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		uc.log.Errorf("Failed to get user by ID: %v", err)
		return "", "", err
	}

	// Generate new tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Username, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate access token: %v", err)
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(user.ID, uc.config)
	if err != nil {
		uc.log.Errorf("Failed to generate refresh token: %v", err)
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// ValidateToken validates an access token and returns the user claims
func (uc *MembershipUsecase) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	claims, err := ValidateAccessToken(tokenString, uc.config)
	if err != nil {
		uc.log.Errorf("Token validation failed: %v", err)
		return nil, err
	}

	return claims, nil
}

// Helper function to generate referral code
func generateReferralCode() string {
	// TODO: Implement proper referral code generation
	return "REF" + time.Now().Format("20060102150405")
}
