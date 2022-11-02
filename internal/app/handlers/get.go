package handlers

import (
	"github.com/go-chi/chi/v5"
	"go-shortener-tpl/internal/app/storage"
	"net/http"
)

type MemStorage struct {
	storage.MemStorage
}

func (m *MemStorage) HandlerGetRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	if id == "" {
		http.Error(w, "ID is not be empty", http.StatusBadRequest)
		return
	}

	columns, err := m.FindByID(storage.DBKey(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Location", columns.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
