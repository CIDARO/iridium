package db

import (
	"fmt"

	"github.com/CIDARO/iridium/internal/config"
	"github.com/bradfitz/gomemcache/memcache"
)

// Memcache is the memcache client for the metrics
var Memcache *memcache.Client

// InitDb creates a new Memcache client used in the StoreMetrics function
func InitDb() {
	Memcache = memcache.New(fmt.Sprintf("%s:%s", config.Config.Memcache.Host, config.Config.Memcache.Port))
}
