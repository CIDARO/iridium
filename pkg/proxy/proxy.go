package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/CIDARO/iridium/pkg/config"
	"github.com/CIDARO/iridium/pkg/metrics"
	"github.com/dgraph-io/badger/v2"
)

var (
	director func(*http.Request)
	cache    badger.DB
)

func HandleDirector(req *http.Request) {
	director(req)
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = req.URL.Host
}

func HandleModifyResponse(res *http.Response) error {
	err := metrics.HandleMetrics(res, &cache)

	if err != nil {
		return err
	}

	return nil
}

func NewReverseProxy(config config.Config, database badger.DB) (*httputil.ReverseProxy, error) {

	proxy := httputil.NewSingleHostReverseProxy(&config.Target)

	director = proxy.Director
	cache = database

	proxy.Director = HandleDirector
	proxy.ModifyResponse = HandleModifyResponse

	return proxy, nil
}
