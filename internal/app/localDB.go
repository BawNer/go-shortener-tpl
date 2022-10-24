package app

import (
	"errors"
)

type DB struct {
	ID       int
	URL      string
	URLShort string
}

var LocalDB []DB

func SaveDB(payload DB) (DB, []DB) {
	LocalDB = append(LocalDB, payload)
	return payload, LocalDB
}

func FindByID(id string) (DB, error) {
	result := DB{}
	err := errors.New("not found")
	for i := 0; i < len(LocalDB); i++ {
		if LocalDB[i].URLShort == id {
			result = LocalDB[i]
			err = nil
			break
		}
	}
	return result, err
}
