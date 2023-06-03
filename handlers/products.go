package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ezratameno/microservices/data"
)

type Products struct {
	log *log.Logger
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log: log}
}

// ServeHTTP is the main entry point for the handler and satisfies the http.Handler interface
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		p.getProducts(w, r)

	case http.MethodPost:
		p.addProducts(w, r)
	case http.MethodPut:
		// expect the id in the URI
		path := r.URL.Path

		reg := regexp.MustCompile(`/([0-9]+)$`)

		groups := reg.FindAllStringSubmatch(path, -1)

		if len(groups) != 1 {
			http.Error(w, "Invalid URI, more than one id", http.StatusBadRequest)
			return
		}

		if len(groups[0]) != 2 {
			http.Error(w, "Invalid URI, more than one capture group", http.StatusBadRequest)
			return
		}

		idString := groups[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid ID, id must be a number", http.StatusBadRequest)
			return
		}

		p.log.Println("got id", id)

		p.updateProducts(id, w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (p *Products) updateProducts(id int, w http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle PUT products.")

	var product data.Product

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, &product)
	if err != nil {
		if err == data.ErrProductNotFound {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update product", http.StatusInternalServerError)
		return
	}

	p.log.Printf("Prod: %+v\n", product)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle GET products.")

	// fetch the products from the datastore
	products := data.GetProducts()

	// serialize the list to JSON
	err := products.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProducts(w http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle POST products.")

	var product data.Product

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	data.AddProduct(&product)
	p.log.Printf("Prod: %+v\n", product)

}
