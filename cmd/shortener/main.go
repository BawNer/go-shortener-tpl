package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := app.NewConfig(&app.Config{
		ServerAddr: os.Getenv("SERVER_ADDRESS"),
		BaseURL:    os.Getenv("BASE_URL"),
	})

	h := &handlers.MemStorage{}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/shorten", h.ShortenerHandler)
	r.Post("/", h.HandlerPostRequest)
	r.Get("/{ID}", h.HandlerGetRequest)

	log.Printf("Server started at %s", cfg.ServerAddr)

	// start server
	err := http.ListenAndServe(cfg.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
