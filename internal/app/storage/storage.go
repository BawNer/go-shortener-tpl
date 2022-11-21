package storage

import (
	"errors"
	"log"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	InMemory *InMemory
	InFile   *InFile
}

func NewMemory(filepath string) *Storage {
	if filepath == "" {
		return &Storage{
			InMemory: NewLocalStorage(),
			InFile:   nil,
		}
	}

	fileStorage, err := NewFileStorage(filepath)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &Storage{
		InMemory: NewLocalStorage(),
		InFile:   fileStorage,
	}
}
