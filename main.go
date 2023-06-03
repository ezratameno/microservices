package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezratameno/microservices/handlers"
	"github.com/gorilla/mux"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {

	log := log.New(os.Stdout, "product-api", log.LstdFlags)

	productHandler := handlers.NewProducts(log)

	// create a new serve mux and register the handlers.
	mux := mux.NewRouter()

	// Get methods.
	getRouter := mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)

	// Put methods.
	putRouter := mux.Methods(http.MethodPut).Subrouter()

	// add middleware
	putRouter.Use(productHandler.MiddlewareProductValidation)

	// creates an id var.
	putRouter.HandleFunc("/{id:[0-9]+$}", productHandler.UpdateProducts)

	// Post methods.
	postRouter := mux.Methods(http.MethodPost).Subrouter()

	// add middleware
	postRouter.Use(productHandler.MiddlewareProductValidation)

	postRouter.HandleFunc("/", productHandler.AddProducts)

	srv := http.Server{
		Handler:      mux,
		Addr:         ":9090",
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	// blocking.
	sig := <-sigChan

	log.Println("Received terminate, graceful shut down", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shutdown the server.
	return srv.Shutdown(tc)
}
