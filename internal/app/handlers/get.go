package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "ID is not be empty", http.StatusBadRequest)
		return
	}

	columns, err := h.storage.GetURL(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if columns.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("Location", columns.URL)

	w.WriteHeader(http.StatusTemporaryRedirect)
}
