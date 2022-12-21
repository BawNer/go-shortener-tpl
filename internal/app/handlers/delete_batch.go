package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

type DataForWorker struct {
	ID     string
	SignID uint32
}

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
	go putJobs(h.inputCh, urlIDs, signID)

	log.Printf("reqID=%s Отдаем ответ со статусом 202", reqID)
	w.WriteHeader(http.StatusAccepted)
}

func putJobs(inputCh chan<- DataForWorker, urlIDs []string, signID uint32) {
	// складируем данные в канал
	for _, urlID := range urlIDs {
		inputCh <- DataForWorker{
			ID:     urlID,
			SignID: signID,
		}
	}
}

func (h *Handler) Worker(inputCh <-chan DataForWorker) {
	log.Printf("Воркер запущен!")
	for {
		log.Printf("генерируем батч")
		batches := map[uint32][]string{}
		log.Printf("Наполняем джобу")
		for job := range inputCh {
			batches[job.SignID] = append(batches[job.SignID], job.ID)
		}
		log.Printf("Отправляем данные в БД!")
		for signID, ids := range batches {
			err := h.writeToDB(ids, signID)
			if err != nil {
				log.Printf("Произошла ошибка при отпрвке в бд  %v", err.Error())
			}
		}
	}
}

func (h *Handler) writeToDB(ids []string, signID uint32) error {
	for _, id := range ids {
		log.Printf("url id to delete is %s", id)
		err := h.storage.DeleteURL(id, true, signID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
