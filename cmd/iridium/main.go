package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/CIDARO/iridium/internal/backend"
	"github.com/CIDARO/iridium/internal/config"
	"github.com/CIDARO/iridium/internal/db"
	"github.com/CIDARO/iridium/internal/pool"
	"github.com/CIDARO/iridium/internal/server"
	"github.com/CIDARO/iridium/internal/utils"
)

func main() {
	utils.SetupLogger()
	// Parse the config path or use the default value
	configPath := flag.String("config", "/usr/local/iridium/configs/", "path to config file")
	flag.Parse()
	// Validates the path
	validatedPath, err := config.ValidatePath(*configPath)
	if err != nil {
		utils.Logger.Fatalf("failed to validate path: %v", err)
	}
	// Retrieve the config from the config path
	config.GetConfig(*validatedPath)
	// Initialize the Memcache DB
	db.InitDb()

	// parse servers
	for _, b := range config.Config.Backends {
		serverURL, err := url.Parse(b)
		if err != nil {
			utils.Logger.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverURL)

		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			utils.Logger.Infof("[%s] %s", serverURL.Host, e.Error())
			retries := server.GetRetryFromContext(request)
			if retries < config.Config.MaxRetries {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), server.Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			pool.Pool.MarkBackendStatus(serverURL, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := server.GetAttemptsFromContext(request)
			utils.Logger.Infof("%s(%s) attempting retry %d", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), server.Attempts, attempts+1)
			server.LoadBalancer(writer, request.WithContext(ctx))
		}

		pool.Pool.AddBackend(&backend.Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxy,
		})
		utils.Logger.Infof("added backend [%v]", serverURL)
	}

	serverPort, err := strconv.Atoi(config.Config.Server.Port)
	if err != nil {
		utils.Logger.Errorf("error while parsing server port: %v", err)
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: http.HandlerFunc(server.LoadBalancer),
	}

	// start health checking
	go pool.HealthCheck()

	utils.Logger.Printf("load balancer started on port %d", serverPort)
	if err := server.ListenAndServe(); err != nil {
		utils.Logger.Fatal(err)
	}
}
