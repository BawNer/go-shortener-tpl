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

var (
	cfg       = app.NewConfigApp()
	ConfigApp = cfg()
)

func (m *MemStorage) HandlerPostRequest(w http.ResponseWriter, r *http.Request) {
	URL, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	URLShort := hex.EncodeToString(shr)

	if ConfigApp.FileStoragePath == "" {
		m.SaveDB(
			storage.DBKey(URLShort),
			storage.MyDB{
				ID:       len(m.Storage),
				URL:      string(URL),
				URLShort: URLShort,
			})
	} else {
		// write url shorten to file
		producer, _ := storage.NewProducer(ConfigApp.FileStoragePath)
		defer producer.Close()
		if err := producer.WriteEvent(&storage.Event{ShortenURL: URLShort}); err != nil {
			log.Fatal(err)
		}
	}

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", ConfigApp.BaseURL, URLShort)))
}
