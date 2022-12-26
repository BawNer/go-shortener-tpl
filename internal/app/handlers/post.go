package handlers

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

func (h *Handler) HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()
	log.Printf("reqID=%s handle request HandlePostRequest", reqID)

	URL, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReqID=%s, can't read body! Err: %v", reqID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	shortURL := hex.EncodeToString(shr)

	//watch cookie
	sign, _ := r.Cookie("sign")
	var (
		signID        uint32
		signDecodeErr error
	)
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
		signID, signDecodeErr = storage.DecodeSign(newSign)
		if signDecodeErr != nil {
			log.Println(signDecodeErr.Error())
		}
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
		URL:       string(URL),
		SignID:    signID,
		IsDeleted: false,
	}

	w.Header().Set("Content-Type", "text/plain")

	log.Printf("ReqID=%s, Start save url", reqID)
	err = h.storage.SaveURL(
		shortURL,
		&evt,
	)
	if err != nil {
		log.Printf("ReqID=%s, Save url with err: %v", reqID, err)
		w.WriteHeader(http.StatusConflict)
		// должны вернуть найденную строку
		log.Printf("ReqID=%s, Start find exist url", reqID)
		finder, err := h.storage.GetByField("url", string(URL))
		log.Printf("ReqID=%s, Start Finded url is %v", reqID, finder)
		if err != nil {
			log.Printf("ReqID=%s, Err when find exist url: %v", reqID, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", app.Config.BaseURL, finder.ID)))
		return
	}
	log.Printf("ReqID=%s, URL Saved success!", reqID)

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", app.Config.BaseURL, shortURL)))
}
