package database

import (
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	repository *pgxpool.Pool
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
	return Insert(d.repository, data)
}

func (d *DB) GetURL(id string) (*storage.LocalShortenData, error) {
	return SelectByID(d.repository, id)
}

func (d *DB) GetAllURL(signID uint32) ([]*storage.LocalShortenData, error) {
	return SelectBySignID(d.repository, signID)
}

func (d *DB) GetByField(field, val string) (*storage.LocalShortenData, error) {
	return SelectByField(d.repository, field, val)
}
