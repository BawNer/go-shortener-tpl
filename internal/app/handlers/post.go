package handlers

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

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

	m.SaveDB(
		storage.DBKey(URLShort),
		storage.MyDB{
			ID:       len(m.Storage),
			URL:      string(URL),
			URLShort: URLShort,
		})

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", URLShort)))
}