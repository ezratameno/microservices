package server

import (
	"context"

	currency "github.com/ezratameno/microservices/app/services/currency/protos/currency/app/services/currency/protos"
	"github.com/hashicorp/go-hclog"
)

// Currency implements the currency server we defined in the .proto file.
type Currency struct {
	log hclog.Logger
}

func NewCurrency(log hclog.Logger) *Currency {
	return &Currency{
		log: log,
	}
}

func (c *Currency) GetRate(ctx context.Context, request *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", request.GetBase(), "destination", request.GetDestination())

	return &currency.RateResponse{Rate: 0.5}, nil
}

// mustEmbedUnimplementedCurrencyServer()
