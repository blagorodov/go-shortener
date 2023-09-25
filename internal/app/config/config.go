package config

import (
	"flag"
	"os"
	"testing"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	LogPath       string
}

var Options = loadConfig()

func loadConfig() Config {
	testing.Init()
	o := Config{}
	flag.StringVar(&o.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&o.BaseURL, "b", "http://localhost:8080", "result server name")
	flag.StringVar(&o.LogPath, "l", "shortener.log", "log file path")
	flag.Parse()

	if e := os.Getenv("SERVER_ADDRESS"); e != "" {
		o.ServerAddress = e
	}
	if e := os.Getenv("BASE_URL"); e != "" {
		o.BaseURL = e
	}
	return o
}
