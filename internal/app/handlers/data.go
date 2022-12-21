package handlers

import (
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Handler struct {
	storage storage.Storage
	inputCh chan DataForWorker
	wg      sync.WaitGroup
}

func NewHandler(repository storage.Storage, inputCh chan DataForWorker) *Handler {
	return &Handler{
		storage: repository,
		inputCh: inputCh,
		wg:      sync.WaitGroup{},
	}
}
