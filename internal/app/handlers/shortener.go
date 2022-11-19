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

func (m *MemStorage) ShortenerHandler(w http.ResponseWriter, r *http.Request) {

	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	URLShort := hex.EncodeToString(shr)

	evt := storage.MyDB{
		ID:  URLShort,
		URL: data.URL,
	}

	if app.Config.FileStoragePath != "" {
		// write url shorten to file
		producer, err := storage.NewProducer(app.Config.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}

		defer producer.Close()

		if err := producer.WriteEvent(&evt); err != nil {
			log.Fatal(err)
		}
	}

	m.SaveDB(
		storage.DBKey(URLShort),
		evt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := ResponseData{
		Result: fmt.Sprintf("%s/%s", app.Config.BaseURL, URLShort),
	}

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(buf.Bytes())
}
