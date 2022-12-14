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
	if app.Config.DSN == "" {
		var errConfInit error
		if app.Config.FileStoragePath != "" {
			repository, errConfInit = file.New(app.Config.FileStoragePath)
			if errConfInit != nil {
				log.Fatal(errConfInit.Error())
			}
			err := repository.Init()
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			repository, errConfInit = memory.New()
			if errConfInit != nil {
				log.Fatal(errConfInit.Error())
			}
		}
	} else {
		var errConfInit error
		repository, errConfInit = database.New()
		if errConfInit != nil {
			log.Fatal(errConfInit.Error())
		}
	}

	h := handlers.NewHandler(repository)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.Decompress)

	r.Get("/api/user/urls", h.HandleUserURLs)
	r.Get("/ping", h.PingDBConn)
	r.Get("/{ID}", h.HandleGetRequest)

	r.Post("/api/shorten", h.HandleShorten)
	r.Post("/api/shorten/batch", h.ShortenBatch)
	r.Post("/", h.HandlePostRequest)

	r.Delete("/api/user/urls", h.HandleDeleteBatchUrls)

	log.Printf("Server started at %s", app.Config.ServerAddr)

	// start server
	err := http.ListenAndServe(app.Config.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
