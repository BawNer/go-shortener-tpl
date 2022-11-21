package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

type RequestData struct {
	URL string `json:"url"`
}

type ResponseData struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	shortURL := hex.EncodeToString(shr)

	evt := storage.LocalShortenData{
		ID:  shortURL,
		URL: data.URL,
	}

	err := h.storage.SaveURL(
		shortURL,
		&evt,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	response := ResponseData{
		Result: fmt.Sprintf("%s/%s", app.Config.BaseURL, shortURL),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(buf.Bytes())
}
