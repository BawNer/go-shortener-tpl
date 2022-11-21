package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandlerGetRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	if id == "" {
		http.Error(w, "ID is not be empty", http.StatusBadRequest)
		return
	}

	columns, err := h.storage.GetURL(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", columns.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
