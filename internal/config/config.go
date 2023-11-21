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
	URLDBPath     string
	DBDataSource  string
}

var Options = loadConfig()

func loadConfig() Config {
	testing.Init()
	o := Config{}
	flag.StringVar(&o.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&o.BaseURL, "b", "http://localhost:8080", "result server name")
	flag.StringVar(&o.LogPath, "l", "shortener.log", "log file path")
	flag.StringVar(&o.URLDBPath, "f", "", "url database file path")
	flag.StringVar(&o.DBDataSource, "d", "", "database data source")
	flag.Parse()

	if e := os.Getenv("SERVER_ADDRESS"); e != "" {
		o.ServerAddress = e
	}
	if e := os.Getenv("BASE_URL"); e != "" {
		o.BaseURL = e
	}
	if e := os.Getenv("FILE_STORAGE_PATH"); e != "" {
		o.URLDBPath = e
	}
	if e := os.Getenv("DATABASE_DSN"); e != "" {
		o.DBDataSource = e
	}
	return o
}
