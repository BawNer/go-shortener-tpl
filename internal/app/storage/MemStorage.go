package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type DBKey string

type MyDB struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type MemStorage struct {
	mu      sync.RWMutex
	Storage map[DBKey]MyDB
}

type MemStorageInterface interface {
	SaveDB(k DBKey, d MyDB) map[DBKey]MyDB
	FindByID(k DBKey) (MyDB, error)
}

func (m *MemStorage) SaveDB(k DBKey, d MyDB) MyDB {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Storage == nil {
		m.Storage = map[DBKey]MyDB{}
	}
	m.Storage[k] = d
	return m.Storage[k]
}

func (m *MemStorage) FindByID(id DBKey) (MyDB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.Storage[id]

	if !ok {
		return MyDB{}, ErrNotFound
	}

	return v, nil
}

func (m *MemStorage) LoadDataFromFile(fileName string) error {
	consumer, err := NewConsumer(fileName)
	if err != nil {
		return err
	}

	fileData, _ := consumer.ReadEventAll()

	for _, data := range fileData {
		m.SaveDB(DBKey(data.ID), MyDB{
			ID:  data.ID,
			URL: data.URL,
		})
	}

	return nil
}
