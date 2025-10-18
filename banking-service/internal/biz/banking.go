package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Card represents a bank card
type Card struct {
	ID               string
	UserID           string
	CardNumberMasked string
	CardType         string
	CardBrand        string
	Status           string
	ExpiryDate       time.Time
	CreatedAt        time.Time
}

// CardBalance represents card balance information
type CardBalance struct {
	CardID             string
	AvailableBalance   float64
	CreditLimit        float64
	OutstandingBalance float64
	Currency           string
	LastUpdated        time.Time
}

// PhysicalCardDetails represents physical card details
type PhysicalCardDetails struct {
	CardID           string
	CardNumberMasked string
	CardholderName   string
	ExpiryDate       string
	CVVMasked        string
	CardType         string
	CardBrand        string
}

// Remittance represents a remittance transaction
type Remittance struct {
	TransactionID   string
	FromAccount     string
	ToAccount       string
	Amount          float64
	FromCurrency    string
	ToCurrency      string
	ExchangeRate    float64
	Fee             float64
	Status          string
	BeneficiaryName string
	BeneficiaryBank string
	Purpose         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ExchangeRate represents currency exchange rate
type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	Timestamp    time.Time
}

// BankingRepo defines the interface for banking data operations
type BankingRepo interface {
	// Card operations
	GetCardByID(ctx context.Context, cardID string) (*Card, error)
	ListCards(ctx context.Context, userID string, page, pageSize int) ([]*Card, int, error)
	
	// Balance operations
	GetCardBalance(ctx context.Context, cardID string) (*CardBalance, error)
	
	// Physical card operations
	GetPhysicalCardDetails(ctx context.Context, cardID string) (*PhysicalCardDetails, error)
	
	// Remittance operations
	CreateRemittance(ctx context.Context, remittance *Remittance) (*Remittance, error)
	GetRemittance(ctx context.Context, transactionID string) (*Remittance, error)
	ListRemittances(ctx context.Context, userID string, page, pageSize int, status string) ([]*Remittance, int, error)
	UpdateRemittanceStatus(ctx context.Context, transactionID, status string) error
}

// UpstreamBankAPI defines the interface for upstream bank API
type UpstreamBankAPI interface {
	// Fetch card information from upstream
	FetchCardInfo(ctx context.Context, cardID string) (*Card, error)
	FetchCardBalance(ctx context.Context, cardID string) (*CardBalance, error)
	
	// Initiate remittance with upstream bank
	InitiateRemittance(ctx context.Context, remittance *Remittance) (*Remittance, error)
	GetRemittanceStatus(ctx context.Context, transactionID string) (string, error)
	
	// Get exchange rate
	GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*ExchangeRate, error)
}

// BankingUsecase handles banking business logic
type BankingUsecase struct {
	repo        BankingRepo
	upstreamAPI UpstreamBankAPI
	log         *log.Helper
}

// NewBankingUsecase creates a new BankingUsecase
func NewBankingUsecase(repo BankingRepo, upstreamAPI UpstreamBankAPI, logger log.Logger) *BankingUsecase {
	return &BankingUsecase{
		repo:        repo,
		upstreamAPI: upstreamAPI,
		log:         log.NewHelper(logger),
	}
}

// GetCardInfo retrieves card information
func (uc *BankingUsecase) GetCardInfo(ctx context.Context, cardID string) (*Card, error) {
	// First try to get from local database
	card, err := uc.repo.GetCardByID(ctx, cardID)
	if err != nil {
		// If not found locally, fetch from upstream
		card, err = uc.upstreamAPI.FetchCardInfo(ctx, cardID)
		if err != nil {
			return nil, err
		}
	}
	return card, nil
}

// ListCards retrieves a list of cards
func (uc *BankingUsecase) ListCards(ctx context.Context, userID string, page, pageSize int) ([]*Card, int, error) {
	return uc.repo.ListCards(ctx, userID, page, pageSize)
}

// GetCardBalance retrieves card balance
func (uc *BankingUsecase) GetCardBalance(ctx context.Context, cardID string) (*CardBalance, error) {
	// Always fetch fresh data from upstream for balance
	return uc.upstreamAPI.FetchCardBalance(ctx, cardID)
}

// GetPhysicalCardDetails retrieves physical card details
func (uc *BankingUsecase) GetPhysicalCardDetails(ctx context.Context, cardID string) (*PhysicalCardDetails, error) {
	return uc.repo.GetPhysicalCardDetails(ctx, cardID)
}

// InitiateRemittance initiates a remittance transaction
func (uc *BankingUsecase) InitiateRemittance(ctx context.Context, fromAccount, toAccount string, amount float64, fromCurrency, toCurrency, beneficiaryName, beneficiaryBank, purpose string) (*Remittance, error) {
	// Get exchange rate
	exchangeRate, err := uc.upstreamAPI.GetExchangeRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}
	
	// Calculate fee (example: 1% of amount)
	fee := amount * 0.01
	
	remittance := &Remittance{
		FromAccount:     fromAccount,
		ToAccount:       toAccount,
		Amount:          amount,
		FromCurrency:    fromCurrency,
		ToCurrency:      toCurrency,
		ExchangeRate:    exchangeRate.Rate,
		Fee:             fee,
		Status:          "pending",
		BeneficiaryName: beneficiaryName,
		BeneficiaryBank: beneficiaryBank,
		Purpose:         purpose,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Initiate with upstream bank
	result, err := uc.upstreamAPI.InitiateRemittance(ctx, remittance)
	if err != nil {
		return nil, err
	}
	
	// Save to local database
	saved, err := uc.repo.CreateRemittance(ctx, result)
	if err != nil {
		return nil, err
	}
	
	return saved, nil
}

// GetRemittanceStatus retrieves remittance status
func (uc *BankingUsecase) GetRemittanceStatus(ctx context.Context, transactionID string) (*Remittance, error) {
	// Get from local database
	remittance, err := uc.repo.GetRemittance(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	
	// Check status with upstream if not completed
	if remittance.Status != "completed" && remittance.Status != "failed" {
		status, err := uc.upstreamAPI.GetRemittanceStatus(ctx, transactionID)
		if err == nil && status != remittance.Status {
			// Update status
			uc.repo.UpdateRemittanceStatus(ctx, transactionID, status)
			remittance.Status = status
			remittance.UpdatedAt = time.Now()
		}
	}
	
	return remittance, nil
}

// ListRemittances retrieves a list of remittances
func (uc *BankingUsecase) ListRemittances(ctx context.Context, userID string, page, pageSize int, status string) ([]*Remittance, int, error) {
	return uc.repo.ListRemittances(ctx, userID, page, pageSize, status)
}

// GetExchangeRate retrieves exchange rate
func (uc *BankingUsecase) GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*ExchangeRate, error) {
	return uc.upstreamAPI.GetExchangeRate(ctx, fromCurrency, toCurrency)
}

