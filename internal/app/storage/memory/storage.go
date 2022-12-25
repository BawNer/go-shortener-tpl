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

func (m *Memory) GetByField(field, val string) (*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch field {
	case "id":
		v, ok := m.storage[val]
		if !ok {
			return &storage.LocalShortenData{}, storage.ErrNotFound
		}
		return v, nil
	case "url":
		for _, v := range m.storage {
			if v.URL == val {
				return v, nil
			}
		}
		return &storage.LocalShortenData{}, storage.ErrNotFound
	default:
		return &storage.LocalShortenData{}, storage.ErrNotFound
	}
}

func (m *Memory) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var (
		urls []*storage.LocalShortenData
	)
	for _, item := range m.storage {
		if item.SignID == signID {
			urls = append(urls, item)
		}
	}

	if len(urls) < 1 {
		return nil, storage.ErrNotFound
	}

	return urls, nil
}

func (m *Memory) DeleteURL(id string, val bool, signID uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.storage[id]
	if !ok {
		return storage.ErrNotFound
	}

	if v.SignID != signID {
		return storage.ErrNotAccepted
	}
	m.storage[id].IsDeleted = val
	return nil
}
