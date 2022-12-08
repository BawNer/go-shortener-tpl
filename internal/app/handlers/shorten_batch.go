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

type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	var (
		data     []RequestBatch
		response []ResponseBatch
	)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//watch cookie
	sign, _ := r.Cookie("sign")
	var signID uint32

	for _, item := range data {
		shr := uuid.New().NodeID()
		shortURL := hex.EncodeToString(shr)

		if sign == nil {
			// create cookie
			newSign := storage.CreateSign(shr[:4], app.Config.Secret)
			cookie := &http.Cookie{
				Name:   "sign",
				Value:  newSign,
				Path:   "/",
				MaxAge: 360,
			}
			http.SetCookie(w, cookie)
			signID, _ = storage.DecodeSign(newSign)
		} else {
			// work with cookie
			v, err := storage.CompareSign(sign.Value, app.Config.Secret)
			if err != nil {
				log.Println(err)
				v = 0
			}
			signID = v
		}

		evt := storage.LocalShortenData{
			ID:     shortURL,
			URL:    item.OriginalURL,
			SignID: signID,
		}

		err := h.storage.SaveURL(
			shortURL,
			&evt,
		)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response = append(response, ResponseBatch{
			CorrelationID: item.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", app.Config.BaseURL, shortURL),
		})
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
