package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Response struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *Handler) UrlsUserHandle(w http.ResponseWriter, r *http.Request) {
	sign, err := r.Cookie("sign")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// work with cookie
	signID, err := storage.CompareSign(sign.Value, app.Config.Secret)
	if err != nil {
		log.Println(err)
		signID = 0
	}

	// find all urls
	v, ok := h.storage.GetAllURL(signID)
	if ok != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var (
		response []Response
		buf      bytes.Buffer
	)
	for _, item := range v {
		response = append(response, Response{
			ShortURL:    fmt.Sprintf("%s/%s", app.Config.BaseURL, item.ID),
			OriginalURL: item.URL,
		})
	}

	if err := json.NewEncoder(&buf).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(buf.Bytes())
}
