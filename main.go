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
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {

	log := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handlers.NewHello(log)
	gh := handlers.NewGoodbye(log)
	mux := http.NewServeMux()

	mux.Handle("/", hh)
	mux.Handle("/goodbye", gh)

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
