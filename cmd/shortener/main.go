package main

import (
	"go-shortener-tpl/internal/app"
	"log"
	"net/http"
)

var LocalDB []app.DB

func main() {

	http.HandleFunc("/", app.HandlerRequest)

	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
