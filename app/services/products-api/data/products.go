package data

import (
	"context"
	"fmt"

	currency "github.com/ezratameno/microservices/app/services/currency/protos/currency/app/services/currency/protos"
	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

type ProductDB struct {
	currencyClient       currency.CurrencyClient
	log                  hclog.Logger
	rates                map[string]float64
	subscribeRatesClient currency.Currency_SubscribeRatesClient
}

func NewProductDB(c currency.CurrencyClient, log hclog.Logger) *ProductDB {
	pb := &ProductDB{
		currencyClient: c,
		log:            log,
		rates:          make(map[string]float64),
	}

	go pb.handleUpdates()

	return pb
}

func (p *ProductDB) handleUpdates() {
	sub, err := p.currencyClient.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("enable to subscribe to rates", err)
		return
	}

	p.subscribeRatesClient = sub

	for {
		rateResp, err := sub.Recv()
		if err != nil {
			p.log.Error("error receiving message", err)
		}
		p.log.Info("received updated rate from server", "dest", rateResp.Base.String())

		// update the rate
		p.rates[rateResp.Destination.String()] = rateResp.Rate

	}

}

func (p *ProductDB) GetProducts(destCurrency string) (Products, error) {

	if destCurrency == "" {
		return productList, nil
	}

	// update the products with the new rates

	rate, err := p.getRate(destCurrency)
	if err != nil {
		return nil, err
	}

	var prodList Products
	for _, prod := range productList {
		tmp := *prod
		tmp.Price = tmp.Price * rate

		prodList = append(prodList, &tmp)
	}

	return prodList, nil
}

// GetProducts returns all products from the database
func GetProducts() Products {
	return productList
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductDB) GetProductByID(id int, destCurrency string) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	// calculate the exchange rate and update the product price
	if destCurrency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(destCurrency)
	if err != nil {
		return nil, err
	}

	prod := *productList[i]
	prod.Price = prod.Price * rate

	return &prod, nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (p *ProductDB) UpdateProduct(prod Product) error {
	i := findIndexByProductID(prod.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &prod

	return nil
}

// AddProduct adds a new product to the database
func AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductDB) getRate(destination string) (float64, error) {

	// check cache
	if rate, ok := p.rates[destination]; ok {
		return rate, nil
	}

	// make a request if not in the cache
	rateReq := &currency.RateRequest{
		Base:        currency.Currencies_EUR,
		Destination: currency.Currencies(currency.Currencies_value[destination]),
	}

	// get initial rate
	rateResp, err := p.currencyClient.GetRate(context.Background(), rateReq)
	if err != nil {
		return 0, fmt.Errorf("unable to get rate: %w", err)
	}

	// update cache
	p.rates[destination] = rateResp.Rate

	// subscribe for updates on this rate
	p.subscribeRatesClient.Send(rateReq)
	return rateResp.Rate, nil
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
