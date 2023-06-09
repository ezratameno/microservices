package handlers

import (
	"net/http"

	"github.com/ezratameno/microservices/app/services/products-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.log.Debug("updating record id", prod.ID)

	err := p.productsDB.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.log.Error("product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
