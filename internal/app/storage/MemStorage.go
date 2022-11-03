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
	ID       int
	URL      string
	URLShort string
}

type MemStorage struct {
	sync.RWMutex
	Storage map[DBKey]MyDB
}

type MemStorageInterface interface {
	SaveDB(k DBKey, d MyDB) map[DBKey]MyDB
	FindByID(k DBKey) (MyDB, error)
}

func (m *MemStorage) SaveDB(k DBKey, d MyDB) MyDB {
	m.Lock()
	defer m.Unlock()
	if m.Storage == nil {
		m.Storage = map[DBKey]MyDB{}
	}
	m.Storage[k] = d
	return m.Storage[k]
}

func (m *MemStorage) FindByID(id DBKey) (MyDB, error) {
	m.RLock()
	defer m.RUnlock()
	result := MyDB{}
	err := ErrNotFound
	for idb := range m.Storage {
		if idb == id {
			return m.Storage[id], nil
		}
	}

	return result, err
}
