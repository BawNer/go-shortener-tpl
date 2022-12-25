package database

import (
	"log"

	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type DB struct {
	repository *PgDB
	memory     *storage.LocalShortenData
}

func (d *DB) Init() error {
	d.memory.InputCh = make(chan storage.DataForWorker, 100)
	return nil
}

func New() (*DB, error) {
	db, err := NewConn()
	if err != nil {
		return nil, err
	}

	return &DB{repository: db}, nil
}

func (d *DB) SaveURL(id string, data *storage.LocalShortenData) error {
	return d.repository.Insert(data)
}

func (d *DB) GetURL(id string) (*storage.LocalShortenData, error) {
	return d.repository.SelectByID(id)
}

func (d *DB) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	return d.repository.SelectBySignID(signID)
}

func (d *DB) GetByField(field, val string) (*storage.LocalShortenData, error) {
	return d.repository.SelectByField(field, val)
}

func (d *DB) DeleteURL(id string, val bool, signID uint32) error {
	return d.repository.DeleteURL(id, val, signID)
}

func (d *DB) RunWorkers(count int) {
	for i := 0; i < count; i++ {
		d.memory.WG.Add(1)
		go func() {
			for {
				data, ok := <-d.memory.InputCh
				if !ok {
					log.Printf("Канал закрылся, завершаем работу")
					d.memory.WG.Done()
					return
				}
				log.Printf("Отправляем данные в БД!")
				err := d.DeleteURL(data.ID, true, data.SignID)
				if err != nil {
					log.Printf("Проблема в бд, %v", err)
				}
			}
		}()
	}
}

func (d *DB) AddJob(urlIDs []string, signID uint32) {
	go func() {
		for _, urlID := range urlIDs {
			d.memory.InputCh <- storage.DataForWorker{
				ID:     urlID,
				SignID: signID,
			}
		}
	}()
}

func (d *DB) Stop() {
	close(d.memory.InputCh)
}

func (d *DB) Wait() {
	d.memory.WG.Wait()
}
