package handlers

import (
	"net/http"

	"github.com/ezratameno/microservices/app/services/products-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.log.Debug("Inserting product: %#v\n", prod)
	data.AddProduct(prod)
}
