package config

import (
	"flag"
	"os"
	"testing"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

var Options = new(Config)

func init() {
	testing.Init()
	flag.StringVar(&Options.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.BaseURL, "b", "http://localhost:8080", "result server name")
	flag.Parse()

	if e := os.Getenv("SERVER_ADDRESS"); e != "" {
		Options.ServerAddress = e
	}
	if e := os.Getenv("BASE_URL"); e != "" {
		Options.BaseURL = e
	}
}
