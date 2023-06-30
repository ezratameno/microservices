package main

import (
	"fmt"
	"net"

	currency "github.com/ezratameno/microservices/app/services/currency/protos/currency/app/services/currency/protos"
	"github.com/ezratameno/microservices/app/services/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	log := hclog.Default()

	grpcServer := grpc.NewServer()
	currencyServer, err := server.NewCurrency(log)
	if err != nil {
		return fmt.Errorf("unable to generate rates: %w", err)
	}

	currency.RegisterCurrencyServer(grpcServer, currencyServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		return err
	}

	return grpcServer.Serve(listener)
}
