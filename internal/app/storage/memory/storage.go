package memory

import (
	"sync"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Memory struct {
	mu      sync.RWMutex
	storage map[string]*storage.LocalShortenData
}

func New() (*Memory, error) {
	return &Memory{
		storage: map[string]*storage.LocalShortenData{},
	}, nil
}

func (m *Memory) Init() error {
	return nil
}

func (m *Memory) SaveURL(id string, data *storage.LocalShortenData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[id] = data
	return nil
}

func (m *Memory) GetURL(id string) (*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.storage[id]
	if !ok {
		return &storage.LocalShortenData{}, storage.ErrNotFound
	}

	return v, nil
}
