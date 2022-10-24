package main

import (
	"go-shortener-tpl/internal/app"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", app.HandlerRequest)
	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
