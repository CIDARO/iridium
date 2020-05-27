package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/CIDARO/iridium/pkg/config"
	"github.com/CIDARO/iridium/pkg/metrics"
)

var director func(*http.Request)

func HandleDirector(req *http.Request) {
	director(req)
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = req.URL.Host
}

func HandleModifyResponse(res *http.Response) error {
	err := metrics.HandleMetrics(res)

	if err != nil {
		return err
	}

	return nil
}

func NewReverseProxy(config config.Config) (*httputil.ReverseProxy, error) {

	proxy := httputil.NewSingleHostReverseProxy(&config.Target)

	director = proxy.Director
	proxy.Director = HandleDirector
	proxy.ModifyResponse = HandleModifyResponse

	return proxy, nil
}
