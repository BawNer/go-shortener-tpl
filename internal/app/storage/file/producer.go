package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
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
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *storage.LocalShortenData) error {
	p.encoder.SetEscapeHTML(false)
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
