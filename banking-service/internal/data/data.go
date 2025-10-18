package data

import (
	"banking-service/internal/biz"
	"banking-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewBankingRepo, NewUpstreamBankAPIFromConfig)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

// NewUpstreamBankAPIFromConfig creates a new upstream bank API client from config
func NewUpstreamBankAPIFromConfig(c *conf.Data, logger log.Logger) biz.UpstreamBankAPI {
	baseURL := "https://api.upstreambank.com"
	apiKey := "demo_api_key"

	if c.UpstreamBank != nil {
		if c.UpstreamBank.BaseUrl != "" {
			baseURL = c.UpstreamBank.BaseUrl
		}
		if c.UpstreamBank.ApiKey != "" {
			apiKey = c.UpstreamBank.ApiKey
		}
	}

	return NewUpstreamBankAPI(baseURL, apiKey, logger)
}
