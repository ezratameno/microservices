tidy:
	@go mod tidy

test:
	go test ./...

run-products-api: swagger
	go run ./app/services/products-api/

run-products-images:
	go run ./app/services/products-images/

run-currency-server:
	go run ./app/services/currency/

swagger:
	swagger generate spec -o ./app/services/products-api/swagger.yaml --scan-models 

# swagger-generate will create a go code based on the swagger docs.
swagger-generate: swagger
	@cd ./app/services/products-api/sdk && \
	swagger generate client -f ../swagger.yaml -A product-api

# will format the proto file to look like go code, using the .clang-format file.
clang-format:
	clang-format -i app/services/currency/protos/*.proto


protos:
	protoc -I ./app/services/currency/protos/ ./app/services/currency/protos/currency.proto --go-grpc_out=app/services/currency/protos/currency
	go mod tidy

# will generate our code.
buf/generate:
	buf generate
	go mod tidy

# will lint our protofiles.
buf/lint:
	buf lint

# will list our services
grpcurl-list:
	grpcurl --plaintext localhost:9092 list


# A call to the GetRate func of Currency server, with the request details.
# grpcurl --plaintext  -d '{"Base": "base dfd", "Destination": "dest"}' localhost:9092 Currency.GetRate