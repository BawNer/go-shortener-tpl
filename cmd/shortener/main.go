package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-shortener-tpl/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	h := &handlers.MemStorage{}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", h.HandlerPostRequest)
	r.Get("/{ID}", h.HandlerGetRequest)

	// start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
