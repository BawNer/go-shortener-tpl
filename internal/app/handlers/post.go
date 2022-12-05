package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

func (h *Handler) PoorPostRequestHandle(w http.ResponseWriter, r *http.Request) {
	URL, err := io.ReadAll(r.Body)
	if err != nil {
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
		URL:    string(URL),
		SignID: signID,
	}

	err = h.storage.SaveURL(
		shortURL,
		&evt,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		// должны вернуть найденную строку
		finder, err := h.storage.GetByField("url", string(URL))
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

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", app.Config.BaseURL, shortURL)))
}
