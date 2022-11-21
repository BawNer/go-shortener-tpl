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

func HandlerPostRequest(w http.ResponseWriter, r *http.Request) {
	URL, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shr := uuid.New().NodeID()
	shortURL := hex.EncodeToString(shr)

	evt := storage.LocalShortenData{
		ID:  shortURL,
		URL: string(URL),
	}

	if app.Config.FileStoragePath != "" {
		// write url shorten to file
		producer, err := storage.NewProducer(app.Config.FileStoragePath)

		if err != nil {
			log.Fatal(err.Error())
		}

		defer producer.Close()

		if err := producer.WriteEvent(&evt); err != nil {
			log.Fatal(err)
		}
	}

	app.Memory.InMemory.Save(
		storage.DBKey(shortURL),
		evt,
	)

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("%s/%s", app.Config.BaseURL, shortURL)))
}
