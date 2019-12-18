package main

import "github.com/alexflint/go-arg"

type config struct {
	ServerAddress string `arg:"--server-address,env:SERVER_ADDRESS,help:Server-address network-interface to bind on e.g.: '127.0.0.1:8080' default ':80'"`
}

func newConfig() (config, error) {
	cfg := config{
		ServerAddress: ":80",
	}
	err := arg.Parse(&cfg)

	return cfg, err
}
