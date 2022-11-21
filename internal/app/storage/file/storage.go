package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type File struct {
	file     *os.File
	encoder  *json.Encoder
	consumer *Consumer
}

func New(fileName string) (*File, error) {
	if _, err := os.Stat(filepath.Dir(fileName)); err != nil {
		err := os.MkdirAll(filepath.Dir(fileName), 0770)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	consumer, err := NewConsumer(fileName)
	if err != nil {
		return nil, err
	}
	return &File{
		file:     file,
		encoder:  json.NewEncoder(file),
		consumer: consumer,
	}, nil
}

func (p *File) SaveURL(id string, data *storage.LocalShortenData) error {
	p.encoder.SetEscapeHTML(false)
	return p.encoder.Encode(&data)
}

func (p *File) GetURL(id string) (*storage.LocalShortenData, error) {
	events, err := p.consumer.ReadEventAll()
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.ID == id {
			return &storage.LocalShortenData{
				ID:  event.ID,
				URL: event.URL,
			}, nil
		}
	}

	return &storage.LocalShortenData{}, nil
}

func (p *File) Close() error {
	return p.file.Close()
}
