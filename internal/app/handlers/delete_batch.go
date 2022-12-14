package handlers

import (
	"io"
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
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
	batchSize = len(string(data))
	inputCh := make(chan dataForWorker, batchSize)
	for _, urlID := range data {
		urlIDs = append(urlIDs, string(urlID))
	}

	go h.worker(inputCh) // init go routine

	go putJobs(inputCh, urlIDs, signID)

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
		err := h.storage.DeleteURL(id, true, signID)
		if err != nil {
			return err
		}
	}

	return nil
}
