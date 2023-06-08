tidy:
	@go mod tidy

test:
	go test ./...

run: swagger
	go run .

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models 

# swagger-generate will create a go code based on the swagger docs.
swagger-generate: swagger
	@cd sdk && \
	swagger generate client -f ../swagger.yaml -A product-api
