package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	cfg       = app.NewConfigApp()
	ConfigApp = cfg()
)

func main() {
	h := &handlers.MemStorage{}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/shorten", h.ShortenerHandler)
	r.Post("/", h.HandlerPostRequest)
	r.Get("/{ID}", h.HandlerGetRequest)

	log.Printf("Server started at %s", ConfigApp.ServerAddr)

	// start server
	err := http.ListenAndServe(ConfigApp.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
