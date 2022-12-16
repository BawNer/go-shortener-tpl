package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BawNer/go-shortener-tpl/internal/app"
	"github.com/BawNer/go-shortener-tpl/internal/app/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDB struct {
	pool *pgxpool.Pool
}

func NewConn() (*PgDB, error) {
	db, err := pgxpool.New(context.Background(), app.Config.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return &PgDB{}, err
	}
	// create table
	query, err := db.Query(context.Background(),
		"CREATE TABLE IF NOT EXISTS shortened_urls (id varchar(20), url varchar(255)  PRIMARY KEY,"+
			" signID bigint NOT NULL, isDeleted BOOLEAN NOT NULL DEFAULT FALSE)",
	)
	if err != nil {
		return &PgDB{}, err
	}
	defer query.Close()

	return &PgDB{pool: db}, err
}

func (d *PgDB) Insert(params *storage.LocalShortenData) error {
	log.Printf("Отправляем данные в бд %v", params)
	query, err := d.pool.Query(context.Background(),
		"INSERT INTO shortened_urls (id, url, signID, isDeleted) VALUES ($1, $2, $3, $4)", params.ID, params.URL, params.SignID, params.IsDeleted)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Данные отпавлены в бд %v", params)

	log.Println("Закрываем запрос")
	query.Close()
	if query.Err() != nil {
		log.Println(err)
		// строка уже существует, необходимо вернуть ошибку и уже существующую строку
		return query.Err()
	}
	log.Println("Соединение закрыто!")
	return nil
}

func (d *PgDB) SelectByID(id string) (*storage.LocalShortenData, error) {
	var (
		data storage.LocalShortenData
	)
	err := d.pool.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE id=$1", id).Scan(&data.ID, &data.URL, &data.SignID, &data.IsDeleted)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (d *PgDB) SelectByField(field string, val string) (*storage.LocalShortenData, error) {
	var (
		data storage.LocalShortenData
	)
	switch field {
	case "url":
		err := d.pool.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE url=$1", val).Scan(&data.ID, &data.URL, &data.SignID, &data.IsDeleted)
		if err != nil {
			return nil, err
		}
	case "id":
		err := d.pool.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE id=$1", val).Scan(&data.ID, &data.URL, &data.SignID, &data.IsDeleted)
		if err != nil {
			return nil, err
		}
	default:
		return &data, nil
	}

	return &data, nil
}

func (d *PgDB) SelectBySignID(signID uint32) ([]*storage.LocalShortenData, error) {
	var (
		data []*storage.LocalShortenData
	)
	query, err := d.pool.Query(context.Background(), "SELECT * FROM shortened_urls WHERE signID = $1", signID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	for query.Next() {
		value, err := query.Values()
		if err != nil {
			return nil, err
		}

		data = append(data, &storage.LocalShortenData{
			ID:        value[0].(string),
			URL:       value[1].(string),
			SignID:    uint32(value[2].(int64)),
			IsDeleted: value[3].(bool),
		})
	}

	return data, nil
}

func (d *PgDB) DeleteURL(id string, value bool, signID uint32) error {
	query, err := d.pool.Query(context.Background(), "UPDATE shortened_urls SET isDeleted=$1 WHERE id = $2 AND signID = $3",
		value, id, signID)
	if err != nil {
		log.Println(err)
		return err
	}
	query.Close()
	if query.Err() != nil {
		log.Println(err)
		return query.Err()
	}

	return nil
}
