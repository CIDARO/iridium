package metrics

import (
	"log"
	"net/http"
)

func HandleMetrics(res *http.Response) error {
	log.Printf("%v", res)
	return nil
}
