package workerpool

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

var (
	ErrRepositoryIsClosing = errors.New("workers chain are closed")
)

type WorkerPool struct {
	wg      sync.WaitGroup
	down    chan struct{}
	inputCh chan DataForWorker
	storage storage.Storage
}

type DataForWorker struct {
	ID     string
	SignID uint32
}

func (w *WorkerPool) AddJob(id string, signID uint32) error {
	select {
	case <-w.down:
		return ErrRepositoryIsClosing
	case w.inputCh <- DataForWorker{
		ID:     id,
		SignID: signID,
	}:
		return nil
	}
}

func (w *WorkerPool) RunWorkers(countWorkers int) {
	for i := 0; i < countWorkers; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-w.down:
					fmt.Println("Exiting")
					return
				case v, ok := <-w.inputCh:
					if !ok {
						return
					}
					// here send data to DB
					err := w.storage.DeleteURL(v.ID, true, v.SignID)
					if err != nil {
						log.Printf("Error wheh url removed %v", err)
					}
				}
			}
		}()
	}
}

func (w *WorkerPool) Stop() {
	once := sync.Once{}

	once.Do(func() {
		close(w.down)
		close(w.inputCh)
	})

	w.wg.Wait()
}

func NewWorkerPool(repository storage.Storage) *WorkerPool {
	return &WorkerPool{
		wg:      sync.WaitGroup{},
		down:    make(chan struct{}),
		inputCh: make(chan DataForWorker, 100),
		storage: repository,
	}
}
