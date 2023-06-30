package server

import (
	"context"
	"fmt"
	"io"
	"time"

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

func (c *Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {

	go func() {
		for {
			// Recv is blocking
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Info("client has closed connection")
				break
			}
			if err != nil {
				c.log.Error("unable to read from client")
				break
			}

			c.log.Info("handle client request", rr)
		}
	}()

	// send updates all the times.
	for {

		// update rates.
		err := src.Send(&currency.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}

}

// mustEmbedUnimplementedCurrencyServer()
