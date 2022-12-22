package database

import (
	"log"

	"github.com/BawNer/go-shortener-tpl/internal/app/handlers"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
)

type DB struct {
	dsn        *PgDB
	repository *storage.Repository
}

func (d *DB) Init() error {
	return nil
}

func New() (*DB, error) {
	db, err := NewConn()
	if err != nil {
		return nil, err
	}

	return &DB{dsn: db}, nil
}

func (d *DB) SaveURL(id string, data *storage.LocalShortenData) error {
	return d.dsn.Insert(data)
}

func (d *DB) GetURL(id string) (*storage.LocalShortenData, error) {
	return d.dsn.SelectByID(id)
}

func (d *DB) GetAllURLsForSignID(signID uint32) ([]*storage.LocalShortenData, error) {
	return d.dsn.SelectBySignID(signID)
}

func (d *DB) GetByField(field, val string) (*storage.LocalShortenData, error) {
	return d.dsn.SelectByField(field, val)
}

func (d *DB) DeleteURL(id string, val bool, signID uint32) error {
	return d.dsn.DeleteURL(id, val, signID)
}

func (d *DB) Wait() {
	d.repository.WG.Wait()
}

func (d *DB) Stop() {
	d.repository.WG.Done()
}

func (d *DB) RunWorkers(count int) {
	d.repository = storage.NewRepository()
	for i := 0; i < count; i++ {
		d.repository.WG.Add(1)
		go func() {
			for {
				data, ok := <-d.repository.InputCh
				if !ok {
					log.Printf("Канал закрылся, завершаем работу")
					d.repository.WG.Done()
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

func (d *DB) PutJob(urlIDs []string, signID uint32) {
	for _, urlID := range urlIDs {
		d.repository.InputCh <- handlers.DataForWorker{
			ID:     urlID,
			SignID: signID,
		}
	}
}
