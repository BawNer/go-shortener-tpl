package file

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
	scanner *bufio.Scanner
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*storage.LocalShortenData, error) {
	event := &storage.LocalShortenData{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Consumer) ReadEventAll() ([]storage.LocalShortenData, error) {
	var (
		event    storage.LocalShortenData
		eventAll []storage.LocalShortenData
	)
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		if err := json.Unmarshal(data, &event); err != nil {
			log.Fatal(err)
		}
		eventAll = append(eventAll, event)
	}

	return eventAll, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
