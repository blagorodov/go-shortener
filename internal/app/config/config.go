package config

import "flag"

type Config struct {
	Server     string
	ResultHost string
}

var Options = new(Config)

func ParseFlags() {
	flag.StringVar(&Options.Server, "a", ":8080", "address and port to run server")
	flag.StringVar(&Options.ResultHost, "b", "http://localhost:8080", "result server name")
	flag.Parse()
}
