package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{
		log:   l,
		rates: make(map[string]float64),
	}

	err := er.getExchangeRates()
	if err != nil {
		return nil, err
	}

	return er, nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	if _, ok := e.rates[base]; !ok {
		return 0, fmt.Errorf("rate note found for currency %s", base)
	}
	if _, ok := e.rates[dest]; !ok {
		return 0, fmt.Errorf("rate note found for currency %s", dest)
	}

	return e.rates[dest] / e.rates[base], nil
}

// getExchangeRates will get the exchange rates from the central bank.
func (e *ExchangeRates) getExchangeRates() error {

	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	md := &Cubes{}

	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil

}
