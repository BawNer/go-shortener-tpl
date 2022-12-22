package storage

import (
	"errors"
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrNotAccepted = errors.New("not accepted")
)

type LocalShortenData struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SignID    uint32 `json:"signID"`
	IsDeleted bool   `json:"-"`
}

type Repository struct {
	InputCh chan handlers.DataForWorker
	WG      sync.WaitGroup
}

type Storage interface {
	SaveURL(id string, data *LocalShortenData) error
	GetURL(id string) (*LocalShortenData, error)
	GetAllURLsForSignID(signID uint32) ([]*LocalShortenData, error)
	GetByField(field, val string) (*LocalShortenData, error)
	DeleteURL(id string, val bool, signID uint32) error
	Init() error
	RunWorkers(count int)
	Wait()
	Stop()
	PutJob(urlIDs []string, signID uint32)
}

func NewRepository() *Repository {
	return &Repository{
		InputCh: make(chan handlers.DataForWorker, 100),
		WG:      sync.WaitGroup{},
	}
}
