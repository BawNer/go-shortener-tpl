package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/file"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/memory"
	"github.com/BawNer/go-shortener-tpl/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var repository storage.Storage

func main() {
	if app.Config.FileStoragePath != "" {
		repository, _ = file.New(app.Config.FileStoragePath)
	}
	repository, _ = memory.New()

	h := handlers.NewHandler(repository)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.Decompress)

	r.Post("/api/shorten", h.ShortenerHandler)
	r.Post("/", h.HandlerPostRequest)
	r.Get("/{ID}", h.HandlerGetRequest)

	log.Printf("Server started at %s", app.Config.ServerAddr)

	// start server
	err := http.ListenAndServe(app.Config.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
