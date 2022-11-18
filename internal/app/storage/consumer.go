package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type consumer struct {
	file    *os.File
	decoder *json.Decoder
	scanner *bufio.Scanner
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) ReadEvent() (*MyDB, error) {
	event := &MyDB{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *consumer) ReadEventAll() ([]MyDB, error) {
	var (
		event    MyDB
		eventAll []MyDB
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

func (c *consumer) Close() error {
	return c.file.Close()
}
