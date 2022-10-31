package main

import (
	"go-shortener-tpl/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	h := &handlers.MemStorage{}

	http.HandleFunc("/", h.HandlerRequest)
	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
