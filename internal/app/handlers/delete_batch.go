package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/google/uuid"
)

type dataForWorker struct {
	ID     string
	SignID uint32
}

var batchSize int

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
	log.Printf("reqID=%s Получаем размера буфера", reqID)
	batchSize = len(urlIDs)
	log.Printf("reqID=%s Буфер установлен на %d", reqID, batchSize)

	log.Printf("reqID=%s Создаем канал с буфером %d", reqID, batchSize)
	inputCh := make(chan dataForWorker, batchSize)
	log.Printf("reqID=%s Канал с буфером %d создан!", reqID, batchSize)

	log.Printf("reqID=%s Складируем в канал ID", reqID)
	go putJobs(inputCh, urlIDs, signID)

	log.Printf("reqID=%s Запускаем рутину на удаление", reqID)
	go h.worker(inputCh) // init go routine

	log.Printf("reqID=%s Отдаем ответ со статусом 202", reqID)
	w.WriteHeader(http.StatusAccepted)
}

func putJobs(inputCh chan<- dataForWorker, urlIDs []string, signID uint32) {
	// складируем данные в канал
	for _, urlID := range urlIDs {
		inputCh <- dataForWorker{
			ID:     urlID,
			SignID: signID,
		}
	}
}

func getFilledChan(inputCh <-chan dataForWorker, size int) <-chan dataForWorker {
	log.Printf("Создаем канал с структурой dataForWorker и буфером %v", size)
	resultCh := make(chan dataForWorker, size)
	for i := 0; i < size; i++ {
		job, ok := <-inputCh
		log.Printf("Chan contains, %v", job)
		if !ok {
			log.Printf("Произошла ОШИБКА при чтении данных из канала, выход из цикла")
			break
		}
		resultCh <- job
	}
	log.Printf("Закрываем канал")
	close(resultCh)
	log.Printf("Возвращаем последние данные %v", resultCh)
	return resultCh
}

func (h *Handler) worker(inputCh <-chan dataForWorker) {
	log.Printf("Воркер запущен!")
	for {
		log.Printf("Получаем заполненный канал")
		filledChan := getFilledChan(inputCh, batchSize)
		log.Printf("генерируем батч")
		batches := map[uint32][]string{}
		log.Printf("Наполняем джобу")
		for job := range filledChan {
			batches[job.SignID] = append(batches[job.SignID], job.ID)
		}
		log.Printf("НОтправляем данные в БД!")
		for signID, ids := range batches {
			err := h.writeToDB(ids, signID)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	time.Sleep(time.Second) // нужно ожидать корректно
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
