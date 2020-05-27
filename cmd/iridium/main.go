package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/CIDARO/iridium/pkg/config"
	"github.com/CIDARO/iridium/pkg/proxy"
	"github.com/CIDARO/iridium/pkg/validation"
)

func main() {
	configPath := flag.String("config", "/usr/local/iridium/configs/", "path to configuration file")
	target := flag.String("target", "http://127.0.0.1:8081", "target url")
	flag.Parse()

	validatedPath, err := validation.ValidatePath(*configPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	config, err := config.GetConfig(*validatedPath, *target)
	if err != nil {
		log.Fatalf("%v", err)
	}

	proxy, err := proxy.NewReverseProxy(*config)
	if err != nil {
		log.Fatalf("%v", err)
	}

	http.HandleFunc("/", proxy.ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, nil))
}
