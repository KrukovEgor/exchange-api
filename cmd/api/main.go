package main

import (
	"flag"
	"log"

	"github.com/KrukovEgor/exchange-api/internal/config"
)

const defaultRawConfigPath = "./configs/backend.yml"

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", defaultRawConfigPath, "config file path")
	flag.Parse()

	_, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
}
