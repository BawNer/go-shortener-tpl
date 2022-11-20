package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
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
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteEvent(event *LocalShortenData) error {
	p.encoder.SetEscapeHTML(false)
	return p.encoder.Encode(&event)
}

func (p *producer) Close() error {
	return p.file.Close()
}
