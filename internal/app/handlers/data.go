package handlers

import (
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/workers"
)

type Handler struct {
	storage storage.Storage
	worker  *workers.WorkerPool
}

func NewHandler(repository storage.Storage, worker *workers.WorkerPool) *Handler {
	return &Handler{
		storage: repository,
		worker:  worker,
	}
}
