tidy:
	@go mod tidy

test:
	go test ./...

run: swagger
	go run .

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models 