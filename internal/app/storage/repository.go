package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type DBKey string

type LocalShortenData struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Repository struct {
	mu      sync.RWMutex
	Storage map[DBKey]LocalShortenData
}

type MemStorageInterface interface {
	SaveDB(k DBKey, d LocalShortenData) map[DBKey]LocalShortenData
	FindByID(k DBKey) (LocalShortenData, error)
}

func (r *Repository) Save(k DBKey, d LocalShortenData) LocalShortenData {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Storage == nil {
		r.Storage = map[DBKey]LocalShortenData{}
	}
	r.Storage[k] = d
	return r.Storage[k]
}

func (r *Repository) FindByID(id DBKey) (LocalShortenData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.Storage[id]

	if !ok {
		return LocalShortenData{}, ErrNotFound
	}

	return v, nil
}

func (r *Repository) LoadDataFromFile(fileName string) error {
	consumer, err := NewConsumer(fileName)
	if err != nil {
		return err
	}

	fileData, err := consumer.ReadEventAll()

	if err != nil {
		return err
	}

	for _, data := range fileData {
		r.Save(DBKey(data.ID), LocalShortenData{
			ID:  data.ID,
			URL: data.URL,
		})
	}

	return nil
}
