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

func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	shortURL := hex.EncodeToString(shr)

	//watch cookie
	sign, _ := r.Cookie("sign")
	var signID uint32
	if sign == nil {
		// create cookie
		newSign := storage.CreateSign(shr[:4], app.Config.Secret)
		cookie := &http.Cookie{
			Name:   "sign",
			Value:  newSign,
			Path:   "/",
			MaxAge: 3600,
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
		ID:        shortURL,
		URL:       data.URL,
		SignID:    signID,
		IsDeleted: false,
	}

	err := h.storage.SaveURL(
		shortURL,
		&evt,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		// должны вернуть найденную строку
		finder, err := h.storage.GetByField("url", data.URL)
		if err != nil {
			log.Println(err.Error())
			return
		}
		response := ResponseData{
			Result: fmt.Sprintf("%s/%s", app.Config.BaseURL, finder.ID),
		}
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(&response); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = w.Write(buf.Bytes())
		return
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
