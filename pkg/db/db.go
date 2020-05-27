package db

import (
	"github.com/CIDARO/iridium/pkg/validation"
	badger "github.com/dgraph-io/badger/v2"
)

func CreateDatabase(databasePath string) (*badger.DB, error) {

	_, err := validation.ValidatePath(databasePath)
	if err != nil {
		return nil, err
	}

	database, err := badger.Open(badger.DefaultOptions(databasePath))
	if err != nil {
		return nil, err
	}

	return database, nil
}
