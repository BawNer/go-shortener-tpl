package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

type MemStorage struct {
	storage.MemStorage
}

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

	m.SaveDB(
		storage.DBKey(URLShort),
		storage.MyDB{
			ID:       len(m.Storage),
			URL:      data.URL,
			URLShort: URLShort,
		})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := ResponseData{
		Result: fmt.Sprintf("http://localhost:8080/%s", URLShort),
	}

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(buf.String()))
}
