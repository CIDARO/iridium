package main

import (
	"flag"
)

func main() {
	configPath := flag.String("config", "/usr/local/iridium/configs/", "path to configuration file")
	flag.Parse()
}
