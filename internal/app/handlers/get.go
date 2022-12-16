package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	log.Printf("reqID=%s handle request HandleGetRequest", reqID)

	log.Printf("reqID=%s Start parse url, find {ID}", reqID)
	id := chi.URLParam(r, "ID")
	if id == "" {
		log.Printf("reqID=%s ID from url not finded %s", reqID, id)
		http.Error(w, "ID is not be empty", http.StatusBadRequest)
		return
	}
	log.Printf("reqID=%s ID from url finded %s", reqID, id)

	log.Printf("reqID=%s Start find %s from storage", reqID, id)
	columns, err := h.storage.GetURL(id)
	if err != nil {
		log.Printf("reqID=%s Not finded %s from storage", reqID, id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("reqID=%s Finded %s from storage", reqID, id)
	if columns.IsDeleted {
		log.Printf("reqID=%s Link %s is deleted!", reqID, id)
		w.WriteHeader(http.StatusGone)
		return
	}

	log.Printf("reqID=%s Link %s set to lacation", reqID, id)
	w.Header().Set("Location", columns.URL)

	w.WriteHeader(http.StatusTemporaryRedirect)
}
