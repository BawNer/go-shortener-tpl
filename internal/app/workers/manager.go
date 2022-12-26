package workers

import (
	"errors"
	"log"
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

var (
	ErrShuttingDown = errors.New("workers chain are closed")
)

type WorkerPool struct {
	wg      sync.WaitGroup
	once    sync.Once
	stop    chan struct{}
	inputCh chan DataForWorker
	storage storage.Storage
}

type DataForWorker struct {
	ID     string
	SignID uint32
}

func (w *WorkerPool) AddJob(id string, signID uint32) error {
	select {
	case <-w.stop:
		return ErrShuttingDown
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
				case <-w.stop:
					log.Println("Exiting")
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
	w.once.Do(func() {
		close(w.stop)
		close(w.inputCh)
	})
	w.wg.Wait()
}

func NewWorkerPool(repository storage.Storage) *WorkerPool {
	return &WorkerPool{
		wg:      sync.WaitGroup{},
		stop:    make(chan struct{}),
		inputCh: make(chan DataForWorker, 100),
		storage: repository,
	}
}
