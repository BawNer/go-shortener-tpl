package storage

import "sync"

type DBKey string

type LocalShortenData struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type InMemory struct {
	mu      sync.RWMutex
	storage map[DBKey]LocalShortenData
}

type MemStorageInterface interface {
	SaveDB(k DBKey, d LocalShortenData) map[DBKey]LocalShortenData
	FindByID(k DBKey) (LocalShortenData, error)
	NewLocalStorage() *InMemory
}

func NewLocalStorage() *InMemory {
	return &InMemory{
		storage: map[DBKey]LocalShortenData{},
	}
}

func (m *InMemory) Save(k DBKey, d LocalShortenData) LocalShortenData {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[k] = d
	return m.storage[k]
}

func (m *InMemory) FindByID(id DBKey) (LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.storage[id]

	if !ok {
		return LocalShortenData{}, ErrNotFound
	}

	return v, nil
}

func (m *InMemory) LoadDataFromFile(fileName string) error {
	consumer, err := NewConsumer(fileName)
	if err != nil {
		return err
	}

	fileData, err := consumer.ReadEventAll()

	if err != nil {
		return err
	}

	for _, data := range fileData {
		m.Save(DBKey(data.ID), LocalShortenData{
			ID:  data.ID,
			URL: data.URL,
		})
	}

	return nil
}
