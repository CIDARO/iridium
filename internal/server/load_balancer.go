package server

import (
	"net/http"

	"github.com/CIDARO/iridium/internal/config"
	"github.com/CIDARO/iridium/internal/metrics"
	"github.com/CIDARO/iridium/internal/pool"
)

const (
	Attempts int = iota
	Retry
)

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

func LoadBalancer(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > config.Config.MaxAttempts {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	peer := pool.Pool.GetNextBackend()
	if peer != nil {
		if config.Config.Metrics {
			metrics.StoreMetrics(r)
		}
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}
