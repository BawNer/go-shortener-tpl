package storage

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type LocalShortenData struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Storage interface {
	SaveURL(id string, data *LocalShortenData) error
	GetURL(id string) (*LocalShortenData, error)
	Init() error
}
