package storage

import "errors"

type DBKey string

type MyDB struct {
	ID       int
	URL      string
	URLShort string
}

type MemStorage struct {
	Storage map[DBKey]MyDB
}

type MemStorageInterface interface {
	SaveDB(k DBKey, d MyDB) map[DBKey]MyDB
	FindByID(k DBKey) (MyDB, error)
}

func (m *MemStorage) SaveDB(k DBKey, d MyDB) MyDB {
	if m.Storage == nil {
		m.Storage = map[DBKey]MyDB{}
	}
	m.Storage[k] = d
	return m.Storage[k]
}

func (m *MemStorage) FindByID(id DBKey) (MyDB, error) {
	result := MyDB{}
	err := errors.New("not found")
	for idb := range m.Storage {
		if idb == id {
			result = m.Storage[id]
			err = nil
			break
		}
	}

	return result, err
}
