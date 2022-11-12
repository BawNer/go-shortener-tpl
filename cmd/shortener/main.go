package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

const (
	defaultServerAddr = "127.0.0.1:8080"
	defaultBaseURL    = "http://localhost:8080"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewConfig(conf Config) Config {
	if conf.ServerAddr == "" {
		conf.ServerAddr = defaultServerAddr
	}

	if conf.BaseURL == "" {
		conf.BaseURL = defaultBaseURL
	}

	return conf
}

func main() {
	viper.AutomaticEnv()

	var cfg = Config{
		ServerAddr: viper.GetString("SERVER_ADDRESS"),
		BaseURL:    viper.GetString("BASE_URL)"),
	}

	cfg = NewConfig(cfg)

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
