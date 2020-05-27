package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/CIDARO/iridium/pkg/config"
	"github.com/CIDARO/iridium/pkg/db"
	"github.com/CIDARO/iridium/pkg/proxy"
	"github.com/CIDARO/iridium/pkg/validation"
)

func main() {
	configPath := flag.String("config", "/usr/local/iridium/configs/", "path to configuration file")
	target := flag.String("target", "", "target url")
	flag.Parse()

	if *target == "" {
		log.Fatalf("missing target service")
	}

	validatedPath, err := validation.ValidatePath(*configPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	config, err := config.GetConfig(*validatedPath, *target)
	if err != nil {
		log.Fatalf("%v", err)
	}

	database, err := db.CreateDatabase(config.DbPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer database.Close()

	proxy, err := proxy.NewReverseProxy(*config, *database)
	if err != nil {
		log.Fatalf("%v", err)
	}

	http.HandleFunc("/", proxy.ServeHTTP)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, nil))
}
