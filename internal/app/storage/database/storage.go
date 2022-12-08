package database

import (
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type DB struct {
	repository *PgDB
}

func (d *DB) Init() error {
	return nil
}

func New() (*DB, error) {
	db, err := NewConn()
	if err != nil {
		return nil, err
	}

	return &DB{repository: db}, nil
}

func (d *DB) SaveURL(id string, data *storage.LocalShortenData) error {
	return d.repository.Insert(data)
}

func (d *DB) GetURL(id string) (*storage.LocalShortenData, error) {
	return d.repository.SelectByID(id)
}

func (d *DB) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	return d.repository.SelectBySignID(signID)
}

func (d *DB) GetByField(field, val string) (*storage.LocalShortenData, error) {
	return d.repository.SelectByField(field, val)
}
