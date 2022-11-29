package handlers

import (
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage/database"
)

func (h *Handler) PingDBConn(w http.ResponseWriter, r *http.Request) {
	_, err := database.New()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
