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

func (m *MemStorage) HandlerPostRequest(w http.ResponseWriter, r *http.Request) {
	URL, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	URLShort := hex.EncodeToString(shr)

	evt := storage.MyDB{
		ID:  URLShort,
		URL: string(URL),
	}

	if app.Config.FileStoragePath != "" {
		// write url shorten to file
		producer, _ := storage.NewProducer(app.Config.FileStoragePath)
		defer producer.Close()

		if err := producer.WriteEvent(&evt); err != nil {
			log.Fatal(err)
		}
	}

	m.SaveDB(
		storage.DBKey(URLShort),
		evt,
	)

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", app.Config.BaseURL, URLShort)))
}
