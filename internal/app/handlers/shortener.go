package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	cfg := app.NewConfig(&app.Config{
		ServerAddr: os.Getenv("SERVER_ADDRESS"),
		BaseURL:    os.Getenv("BASE_URL"),
	})

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
		Result: fmt.Sprintf("%s/%s", cfg.BaseURL, URLShort),
	}

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(buf.Bytes())
}
