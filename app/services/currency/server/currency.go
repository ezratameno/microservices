package server

import (
	"context"
	"fmt"

	"github.com/ezratameno/microservices/app/services/currency/data"
	currency "github.com/ezratameno/microservices/app/services/currency/protos/currency/app/services/currency/protos"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log          hclog.Logger
	exchangeRate *data.ExchangeRates
}

// NewCurrency creates a new Currency server
func NewCurrency(log hclog.Logger) (*Currency, error) {
	exchangeRates, err := data.NewRates(log)
	if err != nil {
		return nil, err
	}
	return &Currency{
		log:          log,
		exchangeRate: exchangeRates,
	}, nil
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, request *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", request.GetBase(), "destination", request.GetDestination())

	rate, err := c.exchangeRate.GetRate(request.Base.String(), request.Destination.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(rate)
	return &currency.RateResponse{Rate: rate}, nil
}

// mustEmbedUnimplementedCurrencyServer()
