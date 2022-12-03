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

func NewConn() (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), app.Config.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	// create table
	query, err := db.Query(context.Background(),
		"CREATE TABLE IF NOT EXISTS shortened_urls (id varchar(20), url varchar(40)  PRIMARY KEY, signID bigint NOT NULL)",
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer query.Close()

	return db, err
}

func Insert(db *pgxpool.Pool, params *storage.LocalShortenData) error {
	query, err := db.Query(context.Background(),
		"INSERT INTO shortened_urls (id, url, signID) VALUES ($1, $2, $3)", params.ID, params.URL, params.SignID)
	if err != nil {
		return err
	}
	query.Close()
	if query.Err() != nil {
		// строка уже существует, необходимо вернуть ошибку и уже существующую строку
		return query.Err()
	}
	return nil
}

func SelectByID(db *pgxpool.Pool, id string) (*storage.LocalShortenData, error) {
	var (
		data storage.LocalShortenData
	)
	err := db.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE id=$1", id).Scan(&data.ID, &data.URL, &data.SignID)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func SelectByField(db *pgxpool.Pool, field string, val string) (*storage.LocalShortenData, error) {
	var (
		data storage.LocalShortenData
	)
	switch field {
	case "url":
		err := db.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE url=$1", val).Scan(&data.ID, &data.URL, &data.SignID)
		if err != nil {
			return nil, err
		}
	case "id":
		err := db.QueryRow(context.Background(), "SELECT * FROM shortened_urls WHERE id=$1", val).Scan(&data.ID, &data.URL, &data.SignID)
		if err != nil {
			return nil, err
		}
	default:
		return &data, nil
	}

	return &data, nil
}

func SelectBySignID(db *pgxpool.Pool, signID uint32) ([]*storage.LocalShortenData, error) {
	var (
		data []*storage.LocalShortenData
	)
	query, err := db.Query(context.Background(), "SELECT * FROM shortened_urls WHERE signID = $1", signID)
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
			ID:     value[0].(string),
			URL:    value[1].(string),
			SignID: uint32(value[2].(int64)),
		})
	}

	return data, nil
}
