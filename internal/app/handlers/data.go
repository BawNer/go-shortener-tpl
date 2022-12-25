package handlers

import (
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/workerpool"
)

type Handler struct {
	storage storage.Storage
	worker  *workerpool.WorkerPool
}

func NewHandler(repository storage.Storage, worker *workerpool.WorkerPool) *Handler {
	return &Handler{
		storage: repository,
		worker:  worker,
	}
}
