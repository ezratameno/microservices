package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ezratameno/microservices/data"
	"github.com/gorilla/mux"
)

type Products struct {
	log *log.Logger
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log: log}
}

func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle PUT products.")

	// Contains the id var.
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	// get the product from the context
	product := r.Context().Value(KeyProduct{}).(data.Product)

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

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
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

func (p *Products) AddProducts(w http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle POST products.")

	// get the product from the context
	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
	p.log.Printf("Prod: %+v\n", product)

}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var product data.Product
		err := product.FromJSON(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal json: %+v", err), http.StatusBadRequest)
			return
		}

		// Validate the product.
		err = product.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error validating product: %+v", err), http.StatusBadRequest)
			return
		}

		// add the product to the context.
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)

		// Call the next handler.
		next.ServeHTTP(w, req)
	})
}
