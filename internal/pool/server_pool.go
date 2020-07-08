package pool

import (
	"log"
	"net"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/CIDARO/iridium/internal/backend"
	"github.com/CIDARO/iridium/internal/utils"
)

// ServerPool struct that holds the backends and the current backend index
type ServerPool struct {
	backends []*backend.Backend
	current  uint32
}

// Pool is the global ServerPool
var Pool ServerPool

// AddBackend adds a new backend to the server pool
func (s *ServerPool) AddBackend(backend *backend.Backend) {
	s.backends = append(s.backends, backend)
}

// NextBackendIndex increments atomically the counter and returns the next backend index
func (s *ServerPool) NextBackendIndex() int {
	return int(atomic.AddUint32(&s.current, uint32(1)) % uint32(len(s.backends)))
}

// MarkBackendStatus changes the status of a backend
func (s *ServerPool) MarkBackendStatus(backendURL *url.URL, alive bool) {
	for _, backend := range s.backends {
		if backend.URL.String() == backendURL.String() {
			backend.SetAlive(alive)
			break
		}
	}
}

// GetNextBackend returns the next active backend that will receive the connection
func (s *ServerPool) GetNextBackend() *backend.Backend {
	// loop all the backends
	next := s.NextBackendIndex()
	length := len(s.backends) + next // start from the next
	for i := next; i < length; i++ {
		index := i % len(s.backends) // take the index by modding the length
		if s.backends[index].IsAlive() {
			if i != next {
				atomic.StoreUint32(&s.current, uint32(index))
			}
			return s.backends[index]
		}
	}
	return nil
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		utils.Logger.Infof("[HealthCheck] %s [%s]\n", b.URL, status)
	}
}

// isBackendAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		utils.Logger.Infof("site unreachable: %v", err)
		return false
	}
	_ = conn.Close()
	return true
}

// HealthCheck runs a routine for check status of the backends every 2 mins
func HealthCheck() {
	t := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-t.C:
			log.Println("[HealthCheck] Starting..")
			Pool.HealthCheck()
			log.Println("[HealthCheck] Completed..")
		}
	}
}
