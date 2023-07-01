package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ezratameno/microservices/app/services/currency/data"
	currency "github.com/ezratameno/microservices/app/services/currency/protos/currency/app/services/currency/protos"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/gommon/log"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log           hclog.Logger
	exchangeRates *data.ExchangeRates

	// subscriptions will tell us who has subscribed to a particular rate and then we will
	// send them a message when there is an update.
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

// NewCurrency creates a new Currency server
func NewCurrency(log hclog.Logger) (*Currency, error) {

	exchangeRates, err := data.NewRates(log)
	if err != nil {
		return nil, err
	}

	c := &Currency{
		log:           log,
		exchangeRates: exchangeRates,
		subscriptions: make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest),
	}

	go c.handleUpdates()

	return c, nil
}

func (c *Currency) handleUpdates() {
	ratesMonitor := c.exchangeRates.MonitorRates(5 * time.Second)

	for range ratesMonitor {
		log.Info("got updated rates")

		// loop over subscribe clients
		for client, rates := range c.subscriptions {

			// loop over subscribed rates
			for _, rateReq := range rates {
				rate, err := c.exchangeRates.GetRate(rateReq.Base.String(), rateReq.Destination.String())
				if err != nil {
					c.log.Error("unable to get updated rate", "base", rateReq.Base.String(),
						"destination", rateReq.Destination.String())
					continue
				}

				err = client.Send(&currency.RateResponse{Rate: rate, Base: rateReq.Base, Destination: rateReq.Destination})
				if err != nil {
					c.log.Error("unable to send updated rate", "base", rateReq.Base.String(),
						"destination", rateReq.Destination.String(), "err", err)
				}
			}
		}

	}

}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, request *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", request.GetBase(), "destination", request.GetDestination())

	rate, err := c.exchangeRates.GetRate(request.Base.String(), request.Destination.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(rate)
	return &currency.RateResponse{Rate: rate, Base: request.Base, Destination: request.Destination}, nil
}

// SubscribeRates shouldn't exist as long as the connection is open.
func (c *Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {

	for {
		// Recv is blocking
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("unable to read from client")
			return err
		}

		c.log.Info("handle client request", rr)

		// make the client subscribe to get updates on this rates
		c.subscriptions[src] = append(c.subscriptions[src], rr)
	}

	return nil

}

// mustEmbedUnimplementedCurrencyServer()
