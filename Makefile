tidy:
	@go mod tidy

test:
	go test ./...

run-products-api: swagger
	go run ./app/services/products-api/

run-products-images:
	go run ./app/services/products-images/

swagger:
	swagger generate spec -o ./app/services/products-api/swagger.yaml --scan-models 

# swagger-generate will create a go code based on the swagger docs.
swagger-generate: swagger
	@cd ./app/services/products-api/sdk && \
	swagger generate client -f ../swagger.yaml -A product-api
