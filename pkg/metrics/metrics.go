package metrics

import (
	"net/http"

	"github.com/dgraph-io/badger/v2"
)

type Metric struct {
	Count           int64
	Duration        int64
	AverageDuration int64
	Path            string
}

func HandleMetrics(res *http.Response, cache *badger.DB) error {

	return nil
}
