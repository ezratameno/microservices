package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return err
	}
	return nil
}
