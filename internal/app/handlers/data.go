package handlers

import "github.com/BawNer/go-shortener-tpl/internal/app/storage"

type Handler struct {
	storage storage.Storage
	inputCh chan DataForWorker
}

func NewHandler(repository storage.Storage, inputCh chan DataForWorker) *Handler {
	return &Handler{
		storage: repository,
		inputCh: inputCh,
	}
}
