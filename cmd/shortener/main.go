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
		if errConfInit != nil {
			log.Fatal(errConfInit.Error())
		}
		repository, errConfInit = database.New()
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
	r.Post("/api/shorten/batch", h.ShortenBatch)
	r.Get("/api/user/urls", h.UrlsUserHandle)
	r.Get("/ping", h.PingDBConn)
	r.Get("/{ID}", h.PoorGetRequestHandle)
	r.Post("/", h.PoorPostRequestHandle)

	log.Printf("Server started at %s", app.Config.ServerAddr)

	// start server
	err := http.ListenAndServe(app.Config.ServerAddr, r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
