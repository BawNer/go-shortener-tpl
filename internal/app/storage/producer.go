package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*producer, error) {
	if _, err := os.Stat(filepath.Dir(filename)); err != nil {
		os.MkdirAll(filepath.Dir(filename), 0770)
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *producer) Close() error {
	return p.file.Close()
}
