package handlers

import "github.com/BawNer/go-shortener-tpl/internal/app/storage"

type Handler struct {
	storage storage.Storage
}

func NewHandler(repository storage.Storage) *Handler {
	return &Handler{
		storage: repository,
	}
}
