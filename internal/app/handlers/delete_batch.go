package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

func (h *Handler) HandleDeleteBatchUrls(w http.ResponseWriter, r *http.Request) {
	var (
		signID uint32
		urlIDs []string
	)

	reqID := uuid.New().String()
	log.Printf("reqID=%s handle request HandleDeleteBatchUrls", reqID)

	log.Printf("reqID=%s Получаем id's из запроса", reqID)
	if err := json.NewDecoder(r.Body).Decode(&urlIDs); err != nil {
		log.Printf("reqID=%s Произошла ошибка при парсинге запроса %v", reqID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("reqID=%s ID получены, их количество %d", reqID, len(urlIDs))

	log.Printf("reqID=%s Получаем sign из куки", reqID)
	sign, err := r.Cookie("sign")
	if err != nil {
		log.Printf("reqID=%s Кука не найдена %v", reqID, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// work with cookie
	log.Printf("reqID=%s Начинаем сравнивать подписи", reqID)
	signID, err = storage.CompareSign(sign.Value, app.Config.Secret)
	if err != nil {
		log.Printf("reqID=%s Подписи не одинаковые, удаление невозможно, %v", reqID, err)
		signID = 0
	}
	// usr not auth
	if signID == 0 {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// main

	log.Printf("reqID=%s Складируем в канал ID", reqID)
	h.storage.AddJob(urlIDs, signID)
	log.Printf("reqID=%s Отдаем ответ со статусом 202", reqID)
	w.WriteHeader(http.StatusAccepted)
}
