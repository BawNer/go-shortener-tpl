package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/database"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/file"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/memory"
	"github.com/BawNer/go-shortener-tpl/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var repository storage.Storage

func main() {
	if app.Config.DB == "" {
		if app.Config.FileStoragePath != "" {
			repository, _ = file.New(app.Config.FileStoragePath)
			err := repository.Init()
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			repository, _ = memory.New()
		}
	} else {
		// TODO: Init postgres conn
		repository, _ = database.New()
	}

	h := handlers.NewHandler(repository)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.Decompress)

	r.Post("/api/shorten", h.ShortenHandle)
	r.Get("/api/user/urls", h.UrlsUserHandle)
	r.Post("/", h.PoorPostRequestHandle)
	r.Get("/{ID}", h.PoorGetRequestHandle)
	r.Get("/ping", h.PingDBConn)

	log.Printf("Server started at %s", app.Config.ServerAddr)

	// start server
	err := http.ListenAndServe(app.Config.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
