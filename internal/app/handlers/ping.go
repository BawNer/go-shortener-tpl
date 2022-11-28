package handlers

import (
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/postgres"
)

func (h *Handler) PingDBConn(w http.ResponseWriter, r *http.Request) {
	_, err := postgres.New()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	return
}
