package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
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
	if err := json.NewDecoder(r.Body).Decode(&urlIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sign, err := r.Cookie("sign")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// work with cookie
	signID, err = storage.CompareSign(sign.Value, app.Config.Secret)
	if err != nil {
		log.Println(err)
		signID = 0
	}
	// usr not auth
	if signID == 0 {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// main
	batchSize = len(urlIDs) + 2
	inputCh := make(chan dataForWorker, batchSize)

	go putJobs(inputCh, urlIDs, signID)

	go h.worker(inputCh) // init go routine

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
	resultCh := make(chan dataForWorker, size)
	for i := 0; i < size; i++ {
		job, ok := <-inputCh
		log.Printf("Chan contains, %v", job)
		if !ok {
			break
		}
		resultCh <- job
	}
	close(resultCh)
	return resultCh
}

func (h *Handler) worker(inputCh <-chan dataForWorker) {
	for {
		filledChan := getFilledChan(inputCh, batchSize)
		batches := map[uint32][]string{}
		for job := range filledChan {
			batches[job.SignID] = append(batches[job.SignID], job.ID)
		}
		for signID, ids := range batches {
			err := h.writeToDB(ids, signID)
			if err != nil {
				log.Println(err.Error())
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
