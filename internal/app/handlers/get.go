package handlers

import (
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

type Repository struct {
	storage.Repository
}

func (m *Repository) HandlerGetRequest(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
