package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	log *log.Logger
}

func NewHello(log *log.Logger) *Hello {
	return &Hello{
		log: log,
	}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Println("Hello world")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello %s", d)
}
