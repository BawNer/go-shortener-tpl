package storage

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrNotAccepted = errors.New("not accepted")
)

type LocalShortenData struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SignID    uint32 `json:"signID"`
	IsDeleted bool   `json:"-"`
}

type Storage interface {
	SaveURL(id string, data *LocalShortenData) error
	GetURL(id string) (*LocalShortenData, error)
	GetAllURLsForSignID(signID uint32) ([]*LocalShortenData, error)
	GetByField(field, val string) (*LocalShortenData, error)
	DeleteURL(id string, val bool, signID uint32) error
	Init() error
}
