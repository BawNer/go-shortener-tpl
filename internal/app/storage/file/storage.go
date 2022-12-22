package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage/memory"
)

type File struct {
	file     *os.File
	encoder  *json.Encoder
	consumer *Consumer
	memory   *memory.Memory
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

	memoryStorage, _ := memory.New()
	return &File{
		file:     file,
		encoder:  json.NewEncoder(file),
		consumer: consumer,
		memory:   memoryStorage,
	}, nil
}

func (f *File) Init() error {
	events, err := f.consumer.ReadEventAll()
	if err != nil {
		return err
	}

	for _, event := range events {
		err := f.memory.SaveURL(event.ID, &event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) SaveURL(id string, data *storage.LocalShortenData) error {
	f.encoder.SetEscapeHTML(false)
	err := f.encoder.Encode(&data)
	if err != nil {
		return err
	}

	return f.memory.SaveURL(data.ID, data)
}

func (f *File) GetURL(id string) (*storage.LocalShortenData, error) {
	return f.memory.GetURL(id)
}

func (f *File) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	return f.memory.GetAllURLsForSignID(signID)
}

func (f *File) GetByField(filed, val string) (*storage.LocalShortenData, error) {
	return f.memory.GetByField(filed, val)
}

func (f *File) DeleteURL(id string, val bool, signID uint32) error {
	return f.memory.DeleteURL(id, val, signID)
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Wait() {
	f.memory.Wait()
}

func (f *File) Stop() {
	f.memory.Stop()
}

func (f *File) RunWorkers(count int) {
	f.memory.RunWorkers(count)
}

func (f *File) PutJob(urlIDs []string, signID uint32) {
	f.memory.PutJob(urlIDs, signID)
}
