package data

import (
	"context"
	"fmt"
	"time"

	"banking-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type upstreamBankAPI struct {
	baseURL string
	apiKey  string
	log     *log.Helper
}

// NewUpstreamBankAPI creates a new upstream bank API client
func NewUpstreamBankAPI(baseURL, apiKey string, logger log.Logger) biz.UpstreamBankAPI {
	return &upstreamBankAPI{
		baseURL: baseURL,
		apiKey:  apiKey,
		log:     log.NewHelper(logger),
	}
}

// FetchCardInfo fetches card information from upstream bank API
func (api *upstreamBankAPI) FetchCardInfo(ctx context.Context, cardID string) (*biz.Card, error) {
	// TODO: Implement actual HTTP request to upstream bank API
	api.log.WithContext(ctx).Infof("FetchCardInfo from upstream: %s", cardID)
	
	// Mock response
	return &biz.Card{
		ID:               cardID,
		UserID:           "user_123",
		CardNumberMasked: "**** **** **** 1234",
		CardType:         "credit",
		CardBrand:        "visa",
		Status:           "active",
		ExpiryDate:       time.Now().AddDate(2, 0, 0),
		CreatedAt:        time.Now().Add(-365 * 24 * time.Hour),
	}, nil
}

// FetchCardBalance fetches card balance from upstream bank API
func (api *upstreamBankAPI) FetchCardBalance(ctx context.Context, cardID string) (*biz.CardBalance, error) {
	// TODO: Implement actual HTTP request to upstream bank API
	api.log.WithContext(ctx).Infof("FetchCardBalance from upstream: %s", cardID)
	
	// Mock response
	return &biz.CardBalance{
		CardID:             cardID,
		AvailableBalance:   5000.00,
		CreditLimit:        10000.00,
		OutstandingBalance: 5000.00,
		Currency:           "USD",
		LastUpdated:        time.Now(),
	}, nil
}

// InitiateRemittance initiates a remittance with upstream bank API
func (api *upstreamBankAPI) InitiateRemittance(ctx context.Context, remittance *biz.Remittance) (*biz.Remittance, error) {
	// TODO: Implement actual HTTP request to upstream bank API
	api.log.WithContext(ctx).Infof("InitiateRemittance with upstream: %s -> %s, amount: %.2f %s",
		remittance.FromAccount, remittance.ToAccount, remittance.Amount, remittance.FromCurrency)
	
	// Mock response
	remittance.TransactionID = fmt.Sprintf("txn_%s", time.Now().Format("20060102150405"))
	remittance.Status = "processing"
	remittance.CreatedAt = time.Now()
	remittance.UpdatedAt = time.Now()
	
	return remittance, nil
}

// GetRemittanceStatus gets remittance status from upstream bank API
func (api *upstreamBankAPI) GetRemittanceStatus(ctx context.Context, transactionID string) (string, error) {
	// TODO: Implement actual HTTP request to upstream bank API
	api.log.WithContext(ctx).Infof("GetRemittanceStatus from upstream: %s", transactionID)
	
	// Mock response - randomly return different statuses
	statuses := []string{"processing", "completed", "failed"}
	return statuses[time.Now().Unix()%int64(len(statuses))], nil
}

// GetExchangeRate gets exchange rate from upstream bank API
func (api *upstreamBankAPI) GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*biz.ExchangeRate, error) {
	// TODO: Implement actual HTTP request to upstream bank API or third-party service
	api.log.WithContext(ctx).Infof("GetExchangeRate from upstream: %s -> %s", fromCurrency, toCurrency)
	
	// Mock exchange rates
	rates := map[string]map[string]float64{
		"USD": {
			"EUR": 0.85,
			"GBP": 0.75,
			"JPY": 110.0,
			"CNY": 6.5,
		},
		"EUR": {
			"USD": 1.18,
			"GBP": 0.88,
		},
		"GBP": {
			"USD": 1.33,
			"EUR": 1.14,
		},
	}
	
	rate := 1.0
	if fromRates, ok := rates[fromCurrency]; ok {
		if r, ok := fromRates[toCurrency]; ok {
			rate = r
		}
	}
	
	return &biz.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
		Timestamp:    time.Now(),
	}, nil
}

