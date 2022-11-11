package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	//set env conf viper2
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	h := &handlers.MemStorage{}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/shorten", h.ShortenerHandler)
	r.Post("/", h.HandlerPostRequest)
	r.Get("/{ID}", h.HandlerGetRequest)

	log.Printf("Server started at %s", viper.GetString("SERVER_ADDRESS"))

	// start server
	err := http.ListenAndServe(viper.GetString("SERVER_ADDRESS"), r)

	// handle err
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
