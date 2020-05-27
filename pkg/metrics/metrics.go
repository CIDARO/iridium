package metrics

import (
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v2"
)

func HandleMetrics(res *http.Response, cache *badger.DB) error {
	err := cache.Update(func(txn *badger.Txn) error {
		status := res.Status
		log.Printf("Status: %s", status)
		err := txn.Set([]byte("test"), []byte(status))
		return err
	})
	if err != nil {
		return err
	}
	err = cache.View(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte("test"))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			log.Printf("%s", val)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})
	return nil
}
