package storage

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type DBKey string

type MyDB struct {
	ID       int
	URL      string
	URLShort string
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
	fmt.Println(d)
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
