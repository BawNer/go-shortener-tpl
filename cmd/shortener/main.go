package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
			err := repository.Init()
			if err != nil {
				log.Fatal(err.Error())
			}
			if errConfInit != nil {
				log.Fatal(errConfInit.Error())
			}
		}
	} else {
		var errConfInit error
		repository, errConfInit = database.New()
		err := repository.Init()
		if err != nil {
			log.Fatal(err.Error())
		}
		if errConfInit != nil {
			log.Fatal(errConfInit.Error())
		}

		repository.RunWorkers(app.Config.Workers)
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

	server := &http.Server{Addr: app.Config.ServerAddr, Handler: r}

	// start server
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Сервер не смог закрыться без ошибок: %v", err)
	}

	repository.Stop()
	log.Printf("Channel has been closed")

	repository.Wait()

}
