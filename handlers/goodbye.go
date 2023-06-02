package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	log *log.Logger
}

func NewGoodbye(log *log.Logger) *Goodbye {
	return &Goodbye{log: log}
}

func (g *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Goodbye"))
}
