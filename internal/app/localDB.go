package app

import (
	"errors"
)

type DB struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	URLShort string `json:"urlShort"`
}

var LocalDB []DB

func SaveDB(payload DB) (DB, []DB) {
	LocalDB = append(LocalDB, payload)
	return payload, LocalDB
}

func FindByID(id int) (DB, error) {
	result := DB{}
	var err = errors.New("not found")
	for _, ldb := range LocalDB {
		if ldb.ID == id {
			result = ldb
			err = nil
		} else {
			err = errors.New("not found")
		}
	}
	return result, err
}
