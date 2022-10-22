package app

import (
	"errors"
)

type DB struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	UrlShort string `json:"urlShort"`
}

var LocalDB []DB

func SaveDB(payload DB) (DB, []DB) {
	LocalDB = append(LocalDB, payload)
	return payload, LocalDB
}

func FindById(id int) (DB, error) {
	result := DB{}
	var err = errors.New("not found")
	for _, ldb := range LocalDB {
		if ldb.Id == id {
			result = ldb
			err = nil
		} else {
			err = errors.New("not found")
		}
	}
	return result, err
}
