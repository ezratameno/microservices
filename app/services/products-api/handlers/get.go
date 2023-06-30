package handlers

import (
	"net/http"

	"github.com/ezratameno/microservices/app/services/products-api/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("get all records")
	cur := r.URL.Query().Get("currency")

	prods, err := p.productsDB.GetProducts(cur)
	if err != nil {
		p.log.Error("failed getting products", err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just in case
		p.log.Error("serializing product", err)
		return
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	cur := r.URL.Query().Get("currency")
	p.log.Debug("get record id", id)

	prod, err := p.productsDB.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.log.Error("fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.log.Error("etching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.log.Error("serializing product", err)
	}
}
